package drs

import (
	"io"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/dynamic"
	"github.com/streamrail/concurrent-map"
)

type Connection struct {
	*Processor
	Stream *protocol.Stream
	Cache  cmap.ConcurrentMap
}

func Dial(proto protocol.Protocol, transport Transport, host string) (*Connection, error) {
	raw, err := transport.Connect(host)
	if err != nil {
		return nil, err
	}
	return Accept(proto, raw)
}

func Accept(proto protocol.Protocol, raw io.ReadWriteCloser) (*Connection, error) {
	conn := &Connection{
		Processor: NewProcessor(),
		Stream:    proto(raw),
		Cache:     cmap.New(),
	}
	return conn, nil
}

func (this *Connection) Read() error {
	var err error
	for {
		cmd := new(Command)
		if err = this.Stream.Decode(cmd); err != nil {
			break
		}
		this.Process(this, cmd)
	}
	return err
}

func (this *Connection) Call(cmd *Command) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = "1234"
	}
	for {
		block := this.Enqueue(cmd.Key)
		err := this.Stream.Encode(cmd)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		result := <-block
		switch result.Action {
		case EXCEPTION:
			time.Sleep(1 * time.Second)
			continue
		case ERROR:
			message := dynamic.String(result.Map(), "message")
			return nil, Error(message)
		case RESPONSE:
			return result.Body, nil
		}
	}
}

func (this *Connection) Close() {
	this.Stream.Close()
}

func (this *Connection) respond(key string, resp interface{}, err error) {
	cmd := &Command{
		Key: key,
	}
	if err == nil {
		cmd.Action = RESPONSE
		cmd.Body = resp
	}
	if _, ok := err.(*DRSError); ok {
		cmd.Action = ERROR
		cmd.Body = err
	}
	if _, ok := err.(*DRSException); ok {
		cmd.Action = EXCEPTION
		cmd.Body = err
	}
	this.Stream.Encode(cmd)
}
