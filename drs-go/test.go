package drs

import (
	"log"
	"testing"

	"github.com/ironbay/drs/drs-go/protocol"
)

func Test(t *testing.T, pipe *Pipe) {
	pipe.Protocol = protocol.GOB
	pipe.Router = func(string) (string, error) {
		return "localhost", nil
	}
	pipe.On("echo", func(cmd *Command, conn *Connection, ctx map[string]interface{}) (interface{}, error) {
		log.Println("Got Request", cmd.Body)
		return cmd.Body, nil
	})
	go pipe.Listen()
	result, err := pipe.Send(&Command{
		Action: "echo",
		Body:   "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Got Response", result)
}
