package scalecodec

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/utiles"
	"golang.org/x/crypto/blake2b"
)

type ExtrinsicsDecoder struct {
	ScaleDecoder
	Metadata            MetadataDecoder          `json:"metadata"`
	ExtrinsicLength     int                      `json:"extrinsic_length"`
	ExtrinsicHash       string                   `json:"extrinsic_hash"`
	VersionInfo         string                   `json:"version_info"`
	ContainsTransaction bool                     `json:"contains_transaction"`
	Address             map[string]string        `json:"address"`
	Signature           string                   `json:"signature"`
	Nonce               int                      `json:"nonce"`
	Era                 string                   `json:"era"`
	CallIndex           string                   `json:"call_index"`
	CallModule          MetadataModules          `json:"call_module"`
	Call                MetadataCalls            `json:"call"`
	Params              []map[string]interface{} `json:"params"`
}

func (e *ExtrinsicsDecoder) Init(data ScaleBytes, args []string) {
	e.TypeMapping = map[string]interface{}{
		"extrinsic_length": "Compact<u32>",
		"version_info":     "u8",
		"address":          "Address",
		"signature":        "Signature",
		"nonce":            "Compact<u32>",
		"era":              "Era",
		"call_index":       "(u8,u8)",
	}
	var metadata MetadataDecoder
	var subType string
	if len(args) > 0 {
		subType = args[0]
	}
	if len(args) > 1 {
		_ = json.Unmarshal([]byte(args[1]), &metadata)
	}
	e.Metadata = metadata
	e.ScaleDecoder.Init(data, subType)
}

func (e *ExtrinsicsDecoder) generateHash() string {
	if e.ContainsTransaction {
		var extrinsicData []byte
		if e.ExtrinsicLength > 0 {
			extrinsicData = e.Data.Data
		} else {
			extrinsicLengthType := CompactU32{}
			extrinsicLengthType.encode(e.Data.Length)
			extrinsicData = append(extrinsicLengthType.Data.Data[:], e.Data.Data[:]...)
		}
		checksum, _ := blake2b.New(32, []byte{})
		checksum.Write(extrinsicData)
		h := checksum.Sum(nil)
		return utiles.BytesToHex(h)
	}
	return ""
}

func (e *ExtrinsicsDecoder) Process() map[string]interface{} {
	e.ExtrinsicLength = int(e.ProcessAndUpdateData("Compact<u32>").Int())
	if e.ExtrinsicLength != e.Data.getRemainingLength() {
		e.ExtrinsicLength = 0
		e.Data.reset()
	}
	e.VersionInfo = utiles.BytesToHex(e.getNextBytes(1))
	e.ContainsTransaction = utiles.U256(e.VersionInfo).Int64() >= 80
	if e.ContainsTransaction {
		e.Address = e.ProcessAndUpdateData("Address").Interface().(map[string]string)
		e.Signature = e.ProcessAndUpdateData("Signature").String()
		e.Nonce = int(e.ProcessAndUpdateData(e.TypeMapping["nonce"].(string)).Int())
		e.Era = e.ProcessAndUpdateData("Era").String()
		e.ExtrinsicHash = e.generateHash()
	}
	e.CallIndex = utiles.BytesToHex(e.getNextBytes(2))
	if e.CallIndex != "" {
		if e.Metadata.CallIndex[e.CallIndex] != nil {
			callIndex := e.Metadata.CallIndex[e.CallIndex].(map[string]interface{})
			bc, _ := json.Marshal(callIndex["call"])
			var call MetadataCalls
			_ = json.Unmarshal(bc, &call)
			e.Call = call
			var CallModule MetadataModules
			bc, _ = json.Marshal(callIndex["module"])
			_ = json.Unmarshal(bc, &CallModule)
			e.CallModule = CallModule
		}
	}
	for _, arg := range e.Call.Args {
		argTypeObj := e.ProcessAndUpdateData(arg["type"].(string)).Interface()
		e.Params = append(e.Params, map[string]interface{}{
			"name":     arg["name"].(string),
			"type":     arg["type"].(string),
			"value":    argTypeObj,
			"valueRaw": "",
		})
	}

	result := map[string]interface{}{
		"valueRaw":         e.RawValue,
		"extrinsic_length": e.ExtrinsicLength,
		"version_info":     e.VersionInfo,
	}

	if e.ContainsTransaction {
		result["account_length"] = e.Address["account_length"]
		result["account_id"] = e.Address["account_id"]
		result["account_index"] = e.Address["account_index"]
		result["account_idx"] = e.Address["account_idx"]
		result["signature"] = e.Signature
		result["nonce"] = e.Nonce
		result["era"] = e.Era
		result["extrinsic_hash"] = e.ExtrinsicHash
	}
	if e.CallIndex != "" {
		result["call_code"] = e.CallIndex
		result["call_module_function"] = e.Call.Name
		result["call_module"] = e.CallModule.Name
	}
	result["params"] = e.Params
	return result
}

