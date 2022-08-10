package types

import (
	"bytes"
	"encoding/binary"
	"io"

	"github.com/huandu/xstrings"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/scale.go/utiles/uint128"
	"github.com/shopspring/decimal"
)

type U8 struct {
	ScaleDecoder
}

func (u *U8) Process() {
	u.Value = u.GetNextU8()
}

type U16 struct {
	Reader io.Reader
	ScaleDecoder
}

func (u *U16) Process() {
	buf := &bytes.Buffer{}
	u.Reader = buf
	_, _ = buf.Write(u.NextBytes(2))
	c := make([]byte, 2)
	_, _ = u.Reader.Read(c)
	u.Value = binary.LittleEndian.Uint16(c)
}

func (u *U16) Encode(value interface{}) string {
	var u16 uint16
	switch v := value.(type) {
	case int:
		u16 = uint16(v)
	case decimal.Decimal:
		u16 = uint16(v.IntPart())
	case uint32:
		u16 = uint16(v)
	case int64:
		u16 = uint16(v)
	case uint16:
		u16 = v
	}
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, u16)
	return utiles.BytesToHex(bs)
}

type U32 struct {
	Reader io.Reader
	ScaleDecoder
}

func (u *U32) Process() {
	buf := &bytes.Buffer{}
	u.Reader = buf
	_, _ = buf.Write(u.NextBytes(4))
	c := make([]byte, 4)
	_, _ = u.Reader.Read(c)
	u.Value = binary.LittleEndian.Uint32(c)
}

func (u *U32) Encode(value interface{}) string {
	var u32 uint32
	switch v := value.(type) {
	case int:
		u32 = uint32(v)
	case decimal.Decimal:
		u32 = uint32(v.IntPart())
	case uint32:
		u32 = v
	}
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, u32)
	return utiles.BytesToHex(bs)
}

type U64 struct {
	ScaleDecoder
	Reader io.Reader
}

func (u *U64) Process() {
	buf := &bytes.Buffer{}
	u.Reader = buf
	_, _ = buf.Write(u.NextBytes(8))
	c := make([]byte, 8)
	_, _ = u.Reader.Read(c)
	u.Value = binary.LittleEndian.Uint64(c)
}

func (u *U64) Encode(value uint64) string {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, value)
	return utiles.BytesToHex(bs)
}

type U128 struct {
	ScaleDecoder
}

func (u *U128) Process() {
	elementBytes := u.NextBytes(16)
	if len(elementBytes) < 16 {
		elementBytes = utiles.HexToBytes(xstrings.LeftJustify(utiles.BytesToHex(elementBytes), 32, "0"))
	}
	u.Value = uint128.FromBytes(elementBytes).String()
}

func (u *U128) Encode(value decimal.Decimal) string {
	bs := make([]byte, 16)
	u128 := uint128.FromBig(value.BigInt())
	u128.PutBytes(bs)
	return utiles.BytesToHex(bs)
}
