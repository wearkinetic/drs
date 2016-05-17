package drs

import (
	"log"
	"testing"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
)

func TestConnection(t *testing.T) {
	transport := ws.New(map[string]interface{}{"token": "djkhaled"})
	conn, err := Dial(protocol.JSON, transport, "delta.inboxtheapp.com")
	log.Println("Connected")
	if err != nil {
		t.Fatal(err)
	}
	go conn.Read()
	result, err := conn.Call(&Command{
		Action: "drs.ping",
	})
	log.Println(result)
}
