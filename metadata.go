package scalecodec

import (
	"encoding/json"
	"github.com/freehere107/scalecodec/types"
)

type MetadataDecoder struct {
	types.ScaleDecoder
	Version    string                 `json:"version"`
	Metadata   string                 `json:"metadata"`
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := types.ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, "")
}

func (m *MetadataDecoder) Process() string {
	magicBytes := m.GetNextBytes(4)
	if string(magicBytes) == "meta" {
		m.Version = m.ProcessAndUpdateData("Enum", "MetadataV0Decoder", "MetadataV1Decoder", "MetadataV2Decoder", "MetadataV3Decoder", "MetadataV4Decoder", "MetadataV5Decoder", "MetadataV6Decoder").String()
		m.Metadata = m.ProcessAndUpdateData(m.Version).String()
		var metadata types.MetadataStruct
		_ = json.Unmarshal([]byte(m.Metadata), &metadata)
		m.CallIndex = metadata.CallIndex
		m.EventIndex = metadata.EventIndex
		bm, _ := json.Marshal(map[string]interface{}{
			"magicNumber": metadata.MagicNumber,
			"metadata":    metadata.Metadata,
		})
		return string(bm)
	} else {
		return ""
	}
}
