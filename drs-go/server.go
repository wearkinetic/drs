package drs

import (
	"io"
	"sync"
	"time"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

type Server struct {
	*Processor
	sync.Mutex
	Protocol  protocol.Protocol
	transport Transport
	inbound   cmap.ConcurrentMap
	closed    bool

	connect    []func(conn *Connection, raw io.ReadWriteCloser) error
	disconnect []func(conn *Connection)
}

func NewServer(transport Transport) *Server {
	return &Server{
		Processor:  newProcessor(),
		Mutex:      sync.Mutex{},
		transport:  transport,
		Protocol:   protocol.JSON,
		inbound:    cmap.New(),
		connect:    make([]func(conn *Connection, raw io.ReadWriteCloser) error, 0),
		disconnect: make([]func(conn *Connection), 0),
	}
}

func (this *Server) OnConnect(cb func(*Connection, io.ReadWriteCloser) error) {
	this.connect = append(this.connect, cb)
}

func (this *Server) OnDisconnect(cb func(*Connection)) {
	this.disconnect = append(this.disconnect, cb)
}

func (this *Server) Broadcast(cmd *Command) int {
	for kv := range this.inbound.Iter() {
		kv.Val.(*Connection).Fire(cmd)
	}
	return len(this.inbound)
}

func (this *Server) Listen() error {
	return this.transport.Listen(func(raw io.ReadWriteCloser) {
		if this.closed {
			return
		}
		conn := NewConnection(this.Protocol)
		key := uuid.Ascending()
		this.inbound.Set(key, conn)
		defer func() {
			conn.Close()
			for _, cb := range this.disconnect {
				cb(conn)
			}
			this.inbound.Remove(key)
		}()

		for _, cb := range this.connect {
			err := cb(conn, raw)
			if err != nil {
				return
			}
		}
		conn.Redirect = this.Processor
		conn.handle(raw)
	})
}

func (this *Server) Close() {
	this.closed = true
	for kv := range this.inbound.IterBuffered() {
		kv.Val.(*Connection).Close()
	}
	for {
		if this.inbound.Count() == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
