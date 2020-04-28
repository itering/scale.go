package scalecodec_test

import (
	"fmt"
	"github.com/freehere107/go-scale-codec"
	"github.com/freehere107/go-scale-codec/source"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/freehere107/go-scale-codec/utiles"
	"io/ioutil"
	"testing"
)

func TestEventsDecoder_Init(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Kusama1055))
	_ = m.Process()
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "source/kusama"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	eventRaw := "0x10020f001400000000000000000000001027000001010000010000000000102700000001000002000000000040420f00000100"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
}

func TestKusama1058EventsDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Kusama1055))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "source/kusama"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata, Spec: 1058}

	eventRaw := "0x0c0000000000000080969800000000000201000001000000000080969800000000000201000002000000000000ca9a3b00000000020100"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
	e.Process()
}
