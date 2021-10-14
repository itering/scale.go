package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/utiles"
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
	Type int    `json:"type"`
}

type SiTypeDef struct {
	Composite      *SiTypeDefComposite   `json:"Composite,omitempty"`
	Variant        *SiTypeDefVariant     `json:"Variant,omitempty"`
	Sequence       *SiTypeDefSequence    `json:"Sequence,omitempty"`
	Array          *SiTypeDefArray       `json:"Array,omitempty"`
	Tuple          *SiTypeDefTuple       `json:"Tuple,omitempty"`
	Primitive      *SiTypeDefPrimitive   `json:"Primitive,omitempty"`
	Compact        *SiTypeDefCompact     `json:"Compact,omitempty"`
	BitSequence    *SiTypeDefBitSequence `json:"BitSequence,omitempty"`
	SiTypeDefRange *SiTypeDefRange       `json:"Range,omitempty"`
	// HistoricMetaCompat compatibility
	HistoricMetaCompat string `json:"HistoricMetaCompat,omitempty"`
}

type SiTypeDefComposite struct {
	Fields []SiField `json:"fields"`
}

type SiField struct {
	Name     string   `json:"name,omitempty"`
	Type     int      `json:"type"`
	TypeName string   `json:"typeName"`
	Docs     []string `json:"docs"`
}

