package scalecodec

import (
	"errors"

	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
)

type MetadataDecoder struct {
	types.ScaleDecoder
	Version       string               `json:"version"`
	VersionNumber int                  `json:"version_number"`
	Metadata      types.MetadataStruct `json:"metadata"`
	CodecTypes    []string             `json:"codec_types"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := scaleBytes.ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, nil)
}

func (m *MetadataDecoder) Process() error {
	magicBytes := m.NextBytes(4)
	if string(magicBytes) == "meta" {
		metadataVersion := utiles.U256(utiles.BytesToHex(m.Data.Data[m.Data.Offset : m.Data.Offset+1]))
		m.Version = m.ProcessAndUpdateData("MetadataVersion").(string)
		m.Metadata = m.ProcessAndUpdateData(m.Version).(types.MetadataStruct)
		m.Metadata.MetadataVersion = int(metadataVersion.Int64())
		return nil
	}
	return errors.New("not metadata")

}

// CheckRegistry deprecated, metadata v14 will auto register all types
func (m *MetadataDecoder) CheckRegistry() (notReg []string) { return }
