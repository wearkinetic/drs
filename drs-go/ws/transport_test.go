package ws

import (
	"testing"

	"github.com/ironbay/drs/drs-go"
)

func TestTransport(t *testing.T) {
	result, err := New()
	if err != nil {
		t.Fatal(err)
	}
	drs.Test(t, result)
}
