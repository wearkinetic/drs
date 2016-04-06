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
	Raw io.ReadWriteCloser
	Encoder
	Decoder
}

func (this *Stream) Close() {
	this.Raw.Close()
}

type Protocol func(raw io.ReadWriteCloser) *Stream

var JSON = func(raw io.ReadWriteCloser) *Stream {
	dc := json.NewDecoder(raw)
	dc.UseNumber()
	return &Stream{
		raw,
		json.NewEncoder(raw),
		dc,
	}
}

var GOB = func(raw io.ReadWriteCloser) *Stream {
	return &Stream{
		raw,
		gob.NewEncoder(raw),
		gob.NewDecoder(raw),
	}
}

var XML = func(raw io.ReadWriteCloser) *Stream {
	return &Stream{
		raw,
		xml.NewEncoder(raw),
		xml.NewDecoder(raw),
	}
}
