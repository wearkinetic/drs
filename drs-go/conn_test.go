package drs

import (
	"log"
	"testing"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
	"github.com/ironbay/dynamic"
	"github.com/ironbay/go-util/actor"
)

func TestConn(t *testing.T) {
	actor.Supervise(func() error {
		transport := ws.New(dynamic.Build("token", "djkhaled"))
		conn := NewConnection()
		if err := conn.Dial(protocol.JSON, transport, "localhost:12000"); err != nil {
			time.Sleep(1 * time.Second)
			return err
		}
		defer conn.Close()

		for i := 0; i < 3; i++ {
			res, err := conn.Request(&Command{
				Action: "drs.ping",
			})
			if err != nil {
				return err
			}
			log.Println(res)
			time.Sleep(1 * time.Second)
		}
		return <-conn.Done
	})
}
