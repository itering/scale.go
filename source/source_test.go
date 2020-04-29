package source_test

import (
	"github.com/freehere107/go-scale-codec/source"
	"testing"
)

func TestLoadTypeRegistry(t *testing.T) {
	source.LoadTypeRegistry([]byte(source.BaseType))
}
