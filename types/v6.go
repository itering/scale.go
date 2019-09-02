package types

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/huandu/xstrings"
	"reflect"
)

type MetadataV6Decoder struct {
	ScaleDecoder
	Version    string                 `json:"version"`
	Modules    []MetadataModules      `json:"modules"`
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

func (m *MetadataV6Decoder) Init(data ScaleBytes, valueList []string) {
	subType := ""
	if len(valueList) > 0 {
		subType = valueList[0]
	}
	m.ScaleDecoder.Init(data, subType)
}

func (m *MetadataV6Decoder) Process() string {
	result := MetadataStruct{
		MagicNumber: 1635018093,
		Metadata: MetadataTag{
			Modules: nil,
		},
	}
	result.CallIndex = make(map[string]interface{})
	result.EventIndex = make(map[string]interface{})
	metadataV6ModuleCall := m.ProcessAndUpdateData("Vec<MetadataV6Module>").Interface().([]interface{})

	callModuleIndex := 0
	eventModuleIndex := 0
	var modules []map[string]interface{}
	for _, v := range metadataV6ModuleCall {
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

type MetadataV6Module struct {
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

func (m *MetadataV6Module) GetIdentifier() string {
	return m.Name
}

func (m *MetadataV6Module) Process() map[string]interface{} {
	m.Name = m.ProcessAndUpdateData("Bytes").String()
	m.Prefix = m.ProcessAndUpdateData("Bytes").String()

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
		storageValue := m.ProcessAndUpdateData("Vec<MetadataV6ModuleStorage>").Interface().([]interface{})
		var storage []MetadataStorage
		for _, v := range storageValue {
			var sv MetadataStorage
			_ = json.Unmarshal([]byte(v.(reflect.Value).Interface().(string)), &sv)
			storage = append(storage, sv)
		}
		result["storage"] = storage
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
	constantValue := m.ProcessAndUpdateData("Vec<MetadataV6ModuleConstants>").Interface().([]interface{})
	var constants []interface{}
	for _, v := range constantValue {
		var sv MetadataV6ModuleConstants
		_ = json.Unmarshal([]byte(v.(reflect.Value).Interface().(string)), &sv)
		constants = append(constants, sv)
	}
	result["constants"] = constants
	return result
}

type MetadataV6ModuleConstants struct {
	ScaleType
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Value string   `json:"value"`
	Docs  []string `json:"docs"`
}

func (m *MetadataV6ModuleConstants) Process() string {
	name := m.ProcessAndUpdateData("Bytes").String()
	cType := ConvertType(m.ProcessAndUpdateData("Bytes").String())
	value := m.ProcessAndUpdateData("HexBytes").String()
	var docsArr []string
	docs := m.ProcessAndUpdateData("Vec<Bytes>").Interface().([]interface{})
	for _, v := range docs {
		docsArr = append(docsArr, v.(reflect.Value).String())
	}
	r := map[string]interface{}{
		"name":  name,
		"type":  cType,
		"value": value,
		"docs":  docsArr,
	}
	br, _ := json.Marshal(r)
	return string(br)
}

type MetadataV6ModuleStorage struct {
	MetadataV5ModuleStorage
}