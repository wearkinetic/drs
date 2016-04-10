package drs

import (
	"testing"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
	"github.com/ironbay/dynamic"
)

func TestConn(t *testing.T) {
	conn := NewConnection2()
	transport := ws.New(dynamic.Empty())
	go conn.Dial(protocol.JSON, transport, "delta.wearkinetic.com")
	count := 0
	for {
		conn.Fire(&Command{
			Action: "drs.ping",
			Body:   count,
		})
		count++
		if count > 1000 {
			break
		}
	}
	conn.Close()
}
