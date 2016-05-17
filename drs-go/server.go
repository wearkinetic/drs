package drs

import (
	"io"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Server struct {
	*Processor
	protocol   protocol.Protocol
	transport  Transport
	connect    []func(*Connection) error
	disconnect []func(*Connection)
}

func New(transport Transport, protocol protocol.Protocol) *Server {
	return &Server{
		Processor:  NewProcessor(),
		protocol:   protocol,
		transport:  transport,
		connect:    []func(*Connection) error{},
		disconnect: []func(*Connection){},
	}
}

func (this *Server) Listen(host string) error {
	return this.transport.Listen(host, func(raw io.ReadWriteCloser) {
		conn := Accept(this.protocol, raw)
		conn.parent = this.Processor
		for _, cb := range this.connect {
			err := cb(conn)
			if err != nil {
				conn.Close()
				return
			}
		}
		defer func() {
			for _, cb := range this.disconnect {
				cb(conn)
			}
		}()
		conn.Read()
	})
}

func (this *Server) OnConnect(cb func(*Connection) error) {
	this.connect = append(this.connect, cb)
}

func (this *Server) OnDisconnect(cb func(*Connection)) {
	this.disconnect = append(this.disconnect, cb)
}
