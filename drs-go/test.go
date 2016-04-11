package drs

import (
	"log"
	"testing"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
)

func Test(t *testing.T, transport Transport) {
	pipe := New(transport)
	pipe.Protocol = protocol.JSON
	pipe.Router = func(string) ([]string, error) {
		return []string{"localhost:12000"}, nil
	}
	pipe.On("echo", func(msg *Message) (interface{}, error) {
		log.Println("Got Request", msg.Command.Body)
		return msg.Command.Body, nil
	})
	go pipe.Listen(":12000")
	log.Println("Sending...")
	result, err := pipe.Request(&Command{
		Action: "echo",
		Body:   "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Got Response", result)
}

func TestConnection(t *testing.T, transport Transport) {
	conn := NewConnection()
	conn.Dial(protocol.JSON, transport, "localhost:12000")
	time.Sleep(1 * time.Minute)
	conn.Close()
}
