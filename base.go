package scalecodec

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/freehere107/scalecodec/utiles"
	"reflect"
	"regexp"
	"strings"
)

type RuntimeConfiguration struct {
}

var regRuntimeStruct map[string]interface{}

func (r *RuntimeConfiguration) init() {
	regRuntimeStruct = make(map[string]interface{})
	regRuntimeStruct["vec<u8>"] = &Bytes{}
	regRuntimeStruct["enum"] = &Enum{}
	regRuntimeStruct["metadatav4decoder"] = &MetadataV4Decoder{}
	regRuntimeStruct["metadatav4module"] = &MetadataV4Module{}
	regRuntimeStruct["bytes"] = &Bytes{}
	regRuntimeStruct["vec"] = &Vec{}
	regRuntimeStruct["compact<u32>"] = &CompactU32{}
	regRuntimeStruct["bool"] = &Bool{}
	regRuntimeStruct["metadatav4modulestorage"] = &MetadataV4ModuleStorage{}
	regRuntimeStruct["storagehasher"] = &StorageHasher{}
	regRuntimeStruct["hexbytes"] = &HexBytes{}
	regRuntimeStruct["metadatamoduleevent"] = &MetadataModuleEvent{}
	regRuntimeStruct["metadatamodulecallargument"] = &MetadataModuleCallArgument{}
	regRuntimeStruct["metadatamodulecall"] = &MetadataModuleCall{}
	regRuntimeStruct["moment"] = &Moment{}
	regRuntimeStruct["compact<moment>"] = &Moment{}
	regRuntimeStruct["eventrecord"] = &EventRecord{}
	regRuntimeStruct["u32"] = &U32{}
}

func (r *RuntimeConfiguration) GetDecoderClass(typeString string) (interface{}, error) {
	fmt.Println("get decoder class", typeString)
	typeString = strings.ToLower(typeString)
	if regRuntimeStruct[typeString] == nil {
		return nil, errors.New("Scalecodec type nil" + typeString)
	}
	return regRuntimeStruct[typeString], nil
}

type ScaleBytes struct {
	Data   []byte `json:"data"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
}

func (s *ScaleBytes) getNextBytes(length int) []byte {
	data := s.Data[s.Offset : s.Offset+length]
	s.Offset = s.Offset + length
	return data
}

func (s *ScaleBytes) getRemainingLength() int {
	return s.Length - s.Offset
}

func (s *ScaleBytes) String() string {
	return utiles.AddHex(utiles.BytesToHex(s.Data))
}

func (s *ScaleBytes) reset() {
	s.Offset = 0
}

type ScaleDecoder struct {
	TypeString  string                 `json:"type_string"`
	Data        ScaleBytes             `json:"data"`
	RawValue    string                 `json:"raw_value"`
	Value       string                 `json:"value"`
	SubType     string                 `json:"sub_type"`
	TypeMapping map[string]interface{} `json:"type_mapping"`
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
		}
		s.TypeMapping = typeMapping
	}
}

func (s *ScaleDecoder) getNextBytes(length int) []byte {
	data := s.Data.getNextBytes(length)
	s.RawValue += utiles.BytesToHex(data)
	return data
}

func (s *ScaleDecoder) ProcessType(typeString string, valueList ...string) reflect.Value { //todo
	c, subType := GetDecoderClass(typeString)
	if subType != "" {
		valueList = append([]string{subType}, valueList...)
	}
	var tt reflect.Value
	switch c.(type) {
	case *Enum:
		class := c.(*Enum)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV3Decoder:
		class := c.(*MetadataV3Decoder)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV4Decoder:
		class := c.(*MetadataV4Decoder)
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
	case *MetadataV4Module:
		class := c.(*MetadataV4Module)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *Bool:
		class := c.(*Bool)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *MetadataV4ModuleStorage:
		class := c.(*MetadataV4ModuleStorage)
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
	case *EventRecord:
		class := c.(*EventRecord)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	case *U32:
		class := c.(*U32)
		class.ScaleDecoder = *s
		tt = reflect.ValueOf(&class).Elem()
	}
	tt.MethodByName("Init").Call([]reflect.Value{reflect.ValueOf(s.Data), reflect.ValueOf(valueList)})
	return tt
}

func (s *ScaleDecoder) getNextU8() int {
	b := s.getNextBytes(1)
	bs := make([]byte, 4-len(b))
	bs = append(b[:], bs...)
	return int(binary.LittleEndian.Uint32(bs))
}

func (s *ScaleDecoder) getNextBool() bool {
	data := s.getNextBytes(1)
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

func GetDecoderClass(typeString string) (interface{}, string) { //todo
	var typeParts []string
	typeString = ConvertType(typeString)
	r := RuntimeConfiguration{}
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
		//todo
		return decoderClass, ""
	}
	return nil, ""

}

func ConvertType(name string) string {
	name = strings.ReplaceAll(name, "T::", "")
	name = strings.ReplaceAll(name, "<T>", "")
	name = strings.ReplaceAll(name, "<T as Trait>::", "")
	if name == "()" {
		return "Null"
	}
	if name == "Vec<u8>" {
		return "Bytes"
	}
	if name == "<Lookup as StaticLookup>::Source" {
		return "Address"
	}
	if name == "<Balance as HasCompact>::Type" {
		return "Compact<Balance>"
	}
	if name == "<BlockNumber as HasCompact>::Type" {
		return "Compact<BlockNumber>"
	}
	if name == "<Balance as HasCompact>::Type" {
		return "Compact<Balance>"
	}
	if name == "<Moment as HasCompact>::Type" {
		return "Compact<Moment>"
	}
	if name == "<InherentOfflineReport as InherentOfflineReport>::Inherent" {
		return "Inherent"
	}
	return name
}

type ScaleType struct {
	ScaleDecoder
	Metadata MetadataDecoder `json:"metadata"`
}

func (s *ScaleType) Init(data ScaleBytes, args []string) {
	subType := ""
	if len(args) > 0 {
		subType = args[0]
	}
	s.ScaleDecoder.Init(data, subType)
}

func (s *ScaleType) Process() {

}
