package drs

import (
	"testing"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
)

func TestConnection(t *testing.T) {
	ws := ws.New()
	conn := NewConnection(protocol.JSON)
	conn.Dial(ws, "localhost:12000")
}
