package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/freehere107/scalecodec/utiles"
	"reflect"
	"regexp"
	"strings"
)

type ScaleDecoder struct {
	TypeString      string                 `json:"type_string"`
	Data            ScaleBytes             `json:"data"`
	RawValue        string                 `json:"raw_value"`
	Value           string                 `json:"value"`
	SubType         string                 `json:"sub_type"`
	TypeMapping     map[string]interface{} `json:"type_mapping"`
	StructTypeOrder []string               `json:"struct_type_order"`
}

func (s *ScaleDecoder) Init(data ScaleBytes, subType string) {
	s.SubType = subType
	if s.TypeMapping == nil && s.TypeString != "" {
		s.buildTypeMapping()
	}
	s.Data = data
	s.RawValue = ""
	s.Value = ""
}

func (s *ScaleDecoder) buildTypeMapping() {
	if s.TypeString != "" && string(s.TypeString[0]) == "(" && string(s.TypeString[len(s.TypeString)-1:]) == ")" {
		typeMapping := make(map[string]interface{})
		n := 1
		for _, v := range strings.Split(s.TypeString[1:len(s.TypeString)-1], ",") {
			typeMapping[fmt.Sprintf("col%d", n)] = strings.TrimSpace(v)
			n += 1
		}
		s.TypeMapping = typeMapping
	}
}

func (s *ScaleDecoder) GetNextBytes(length int) []byte {
	data := s.Data.GetNextBytes(length)
	s.RawValue += utiles.BytesToHex(data)
	return data
}

func (s *ScaleDecoder) ProcessType(typeString string, valueList ...string) reflect.Value { //todo
	c, subType := getDecoderClass(typeString)
	if c == nil {
		panic(fmt.Sprintf("not found decoder class %s", typeString))
	}
	if subType != "" {
		valueList = append([]string{subType}, valueList...)
	}
	var tt reflect.Value
	switch c.(type) {
	case *Enum:
		class := c.(*Enum)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Vec:
		class := c.(*Vec)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Bytes:
		class := c.(*Bytes)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *CompactU32:
		class := c.(*CompactU32)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Bool:
		class := c.(*Bool)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *StorageHasher:
		class := c.(*StorageHasher)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *HexBytes:
		class := c.(*HexBytes)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataModuleEvent:
		class := c.(*MetadataModuleEvent)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataModuleCallArgument:
		class := c.(*MetadataModuleCallArgument)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataModuleCall:
		class := c.(*MetadataModuleCall)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Moment:
		class := c.(*Moment)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *U32:
		class := c.(*U32)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *EraIndex:
		class := c.(*EraIndex)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *BlockNumber:
		class := c.(*BlockNumber)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *AccountId:
		class := c.(*AccountId)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *SessionIndex:
		class := c.(*SessionIndex)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV6Decoder:
		class := c.(*MetadataV6Decoder)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV6Module:
		class := c.(*MetadataV6Module)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV6ModuleStorage:
		class := c.(*MetadataV6ModuleStorage)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV6ModuleConstants:
		class := c.(*MetadataV6ModuleConstants)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV7Decoder:
		class := c.(*MetadataV7Decoder)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV7Module:
		class := c.(*MetadataV7Module)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV7ModuleStorage:
		class := c.(*MetadataV7ModuleStorage)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV7ModuleConstants:
		class := c.(*MetadataV7ModuleConstants)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV7ModuleStorageEntry:
		class := c.(*MetadataV7ModuleStorageEntry)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *StakingLedgers:
		class := c.(*StakingLedgers)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *RingBalanceOf:
		class := c.(*RingBalanceOf)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *KtonBalanceOf:
		class := c.(*KtonBalanceOf)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *UnlockChunk:
		class := c.(*UnlockChunk)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Compact:
		class := c.(*Compact)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *ExtendedBalance:
		class := c.(*ExtendedBalance)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *RegularItem:
		class := c.(*RegularItem)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *StakingBalance:
		class := c.(*StakingBalance)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Keys:
		class := c.(*Keys)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	}
	tt.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(s.Data), reflect.ValueOf(valueList)})
	return tt
}

func (s *ScaleDecoder) GetNextU8() int {
	b := s.GetNextBytes(1)
	data := make([]byte, len(s.Data.Data))
	copy(data, s.Data.Data)
	bs := make([]byte, 4-len(b))
	bs = append(b[:], bs...)
	s.Data.Data = data
	return int(binary.LittleEndian.Uint32(bs))
}

