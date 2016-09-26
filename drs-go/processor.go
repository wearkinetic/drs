package drs

import (
	"fmt"
	"log"
	"runtime/debug"

	"github.com/ironbay/dynamic"
	"github.com/ironbay/go-util/console"
	"github.com/streamrail/concurrent-map"
)

type Message struct {
	Conn    *Connection
	Command *Command
	Context map[string]interface{}
}

type Processor struct {
	handlers map[string][]func(*Message) (interface{}, error)
	pending  cmap.ConcurrentMap
	stats    cmap.ConcurrentMap

	Before func(msg *Message)
	After  func(msg *Message, result interface{}, err error)
}

type Stats struct {
	Errors     int64 `json:"errors"`
	Exceptions int64 `json:"exceptions"`
	Success    int64 `json:"success"`
}

func NewProcessor() *Processor {
	return &Processor{
		handlers: make(map[string][]func(*Message) (interface{}, error)),
		pending:  cmap.New(),
		stats:    cmap.New(),
		Before:   func(msg *Message) {},
		After:    func(msg *Message, result interface{}, err error) {},
	}
}

func (this *Processor) Enqueue(key string) chan *Command {
	block := make(chan *Command)
	this.pending.Set(key, block)
	return block
}

func (this *Processor) On(action string, cb ...func(*Message) (interface{}, error)) {
	this.handlers[action] = cb
}

func (this *Processor) Off(action string) {
	delete(this.handlers, action)
}

func (this *Processor) Process(conn *Connection, cmd *Command) {
	if cmd.Action == ERROR || cmd.Action == EXCEPTION || cmd.Action == RESPONSE {
		match, ok := this.pending.Get(cmd.Key)
		if !ok {
			return
		}
		this.pending.Remove(cmd.Key)
		match.(chan *Command) <- cmd
		return
	}
	// if this.parent != nil {
	// 	this.parent.Process(conn, cmd)
	// 	return
	// }
	resp, err := this.Invoke(conn, cmd)
	conn.respond(cmd.Key, resp, err)
}

func (this *Processor) Invoke(conn *Connection, cmd *Command) (result interface{}, err error) {
	message := &Message{
		Conn:    conn,
		Command: cmd,
		Context: dynamic.Empty(),
	}
	defer func() {
		this.After(message, result, err)
	}()
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			console.JSON(cmd)
			log.Println(string(debug.Stack()))
			var ok bool
			if err, ok = r.(error); !ok {
				err = Exception(fmt.Sprint(r))
			}
		}
	}()
	this.Before(message)
	handlers := this.handlers[cmd.Action]
	if handlers == nil {
		return nil, Error("No handlers for " + cmd.Action)
	}
	for _, cb := range handlers {
		result, err = cb(message)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (this *Processor) Clear() {
	for kv := range this.pending.Iter() {
		kv.Val.(chan *Command) <- &Command{
			Key:    kv.Key,
			Action: ERROR,
			Body: dynamic.Build(
				"message", "Connection closed",
			),
		}
		this.pending.Remove(kv.Key)
	}
}