type EventsDecoder struct {
	Vec
	Metadata MetadataDecoder          `json:"metadata"`
	Elements []map[string]interface{} `json:"elements"`
}

func (e *EventsDecoder) Init(data ScaleBytes, args []string) {
	e.TypeString = "Vec<EventRecord>"
	var metadata MetadataDecoder
	var subType string
	if len(args) > 0 {
		subType = args[0]
	}
	if len(args) > 1 {
		_ = json.Unmarshal([]byte(args[1]), &metadata)
	}
	e.Metadata = metadata
	e.ScaleDecoder.Init(data, subType)
}

func (e *EventsDecoder) Process() []map[string]interface{} {
	elementCount := int(e.ProcessAndUpdateData("Compact<u32>").Int())
	bm, _ := json.Marshal(e.Metadata)
	for i := 0; i < elementCount; i++ {
		element := e.ProcessAndUpdateData("EventRecord", "", string(bm)).Interface().(map[string]interface{})
		e.Elements = append(e.Elements, element)
	}
	return e.Elements
}

type EventRecord struct {
	ScaleDecoder
	Metadata     MetadataDecoder          `json:"metadata"`
	Phase        int                      `json:"phase"`
	ExtrinsicIdx int                      `json:"extrinsic_idx"`
	Type         string                   `json:"type"`
	Params       []map[string]interface{} `json:"params"`
	Event        MetadataEvents           `json:"event"`
	EventModule  MetadataModules          `json:"event_module"`
}

func (e *EventRecord) Init(data ScaleBytes, args []string) {
	var metadata MetadataDecoder
	var subType string
	if len(args) > 0 {
		subType = args[0]
	}
	if len(args) > 1 {
		_ = json.Unmarshal([]byte(args[1]), &metadata)
	}
	e.Metadata = metadata
	e.ScaleDecoder.Init(data, subType)
}

func (e *EventRecord) Process() map[string]interface{} {
	e.Phase = e.getNextU8()
	if e.Phase == 0 {
		e.ExtrinsicIdx = int(e.ProcessAndUpdateData("U32").Uint())
	}
	e.Type = utiles.BytesToHex(e.getNextBytes(2))
	if e.Metadata.EventIndex[e.Type] != nil {
		eventIndex := e.Metadata.EventIndex[e.Type].(map[string]interface{})
		bc, _ := json.Marshal(eventIndex["call"])
		var event MetadataEvents
		_ = json.Unmarshal(bc, &event)
		e.Event = event
		var EventModule MetadataModules
		bc, _ = json.Marshal(eventIndex["module"])
		_ = json.Unmarshal(bc, &EventModule)
		e.EventModule = EventModule
	}
	for _, argType := range e.Event.Args {
		argTypeObj := e.ProcessAndUpdateData(argType).Interface()
		e.Params = append(e.Params, map[string]interface{}{
			"type":     argType,
			"value":    argTypeObj,
			"valueRaw": "",
		})
	}
	return map[string]interface{}{
		"phase":         e.Phase,
		"extrinsic_idx": e.ExtrinsicIdx,
		"type":          e.Type,
		"module_id":     e.EventModule.Name,
		"event_id":      e.Event.Name,
		"params":        e.Params,
	}

}
