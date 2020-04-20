package types

import (
	"encoding/binary"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/freehere107/scalecodec/utiles/uint128"
	"github.com/huandu/xstrings"
	"reflect"
	"strconv"
)

type Compact struct {
	ScaleDecoder
	CompactLength int    `json:"compact_length"`
	CompactBytes  []byte `json:"compact_bytes"`
}

func (c *Compact) ProcessCompactBytes() []byte {
	compactByte := c.NextBytes(1)
	byteMod := compactByte[0] % 4
	if byteMod == 0 {
		c.CompactLength = 1
	} else if byteMod == 1 {
		c.CompactLength = 2
	} else if byteMod == 2 {
		c.CompactLength = 4
	} else {
		c.CompactLength = 5 + ((int(compactByte[0]) - 3) / 4)
	}
	if c.CompactLength == 1 {
		c.CompactBytes = compactByte
	} else if utiles.IntInSlice(c.CompactLength, []int{2, 4}) {
		c.CompactBytes = append(compactByte[:], c.NextBytes(c.CompactLength - 1)[:]...)
	} else {
		c.CompactBytes = c.NextBytes(c.CompactLength - 1)
	}
	return c.CompactBytes
}

func (c *Compact) Process() {
	c.ProcessCompactBytes()
	if c.SubType != "" {
		s := ScaleDecoder{TypeString: c.SubType, Data: ScaleBytes{Data: c.CompactBytes}}
		byteData := s.ProcessAndUpdateData(c.SubType)
		if reflect.TypeOf(byteData).Kind() == reflect.Int && c.CompactLength <= 4 {
			c.Value = []byte(strconv.Itoa(byteData.(int) / 4))
		} else {
			c.Value = byteData
		}
	} else {
		c.Value = c.CompactBytes
	}
}

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
	ScaleDecoder
}

func (u *U16) Process() {
	u.Value = uint16(binary.LittleEndian.Uint16(u.NextBytes(2)))
}

type U32 struct {
	ScaleDecoder
}

func (u *U32) Process() {
	u.Value = uint32(binary.LittleEndian.Uint32(u.NextBytes(4)))
}

func (u *U32) Encode(value int) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(value))
	u.Data = ScaleBytes{Data: bs}
}

type U64 struct {
	ScaleDecoder
}

func (u *U64) Process() {
	u.Value = uint64(binary.LittleEndian.Uint64(u.NextBytes(8)))
}

type U128 struct {
	ScaleDecoder
}

func (u *U128) Process() {
	if len(u.Data.Data) < 16 {
		u.Data.Data = utiles.HexToBytes(xstrings.RightJustify(utiles.BytesToHex(u.Data.Data), 32, "0"))
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
	intValue := m.ProcessAndUpdateData("Compact<u32>").(int)
	m.Value = intValue
}

type Struct struct {
	ScaleDecoder
}

func (s *Struct) Init(data ScaleBytes, option *ScaleDecoderOption) {
	s.ScaleDecoder.Init(data, option)
}

func (s *Struct) Process() {
	result := make(map[string]interface{})
	t := reflect.TypeOf(s)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Tag.Get("scale") != "" && field.Tag.Get("json") != "" {
			result[field.Tag.Get("json")] = s.ProcessAndUpdateData(field.Tag.Get("scale"))
		}
	}
	s.Value = result
}

type BlockNumber struct {
	U64
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
	AccountId     string `json:"account_id"`
	AccountIndex  string `json:"account_index"`
	AccountIdx    string `json:"account_idx"`
}

func (a *Address) Process() {
	AccountLength := a.NextBytes(1)
	a.AccountLength = utiles.BytesToHex(AccountLength)
	if a.AccountLength == "ff" {
		a.AccountId = utiles.BytesToHex(a.NextBytes(32))
	} else {
		var AccountIndex []byte
		if a.AccountLength == "fc" {
			AccountIndex = a.NextBytes(2)
		} else if a.AccountLength == "fd" {
			AccountIndex = a.NextBytes(4)
		} else if a.AccountLength == "fe" {
			AccountIndex = a.NextBytes(8)
		} else {
			AccountIndex = AccountLength
		}
		a.AccountIndex = utiles.BytesToHex(AccountIndex)
		a.AccountIdx = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(AccountIndex)), 10)
	}
	a.Value = map[string]string{"account_length": a.AccountLength, "account_id": a.AccountId, "account_index": a.AccountIndex, "account_idx": a.AccountIdx}
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
	if option != nil {
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
