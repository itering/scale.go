package types

import (
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
