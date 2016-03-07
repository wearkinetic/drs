package drs

import (
	"io"
	"time"

	"github.com/ironbay/delta/uuid"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/streamrail/concurrent-map"
)

type Server struct {
	*Processor
	Protocol  protocol.Protocol
	handlers  map[string][]CommandHandler
	transport Transport
	inbound   cmap.ConcurrentMap

	OnConnect    func(conn *Connection) error
	OnDisconnect func(conn *Connection)
}

func NewServer(transport Transport) *Server {
	return &Server{
		Processor: NewProcessor(),
		handlers:  make(map[string][]CommandHandler),
		transport: transport,
		Protocol:  protocol.JSON,
	}
}

func (this *Server) Broadcast(cmd *Command) {
	for kv := range this.inbound.Iter() {
		kv.Val.(*Connection).Fire(cmd)
	}
}

func (this *Server) Listen() error {
	this.On("drs.ping", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
		return time.Now().UnixNano() / 1000, nil
	})
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		conn := NewConnection(rw, this.Protocol)
		id := uuid.Ascending()
		this.inbound.Set(id, conn)
		defer this.inbound.Remove(id)

		if this.OnConnect != nil {
			if err := this.OnConnect(conn); err != nil {
				return
			}
		}
		conn.Redirect = this.Processor
		conn.Read()
		if this.OnDisconnect != nil {
			this.OnDisconnect(conn)
		}
	})
}
