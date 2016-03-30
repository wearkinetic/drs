package drs

import (
	"io"
	"sync"
	"time"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
)

type Server struct {
	*Processor
	sync.Mutex
	Protocol  protocol.Protocol
	transport Transport
	inbound   map[string]*Connection

	connect    []func(conn *Connection, raw io.ReadWriteCloser) error
	disconnect []func(conn *Connection)
}

func NewServer(transport Transport) *Server {
	return &Server{
		Processor:  newProcessor(),
		Mutex:      sync.Mutex{},
		transport:  transport,
		Protocol:   protocol.JSON,
		inbound:    map[string]*Connection{},
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
	for _, value := range this.inbound {
		value.Fire(cmd)
	}
	return len(this.inbound)
}

func (this *Server) Listen() error {
	return this.transport.Listen(func(raw io.ReadWriteCloser) {
		conn := NewConnection(this.Protocol)
		id := uuid.Ascending()
		this.Lock()
		this.inbound[id] = conn
		this.Unlock()
		defer func() {
			conn.Close()
			for _, cb := range this.disconnect {
				cb(conn)
			}
			this.Lock()
			delete(this.inbound, id)
			this.Unlock()
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
	for _, value := range this.inbound {
		value.Close()
	}
	for {
		if len(this.inbound) == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
