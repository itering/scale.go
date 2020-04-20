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
	r.Runtime = map[string]interface{}{
		"index":                        &U64{},
		"u64":                          &U64{},
		"address":                      &Address{},
		"option":                       &Option{},
		"struct":                       &Struct{},
		"vec<u8>":                      &Bytes{},
		"enum":                         &Enum{},
		"bytes":                        &Bytes{},
		"vec":                          &Vec{},
		"compact<u32>":                 &CompactU32{},
		"bool":                         &Bool{},
		"storagehasher":                &StorageHasher{},
		"hexbytes":                     &HexBytes{},
		"moment":                       &Moment{},
		"compact<moment>":              &Moment{},
		"u32":                          &U32{},
		"blocknumber":                  &BlockNumber{},
		"compact":                      &Compact{},
		"metadatamoduleevent":          &MetadataModuleEvent{},
		"metadatamodulecallargument":   &MetadataModuleCallArgument{},
		"metadatamodulecall":           &MetadataModuleCall{},
		"metadatav6decoder":            &MetadataV6Decoder{},
		"metadatav6module":             &MetadataV6Module{},
		"metadatav6modulestorage":      &MetadataV6ModuleStorage{},
		"metadatav6moduleconstants":    &MetadataV6ModuleConstants{},
		"metadatav7decoder":            &MetadataV7Decoder{},
		"metadatav7module":             &MetadataV7Module{},
		"metadatav7modulestorage":      &MetadataV7ModuleStorage{},
		"metadatav7moduleconstants":    &MetadataV7ModuleConstants{},
		"metadatav7modulestorageentry": &MetadataV7ModuleStorageEntry{},
		"metadatav8module":             &MetadataV8Module{},
		"metadatav8decoder":            &MetadataV8Decoder{},
		"metadatav9decoder":            &MetadataV9Decoder{},
		"metadatav10decoder":           &MetadataV10Decoder{},
		"metadatav11decoder":           &MetadataV11Decoder{},
		"metadatamoduleerror":          &MetadataModuleError{},
	}
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

func (r *RuntimeType) getDecoderClass(typeString string) (reflect.Type, reflect.Value, string) {
	var typeParts []string
	typeString = ConvertType(typeString)
	if typeString[len(typeString)-1:] == ">" {
		decoderClass, rcvr, err := r.getCodecInstant(typeString)
		if err == nil {
			return decoderClass, rcvr, ""
		}
		reg := regexp.MustCompile("^([^<]*)<(.+)>$")
		typeParts = reg.FindStringSubmatch(typeString)
	}
	if len(typeParts) > 0 {
		class, rcvr, err := r.getCodecInstant(typeParts[1])
		if err == nil {
			return class, rcvr, typeParts[2]
		}
	} else {
		class, rcvr, err := r.getCodecInstant(typeString)
		if err == nil {
			return class, rcvr, ""
		}
	}
	if typeString != "()" && string(typeString[0]) == "(" && string(typeString[len(typeString)-1:]) == ")" {
		decoderClass, rcvr, _ := r.getCodecInstant("Struct")
		s := rcvr.Interface().(*Struct)
		s.TypeString = typeString
		s.buildTypeMapping()
		return decoderClass, rcvr, ""
	}
	return nil, reflect.ValueOf((*error)(nil)).Elem(), ""
}
