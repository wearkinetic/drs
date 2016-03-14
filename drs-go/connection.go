package drs

import (
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

type Connection struct {
	*Processor
	Cache      cmap.ConcurrentMap
	status     int32
	protocol   protocol.Protocol
	stream     *protocol.Stream
	raw        io.ReadWriteCloser
	Raw        io.ReadWriteCloser
	connecting sync.Mutex
}

const (
	OFFLINE      = int32(0)
	ONLINE       = int32(1)
	RECONNECTING = int32(2)
	CLOSED       = int32(3)
)

func NewConnection(protocol protocol.Protocol) *Connection {
	result := &Connection{
		Processor: newProcessor(),
		Cache:     cmap.New(),
		status:    OFFLINE,
		protocol:  protocol,
	}
	return result
}

func (this *Connection) Dial(transport Transport, host string, reconnect bool) {
	for {
		this.connecting.Lock()
		if atomic.LoadInt32(&this.status) == CLOSED {
			return
		}
		atomic.StoreInt32(&this.status, RECONNECTING)
		raw, err := transport.Connect(host)
		if err != nil {
			this.connecting.Unlock()
			time.Sleep(1 * time.Second)
			continue
		}
		this.accept(raw)
		this.connecting.Unlock()
		err = this.handle()
		if !reconnect {
			return
		}
		time.Sleep(1 * time.Second)
	}
}

func (this *Connection) accept(raw io.ReadWriteCloser) {
	this.raw = raw
	this.Raw = raw
	this.stream = this.protocol(raw)
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
		snap := atomic.LoadInt32(&this.status)
		if snap == CLOSED {
			return errors.New("Connection has been closed")
		}
		if snap == ONLINE {
			var err error
			if err = this.stream.Encode(cmd); err == nil {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func (this *Connection) handle() error {
	atomic.StoreInt32(&this.status, ONLINE)
	var err error
	for {
		cmd := new(Command)
		err = this.stream.Decode(cmd)
		if err != nil {
			break
		}
		this.process(cmd, this)
	}
	atomic.StoreInt32(&this.status, OFFLINE)
	return err
}

func (this *Connection) Close() {
	this.connecting.Lock()
	atomic.StoreInt32(&this.status, CLOSED)
	this.raw.Close()
	this.connecting.Unlock()
}
