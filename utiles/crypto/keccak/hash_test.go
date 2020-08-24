package keccak

import (
	"github.com/itering/scale.go/utiles"
	"testing"
)

func TestKeccak256(t *testing.T) {
	input := "test value"
	out := "0x2d07364b5c231c56ce63d49430e085ea3033c750688ba532b24029124c26ca5e"
	if utiles.AddHex(utiles.BytesToHex(Keccak256([]byte(input)))) != out {
		t.Error("TestKeccak256 fail")
	}
}
