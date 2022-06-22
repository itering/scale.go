package types

import (
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
)

type Bytes struct{ ScaleDecoder }

func (b *Bytes) Init(data scaleBytes.ScaleBytes, option *ScaleDecoderOption) {
	b.TypeString = "Vec<u8>"
	b.ScaleDecoder.Init(data, option)
}

func (b *Bytes) Process() {
	length := b.ProcessAndUpdateData("Compact<u32>").(int)
	if value := b.NextBytes(length); utiles.IsASCII(value) {
		b.Value = string(value)
	} else {
		b.Value = utiles.BytesToHex(value)
	}
}

type HexBytes struct{ ScaleDecoder }

func (h *HexBytes) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(h.ProcessAndUpdateData("Compact<u32>").(int))))
}

type String struct{ Bytes }
