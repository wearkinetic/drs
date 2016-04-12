package auth

import (
	"github.com/ironbay/drs/drs-go"
	"github.com/ironbay/go-util/actor"
)
import "log"
import "golang.org/x/net/websocket"
import "io"
import "errors"

func Attach(server *drs.Server, cb func(string) (string, error)) {
	server.On(
		"auth.upgrade",
		func(msg *drs.Message) (interface{}, error) {
			token, ok := msg.Command.Body.(string)
			if !ok {
				return nil, actor.Error("Token must be a string")
			}
			user, err := cb(token)
			if err != nil {
				return nil, actor.Error(err.Error())
			}
			msg.Connection.Cache.Set("user", user)
			return user, nil
		},
	)

	// LEGACY
	server.OnConnect(func(conn *drs.Connection, raw io.ReadWriteCloser) error {
		ws := raw.(*websocket.Conn)
		query := ws.Request().URL.Query()
		token := query.Get("token")
		user, err := cb(token)
		if err != nil {
			return err
		}
		conn.Cache.Set("user", user)
		log.Println("Legacy Client", user)
		return nil
	})
}

func Validator(msg *drs.Message) (interface{}, error) {
	user, ok := msg.Connection.Cache.Get("user")
	if ok {
		msg.Context["user"] = user
		return nil, nil
	}
	return nil, errors.New("Authentication required")
}
