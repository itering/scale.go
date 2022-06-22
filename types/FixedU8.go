package types

import "github.com/itering/scale.go/utiles"

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