func (s *ScaleDecoder) getNextBool() bool {
	data := s.GetNextBytes(1)
	return utiles.BytesToHex(data) == "01"
}

func (s *ScaleDecoder) Process() {

}

func (s *ScaleDecoder) String() string {
	return s.Value
}

func (s *ScaleDecoder) UpdateData(v reflect.Value) {
	s.Data.Offset = int(v.Elem().FieldByName("Data").FieldByName("Offset").Int())
	s.Data.Data = v.Elem().FieldByName("Data").FieldByName("Data").Bytes()
}

func (s *ScaleDecoder) ProcessAndUpdateData(typeString string, args ...string) reflect.Value {
	v := s.ProcessType(typeString, args...)
	vm := v.MethodByName("Process").Call(nil)[0]
	s.UpdateData(v)
	return vm
}

func getDecoderClass(typeString string) (interface{}, string) { //todo
	var typeParts []string
	typeString = ConvertType(typeString)
	r := RuntimeConfig{}
	r.init()
	if typeString[len(typeString)-1:] == ">" {
		decoderClass, err := r.GetDecoderClass(typeString)
		if err == nil {
			return decoderClass, ""
		}
		reg := regexp.MustCompile("^([^<]*)<(.+)>$")
		typeParts = reg.FindStringSubmatch(typeString)
	}
	if len(typeParts) > 0 {
		class, err := r.GetDecoderClass(typeParts[1])
		if err == nil {
			return class, typeParts[2]
		}
	} else {
		class, err := r.GetDecoderClass(typeString)
		if err == nil {
			return class, ""
		}
	}
	if typeString != "()" && string(typeString[0]) == "(" && string(typeString[len(typeString)-1:]) == ")" {
		decoderClass, _ := r.GetDecoderClass("struct")
		return decoderClass, ""
	}
	return nil, ""
}

type ScaleType struct {
	ScaleDecoder
	//Metadata MetadataDecoder `json:"metadata"`
}

func (s *ScaleType) Init(data ScaleBytes, args []string) {
	subType := ""
	if len(args) > 0 {
		subType = args[0]
	}
	s.ScaleDecoder.Init(data, subType)
}

type MetadataModuleCall struct {
	ScaleType
	Name string                   `json:"name"`
	Args []map[string]interface{} `json:"args"`
	Docs []string                 `json:"docs"`
}

func (m *MetadataModuleCall) Init(data ScaleBytes, valueList []string) {
	subType := ""
	if len(valueList) > 0 {
		subType = valueList[0]
	}
	m.ScaleDecoder.Init(data, subType)
}

func (m *MetadataModuleCall) Process() string {
	m.Name = m.ProcessAndUpdateData("Bytes").String()
	argsValue := m.ProcessAndUpdateData("Vec<MetadataModuleCallArgument>").Interface().([]interface{})
	var args []map[string]interface{}
	for _, v := range argsValue {
		sv := v.(reflect.Value).Interface().(map[string]interface{})
		args = append(args, sv)
	}
	m.Args = args
	docs := m.ProcessAndUpdateData("Vec<Bytes>").Interface().([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(reflect.Value).String())
	}
	r := map[string]interface{}{
		"name": m.Name,
		"args": m.Args,
		"docs": m.Docs,
	}
	br, _ := json.Marshal(r)
	return string(br)
}

type MetadataModuleCallArgument struct {
	ScaleType
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m *MetadataModuleCallArgument) Process() map[string]interface{} {
	m.Name = m.ProcessAndUpdateData("Bytes").String()
	m.Type = ConvertType(m.ProcessAndUpdateData("Bytes").String())
	return map[string]interface{}{
		"name": m.Name,
		"type": m.Type,
	}
}

type MetadataModuleEvent struct {
	ScaleType
	Name string   `json:"name"`
	Docs []string `json:"docs"`
	Args []string `json:"args"`
}

func (m *MetadataModuleEvent) Process() string {
	m.Name = m.ProcessAndUpdateData("Bytes").String()
	args := m.ProcessAndUpdateData("Vec<Bytes>").Interface().([]interface{})
	for _, v := range args {
		m.Args = append(m.Args, v.(reflect.Value).String())
	}
	docs := m.ProcessAndUpdateData("Vec<Bytes>").Interface().([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(reflect.Value).String())
	}
	r := map[string]interface{}{
		"name": m.Name,
		"args": m.Args,
		"docs": m.Docs,
	}
	br, _ := json.Marshal(r)
	return string(br)
}
