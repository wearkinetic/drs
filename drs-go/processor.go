package drs

import (
	"fmt"
	"log"
	"runtime/debug"
	"sync/atomic"

	"github.com/ironbay/delta/uuid"
	"github.com/streamrail/concurrent-map"
)

type CommandHandler func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error)

type Processor struct {
	handlers   map[string][]CommandHandler
	pending    cmap.ConcurrentMap
	Redirect   *Processor
	errors     int64
	exceptions int64
	total      int64
}

func newProcessor() *Processor {
	return &Processor{
		handlers: map[string][]CommandHandler{},
		pending:  cmap.New(),
	}
}

func (this *Processor) On(action string, handlers ...CommandHandler) error {
	this.handlers[action] = handlers
	return nil
}

func (this *Processor) wait(cmd *Command, cb func() error) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	wait := make(chan *Command)
	this.pending.Set(cmd.Key, wait)
	err := cb()
	if err != nil {
		this.pending.Remove(cmd.Key)
		return nil, err
	}
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

func (this *Processor) process(cmd *Command, conn *Connection) error {
	if cmd.Action == RESPONSE || cmd.Action == ERROR || cmd.Action == EXCEPTION {
		waiting, ok := this.pending.Get(cmd.Key)
		if ok {
			waiting.(chan *Command) <- cmd
			this.pending.Remove(cmd.Key)
			return nil
		}
		return nil
	}

	if this.Redirect != nil {
		return this.Redirect.process(cmd, conn)
	}

	atomic.AddInt64(&this.total, 1)
	handlers, ok := this.handlers[cmd.Action]
	if ok {
		result, err := this.trigger(cmd, conn, handlers...)
		this.respond(cmd, conn, result, err)
		return nil
	}
	return nil
}

func (this *Processor) respond(cmd *Command, conn *Connection, result interface{}, err error) {
	if err != nil {
		log.Println(err)
		log.Println(string(debug.Stack()))
		response := &Command{
			Key:    cmd.Key,
			Action: EXCEPTION,
			Body: map[string]interface{}{
				"message": err.Error(),
			},
		}
		if _, ok := err.(*DRSError); ok {
			response.Action = ERROR
			atomic.AddInt64(&this.errors, 1)
		} else {
			atomic.AddInt64(&this.exceptions, 1)
		}
		conn.Fire(response)
		return
	}
	conn.Fire(&Command{
		Key:    cmd.Key,
		Action: RESPONSE,
		Body:   result,
	})
}

func (this *Processor) trigger(cmd *Command, conn *Connection, handlers ...CommandHandler) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	ctx := make(map[string]interface{})
	for _, h := range handlers {
		result, err = h(cmd, conn, ctx)
		if err != nil {
			break
		}
	}
	return result, err
}
