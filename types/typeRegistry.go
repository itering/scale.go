package types

import (
	"errors"
	"fmt"
	"github.com/freehere107/go-scale-codec/source"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var typeRegistry map[string]interface{}

type RuntimeType struct{}

func (r RuntimeType) Reg() *RuntimeType {
	registry := make(map[string]interface{})
	scales := []interface{}{
		&Null{},
		&U8{},
		&U32{},
		&U64{},
		&Compact{},
		&H256{},
		&Address{},
		&Option{},
		&Struct{},
		&Bytes{},
		&Enum{},
		&Bytes{},
		&Vec{},
		&Set{},
		&CompactU32{},
		&Bool{},
		&StorageHasher{},
		&HexBytes{},
		&Moment{},
		&BlockNumber{},
		&AccountId{},
		&BoxProposal{},
		&Signature{},
		&Era{},
		&Balance{},
		&Index{},
		&SessionIndex{},
		&EraIndex{},
		&ParaId{},
		&LogDigest{},
		&Other{},
		&ChangesTrieRoot{},
		&AuthoritiesChange{},
		&SealV0{},
		&Consensus{},
		&Seal{},
		&PreRuntime{},
		&Exposure{},
		&RawAuraPreDigest{},
		&RawBabePreDigest{},
		&RawBabePreDigestPrimary{},
		&RawBabePreDigestSecondary{},
		&SlotNumber{},
		&BabeBlockWeight{},
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
	registry["compact<moment>"] = &Moment{}
	registry["hash"] = &H256{}
	registry["[u8; 32]"] = &VecU8FixedLength{FixedLength: 32}
	registry["[u8; 16]"] = &VecU8FixedLength{FixedLength: 16}
	registry["[u8; 8]"] = &VecU8FixedLength{FixedLength: 8}
	registry["[u8; 4]"] = &VecU8FixedLength{FixedLength: 4}
	registry["[u8; 256]"] = &VecU8FixedLength{FixedLength: 256}

	typeRegistry = registry

	_, filename, _, _ := runtime.Caller(0)
	RegCustomTypes(source.LoadTypeRegistry(fmt.Sprintf("%s/../source/base", path.Dir(filename))))
	return &r
}

func (r *RuntimeType) getCodecInstant(t string) (reflect.Type, reflect.Value, error) {
	rt := typeRegistry[strings.ToLower(t)]
	if rt == nil {
		return nil, reflect.ValueOf((*error)(nil)).Elem(), errors.New("Scale codec type nil" + t)
	}
	value := reflect.ValueOf(rt)
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}
	p := reflect.New(value.Type())
	p.Elem().Set(value)
	return p.Type(), p, nil
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
