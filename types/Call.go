package types

import "github.com/itering/scale.go/utiles"

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

type BoxProposal struct {
	Call
}
