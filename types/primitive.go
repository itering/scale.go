package types

import (
	"encoding/binary"
	"github.com/freehere107/scalecodec/utiles"
	"unicode/utf8"
)

type CompactU32 struct {
	Compact
}

func (c *CompactU32) Init(data ScaleBytes, option *ScaleDecoderOption) {
	c.TypeString = "Compact<u32>"
	c.ScaleDecoder.Init(data, option)
}

func (c *CompactU32) Process() {
	c.ProcessCompactBytes()
	if c.CompactLength <= 4 {
		data := make([]byte, len(c.Data.Data))
		copy(data, c.Data.Data)

		compactBytes := c.CompactBytes
		bs := make([]byte, 4-len(compactBytes))
		compactBytes = append(compactBytes[:], bs...)
		c.Data.Data = data
		c.Value = int(binary.LittleEndian.Uint32(compactBytes)) / 4
	} else {
		c.Value = int(binary.LittleEndian.Uint32(c.CompactBytes))
	}
}

func (c *CompactU32) Encode(value int) ScaleBytes {
	if value <= 63 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2))
		c.Data.Data = bs[0:1]
	} else if value <= 16383 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2)|1)
		c.Data.Data = bs[0:2]
	} else if value <= 1073741823 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2)|2)
		c.Data.Data = bs
	}
	return c.Data
}

type Option struct {
	ScaleDecoder
}

func (o *Option) Process() {
	optionType := o.NextBytes(1)
	if o.SubType != "" && utiles.BytesToHex(optionType) != "00" {
		o.Value = o.ProcessAndUpdateData(o.SubType)
	}
}

type Bytes struct {
	ScaleDecoder
}

func (b *Bytes) Init(data ScaleBytes, option *ScaleDecoderOption) {
	b.TypeString = "Vec<u8>"
	b.ScaleDecoder.Init(data, option)
}

func (b *Bytes) Process() {
	length := b.ProcessAndUpdateData("Compact<u32>").(int)
	value := b.NextBytes(int(length))
	if utf8.Valid(value) {
		b.Value = string(value)
	} else {
		b.Value = utiles.BytesToHex(value)
	}
}
