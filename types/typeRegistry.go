package types

import (
	"errors"
	"reflect"
	"regexp"
	"strings"
)

type RuntimeType struct {
	Runtime map[string]interface{}
}

func (r *RuntimeType) reg() *RuntimeType {
	registry := make(map[string]interface{})
	scales := []interface{}{
		&U64{},
		&Address{},
		&Option{},
		&Struct{},
		&Bytes{},
		&Enum{},
		&Bytes{},
		&Vec{},
		&CompactU32{},
		&Bool{},
		&StorageHasher{},
		&HexBytes{},
		&Moment{},
		&Moment{},
		&U32{},
		&BlockNumber{},
		&Compact{},
		&MetadataModuleEvent{},
		&MetadataModuleCallArgument{},
		&MetadataModuleCall{},
		&MetadataV6Decoder{},
		&MetadataV6Module{},
		&MetadataV6ModuleStorage{},
		&MetadataV6ModuleConstants{},
		&MetadataV7Decoder{},
		&MetadataV7Module{},
		&MetadataV7ModuleStorage{},
		&MetadataV7ModuleConstants{},
		&MetadataV7ModuleStorageEntry{},
		&MetadataV8Module{},
		&MetadataV8Decoder{},
		&MetadataV9Decoder{},
		&MetadataV10Decoder{},
		&MetadataV11Decoder{},
		&MetadataModuleError{},
	}
	for _, class := range scales {
		valueOf := reflect.ValueOf(class)
		if valueOf.Type().Kind() == reflect.Ptr {
			registry[strings.ToLower(reflect.Indirect(valueOf).Type().Name())] = class
		} else {
			registry[strings.ToLower(valueOf.Type().Name())] = class
		}
	}
	registry["compact<u32>"] = &CompactU32{}
	r.Runtime = registry
	return r
}

func (r *RuntimeType) getCodecInstant(t string) (reflect.Type, reflect.Value, error) {
	t = strings.ToLower(t)
	rt := r.Runtime[t]
	if rt == nil {
		return nil, reflect.ValueOf((*error)(nil)).Elem(), errors.New("Scalecodec type nil" + t)
	}
	value := reflect.ValueOf(rt)
	return value.Type(), value, nil
}

func (r *RuntimeType) decoderClass(typeString string) (reflect.Type, reflect.Value, string) {
	var typeParts []string
	typeString = ConvertType(typeString)
	if typeString[len(typeString)-1:] == ">" {
		decoderClass, rc, err := r.getCodecInstant(typeString)
		if err == nil {
			return decoderClass, rc, ""
		}
		reg := regexp.MustCompile("^([^<]*)<(.+)>$")
		typeParts = reg.FindStringSubmatch(typeString)
	}
	if len(typeParts) > 0 {
		class, rc, err := r.getCodecInstant(typeParts[1])
		if err == nil {
			return class, rc, typeParts[2]
		}
	} else {
		class, rc, err := r.getCodecInstant(typeString)
		if err == nil {
			return class, rc, ""
		}
	}
	if typeString != "()" && string(typeString[0]) == "(" && string(typeString[len(typeString)-1:]) == ")" {
		decoderClass, rc, _ := r.getCodecInstant("Struct")
		s := rc.Interface().(*Struct)
		s.TypeString = typeString
		s.buildStruct()
		return decoderClass, rc, ""
	}
	return nil, reflect.ValueOf((*error)(nil)).Elem(), ""
}
