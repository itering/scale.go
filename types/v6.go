package types

import (
	"encoding/json"
	"github.com/huandu/xstrings"
	"github.com/itering/scale.go/utiles"
)

type MetadataV6Decoder struct {
	ScaleDecoder
}

func (m *MetadataV6Decoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV6Decoder) Process() {
	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	result.CallIndex = make(map[string]CallIndex)
	result.EventIndex = make(map[string]EventIndex)
	metadataV6ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV6Module>").([]interface{})

	callModuleIndex := 0
	eventModuleIndex := 0
	bm, _ := json.Marshal(metadataV6ModuleCall)
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

type MetadataV6Module struct {
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
}

func (m *MetadataV6Module) GetIdentifier() string {
	return m.Name
}

func (m *MetadataV6Module) Process() {
	cm := MetadataV6Module{}
	cm.Name = m.ProcessAndUpdateData("String").(string)
	cm.Prefix = m.ProcessAndUpdateData("String").(string)
	cm.HasStorage = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasStorage {
		storageValue := m.ProcessAndUpdateData("Vec<MetadataV6ModuleStorage>").([]interface{})
		var storage []MetadataStorage
		for _, v := range storageValue {
			storage = append(storage, v.(MetadataStorage))
		}
		cm.Storage = storage
	}

	cm.HasCalls = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasCalls {
		callValue := m.ProcessAndUpdateData("Vec<MetadataModuleCall>").([]interface{})
		var calls []MetadataModuleCall
		for _, v := range callValue {
			calls = append(calls, v.(MetadataModuleCall))
		}
		cm.Calls = calls
	}
	cm.HasEvents = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasEvents {
		eventValue := m.ProcessAndUpdateData("Vec<MetadataModuleEvent>").([]interface{})
		var events []MetadataEvents
		for _, v := range eventValue {
			events = append(events, v.(MetadataEvents))
		}
		cm.Events = events
	}
	constantValue := m.ProcessAndUpdateData("Vec<MetadataV6ModuleConstants>").([]interface{})
	var constants []map[string]interface{}
	for _, v := range constantValue {
		constants = append(constants, v.(map[string]interface{}))
	}
	cm.Constants = constants
	m.Value = cm
}

type MetadataV6ModuleConstants struct {
	ScaleDecoder
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	ConstantsValue string   `json:"constants_value"`
	Docs           []string `json:"docs"`
}

func (m *MetadataV6ModuleConstants) Process() {
	name := m.ProcessAndUpdateData("String").(string)
	cType := ConvertType(m.ProcessAndUpdateData("String").(string))
	ConstantsValue := m.ProcessAndUpdateData("HexBytes").(string)
	var docsArr []string
	docs := m.ProcessAndUpdateData("Vec<String>").([]interface{})
	for _, v := range docs {
		docsArr = append(docsArr, v.(string))
	}
	r := map[string]interface{}{
		"name":            name,
		"type":            cType,
		"constants_value": ConstantsValue,
		"docs":            docsArr,
	}
	m.Value = r
}

type MetadataV6ModuleStorage struct {
	ScaleDecoder
	Name     string                 `json:"name"`
	Modifier string                 `json:"modifier"`
	Type     map[string]interface{} `json:"type"`
	Fallback string                 `json:"fallback"`
	Docs     []string               `json:"docs"`
	Hasher   string                 `json:"hasher"`
}

func (m *MetadataV6ModuleStorage) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV6ModuleStorage) Process() {
	cm := MetadataStorage{}
	cm.Name = m.ProcessAndUpdateData("String").(string)
	cm.Modifier = m.ProcessAndUpdateData("StorageModify").(string)
	storageFunctionType := m.ProcessAndUpdateData("StorageFunctionType").(string)
	if storageFunctionType == "MapType" {
		cm.Hasher = m.ProcessAndUpdateData("StorageHasher").(string)
		cm.Type = StorageType{
			Origin: "MapType",
			MapType: &MapType{
				Hasher:   cm.Hasher,
				Key:      ConvertType(m.ProcessAndUpdateData("String").(string)),
				Value:    ConvertType(m.ProcessAndUpdateData("String").(string)),
				IsLinked: m.ProcessAndUpdateData("bool").(bool),
			},
		}
	} else if storageFunctionType == "DoubleMapType" {
		cm.Hasher = m.ProcessAndUpdateData("StorageHasher").(string)
		key1 := ConvertType(m.ProcessAndUpdateData("String").(string))
		key2 := ConvertType(m.ProcessAndUpdateData("String").(string))
		value := ConvertType(m.ProcessAndUpdateData("String").(string))
		key2Hasher := m.ProcessAndUpdateData("StorageHasher").(string)
		cm.Type = StorageType{
			Origin: "DoubleMapType",
			DoubleMapType: &MapType{
				Hasher:     cm.Hasher,
				Key:        key1,
				Key2:       key2,
				Value:      value,
				Key2Hasher: key2Hasher,
			},
		}
	} else if storageFunctionType == "PlainType" {
		plainType := ConvertType(m.ProcessAndUpdateData("String").(string))
		cm.Type = StorageType{Origin: "PlainType", PlainType: &plainType}
	}
	cm.Fallback = m.ProcessAndUpdateData("HexBytes").(string)
	docs := m.ProcessAndUpdateData("Vec<String>").([]interface{})
	for _, v := range docs {
		cm.Docs = append(m.Docs, v.(string))
	}
	m.Value = cm
}
