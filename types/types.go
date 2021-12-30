package types

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/huandu/xstrings"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/scale.go/utiles/crypto/ethereum"
	"github.com/itering/scale.go/utiles/uint128"
	"github.com/shopspring/decimal"
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

func (u *U32) Encode(value int) string {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(value))
	return utiles.BytesToHex(bs)
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
	elementBytes := u.NextBytes(16)
	if len(elementBytes) < 16 {
		elementBytes = utiles.HexToBytes(xstrings.LeftJustify(utiles.BytesToHex(elementBytes), 32, "0"))
	}
	u.Value = uint128.FromBytes(elementBytes).String()
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

type EraExtrinsic struct {
	ScaleDecoder
}

func (e *EraExtrinsic) Process() {
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

type CompactMoment struct {
	CompactU32
}

func (m *CompactMoment) Process() {
	m.CompactU32.Process()
	if m.Value.(int) > 10000000000 {
		m.Value = m.Value.(int) / 1000
	}
}

type Moment struct {
	U64
}

func (m *Moment) Process() {
	m.U64.Process()
	if m.Value.(uint64) > 10000000000 {
		m.Value = m.Value.(uint64) / 1000
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
}

func (v *Vec) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if v.SubType != "" && option != nil {
		option.SubType = v.SubType
	}
	v.ScaleDecoder.Init(data, option)
}

func (v *Vec) Process() {
	elementCount := v.ProcessAndUpdateData("Compact<u32>").(int)
	var result []interface{}
	if elementCount > 50000 {
		panic(fmt.Sprintf("Vec length %d exceeds %d with subType %s", elementCount, 50000, v.SubType))
	}
	subType := v.SubType
	for i := 0; i < elementCount; i++ {
		result = append(result, v.ProcessAndUpdateData(subType))
	}
	v.Value = result
}

type BoundedVec struct {
	Vec
}

// BoundedVec<Type, Size> to Vec<Type>
// https://github.com/paritytech/substrate/pull/8556
func (v *BoundedVec) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option != nil {
		if BoundedArr := strings.Split(option.SubType, ","); len(BoundedArr) >= 2 {
			size := BoundedArr[len(BoundedArr)-1]
			v.SubType = strings.Replace(option.SubType, fmt.Sprintf(",%s", size), "", 1)
			option.SubType = v.SubType
		}
	}
	v.ScaleDecoder.Init(data, option)
}

// https://github.com/paritytech/substrate/pull/8842
type WeakBoundedVec struct {
	BoundedVec
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
		return
	}
	a.Value = utiles.BytesToHex(append(AccountLength, a.NextBytes(31)...))
}

type GenericAddress struct {
	ScaleDecoder
	AccountLength string `json:"account_length"`
}

func (a *GenericAddress) Process() {
	AccountLength := a.NextBytes(1)
	a.AccountLength = utiles.BytesToHex(AccountLength)
	if a.AccountLength == "ff" {
		a.Value = utiles.BytesToHex(a.NextBytes(32))
		return
	}
	switch a.AccountLength {
	case "fc":
		a.NextBytes(2)
	case "fd":
		a.NextBytes(4)
	case "fe":
		a.NextBytes(8)
	default:
		a.Value = utiles.BytesToHex(append(AccountLength, a.NextBytes(31)...))
	}
}

type Signature struct {
	ScaleDecoder
}

func (s *Signature) Process() {
	s.Value = utiles.BytesToHex(s.NextBytes(64))
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
		// check c-like enum
		isCLikeEnum := true
		for _, subType := range e.TypeMapping.Types {
			if !regexp.MustCompile("^[0-9]+$").MatchString(subType) {
				isCLikeEnum = false
				break
			}
		}
		if isCLikeEnum {
			rustEnum := make(map[int]string)
			for index, v := range e.TypeMapping.Names {
				rustEnum[utiles.StringToInt(e.TypeMapping.Types[index])] = v
			}
			e.Value = rustEnum[e.Index]
			return
		}
		if subType := e.TypeMapping.Types[e.Index]; subType != "" {
			// struct subType
			var typeMap [][]string
			if len(subType) > 4 && subType[0:2] == "[[" && json.Unmarshal([]byte(subType), &typeMap) == nil && len(typeMap) > 0 && len(typeMap[0]) == 2 {
				result := make(map[string]interface{})
				for _, v := range typeMap {
					result[v[0]] = e.ProcessAndUpdateData(v[1])
				}
				e.Value = map[string]interface{}{e.TypeMapping.Names[e.Index]: result}
				return
			}
			e.Value = map[string]interface{}{e.TypeMapping.Names[e.Index]: e.ProcessAndUpdateData(subType)}
			return
		}
	}
	if e.ValueList[e.Index] != "" {
		e.Value = e.ValueList[e.Index]
	}
}

