package types

import (
	"strings"

	"github.com/itering/scale.go/utiles"
)

type Result struct {
	ScaleDecoder
}

func (b *Result) Process() {
	optionValue := utiles.BytesToHex(b.NextBytes(1))
	subType := strings.Split(b.SubType, ",")
	if len(subType) != 2 {
		panic("Result subType not illegal")
	}
	if optionValue == "00" || optionValue == "" {
		b.Value = map[string]interface{}{"Ok": b.ProcessAndUpdateData(subType[0])}
	} else if optionValue == "01" {
		b.Value = map[string]interface{}{"Error": b.ProcessAndUpdateData(subType[1])}
	} else {
		panic("illegal Result data")
	}
}

func (b *Result) Encode(value map[string]interface{}) string {
	subType := strings.Split(b.SubType, ",")
	if len(subType) != 2 {
		panic("Result subType not illegal")
	}
	if data, ok := value["Ok"]; ok {
		return "01" + Encode(subType[0], data)
	}
	if data, ok := value["Error"]; ok {
		return "00" + Encode(subType[1], data)
	}
	panic("illegal Result data")
}
