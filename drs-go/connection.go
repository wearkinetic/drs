package drs

import (
	"io"
	"log"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Connection struct {
	protocol *protocol.Stream
	cache    map[string]interface{}
	Raw      io.ReadWriteCloser
}

func (this *Connection) Set(key string, value interface{}) {
	this.cache[key] = value
}

func (this *Connection) Get(key string) (interface{}, bool) {
	result, ok := this.cache[key]
	return result, ok
}

func (this *Pipe) connect(host string) (*Connection, error) {
	conn, ok := this.connections[host]
	if ok {
		return conn, nil
	}
	rw, err := this.transport.Connect(host)
	if err != nil {
		return nil, err
	}
	conn = this.newConnection(rw)
	this.connections[host] = conn
	go func() {
		this.handle(conn)
		delete(this.connections, host)
	}()
	return conn, nil
}

func (this *Pipe) newConnection(rw io.ReadWriteCloser) *Connection {
	conn := &Connection{
		protocol: this.Protocol(rw),
		cache:    map[string]interface{}{},
		Raw:      rw,
	}
	return conn
}

func (this *Pipe) handle(conn *Connection) {
	for {
		cmd := new(Command)
		err := conn.protocol.Decode(&cmd)
		if err != nil && err.Error() == "EOF" {
			log.Println(err)
			break
		}
		this.Process(conn, cmd)
	}
}
