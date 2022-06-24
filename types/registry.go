package types

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types/convert"
	"github.com/itering/scale.go/types/override"
	"github.com/itering/scale.go/utiles"
)

type RuntimeType struct {
	Module string
}

type Special struct {
	Version  []int
	Registry interface{}
}

var (
	TypeRegistry    map[string]interface{}
	specialRegistry map[string][]Special
)

func HasReg(typeName string) bool {
	_, ok := TypeRegistry[strings.ToLower(typeName)]
	return ok
}

func regDefaultType() {
	registry := make(map[string]interface{})
	scales := []interface{}{
		&Null{},
		&U8{},
		&U16{},
		&U32{},
		&U64{},
		&Float64{},
		&Float32{},
		&U128{},
		&H160{},
		&H256{},
		&H512{},
		&Address{},
		&Option{},
		&Struct{},
		&Enum{},
		&Bytes{},
		&Vec{},
		&BoundedVec{},
		&WeakBoundedVec{},
		&Set{},
		&Compact{},
		&CompactU32{},
		&Bool{},
		&HexBytes{},
		&Moment{},
		&BlockNumber{},
		&AccountId{},
		&BoxProposal{},
		&Signature{},
		&Era{},
		&EraExtrinsic{},
		&Balance{},
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
		&RawBabePreDigestSecondaryVRF{},
		&SlotNumber{},
		&LockIdentifier{},
		&BabeBlockWeight{},
		&Call{},
		&EcdsaSignature{},
		&EthereumAddress{},
		&Data{},
		&VoteOutcome{},
		&RawBabeLabel{},
		&String{},
		&GenericAddress{},
		&OpaqueCall{},
		&BitVec{},
		&MetadataModuleEvent{},
		&MetadataModuleCallArgument{},
		&MetadataModuleCall{},
		&MetadataV13ModuleStorage{},
		&MetadataV7ModuleConstants{},
		&MetadataV13ModuleStorageEntry{},
		&MetadataV8Module{},
		&MetadataV12Decoder{},
		&MetadataV13Decoder{},
		&MetadataV14Decoder{},
		&MetadataV12Module{},
		&MetadataV13Module{},
		&MetadataV14Module{},
		&MetadataV14ModuleStorage{},
		&MetadataV14ModuleStorageEntry{},
		&PalletConstantMetadataV14{},
		&MetadataModuleError{},
		&GenericLookupSource{},
		&BTreeMap{},
		&BTreeSet{},
		&Box{},
		&Result{},
		&RuntimeEnvironmentUpdated{},
		&WrapperOpaque{},
		&Range{},
		&RangeInclusive{},
		&SubstrateFixedU64{},
		&SubstrateFixedI128{},
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
	registry["compact<moment>"] = &CompactMoment{}
	registry["str"] = &String{}
	registry["hash"] = &H256{}
	registry["blockhash"] = &H256{}
	registry["i8"] = &IntFixed{FixedLength: 1}
	registry["i16"] = &IntFixed{FixedLength: 2}
	registry["i32"] = &IntFixed{FixedLength: 4}
	registry["i64"] = &IntFixed{FixedLength: 8}
	registry["i128"] = &IntFixed{FixedLength: 16}
	registry["[u8; 32]"] = &FixedU8{FixedLength: 32}
	registry["[u8; 64]"] = &FixedU8{FixedLength: 64}
	registry["[u8; 65]"] = &FixedU8{FixedLength: 65}
	registry["[u8; 16]"] = &FixedU8{FixedLength: 16}
	registry["[u8; 20]"] = &FixedU8{FixedLength: 20}
	registry["[u8; 8]"] = &FixedU8{FixedLength: 8}
	registry["[u8; 4]"] = &FixedU8{FixedLength: 4}
	registry["[u8; 2]"] = &FixedU8{FixedLength: 2}
	registry["[u8; 256]"] = &FixedU8{FixedLength: 256}
	registry["[u128; 3]"] = &FixedArray{FixedLength: 3, SubType: "u128"}
	TypeRegistry = registry
	// todo change load source pallet type to lazy load
	RegCustomTypes(source.LoadTypeRegistry([]byte(source.BaseType)))
}

func (r *RuntimeType) getCodecInstant(t string, spec int) (reflect.Type, reflect.Value, error) {
	t = override.ModuleType(strings.ToLower(t), r.Module)
	rt, err := r.specialVersionCodec(t, spec)

	if err != nil {
		rt = TypeRegistry[strings.ToLower(t)]
		// fixed array
		if rt == nil && t != "[]" && string(t[0]) == "[" && t[len(t)-1:] == "]" {
			if typePart := strings.Split(t[1:len(t)-1], ";"); len(typePart) >= 2 {
				remainPart := typePart[0 : len(typePart)-1]
				fixed := FixedArray{
					FixedLength: utiles.StringToInt(strings.TrimSpace(typePart[len(typePart)-1])),
					SubType:     strings.TrimSpace(strings.Join(remainPart, ";")),
				}
				rt = &fixed
			}
		}
		if rt == nil {
			return nil, reflect.ValueOf((*error)(nil)).Elem(), errors.New("Scale codec type nil" + t)
		}
	}

	value := reflect.ValueOf(rt)
	if value.Kind() == reflect.Ptr {
		value = reflect.Indirect(value)
	}
	p := reflect.New(value.Type())
	p.Elem().Set(value)
	return p.Type(), p, nil
}

func (r *RuntimeType) GetCodecClass(typeString string, spec int) (reflect.Type, reflect.Value, string) {
	var typeParts []string
	typeString = convert.ConvertType(typeString)

	// complex
	if typeString[len(typeString)-1:] == ">" {
		decoderClass, rc, err := r.getCodecInstant(typeString, spec)
		if err == nil {
			return decoderClass, rc, ""
		}
		reg := regexp.MustCompile("^([^<]*)<(.+)>$")
		typeParts = reg.FindStringSubmatch(typeString)
	}

	if len(typeParts) > 0 {
		class, rc, err := r.getCodecInstant(typeParts[1], spec)
		if err == nil {
			return class, rc, typeParts[2]
		}
	} else {
		class, rc, err := r.getCodecInstant(typeString, spec)
		if err == nil {
			return class, rc, ""
		}
	}

	// Tuple
	if typeString != "()" && string(typeString[0]) == "(" && typeString[len(typeString)-1:] == ")" {
		decoderClass, rc, _ := r.getCodecInstant("Struct", spec)
		s := rc.Interface().(*Struct)
		s.TypeString = typeString
		s.buildStruct()
		return decoderClass, rc, ""
	}

	// namespace
	if strings.Contains(typeString, "::") && typeString != "::" {
		namespaceSlice := strings.Split(typeString, "::")
		return r.GetCodecClass(namespaceSlice[len(namespaceSlice)-1], spec)
	}

	return nil, reflect.ValueOf((*error)(nil)).Elem(), ""
}

func (r *RuntimeType) specialVersionCodec(t string, spec int) (interface{}, error) {
	var rt interface{}

	if specials, ok := specialRegistry[t]; ok {
		for _, special := range specials {
			if spec >= special.Version[0] && spec <= special.Version[1] {
				rt = special.Registry
				return rt, nil
			}
		}
	}
	return rt, fmt.Errorf("not found")
}
