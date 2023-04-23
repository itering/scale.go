package types

import (
	"fmt"

	"github.com/itering/scale.go/utiles"
)

type FixedU8 struct {
	ScaleDecoder
	FixedLength int
}

func (s *FixedU8) Process() {
	value := s.NextBytes(s.FixedLength)
	if utiles.IsASCII(value) {
		s.Value = string(value)
	} else {
		s.Value = utiles.AddHex(utiles.BytesToHex(value))
	}
}

func (s *FixedU8) Encode(value string) string {
	return utiles.TrimHex(value)
}

func (s *FixedU8) TypeStructString() string {
	return fmt.Sprintf("[%d;U8]", s.FixedLength)
}
