package ping

import (
	"time"

	"github.com/ironbay/drs/drs-go"
)

func Attach(processor *drs.Processor) {
	processor.On(
		"ping",
		func(cmd *drs.Command, conn *drs.Connection, ctx map[string]interface{}) (interface{}, error) {
			return time.Now().UnixNano() / int64(time.Millisecond), nil
		},
	)
}
