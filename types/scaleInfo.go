package types

import (
	"encoding/json"
	"fmt"
	"github.com/itering/scale.go/source"
	"strings"
)

type PortableType struct {
	Id   int    `json:"id"`
	Type SiType `json:"type"`
}

type SiType struct {
	Path   []string          `json:"path"`
	Params []SiTypeParameter `json:"params,omitempty"`
	Def    SiTypeDef         `json:"def"`
	Docs   []string          `json:"docs"`
}

type SiTypeParameter struct {
	Name string `json:"name"`
	Type int    `json:"type,omitempty"`
}

type SiTypeDef struct {
	Composite   *SiTypeDefComposite   `json:"Composite,omitempty"`
	Variant     *SiTypeDefVariant     `json:"Variant,omitempty"`
	Sequence    *SiTypeDefSequence    `json:"Sequence,omitempty"`
	Array       *SiTypeDefArray       `json:"Array,omitempty"`
	Tuple       *SiTypeDefTuple       `json:"Tuple,omitempty"`
	Primitive   *SiTypeDefPrimitive   `json:"Primitive,omitempty"`
	Compact     *SiTypeDefCompact     `json:"Compact,omitempty"`
	BitSequence *SiTypeDefBitSequence `json:"BitSequence,omitempty"`
	// HistoricMetaCompat compatibility
	HistoricMetaCompat string `json:"HistoricMetaCompat,omitempty"`
}

type SiTypeDefComposite struct {
	Fields []SiField `json:"fields"`
}

type SiField struct {
	Name     string   `json:"name"`
	Type     int      `json:"type"`
	TypeName string   `json:"typeName"`
	Docs     []string `json:"docs"`
}

type SiTypeDefVariant struct {
	Variants []SiVariant `json:"variants"`
}

type SiVariant struct {
	Name   string    `json:"name"`
	Fields []SiField `json:"fields"`
	Index  int       `json:"index"`
	Docs   []string  `json:"docs"`
}

type SiTypeDefSequence struct {
	Type int `json:"type"`
}

type SiTypeDefCompact struct {
	Type int `json:"type"`
}

type SiTypeDefArray struct {
	Len  int `json:"len"`
	Type int `json:"type"`
}

type SiTypeDefTuple []int

type SiTypeDefPrimitive string

type SiTypeDefBitSequence struct {
	BitStoreType int `json:"bitStoreType"`
	BitOrderType int `json:"bitOrderType"`
}

func initPortableRaw(raw []interface{}) map[int]SiType {
	var portables []PortableType
	bm, _ := json.Marshal(raw)
	if err := json.Unmarshal(bm, &portables); err != nil {
		panic(err)
	}
	id2Portable := make(map[int]SiType)
	for _, v := range portables {
		id2Portable[v.Id] = v.Type
	}
	return id2Portable
}

type SiTypeWithName struct {
	TypeString string
	Structure  interface{}
}

var registeredSiType = make(map[int]string)

func (s *ScaleDecoder) processSiType(id2Portable map[int]SiType) {
	for id, v := range id2Portable {
		if v.Def.Primitive != nil {
			s.dealPrimitiveSiType(id, v)
		}
	}
	for id, v := range id2Portable {
		if _, ok := registeredSiType[id]; ok {
			continue
		}
		s.dealOneSiType(id, v, id2Portable)
	}
}

func (s *ScaleDecoder) dealPrimitiveSiType(id int, SiTyp SiType) {
	if _, ok := TypeRegistry[strings.ToLower(string(*SiTyp.Def.Primitive))]; !ok {
		panic(fmt.Sprintf("lack Primitive type %v", *SiTyp.Def.Primitive))
	}
	registeredSiType[id] = string(*SiTyp.Def.Primitive)
}

