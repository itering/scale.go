package convert

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertType(t *testing.T) {
	assert.Equal(t, ConvertType("&'static[u8]"), "Bytes")
}
