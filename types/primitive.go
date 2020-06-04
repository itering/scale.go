package types

import (
	"bytes"
	"encoding/binary"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"io"
	"unicode/utf8"
)

type Compact struct {
	ScaleDecoder
	CompactLength int    `json:"compact_length"`
	CompactBytes  []byte `json:"compact_bytes"`
}

func (c *Compact) ProcessCompactBytes() []byte {
	compactByte := c.NextBytes(1)
	var byteMod = 0
	if len(compactByte) != 0 {
		byteMod = int(compactByte[0]) % 4
	}
	if byteMod == 0 {
		c.CompactLength = 1
	} else if byteMod == 1 {
		c.CompactLength = 2
	} else if byteMod == 2 {
		c.CompactLength = 4
	} else {
		c.CompactLength = 5 + ((int(compactByte[0]) - 3) / 4)
	}
	if c.CompactLength == 1 {
		c.CompactBytes = compactByte
	} else if utiles.IntInSlice(c.CompactLength, []int{2, 4}) {
		c.CompactBytes = append(compactByte[:], c.NextBytes(c.CompactLength - 1)[:]...)
	} else {
		c.CompactBytes = c.NextBytes(c.CompactLength - 1)
	}
	return c.CompactBytes
}

func (c *Compact) Process() {
	c.ProcessCompactBytes()
	if c.SubType != "" {
		s := ScaleDecoder{TypeString: c.SubType, Data: ScaleBytes{Data: c.CompactBytes}}
		byteData := s.ProcessAndUpdateData(c.SubType)
		if c.CompactLength <= 4 {
			switch byteData.(type) {
			case uint64:
				c.Value = uint64(byteData.(uint64) / 4)
			case uint32:
				c.Value = uint64(byteData.(uint32) / 4)
			case int:
				c.Value = uint64(byteData.(int) / 4)
			case decimal.Decimal:
				c.Value = byteData.(decimal.Decimal).Div(decimal.New(4, 0)).Floor()
			default:
				c.Value = byteData
			}
		} else {
			c.Value = byteData
		}
	} else {
		c.Value = c.CompactBytes
	}
}

type CompactU32 struct {
	Compact
	Reader io.Reader
}

func (c *CompactU32) Init(data ScaleBytes, option *ScaleDecoderOption) {
	c.TypeString = "Compact<u32>"
	c.ScaleDecoder.Init(data, option)
}

func (c *CompactU32) Process() {
	c.ProcessCompactBytes()
	buf := &bytes.Buffer{}
	c.Reader = buf
	_, _ = buf.Write(c.CompactBytes)
	b := make([]byte, 8)
	_, _ = c.Reader.Read(b)
	c.Value = int(binary.LittleEndian.Uint64(b))
	if c.CompactLength <= 4 {
		c.Value = int(c.Value.(int)) / 4
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
	b.Value = utiles.BytesToHex(value)
}

type String struct {
	ScaleDecoder
}

func (b *String) Process() {
	length := b.ProcessAndUpdateData("Compact<u32>").(int)
	value := b.NextBytes(int(length))
	if utf8.Valid(value) {
		b.Value = string(value)
	} else {
		b.Value = utiles.BytesToHex(value)
	}
}
