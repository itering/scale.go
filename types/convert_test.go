package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertType(t *testing.T) {
	assert.Equal(t, ConvertType("&'static[u8]"), "Bytes")
}
