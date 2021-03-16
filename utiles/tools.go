package utiles

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

func StringToInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0

}

func IntInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
func AddHex(s string) string {
	if strings.HasPrefix(s, "0x") {
		return s
	}
	return "0x" + s
}

func U256(v string) *big.Int {
	v = strings.TrimPrefix(v, "0x")
	bn := new(big.Int)
	n, _ := bn.SetString(v, 16)
	return n
}

func HexToBytes(s string) []byte {
	s = strings.TrimPrefix(s, "0x")
	c := make([]byte, hex.DecodedLen(len(s)))
	_, _ = hex.Decode(c, []byte(s))
	return c
}

func BytesToHex(b []byte) string {
	c := make([]byte, hex.EncodedLen(len(b)))
	hex.Encode(c, b)
	return string(c)
}

func IntToHex(i interface{}) string {
	return fmt.Sprintf("%x", i)
}

func UniqueSlice(s []string) (list []string) {
	keys := make(map[string]bool)
	for _, entry := range s {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return
}

func ReverseBytes(a []byte) []byte {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func ToString(i interface{}) string {
	var val string
	switch i := i.(type) {
	case string:
		val = i
	case []byte:
		val = string(i)
	default:
		b, _ := json.Marshal(i)
		val = string(b)
	}
	return val
}

func TrimHex(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func BytesToBnHex(value []byte) string {
	var h string
	for _, b := range value {
		h += fmt.Sprintf("%d", int(b))
	}
	return h
}

func Debug(i interface{}) {
	var val string
	switch i := i.(type) {
	case string:
		val = i
	case []byte:
		val = string(i)
	case error:
		val = i.Error()
	default:
		b, _ := json.MarshalIndent(i, "", "  ")
		val = string(b)
	}
	fmt.Println(val)
}

func IsASCII(b []byte) bool {
	for _, c := range b {
		if c > 127 || (c < 32 && !IntInSlice(int(c), []int{9, 10, 13})) {
			return false
		}
	}
	return true
}

func SliceIndex(a string, list []string) int {
	for index, b := range list {
		if b == a {
			return index
		}
	}
	return -1
}
