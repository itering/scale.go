package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/freehere107/go-scale-codec/utiles"
	"github.com/freehere107/go-scale-codec/utiles/uint128"
	"github.com/huandu/xstrings"
	"github.com/shopspring/decimal"
	"io"
	"math"
)

type HexBytes struct {
	ScaleDecoder
}

func (h *HexBytes) Process() {
	length := h.ProcessAndUpdateData("Compact<u32>").(int)
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(int(length))))
}

type U8 struct {
	ScaleDecoder
}

func (u *U8) Process() {
	u.Value = u.GetNextU8()
}

type U16 struct {
	Reader io.Reader
	ScaleDecoder
}

func (u *U16) Process() {
	buf := &bytes.Buffer{}
	u.Reader = buf
	_, _ = buf.Write(u.NextBytes(2))
	c := make([]byte, 2)
	_, _ = u.Reader.Read(c)
	u.Value = binary.LittleEndian.Uint16(c)
}

type U32 struct {
	Reader io.Reader
	ScaleDecoder
}

func (u *U32) Process() {
	buf := &bytes.Buffer{}
	u.Reader = buf
	_, _ = buf.Write(u.NextBytes(4))
	c := make([]byte, 4)
	_, _ = u.Reader.Read(c)
	u.Value = binary.LittleEndian.Uint32(c)
}

func (u *U32) Encode(value int) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(value))
	u.Data = ScaleBytes{Data: bs}
}

type U64 struct {
	ScaleDecoder
	Reader io.Reader
}

func (u *U64) Process() {
	buf := &bytes.Buffer{}
	u.Reader = buf
	_, _ = buf.Write(u.NextBytes(8))
	c := make([]byte, 8)
	_, _ = u.Reader.Read(c)
	u.Value = binary.LittleEndian.Uint64(c)
}

type U128 struct {
	ScaleDecoder
}

func (u *U128) Process() {
	if len(u.Data.Data) < 16 {
		u.Data.Data = utiles.HexToBytes(xstrings.LeftJustify(utiles.BytesToHex(u.Data.Data), 32, "0"))
	}
	u.Value = uint128.FromBytes(u.NextBytes(16)).String()
}

type H256 struct {
	ScaleDecoder
}

func (h *H256) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(32)))
}

type H512 struct {
	ScaleDecoder
}

func (h *H512) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(64)))
}

type Era struct {
	ScaleDecoder
}

func (e *Era) Process() {
	optionHex := utiles.BytesToHex(e.NextBytes(1))
	if optionHex == "00" {
		e.Value = optionHex
	} else {
		e.Value = optionHex + utiles.BytesToHex(e.NextBytes(1))
	}
}

type Bool struct {
	ScaleDecoder
}

func (b *Bool) Process() {
	b.Value = b.getNextBool()
}

type Moment struct {
	CompactU32
}

func (m *Moment) Init(data ScaleBytes, option *ScaleDecoderOption) {
	m.TypeString = "Compact<Moment>"
	m.ScaleDecoder.Init(data, option)
}

func (m *Moment) Process() {
	m.CompactU32.Process()
	if m.Value.(int) > 10000000000 {
		m.Value = m.Value.(int) / 1000
	}
}

type Struct struct {
	ScaleDecoder
}

func (s *Struct) Process() {
	result := make(map[string]interface{})
	if s.TypeMapping != nil {
		for k, v := range s.TypeMapping.Names {
			result[v] = s.ProcessAndUpdateData(s.TypeMapping.Types[k])
		}
	}
	s.Value = result
}

type BlockNumber struct {
	U32
}

type Vec struct {
	ScaleDecoder
	Elements []interface{} `json:"elements"`
}

func (v *Vec) Init(data ScaleBytes, option *ScaleDecoderOption) {
	v.Elements = []interface{}{}
	if v.SubType != "" && option != nil {
		option.SubType = v.SubType
	}
	v.ScaleDecoder.Init(data, option)
}

func (v *Vec) Process() {
	elementCount := v.ProcessAndUpdateData("Compact<u32>").(int)
	var result []interface{}
	for i := 0; i < elementCount; i++ {
		element := v.ProcessAndUpdateData(v.SubType)
		v.Elements = append(v.Elements, element)
		result = append(result, element)
	}
	v.Value = result
}

type Address struct {
	ScaleDecoder
	AccountLength string `json:"account_length"`
}

// only support latest address type
func (a *Address) Process() {
	AccountLength := a.NextBytes(1)
	a.AccountLength = utiles.BytesToHex(AccountLength)
	if a.AccountLength == "ff" {
		a.Value = utiles.BytesToHex(a.NextBytes(32))
	}
	a.Value = utiles.BytesToHex(append(AccountLength, a.NextBytes(31)...))
}

type Signature struct {
	ScaleDecoder
}

func (s *Signature) Process() {
	s.Value = utiles.BytesToHex(s.NextBytes(64))
}

type RawAddress struct {
	Address
}

