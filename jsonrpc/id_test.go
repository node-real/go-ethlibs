package jsonrpc

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func Test_id_null(t *testing.T) {
	id := StringID("")
	json, err := id.MarshalJSON()
	assert.Equal(t, err, nil)
	assert.Equal(t, json, []uint8{110, 117, 108, 108})
}
