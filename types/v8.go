package types

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/huandu/xstrings"
)

type MetadataV8Decoder struct {
	ScaleDecoder
}

func (m *MetadataV8Decoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV8Decoder) Process() {
	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	metadataV8ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV8Module>").([]interface{})
	callModuleIndex := 0
	eventModuleIndex := 0
	result.CallIndex = make(map[string]CallIndex)
	result.EventIndex = make(map[string]EventIndex)
	bm, _ := json.Marshal(metadataV8ModuleCall)
	var modulesType []MetadataModules
	_ = json.Unmarshal(bm, &modulesType)
	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex, call := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(callModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
				result.CallIndex[modulesType[k].Calls[callIndex].Lookup] = CallIndex{
					Module: module,
					Call:   call,
				}
			}
			callModuleIndex++
		}
		if module.Events != nil {
			for eventIndex, event := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(eventModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
				result.EventIndex[modulesType[k].Events[eventIndex].Lookup] = EventIndex{
					Module: module,
					Call:   event,
				}
			}
			eventModuleIndex++
		}
	}

	result.Metadata.Modules = modulesType
	m.Value = result
}

type MetadataV8Module struct {
	ScaleDecoder
	Name       string                   `json:"name"`
	Prefix     string                   `json:"prefix"`
	CallIndex  string                   `json:"call_index"`
	HasStorage bool                     `json:"has_storage"`
	Storage    []MetadataStorage        `json:"storage"`
	HasCalls   bool                     `json:"has_calls"`
	Calls      []MetadataModuleCall     `json:"calls"`
	HasEvents  bool                     `json:"has_events"`
	Events     []MetadataEvents         `json:"events"`
	Constants  []map[string]interface{} `json:"constants"`
	Errors     []MetadataModuleError    `json:"errors"`
}

func (m *MetadataV8Module) GetIdentifier() string {
	return m.Name
}

func (m *MetadataV8Module) Process() {
	cm := MetadataV8Module{}
	cm.Name = m.ProcessAndUpdateData("Bytes").(string)

	cm.HasStorage = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasStorage {
		storageValue := m.ProcessAndUpdateData("MetadataV7ModuleStorage").(MetadataV7ModuleStorage)
		cm.Storage = storageValue.Items
		cm.Prefix = storageValue.Prefix
	}
	cm.HasCalls = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasCalls {
		callValue := m.ProcessAndUpdateData("Vec<MetadataModuleCall>").([]interface{})
		calls := []MetadataModuleCall{}
		for _, v := range callValue {
			calls = append(calls, v.(MetadataModuleCall))
		}
		cm.Calls = calls
	}
	cm.HasEvents = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasEvents {
		eventValue := m.ProcessAndUpdateData("Vec<MetadataModuleEvent>").([]interface{})
		events := []MetadataEvents{}
		for _, v := range eventValue {
			events = append(events, v.(MetadataEvents))
		}
		cm.Events = events
	}
	constantValue := m.ProcessAndUpdateData("Vec<MetadataV7ModuleConstants>").([]interface{})
	var constants []map[string]interface{}
	for _, v := range constantValue {
		constants = append(constants, v.(map[string]interface{}))
	}
	cm.Constants = constants

	errorValue := m.ProcessAndUpdateData("Vec<MetadataModuleError>").([]interface{})
	var errors []MetadataModuleError
	for _, v := range errorValue {
		errors = append(errors, v.(MetadataModuleError))
	}
	cm.Errors = errors
	m.Value = cm
}

type MetadataModuleError struct {
	ScaleDecoder `json:"-"`
	Name         string   `json:"name"`
	Doc          []string `json:"doc"`
}

func (m *MetadataModuleError) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.Name = ""
	m.Doc = []string{}
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataModuleError) Process() {
	cm := MetadataModuleError{}
	cm.Name = m.ProcessAndUpdateData("Bytes").(string)
	var docsArr []string
	docs := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range docs {
		docsArr = append(docsArr, v.(string))
	}
	cm.Doc = docsArr
	m.Value = cm
}
