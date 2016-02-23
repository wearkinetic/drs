package drs

import (
	"io"
	"log"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Connection struct {
	*protocol.Stream
	Cache map[string]interface{}
	Raw   io.ReadWriteCloser
	block chan bool
}

func (this *Connection) Set(key string, value interface{}) {
	this.Cache[key] = value
}

func (this *Connection) Get(key string) interface{} {
	return this.Cache[key]
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
		Cache:  map[string]interface{}{},
		Raw:    rw,
		block:  make(chan bool, 1),
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
			log.Println(err)
			break
		}
		go this.process(conn, cmd)
	}
	conn.Raw.Close()
}
