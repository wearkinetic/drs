package ws

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ironbay/drs/drs-go"
	"golang.org/x/net/websocket"
)

type Transport struct {
	query map[string]interface{}
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(ch drs.ConnectionHandler) error {
	ws := websocket.Server{
		Handler: websocket.Handler(func(ws *websocket.Conn) {
			ch(ws)
		}),
	}
	http.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		ws.ServeHTTP(w, req)
	})
	return http.ListenAndServe(":12000", nil)
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
