package types

import (
	"fmt"
	"github.com/Ompluscator/dynamic-struct"
)

type StructType dynamicstruct.Builder

func NewStruct(names, typeString []string) interface{} {
	instance := dynamicstruct.NewStruct()
	if len(names) != len(typeString) {
		panic("init NewStruct names and typeString length not equal")
	}
	for index, name := range names {
		instance = instance.AddField("Text", "", fmt.Sprintf(`json:"%s" scale:"%s",`, name, typeString[index]))
	}
	return instance.Build().New()
}
