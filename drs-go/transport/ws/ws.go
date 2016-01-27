package ws

import (
	"io"
	"net/http"

	"github.com/ironbay/drs/drs-go"
	"golang.org/x/net/websocket"
)

type Transport struct {
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(ch drs.ConnectionHandler) error {
	http.HandleFunc("/socket", func(w http.ResponseWriter, req *http.Request) {
		s := websocket.Server{
			Handler: websocket.Handler(func(ws *websocket.Conn) {
				ch(ws)
			}),
		}
		s.ServeHTTP(w, req)
	})
	return http.ListenAndServe(":12000", nil)
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	ws, err := websocket.Dial("ws://"+host+":12000/socket", "", "http://"+host)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func New() (*drs.Pipe, error) {
	transport := new(Transport)
	return drs.New(transport)
}
