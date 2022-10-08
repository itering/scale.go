package utiles

import "testing"

func TestIsASCII(t *testing.T) {
	if !IsASCII([]byte("hello world")) {
		t.Error("hello world is ASCII")
	}
	if IsASCII([]byte{5}) {
		t.Error("[]byte{5} not  ASCII")
	}
	if IsASCII(HexToBytes("1c08144900000000000000")) {
		t.Error("hex 1c08144900000000000000 not ASCII")
	}
	if IsASCII(HexToBytes("080b2501")) {
		t.Error("hex 080b2501 not ASCII")
	}
}
