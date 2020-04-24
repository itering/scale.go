package types

import (
	"encoding/json"
	"fmt"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/huandu/xstrings"
)

type MetadataV7Decoder struct {
	ScaleDecoder
}

func (m *MetadataV7Decoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV7Decoder) Process() {
	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	result.CallIndex = make(map[string]CallIndex)
	result.EventIndex = make(map[string]EventIndex)
	metadataV7ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV7Module>").([]interface{})
	callModuleIndex := 0
	eventModuleIndex := 0
	bm, _ := json.Marshal(metadataV7ModuleCall)
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

type MetadataV7Module struct {
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

func (m *MetadataV7Module) Process() {
	cm := MetadataV7Module{}
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
	constantValue := m.ProcessAndUpdateData("Vec<MetadataV7ModuleConstants>").([]interface{})
	var constants []map[string]interface{}
	for _, v := range constantValue {
		constants = append(constants, v.(map[string]interface{}))
	}
	cm.Constants = constants
	m.Value = cm
}

type MetadataV7ModuleConstants struct {
	MetadataV6ModuleConstants
}

type MetadataV7ModuleStorage struct {
	MetadataV6ModuleStorage
	Prefix string            `json:"prefix"`
	Items  []MetadataStorage `json:"items"`
}

func (m *MetadataV7ModuleStorage) Process() {
	cm := MetadataV7ModuleStorage{}
	cm.Prefix = m.ProcessAndUpdateData("Bytes").(string)
	itemsValue := m.ProcessAndUpdateData("Vec<MetadataV7ModuleStorageEntry>").([]interface{})
	var items []MetadataStorage
	for _, v := range itemsValue {
		items = append(items, v.(MetadataStorage))
	}
	cm.Items = items
	m.Value = cm
}

type MetadataV7ModuleStorageEntry struct {
	ScaleDecoder
	Name     string      `json:"name"`
	Modifier string      `json:"modifier"`
	Type     StorageType `json:"type"`
	Fallback string      `json:"fallback"`
	Docs     []string    `json:"docs"`
	Hasher   string      `json:"hasher"`
}

func (m *MetadataV7ModuleStorageEntry) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV7ModuleStorageEntry) Process() {
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	m.Modifier = m.ProcessAndUpdateData("Enum", "Optional", "Default").(string)
	storageFunctionType := m.ProcessAndUpdateData("Enum", "PlainType", "MapType", "DoubleMapType").(string)
	if storageFunctionType == "MapType" {
		m.Hasher = m.ProcessAndUpdateData("StorageHasher").(string)
		m.Type = StorageType{
			Origin: "MapType",
			MapType: &MapType{
				Hasher:   m.Hasher,
				Key:      ConvertType(m.ProcessAndUpdateData("Bytes").(string)),
				Value:    ConvertType(m.ProcessAndUpdateData("Bytes").(string)),
				IsLinked: m.ProcessAndUpdateData("bool").(bool)},
		}
	} else if storageFunctionType == "DoubleMapType" {
		m.Hasher = m.ProcessAndUpdateData("StorageHasher").(string)
		key1 := ConvertType(m.ProcessAndUpdateData("Bytes").(string))
		key2 := ConvertType(m.ProcessAndUpdateData("Bytes").(string))
		value := ConvertType(m.ProcessAndUpdateData("Bytes").(string))
		key2Hasher := m.ProcessAndUpdateData("StorageHasher").(string)
		m.Type = StorageType{
			Origin: "DoubleMapType",
			DoubleMapType: &MapType{
				Hasher:     m.Hasher,
				Key:        key1,
				Key2:       key2,
				Value:      value,
				Key2Hasher: key2Hasher,
			},
		}
	} else if storageFunctionType == "PlainType" {
		plainType := ConvertType(m.ProcessAndUpdateData("Bytes").(string))
		m.Type = StorageType{
			Origin:    "PlainType",
			PlainType: &plainType}
	}
	m.Fallback = m.ProcessAndUpdateData("HexBytes").(string)
	docs := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(string))
	}
	m.Value = MetadataStorage{
		Name:     m.Name,
		Modifier: m.Modifier,
		Type:     m.Type,
		Fallback: m.Fallback,
		Docs:     m.Docs,
	}
	fmt.Println(m.Value)

}
