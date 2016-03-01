package drs

import (
	"io"
	"log"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Pipe struct {
	Protocol  protocol.Protocol
	Router    RouterHandler
	transport Transport
	Events    *Events
	*Processor
	connections map[string]*Connection
}

type Events struct {
	Connect    func(conn *Connection) error
	Disconnect func(conn *Connection)
}

func New(transport Transport) *Pipe {
	return &Pipe{
		Processor:   NewProcessor(),
		Protocol:    protocol.JSON,
		Events:      new(Events),
		connections: map[string]*Connection{},
		transport:   transport,
	}
}

func (this *Pipe) Send(cmd *Command) (interface{}, error) {
	for {
		conn, err := this.route(cmd.Action)
		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}
		result, err := conn.Send(cmd)
		if err != nil {
			casted, ok := err.(*DRSError)
			if !ok || casted.Kind == "exception" {
				time.Sleep(1 * time.Second)
				continue
			}
		}
		log.Println(result)
		return result, err
	}
}

func (this *Pipe) route(action string) (*Connection, error) {
	host, err := this.Router(action)
	if err != nil {
		return nil, err
	}
	match, ok := this.connections[host]
	if ok {
		return match, nil
	}
	conn, err := Dial(this.transport, this.Protocol, host)
	if err != nil {
		return nil, err
	}
	conn.Redirect = this.Processor
	this.connections[host] = conn
	go func() {
		conn.Read()
		delete(this.connections, host)
	}()
	return conn, nil
}

func (this *Pipe) Listen() error {
	this.On("drs.ping", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
		return time.Now().Unix() / 1000, nil
	})
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		conn := NewConnection(rw, this.Protocol)
		if this.Events.Connect != nil {
			if err := this.Events.Connect(conn); err != nil {
				return
			}
		}
		conn.Redirect = this.Processor
		conn.Read()
		if this.Events.Disconnect != nil {
			this.Events.Disconnect(conn)
		}
	})
}

/*
type Pipe struct {
	transport   Transport
	Router      RouterHandler
	Protocol    protocol.Protocol
	handlers    map[string][]CommandHandler
	connections map[string]*Connection
	pending     map[string]chan *Command
	Events      *Events
	block       chan bool
}

type Events struct {
	Connect    func(conn *Connection) error
	Disconnect func(conn *Connection) error
}

func New(transport Transport) (*Pipe, error) {
	return &Pipe{
		transport: transport,
		Router: func(action string) (string, error) {
			return ":12000", nil
		},
		Protocol:    protocol.JSON,
		handlers:    make(map[string][]CommandHandler),
		connections: make(map[string]*Connection),
		Events:      new(Events),
		block:       make(chan bool, 1),
	}, nil
}

func (this *Pipe) On(action string, handlers ...CommandHandler) error {
	this.handlers[action] = handlers
}

func (this *Pipe) Send(cmd *Command) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	this.block <- true
	for {
		conn, err := this.route(cmd.Action)
		if err != nil {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		wait := make(chan *Command)
		this.pending[cmd.Key] = wait
		err = conn.Send(cmd)
		if err != nil {
			continue
		}
		<-this.block
		response := <-wait
		if response.Action == ERROR {
			return nil, &DRSError{
				Message: "Error",
			}
		}
		if response.Action == EXCEPTION {
			log.Println(response.Body)
			time.Sleep(time.Second)
			this.block <- true
			continue
		}
		// TODO: Handle exceptions vs errors
		return response.Body, err
	}
}

func (this *Pipe) Listen() error {
	this.On("drs.ping", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
			conn.Encode(&Command{
				Action: "ping",
			})
		return time.Now().Unix(), nil
	})
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		conn := this.connect(rw)
		if this.Events.Connect != nil {
			err := this.Events.Connect(conn)
			if err != nil {
				return
			}
		}
		this.handle(conn)
		if this.Events.Disconnect != nil {
			this.Events.Disconnect(conn)
		}
	})
}

func (this *Pipe) process(conn *Connection, cmd *Command) {
	defer func() {
		if r := recover(); r != nil {
			response := &Command{
				Key:    cmd.Key,
				Action: EXCEPTION,
				Body:   fmt.Sprint(r),
			}
			log.Println(r)
			conn.Send(response)
		}
	}()

	if cmd.Action == RESPONSE || cmd.Action == ERROR || cmd.Action == EXCEPTION {
		waiting, ok := this.pending[cmd.Key]
		if ok {
			waiting <- cmd
			delete(this.pending, cmd.Key)
			return
		}
		return
	}

	handlers, ok := this.handlers[cmd.Action]
	if !ok {
		return
	}
	ctx := make(map[string]interface{})
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
		conn.Send(response)
		return
	}
	conn.Send(&Command{
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
*/
