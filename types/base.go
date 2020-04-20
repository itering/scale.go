package types

import (
	"fmt"
	"github.com/freehere107/scalecodec/utiles"
	"reflect"
	"strings"
)

type IScaleDecoder interface {
	Init(data ScaleBytes, subType string, arg ...interface{})
	Process()
	buildTypeMapping()
	NextBytes(int) []byte
	GetNextU8() int
	reset()
}

type ScaleDecoder struct {
	TypeString       string            `json:"-"`
	Data             ScaleBytes        `json:"-"`
	RawValue         string            `json:"-"`
	Value            interface{}       `json:"-"`
	SubType          string            `json:"-"`
	TypeMapping      map[string]string `json:"-"`
	StructOrderField []string          `json:"-"`
}

func (s *ScaleDecoder) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.SubType = subType
	s.Data = data
	s.RawValue = ""
	s.Value = nil
	if s.TypeMapping == nil && s.TypeString != "" {
		s.buildTypeMapping()
	}
}

func (s *ScaleDecoder) Process() {}

func (s *ScaleDecoder) NextBytes(length int) []byte {
	data := s.Data.GetNextBytes(length)
	s.RawValue += utiles.BytesToHex(data)
	return data
}

func (s *ScaleDecoder) GetNextU8() int {
	b := s.NextBytes(1)
	return int(b[0])
}

func (s *ScaleDecoder) getNextBool() bool {
	data := s.NextBytes(1)
	return utiles.BytesToHex(data) == "01"
}

func (s *ScaleDecoder) reset() {
	s.Data.Data = []byte{}
	s.Data.Offset = 0
}

func (s *ScaleDecoder) buildTypeMapping() {
	if s.TypeString != "" && string(s.TypeString[0]) == "(" && string(s.TypeString[len(s.TypeString)-1:]) == ")" {
		s.TypeMapping = make(map[string]string)
		for k, v := range strings.Split(s.TypeString[1:len(s.TypeString)-1], ",") {
			s.TypeMapping[fmt.Sprintf("col%d", k+1)] = strings.TrimSpace(v)
			s.StructOrderField = append(s.StructOrderField, fmt.Sprintf("col%d", k+1))
		}
	}
}

func (s *ScaleDecoder) ProcessAndUpdateData(typeString string, args ...string) interface{} {
	r := RuntimeType{}
	c, rc, subType := r.reg().getDecoderClass(typeString)
	if c == nil {
		panic(fmt.Sprintf("not found decoder class %s", typeString))
	}
	method, exist := c.MethodByName("Init")
	if exist == false {
		panic(fmt.Sprintf("%s not implement init function", typeString))
	}

	method.Func.Call([]reflect.Value{rc, reflect.ValueOf(s.Data), reflect.ValueOf(subType), reflect.ValueOf(args)})
	rc.MethodByName("Process").Call(nil)
	s.Data.Offset = int(rc.Elem().FieldByName("Data").FieldByName("Offset").Int())
	s.Data.Data = rc.Elem().FieldByName("Data").FieldByName("Data").Bytes()
	return rc.Elem().FieldByName("Value").Interface()
}
