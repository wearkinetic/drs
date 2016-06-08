package ipc

import (
	"io"
	"net/http"
	"os"

	"golang.org/x/net/websocket"
)

type Transport struct {
	io.ReadCloser
	io.WriteCloser
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(host string, ch func(raw io.ReadWriteCloser)) error {
	ws := websocket.Server{
		Handler: websocket.Handler(func(ws *websocket.Conn) {
			ch(ws)
		}),
	}
	http.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		ws.ServeHTTP(w, req)
	})
	return http.ListenAndServe(host, nil)
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
