package utiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsASCII(t *testing.T) {
	assert.True(t, IsASCII([]byte("hello world")))
	assert.False(t, IsASCII([]byte{5}))
	assert.False(t, IsASCII(HexToBytes("1c08144900000000000000")))
	assert.False(t, IsASCII(HexToBytes("080b2501")))
	assert.True(t, IsASCII(HexToBytes("0x73686f7274206d656d6f00000000000000000000000000000000000000000000")))
}
