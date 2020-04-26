package scalecodec_test

import (
	"encoding/json"
	"github.com/freehere107/go-scale-codec"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/freehere107/go-scale-codec/utiles"
	"testing"
)

func TestExtrinsicDecoder_Init(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Kusama1055))
	_ = m.Process()

	e := scalecodec.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes("0x280402000b900b7aa47101")}, &option)
	if utiles.BytesToHex(e.Data.Data) != "280402000b900b7aa47101" {
		t.Errorf("Test TestExtrinsicDecoder_Init fail")
	}
}

func TestExtrinsicDecoder(t *testing.T) {
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(Kusama1055))
	_ = m.Process()

	e := scalecodec.ExtrinsicDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}

	extrinsicRaw := "0x280402000b900b7aa47101"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsicRaw)}, &option)
	e.Process()
	b, _ := json.Marshal(e.Value)
	r := `{"call_code":"0200","call_module":"Timestamp","call_module_function":"set","era":"","extrinsic_length":10,"nonce":0,"params":[{"name":"now","type":"Compact\u003cMoment\u003e","value":1587602394,"value_raw":"0b900b7aa47101"}],"tip":"0","version_info":"04"}`
	if string(b) != r {
		t.Errorf("Test TestExtrinsicDecoder Process fail, decode return %s", string(b))
	}

	// sign
	extrinsicRaw = "0x39028430e078868ac70f8959c09deaa9a58880a0eaabb3c6145d042a54e13e70a7701801b0d317678364078aa289f8ec04851a7ae51d8c3a1dba7b121f62c56ae505b37208a1428fce37db2e9cf04f08a97569b0aada4185cb4ce92b5bc25e500dc49d83650010000400ee0b2c10108fe3d074ecd0f0bbef1aad46eaf865c0e00bd533878e1c14b5c0350700d0ed902e"
	e.Init(types.ScaleBytes{Data: utiles.HexToBytes(extrinsicRaw)}, &option)
	e.Process()
}
