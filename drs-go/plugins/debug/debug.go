package debug

import (
	"time"

	"github.com/ironbay/drs/drs-go"
)

func Attach(processor *drs.Processor) {
	processor.On(
		"drs.kaboom",
		func(msg *drs.Message) (interface{}, error) {
			return nil, drs.Exception(rune(0x1F4A5))
		},
	)
}
