package drs

/*
func TestConnection(t *testing.T) {
	transport := ws.New(map[string]interface{}{"token": "djkhaled"})
	server := NewServer(transport)
	go server.Listen()
	conn := NewConnection(protocol.JSON)
	go conn.Dial(transport, "localhost:12000", true)
	result, _ := conn.Request(&Command{
		Action: "drs.ping",
	})
	log.Println("Pinged", result)
}
*/