func (e *Enum) Encode(data interface{}) string {
	// struct
	if e.TypeMapping != nil {
		isCLikeEnum := true
		for _, subType := range e.TypeMapping.Types {
			if !regexp.MustCompile("^[0-9]+$").MatchString(subType) {
				isCLikeEnum = false
				break
			}
		}
		if isCLikeEnum {
			for k, v := range e.TypeMapping.Names {
				if v == utiles.ToString(data) {
					return utiles.U8Encode(utiles.StringToInt(e.TypeMapping.Types[k]))
				}
			}
		}
		for enumKey, value := range data.(map[string]interface{}) {
			index := 0
			for k, v := range e.TypeMapping.Names {
				if v == enumKey {
					return utiles.U8Encode(index) + Encode(e.TypeMapping.Types[k], value)
				}
				index++
			}
		}
	}
	for index, v := range e.ValueList {
		if utiles.ToString(data) == v {
			return utiles.U8Encode(index)
		}
	}
	return ""
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
	value := s.NextBytes(s.FixedLength)
	if utiles.IsASCII(value) {
		s.Value = string(value)
	} else {
		s.Value = utiles.AddHex(utiles.BytesToHex(value))
	}
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
	s.Module = callModule.Module.Name
	for _, arg := range callModule.Call.Args {
		param = append(param, ExtrinsicParam{
			Name:  arg.Name,
			Type:  arg.Type,
			Value: s.ProcessAndUpdateData(arg.Type),
		})
	}
	result["params"] = param
	s.Value = result
}

type Balance struct {
	ScaleDecoder
	Reader io.Reader
}

func (b *Balance) Process() {
	buf := &bytes.Buffer{}
	b.Reader = buf
	_, _ = buf.Write(b.NextBytes(16))
	c := make([]byte, 16)
	_, _ = b.Reader.Read(c)
	if utiles.BytesToHex(c) == "ffffffffffffffffffffffffffffffff" {
		b.Value = decimal.Zero
		return
	}
	b.Value = decimal.NewFromBigInt(uint128.FromBytes(c).Big(), 0)
}

type Index struct{ U64 }

type SessionIndex struct{ U32 }

type EraIndex struct{ U32 }

type ParaId struct{ U32 }

type Set struct {
	ScaleDecoder
	SetValue  int
	ValueList []string
	BitLength int
}

func (s *Set) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.SetValue = 0
	if option.ValueList != nil {
		s.ValueList = option.ValueList
	}
	s.ScaleDecoder.Init(data, option)
}

func (s *Set) Process() {
	setValue := s.ProcessAndUpdateData(fmt.Sprintf("U%d", s.BitLength))
	switch v := setValue.(type) {
	case uint64:
		s.SetValue = int(v)
	case uint8:
		s.SetValue = int(v)
	case uint32:
		s.SetValue = int(v)
	case uint16:
		s.SetValue = int(v)
	}
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
	l.ValueList = []string{"Other", "AuthoritiesChange", "ChangesTrieRoot", "SealV0", "Consensus", "Seal", "PreRuntime", "ChangesTrieSignal", "RuntimeEnvironmentUpdated"}
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
	r.Struct.TypeMapping = &TypeMapping{Names: []string{"isPhantom", "Primary", "Secondary", "VRF"}, Types: []string{"bool", "RawBabePreDigestPrimary", "RawBabePreDigestSecondary", "RawBabePreDigestSecondaryVRF"}}
	r.Struct.Init(data, option)
}

