package drs

import (
	"log"
	"testing"
	"time"

	"github.com/wearkinetic/drs/drs-go/protocol"
	"github.com/wearkinetic/drs/drs-go/transports/ws"
	"github.com/wearkinetic/dynamic"
)

func TestConnection(t *testing.T) {
	server := New(ws.New(dynamic.Empty()), protocol.JSON)
	server.On("drs.ping", func(msg *Message) (interface{}, error) {
		return time.Now().UnixNano() / int64(time.Millisecond), nil
	})
	go server.Listen(":12000")

	transport := ws.New(map[string]interface{}{"token": "djkhaled"})
	conn, err := Dial(protocol.JSON, transport, "localhost:12000")
	log.Println("Connected")
	if err != nil {
		t.Fatal(err)
	}
	go conn.Read()

	result, err := conn.Call(&Command{
		Action: "drs.ping",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
}
