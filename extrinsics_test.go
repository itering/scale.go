package scalecodec_test

import (
	"encoding/json"
	"github.com/freehere107/scalecodec"
	"github.com/freehere107/scalecodec/types"
	"github.com/freehere107/scalecodec/utiles"
	"testing"
)

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
	r := `{"call_code":"0200","call_module":"Timestamp","call_module_function":"set","era":"","extrinsic_length":10,"nonce":0,"params":[{"name":"now","type":"Compact\u003cMoment\u003e","value":1587602394,"value_raw":"040200"}],"tip":"","valueRaw":"040200","version_info":"04"}`
	if string(b) != r {
		t.Errorf("Test TestExtrinsicDecoder Process fail, decode return %s", string(b))
	}
}
