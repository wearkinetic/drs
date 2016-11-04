package ws

import (
	"log"
	"testing"
	"time"

	"github.com/ironbay/drs/drs-go"
	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/dynamic"
)

func TestConnection(t *testing.T) {
	server := drs.New(New(dynamic.Empty()), protocol.JSON)
	server.On("drs.ping", func(msg *drs.Message) (interface{}, error) {
		return time.Now().UnixNano() / int64(time.Millisecond), nil
	})
	go server.Listen(":12000")

	transport := New(map[string]interface{}{"token": "djkhaled"})
	conn, err := drs.Dial(protocol.JSON, transport, "localhost:12000")
	log.Println("Connected")
	if err != nil {
		t.Fatal(err)
	}
	go conn.Read()

	result, err := conn.Call(&drs.Command{
		Action: "drs.ping",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
}
