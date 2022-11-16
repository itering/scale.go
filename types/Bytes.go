package types

import (
	"github.com/itering/scale.go/utiles"
)

type Bytes struct{ ScaleDecoder }

func (b *Bytes) Process() {
	length := b.ProcessAndUpdateData("Compact<u32>").(int)
	if value := b.NextBytes(length); utiles.IsASCII(value) {
		b.Value = string(value)
	} else {
		b.Value = utiles.AddHex(utiles.BytesToHex(value))
	}
}

func (b *Bytes) Encode(value string) string {
	value = utiles.TrimHex(value)
	if len(value)%2 == 1 {
		value += "0"
	}
	bytes := utiles.HexToBytes(value)
	return Encode("Compact<u32>", len(bytes)) + value
}

type HexBytes struct{ ScaleDecoder }

func (h *HexBytes) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(h.ProcessAndUpdateData("Compact<u32>").(int))))
}

func (h *HexBytes) Encode(value string) string {
	value = utiles.TrimHex(value)
	if len(value)%2 == 1 {
		value += "0"
	}
	bytes := utiles.HexToBytes(value)
	return Encode("Compact<u32>", len(bytes)) + value
}

type String struct{ Bytes }

func (s *String) Encode(value string) string {
	bytes := []byte(value)
	return Encode("Compact<u32>", len(bytes)) + utiles.BytesToHex(bytes)
}
