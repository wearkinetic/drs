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
	pipe.On("echo", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
		log.Println("Got Request", cmd.Body)
		return cmd.Body, nil
	})
	go pipe.Listen()
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
	conn := NewConnection(protocol.JSON)
	conn.Dial(transport, "localhost:12000", false)
	time.Sleep(1 * time.Minute)
	conn.Close()
}
