package tcp

import (
	"testing"

	"github.com/wearkinetic/drs/drs-go"
)

func TestTransport(t *testing.T) {
	result, err := New()
	if err != nil {
		t.Fatal(err)
	}
	drs.Test(t, result)
}
