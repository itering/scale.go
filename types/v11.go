package types

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
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

	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(callModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
			}
			callModuleIndex++
		}
		if module.Events != nil {
			for eventIndex := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(eventModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
			}
			eventModuleIndex++
		}
	}

	result.Metadata.Modules = modulesType
	m.Value = result
}
