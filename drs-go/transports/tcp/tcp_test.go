package tcp

import (
	"testing"

	"github.com/ironbay/drs/drs-go"
)

func TestTransport(t *testing.T) {
	result := New()
	drs.TestConnection(t, result)
}
