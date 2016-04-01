package auth

import "github.com/ironbay/drs/drs-go"
import "golang.org/x/net/websocket"
import "io"
import "errors"

func Attach(server *drs.Server, cb func(string) (string, error)) {
	server.On(
		"auth.upgrade",
		func(cmd *drs.Command, conn *drs.Connection, ctx map[string]interface{}) (interface{}, error) {
			token, ok := cmd.Body.(string)
			if !ok {
				return nil, drs.Error("Token must be a string")
			}
			user, err := cb(token)
			if err != nil {
				return nil, drs.Error(err.Error())
			}
			conn.Cache.Set("user", user)
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
		return nil
	})
}

func Validator(cmd *drs.Command, conn *drs.Connection, ctx map[string]interface{}) (interface{}, error) {
	user, ok := conn.Cache.Get("user")
	if ok {
		ctx["user"] = user
		return nil, nil
	}
	return nil, errors.New("Authentication required")
}
