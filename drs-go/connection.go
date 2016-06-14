package drs

import (
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/dynamic"
	"github.com/satori/go.uuid"
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
	return Accept(proto, raw), nil
}

func Accept(proto protocol.Protocol, raw io.ReadWriteCloser) *Connection {
	conn := &Connection{
		Processor: NewProcessor(),
		Stream:    proto(raw),
		Cache:     cmap.New(),
	}
	return conn
}

func (this *Connection) Read() error {
	var err error
	for {
		cmd := new(Command)
		if err = this.Stream.Decode(cmd); err != nil {
			break
		}
		go this.Process(this, cmd)
	}
	return err
}

func (this *Connection) Call(cmd *Command) (interface{}, error) {
	if cmd.Key == "" {
		cmd.Key = uuid.NewV4().String()
	}
	match, ok := this.stats.Get(cmd.Action)
	if !ok {
		match = new(Stats)
		this.stats.Set(cmd.Action, match)
	}
	stats := match.(*Stats)
	for {
		block := this.Enqueue(cmd.Key)
		err := this.Fire(cmd)
		if err != nil {
			return nil, err
		}
		result := <-block
		switch result.Action {
		case EXCEPTION:
			atomic.AddInt64(&stats.Exceptions, 1)
			time.Sleep(1 * time.Second)
			continue
		case ERROR:
			message := dynamic.String(result.Map(), "message")
			atomic.AddInt64(&stats.Errors, 1)
			return nil, Error(message)
		case RESPONSE:
			atomic.AddInt64(&stats.Errors, 1)
			return result.Body, nil
		}
	}
}

func (this *Connection) Fire(cmd *Command) error {
	return this.Stream.Encode(cmd)
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
	} else {
		if _, ok := err.(*DRSError); ok {
			cmd.Action = ERROR
			cmd.Body = err
		} else {
			log.Println(err)
			cmd.Action = EXCEPTION
			cmd.Body = err
		}
	}
	this.Stream.Encode(cmd)
}