type Enum struct {
	ScaleDecoder
	ValueList []string `json:"value_list"`
	Index     int      `json:"index"`
}

func (e *Enum) Init(data ScaleBytes, option *ScaleDecoderOption) {
	e.Index = 0
	if option != nil && len(e.ValueList) == 0 {
		e.ValueList = option.ValueList
	}
	e.ScaleDecoder.Init(data, option)
}

func (e *Enum) Process() {
	index := utiles.BytesToHex(e.NextBytes(1))
	if utiles.U256(index) != nil {
		e.Index = int(utiles.U256(index).Uint64())
	}
	if e.TypeMapping != nil {
		if e.TypeMapping.Types[e.Index] != "" {
			e.Value = e.ProcessAndUpdateData(e.TypeMapping.Types[e.Index])
			return
		}
	}
	if e.ValueList[e.Index] != "" {
		e.Value = e.ValueList[e.Index]
	}
}

type StorageHasher struct {
	Enum
}

func (s *StorageHasher) Init(data ScaleBytes, option *ScaleDecoderOption) {
	option.ValueList = []string{"Blake2_128", "Blake2_256", "Blake2_128Concat", "Twox128", "Twox256", "Twox64Concat", "Identity"}
	s.Enum.Init(data, option)
}

type Null struct {
	ScaleDecoder
}

type VecU8FixedLength struct {
	ScaleDecoder
	FixedLength int
}

func (s *VecU8FixedLength) Process() {
	s.Value = utiles.AddHex(utiles.BytesToHex(s.NextBytes(s.FixedLength)))
}

type AccountId struct {
	ScaleDecoder
}

func (s *AccountId) Process() {
	s.Value = xstrings.RightJustify(utiles.BytesToHex(s.NextBytes(32)), 64, "0")
}

type BoxProposal struct {
	ScaleDecoder
}

func (s *BoxProposal) Process() {
	callIndex := utiles.BytesToHex(s.NextBytes(2))
	callModule := s.Metadata.CallIndex[callIndex]
	result := map[string]interface{}{
		"call_index":  callIndex,
		"call_name":   callModule.Call.Name,
		"call_module": callModule.Module.Name,
	}
	var param []ExtrinsicParam
	for _, arg := range callModule.Call.Args {
		param = append(param, ExtrinsicParam{
			Name:     arg.Name,
			Type:     arg.Type,
			Value:    s.ProcessAndUpdateData(arg.Type),
			ValueRaw: "",
		})
	}
	s.Value = result
}

type Balance struct{ ScaleDecoder }

func (b *Balance) Process() {
	if len(b.Data.Data) < 16 {
		b.Data.Data = utiles.HexToBytes(xstrings.LeftJustify(utiles.BytesToHex(b.Data.Data), 32, "0"))
	}
	b.Value = decimal.NewFromBigInt(uint128.FromBytes(b.NextBytes(16)).Big(), 0)
}

type Index struct{ U64 }

type SessionIndex struct{ U32 }

type EraIndex struct{ U32 }

type ParaId struct{ U32 }

type Set struct {
	ScaleDecoder
	SetValue  int
	ValueList []string
}

func (s *Set) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.SetValue = 0
	if option.ValueList != nil {
		s.ValueList = option.ValueList
	}
	s.ScaleDecoder.Init(data, option)
}

func (s *Set) Process() {
	s.SetValue = s.ProcessAndUpdateData("U8").(int)
	var result []string
	if s.SetValue > 0 {
		for k, value := range s.ValueList {
			if s.SetValue&int(math.Pow(2, float64(k))) > 0 {
				result = append(result, value)
			}
		}
	}
	s.Value = result
}

type LogDigest struct{ Enum }

func (l *LogDigest) Init(data ScaleBytes, option *ScaleDecoderOption) {
	l.ValueList = []string{"Other", "AuthoritiesChange", "ChangesTrieRoot", "SealV0", "Consensus", "Seal", "PreRuntime"}
	l.Enum.Init(data, option)
}

func (l *LogDigest) Process() {
	index := utiles.BytesToHex(l.NextBytes(1))
	if utiles.U256(index) != nil {
		l.Index = int(utiles.U256(index).Uint64())
	}
	indexType := l.ValueList[l.Index]
	if indexType == "" {
		panic(fmt.Sprintf("LogDigest index %d not in list", l.Index))
	}
	l.Value = map[string]interface{}{
		"type":  indexType,
		"value": l.ProcessAndUpdateData(indexType),
	}
}

type Other struct{ HexBytes }

type ChangesTrieRoot struct{ HexBytes }

type AuthoritiesChange struct{ Vec }

func (l *AuthoritiesChange) Init(data ScaleBytes, option *ScaleDecoderOption) {
	option.SubType = "AccountId"
	l.Vec.Init(data, option)
}

type SealV0 struct{ Struct }

func (s *SealV0) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"slot", "signature"}, Types: []string{"u64", "Signature"}}
	s.Struct.Init(data, option)
}

