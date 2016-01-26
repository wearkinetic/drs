package drs

import (
	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
)

type Pipe struct {
	transport   Transport
	Router      RouterHandler
	Protocol    protocol.Protocol
	handlers    map[string][]CommandHandler
	connections map[string]*Connection
	pending     map[string]chan *Command
}

func New(transport Transport) (*Pipe, error) {
	return &Pipe{
		transport: transport,
		Router: func(action string) (string, error) {
			return "", nil
		},
		Protocol:    protocol.JSON,
		handlers:    make(map[string][]CommandHandler),
		connections: make(map[string]*Connection),
		pending:     map[string]chan *Command{},
	}, nil
}

func (this *Pipe) On(action string, handlers ...CommandHandler) error {
	this.handlers[action] = handlers
	return this.transport.On(action)
}

func (this *Pipe) Send(cmd *Command) (interface{}, error) {
	cmd.Key = uuid.Ascending()
	conn, err := this.route(cmd.Action)
	if err != nil {
		return nil, err
	}
	wait := make(chan *Command)
	this.pending[cmd.Key] = wait
	err = conn.protocol.Encode(cmd)
	response := <-wait
	// TODO: Handle exceptions vs errors
	return response.Body, err
}

func (this *Pipe) Listen() error {
	return this.transport.Listen(this.register)
}

func (this *Pipe) Process(conn *Connection, cmd *Command) {
	if cmd.Action == "response" || cmd.Action == "error" {
		waiting, ok := this.pending[cmd.Key]
		if ok {
			waiting <- cmd
			delete(this.pending, cmd.Key)
			return
		}
	}

	handlers, ok := this.handlers[cmd.Action]
	if !ok {
		return
	}
	result, _ := handlers[0](cmd, conn)
	response := &Command{
		Key:    cmd.Key,
		Action: "response",
		Body:   result,
	}
	conn.protocol.Encode(response)
}

func (this *Pipe) route(action string) (*Connection, error) {
	host, err := this.Router(action)
	if err != nil {
		return nil, err
	}
	return this.connect(host)
}
