package encointer

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/huandu/xstrings"
	"github.com/shopspring/decimal"
)

// https://github.com/encointer/encointer-js/tree/master/packages/util/src

// Function to produce function to convert fixed-point to Number
//
// Fixed interpretation of u<N> place values
// ... ___ ___ ___ ___ . ___ ___ ___ ___ ...
// ...  8   4   2   1    1/2 1/4 1/8 1/16...
//
// Parameters:
// upper: 0..128 - number of bits in decimal part
// lower: 0..128 - number of bits in fractional part
//
// Produced function parameters:
// raw: substrate_fixed::types::I<upper>F<lower> as I<upper+lower>
// precision: 0..lower number bits in fractional part to process

func parseI32F32(value decimal.Decimal, precision int) decimal.Decimal {
	return parserFixPoint(value, precision, 32, 32)
}

func parseI16F16(value decimal.Decimal, precision int) decimal.Decimal {
	return parserFixPoint(value, precision, 16, 16)
}

func assertLength(upper, lower int) (int, error) {
	length := upper + lower
	if length < 8 || length > 128 {
		return 0, fmt.Errorf("bit length can't be less than 8 or bigger than 128, provided %d", length)
	}
	if length&(length-1) != 0 {
		return 0, fmt.Errorf("bit length should be power of 2, provided %d", length)
	}
	return length, nil
}

func parserFixPoint(value decimal.Decimal, precision int, upper, lower int) decimal.Decimal {
	bn := value.BigInt()
	length, err := assertLength(upper, lower)
	zero := big.NewInt(0)
	if err != nil {
		return decimal.Zero
	}
	if len(fmt.Sprintf("%b", bn)) > length {
		return decimal.Zero
	}

	bits := xstrings.RightJustify(fmt.Sprintf("%b", bn), length, "0")
	var lowerBits string
	if lower > len(bits) {
		lowerBits = xstrings.RightJustify(bits, lower, "0")
	} else {
		lowerBits = substr(bits, lower, lower-precision)
	}

	floatPart := big.NewFloat(0)
	for k, v := range strings.Split(lowerBits, "") {
		if v == "1" {
			floatPart.Add(floatPart, big.NewFloat(0).Quo(big.NewFloat(1), big.NewFloat(float64(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(k+1)), big.NewInt(0)).Int64()))))
		}
	}
	floatValue, _ := floatPart.Float64()
	upperBits := bits[0 : len(bits)-lower]
	decimalPart := zero
	if upperBits != "" {
		decimalPart, _ = big.NewInt(0).SetString(upperBits, 2)
	}
	if bn.Cmp(zero) == 1 {
		return decimal.NewFromBigInt(decimalPart, 0).Add(decimal.NewFromFloat(floatValue))
	}
	return decimal.NewFromBigInt(decimalPart, 0).Sub(decimal.NewFromFloat(floatValue))

}

func substr(input string, length int, endIndex int) string {
	asRunes := []rune(input)
	if endIndex > 0 {
		return string(asRunes[len(asRunes)-length : len(asRunes)-endIndex])
	}
	return string(asRunes[len(asRunes)-length:])
}
