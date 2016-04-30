package drs

import (
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

type Connection struct {
	*Processor
	Cache     cmap.ConcurrentMap
	Raw       io.ReadWriteCloser
	outgoing  chan *Command
	close     chan bool
	bootstrap []*Command
}

func NewConnection() *Connection {
	result := &Connection{
		Processor: newProcessor(),
		Cache:     cmap.New(),
		outgoing:  make(chan *Command, 500),
		close:     make(chan bool, 1),
		bootstrap: []*Command{},
	}
	return result
}

func (this *Connection) Dial(proto protocol.Protocol, transport Transport, host string, reconnect bool) {
	for {
		raw, err := transport.Connect(host)
		if err == nil {
			if this.handle(proto(raw)) {
				break
			}
		}
		if !reconnect {
			break
		}
		time.Sleep(1 * time.Second)
	}
	close(this.outgoing)
}

func (this *Connection) handle(stream *protocol.Stream) bool {
	defer stream.Close()
	this.Raw = stream.Raw
	incoming := make(chan *Command)

	go func() {
		for {
			cmd := new(Command)
			if err := stream.Decode(cmd); err != nil {
				break
			}
			incoming <- cmd
		}
		close(incoming)
	}()

	for _, cmd := range this.bootstrap {
		if err := stream.Encode(cmd); err != nil {
			return false
		}
	}

	for {
		select {

		case cmd := <-incoming:
			if cmd == nil {
				return false
			}
			go func() {
				res, err := this.Process(cmd, this)
				this.respond(cmd.Key, res, err)
			}()

		case cmd := <-this.outgoing:
			if err := stream.Encode(cmd); err != nil {
				go this.Fire(cmd)
			}

		case <-this.close:
			return true
		}
	}
}

func (this *Connection) Fire(cmd *Command) {
	this.outgoing <- cmd
}

func (this *Connection) Bootstrap(cmd *Command) (interface{}, error) {
	this.bootstrap = append(this.bootstrap, cmd)
	return this.Request(cmd)
}

func (this *Connection) Request(cmd *Command) (interface{}, error) {
	for {
		res, err := this.wait(cmd, func() {
			this.Fire(cmd)
		})
		if err != nil {
			return res, nil
		}
		if _, ok := err.(*DRSException); ok {
			continue
		}
		return nil, err
	}
}

func (this *Connection) respond(key string, res interface{}, err error) {
	if err != nil {
		log.Println(err)
		response := &Command{
			Key:    key,
			Action: EXCEPTION,
			Body: map[string]interface{}{
				"message": err.Error(),
			},
		}
		if _, ok := err.(*DRSError); ok {
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
	this.close <- true
}
