package drs

import (
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	_ "net/http/pprof"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/dynamic"
	"github.com/streamrail/concurrent-map"
)

const (
	ERROR     = "drs.error"
	RESPONSE  = "drs.response"
	EXCEPTION = "drs.exception"
)

type Connection struct {
	*Processor
	Cache    cmap.ConcurrentMap
	protocol protocol.Protocol
	closed   bool
	raw      io.ReadWriteCloser
	Raw      io.ReadWriteCloser
	stream   *protocol.Stream
	mutex    sync.Mutex
	pending  cmap.ConcurrentMap
	write    chan *Command
}

func NewConnection(proto protocol.Protocol) *Connection {
	return &Connection{
		Processor: NewProcessor(),
		Cache:     cmap.New(),
		// mutex:     sync.Mutex{},
		pending:  cmap.New(),
		protocol: proto,
	}
}

func Dial(proto protocol.Protocol, transport Transport, host string) *Connection {
	this := NewConnection(proto)
	go func() {
		for {
			this.mutex.Lock()
			if this.closed {
				return
			}
			raw, err := transport.Connect(host)
			if err == nil {
				this.accept(raw)
				this.mutex.Unlock()
				this.Read()
			} else {
				this.mutex.Unlock()
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return this
}

func (this *Connection) accept(raw io.ReadWriteCloser) {
	this.raw = raw
	this.Raw = raw
	this.stream = this.protocol(raw)
}

func (this *Connection) clear() {
	for value := range this.pending.Iter() {
		value.Val.(chan *Command) <- &Command{
			Key:    value.Key,
			Action: EXCEPTION,
			Body:   "Disconnected",
		}
	}
	this.pending = cmap.New()
}

func (this *Connection) Send(cmd *Command) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	wait := make(chan *Command)
	this.pending.Set(cmd.Key, wait)
	err := this.stream.Encode(cmd)
	if err != nil {
		this.pending.Remove(cmd.Key)
		return nil, err
	}
	response := <-wait
	if response.Action == ERROR {
		return nil, &DRSError{
			Message: dynamic.String(response.Body.(map[string]interface{}), "message"),
		}
	}
	if response.Action == EXCEPTION {
		return nil, &DRSException{
			Message: fmt.Sprint(response.Body),
		}
	}
	return response.Body, nil
}

func (this *Connection) Fire(cmd *Command) error {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	return this.stream.Encode(cmd)
}

func (this *Connection) Read() {
	for {
		cmd := new(Command)
		err := this.stream.Decode(cmd)
		if err != nil {
			log.Println("Connection closing because", err)
			break
		}
		go this.process(cmd, this)
	}
	this.clear()
}

func (this *Connection) Close() {
	this.mutex.Lock()
	this.closed = true
	if this.raw != nil {
		this.raw.Close()
	}
	this.mutex.Unlock()
}
