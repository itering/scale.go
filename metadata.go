package scalecodec

import (
	"errors"
	"github.com/freehere107/scalecodec/types"
	"github.com/freehere107/scalecodec/utiles"
)

type MetadataDecoder struct {
	types.ScaleDecoder
	Version       string               `json:"version"`
	VersionNumber int                  `json:"version_number"`
	Metadata      types.MetadataStruct `json:"metadata"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := types.ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, nil)
}

func (m *MetadataDecoder) Process() error {
	magicBytes := m.NextBytes(4)
	if string(magicBytes) == "meta" {
		metadataVersion := utiles.U256(utiles.BytesToHex(m.Data.Data[m.Data.Offset : m.Data.Offset+1]))
		m.Version = m.ProcessAndUpdateData(
			"Enum",
			"MetadataV0Decoder",
			"MetadataV1Decoder",
			"MetadataV2Decoder",
			"MetadataV3Decoder",
			"MetadataV4Decoder",
			"MetadataV5Decoder",
			"MetadataV6Decoder",
			"MetadataV7Decoder",
			"MetadataV8Decoder",
			"MetadataV9Decoder",
			"MetadataV10Decoder",
			"MetadataV11Decoder",
		).(string)
		m.Metadata = m.ProcessAndUpdateData(m.Version).(types.MetadataStruct)
		m.Metadata.MetadataVersion = int(metadataVersion.Int64())
		return nil
	}
	return errors.New("not metadata")
}
