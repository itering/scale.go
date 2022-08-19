package types

import (
	"fmt"
	"math"

	"github.com/itering/scale.go/types/scaleBytes"
)

type Set struct {
	ScaleDecoder
	SetValue  int
	ValueList []string
	BitLength int
}

func (s *Set) Init(data scaleBytes.ScaleBytes, option *ScaleDecoderOption) {
	s.SetValue = 0
	if option.ValueList != nil {
		s.ValueList = option.ValueList
	}
	s.ScaleDecoder.Init(data, option)
}

func (s *Set) Process() {
	setValue := s.ProcessAndUpdateData(fmt.Sprintf("U%d", s.BitLength))
	switch v := setValue.(type) {
	case uint64:
		s.SetValue = int(v)
	case uint8:
		s.SetValue = int(v)
	case uint32:
		s.SetValue = int(v)
	case uint16:
		s.SetValue = int(v)
	}
	var result []string
	if s.SetValue > 0 {
		for k, value := range s.ValueList {
			if s.SetValue&int(math.Pow(2, float64(k))) > 0 {
				result = append(result, value)
			}
		}
	}
	s.Value = result
}
