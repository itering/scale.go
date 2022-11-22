package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
)

type Enum struct {
	ScaleDecoder
	ValueList []string `json:"value_list"`
	Index     int      `json:"index"`
}

func (e *Enum) WithValueList(valueList []string) *Enum {
	e.ValueList = valueList
	return e
}

func (e *Enum) Init(data scaleBytes.ScaleBytes, option *ScaleDecoderOption) {
	e.Index = 0
	if option != nil && len(e.ValueList) == 0 {
		e.ValueList = option.ValueList
	}
	e.ScaleDecoder.Init(data, option)
}

func (e *Enum) Process() {
	index := utiles.BytesToHex(e.NextBytes(1))
	if utiles.U256(index) != nil {
		e.Index = int(utiles.U256(index).Uint64())
	}
	if e.TypeMapping != nil {
		// check c-like enum
		isCLikeEnum := true
		for _, subType := range e.TypeMapping.Types {
			if !regexp.MustCompile("^[0-9]+$").MatchString(subType) {
				isCLikeEnum = false
				break
			}
		}
		if isCLikeEnum {
			rustEnum := make(map[int]string)
			for index, v := range e.TypeMapping.Names {
				rustEnum[utiles.StringToInt(e.TypeMapping.Types[index])] = v
			}
			e.Value = rustEnum[e.Index]
			return
		}
		if len(e.TypeMapping.Types) <= e.Index {
			panic(fmt.Errorf("type %s index out of range [%d] with length %d", e.TypeName, e.Index, len(e.TypeMapping.Types)))
		}
		if subType := e.TypeMapping.Types[e.Index]; subType != "" {
			// struct subType
			var typeMap [][]string
			if len(subType) > 4 && subType[0:2] == "[[" && json.Unmarshal([]byte(subType), &typeMap) == nil && len(typeMap) > 0 && len(typeMap[0]) == 2 {
				result := make(map[string]interface{})
				for _, v := range typeMap {
					result[v[0]] = e.ProcessAndUpdateData(v[1])
				}
				e.Value = map[string]interface{}{e.TypeMapping.Names[e.Index]: result}
				return
			}
			e.Value = map[string]interface{}{e.TypeMapping.Names[e.Index]: e.ProcessAndUpdateData(subType)}
			return
		}
	}
	if e.ValueList[e.Index] != "" {
		e.Value = e.ValueList[e.Index]
	}
}

func (e *Enum) Encode(data interface{}) string {
	// struct
	if e.TypeMapping != nil {
		isCLikeEnum := true
		for _, subType := range e.TypeMapping.Types {
			if !regexp.MustCompile("^[0-9]+$").MatchString(subType) {
				isCLikeEnum = false
				break
			}
		}
		if isCLikeEnum {
			for k, v := range e.TypeMapping.Names {
				if v == utiles.ToString(data) {
					return utiles.U8Encode(utiles.StringToInt(e.TypeMapping.Types[k]))
				}
			}
		}
		for enumKey, value := range data.(map[string]interface{}) {
			index := 0
			for k, v := range e.TypeMapping.Names {
				if v == enumKey {
					subType := e.TypeMapping.Types[k]
					var typeMap [][]string
					if len(subType) > 4 && subType[0:2] == "[[" && json.Unmarshal([]byte(subType), &typeMap) == nil && len(typeMap) > 0 && len(typeMap[0]) == 2 {
						var raw string
						valueStruct := value.(map[string]interface{})
						for _, st := range typeMap {
							raw += EncodeWithOpt(st[1], valueStruct[st[0]], &ScaleDecoderOption{Spec: e.Spec})
						}
						return raw
					}
					return utiles.U8Encode(index) + EncodeWithOpt(subType, value, &ScaleDecoderOption{Spec: e.Spec})
				}
				index++
			}
		}
	}
	for index, v := range e.ValueList {
		if utiles.ToString(data) == v {
			return utiles.U8Encode(index)
		}
	}
	return ""
}

func (e *Enum) ToString() string {
	if e.TypeMapping != nil {
		return strings.Join(e.TypeMapping.Types, "")
	}
	if e.ValueList != nil {
		return strings.Join(e.ValueList, "")
	}
	return e.ScaleDecoder.TypeString
}
