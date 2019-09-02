package types

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/huandu/xstrings"
	"reflect"
)

type MetadataV7Decoder struct {
	ScaleDecoder
	Version    string                 `json:"version"`
	Modules    []MetadataModules      `json:"modules"`
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

func (m *MetadataV7Decoder) Init(data ScaleBytes, valueList []string) {
	subType := ""
	if len(valueList) > 0 {
		subType = valueList[0]
	}
	m.ScaleDecoder.Init(data, subType)
}

func (m *MetadataV7Decoder) Process() string {
	result := MetadataStruct{
		MagicNumber: 1635018093,
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	result.CallIndex = make(map[string]interface{})
	result.EventIndex = make(map[string]interface{})
	metadataV7ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV7Module>").Interface().([]interface{})
	callModuleIndex := 0
	eventModuleIndex := 0
	var modules []map[string]interface{}
	for _, v := range metadataV7ModuleCall {
		s := v.(reflect.Value).Interface().(map[string]interface{})
		modules = append(modules, s)
	}
	bm, _ := json.Marshal(modules)
	var modulesType []MetadataModules
	_ = json.Unmarshal(bm, &modulesType)
	for k, module := range modulesType {
		if module.Calls != nil {
			for callIndex, call := range module.Calls {
				modulesType[k].Calls[callIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(callModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(callIndex), 2, "0")
				result.CallIndex[modulesType[k].Calls[callIndex].Lookup] = map[string]interface{}{
					"module": module,
					"call":   call,
				}
			}
			callModuleIndex += 1
		}
		if module.Events != nil {
			for eventIndex, event := range module.Events {
				modulesType[k].Events[eventIndex].Lookup = xstrings.RightJustify(utiles.IntToHex(eventModuleIndex), 2, "0") + xstrings.RightJustify(utiles.IntToHex(eventIndex), 2, "0")
				result.EventIndex[modulesType[k].Events[eventIndex].Lookup] = map[string]interface{}{
					"module": module,
					"call":   event,
				}
			}
			eventModuleIndex += 1
		}
	}

	result.Metadata.Modules = modulesType
	bResult, _ := json.Marshal(result)
	return string(bResult)
}

type MetadataV7Module struct {
	ScaleType
	Name       string                 `json:"name"`
	Prefix     string                 `json:"prefix"`
	CallIndex  string                 `json:"call_index"`
	HasStorage bool                   `json:"has_storage"`
	Storage    map[string]interface{} `json:"storage"`
	HasCalls   bool                   `json:"has_calls"`
	Calls      map[string]interface{} `json:"calls"`
	HasEvents  bool                   `json:"has_events"`
	Events     map[string]interface{} `json:"events"`
	Constants  []interface{}          `json:"constants"`
}

func (m *MetadataV7Module) GetIdentifier() string {
	return m.Name
}

func (m *MetadataV7Module) Process() map[string]interface{} {
	m.Name = m.ProcessAndUpdateData("Bytes").String()

	result := map[string]interface{}{
		"name":      m.Name,
		"prefix":    m.Prefix,
		"storage":   m.Storage,
		"calls":     m.Calls,
		"events":    m.Events,
		"constants": m.Constants,
	}
	m.HasStorage = m.ProcessAndUpdateData("bool").Bool()
	if m.HasStorage {
		storageValue := m.ProcessAndUpdateData("MetadataV7ModuleStorage").Interface().(map[string]interface{})
		result["storage"] = storageValue["items"].([]MetadataStorage)
		result["prefix"] = storageValue["prefix"].(string)
	}
	m.HasCalls = m.ProcessAndUpdateData("bool").Bool()
	if m.HasCalls {
		callValue := m.ProcessAndUpdateData("Vec<MetadataModuleCall>").Interface().([]interface{})
		var calls []MetadataModuleCall
		for _, v := range callValue {
			var sv MetadataModuleCall
			_ = json.Unmarshal([]byte(v.(reflect.Value).Interface().(string)), &sv)
			calls = append(calls, sv)
		}
		result["calls"] = calls
	}
	m.HasEvents = m.ProcessAndUpdateData("bool").Bool()
	if m.HasEvents {
		eventValue := m.ProcessAndUpdateData("Vec<MetadataModuleEvent>").Interface().([]interface{})
		var events []MetadataEvents
		for _, v := range eventValue {
			var sv MetadataEvents
			_ = json.Unmarshal([]byte(v.(reflect.Value).Interface().(string)), &sv)
			events = append(events, sv)
		}
		result["events"] = events
	}
	constantValue := m.ProcessAndUpdateData("Vec<MetadataV7ModuleConstants>").Interface().([]interface{})
	var constants []interface{}
	for _, v := range constantValue {
		var sv MetadataV7ModuleConstants
		_ = json.Unmarshal([]byte(v.(reflect.Value).Interface().(string)), &sv)
		constants = append(constants, sv)
	}
	result["constants"] = constants
	return result
}

type MetadataV7ModuleConstants struct {
	MetadataV6ModuleConstants
}

type MetadataV7ModuleStorage struct {
	MetadataV6ModuleStorage
	Prefix string            `json:"prefix"`
	Items  []MetadataStorage `json:"items"`
}

func (m *MetadataV7ModuleStorage) Process() map[string]interface{} {
	m.Prefix = m.ProcessAndUpdateData("Bytes").String()
	itemsValue := m.ProcessAndUpdateData("Vec<MetadataV7ModuleStorageEntry>").Interface().([]interface{})
	var items []MetadataStorage
	for _, v := range itemsValue {
		var sv MetadataStorage
		_ = json.Unmarshal([]byte(v.(reflect.Value).Interface().(string)), &sv)
		items = append(items, sv)
	}
	m.Items = items
	r := map[string]interface{}{
		"prefix": m.Prefix,
		"items":  m.Items,
	}
	return r
}

type MetadataV7ModuleStorageEntry struct {
	ScaleDecoder
	Name     string                 `json:"name"`
	Modifier string                 `json:"modifier"`
	Type     map[string]interface{} `json:"type"`
	Fallback string                 `json:"fallback"`
	Docs     []string               `json:"docs"`
	Hasher   string                 `json:"hasher"`
}

func (m *MetadataV7ModuleStorageEntry) Init(data ScaleBytes, valueList []string) {
	subType := ""
	if len(valueList) > 0 {
		subType = valueList[0]
	}
	m.ScaleDecoder.Init(data, subType)
}

func (m *MetadataV7ModuleStorageEntry) Process() string {
	m.Name = m.ProcessAndUpdateData("Bytes").String()
	m.Modifier = m.ProcessAndUpdateData("Enum", "Optional", "Default").String()
	storageFunctionType := m.ProcessAndUpdateData("Enum", "PlainType", "MapType", "DoubleMapType").String()
	if storageFunctionType == "MapType" {
		m.Hasher = m.ProcessAndUpdateData("StorageHasher").String()
		m.Type = map[string]interface{}{
			"MapType": map[string]interface{}{
				"hasher":   m.Hasher,
				"key":      ConvertType(m.ProcessAndUpdateData("Bytes").String()),
				"value":    ConvertType(m.ProcessAndUpdateData("Bytes").String()),
				"isLinked": m.ProcessAndUpdateData("bool").Bool(),
			},
		}
	} else if storageFunctionType == "DoubleMapType" {
		m.Hasher = m.ProcessAndUpdateData("StorageHasher").String()
		key1 := ConvertType(m.ProcessAndUpdateData("Bytes").String())
		key2 := ConvertType(m.ProcessAndUpdateData("Bytes").String())
		value := ConvertType(m.ProcessAndUpdateData("Bytes").String())
		key2Hasher := m.ProcessAndUpdateData("StorageHasher").String()
		m.Type = map[string]interface{}{
			"DoubleMapType": map[string]interface{}{
				"hasher":     m.Hasher,
				"key1":       key1,
				"key2":       key2,
				"value":      value,
				"key2Hasher": key2Hasher,
			},
		}
	} else if storageFunctionType == "PlainType" {
		m.Type = map[string]interface{}{
			"PlainType": ConvertType(m.ProcessAndUpdateData("Bytes").String()),
		}
	}
	m.Fallback = m.ProcessAndUpdateData("HexBytes").String()
	docs := m.ProcessAndUpdateData("Vec<Bytes>").Interface().([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(reflect.Value).String())
	}
	r := map[string]interface{}{
		"name":     m.Name,
		"modifier": m.Modifier,
		"type":     m.Type,
		"fallback": m.Fallback,
		"docs":     m.Docs,
	}
	br, _ := json.Marshal(r)
	return string(br)
}
