package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeStructStr(t *testing.T) {
	cases := []struct {
		input  string
		output string
	}{
		{"bool", "Bool"},
		{"string", "String"},
		{"BTreeMap<bool,string>", "BTreeMap<Bool,String>"},
		{"Data", "Enum(NULL,String,H256,H256,H256,H256)"},
		{"Call", "Call"},
		{"Bytes", "Bytes"},
		{"i8", "Int<8>"},
		{"i128", "Int<128>"},
		{"U128", "U128"},
		{"compact<U128>", "Compact<U128>"},
		{"Vec<Call>", "Vec<Call>"},
		{"Option<Data>", "Option<Enum(NULL,String,H256,H256,H256,H256)>"},
		{"Seal", "Struct(U32,HexBytes)"},
		{"H256", "H256"},
		{"AssetBalance", "Struct(U64,Bool,Bool)"},
		{"Result<U64,String>", "Results<U64,String>"},
	}
	for _, c := range cases {
		assert.Equal(t, getTypeStructString(c.input, 0), c.output)
	}
}