func (s *ScaleDecoder) dealOneSiType(id int, SiTyp SiType, id2Portable map[int]SiType, notVariant ...bool) string {
	if SiTyp.Def.Composite != nil {
		if len(SiTyp.Def.Composite.Fields) == 1 { // single
			subTypeValue := SiTyp.Def.Composite.Fields[0].Type
			if instant, ok := registeredSiType[subTypeValue]; ok {
				registeredSiType[id] = instant
			} else {
				registeredSiType[id] = s.dealOneSiType(subTypeValue, id2Portable[subTypeValue], id2Portable)
			}
			return registeredSiType[id]
		} else { // struct
			var types [][]string
			for _, field := range SiTyp.Def.Composite.Fields {
				var typeName string
				if instant, ok := registeredSiType[field.Type]; ok {
					typeName = instant
				} else {
					typeName = s.dealOneSiType(field.Type, id2Portable[field.Type], id2Portable)
				}
				types = append(types, []string{field.Name, typeName})
			}
			typeString := SiTyp.Path[len(SiTyp.Path)-1]
			RegCustomTypes(map[string]source.TypeStruct{typeString: {Type: "struct", TypeMapping: types}})
			registeredSiType[id] = typeString
			return typeString
		}

	} else if SiTyp.Def.Array != nil {
		subTypeValue := SiTyp.Def.Array.Type

		if instant, ok := registeredSiType[subTypeValue]; ok {
			registeredSiType[id] = fmt.Sprintf("[%s; %d]", instant, SiTyp.Def.Array.Len)
		} else {
			registeredSiType[id] = fmt.Sprintf("[%s; %d]", s.dealOneSiType(id, id2Portable[subTypeValue], id2Portable),
				SiTyp.Def.Array.Len)
		}
		return registeredSiType[id]

	} else if SiTyp.Def.Sequence != nil {
		subTypeValue := SiTyp.Def.Sequence.Type
		if instant, ok := registeredSiType[subTypeValue]; ok {
			registeredSiType[id] = fmt.Sprintf("Vec<%s>", instant)
		} else {
			registeredSiType[id] = fmt.Sprintf("Vec<%s>", s.dealOneSiType(id, id2Portable[subTypeValue], id2Portable))
		}
		return registeredSiType[id]

	} else if SiTyp.Def.Tuple != nil {
		var tupleSlice []string
		if len(*SiTyp.Def.Tuple) == 0 {
			return "Null"
		}
		for _, field := range *SiTyp.Def.Tuple {
			if instant, ok := registeredSiType[field]; ok {
				tupleSlice = append(tupleSlice, instant)
			} else {
				tupleSlice = append(tupleSlice, s.dealOneSiType(id, id2Portable[field], id2Portable))
			}
		}
		registeredSiType[id] = fmt.Sprintf("(%s)", strings.Join(tupleSlice, ","))
		return registeredSiType[id]

	} else if SiTyp.Def.Compact != nil {
		subTypeValue := SiTyp.Def.Compact.Type
		var compactType string
		if instant, ok := registeredSiType[subTypeValue]; ok {
			compactType = instant
		} else {
			compactType = s.dealOneSiType(id, id2Portable[subTypeValue], id2Portable)
		}
		registeredSiType[id] = fmt.Sprintf("compact<%s>", compactType)
		return registeredSiType[id]
	} else if SiTyp.Def.BitSequence != nil { // Todo

	} else if SiTyp.Def.Variant != nil {
		specialVariant := SiTyp.Path[0]

		if strings.EqualFold(specialVariant, "option") {
			subTypeValue := SiTyp.Params[0].Type
			var subType string
			if instant, ok := registeredSiType[subTypeValue]; ok {
				subType = instant
			} else {
				subType = s.dealOneSiType(id, id2Portable[subTypeValue], id2Portable)
			}
			registeredSiType[id] = fmt.Sprintf("option<%s>", subType)
			return registeredSiType[id]
		} else if strings.EqualFold(specialVariant, "result") {
			// Results<u8, bool>
			resultOk := SiTyp.Params[0].Type
			resultErr := SiTyp.Params[1].Type
			var okType string
			var errType string
			if instant, ok := registeredSiType[resultOk]; ok {
				okType = instant
			} else {
				okType = s.dealOneSiType(id, id2Portable[resultOk], id2Portable)
			}
			if instant, ok := registeredSiType[resultErr]; ok {
				errType = instant
			} else {
				errType = s.dealOneSiType(id, id2Portable[resultErr], id2Portable)
			}
			registeredSiType[id] = fmt.Sprintf("Results<%s,%s>", okType, errType)
			return registeredSiType[id]
		} else {
			var types [][]string
			for _, variant := range SiTyp.Def.Variant.Variants {
				typeName := "NULL"
				// if len(variant.Fields) > 0 {
				// 	if instant, ok := registeredSiType[variant.Fields[0].Type]; ok {
				// 		typeName = instant
				// 	} else {
				// 		fmt.Println("get unknown type", variant.Fields[0].TypeName, variant.Fields[0].Type)
				// 		// utiles.Debug(variant.Fields[0])
				// 		// typeName = "NULL"
				// 		typeName = s.dealOneSiType(variant.Fields[0].Type, id2Portable[variant.Fields[0].Type], id2Portable)
				// 	}
				// }
				types = append(types, []string{variant.Name, typeName})
			}
			typeString := SiTyp.Path[len(SiTyp.Path)-1]
			RegCustomTypes(map[string]source.TypeStruct{typeString: {Type: "enum", TypeMapping: types}})
			registeredSiType[id] = typeString
			return typeString
		}

	}
	return ""
}
