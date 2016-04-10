package drs

import (
	"log"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/go-util/console"
)

type Connection2 struct {
	*Processor
	outgoing  chan *Command
	close     chan bool
	bootstrap []*Command
}

func NewConnection2() *Connection2 {
	result := &Connection2{
		Processor: newProcessor(),
		outgoing:  make(chan *Command, 1),
		close:     make(chan bool),
		bootstrap: []*Command{},
	}
	return result
}

func (this *Connection2) Dial(proto protocol.Protocol, transport Transport, host string) {
	for {
		raw, err := transport.Connect(host)
		if err == nil {
			if this.handle(proto(raw)) {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
	log.Println("Closing")
	close(this.outgoing)
}

func (this *Connection2) handle(stream *protocol.Stream) bool {
	defer stream.Close()
	incoming := make(chan *Command)

	go func() {
		for {
			cmd := new(Command)
			if err := stream.Decode(cmd); err != nil {
				log.Println(err)
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
			this.process(cmd, this)

		case cmd := <-this.outgoing:
			if err := stream.Encode(cmd); err != nil {
				go this.Fire(cmd)
			}
			console.JSON(cmd)

		case <-this.close:
			return true
		}
	}
}

func (this *Connection2) Fire(cmd *Command) {
	this.outgoing <- cmd
}

func (this *Connection2) Close() {
	this.close <- true
}
