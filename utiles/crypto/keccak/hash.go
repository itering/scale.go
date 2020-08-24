package keccak

import "github.com/itering/scale.go/pkg/go-ethereum/crypto/sha3"

func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