func (r *RawBabePreDigest) Process() {
	label := r.ProcessAndUpdateData("RawBabeLabel").(string)
	for k, name := range r.Struct.TypeMapping.Names {
		if name == label {
			r.Value = map[string]interface{}{label: r.ProcessAndUpdateData(r.TypeMapping.Types[k])}
			break
		}
	}
}

type RawBabeLabel struct {
	Enum
}

func (s *RawBabeLabel) Init(data ScaleBytes, option *ScaleDecoderOption) {
	option.ValueList = []string{"isPhantom", "Primary", "Secondary", "VRF"}
	s.Enum.Init(data, option)
}

type RawBabePreDigestPrimary struct{ Struct }

type RawBabePreDigestSecondary struct{ Struct }

type RuntimeEnvironmentUpdated struct{ Null }

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

type RawBabePreDigestSecondaryVRF struct{ Struct }

func (r *RawBabePreDigestSecondaryVRF) Init(data ScaleBytes, option *ScaleDecoderOption) {
	r.Struct.TypeMapping = &TypeMapping{Names: []string{"authorityIndex", "slotNumber", "vrfOutput", "vrfProof"}, Types: []string{"u32", "SlotNumber", "VrfData", "VrfProof"}}
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
		if strings.EqualFold(f.SubType, "u8") {
			value := f.NextBytes(f.FixedLength)
			if utiles.IsASCII(value) {
				f.Value = string(value)
			} else {
				f.Value = utiles.BytesToHex(value)
			}
			return
		}
		for i := 0; i < f.FixedLength; i++ {
			result = append(result, f.ProcessAndUpdateData(f.SubType))
		}
		f.Value = result
	} else {
		f.GetNextU8()
	}
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
		"call_index":  callIndex,
		"call_name":   callModule.Call.Name,
		"call_module": callModule.Module.Name,
	}
	var param []ExtrinsicParam
	s.Module = callModule.Module.Name
	for _, arg := range callModule.Call.Args {
		param = append(param, ExtrinsicParam{
			Name:  arg.Name,
			Type:  arg.Type,
			Value: s.ProcessAndUpdateData(arg.Type),
		})
	}
	result["params"] = param
	s.Value = result

}

type ReferendumIndex struct{ U32 }

type PropIndex struct{ U32 }

type EthereumAddress struct {
	ScaleDecoder
}

func (e *EthereumAddress) Process() {
	e.Value = ethereum.Encode(utiles.BytesToHex(e.NextBytes(20)))
}

type Data struct {
	Enum
}

func (d *Data) Init(data ScaleBytes, option *ScaleDecoderOption) {
	d.TypeMapping = &TypeMapping{
		Names: []string{"None", "Raw", "BlakeTwo256", "Sha256", "Keccak256", "ShaThree256"},
		Types: []string{"Null", "String", "H256", "H256", "H256", "H256"},
	}
	d.Enum.Init(data, option)
}

func (d *Data) Process() {
	c := utiles.BytesToHex(d.NextBytes(1))
	if c == "" || utiles.U256(c).Uint64() == 0 {
		d.Value = map[string]interface{}{"None": nil}
		return
	}

	d.Index = int(utiles.U256(c).Uint64())

	if d.Index >= 1 && d.Index <= 33 {
		data := d.NextBytes(d.Index - 1)
		if utf8.Valid(data) {
			d.Value = map[string]interface{}{"Raw": string(data)}
		} else {
			d.Value = map[string]interface{}{"Raw": utiles.BytesToHex(data)}
		}

	}
	if d.Index >= 34 && d.Index <= 37 {
		enumKey := d.TypeMapping.Names[d.Index-32]
		d.Value = map[string]interface{}{enumKey: d.ProcessAndUpdateData(d.TypeMapping.Types[d.Index-32])}
	}
}

type Vote struct {
	U8
}

type VoteOutcome struct {
	ScaleDecoder
}

func (v *VoteOutcome) Process() {
	v.Value = utiles.BytesToHex(v.NextBytes(32))
}

