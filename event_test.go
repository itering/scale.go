package scalecodec_test

import (
	"encoding/json"
	"fmt"
	"github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
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

func TestCrabEventsDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Crab2))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "source/crab"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}

	eventRaw := "0x180000000000000080e36a090000000002000000010000000000000000000000000002000000020000000e022ac9219ace40f5846ed675dded4e25a1997da7eabdea2f78597a71d6f38031487089d481874e06bd4026807fd0464f5c7a1c691c21237ed9b912ac0443a2bc2200ca9a3b000000000000000000000000000002000000150600b4c4040000000000000000000000000000020000000e04b4f7f03bebc56ebe96bc52ea5ed3159d45a0ce3a8d7f082983c33ef133274747002d31010000000000000000000000000000020000000000401b5f1300000000000000"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
	e.Process()
	b, _ := json.Marshal(e.Value)
	fmt.Println(string(b))

}

func TestPlasmEventsDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Plasm))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "source/plasm"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}

	eventRaw := "0x1000000000000000102700000201000001000000000010270000020100000200000002020ce5142b4246d526a9cddc3e20201450294c7815227e4aa356e15c687e1d325c16b9129a07f70d3fa66666d27dec4fd25d7261c7f08a23883b082e945ea504620000c52ebca2b1000000000000000000000002000000000040420f00000100"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
	e.Process()
}

func TestWestendEventsDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(westend))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "source/westend"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata, Spec: 1}

	eventRaw := "0x180000000000000080e36a09000000000200000001000000000000000000000000000200000002000000000000ca9a3b0000000002000000030000000300c601b3e5d664c8cfd77f6713be93b8b4364e6a1e93217d04888b3c7cc21ee235005039278c04000000000000000000000000030000001100c601b3e5d664c8cfd77f6713be93b8b4364e6a1e93217d04888b3c7cc21ee235d0ee80934b74c7f0f25c7a137a8a16e58e713283005039278c0400000000000000000000000003000000000000ab904100000000000100"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
	e.Process()
}
