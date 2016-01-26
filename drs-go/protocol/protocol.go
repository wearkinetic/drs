package protocol

import (
	"encoding/gob"
	"encoding/json"
	"io"
)

type Encoder interface {
	Encode(v interface{}) error
}

type Decoder interface {
	Decode(v interface{}) error
}

type Stream struct {
	Encoder
	Decoder
}

type Protocol func(rw io.ReadWriteCloser) *Stream

var JSON = func(rw io.ReadWriteCloser) *Stream {
	return &Stream{
		json.NewEncoder(rw),
		json.NewDecoder(rw),
	}
}

var GOB = func(rw io.ReadWriteCloser) *Stream {
	return &Stream{
		gob.NewEncoder(rw),
		gob.NewDecoder(rw),
	}
}
