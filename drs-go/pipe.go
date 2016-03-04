package drs

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

var connections = int64(0)
var pending = int64(0)
var total = int64(0)

func init() {
	http.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
		response(w, 200, map[string]interface{}{
			"connections": connections,
			"commands": map[string]interface{}{
				"total": total,
			},
		})
	})
}

func response(w http.ResponseWriter, status int, input interface{}) {
	data, _ := json.Marshal(input)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

type Pipe struct {
	Protocol  protocol.Protocol
	Router    RouterHandler
	transport Transport
	Events    *Events
	*Processor
	mutex    sync.Mutex
	outbound cmap.ConcurrentMap
}

type Events struct {
	Connect    func(conn *Connection) error
	Disconnect func(conn *Connection)
}

func New(transport Transport) *Pipe {
	return &Pipe{
		Processor: NewProcessor(),
		Protocol:  protocol.JSON,
		Events:    new(Events),
		outbound:  cmap.New(),
		transport: transport,
		mutex:     sync.Mutex{},
	}
}

func (this *Pipe) Send(cmd *Command) (interface{}, error) {
	for {
		conn, err := this.route(cmd.Action)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		result, err := conn.Send(cmd)
		if err != nil {
			if _, ok := err.(*DRSException); ok {
				time.Sleep(1 * time.Second)
				continue
			}
			if casted, ok := err.(*DRSError); ok {
				return nil, casted
			}
		}
		return result, err
	}
}

func (this *Pipe) route(action string) (*Connection, error) {
	host, err := this.Router(action)
	if err != nil {
		return nil, err
	}
	{
		match, ok := this.outbound.Get(host)
		if ok {
			return match.(*Connection), nil
		}
	}
	{
		this.mutex.Lock()
		defer this.mutex.Unlock()
		match, ok := this.outbound.Get(host)
		if ok {
			return match.(*Connection), nil
		}

		conn, err := Dial(this.transport, this.Protocol, host)
		if err != nil {
			return nil, err
		}
		conn.Redirect = this.Processor
		this.outbound.Set(host, conn)
		go func() {
			conn.Read()
			this.outbound.Remove(host)
		}()
		return conn, nil
	}
}

func (this *Pipe) Listen() error {
	this.On("drs.ping", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
		return time.Now().UnixNano() / 1000, nil
	})
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		atomic.AddInt64(&connections, 1)
		defer atomic.AddInt64(&connections, -1)

		conn := NewConnection(rw, this.Protocol)
		if this.Events.Connect != nil {
			if err := this.Events.Connect(conn); err != nil {
				return
			}
		}
		conn.Redirect = this.Processor
		conn.Read()
		if this.Events.Disconnect != nil {
			this.Events.Disconnect(conn)
		}
	})
}

func (this *Pipe) Close() {
	for value := range this.outbound.Iter() {
		value.Val.(*Connection).Close()
	}
	this.outbound = cmap.New()
}
