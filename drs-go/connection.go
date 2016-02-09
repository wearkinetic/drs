package drs

import (
	"io"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Connection struct {
	*protocol.Stream
	cache map[string]interface{}
	Raw   io.ReadWriteCloser
}

func (this *Connection) Set(key string, value interface{}) {
	this.cache[key] = value
}

func (this *Connection) Get(key string) interface{} {
	return this.cache[key]
}

func (this *Pipe) dial(host string) (*Connection, error) {
	conn, ok := this.connections[host]
	if ok {
		return conn, nil
	}
	rw, err := this.transport.Connect(host)
	if err != nil {
		return nil, err
	}
	conn = this.connect(rw)
	this.connections[host] = conn
	go func() {
		this.handle(conn)
		delete(this.connections, host)
	}()
	return conn, nil
}

func (this *Pipe) connect(rw io.ReadWriteCloser) *Connection {
	conn := &Connection{
		Stream: this.Protocol(rw),
		cache:  map[string]interface{}{},
		Raw:    rw,
	}
	return conn
}

func (this *Pipe) handle(conn *Connection) {
	for {
		cmd := new(Command)
		err := conn.Decode(&cmd)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			break
		}
		go this.process(conn, cmd)
	}
}