type Consensus struct{ Struct }

func (s *Consensus) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"engine", "data"}, Types: []string{"u32", "Vec<u8>"}}
	s.Struct.Init(data, option)
}

type Seal struct{ Struct }

func (s *Seal) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"engine", "data"}, Types: []string{"u32", "HexBytes"}}
	s.Struct.Init(data, option)
}

type PreRuntime struct{ Struct }

func (s *PreRuntime) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"engine", "data"}, Types: []string{"u32", "HexBytes"}}
	s.Struct.Init(data, option)
}

type Exposure struct{ Struct }

func (s *Exposure) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"total", "own", "others"}, Types: []string{"Compact<Balance>", "Compact<Balance>", "Vec<IndividualExposure<AccountId, Balance>>"}}
	s.Struct.Init(data, option)
}

type IndividualExposure struct{ Struct }

func (s *IndividualExposure) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"who", "value"}, Types: []string{"AccountId", "Compact<Balance>"}}
	s.Struct.Init(data, option)
}

type RawAuraPreDigest struct{ Struct }

func (s *RawAuraPreDigest) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.Struct.TypeMapping = &TypeMapping{Names: []string{"slotNumber"}, Types: []string{"u64"}}
	s.Struct.Init(data, option)
}

type RawBabePreDigest struct {
	Struct
}

func (r *RawBabePreDigest) Init(data ScaleBytes, option *ScaleDecoderOption) {
	r.Struct.TypeMapping = &TypeMapping{Names: []string{"isPhantom", "Primary", "Secondary"}, Types: []string{"bool", "RawBabePreDigestPrimary", "RawBabePreDigestSecondary"}}
	r.Struct.Init(data, option)
}

func (r *RawBabePreDigest) Process() {
	label := r.ProcessAndUpdateData("Enum", []string{"isPhantom", "Primary", "Secondary"}...).(string)
	for k, name := range r.Struct.TypeMapping.Names {
		if name == label {
			r.Value = map[string]interface{}{label: r.ProcessAndUpdateData(r.TypeMapping.Types[k])}
			break
		}
	}
}

type RawBabePreDigestPrimary struct{ Struct }

type RawBabePreDigestSecondary struct{ Struct }

type SlotNumber struct{ U64 }

type BabeBlockWeight struct{ U32 }

func (r *RawBabePreDigestPrimary) Init(data ScaleBytes, option *ScaleDecoderOption) {
	r.Struct.TypeMapping = &TypeMapping{
		Names: []string{"authorityIndex", "slotNumber", "weight", "vrfOutput", "vrfProof"},
		Types: []string{"u32", "SlotNumber", "BabeBlockWeight", "H256", "H256"},
	}
	r.Struct.Init(data, option)
}

func (r *RawBabePreDigestSecondary) Init(data ScaleBytes, option *ScaleDecoderOption) {
	r.Struct.TypeMapping = &TypeMapping{Names: []string{"authorityIndex", "slotNumber", "weight"}, Types: []string{"u32", "SlotNumber", "BabeBlockWeight"}}
	r.Struct.Init(data, option)
}

type LockIdentifier struct {
	ScaleDecoder
}

func (l *LockIdentifier) Process() {
	l.Value = utiles.AddHex(utiles.BytesToHex(l.NextBytes(8)))
}

type AccountIndex struct{ U32 }

type FixedLengthArray struct {
	ScaleDecoder
	FixedLength int
	SubType     string
}

func (f *FixedLengthArray) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option != nil && option.FixedLength != 0 {
		f.FixedLength = option.FixedLength
	}
	f.ScaleDecoder.Init(data, option)
}

func (f *FixedLengthArray) Process() {
	var result []interface{}
	if f.FixedLength > 0 {
		for i := 0; i < f.FixedLength; i++ {
			result = append(result, f.ProcessAndUpdateData(f.SubType))
		}
	} else {
		f.GetNextU8()
	}
	f.Value = result
}

type AuthorityId struct{ H256 }

type EcdsaSignature struct {
	ScaleDecoder
}

func (e *EcdsaSignature) Process() {
	e.Value = utiles.AddHex(utiles.BytesToHex(e.NextBytes(65)))
}

type Call struct {
	ScaleDecoder
}

func (s *Call) Process() {
	callIndex := utiles.BytesToHex(s.NextBytes(2))
	callModule := s.Metadata.CallIndex[callIndex]
	result := map[string]interface{}{
		"call_index":    callIndex,
		"call_function": callModule.Call.Name,
		"call_module":   callModule.Module.Name,
	}
	var param []ExtrinsicParam
	for _, arg := range callModule.Call.Args {
		param = append(param, ExtrinsicParam{
			Name:     arg.Name,
			Type:     arg.Type,
			Value:    s.ProcessAndUpdateData(arg.Type),
			ValueRaw: s.RawValue,
		})
	}
	result["call_args"] = param
	s.Value = result

}
