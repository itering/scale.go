package types

import (
	"github.com/itering/scale.go/utiles"
)

type Option struct {
	ScaleDecoder
}

func (o *Option) Process() {
	optionFlag := o.NextBytes(1)
	if o.SubType != "" && utiles.BytesToHex(optionFlag) != "00" {
		if o.SubType == "bool" {
			o.Value = utiles.BytesToHex(optionFlag) == "01"
			return
		}
		o.Value = o.ProcessAndUpdateData(o.SubType)
	}
}

func (o *Option) Encode(value interface{}) string {
	if v, ok := value.(string); ok && v == "" {
		return "00"
	}
	if utiles.IsNil(value) {
		return "00"
	}
	if o.SubType == "bool" {
		if value.(bool) {
			return "01"
		}
		return "02"
	}
	return EncodeWithOpt(o.SubType, value, &ScaleDecoderOption{Spec: o.Spec, Metadata: o.Metadata})
}
