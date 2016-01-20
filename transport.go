package drs

import "io"

type Transport interface {
	On(action string) error
	Listen(ch ConnectionHandler) error
	Connect(host string) (io.ReadWriteCloser, error)
	Frame(rw io.ReadWriteCloser) ([]byte, error)
}

type Command struct {
	Key    string
	Action string
	Body   interface{}
}

type CommandHandler func(cmd *Command, conn *Connection) (interface{}, error)
type RouterHandler func(action string) (string, error)
type ConnectionHandler func(rw io.ReadWriteCloser) (chan bool, *Connection)
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
