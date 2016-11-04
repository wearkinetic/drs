package ws

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Transport struct {
	query map[string]interface{}
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(host string, ch func(raw io.ReadWriteCloser)) error {
	ws := websocket.Handler(func(w *websocket.Conn) {
		ch(w)
	})
	http.HandleFunc("/socket", ws.ServeHTTP)
	return http.ListenAndServe(host, nil)
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	query := ""
	for key, value := range this.query {
		query += fmt.Sprintf("%v=%v&", key, value)
	}
	ws, err := websocket.Dial("ws://"+host+"/socket?"+query, "", "http://"+host)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func New(query map[string]interface{}) *Transport {
	transport := new(Transport)
	transport.query = query
	return transport
}