type H160 struct{ ScaleDecoder }

func (h *H160) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(20)))
}

type IntFixed struct {
	ScaleDecoder
	FixedLength int
	Reader      io.Reader
}

func (f *IntFixed) Init(data ScaleBytes, option *ScaleDecoderOption) {
	if option != nil && option.FixedLength != 0 {
		f.FixedLength = option.FixedLength
	}
	f.ScaleDecoder.Init(data, option)
}

func (f *IntFixed) Process() {
	value := utiles.U256(utiles.BytesToHex(utiles.ReverseBytes(f.NextBytes(f.FixedLength))))
	var i, e, b = big.NewInt(2), big.NewInt(int64(f.FixedLength*8) - 1), big.NewInt(int64(f.FixedLength * 8))
	unsignedMid := big.NewInt(2).Exp(i, e, nil)
	ceiling := big.NewInt(2).Exp(i, b, nil)
	if value.Cmp(unsignedMid) > 0 {
		f.Value = value.Sub(value, ceiling)
		return
	}
	f.Value = value
}

type Key struct{ Bytes }

type OpaqueCall struct {
	Bytes
}

func (f *OpaqueCall) Process() {
	f.Bytes.Process()
	e := ScaleDecoder{}
	option := ScaleDecoderOption{Metadata: f.Metadata, Spec: f.Spec}
	e.Init(ScaleBytes{Data: utiles.HexToBytes(f.Value.(string))}, &option)
	value := e.ProcessAndUpdateData("Call")
	f.Value = value
}

type GenericLookupSource struct {
	ScaleDecoder
}

func (g *GenericLookupSource) Process() {
	if len(g.Data.Data) == 32 {
		g.Value = utiles.BytesToHex(g.NextBytes(32))
	}
	AccountLength := g.NextBytes(1)
	accountLength := utiles.BytesToHex(AccountLength)
	if accountLength == "ff" {
		g.Value = utiles.BytesToHex(g.NextBytes(32))
		return
	}
	offset, length := readLength(AccountLength)
	e := ScaleDecoder{}
	e.Init(ScaleBytes{Data: g.Data.Data[offset : offset+length]}, nil)
	g.Value = strconv.Itoa(int(e.ProcessAndUpdateData("U32").(uint32)))
}

func readLength(value []byte) (int, int) {
	first := value[0]
	switch first {
	case 0xfc:
		return 1, 2
	case 0xfd:
		return 1, 4
	case 0xfe:
		return 1, 8
	}
	return 0, 1
}

type BTreeMap struct{ ScaleDecoder }

func (b *BTreeMap) Process() {
	elementCount := b.ProcessAndUpdateData("Compact<u32>").(int)
	var result []interface{}
	if elementCount > 1000 {
		panic(fmt.Sprintf("BTreeMap length %d exceeds %d", elementCount, 1000))
	}
	for i := 0; i < elementCount; i++ {
		subType := strings.Split(b.SubType, ",")
		key := utiles.ToString(b.ProcessAndUpdateData(subType[0]))
		result = append(result, map[string]interface{}{
			key: b.ProcessAndUpdateData(subType[1]),
		})
	}
	b.Value = result
}

type Box struct{ ScaleDecoder }

func (b *Box) Process() {
	b.Value = b.ProcessAndUpdateData(b.SubType)
}

type BTreeSet struct{ Vec }

type WrapperOpaque struct{ ScaleDecoder }

func (w *WrapperOpaque) Process() {
	w.ProcessAndUpdateData("Compact<u32>")
	w.Value = w.ProcessAndUpdateData(w.SubType)
}

type Range struct{ ScaleDecoder }
type RangeInclusive struct{ Range }

func (r *Range) Process() {
	subTypeSlice := strings.Split(r.SubType, ",")
	start := r.SubType
	end := r.SubType
	if len(subTypeSlice) == 2 {
		start = subTypeSlice[0]
		end = subTypeSlice[1]
	}
	r.Value = map[string]interface{}{
		"start": r.ProcessAndUpdateData(start),
		"end":   r.ProcessAndUpdateData(end),
	}
}
