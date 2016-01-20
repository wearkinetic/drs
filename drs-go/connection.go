package drs

import (
	"encoding/json"
	"io"
)

type Connection struct {
	cache map[string]interface{}
	rw    io.ReadWriteCloser
}

func (this *Connection) Set(key string, value interface{}) {
	this.cache[key] = value
}

func (this *Connection) Get(key string) (interface{}, bool) {
	result, ok := this.cache[key]
	return result, ok
}

func (this *Connection) Send(cmd *Command) error {
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}
	_, err = this.rw.Write(data)
	return err
}

func NewConnection(rw io.ReadWriteCloser) *Connection {
	return &Connection{
		cache: map[string]interface{}{},
		rw:    rw,
	}
}
