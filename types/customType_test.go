package types_test

import (
	"github.com/freehere107/go-scale-codec/source"
	"github.com/freehere107/go-scale-codec/types"
	"testing"
)

func TestRegCustomTypes(t *testing.T) {
	types.RuntimeType{}.Reg()
	types.RegCustomTypes(source.LoadTypeRegistry("../source/base"))
}
