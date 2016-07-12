package tcp

import (
	"io"
	"net"
)

type Transport struct {
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(host string, ch func(raw io.ReadWriteCloser)) error {
	ln, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go func() {
			ch(conn)
			conn.Close()
		}()
	}
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func New() *Transport {
	transport := new(Transport)
	return transport
}
