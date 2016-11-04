package ws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

type Transport struct {
	query map[string]interface{}
}

func (this *Transport) On(action string) error {
	return nil
}

var u = websocket.Upgrader{}

func (this *Transport) Listen(host string, cb func(raw io.ReadWriteCloser)) error {
	http.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		conn, err := u.Upgrade(w, req, nil)
		if err != nil {
			response(w, 500, err)
			return
		}
		cb(conn.UnderlyingConn())
	})
	return http.ListenAndServe(host, nil)
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	query := ""
	for key, value := range this.query {
		query += fmt.Sprintf("%v=%v&", key, value)
	}
	ws, _, err := websocket.DefaultDialer.Dial("ws://"+host+"/socket?"+query, nil)
	if err != nil {
		return nil, err
	}
	return ws.UnderlyingConn(), nil
}

func New(query map[string]interface{}) *Transport {
	transport := new(Transport)
	transport.query = query
	return transport
}

func response(w http.ResponseWriter, status int, input interface{}) {
	data, _ := json.Marshal(input)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
