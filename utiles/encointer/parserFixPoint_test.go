package encointer

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/huandu/xstrings"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func Test_bytes(t *testing.T) {
	bn := new(big.Int)
	bn.SetString("110000000000000000000000000000001", 2)
	assert.Equal(t, "180000001", fmt.Sprintf("%x", bn))
	assert.Equal(t, "0000000000000000000000000000000110000000000000000000000000000001", xstrings.RightJustify(fmt.Sprintf("%b", bn), 64, "0"))
	assert.Equal(t, 33, len(fmt.Sprintf("%b", bn)))
	assert.Equal(t, "000001", substr("0000000000000000000000000000000110000000000000000000000000000001", 6, 0))

	floatPart := big.NewFloat(0)
	for k, v := range strings.Split("000001", "") {
		if v == "1" {
			floatPart.Add(floatPart, big.NewFloat(0).Quo(big.NewFloat(1), big.NewFloat(float64(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(k+1)), big.NewInt(0)).Int64()))))
		}
	}
	bits := "0000000000000000000000000000000110000000000000000000000000000001"
	assert.Equal(t, "0000000000000000000000000000000110000000000000000000000000", bits[0:len(bits)-6])
	b, _ := big.NewInt(0).SetString("0000000000000000000000000000000110000000000000000000000000", 2)
	assert.Equal(t, b.String(), "100663296")

	assert.Equal(t, "18.543548583984375", parserFixPoint(decimal.New(79643934720, 0), 18, 32, 32).String())
	assert.Equal(t, "1.9199981689453125", parserFixPoint(decimal.New(8246337024, 0), 18, 32, 32).String())
	assert.Equal(t, "1", parserFixPoint(decimal.New(65536, 0), 0, 16, 16).String())
	assert.Equal(t, "0", parserFixPoint(decimal.New(0, 0), 0, 16, 16).String())
	assert.Equal(t, "0", parserFixPoint(decimal.New(0, 0), 0, 16, 16).String())
	assert.Equal(t, "-18.1233978271484375", parserFixPoint(decimal.RequireFromString("-334317721545467687680"), 18, 64, 64).String())

}
