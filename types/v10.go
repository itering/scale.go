package types

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/huandu/xstrings"
)

type MetadataV10Decoder struct {
	ScaleDecoder
	Version    string                 `json:"version"`
	Modules    []MetadataModules      `json:"modules"`
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

func (m *MetadataV10Decoder) Init(data ScaleBytes, subType string, arg ...interface{}) {
	m.ScaleDecoder.Init(data, subType, arg...)
}

func (m *MetadataV10Decoder) Process() {
	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	// result.CallIndex = make(map[string]interface{})
	// result.EventIndex = make(map[string]interface{})
	MetadataV10ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV8Module>").([]interface{})
	callModuleIndex := 0
	eventModuleIndex := 0
	bm, _ := json.Marshal(MetadataV10ModuleCall)
	var modulesType []MetadataModules
	_ = json.Unmarshal(bm, &modulesType)
	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(callModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
				// result.CallIndex[modulesType[k].Calls[callIndex].Lookup] = map[string]interface{}{
				// 	"module": module,
				// 	"call":   call,
				// }
			}
			callModuleIndex++
		}
		if module.Events != nil {
			for eventIndex := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(eventModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
				// result.EventIndex[modulesType[k].Events[eventIndex].Lookup] = map[string]interface{}{
				// 	"module": module,
				// 	"call":   event,
				// }
			}
			eventModuleIndex++
		}
	}

	result.Metadata.Modules = modulesType
	m.Value = result
}
