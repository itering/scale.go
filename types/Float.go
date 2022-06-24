package types

import (
	"bytes"
	"encoding/binary"

	"github.com/itering/scale.go/types/scaleBytes"
)

type Float struct {
	ScaleDecoder
	bitLength     int
	encodedLength int
}

type Float64 struct {
	Float
}

func (f *Float64) Init(data scaleBytes.ScaleBytes, option *ScaleDecoderOption) {
	f.bitLength = 64
	f.encodedLength = f.bitLength / 8
	f.ScaleDecoder.Init(data, option)
}

func (f *Float64) Process() {
	var f64 float64
	buf := bytes.NewReader(f.Data.Data)
	err := binary.Read(buf, binary.LittleEndian, &f64)
	if err != nil {
		panic(err)
	}
	f.Value = f64
}

type Float32 struct {
	Float
}

func (f *Float32) Init(data scaleBytes.ScaleBytes, option *ScaleDecoderOption) {
	f.bitLength = 32
	f.encodedLength = f.bitLength / 8
	f.ScaleDecoder.Init(data, option)
}

func (f *Float32) Process() {
	var f32 float32
	buf := bytes.NewReader(f.Data.Data)
	err := binary.Read(buf, binary.LittleEndian, &f32)
	if err != nil {
		panic(err)
	}
	f.Value = f32
}
