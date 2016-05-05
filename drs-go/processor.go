package drs

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/dynamic"
	"github.com/streamrail/concurrent-map"
)

type Message struct {
	Command    *Command
	Connection *Connection
	Context    map[string]interface{}
}

type Processor struct {
	handlers map[string][]func(*Message) (interface{}, error)
	pending  cmap.ConcurrentMap
	stats    cmap.ConcurrentMap
}

type Stats struct {
	Errors     int64 `json:"errors"`
	Exceptions int64 `json:"exceptions"`
	Success    int64 `json:"success"`
}

func newProcessor() *Processor {
	return &Processor{
		handlers: map[string][]func(*Message) (interface{}, error){},
		pending:  cmap.New(),
		stats:    cmap.New(),
	}
}

func (this *Processor) On(action string, handlers ...func(*Message) (interface{}, error)) error {
	exists, ok := this.handlers[action]
	if ok {
		exists = append(exists, handlers...)
		return nil
	}
	this.handlers[action] = handlers
	return nil
}

func (this *Processor) wait(cmd *Command, cb func()) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	wait := make(chan *Command, 1)
	this.pending.Set(cmd.Key, wait)
	cb()
	response := <-wait
	if response.Action == ERROR {
		return nil, &DRSError{
			Message: fmt.Sprint(response.Body),
		}
	}
	if response.Action == EXCEPTION {
		return nil, &DRSException{
			Message: fmt.Sprint(response.Body),
		}
	}
	return response.Body, nil
}

func (this *Processor) Process(cmd *Command, conn *Connection) (interface{}, error) {
	if cmd.Action == RESPONSE || cmd.Action == ERROR || cmd.Action == EXCEPTION {
		waiting, ok := this.pending.Get(cmd.Key)
		if ok {
			waiting.(chan *Command) <- cmd
			this.pending.Remove(cmd.Key)
		}
		return nil, nil
	}

	// atomic.AddInt64(&this.total, 1)
	return this.trigger(cmd, conn)
}

func (this *Processor) trigger(cmd *Command, conn *Connection) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			log.Println(string(debug.Stack()))
			err = r.(error)
		}
	}()
	handlers, ok := this.handlers[cmd.Action]
	if !ok {
		handlers, ok = this.handlers["*"]
		if !ok {
			return nil, Error("No handlers for this action")
		}
	}
	msg := &Message{
		Context:    dynamic.Empty(),
		Command:    cmd,
		Connection: conn,
	}
	for _, h := range handlers {
		result, err = h(msg)
		if err != nil {
			return nil, err
		}
	}
	return result, err
}

func (this *Processor) clear() {
	for kv := range this.pending.Iter() {
		kv.Val.(chan *Command) <- &Command{
			Action: "drs.exception",
			Body: map[string]interface{}{
				"message": "Connection closed",
			},
		}
		this.pending.Remove(kv.Key)
	}
}
