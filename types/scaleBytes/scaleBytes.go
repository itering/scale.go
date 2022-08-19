package scaleBytes

import (
	"github.com/itering/scale.go/utiles"
)

type ScaleBytes struct {
	Data   []byte `json:"data"`
	Offset int    `json:"offset"`
}

func EmptyScaleBytes() ScaleBytes {
	return ScaleBytes{
		Data:   []byte{},
		Offset: 0,
	}
}

func (s *ScaleBytes) GetNextBytes(length int) []byte {
	if s.Offset+length > len(s.Data) {
		data := s.Data[s.Offset:]
		s.Offset = len(s.Data)
		return data
	}
	data := s.Data[s.Offset : s.Offset+length]
	s.Offset = s.Offset + length
	return data
}

func (s *ScaleBytes) GetRemainingLength() int {
	return len(s.Data) - s.Offset
}

func (s *ScaleBytes) String() string {
	return utiles.AddHex(utiles.BytesToHex(s.Data))
}

func (s *ScaleBytes) Reset() {
	s.Offset = 0
}
