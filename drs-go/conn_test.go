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
	actor.Supervise(func(s *actor.Session) {
		transport := ws.New(dynamic.Build("token", "djkhaled"))
		conn := NewConnection()
		if err := conn.Dial(protocol.JSON, transport, "localhost:12000"); err != nil {
			time.Sleep(1 * time.Second)
			s.Stop <- err
			return
		}
		s.Cleanup(conn.Close)
		conn.OnDisconnect = func(err error) {
			s.Stop <- err
		}

		for i := 0; i < 3; i++ {
			res, err := conn.Request(&Command{
				Action: "drs.ping",
			})
			if err != nil {
				s.Stop <- err
				return
			}
			log.Println(res)
			time.Sleep(1 * time.Second)
		}

		s.Stop <- nil
	})
}
