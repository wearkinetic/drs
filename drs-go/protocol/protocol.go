package protocol

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
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

var XML = func(rw io.ReadWriteCloser) *Stream {
	return &Stream{
		xml.NewEncoder(rw),
		xml.NewDecoder(rw),
	}
}
