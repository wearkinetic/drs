package drs

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/streamrail/concurrent-map"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func response(w http.ResponseWriter, status int, input interface{}) {
	data, _ := json.Marshal(input)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

type Pipe struct {
	*Server
	Router    RouterHandler
	transport Transport
	mutex     sync.Mutex
	outbound  cmap.ConcurrentMap
}

func New(transport Transport) *Pipe {
	result := &Pipe{
		Server:    NewServer(transport),
		outbound:  cmap.New(),
		transport: transport,
		mutex:     sync.Mutex{},
	}
	http.HandleFunc("/stats", func(w http.ResponseWriter, req *http.Request) {
		functions := []string{}
		for key, _ := range result.Processor.handlers {
			functions = append(functions, key)
		}
		response(w, 200, map[string]interface{}{
			"connections": map[string]interface{}{
				"inbound":  result.inbound.Count(),
				"outbound": result.outbound.Count(),
			},
			"commands": map[string]interface{}{
				"total":      result.total,
				"exceptions": result.exceptions,
				"errors":     result.errors,
			},
			"functions": functions,
		})
	})
	return result
}

func (this *Pipe) Request(cmd *Command) (interface{}, error) {
	for {
		conn, err := this.route(cmd.Action)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		result, err := conn.Request(cmd)
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

func (this *Pipe) Broadcast(cmd *Command) error {
	return errors.New("Not implemented")
}

func (this *Pipe) route(action string) (*Connection, error) {
	all, err := this.Router(action)
	if err != nil {
		return nil, err
	}
	host := all[rand.Intn(len(all))]
	// TODO: Find out whether double checked locking is bad
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
		conn := NewConnection(this.Protocol)
		conn.Redirect = this.Processor
		go func() {
			conn.Dial(this.transport, host, false)
			this.outbound.Remove(host)
		}()
		this.outbound.Set(host, conn)
		return conn, nil
	}
}

func (this *Pipe) Close() {
	for value := range this.outbound.Iter() {
		value.Val.(*Connection).Close()
	}
	this.outbound = cmap.New()
	this.Server.Close()
}
