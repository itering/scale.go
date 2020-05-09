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

	eventRaw := "0x1800000000000000000000000000000002000000010000000000000000000000000002000000020000001401b62d88e3f439fe9b5ea799b27bf7c6db5e795de1784f27b1bc051553499e420f912e313b310c2bfe1bf0889aaea95f4d70a48e31ce617ae35b5ab189affd2ce44ec8a85e000000007b0a980000000000829bd824b016326a401d083b33d092293333a830fb5a9e8944ed036f6f0c169472c5b88f0dec77cfc3ae12c07e1b0c6ca4fc3d491dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347607070796520e4b883e5bda9e7a59ee4bb99e9b1bcbb249e98b2230bf5c223669791eddb0348a070e3fc088acb3ddedd9001dcbacf06d1717711dd1b514308ba9d12a75c6a0a785e25be6b65fc21332682c3f3d7371e2bab04084a11200130200824848800600042000282080014900010200723101011000820104c102008000810600000020041011010945c42000808004012801030040800012000900044219100d10e2400a4840900000200200402800000884000900074000001020010010001c100000498400c4000a00a42200108b02214b03070000020800700401010640002906200401000010d080182000040505002083000130250201040244000020420920083800100080005a0422a2200000088004003000c491842000000004418002400000820000442420208800840900000105020080014241100109018800021104018e006000900b0000005a4c00032200440000c05253700000000000000000000000000000000000000000000000000000000001e53980000000000000000000000000000000000000000000000000000000000364ef678adc307000000000000000000000000000000000000000000000000000884a0eb20c32d420fbfa5216c2a1966cc8b4ff6a069cdc63cc4eee9bd72eda335e5bf248816be29e804403ffd01a9f234777acf9a7fcc3975d8e096c697bc5b633f49820a6aa68afeecdfd3611e0000020000001506f6dea00e0000000000000000000000000000020000000e04b4f7f03bebc56ebe96bc52ea5ed3159d45a0ce3a8d7f082983c33ef133274747be37a803000000000000000000000000000002000000000000c2eb0b00000000000000"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(eventRaw)}, &option)
	e.Process()
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
