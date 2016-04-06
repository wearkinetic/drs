package drs

import (
	"testing"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
	"github.com/ironbay/dynamic"
)

func TestConn(t *testing.T) {
	conn := NewConnection2()
	transport := ws.New(dynamic.Empty())
	go conn.Dial(protocol.JSON, transport, "localhost:12000")
	/*
		count := 0
			for {
				conn.Fire(&Command{
					Action: "drs.ping",
				})
				count++
				if count > 5 {
					break
				}
			}
	*/
	time.Sleep(5 * time.Second)
	conn.Close()
}
