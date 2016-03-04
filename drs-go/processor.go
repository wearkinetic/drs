package drs

import (
	"log"
	"runtime/debug"
	"sync/atomic"
)

type CommandHandler func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error)

type Processor struct {
	handlers map[string][]CommandHandler
	Redirect *Processor
}

func NewProcessor() *Processor {
	return &Processor{
		handlers: map[string][]CommandHandler{},
	}
}

func (this *Processor) On(action string, handlers ...CommandHandler) error {
	this.handlers[action] = handlers
	return nil
}

func (this *Processor) process(cmd *Command, conn *Connection) (interface{}, error) {
	if this.Redirect != nil {
		return this.Redirect.process(cmd, conn)
	}
	atomic.AddInt64(&total, 1)
	{
		handlers, ok := this.handlers[cmd.Action]
		if ok {
			return this.trigger(cmd, conn, handlers...)
		}
	}
	return nil, nil
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
			atomic.AddInt64(&cerr, 1)
		} else {
			atomic.AddInt64(&exceptions, 1)
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
			log.Println(err)
			log.Println(string(debug.Stack()))
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
