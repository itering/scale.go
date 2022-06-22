package types

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types/convert"
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

func (s *SiType) FindParameter(name string) *SiTypeParameter {
	for _, param := range s.Params {
		if strings.EqualFold(param.Name, name) {
			return &param
		}
	}
	return nil
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

var registeredSiType = make(map[string]map[int]string)

func (s *ScaleDecoder) processSiType(id2Portable map[int]SiType, uniqueHash string) {
	elementCount := len(id2Portable)

	registeredSiType[uniqueHash] = make(map[int]string)
	// deal Primitive
	for i := 0; i < elementCount; i++ {
		if id2Portable[i].Def.Primitive != nil {
			s.dealPrimitiveSiType(i, id2Portable[i], uniqueHash)
		}
	}
	// deal primitive_types
	for i := 0; i < elementCount; i++ {
		if len(id2Portable[i].Path) > 0 && id2Portable[i].Path[0] == "primitive_types" {
			s.expandPrimitiveTypes(i, id2Portable[i], uniqueHash)
		}
	}

	// deal sp_core
	for i := 0; i < elementCount; i++ {
		if len(id2Portable[i].Path) > 0 && id2Portable[i].Path[0] == "sp_core" {
			s.dealOneSiType(i, id2Portable[i], id2Portable, uniqueHash)
		}
	}
	for i := 0; i < elementCount; i++ {
		if _, ok := registeredSiType[uniqueHash][i]; ok {
			continue
		}
		s.dealOneSiType(i, id2Portable[i], id2Portable, uniqueHash)
	}
}

func nameSiType(SiTyp SiType, id int) string {
	var generateName = func(s SiType) string {
		pathName := strings.Join(s.Path, ":")
		if _, ok := V14Types[strings.Join(s.Path, ":")]; ok {
			return fmt.Sprintf("%s@%d", pathName, id)
		}
		return pathName
	}
	if SiTyp.Def.Composite != nil {
		return generateName(SiTyp)
	}
	if SiTyp.Def.Variant != nil {
		specialVariant := SiTyp.Path[0]
		if strings.EqualFold(specialVariant, "Option") || strings.EqualFold(specialVariant, "Result") {
			return ""
		}
		if len(SiTyp.Path) >= 2 && (SiTyp.Path[len(SiTyp.Path)-2] == "pallet" && SiTyp.Path[len(SiTyp.Path)-1] == "Call") {
			return ""
		} else if utiles.SliceIndex(SiTyp.Path[len(SiTyp.Path)-1], []string{"Call", "Event"}) != -1 {
			return ""
		} else {
			return generateName(SiTyp)
		}
	}
	return ""
}

func (s *ScaleDecoder) expandPrimitiveTypes(id int, SiTyp SiType, uniqueHash string) {
	typeString := SiTyp.Path[len(SiTyp.Path)-1]
	if _, ok := TypeRegistry[strings.ToLower(typeString)]; !ok {
		panic(fmt.Sprintf("lack Primitive type %s", SiTyp.Path[len(SiTyp.Path)-1]))
	}
	registeredSiType[uniqueHash][id] = typeString
	RegCustomTypes(map[string]source.TypeStruct{
		fmt.Sprintf("%s:%s", "primitive_types", typeString): {Type: "string", TypeString: typeString, V14: true},
	})
}

func (s *ScaleDecoder) dealPrimitiveSiType(id int, SiTyp SiType, uniqueHash string) {
	if _, ok := TypeRegistry[strings.ToLower(string(*SiTyp.Def.Primitive))]; !ok {
		panic(fmt.Sprintf("lack Primitive type %v", *SiTyp.Def.Primitive))
	}
	typeString := string(*SiTyp.Def.Primitive)
	registeredSiType[uniqueHash][id] = typeString
	RegCustomTypes(map[string]source.TypeStruct{
		fmt.Sprintf("%s:%s", "primitive_types", typeString): {Type: "string", TypeString: typeString, V14: true},
	})
}

func (s *ScaleDecoder) expandComposite(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	// SpRuntimeUncheckedExtrinsic
	if len(SiTyp.Path) >= 2 && SiTyp.Path[0] == "sp_runtime" && SiTyp.Path[len(SiTyp.Path)-1] == "UncheckedExtrinsic" {
		if param := SiTyp.FindParameter("Signature"); param != nil {
			RegCustomTypes(map[string]source.TypeStruct{"ExtrinsicSignature": {Type: "string", TypeString: s.dealOneSiType(param.Type, id2Portable[param.Type], id2Portable, uniqueHash)}})
		}
		if param := SiTyp.FindParameter("Address"); param != nil {
			RegCustomTypes(map[string]source.TypeStruct{"ExtrinsicSigner": {Type: "string", TypeString: s.dealOneSiType(param.Type, id2Portable[param.Type], id2Portable, uniqueHash)}})
		}
	}

	if len(SiTyp.Path) >= 2 && utiles.SliceIndex(SiTyp.Path[len(SiTyp.Path)-1], []string{"WrapperKeepOpaque", "WrapperOpaque"}) != -1 {
		subType := fmt.Sprintf("WrapperOpaque<%s>", s.dealOneSiType(SiTyp.Params[0].Type, id2Portable[SiTyp.Params[0].Type], id2Portable, uniqueHash))
		registeredSiType[uniqueHash][id] = subType
		return registeredSiType[uniqueHash][id]
	}

	if len(SiTyp.Def.Composite.Fields) == 1 {
		subTypeId := SiTyp.Def.Composite.Fields[0].Type
		subType, ok := registeredSiType[uniqueHash][subTypeId]
		if !ok {
			subType = s.dealOneSiType(subTypeId, id2Portable[subTypeId], id2Portable, uniqueHash)
		}
		registeredSiType[uniqueHash][id] = subType
		RegCustomTypes(map[string]source.TypeStruct{nameSiType(SiTyp, id): {Type: "string", TypeString: subType, V14: true}})
		return registeredSiType[uniqueHash][id]
	}
	var typeMapping [][]string
	for _, field := range SiTyp.Def.Composite.Fields {
		subType, ok := registeredSiType[uniqueHash][field.Type]
		if !ok {
			subType = s.dealOneSiType(field.Type, id2Portable[field.Type], id2Portable, uniqueHash, RecursiveOption())
		}
		typeMapping = append(typeMapping, []string{field.Name, subType})
	}
	typeString := nameSiType(SiTyp, id)
	RegCustomTypes(map[string]source.TypeStruct{typeString: {Type: "struct", TypeMapping: typeMapping, V14: true}})
	registeredSiType[uniqueHash][id] = typeString
	return typeString
}

func (s *ScaleDecoder) expandArray(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	subTypeId := SiTyp.Def.Array.Type
	if SiTyp.Def.Array.Len == 0 {
		registeredSiType[uniqueHash][id] = "NULL"
		return "NULL"
	}
	subType, ok := registeredSiType[uniqueHash][subTypeId]
	if !ok {
		subType = s.dealOneSiType(subTypeId, id2Portable[subTypeId], id2Portable, uniqueHash, RecursiveOption())
	}
	registeredSiType[uniqueHash][id] = fmt.Sprintf("[%s; %d]", subType, SiTyp.Def.Array.Len)
	return registeredSiType[uniqueHash][id]
}

func (s *ScaleDecoder) expandSequence(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	subTypeId := SiTyp.Def.Sequence.Type
	subType, ok := registeredSiType[uniqueHash][subTypeId]
	if !ok {
		subType = s.dealOneSiType(subTypeId, id2Portable[subTypeId], id2Portable, uniqueHash, RecursiveOption())
	}
	registeredSiType[uniqueHash][id] = fmt.Sprintf("Vec<%s>", subType)
	return registeredSiType[uniqueHash][id]
}

func (s *ScaleDecoder) expandTuple(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	var tupleSlice []string
	var tupleStruct [][]string
	if len(*SiTyp.Def.Tuple) == 0 {
		return "NULL"
	}
	for index, field := range *SiTyp.Def.Tuple {
		subType, ok := registeredSiType[uniqueHash][field]
		if !ok {
			subType = s.dealOneSiType(field, id2Portable[field], id2Portable, uniqueHash, RecursiveOption())
		}
		tupleStruct = append(tupleStruct, []string{fmt.Sprintf("col%d", index+1), subType})
		tupleSlice = append(tupleSlice, subType)
	}
	tupleTypeName := fmt.Sprintf("Tuple:%s", strings.Join(tupleSlice, ""))
	RegCustomTypes(map[string]source.TypeStruct{tupleTypeName: {Type: "struct", TypeMapping: tupleStruct, V14: true}})
	registeredSiType[uniqueHash][id] = tupleTypeName
	return registeredSiType[uniqueHash][id]
}

func (s *ScaleDecoder) expandCompact(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	subTypeId := SiTyp.Def.Compact.Type
	subType, ok := registeredSiType[uniqueHash][subTypeId]
	if !ok {
		subType = s.dealOneSiType(subTypeId, id2Portable[subTypeId], id2Portable, uniqueHash, RecursiveOption())
	}
	registeredSiType[uniqueHash][id] = fmt.Sprintf("compact<%s>", subType)
	return registeredSiType[uniqueHash][id]
}

func (s *ScaleDecoder) expandOption(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	subTypeId := SiTyp.Params[0].Type
	subType, ok := registeredSiType[uniqueHash][subTypeId]
	if !ok {
		subType = s.dealOneSiType(subTypeId, id2Portable[subTypeId], id2Portable, uniqueHash, RecursiveOption())
	}
	registeredSiType[uniqueHash][id] = fmt.Sprintf("option<%s>", subType)
	return registeredSiType[uniqueHash][id]
}

func (s *ScaleDecoder) expandResult(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	// Result<u8, bool>
	resultOk := SiTyp.Params[0].Type
	resultErr := SiTyp.Params[1].Type
	okType, ok := registeredSiType[uniqueHash][resultOk]
	if !ok {
		okType = s.dealOneSiType(resultOk, id2Portable[resultOk], id2Portable, uniqueHash, RecursiveOption())
	}
	errType, ok := registeredSiType[uniqueHash][resultErr]
	if !ok {
		errType = s.dealOneSiType(resultErr, id2Portable[resultErr], id2Portable, uniqueHash, RecursiveOption())
	}
	registeredSiType[uniqueHash][id] = fmt.Sprintf("Result<%s,%s>", okType, errType)
	return registeredSiType[uniqueHash][id]
}

func (s *ScaleDecoder) expandEnum(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string) string {
	var types, enumValueList [][]string
	var ok bool

	valueEnum := false
	// sort enum index
	sort.Slice(SiTyp.Def.Variant.Variants[:], func(i, j int) bool {
		return SiTyp.Def.Variant.Variants[i].Index < SiTyp.Def.Variant.Variants[j].Index
	})

	// deal enum element
	for index, variant := range SiTyp.Def.Variant.Variants {
		typeName := "NULL"
		var StructTypes [][]string
		switch len(variant.Fields) {
		case 0:
			enumValueList = append(enumValueList, []string{variant.Name, fmt.Sprintf("%d", variant.Index)})
		case 1:
			valueEnum = true
			typeName, ok = registeredSiType[uniqueHash][variant.Fields[0].Type]
			if !ok {
				typeName = s.dealOneSiType(variant.Fields[0].Type, id2Portable[variant.Fields[0].Type], id2Portable, uniqueHash, RecursiveOption())
			}
		default:
			valueEnum = true
			for i, v := range variant.Fields {
				typeName, ok = registeredSiType[uniqueHash][v.Type]
				if !ok {
					typeName = s.dealOneSiType(v.Type, id2Portable[v.Type], id2Portable, uniqueHash, RecursiveOption())
				}
				if v.Name == "" {
					v.Name = fmt.Sprintf("col%d", i)
				}
				StructTypes = append(StructTypes, []string{v.Name, convert.ConvertType(typeName)})
			}
			if len(StructTypes) > 0 {
				typeName = utiles.ToString(StructTypes)
			}
		}

		// fill enum element, avoid empty enum
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
	if !valueEnum { // only value,like {"a":1,"b":2,"c":3}
		types = enumValueList
	}
	typeString := nameSiType(SiTyp, id)
	RegCustomTypes(map[string]source.TypeStruct{typeString: {Type: "enum", TypeMapping: types, V14: true}})
	registeredSiType[uniqueHash][id] = typeString
	return typeString
}

type SiTypeOption struct {
	Recursive bool
}

func RecursiveOption() SiTypeOption {
	return SiTypeOption{Recursive: true}
}

func (s *ScaleDecoder) dealOneSiType(id int, SiTyp SiType, id2Portable map[int]SiType, uniqueHash string, opt ...SiTypeOption) string {
	// Intercept some fixed types
	if ot := overriderTypesName(SiTyp); ot != "" {
		registeredSiType[uniqueHash][id] = ot
		return ot
	}

	if len(opt) > 0 && opt[0].Recursive {
		if siTypName := nameSiType(SiTyp, id); siTypName != "" {
			return siTypName
		}
	}
	if SiTyp.Def.Composite != nil {
		return s.expandComposite(id, SiTyp, id2Portable, uniqueHash)

	} else if SiTyp.Def.Array != nil {
		return s.expandArray(id, SiTyp, id2Portable, uniqueHash)

	} else if SiTyp.Def.Sequence != nil {
		return s.expandSequence(id, SiTyp, id2Portable, uniqueHash)

	} else if SiTyp.Def.Tuple != nil {
		return s.expandTuple(id, SiTyp, id2Portable, uniqueHash)

	} else if SiTyp.Def.Compact != nil {
		return s.expandCompact(id, SiTyp, id2Portable, uniqueHash)

	} else if SiTyp.Def.BitSequence != nil {
		registeredSiType[uniqueHash][id] = "BitVec"
		return registeredSiType[uniqueHash][id]
	} else if SiTyp.Def.SiTypeDefRange != nil { // Todo

	} else if SiTyp.Def.Variant != nil {
		specialVariant := SiTyp.Path[0]
		// Option
		if strings.EqualFold(specialVariant, "Option") {
			// Option
			return s.expandOption(id, SiTyp, id2Portable, uniqueHash)

		} else if strings.EqualFold(specialVariant, "Result") {
			// Result
			return s.expandResult(id, SiTyp, id2Portable, uniqueHash)

		} else if len(SiTyp.Path) >= 2 && (SiTyp.Path[len(SiTyp.Path)-2] == "pallet" && SiTyp.Path[len(SiTyp.Path)-1] == "Call") { // Call Extrinsic
			registeredSiType[uniqueHash][id] = "Call"
			return "Call"
		} else if utiles.SliceIndex(SiTyp.Path[len(SiTyp.Path)-1], []string{"Call", "Event"}) != -1 {
			registeredSiType[uniqueHash][id] = "Call" // tag
			return "Call"
		} else {
			// enum
			return s.expandEnum(id, SiTyp, id2Portable, uniqueHash)
		}
	} else {
		registeredSiType[uniqueHash][id] = "NULL"
		return "NULL"
	}
	return ""
}
