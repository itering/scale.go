package scalecodec_test

import (
	"encoding/json"
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

func TestCrabEventsDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Crab2))
	_ = m.Process()
	c, err := ioutil.ReadFile(fmt.Sprintf("%s.json", "source/crab"))
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(c))
	e := types.ScaleDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}

	eventRaw := "0x02286bee00b3e60100b3e6010000"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
	v := e.ProcessAndUpdateData("ExposureT")
	cb, _ := json.Marshal(v)
	fmt.Println(string(cb))
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
