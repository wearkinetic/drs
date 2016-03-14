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
	Protocol  protocol.Protocol
	transport Transport
	inbound   map[string]*Connection
	mutex     sync.Mutex

	OnConnect    func(conn *Connection) error
	OnDisconnect func(conn *Connection)
}

func NewServer(transport Transport) *Server {
	return &Server{
		Processor: newProcessor(),
		transport: transport,
		Protocol:  protocol.JSON,
		inbound:   map[string]*Connection{},
		mutex:     sync.Mutex{},
	}
}

func (this *Server) Broadcast(cmd *Command) int {
	for _, value := range this.inbound {
		value.Fire(cmd)
	}
	return len(this.inbound)
}

func (this *Server) Listen() error {
	this.On("drs.ping", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
		return time.Now().UnixNano() / int64(time.Millisecond), nil
	})
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		conn := NewConnection(this.Protocol)
		conn.accept(rw)
		id := uuid.Ascending()
		this.mutex.Lock()
		this.inbound[id] = conn
		this.mutex.Unlock()
		defer func() {
			this.mutex.Lock()
			delete(this.inbound, id)
			this.mutex.Unlock()
		}()

		if this.OnConnect != nil {
			if err := this.OnConnect(conn); err != nil {
				return
			}
		}
		conn.Redirect = this.Processor
		conn.handle()
		conn.Close()
		if this.OnDisconnect != nil {
			this.OnDisconnect(conn)
		}
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
