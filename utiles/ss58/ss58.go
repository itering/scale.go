package ss58

import (
	"github.com/freehere107/scalecodec/utiles"
	"github.com/freehere107/scalecodec/utiles/base58"
	"golang.org/x/crypto/blake2b"
)

func SS58Decode(address string) string {
	checksumPrefix := []byte("SS58PRE")
	ss58Format := base58.Decode(address)
	if ss58Format[0] != 42 {
		return ""
	}
	var checksumLength int
	if len(ss58Format) == 35 {
		checksumLength = 2
	} else {
		checksumLength = 1
	}
	bss := ss58Format[0 : len(ss58Format)-checksumLength]
	checksum, _ := blake2b.New(64, []byte{})
	w := append(checksumPrefix[:], bss[:]...)
	checksum.Write(w)
	h := checksum.Sum(nil)
	if utiles.BytesToHex(h[0:checksumLength]) != utiles.BytesToHex(ss58Format[len(ss58Format)-checksumLength:]) {
		return ""
	}
	return utiles.BytesToHex(ss58Format[1:33])
}

func SS58Encode(address string) string {
	checksumPrefix := []byte("SS58PRE")
	addressBytes := utiles.HexToBytes(address)
	var checksumLength int
	if len(addressBytes) == 32 {
		checksumLength = 2
	} else if utiles.IntInSlice(len(addressBytes), []int{1, 2, 4, 8}) {
		checksumLength = 1
	} else {
		return ""
	}
	addressFormat := append(utiles.HexToBytes("2a")[:], addressBytes[:]...)
	checksum, _ := blake2b.New(64, []byte{})
	w := append(checksumPrefix[:], addressFormat[:]...)
	checksum.Write(w)
	h := checksum.Sum(nil)
	b := append(addressFormat[:], h[:checksumLength][:]...)
	return base58.Encode(b)
}
