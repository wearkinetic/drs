package drs

import (
	"encoding/json"
	"io"
	"log"

	"github.com/ironbay/delta/uuid"
)

type DRS struct {
	transport   Transport
	Router      RouterHandler
	handlers    map[string][]CommandHandler
	connections map[string]*Connection
	pending     map[string]chan *Command
}

func New(transport Transport) (*DRS, error) {
	result := new(DRS)
	result.transport = transport
	result.Router = func(action string) (string, error) {
		return "", nil
	}
	return &DRS{
		transport: transport,
		Router: func(action string) (string, error) {
			return "", nil
		},
		handlers:    make(map[string][]CommandHandler),
		connections: make(map[string]*Connection),
		pending:     map[string]chan *Command{},
	}, nil
}

func (this *DRS) On(action string, handlers ...CommandHandler) error {
	this.handlers[action] = handlers
	return this.transport.On(action)
}

func (this *DRS) Send(cmd *Command) (interface{}, error) {
	cmd.Key = uuid.Ascending()
	conn, err := this.route(cmd.Action)
	if err != nil {
		return nil, err
	}
	wait := make(chan *Command)
	this.pending[cmd.Key] = wait
	err = conn.Send(cmd)
	response := <-wait
	return response.Body, err
}

func (this *DRS) Listen() error {
	return this.transport.Listen(this.connection)
}

func (this *DRS) connection(rw io.ReadWriteCloser) (chan bool, *Connection) {
	conn := NewConnection(rw)
	done := make(chan bool)
	go func() {
		for {
			data, err := this.transport.Frame(rw)
			if err != nil && err.Error() == "EOF" {
				log.Println(err)
				break
			}
			cmd := new(Command)
			json.Unmarshal(data, cmd)
			this.Process(conn, cmd)
		}
		done <- true
	}()
	return done, conn
}

func (this *DRS) Process(conn *Connection, cmd *Command) {
	if cmd.Action == "response" || cmd.Action == "error" {
		waiting, ok := this.pending[cmd.Key]
		if ok {
			waiting <- cmd
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
	conn.Send(response)
}

func (this *DRS) route(action string) (*Connection, error) {
	host, err := this.Router(action)
	if err != nil {
		return nil, err
	}
	connection, ok := this.connections[action]
	if ok {
		return connection, nil
	}
	rw, err := this.transport.Connect(host)
	if err != nil {
		return nil, err
	}
	_, conn := this.connection(rw)
	this.connections[action] = conn
	return conn, nil
}
