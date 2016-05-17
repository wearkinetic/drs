package drs

import (
	"github.com/ironbay/dynamic"
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
}

func NewProcessor() *Processor {
	return &Processor{
		handlers: make(map[string][]func(*Message) (interface{}, error)),
		pending:  cmap.New(),
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

func (this *Processor) Process(conn *Connection, cmd *Command) {
	if cmd.Action == ERROR || cmd.Action == EXCEPTION || cmd.Action == RESPONSE {
		match, ok := this.pending.Get(cmd.Key)
		if !ok {
			return
		}
		this.pending.Remove(cmd.Key)
		match.(chan *Command) <- cmd
	}
	resp, err := this.Invoke(conn, cmd)
	conn.respond(cmd.Key, resp, err)
}

func (this *Processor) Invoke(conn *Connection, cmd *Command) (interface{}, error) {
	message := &Message{
		Conn:    conn,
		Command: cmd,
		Context: dynamic.Empty(),
	}
	handlers := this.handlers[cmd.Action]
	if handlers == nil {
		return nil, Error("No handlers for " + cmd.Action)
	}
	var result interface{}
	var err error
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
