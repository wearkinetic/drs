package drs

import (
	"errors"
	"io"
	"sync"
	"time"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

type Connection struct {
	*Processor
	sync.RWMutex
	Cache    cmap.ConcurrentMap
	closed   bool
	protocol protocol.Protocol
	stream   *protocol.Stream

	connect []func() error
}

func NewConnection(protocol protocol.Protocol) *Connection {
	result := &Connection{
		Processor: newProcessor(),
		Cache:     cmap.New(),
		protocol:  protocol,
		RWMutex:   sync.RWMutex{},
		connect:   []func() error{},
	}
	return result
}

func (this *Connection) Raw() io.ReadWriteCloser {
	this.RLock()
	defer this.RUnlock()
	if this.stream == nil {
		return nil
	}
	return this.stream.Raw
}

func (this *Connection) Closed() bool {
	this.RLock()
	defer this.RUnlock()
	return this.closed
}

func (this *Connection) Open() bool {
	this.RLock()
	defer this.RUnlock()
	return this.stream != nil
}

func (this *Connection) Dial(transport Transport, host string, reconnect bool) {
	for {
		raw, err := transport.Connect(host)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		err = this.handle(raw)
		if this.closed || !reconnect {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (this *Connection) Request(cmd *Command) (interface{}, error) {
	return this.wait(cmd, func() error {
		return this.Fire(cmd)
	})
}

func (this *Connection) Fire(cmd *Command) error {
	if cmd.Key == "" {
		cmd.Key = uuid.Ascending()
	}
	for {
		if this.Closed() {
			return errors.New("Connection has been closed")
		}
		if this.Open() {
			this.RLock()
			err := this.stream.Encode(cmd)
			this.RUnlock()
			if err == nil {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (this *Connection) handle(raw io.ReadWriteCloser) error {
	this.Lock()
	if this.closed {
		this.Unlock()
		return errors.New("Connection has been closed")
	}
	this.stream = this.protocol(raw)
	this.Unlock()
	// TODO: Considering using channels properly
	var err error
	buffer := make(chan bool, 500)
	for {
		cmd := new(Command)
		err = this.stream.Decode(cmd)
		if err != nil {
			break
		}
		buffer <- true
		go func() {
			this.process(cmd, this)
			<-buffer
		}()
	}
	return err
}

func (this *Connection) Close() {
	this.Lock()
	this.closed = true
	if this.stream != nil {
		this.stream.Close()
	}
	this.clear()
	this.Unlock()
}
