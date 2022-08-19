package types

type Struct struct {
	ScaleDecoder
}

func (s *Struct) WithTypeMapping(typeMapping *TypeMapping) *Struct {
	s.TypeMapping = typeMapping
	return s
}

func (s *Struct) Process() {
	result := make(map[string]interface{})
	if s.TypeMapping != nil {
		for k, v := range s.TypeMapping.Names {
			result[v] = s.ProcessAndUpdateData(s.TypeMapping.Types[k])
		}
	}
	s.Value = result
}

func (s *Struct) Encode(value map[string]interface{}) string {
	var raw string
	if s.TypeMapping != nil {
		for k, v := range s.TypeMapping.Names {
			raw += Encode(s.TypeMapping.Types[k], value[v])
		}
	}
	return raw
}
