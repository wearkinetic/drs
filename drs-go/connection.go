package drs

import (
	"encoding/json"
	"io"
)

type Connection struct {
	Encoder
	Decoder
	cache   map[string]interface{}
	rw      io.ReadWriteCloser
	encoder Encoder
	decoder Decoder
}

func (this *Connection) Set(key string, value interface{}) {
	this.cache[key] = value
}

func (this *Connection) Get(key string) (interface{}, bool) {
	result, ok := this.cache[key]
	return result, ok
}

func NewConnection(rw io.ReadWriteCloser) *Connection {
	return &Connection{
		Encoder: json.NewEncoder(rw),
		Decoder: json.NewDecoder(rw),
		cache:   map[string]interface{}{},
		rw:      rw,
	}
}
