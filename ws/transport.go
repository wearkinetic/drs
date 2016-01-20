package ws

import (
	"io"
	"net/http"

	"github.com/ironbay/drs-go"
	"golang.org/x/net/websocket"
)

type Transport struct {
}

func (this *Transport) On(action string) error {
	return nil
}

func (this *Transport) Listen(ch drs.ConnectionHandler) error {
	http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		done, _ := ch(ws)
		<-done
	}))
	go http.ListenAndServe(":12000", nil)
	return nil
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	ws, err := websocket.Dial("ws://"+host+":12000", "", "http://"+host)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func (this *Transport) Frame(rw io.ReadWriteCloser) ([]byte, error) {
	ws := rw.(*websocket.Conn)
	var data []byte
	err := websocket.Message.Receive(ws, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func New() (*drs.DRS, error) {
	transport := new(Transport)
	return drs.New(transport)
}
