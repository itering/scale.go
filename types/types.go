package types

import (
	"encoding/binary"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/freehere107/scalecodec/utiles/uint128"
	"github.com/huandu/xstrings"
	"reflect"
	"strconv"
	"unicode/utf8"
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

type CompactU32 struct {
	Compact
}

func (c *CompactU32) Init(data ScaleBytes, subType string, arg ...interface{}) {
	c.TypeString = "Compact<u32>"
	c.ScaleDecoder.Init(data, "", arg...)
}

func (c *CompactU32) Process() {
	c.ProcessCompactBytes()
	if c.CompactLength <= 4 {
		data := make([]byte, len(c.Data.Data))
		copy(data, c.Data.Data)

		compactBytes := c.CompactBytes
		bs := make([]byte, 4-len(compactBytes))
		compactBytes = append(compactBytes[:], bs...)
		c.Data.Data = data
		c.Value = int(binary.LittleEndian.Uint32(compactBytes)) / 4
	} else {
		c.Value = int(binary.LittleEndian.Uint32(c.CompactBytes))
	}
}

func (c *CompactU32) Encode(value int) ScaleBytes {
	if value <= 63 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2))
		c.Data.Data = bs[0:1]
	} else if value <= 16383 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2)|1)
		c.Data.Data = bs[0:2]
	} else if value <= 1073741823 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2)|2)
		c.Data.Data = bs
	}
	return c.Data
}

type Option struct {
	ScaleDecoder
}

func (o *Option) Process() {
	optionType := o.NextBytes(1)
	if o.SubType != "" && utiles.BytesToHex(optionType) != "00" {
		o.Value = o.ProcessAndUpdateData(o.SubType)
	}
}

type Bytes struct {
	ScaleDecoder
}

func (b *Bytes) Init(data ScaleBytes, subType string, arg ...interface{}) {
	b.TypeString = "Vec<u8>"
	b.ScaleDecoder.Init(data, "", arg...)
}

func (b *Bytes) Process() {
	length := b.ProcessAndUpdateData("Compact<u32>").(int)
	value := b.NextBytes(int(length))
	if utf8.Valid(value) {
		b.Value = string(value)
	} else {
		b.Value = utiles.BytesToHex(value)
	}
}

type OptionBytes struct {
	ScaleDecoder
}

func (b *OptionBytes) Init(data ScaleBytes, subType string, arg ...interface{}) {
	b.TypeString = "Option<Vec<u8>>"
	b.ScaleDecoder.Init(data, "", arg...)
}

func (b *OptionBytes) Process() {
	optionByte := b.NextBytes(1)
	if utiles.BytesToHex(optionByte) != "00" {
		b.Value = b.ProcessAndUpdateData("Bytes").(string)
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

func (m *Moment) Init(data ScaleBytes, subType string, arg ...interface{}) {
	m.TypeString = "Compact<Moment>"
	m.ScaleDecoder.Init(data, subType, arg...)
}

func (m *Moment) Process() {
	intValue := m.ProcessAndUpdateData("Compact<u32>").(int)
	m.Value = intValue
}

type Struct struct {
	ScaleDecoder
}

func (s *Struct) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.ScaleDecoder.Init(data, subType, arg...)
}

func (s *Struct) Process() {
	result := make(map[string]interface{})
	for _, dataType := range s.StructOrderField {
		result[dataType] = s.ProcessAndUpdateData(s.TypeMapping[dataType])
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

func (v *Vec) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.Elements = []interface{}{}
	v.ScaleDecoder.Init(data, subType, arg...)
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

func (e *Enum) Init(data ScaleBytes, subType string, arg ...interface{}) {
	e.Index = 0
	e.ValueList = arg[0].([]string)
	e.ScaleDecoder.Init(data, subType, arg...)
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

func (s *StorageHasher) Init(data ScaleBytes, subType string, arg ...interface{}) {
	valueList := []string{"Blake2_128", "Blake2_256", "Blake2_128Concat", "Twox128", "Twox256", "Twox64Concat", "Identity"}
	s.Enum.Init(data, "", valueList)
}

type DigestItem struct {
	Enum
}

func (s *DigestItem) Init(data ScaleBytes, subType string, arg ...interface{}) {
	valueList := []string{"Other", "AuthoritiesChange", "ChangesTrieRoot", "SealV0", "Consensus", "Seal", "PreRuntime"}
	s.Enum.Init(data, subType, valueList)
}
