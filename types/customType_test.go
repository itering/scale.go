package types_test

import (
	"github.com/freehere107/scalecodec/source"
	"github.com/freehere107/scalecodec/types"
	"testing"
)

func TestRegCustomTypes(t *testing.T) {
	types.RuntimeType{}.Reg()
	types.RegCustomTypes(source.LoadTypeRegistry("../source/base"))
}
