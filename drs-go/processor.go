package drs

import (
	"errors"
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
		handlers:   map[string][]CommandHandler{},
		pending:    cmap.New(),
		total:      0,
		exceptions: 0,
		errors:     0,
	}
}

func (this *Processor) On(action string, handlers ...CommandHandler) error {
	exists, ok := this.handlers[action]
	if ok {
		exists = append(exists, handlers...)
		return nil
	}
	this.handlers[action] = handlers
	return nil
}

func (this *Processor) wait(cmd *Command, cb func() error) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	wait := make(chan *Command, 1)
	err := cb()
	if err != nil {
		return nil, err
	}
	this.pending.Set(cmd.Key, wait)
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
		}
		return nil
	}

	if this.Redirect != nil {
		return this.Redirect.process(cmd, conn)
	}

	// atomic.AddInt64(&this.total, 1)
	result, err := this.Trigger(cmd, conn)
	this.respond(cmd, conn, result, err)
	return nil
}

func (this *Processor) respond(cmd *Command, conn *Connection, result interface{}, err error) {
	if err != nil {
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
			log.Println(cmd.Action, err)
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

func (this *Processor) Trigger(cmd *Command, conn *Connection) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			log.Println(err)
			log.Println(string(debug.Stack()))
		}
	}()
	handlers, ok := this.handlers[cmd.Action]
	if !ok {
		return nil, errors.New("No handlers for this action")
	}
	ctx := make(map[string]interface{})
	for _, h := range handlers {
		result, err = h(cmd, conn, ctx)
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
