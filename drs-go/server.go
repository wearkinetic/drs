package drs

import (
	"io"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Server struct {
	*Processor
	handlers  map[string][]CommandHandler
	transport Transport
	proto     protocol.Protocol
}

func NewServer(transport Transport, proto protocol.Protocol) *Server {
	return &Server{
		Processor: NewProcessor(),
		handlers:  make(map[string][]CommandHandler),
		transport: transport,
		proto:     proto,
	}
}

func (this *Server) Listen() error {
	return this.transport.Listen(func(rw io.ReadWriteCloser) {
		conn := NewConnection(rw, this.proto)
		conn.Redirect = this.Processor
		conn.Read()
	})
}
