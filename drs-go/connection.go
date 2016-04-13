package drs

import (
	"io"
	"sync/atomic"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/go-util/actor"
	"github.com/streamrail/concurrent-map"
)

type Connection struct {
	*Processor
	stream *protocol.Stream
	Raw    io.ReadWriteCloser
	Cache  cmap.ConcurrentMap

	OnDisconnect func(err error)
}

func NewConnection() *Connection {
	result := &Connection{
		Processor: newProcessor(),
		Cache:     cmap.New(),
	}
	return result
}

func (this *Connection) Dial(proto protocol.Protocol, transport Transport, host string) error {
	raw, err := transport.Connect(host)
	if err != nil {
		return err
	}
	this.Raw = raw
	this.stream = proto(raw)
	go this.handle()
	return nil
}

func (this *Connection) handle() {
	defer this.stream.Close()
	defer this.clear()
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			this.Request(&Command{
				Action: "drs.ping",
			})
		}
	}()
	for {
		cmd := new(Command)
		if err := this.stream.Decode(cmd); err != nil {
			if this.OnDisconnect != nil {
				this.OnDisconnect(err)
			}
			return
		}
		go func() {
			result, err := this.Process(cmd, this)
			if result != nil {
				this.respond(cmd.Key, result, err)
			}
		}()
	}
}

func (this *Connection) Fire(cmd *Command) error {
	return this.stream.Encode(cmd)
}

func (this *Connection) Request(cmd *Command) (interface{}, error) {
	return this.wait(cmd, func() error {
		return this.Fire(cmd)
	})
}

func (this *Connection) respond(key string, res interface{}, err error) {
	if err != nil {
		response := &Command{
			Key:    key,
			Action: EXCEPTION,
			Body: map[string]interface{}{
				"message": err.Error(),
			},
		}
		if _, ok := err.(*actor.ActorError); ok {
			response.Action = ERROR
			atomic.AddInt64(&this.errors, 1)
		} else {
			atomic.AddInt64(&this.exceptions, 1)
		}
		this.Fire(response)
		return
	}
	this.Fire(&Command{
		Key:    key,
		Action: RESPONSE,
		Body:   res,
	})
}

func (this *Connection) Close() {
	if this.stream == nil {
		return
	}
	this.stream.Close()
}
