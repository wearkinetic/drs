package drs

import "io"

type Transport interface {
	On(action string) error
	Listen(ch ConnectionHandler) error
	Connect(host string) (io.ReadWriteCloser, error)
}

type Command struct {
	Key    string      `json:"key"`
	Action string      `json:"action"`
	Body   interface{} `json:"body"`
	Meta   Dynamic     `json:"meta"`
}

func (this *Command) Dynamic() Dynamic {
	return Dynamic(this.Body.(map[string]interface{}))
}

type CommandHandler func(cmd *Command, conn *Connection, context Dynamic) (interface{}, error)
type RouterHandler func(action string) (string, error)
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
