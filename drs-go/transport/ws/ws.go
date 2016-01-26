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
	http.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		request := ws.Request()
		ip := request.RemoteAddr
		headers := request.Header["X-Forwarded-For"]
		if len(headers) > 0 {
			ip = headers[0]
		}
		_, done := ch(ip, ws)
		<-done
	}))
	return http.ListenAndServe(":12000", nil)
}

func (this *Transport) Connect(host string) (io.ReadWriteCloser, error) {
	ws, err := websocket.Dial("ws://"+host+":12000", "", "http://"+host)
	if err != nil {
		return nil, err
	}
	return ws, nil
}

func New() (*drs.Pipe, error) {
	transport := new(Transport)
	return drs.New(transport)
}
