package ws

import (
	"testing"

	"github.com/ironbay/drs/drs-go"
)

func TestTransport(t *testing.T) {
	result := New(map[string]interface{}{"token": "djkhaled"})
	drs.TestConnection(t, result)
}
