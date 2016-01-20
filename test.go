package drs

import (
	"log"
	"testing"
)

func Test(t *testing.T, pipe *DRS) {
	pipe.Router = func(string) (string, error) {
		return "localhost", nil
	}
	pipe.On("echo", func(cmd *Command, conn *Connection) (interface{}, error) {
		log.Println("Got Request", cmd.Body)
		return cmd.Body, nil
	})
	err := pipe.Listen()
	if err != nil {
		t.Fatal(err)
	}
	result, err := pipe.Send(&Command{
		Action: "echo",
		Body:   "hello",
	})
	if err != nil {
		t.Fatal(err)
	}
	log.Println("Got Response", result)
}
