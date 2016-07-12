package ipc

import (
	"errors"
	"io"
	"os"
)

type Transport struct {
	io.ReadCloser
	io.WriteCloser
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(host string, ch func(raw io.ReadWriteCloser)) error {
	return errors.New("Listen: not supported for ipc")
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	return this, nil
}

func New() *Transport {
	return &Transport{
		os.Stdin,
		os.Stdout,
	}
}

func (this *Transport) Close() error {
	if err := this.ReadCloser.Close(); err != nil {
		return err
	}
	if err := this.WriteCloser.Close(); err != nil {
		return err
	}
	return nil
}
