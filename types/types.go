package types

import (
	"bytes"
	"encoding/binary"
	"github.com/freehere107/go-scale-codec/utiles"
	"github.com/freehere107/go-scale-codec/utiles/uint128"
	"github.com/huandu/xstrings"
	"github.com/shopspring/decimal"
	"io"
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
	u.Value = binary.LittleEndian.Uint32(c)
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
	u.Value = uint128.FromBytes(u.NextBytes(16))
}

type H256 struct {
	ScaleDecoder
}

func (h *H256) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.NextBytes(32)))
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
	e.Index = int(utiles.U256(index).Int64())
	if e.ValueList[e.Index] != "" {
		e.Value = e.ValueList[e.Index]
	} else {
		e.Value = ""
	}
}

type StorageHasher struct {
	Enum
}

func (s *StorageHasher) Init(data ScaleBytes, option *ScaleDecoderOption) {
	option.ValueList = []string{"Blake2_128", "Blake2_256", "Blake2_128Concat", "Twox128", "Twox256", "Twox64Concat", "Identity"}
	s.Enum.Init(data, option)
}

type DigestItem struct {
	Enum
}

func (s *DigestItem) Init(data ScaleBytes, option *ScaleDecoderOption) {
	option.ValueList = []string{"Other", "AuthoritiesChange", "ChangesTrieRoot", "SealV0", "Consensus", "Seal", "PreRuntime"}
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
	s.Value = utiles.AddHex(xstrings.RightJustify(utiles.BytesToHex(s.NextBytes(32)), 64, "0"))
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

type Balance struct {
	ScaleDecoder
}

func (b *Balance) Process() {
	if len(b.Data.Data) < 16 {
		b.Data.Data = utiles.HexToBytes(xstrings.LeftJustify(utiles.BytesToHex(b.Data.Data), 32, "0"))
	}
	b.Value = decimal.NewFromBigInt(uint128.FromBytes(b.NextBytes(16)).Big(), 0)
}
