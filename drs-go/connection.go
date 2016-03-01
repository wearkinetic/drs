package drs

import (
	"io"
	"log"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
)

const (
	ERROR     = "drs.error"
	RESPONSE  = "drs.response"
	EXCEPTION = "drs.exception"
)

type Connection struct {
	*Processor
	Raw     io.ReadWriteCloser
	Cache   map[string]interface{}
	stream  *protocol.Stream
	pending map[string]chan *Command
}

func Dial(transport Transport, proto protocol.Protocol, host string) (*Connection, error) {
	rw, err := transport.Connect(host)
	if err != nil {
		return nil, err
	}
	conn := NewConnection(rw, proto)
	return conn, nil
}

func NewConnection(rw io.ReadWriteCloser, proto protocol.Protocol) *Connection {
	return &Connection{
		Processor: NewProcessor(),
		Raw:       rw,
		Cache:     map[string]interface{}{},
		stream:    proto(rw),
		pending:   map[string]chan *Command{},
	}
}

func (this *Connection) Send(cmd *Command) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	wait := make(chan *Command)
	this.pending[cmd.Key] = wait
	err := this.stream.Encode(cmd)
	if err != nil {
		return nil, err
	}
	response := <-wait
	if response.Action == ERROR {
		return nil, &DRSError{
			Message: response.Body.(string),
		}
	}
	if response.Action == EXCEPTION || response.Action == ERROR {
		args := cmd.Map()
		return nil, &DRSError{
			Message: args["message"].(string),
			Kind:    args["kind"].(string),
		}
	}
	return response.Body, nil
}

func (this *Connection) Read() {
	for {
		cmd := new(Command)
		err := this.stream.Decode(cmd)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Println(err)
			continue
		}
		if cmd.Action == RESPONSE || cmd.Action == ERROR || cmd.Action == EXCEPTION {
			waiting, ok := this.pending[cmd.Key]
			if ok {
				waiting <- cmd
				delete(this.pending, cmd.Key)
				continue
			}
		}
		go func() {
			result, err := this.process(cmd, this)
			this.respond(this, cmd, result, err)
		}()
	}
}
