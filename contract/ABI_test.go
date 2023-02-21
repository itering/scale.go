package contract

import (
	"io/ioutil"
	"testing"

	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/stretchr/testify/assert"
)

func Test_AbiParse(t *testing.T) {
	c, err := ioutil.ReadFile("metadata.json")
	assert.NoError(t, err)

	abi, err := InitAbi(c)
	assert.NoError(t, err)
	assert.Greater(t, len(abi.Types), 1)
	utiles.Debug(abi)

	sc := types.ScaleDecoder{DuplicateName: make(map[string]int), RegisteredSiType: make(map[int]string)}
	abi.Register(&sc, "pre")

	utiles.Debug(abi.RegisteredSiType)
	assert.Equal(t, len(abi.RegisteredSiType), 8)

}
