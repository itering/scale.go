package types_test

import (
	"github.com/freehere107/go-scale-codec/types"
	"github.com/freehere107/go-scale-codec/utiles"
	"testing"
)

func TestCompactU64(t *testing.T) {
	raw := "10"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Compact<U64>").(int)
	if r != 4 {
		t.Errorf("Test TestCompactU64 Process fail, decode return %d", r)
	}
}
