package ws

import (
	"testing"

	"github.com/ironbay/drs/drs-go"
)

func TestTransport(t *testing.T) {
	result := New(map[string]interface{}{})
	drs.Test(t, result)
}
