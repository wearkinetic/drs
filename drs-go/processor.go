package drs

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

func (this *Processor) respond(conn *Connection, cmd *Command, result interface{}, err error) {
	if err != nil {
		response := &Command{
			Key:    cmd.Key,
			Action: EXCEPTION,
			Body: &DRSError{
				Message: err.Error(),
			},
		}
		if casted, ok := err.(*DRSError); ok {
			response.Action = ERROR
			response.Body = casted
		}
		conn.Send(response)
		return
	}
	conn.Send(&Command{
		Key:    cmd.Key,
		Action: RESPONSE,
		Body:   result,
	})
}

func (this *Processor) trigger(cmd *Command, conn *Connection, handlers ...CommandHandler) (interface{}, error) {
	ctx := make(map[string]interface{})
	var result interface{}
	var err error
	for _, h := range handlers {
		result, err = h(cmd, conn, ctx)
		if err != nil {
			break
		}
	}
	return result, err
}
