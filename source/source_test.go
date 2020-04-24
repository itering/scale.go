package source_test

import (
	"github.com/freehere107/scalecodec/source"
	"testing"
)

func TestLoadTypeRegistry(t *testing.T) {
	fileName := "default"
	source.LoadTypeRegistry(fileName)
}
