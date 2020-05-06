package scalecodec

import (
	"errors"
	"github.com/freehere107/go-scale-codec/types"
	"github.com/freehere107/go-scale-codec/utiles"
)

type MetadataDecoder struct {
	types.ScaleDecoder
	Version       string               `json:"version"`
	VersionNumber int                  `json:"version_number"`
	Metadata      types.MetadataStruct `json:"metadata"`
	CodecTypes    []string             `json:"codec_types"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := types.ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, nil)
}

func (m *MetadataDecoder) Process() error {
	magicBytes := m.NextBytes(4)
	if string(magicBytes) == "meta" {
		metadataVersion := utiles.U256(utiles.BytesToHex(m.Data.Data[m.Data.Offset : m.Data.Offset+1]))
		m.Version = m.ProcessAndUpdateData("MetadataVersion").(string)
		m.Metadata = m.ProcessAndUpdateData(m.Version).(types.MetadataStruct)
		m.Metadata.MetadataVersion = int(metadataVersion.Int64())
		m.CodecTypes = utiles.UniqueSlice(types.RuntimeCodecType)
		return nil
	}
	return errors.New("not metadata")

}
