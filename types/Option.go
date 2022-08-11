package types

import "github.com/itering/scale.go/utiles"

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

func (o *Option) Encode(value interface{}) {
	if utiles.IsNil(value) {
		o.Value = "00"
		return
	}
	if o.SubType == "bool" {
		if value.(bool) {
			o.Value = "01"
		} else {
			o.Value = "02"
		}
		return
	}
	o.Value = Encode(o.SubType, value)
}
