package types

import (
	"math/big"
	"strings"

	"github.com/itering/scale.go/utiles/encointer"
	"github.com/shopspring/decimal"
)

// SubstrateFixedU64 https://github.com/encointer/encointer-js/blob/master/packages/util/src/parserFixPoint.spec.ts#L17
// substrate_fixed:FixedU64
type SubstrateFixedU64 struct {
	U64
}

func (s *SubstrateFixedU64) Process() {
	s.U64.Process()
	value := s.U64.Value
	s.Value = encointer.ParseI32F32(decimal.New(int64(value.(uint64)), 0), 9)
}

type SubstrateFixedI128 struct {
	ScaleDecoder
}

func (s *SubstrateFixedI128) Process() {
	value := s.ProcessAndUpdateData("i128")
	s.Value = encointer.ParseI64F64(decimal.NewFromBigInt(value.(*big.Int), 0), 9)
}

func overriderTypesName(st SiType) string {
	complexPathName := strings.Join(st.Path, ":")
	switch complexPathName {
	case "substrate_fixed:FixedU64":
		if st.Def.Composite != nil && len(st.Def.Composite.Fields) == 1 && st.Def.Composite.Fields[0].TypeName == "u64" {
			return "SubstrateFixedU64"
		}
	case "substrate_fixed:FixedI128":
		if st.Def.Composite != nil && len(st.Def.Composite.Fields) == 1 && st.Def.Composite.Fields[0].TypeName == "i128" {
			return "SubstrateFixedI128"
		}
	}
	return ""
}
