package ping

import (
	"time"

	"github.com/wearkinetic/drs/drs-go"
)

func Attach(processor *drs.Processor) {
	processor.On(
		"drs.ping",
		func(msg *drs.Message) (interface{}, error) {
			return time.Now().UnixNano() / int64(time.Millisecond), nil
		},
	)
}
