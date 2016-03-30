package auth

import "github.com/ironbay/drs/drs-go"

func Attach(processor *drs.Processor, cb func(string) (string, error)) {
	processor.On(
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
}

func Validator(cmd *drs.Command, conn *drs.Connection, ctx map[string]interface{}) (interface{}, error) {
	user, ok := conn.Cache.Get("user")
	if !ok {
		return nil, drs.Error("Requires authentication")
	}
	ctx["user"] = user
	return nil, nil
}
