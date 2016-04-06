package drs

import "io"

type Transport interface {
	Listen(host string, ch ConnectionHandler) error
	Connect(host string) (io.ReadWriteCloser, error)
}

type RouterHandler func(action string) ([]string, error)
type ConnectionHandler func(rw io.ReadWriteCloser)
type ByteWriter func(data []byte) error

/*

transport := websocket.New()

transporter.On("myqction", (conn, replies) => {
	replies <- someCommand1
	replies <- someCommand2
	close(replies)
})

transporter.Send(myCommand)

transporter.Listen()

*/
