package drs

import (
	"io"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Server struct {
	protocol  protocol.Protocol
	transport Transport
	Events    struct {
		Connect    []func(*Connection) error
		Disconnect []func(*Connection)
	}
}

func New(transport Transport, protocol protocol.Protocol) *Server {
	return &Server{
		protocol:  protocol,
		transport: transport,
		Events: struct {
			Connect    []func(*Connection) error
			Disconnect []func(*Connection)
		}{
			[]func(*Connection) error{},
			[]func(*Connection){},
		},
	}
}

func (this *Server) Listen(host string) error {
	return this.transport.Listen(host, func(raw io.ReadWriteCloser) {
		conn := Accept(this.protocol, raw)

		for _, cb := range this.Events.Connect {
			err := cb(conn)
			if err != nil {
				conn.Close()
				return
			}
		}
		defer func() {
			for _, cb := range this.Events.Disconnect {
				cb(conn)
			}
		}()
		conn.Read()
	})
}
