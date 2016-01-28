package drs

import (
	"fmt"
	"io"
	"time"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
)

const PORT = 12000

const (
	ERROR     = "drs.error"
	RESPONSE  = "drs.response"
	EXCEPTION = "drs.exception"
)

type Pipe struct {
	transport   Transport
	Router      RouterHandler
	Protocol    protocol.Protocol
	handlers    map[string][]CommandHandler
	connections map[string]*Connection
	pending     map[string]chan *Command
	Events      *Events
}

type Events struct {
	Connect func(conn *Connection) error
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
		Events:      new(Events),
	}, nil
}

func (this *Pipe) On(action string, handlers ...CommandHandler) error {
	this.handlers[action] = handlers
	return this.transport.On(action)
}

func (this *Pipe) Send(cmd *Command) (interface{}, error) {
	cmd.Key = uuid.Ascending()
	for {
		conn, err := this.route(cmd.Action)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		wait := make(chan *Command)
		this.pending[cmd.Key] = wait
		err = conn.Encode(cmd)
		response := <-wait
		if response.Action == ERROR {
			return nil, &DRSError{
				Message: response.Dynamic().String("message"),
			}
		}
		if response.Action == EXCEPTION {
			time.Sleep(time.Second)
			continue
		}
		// TODO: Handle exceptions vs errors
		return response.Body, err
	}
}

func (this *Pipe) Listen() error {
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		conn := this.connect(rw)
		if this.Events.Connect != nil {
			err := this.Events.Connect(conn)
			if err != nil {
				return
			}
		}
		this.handle(conn)
	})
}

func (this *Pipe) Process(conn *Connection, cmd *Command) {
	defer func() {
		if r := recover(); r != nil {
			response := &Command{
				Key:    cmd.Key,
				Action: EXCEPTION,
				Body:   fmt.Sprint(r),
			}
			conn.Encode(response)
		}
	}()

	if cmd.Action == RESPONSE || cmd.Action == ERROR || cmd.Action == EXCEPTION {
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
	ctx := make(Dynamic)
	var result interface{}
	var err error
	for _, h := range handlers {
		result, err = h(cmd, conn, ctx)
		if err != nil {
			break
		}
	}
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
		conn.Encode(response)
		return
	}
	conn.Encode(&Command{
		Key:    cmd.Key,
		Action: RESPONSE,
		Body:   result,
	})
}

func (this *Pipe) route(action string) (*Connection, error) {
	host, err := this.Router(action)
	if err != nil {
		return nil, err
	}
	return this.dial(host)
}
