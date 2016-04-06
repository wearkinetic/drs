package drs

import (
	"log"
	"sync"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/go-util/console"
)

type Connection2 struct {
	sync.RWMutex
	stream   *protocol.Stream
	incoming chan *Command
	outgoing chan *Command
	closed   bool
}

func NewConnection2() *Connection2 {
	result := &Connection2{
		incoming: make(chan *Command, 1),
		outgoing: make(chan *Command, 1),
	}
	go result.read()
	go result.write()
	return result
}

func (this *Connection2) read() {
	for cmd := range this.incoming {
		if this.closed {
			break
		}
		console.JSON(cmd)
	}
}

func (this *Connection2) write() {
	for cmd := range this.outgoing {
		for {
			if this.stream != nil {
				time.Sleep(1 * time.Second)
				err := this.stream.Encode(cmd)
				if err == nil {
					console.JSON(cmd)
					break
				}
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (this *Connection2) Dial(proto protocol.Protocol, transport Transport, host string) {
	for {
		this.Lock()
		if this.closed {
			this.Unlock()
			break
		}
		raw, err := transport.Connect(host)
		if err != nil {
			this.Unlock()
		}
		if err == nil {
			this.stream = proto(raw)
			this.Unlock()
			for {
				cmd := new(Command)
				if err = this.stream.Decode(cmd); err != nil {
					log.Println(err)
					break
				}
				this.incoming <- cmd
			}
		}
		// time.Sleep(1 * time.Second)
	}
	log.Println("Closing")
	close(this.incoming)
	close(this.outgoing)
}

func (this *Connection2) Fire(cmd *Command) {
	this.outgoing <- cmd
}

func (this *Connection2) Close() {
	this.Lock()
	this.closed = true
	if this.stream != nil {
		this.stream.Close()
	}
	this.Unlock()
}
