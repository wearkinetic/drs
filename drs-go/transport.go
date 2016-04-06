package drs

import "io"

type Transport interface {
	Listen(string, func(raw io.ReadWriteCloser)) error
	Connect(host string) (io.ReadWriteCloser, error)
}

type RouterHandler func(action string) ([]string, error)
type ByteWriter func(data []byte) error
