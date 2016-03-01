package drs

import (
	"io"
	"log"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/dynamic"
	"github.com/streamrail/concurrent-map"
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
	write   chan *Command
	pending cmap.ConcurrentMap
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
		write:     make(chan *Command),
		pending:   cmap.New(),
	}
}

func (this *Connection) Send(cmd *Command) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	wait := make(chan *Command)
	this.pending.Set(cmd.Key, wait)
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
	if response.Action == EXCEPTION {
		args := response.Map()
		return nil, &DRSError{
			Message: dynamic.String(args, "message"),
			Kind:    dynamic.String(args, "kind"),
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
			continue
		}
		log.Println(cmd)
		if cmd.Action == RESPONSE || cmd.Action == ERROR || cmd.Action == EXCEPTION {
			waiting, ok := this.pending.Get(cmd.Key)
			if ok {
				waiting.(chan *Command) <- cmd
				this.pending.Remove(cmd.Key)
				continue
			}
		}
		go func() {
			result, err := this.process(cmd, this)
			log.Println(result, err)
			this.respond(this, cmd, result, err)
		}()
	}
}
