package drs

import (
	"log"
	"testing"
	"time"

	"github.com/ironbay/drs/drs-go/protocol"
	"github.com/ironbay/drs/drs-go/transports/ws"
	"github.com/ironbay/dynamic"
	"github.com/ironbay/go-util/console"
)

func TestConn(t *testing.T) {
	conn := NewConnection()
	transport := ws.New(dynamic.Build("token", "djkhaled"))
	go conn.Dial(protocol.JSON, transport, "localhost:12000", true)
	count := 0
	for {
		go func() {
			now := time.Now()
			result, _ := conn.Request(&Command{
				Action: "drs.ping",
				Body:   count,
			})
			log.Println(time.Since(now).Seconds() * 1000)
			console.JSON(result)
		}()
		count++
		if count > 1000 {
			break
		}
	}
	log.Println("Sleeping")
	time.Sleep(1 * time.Minute)
	conn.Close()
}
