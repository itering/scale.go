package types_test

import (
	"encoding/json"
	"fmt"
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

func TestRawBabePreDigest(t *testing.T) {
	raw := "0x02020000008b86750900000000"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("RawBabePreDigest")
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

// 0x025ed0b2 Compact<Balance>
func TestCompactBalance(t *testing.T) {
	raw := "0x025ed0b2"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	m.ProcessAndUpdateData("Compact<Balance>")

}

//
func TestRegistration(t *testing.T) {
	raw := "0x04010000000200a0724e180900000000000000000000000d505552455354414b452d30310e507572655374616b65204c74641b68747470733a2f2f7777772e707572657374616b652e636f6d2f000000000d40707572657374616b65636f"
	m := types.ScaleDecoder{}
	m.Init(types.ScaleBytes{Data: utiles.HexToBytes(raw)}, nil)
	r := m.ProcessAndUpdateData("Registration<BalanceOf>")
	rb, _ := json.Marshal(r)
	fmt.Println(string(rb))
}