type SiTypeDefRange struct {
	Start     int  `json:"start"`
	End       int  `json:"end"`
	Inclusive bool `json:"inclusive"`
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

var (
	registeredSiType  = make(map[string]map[int]string)
	prepareRegistered = make(map[string]source.TypeStruct)
)

func (s *ScaleDecoder) processSiType(id2Portable map[int]SiType, uniqueHash string) {
	registeredSiType[uniqueHash] = make(map[int]string)
	for id, v := range id2Portable {
		if v.Def.Primitive != nil {
			s.dealPrimitiveSiType(id, v, uniqueHash)
		}
	}
	for id, v := range id2Portable {
		if _, ok := registeredSiType[uniqueHash][id]; ok {
			continue
		}
		s.dealOneSiType(id, v, id2Portable, uniqueHash)
	}
	RegCustomTypes(prepareRegistered)
}

func (s *ScaleDecoder) dealPrimitiveSiType(id int, SiTyp SiType, uniqueHash string) {
	if _, ok := TypeRegistry[strings.ToLower(string(*SiTyp.Def.Primitive))]; !ok {
		panic(fmt.Sprintf("lack Primitive type %v", *SiTyp.Def.Primitive))
	}
	registeredSiType[uniqueHash][id] = string(*SiTyp.Def.Primitive)
}

type SiTypeOption struct {
	SkipEnum bool
}

func NameSiType(SiTyp SiType) string {
	if SiTyp.Def.Composite != nil || SiTyp.Def.Variant != nil {
		return strings.ToLower(strings.Join(SiTyp.Path, ":"))
	}
	return ""
}

func (s *ScaleDecoder) dealOneSiType(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	if SiTyp.Def.Composite != nil {
		if len(SiTyp.Def.Composite.Fields) == 1 { // single
			subTypeValue := SiTyp.Def.Composite.Fields[0].Type
			if instant, ok := registeredSiType[uniqueHash][subTypeValue]; ok {
				registeredSiType[uniqueHash][id] = instant
			} else {
				if NameSiType(id2Portable[subTypeValue]) == "" {
					registeredSiType[uniqueHash][id] = s.dealOneSiType(subTypeValue, id2Portable[subTypeValue], id2Portable, uniqueHash)
				} else {
					registeredSiType[uniqueHash][id] = NameSiType(id2Portable[subTypeValue])
				}
			}
			typeString := strings.Join(SiTyp.Path, ":")
			// fmt.Println("typeString", typeString, registeredSiType[uniqueHash][id])
			prepareRegistered[typeString] = source.TypeStruct{V14: true, Type: "string", TypeString: registeredSiType[uniqueHash][id]}
			return registeredSiType[uniqueHash][id]
		} else { // struct
			var types [][]string
			for _, field := range SiTyp.Def.Composite.Fields {
				var typeName string
				if instant, ok := registeredSiType[uniqueHash][field.Type]; ok {
					typeName = instant
				} else {
					if NameSiType(id2Portable[field.Type]) == "" {
						typeName = s.dealOneSiType(field.Type, id2Portable[field.Type], id2Portable, uniqueHash)
					} else {
						typeName = NameSiType(id2Portable[field.Type])
					}
				}
				fieldName := field.Name
				if fieldName == "" {
					fieldName = field.TypeName
				}
				types = append(types, []string{fieldName, typeName})
			}
			typeString := strings.ToLower(strings.Join(SiTyp.Path, ":"))
			prepareRegistered[typeString] = source.TypeStruct{V14: true, Type: "struct", TypeMapping: types}
			registeredSiType[uniqueHash][id] = typeString
			return typeString
		}

	} else if SiTyp.Def.Array != nil {
		subTypeValue := SiTyp.Def.Array.Type
		if SiTyp.Def.Array.Len == 0 {
			registeredSiType[uniqueHash][id] = "NULL"
			return "NULL"
		}
		if instant, ok := registeredSiType[uniqueHash][subTypeValue]; ok {
			registeredSiType[uniqueHash][id] = fmt.Sprintf("[%s; %d]", instant, SiTyp.Def.Array.Len)
		} else {
			registeredSiType[uniqueHash][id] = fmt.Sprintf("[%s; %d]", s.dealOneSiType(subTypeValue, id2Portable[subTypeValue], id2Portable, uniqueHash),
				SiTyp.Def.Array.Len)
		}
		return registeredSiType[uniqueHash][id]

	} else if SiTyp.Def.Sequence != nil {
		subTypeValue := SiTyp.Def.Sequence.Type
		if instant, ok := registeredSiType[uniqueHash][subTypeValue]; ok {
			registeredSiType[uniqueHash][id] = fmt.Sprintf("Vec<%s>", instant)
		} else {
			subTypeName := NameSiType(id2Portable[subTypeValue])
			if subTypeName == "" {
				subTypeName = s.dealOneSiType(subTypeValue, id2Portable[subTypeValue], id2Portable, uniqueHash)
			}
			registeredSiType[uniqueHash][id] = fmt.Sprintf("Vec<%s>", subTypeName)
		}
		return registeredSiType[uniqueHash][id]

	} else if SiTyp.Def.Tuple != nil {
		var tupleSlice []string
		var tupleStruct [][]string
		if len(*SiTyp.Def.Tuple) == 0 {
			return "NULL"
		}
		for index, field := range *SiTyp.Def.Tuple {
			if instant, ok := registeredSiType[uniqueHash][field]; ok {
				tupleSlice = append(tupleSlice, instant)
				tupleStruct = append(tupleStruct, []string{fmt.Sprintf("col%d", index+1), instant})
			} else {
				instant := s.dealOneSiType(field, id2Portable[field], id2Portable, uniqueHash)
				tupleStruct = append(tupleStruct, []string{fmt.Sprintf("col%d", index+1), instant})
				tupleSlice = append(tupleSlice, instant)
			}
		}
		tupleTypeName := fmt.Sprintf("Tuple%s", strings.ToLower(strings.Join(tupleSlice, ":")))
		prepareRegistered[tupleTypeName] = source.TypeStruct{Type: "struct", TypeMapping: tupleStruct, V14: true}
		registeredSiType[uniqueHash][id] = tupleTypeName

		return registeredSiType[uniqueHash][id]

	} else if SiTyp.Def.Compact != nil {
		subTypeValue := SiTyp.Def.Compact.Type
		var compactType string
		if instant, ok := registeredSiType[uniqueHash][subTypeValue]; ok {
			compactType = instant
		} else {
			compactType = s.dealOneSiType(subTypeValue, id2Portable[subTypeValue], id2Portable, uniqueHash)
		}
		registeredSiType[uniqueHash][id] = fmt.Sprintf("compact<%s>", compactType)
		return registeredSiType[uniqueHash][id]
	} else if SiTyp.Def.BitSequence != nil { //
		// https://github.com/polkadot-js/api/blob/662aef203356b8f391415e548e1eb5339f68f828/packages/types/src/generic/PortableRegistry.ts#L293
		registeredSiType[uniqueHash][id] = fmt.Sprintf("BitVec")
		return registeredSiType[uniqueHash][id]
	} else if SiTyp.Def.SiTypeDefRange != nil {
		// todo
	} else if SiTyp.Def.Variant != nil {
		specialVariant := SiTyp.Path[0]

		if strings.EqualFold(specialVariant, "option") {
			subTypeValue := SiTyp.Params[0].Type
			var subType string
			if instant, ok := registeredSiType[uniqueHash][subTypeValue]; ok {
				subType = instant
			} else {
				subType = s.dealOneSiType(subTypeValue, id2Portable[subTypeValue], id2Portable, uniqueHash)
			}
			registeredSiType[uniqueHash][id] = fmt.Sprintf("option<%s>", subType)
			return registeredSiType[uniqueHash][id]
		} else if strings.EqualFold(specialVariant, "Result") {
			// Results<u8, bool>
			resultOk := SiTyp.Params[0].Type
			resultErr := SiTyp.Params[1].Type
			var okType string
			var errType string
			if instant, ok := registeredSiType[uniqueHash][resultOk]; ok {
				okType = instant
			} else {
				okType = s.dealOneSiType(resultOk, id2Portable[resultOk], id2Portable, uniqueHash)
			}
			if instant, ok := registeredSiType[uniqueHash][resultErr]; ok {
				errType = instant
			} else {
				errType = s.dealOneSiType(resultErr, id2Portable[resultErr], id2Portable, uniqueHash)
			}
			registeredSiType[uniqueHash][id] = fmt.Sprintf("Results<%s,%s>", okType, errType)
			return registeredSiType[uniqueHash][id]
		} else if len(SiTyp.Path) >= 2 && (SiTyp.Path[len(SiTyp.Path)-2] == "pallet" && SiTyp.Path[len(SiTyp.Path)-1] == "Call") { // Call Extrinsic
			registeredSiType[uniqueHash][id] = "Call"
			return "Call"
		} else if utiles.SliceIndex(SiTyp.Path[len(SiTyp.Path)-1], []string{"Call", "Event"}) != -1 {
			registeredSiType[uniqueHash][id] = "Call" // tag
			return "Call"
		} else { // enum
			var types, enumValueList [][]string
			ValueEnum := false

			sort.Slice(SiTyp.Def.Variant.Variants[:], func(i, j int) bool {
				return SiTyp.Def.Variant.Variants[i].Index < SiTyp.Def.Variant.Variants[j].Index
			})

			for index, variant := range SiTyp.Def.Variant.Variants {
				typeName := "NULL"
				var StructTypes [][]string
				if len(variant.Fields) == 1 {
					ValueEnum = true
					if instant, ok := registeredSiType[uniqueHash][variant.Fields[0].Type]; ok {
						typeName = instant
					} else {
						if NameSiType(id2Portable[variant.Fields[0].Type]) == "" {
							typeName = s.dealOneSiType(variant.Fields[0].Type, id2Portable[variant.Fields[0].Type], id2Portable, uniqueHash)
						} else {
							typeName = NameSiType(id2Portable[variant.Fields[0].Type])
						}
					}
				} else if len(variant.Fields) > 1 {
					ValueEnum = true
					for i, v := range variant.Fields {
						VName := "NULL"
						if instant, ok := registeredSiType[uniqueHash][v.Type]; ok {
							VName = instant
						} else {
							if NameSiType(id2Portable[v.Type]) == "" {
								typeName = s.dealOneSiType(v.Type, id2Portable[v.Type], id2Portable, uniqueHash)
							} else {
								typeName = NameSiType(id2Portable[v.Type])
							}
						}
						if v.Name == "" {
							v.Name = fmt.Sprintf("col%d", i)
						}
						StructTypes = append(StructTypes, []string{v.Name, ConvertType(VName)})
					}
					if len(StructTypes) > 0 {
						typeName = utiles.ToString(StructTypes)
					}
				} else {
					enumValueList = append(enumValueList, []string{variant.Name, fmt.Sprintf("%d", variant.Index)})
				}
				// fill enum element
				var interval = variant.Index
				if index > 0 {
					interval = variant.Index - SiTyp.Def.Variant.Variants[index-1].Index - 1
				}
				for {
					if interval > 0 {
						types = append(types, []string{"empty", "NULL"})
						interval = interval - 1
					} else {
						break
					}
				}
				types = append(types, []string{variant.Name, typeName})
			}
			typeString := strings.ToLower(strings.Join(SiTyp.Path, ":"))
			if !ValueEnum {
				types = enumValueList
			}
			prepareRegistered[typeString] = source.TypeStruct{Type: "enum", TypeMapping: types, V14: true}
			registeredSiType[uniqueHash][id] = typeString
			return typeString
		}
	} else {
		registeredSiType[uniqueHash][id] = "NULL"
		return "NULL"
	}
	return ""
}
