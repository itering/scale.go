package scalecodec

import (
	"fmt"

	scaleType "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
)

type EventsDecoder struct {
	scaleType.Vec
	Metadata *scaleType.MetadataStruct
}

func (e *EventsDecoder) Init(data scaleBytes.ScaleBytes, option *scaleType.ScaleDecoderOption) {
	if option.Metadata == nil {
		panic("ExtrinsicDecoder option metadata required")
	}
	e.TypeString = "Vec<EventRecord>"
	e.Metadata = option.Metadata
	e.Vec.Init(data, option)
}

type EventParam struct {
	Type     string      `json:"type"`
	TypeName string      `json:"type_name"`
	Value    interface{} `json:"value"`
}

func (e *EventsDecoder) Process() {
	elementCount := e.ProcessAndUpdateData("Compact<u32>").(int)

	er := EventRecord{Metadata: e.Metadata}
	er.Data = e.Data
	er.Spec = e.Spec
	var result []interface{}
	for i := 0; i < elementCount; i++ {
		element := er.Process()
		element["event_idx"] = i
		result = append(result, element)
	}
	e.Value = result
}

type EventRecord struct {
	scaleType.ScaleDecoder
	Metadata     *scaleType.MetadataStruct
	Phase        int                      `json:"phase"`
	ExtrinsicIdx int                      `json:"extrinsic_idx"`
	Type         string                   `json:"type"`
	Params       []EventParam             `json:"params"`
	Event        scaleType.MetadataEvents `json:"event"`
	Topic        []string                 `json:"topic"`
}

func (e *EventRecord) Process() map[string]interface{} {
	e.Params = []EventParam{}
	e.ExtrinsicIdx = 0
	e.Topic = []string{}

	e.Phase = e.GetNextU8()
	// enum
	// isApplyExtrinsic: bool;
	// asApplyExtrinsic: u32;
	// isFinalization: boolean;
	// isInitialization: boolean;

	if e.Phase == 0 {
		e.ExtrinsicIdx = int(e.ProcessAndUpdateData("U32").(uint32))
	}
	e.Type = utiles.BytesToHex(e.NextBytes(2))

	call, ok := e.Metadata.EventIndex[e.Type]
	if !ok {
		panic(fmt.Sprintf("Not find Event Lookup %s, please check metadata info", e.Type))
	}

	e.Event = call.Call
	e.Module = call.Module.Name
	for index, argType := range e.Event.Args {
		value := e.ProcessAndUpdateData(argType)
		param := EventParam{Type: argType, Value: value}
		if len(e.Event.ArgsTypeName) == len(e.Event.Args) {
			param.TypeName = e.Event.ArgsTypeName[index]
		}
		e.Params = append(e.Params, param)
	}

	if e.Metadata.MetadataVersion >= 5 {
		if topic := e.ProcessAndUpdateData("Vec<Hash>"); topic != nil {
			topicValue := topic.([]interface{})
			for _, v := range topicValue {
				e.Topic = append(e.Topic, v.(string))
			}
		}
	}

	return map[string]interface{}{
		"phase":         e.Phase,
		"extrinsic_idx": e.ExtrinsicIdx,
		"type":          e.Type,
		"module_id":     call.Module.Name,
		"event_id":      e.Event.Name,
		"params":        e.Params,
		"topic":         e.Topic,
	}

}
