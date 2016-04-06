package drs

import "io"

type Transport interface {
	Listen(func(raw io.ReadWriteCloser)) error
	Connect(host string) (io.ReadWriteCloser, error)
}

type RouterHandler func(action string) ([]string, error)
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
