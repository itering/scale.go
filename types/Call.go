package types

import (
	"encoding/json"
	"fmt"

	"github.com/itering/scale.go/utiles"
)

type Call struct {
	ScaleDecoder
}

func (s *Call) Process() {
	callIndex := utiles.BytesToHex(s.NextBytes(2))
	callModule := s.Metadata.CallIndex[callIndex]
	result := map[string]interface{}{
		"call_index":  callIndex,
		"call_name":   callModule.Call.Name,
		"call_module": callModule.Module.Name,
	}
	var param []ExtrinsicParam
	s.Module = callModule.Module.Name
	for _, arg := range callModule.Call.Args {
		param = append(param, ExtrinsicParam{Name: arg.Name, Type: arg.Type, Value: s.ProcessAndUpdateData(arg.Type)})
	}
	result["params"] = param
	s.Value = result
}

type BoxCall struct {
	CallIndex string           `json:"call_Index"`
	Params    []ExtrinsicParam `json:"params"`
}

func (s *Call) Encode(value interface{}) string {
	if s.Metadata == nil {
		panic("nil metadata")
	}
	var boxCall BoxCall
	switch v := value.(type) {
	case BoxCall:
		boxCall = v
	case map[string]interface{}:
		b, _ := json.Marshal(v)
		if err := json.Unmarshal(b, &boxCall); err != nil {
			panic("input value is not valid boxCall")
		}
	default:
		panic("input value is not valid boxCall")
	}
	callIndex := utiles.TrimHex(boxCall.CallIndex)
	raw := callIndex
	callModule, ok := s.Metadata.CallIndex[callIndex]
	if !ok {
		panic(fmt.Sprintf("Not find Extrinsic Lookup %s, please check metadata info", callIndex))
	}
	for index, arg := range callModule.Call.Args {
		raw += EncodeWithOpt(arg.Type, boxCall.Params[index].Value, &ScaleDecoderOption{Spec: s.Spec})
	}
	return raw
}

type BoxProposal struct {
	Call
}
