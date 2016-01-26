package drs

import (
	"io"
	"log"

	"github.com/ironbay/drs/drs-go/protocol"
)

type Connection struct {
	protocol *protocol.Stream
	cache    map[string]interface{}
	rw       io.ReadWriteCloser
}

func (this *Connection) Set(key string, value interface{}) {
	this.cache[key] = value
}

func (this *Connection) Get(key string) (interface{}, bool) {
	result, ok := this.cache[key]
	return result, ok
}

func (this *DRS) connect(host string) (*Connection, error) {
	conn, ok := this.connections[host]
	if ok {
		return conn, nil
	}
	rw, err := this.transport.Connect(host)
	if err != nil {
		return nil, err
	}
	conn, _ = this.register(host, rw)
	return conn, nil
}

func (this *DRS) register(host string, rw io.ReadWriteCloser) (*Connection, chan bool) {
	conn := &Connection{
		protocol: this.Protocol(rw),
		cache:    map[string]interface{}{},
		rw:       rw,
	}
	done := make(chan bool)
	this.connections[host] = conn
	log.Println("Connected to", host)
	go func() {
		for {
			cmd := new(Command)
			err := conn.protocol.Decode(&cmd)
			if err != nil && err.Error() == "EOF" {
				log.Println(err)
				break
			}
			this.Process(conn, cmd)
		}
		delete(this.connections, host)
		done <- true
	}()
	return conn, done
}
