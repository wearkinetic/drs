package drs

import "log"

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
			Body:   err.Error(),
		}
		if _, ok := err.(*DRSError); ok {
			response.Action = ERROR
		}
		conn.stream.Encode(response)
		return
	}
	conn.stream.Encode(&Command{
		Key:    cmd.Key,
		Action: RESPONSE,
		Body:   result,
	})
}

func (this *Processor) trigger(cmd *Command, conn *Connection, handlers ...CommandHandler) (result interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(err)
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
