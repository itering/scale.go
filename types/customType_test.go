package types_test

import (
	"github.com/freehere107/go-scale-codec/types"
	"testing"
)

func TestRegCustomTypes(t *testing.T) {
	types.RuntimeType{}.Reg()
}
