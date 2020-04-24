package types

import (
	"fmt"
	"github.com/Ompluscator/dynamic-struct"
	"github.com/freehere107/scalecodec/source"
	"strings"
)

type StructType dynamicstruct.Builder

var customTypeRegistry map[string]interface{}

func newStruct(names, typeString []string) interface{} {
	instance := dynamicstruct.NewStruct()
	if len(names) != len(typeString) {
		panic("init newStruct names and typeString length not equal")
	}
	for index, name := range names {
		instance = instance.AddField("Text", "", fmt.Sprintf(`json:"%s" scale:"%s",`, name, typeString[index]))
	}
	return instance.Build().New()
}

func RegCustomTypes(registry map[string]source.TypeStruct) {
	for key, typeStruct := range registry {
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
			typeRegistry[key] = newStruct(names, typeStrings)
		case "enum":
			typeRegistry[key] = Enum{ValueList: typeStruct.ValueList}
		}
	}
}
