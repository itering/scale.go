package types

import (
	"fmt"
	"github.com/freehere107/go-scale-codec/source"
	"strings"
)

func newStruct(names, typeString []string) *TypeMapping {
	if len(names) != len(typeString) {
		panic("init newStruct names and typeString length not equal")
	}
	return &TypeMapping{Names: names, Types: typeString}
}

func RegCustomTypes(registry map[string]source.TypeStruct) {
	for key, typeStruct := range registry {
		key = strings.ToLower(key)

		switch typeStruct.Type {
		case "string":
			instant := typeRegistry[strings.ToLower(typeStruct.TypeString)]
			if instant != nil {
				typeRegistry[key] = instant
			} else {
				fmt.Println("type", strings.ToLower(typeStruct.TypeString), "not reg")
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
			typeRegistry[key] = &Enum{ValueList: typeStruct.ValueList}
		}
	}
}
