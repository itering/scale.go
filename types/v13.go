package types

import (
	"encoding/json"
	"github.com/huandu/xstrings"
	"github.com/itering/scale.go/utiles"
)

type MetadataV13Decoder struct {
	ScaleDecoder
}

func (m *MetadataV13Decoder) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV13Decoder) Process() {
	result := MetadataStruct{
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	MetadataV11ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV13Module>").([]interface{})
	bm, _ := json.Marshal(MetadataV11ModuleCall)

	var modulesType []MetadataModules
	_ = json.Unmarshal(bm, &modulesType)
	result.CallIndex = make(map[string]CallIndex)
	result.EventIndex = make(map[string]EventIndex)
	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex, call := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(module.Index), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
				result.CallIndex[modulesType[k].Calls[callIndex].Lookup] = CallIndex{Module: module, Call: call}
			}
		}
		if module.Events != nil {
			for eventIndex, event := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(module.Index), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
				result.EventIndex[modulesType[k].Events[eventIndex].Lookup] = EventIndex{Module: module, Call: event}
			}
		}
	}

	result.Metadata.Modules = modulesType
	extrinsicMetadata := m.ProcessAndUpdateData("ExtrinsicMetadata").(map[string]interface{})
	bm, _ = json.Marshal(extrinsicMetadata)
	_ = json.Unmarshal(bm, &result.Extrinsic)

	m.Value = result
}

type MetadataV13Module struct {
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
	Index      int                      `json:"index"`
}

func (m *MetadataV13Module) Process() {
	cm := MetadataV13Module{}
	cm.Name = m.ProcessAndUpdateData("String").(string)

	// storage
	cm.HasStorage = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasStorage {
		storageValue := m.ProcessAndUpdateData("MetadataV13ModuleStorage").(MetadataV7ModuleStorage)
		cm.Storage = storageValue.Items
		cm.Prefix = storageValue.Prefix
	}

	// call
	cm.HasCalls = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasCalls {
		callValue := m.ProcessAndUpdateData("Vec<MetadataModuleCall>").([]interface{})
		var calls []MetadataModuleCall
		for _, v := range callValue {
			calls = append(calls, v.(MetadataModuleCall))
		}
		cm.Calls = calls
	}

	// event
	cm.HasEvents = m.ProcessAndUpdateData("bool").(bool)
	if cm.HasEvents {
		eventValue := m.ProcessAndUpdateData("Vec<MetadataModuleEvent>").([]interface{})
		var events []MetadataEvents
		for _, v := range eventValue {
			events = append(events, v.(MetadataEvents))
		}
		cm.Events = events
	}

	// constant
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
	cm.Index = m.ProcessAndUpdateData("U8").(int)
	m.Value = cm
}

type MetadataV13ModuleStorage struct {
	MetadataV6ModuleStorage
	Prefix string            `json:"prefix"`
	Items  []MetadataStorage `json:"items"`
}

func (m *MetadataV13ModuleStorage) Process() {
	cm := MetadataV7ModuleStorage{}
	cm.Prefix = m.ProcessAndUpdateData("String").(string)
	itemsValue := m.ProcessAndUpdateData("Vec<MetadataV13ModuleStorageEntry>").([]interface{})
	var items []MetadataStorage
	for _, v := range itemsValue {
		items = append(items, v.(MetadataStorage))
	}
	cm.Items = items
	m.Value = cm
}

type MetadataV13ModuleStorageEntry struct {
	ScaleDecoder
	Name     string      `json:"name"`
	Modifier string      `json:"modifier"`
	Type     StorageType `json:"type"`
	Fallback string      `json:"fallback"`
	Docs     []string    `json:"docs"`
	Hasher   string      `json:"hasher"`
}

func (m *MetadataV13ModuleStorageEntry) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.ScaleDecoder.Init(data, option)
}

func (m *MetadataV13ModuleStorageEntry) Process() {
	m.Name = m.ProcessAndUpdateData("String").(string)
	m.Modifier = m.ProcessAndUpdateData("StorageModify").(string)
	storageFunctionType := m.ProcessAndUpdateData("StorageFunctionType").(string)

	switch storageFunctionType {
	case "MapType":
		m.Hasher = m.ProcessAndUpdateData("StorageHasher").(string)
		m.Type = StorageType{
			Origin: "MapType",
			MapType: &MapType{
				Hasher:   m.Hasher,
				Key:      ConvertType(m.ProcessAndUpdateData("String").(string)),
				Value:    ConvertType(m.ProcessAndUpdateData("String").(string)),
				IsLinked: m.ProcessAndUpdateData("bool").(bool)},
		}
	case "DoubleMapType":
		m.Hasher = m.ProcessAndUpdateData("StorageHasher").(string)
		key1 := ConvertType(m.ProcessAndUpdateData("String").(string))
		key2 := ConvertType(m.ProcessAndUpdateData("String").(string))
		value := ConvertType(m.ProcessAndUpdateData("String").(string))
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
	case "PlainType":
		plainType := ConvertType(m.ProcessAndUpdateData("String").(string))
		m.Type = StorageType{
			Origin:    "PlainType",
			PlainType: &plainType}
	case "NMap":
		hashers := m.ProcessAndUpdateData("Vec<StorageHasher>").([]string)
		var KeyVec []string
		for _, v := range m.ProcessAndUpdateData("Vec<String>").([]string) {
			KeyVec = append(KeyVec, ConvertType(v))
		}
		m.Type = StorageType{
			Origin: "NMapType",
			NMapType: &NMapType{
				Hashers: hashers,
				KeyVec:  KeyVec,
				Value:   ConvertType(m.ProcessAndUpdateData("String").(string)),
			}}
	}

	m.Fallback = m.ProcessAndUpdateData("HexBytes").(string)
	docs := m.ProcessAndUpdateData("Vec<String>").([]interface{})
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
}
