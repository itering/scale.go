package types_test

import (
	"github.com/freehere107/go-scale-codec/source"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/freehere107/go-scale-codec/utiles"
	"strings"
	"testing"
)

func TestCompactU64(t *testing.T) {
	raw := "10"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Compact<U64>").(uint64)
	if r != 4 {
		t.Errorf("Test TestCompactU64 Process fail, decode return %d", r)
	}
}

func TestEnum_Process(t *testing.T) {
	types.RuntimeType{}.Reg()
	types.RegCustomTypes(source.LoadTypeRegistry("../source/base"))

}

func TestSet_Process(t *testing.T) {
	types.RuntimeType{}.Reg()
	types.RegCustomTypes(map[string]source.TypeStruct{
		"CustomSet": {
			Type:      "set",
			ValueList: []string{"Value1", "Value2", "Value3", "Value4", "Value5"},
		},
	})
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes("0x03000000")}, nil)
	r := m.ProcessAndUpdateData("CustomSet")
	if strings.Join(r.([]string), "") != "Value1Value2" {
		t.Errorf("Test TestSet_Process Process fail, decode return %v", r.([]string))
	}
}
