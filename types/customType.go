package types

import (
	"github.com/freehere107/go-scale-codec/source"
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
				typeRegistry[key] = instant
				continue
			}
			// Vec
			if typeString[len(typeString)-1:] == ">" {
				reg := regexp.MustCompile("^([^<]*)<(.+)>$")
				typeParts := reg.FindStringSubmatch(typeString)
				if len(typeParts) > 2 && strings.ToLower(typeParts[0]) == "vec" {
					v := &Vec{}
					v.SubType = typeParts[1]
					typeRegistry[key] = &v
					continue
				}
			}

			// Tuple
			if typeString != "()" && string(typeString[0]) == "(" && string(typeString[len(typeString)-1:]) == ")" {
				s := Struct{}
				s.TypeString = typeString
				s.buildStruct()
				typeRegistry[key] = &s
			}

		case "struct":
			var names, typeStrings []string
			for _, v := range typeStruct.TypeMapping {
				names = append(names, v[0])
				typeStrings = append(typeStrings, v[1])
			}
			s := Struct{}
			s.TypeMapping = newStruct(names, typeStrings)
			typeRegistry[key] = &s

		case "enum":
			var names, typeStrings []string
			for _, v := range typeStruct.TypeMapping {
				names = append(names, v[0])
				typeStrings = append(typeStrings, v[1])
			}
			e := Enum{ValueList: typeStruct.ValueList}
			e.TypeMapping = newStruct(names, typeStrings)
			typeRegistry[key] = &e

		case "set":
			typeRegistry[key] = &Set{ValueList: typeStruct.ValueList}
		}
	}
}
