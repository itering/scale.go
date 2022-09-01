package types

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
)

type FixedArray struct {
	ScaleDecoder
	FixedLength int
	SubType     string
}

func (f *FixedArray) Init(data scaleBytes.ScaleBytes, option *ScaleDecoderOption) {
	if option != nil && option.FixedLength != 0 {
		f.FixedLength = option.FixedLength
	}
	f.ScaleDecoder.Init(data, option)
}

func (f *FixedArray) Process() {
	var result []interface{}
	if f.FixedLength > 0 {
		if strings.EqualFold(f.SubType, "u8") {
			value := f.NextBytes(f.FixedLength)
			if utiles.IsASCII(value) {
				f.Value = string(value)
			} else {
				f.Value = utiles.BytesToHex(value)
			}
			return
		}
		for i := 0; i < f.FixedLength; i++ {
			result = append(result, f.ProcessAndUpdateData(f.SubType))
		}
		f.Value = result
	} else {
		f.GetNextU8()
	}
}

func (f *FixedArray) Encode(value interface{}) string {
	var raw string
	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(value)
		if s.Len() != f.FixedLength {
			panic("fixed length not match")
		}
		subType := f.SubType
		for i := 0; i < s.Len(); i++ {
			raw += Encode(subType, s.Index(i).Interface())
		}
		return raw
	default:
		panic(fmt.Errorf("invalid vec input"))
	}
}
