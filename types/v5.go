package types

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/huandu/xstrings"
	"reflect"
)

type MetadataV5Decoder struct {
	ScaleDecoder
	Version    string                 `json:"version"`
	Modules    []MetadataModules      `json:"modules"`
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

func (m *MetadataV5Decoder) Init(data ScaleBytes, valueList []string) {
	subType := ""
	if len(valueList) > 0 {
		subType = valueList[0]
	}
	m.ScaleDecoder.Init(data, subType)
}

func (m *MetadataV5Decoder) Process() string {
	result := MetadataStruct{
		MagicNumber: 1635018093,
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	result.CallIndex = make(map[string]interface{})
	result.EventIndex = make(map[string]interface{})
	metadataV5ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV5Module>").Interface().([]interface{})

	callModuleIndex := 0
	eventModuleIndex := 0
	var modules []map[string]interface{}
	for _, v := range metadataV5ModuleCall {
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

type MetadataV5Module struct {
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
}

func (m *MetadataV5Module) GetIdentifier() string {
	return m.Name
}

func (m *MetadataV5Module) Process() map[string]interface{} {
	m.Name = m.ProcessAndUpdateData("Bytes").String()
	m.Prefix = m.ProcessAndUpdateData("Bytes").String()

	result := map[string]interface{}{
		"name":    m.Name,
		"prefix":  m.Prefix,
		"storage": m.Storage,
		"calls":   m.Calls,
		"events":  m.Events,
	}
	m.HasStorage = m.ProcessAndUpdateData("bool").Bool()
	if m.HasStorage {
		storageValue := m.ProcessAndUpdateData("Vec<MetadataV5ModuleStorage>").Interface().([]interface{})
		var storage []MetadataStorage
		for _, v := range storageValue {
			s := v.(reflect.Value).Interface().(string)
			var sv MetadataStorage
			_ = json.Unmarshal([]byte(s), &sv)
			storage = append(storage, sv)
		}
		result["storage"] = storage
	}
	m.HasCalls = m.ProcessAndUpdateData("bool").Bool()
	if m.HasCalls {
		callValue := m.ProcessAndUpdateData("Vec<MetadataModuleCall>").Interface().([]interface{})
		var calls []MetadataModuleCall
		for _, v := range callValue {
			s := v.(reflect.Value).Interface().(string)
			var sv MetadataModuleCall
			_ = json.Unmarshal([]byte(s), &sv)
			calls = append(calls, sv)
		}
		result["calls"] = calls
	}
	m.HasEvents = m.ProcessAndUpdateData("bool").Bool()
	if m.HasEvents {
		eventValue := m.ProcessAndUpdateData("Vec<MetadataModuleEvent>").Interface().([]interface{})
		var events []MetadataEvents
		for _, v := range eventValue {
			s := v.(reflect.Value).Interface().(string)
			var sv MetadataEvents
			_ = json.Unmarshal([]byte(s), &sv)
			events = append(events, sv)
		}
		result["events"] = events
	}
	return result

}

type MetadataV5ModuleStorage struct {
	ScaleType
	Name     string                 `json:"name"`
	Modifier string                 `json:"modifier"`
	Type     map[string]interface{} `json:"type"`
	Fallback string                 `json:"fallback"`
	Docs     []string               `json:"docs"`
	Hasher   string                 `json:"hasher"`
}

func (m *MetadataV5ModuleStorage) Init(data ScaleBytes, valueList []string) {
	subType := ""
	if len(valueList) > 0 {
		subType = valueList[0]
	}
	m.ScaleDecoder.Init(data, subType)
}

func (m *MetadataV5ModuleStorage) Process() string {
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
