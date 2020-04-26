package types

import (
	"encoding/json"
	"github.com/freehere107/go-scale-codec/utiles"
	"github.com/huandu/xstrings"
)

type MetadataV11Decoder struct {
	ScaleDecoder
}

func (m *MetadataV11Decoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV11Decoder) Process() {
	var callModuleIndex, eventModuleIndex int

	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	MetadataV11ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV8Module>").([]interface{})
	bm, _ := json.Marshal(MetadataV11ModuleCall)

	var modulesType []MetadataModules
	_ = json.Unmarshal(bm, &modulesType)
	result.CallIndex = make(map[string]CallIndex)
	result.EventIndex = make(map[string]EventIndex)
	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex, call := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(callModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
				result.CallIndex[modulesType[k].Calls[callIndex].Lookup] = CallIndex{Module: module, Call: call}
			}
			callModuleIndex++
		}
		if module.Events != nil {
			for eventIndex, event := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(eventModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
				result.EventIndex[modulesType[k].Events[eventIndex].Lookup] = EventIndex{Module: module, Call: event}
			}
			eventModuleIndex++
		}
	}

	result.Metadata.Modules = modulesType
	m.Value = result
}
