package debug

import (
	"github.com/ironbay/drs/drs-go"
)

const KABOOM = rune(0x1F4A5)

func Attach(processor *drs.Processor) {
	processor.On(
		"drs.kaboom",
		func(msg *drs.Message) (interface{}, error) {
			return nil, drs.Exception(string(KABOOM))
		},
	)
}
