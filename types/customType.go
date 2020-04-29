package types

import (
	"github.com/freehere107/go-scale-codec/source"
	"github.com/freehere107/go-scale-codec/utiles"
	"regexp"
	"strings"
)

func newStruct(names, typeString []string) *TypeMapping {
	if len(names) != len(typeString) {
		panic("init newStruct names and typeString length not equal")
	}
	if len(names) == 0 {
		return nil
	}
	return &TypeMapping{Names: names, Types: typeString}
}

func RegCustomTypes(registry map[string]source.TypeStruct) {
	for key, typeStruct := range registry {
		key = strings.ToLower(key)

		switch typeStruct.Type {
		case "string":
			typeString := typeStruct.TypeString
			instant := typeRegistry[strings.ToLower(typeString)]
			if instant != nil {
				regCustomKey(key, instant)
				continue
			}

			// Vec
			if typeString[len(typeString)-1:] == ">" {
				reg := regexp.MustCompile("^([^<]*)<(.+)>$")
				typeParts := reg.FindStringSubmatch(typeString)
				if len(typeParts) > 2 && strings.ToLower(typeParts[1]) == "vec" {
					v := &Vec{}
					v.SubType = typeParts[1]
					regCustomKey(key, &v)
					continue
				}
			}

			// Tuple
			if typeString != "()" && string(typeString[0]) == "(" && string(typeString[len(typeString)-1:]) == ")" {
				s := Struct{}
				s.TypeString = typeString
				s.buildStruct()
				regCustomKey(key, &s)
			}

		case "struct":
			var names, typeStrings []string
			for _, v := range typeStruct.TypeMapping {
				names = append(names, v[0])
				typeStrings = append(typeStrings, v[1])
			}
			s := Struct{}
			s.TypeMapping = newStruct(names, typeStrings)

			regCustomKey(key, &s)
		case "enum":
			var names, typeStrings []string
			for _, v := range typeStruct.TypeMapping {
				names = append(names, v[0])
				typeStrings = append(typeStrings, v[1])
			}
			e := Enum{ValueList: typeStruct.ValueList}
			e.TypeMapping = newStruct(names, typeStrings)
			regCustomKey(key, &e)
		case "set":
			regCustomKey(key, &Set{ValueList: typeStruct.ValueList})
		}
	}
}

func regCustomKey(key string, rt interface{}) {
	slice := strings.Split(key, "#")
	if len(slice) == 2 { // for Special
		special := Special{Registry: rt, Version: []int{0, 99999999}}
		if version := strings.Split(slice[1], "-"); len(version) == 2 {
			special.Version[0] = utiles.StringToInt(version[0])
			if version[1] != "?" {
				special.Version[1] = utiles.StringToInt(version[1])
			}
		}
		if specialRegistry == nil {
			specialRegistry = make(map[string]Special)
		}
		specialRegistry[slice[0]] = special
	} else {
		typeRegistry[key] = rt
	}

}
