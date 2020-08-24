package ethereum

import (
	"testing"
)

func Test_Encode(t *testing.T) {
	if Encode("0x4119b2e6c3cb618f4f0B93ac77f9Beec7ff02887") != "0x4119b2E6C3Cb618F4f0B93aC77f9Beec7FF02887" {
		t.Error("Test_Encode error", Encode("0x4119b2e6c3cb618f4f0B93ac77f9Beec7ff02887"))
	}
}
