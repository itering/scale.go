package types

import (
	"fmt"
	"strings"

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

func (s *FixedU8) Encode(value interface{}) string {
	switch v := value.(type) {
	case string:
		valueStr := v
		if strings.HasPrefix(valueStr, "0x") {
			return utiles.TrimHex(valueStr)
		} else {
			return utiles.BytesToHex([]byte(valueStr))
		}
	case []byte:
		return utiles.TrimHex(utiles.BytesToHex(v))
	default:
		panic("type error,only support string or []byte")
	}
}

func (s *FixedU8) TypeStructString() string {
	return fmt.Sprintf("[%d;U8]", s.FixedLength)
}
