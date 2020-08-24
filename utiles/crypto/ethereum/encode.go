package ethereum

import (
	"fmt"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/scale.go/utiles/crypto/keccak"
	"strings"
)

func Encode(address string) string {
	bytesAddress := utiles.HexToBytes(address)
	if len(bytesAddress) != 20 {
		return ""
	}
	address = utiles.TrimHex(address)
	result := ""
	hash := utiles.BytesToHex(keccak.Keccak256([]byte(address)))
	for index := 0; index < 40; index++ {
		sum := ""
		if utiles.U256(string(hash[index])).Int64() > 7 {
			sum = strings.ToUpper(string(address[index]))
		} else {
			sum = string(address[index])
		}
		result = fmt.Sprintf("%s%s", result, sum)
	}
	return utiles.AddHex(result)
}
