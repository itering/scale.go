package scalecodec

import (
	"encoding/binary"
	"github.com/freehere107/scalecodec/utiles"
	"github.com/freehere107/scalecodec/utiles/uint128"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"
)

type Compact struct {
	ScaleType
	CompactLength int    `json:"compact_length"`
	CompactBytes  []byte `json:"compact_bytes"`
}

func (c *Compact) ProcessCompactBytes() []byte {
	compactByte := c.getNextBytes(1)
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
		c.CompactBytes = append(compactByte[:], c.getNextBytes(c.CompactLength - 1)[:]...)
	} else {
		c.CompactBytes = c.getNextBytes(c.CompactLength - 1)
	}
	return c.CompactBytes
}

func (c *Compact) Process() []byte {
	c.ProcessCompactBytes()
	if c.SubType != "" {
		s := ScaleDecoder{TypeString: c.SubType, Data: ScaleBytes{Data: c.CompactBytes}}
		byteData := s.ProcessAndUpdateData(c.SubType).Interface()
		if reflect.TypeOf(byteData).Kind() == reflect.Int && c.CompactLength <= 4 {
			return []byte(strconv.Itoa(byteData.(int) / 4))
		} else {
			return byteData.([]byte)
		}

	} else {
		return c.CompactBytes
	}
}

type CompactU32 struct {
	Compact
}

func (c *CompactU32) Init(data ScaleBytes, args []string) {
	c.TypeString = "Compact<u32>"
	c.ScaleDecoder.Init(data, "")
}

func (c *CompactU32) Process() int {
	c.ProcessCompactBytes()
	if c.CompactLength <= 4 {
		data := make([]byte, len(c.Data.Data))
		copy(data, c.Data.Data)

		compactBytes := c.CompactBytes
		bs := make([]byte, 4-len(compactBytes))
		compactBytes = append(compactBytes[:], bs...)
		c.Data.Data = data
		return int(binary.LittleEndian.Uint32(compactBytes)) / 4
	} else {
		return int(binary.LittleEndian.Uint32(c.CompactBytes))
	}
}

func (c *CompactU32) encode(value int) ScaleBytes {
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
	ScaleType
}

func (o *Option) Process() interface{} {
	optionType := o.getNextBytes(1)
	if o.SubType != "" && utiles.BytesToHex(optionType) != "00" {
		return o.ProcessAndUpdateData(o.SubType).Interface()
	}
	return nil
}

type Bytes struct {
	ScaleType
}

func (b *Bytes) Init(data ScaleBytes, args []string) {
	b.TypeString = "Vec<u8>"
	b.ScaleDecoder.Init(data, "")
}

func (b *Bytes) Process() string {
	length := b.ProcessAndUpdateData("Compact<u32>").Int()
	value := b.getNextBytes(int(length))
	if utf8.Valid(value) {
		return string(value)
	} else {
		return utiles.BytesToHex(value)
	}
}

type OptionBytes struct {
	ScaleType
}

func (b *OptionBytes) Init() {
	b.TypeString = "Option<Vec<u8>>"
}

func (b *OptionBytes) Process() string {
	optionByte := b.getNextBytes(1)
	if utiles.BytesToHex(optionByte) != "00" {
		return b.ProcessAndUpdateData("Bytes").String()
	}
	return ""

}

type String struct {
	ScaleType
}

func (s *String) Process() string {
	length := s.ProcessAndUpdateData("Compact<u32>").Int()
	value := s.getNextBytes(int(length))
	return string(value)
}

type HexBytes struct {
	ScaleType
}

func (h *HexBytes) Process() string {
	length := h.ProcessAndUpdateData("Compact<u32>").Int()
	return utiles.AddHex(utiles.BytesToHex(h.getNextBytes(int(length))))
}

type U8 struct {
	ScaleType
}

func (u *U8) Process() int {
	return u.getNextU8()
}

type U32 struct {
	ScaleType
}

func (u *U32) Process() uint32 {
	return uint32(binary.LittleEndian.Uint32(u.getNextBytes(4)))
}

type U64 struct {
	ScaleType
}

func (u *U64) Process() uint64 {
	return uint64(binary.LittleEndian.Uint64(u.getNextBytes(8)))
}

type U128 struct {
	ScaleType
}

func (u *U128) Process() uint128.Uint128 {
	return uint128.FromBytes(u.getNextBytes(16))
}

type H256 struct {
	ScaleType
}

func (h *H256) Process() string {
	return utiles.AddHex(utiles.BytesToHex(h.getNextBytes(32)))
}

type Era struct {
	ScaleType
}

func (e *Era) Process() string {
	optionHex := utiles.BytesToHex(e.getNextBytes(1))
	if optionHex == "00" {
		return optionHex
	} else {
		return optionHex + utiles.BytesToHex(e.getNextBytes(1))
	}
}

type Bool struct {
	ScaleType
}

func (b *Bool) Process() bool {
	return b.getNextBool()
}

type Moment struct {
	CompactU32
}

func (m *Moment) Init(data ScaleBytes, args []string) {
	m.TypeString = "Compact<Moment>"
	m.CompactU32.Init(data, args)
}

func (m *Moment) Process() time.Time {
	intValue := m.ProcessAndUpdateData("Compact<u32>").Int()
	return time.Unix(intValue, 0)
}

func (m *Moment) serialize() {
	m.Process().Format("2006-01-02 15:04:05")
}

type BoxProposal struct {
	ScaleType
	CallIndex  string                   `json:"call_index"`
	Call       interface{}              `json:"call"`
	CallModule interface{}              `json:"call_module"`
	Params     []map[string]interface{} `json:"params"`
}

func (b *BoxProposal) Init() {
	b.TypeString = "Box<Proposal>"
}

func (b *BoxProposal) Process() map[string]interface{} {
	b.CallIndex = utiles.BytesToHex(b.getNextBytes(2))
	callIndex := b.Metadata.CallIndex
	b.CallModule = callIndex["module"]
	b.Call = callIndex["call"]
	//todo
	return map[string]interface{}{
		"call_index":  b.CallIndex,
		"call_name":   b.Call,
		"call_module": b.CallModule,
		"params":      b.Params,
	}
}

type Struct struct {
	ScaleType
}

func (s *Struct) Init() {
	s.TypeMapping = map[string]interface{}{}
}

func (s *Struct) Process() map[string]interface{} {
	result := map[string]interface{}{}
	for key, dataType := range s.TypeMapping {
		result[key] = s.ProcessAndUpdateData(dataType.(string)).Interface()
	}
	return result
}

type ValidatorPrefs struct {
	Struct
}

func (v *ValidatorPrefs) Init() {
	v.TypeMapping = map[string]interface{}{
		"col1": "Compact<u32>",
		"col2": "Compact<Balance>",
	}
	v.TypeString = "(Compact<u32>,Compact<Balance>)"
}

type AccountId struct {
	H256
}

type AccountIndex struct {
	U32
}

type ReferendumIndex struct {
	U32
}

type PropIndex struct {
	U32
}

type Vote struct {
	U8
}

type SessionKey struct {
	H256
}

type AttestedCandidate struct {
	H256
}

type Balance struct {
	U128
}

type ParaId struct {
	U32
}

type Key struct {
	Bytes
}

type KeyValue struct {
	Struct
}

func (k *KeyValue) Init() {
	k.TypeString = "(Vec<u8>, Vec<u8>)"
	k.TypeMapping = map[string]interface{}{
		"key":   "Vec<u8>",
		"value": "Vec<u8>",
	}
}

type Signature struct {
	ScaleType
}

func (s *Signature) Process() string {
	return utiles.BytesToHex(s.getNextBytes(64))
}

type BalanceOf struct {
	Balance
}

type BlockNumber struct {
	U64
}

type NewAccountOutcome struct {
	CompactU32
}

type Vec struct {
	ScaleType
	Elements []interface{} `json:"elements"`
}

func (v *Vec) Process() interface{} {
	elementCount := int(v.ProcessAndUpdateData("Compact<u32>").Int())
	var result []interface{}
	for i := 0; i < elementCount; i++ {
		element := v.ProcessAndUpdateData(v.SubType)
		v.Elements = append(v.Elements, element)
		result = append(result, element)
	}
	return result
}

type Address struct {
	ScaleType
	AccountLength string `json:"account_length"`
	AccountId     string `json:"account_id"`
	AccountIndex  string `json:"account_index"`
	AccountIdx    string `json:"account_idx"`
}

func (a *Address) Process() map[string]string {
	AccountLength := a.getNextBytes(1)
	a.AccountLength = utiles.BytesToHex(AccountLength)
	if a.AccountLength == "ff" {
		a.AccountId = utiles.BytesToHex(a.getNextBytes(32))
	} else {
		var AccountIndex []byte
		if a.AccountLength == "fc" {
			AccountIndex = a.getNextBytes(2)
		} else if a.AccountLength == "fd" {
			AccountIndex = a.getNextBytes(4)
		} else if a.AccountLength == "fe" {
			AccountIndex = a.getNextBytes(8)
		} else {
			AccountIndex = AccountLength
		}
		a.AccountIndex = utiles.BytesToHex(AccountIndex)
		a.AccountIdx = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(AccountIndex)), 10)
	}
	return map[string]string{"account_length": a.AccountLength, "account_id": a.AccountId, "account_index": a.AccountIndex, "account_idx": a.AccountIdx}
}

type RawAddress struct {
	Address
}

type Enum struct {
	ScaleType
	ValueList []string `json:"value_list"`
	Index     int      `json:"index"`
}

func (e *Enum) Init(data ScaleBytes, valueList []string) {
	e.Index = 0
	e.ValueList = valueList
	e.ScaleDecoder.Init(data, "")
}

func (e *Enum) Process() string {
	index := utiles.BytesToHex(e.getNextBytes(1))
	e.Index = utiles.StringToInt(index)
	if e.ValueList[e.Index] != "" {
		return e.ValueList[e.Index]
	}
	return ""
}

type RewardDestination struct {
	Enum
}

func (r *RewardDestination) Init() {
	r.ValueList = []string{"Staked", "Stash", "Controller"}
}

type VoteThreshold struct {
	Enum
}

func (v *VoteThreshold) Init() {
	v.ValueList = []string{"SuperMajorityApprove", "SuperMajorityAgainst", "SimpleMajority"}
}

type Inherent struct {
	Bytes
}

type LockPeriods struct {
	U8
}

type Hash struct {
	H256
}

type VoteIndex struct {
	U32
}

type ProposalIndex struct {
	U32
}

type Permill struct {
	U32
}

type IdentityType struct {
	Bytes
}

type VoteType struct {
	Enum
}

func (v *VoteType) Init() {
	v.TypeString = "voting::VoteType"
	v.ValueList = []string{"Binary", "MultiOption"}
}

type VoteOutcome struct {
	ScaleType
}

func (v *VoteOutcome) Process() []byte {
	return v.getNextBytes(32)
}

type Identity struct {
	Bytes
}

type ProposalTitle struct {
	Bytes
}

type ProposalContents struct {
	Bytes
}

type ProposalStage struct {
	Enum
}

func (p *ProposalStage) Init() {
	p.ValueList = []string{"PreVoting", "Voting", "Completed"}
}

type ProposalCategory struct {
	Enum
}

func (p *ProposalCategory) Init() {
	p.ValueList = []string{"Signaling"}
}

type VoteStage struct {
	Enum
}

func (p *VoteStage) Init() {
	p.ValueList = []string{"PreVoting", "Commit", "Voting", "Completed"}
}

type TallyType struct {
	Enum
}

func (t *TallyType) Init() {
	t.TypeString = "voting::TallyType"
	t.ValueList = []string{"OnePerson", "OneCoin"}
}

type Attestation struct {
	Bytes
}

type ContentId struct {
	H256
}

type MemberId struct {
	U64
}

type PaidTermId struct {
	U64
}

type SubscriptionId struct {
	U64
}

type SchemaId struct {
	U64
}

type DownloadSessionId struct {
	U64
}

type UserInfo struct {
	Struct
}

func (u *UserInfo) Init() {
	u.TypeMapping = map[string]interface{}{
		"handle":     "Option<Vec<u8>>",
		"avatar_uri": "Option<Vec<u8>>",
		"about":      "Option<Vec<u8>>",
	}
}

type Role struct {
	Enum
}

func (r *Role) Init() {
	r.ValueList = []string{"Storage"}
}

type ContentVisibility struct {
	Enum
}

func (c *ContentVisibility) Init() {
	c.ValueList = []string{"Draft", "Public"}
}

type ContentMetadata struct {
	Struct
}

func (c *ContentMetadata) Init() {
	c.TypeMapping = map[string]interface{}{
		"owner":        "AccountId",
		"added_at":     "BlockAndTime",
		"children_ids": "Vec<ContentId>",
		"visibility":   "ContentVisibility",
		"schema":       "SchemaId",
		"json":         "Vec<u8>",
	}
}

type ContentMetadataUpdate struct {
	Struct
}

func (c *ContentMetadataUpdate) Init() {
	c.TypeMapping = map[string]interface{}{
		"children_ids": "Option<Vec<ContentId>>",
		"visibility":   "Option<ContentVisibility>",
		"schema":       "Option<SchemaId>",
		"json":         "Option<Vec<u8>>",
	}
}

type LiaisonJudgement struct {
	Enum
}

func (l *LiaisonJudgement) Init() {
	l.ValueList = []string{"Pending", "Accepted", "Rejected"}
}

type BlockAndTime struct {
	Struct
}

func (b *BlockAndTime) Init() {
	b.TypeMapping = map[string]interface{}{
		"block": "BlockNumber",
		"time":  "Moment",
	}
}

type DataObjectTypeId struct {
	U64
}

func (d *DataObjectTypeId) Init() {
	d.TypeString = "<T as DOTRTrait>::DataObjectTypeId"
}

type DataObject struct {
	Struct
}

func (d *DataObject) Init() {
	d.TypeMapping = map[string]interface{}{
		"owner":             "AccountId",
		"added_at":          "BlockAndTime",
		"type_id":           "DataObjectTypeId",
		"size":              "u64",
		"liaison":           "AccountId",
		"liaison_judgement": "LiaisonJudgement",
	}
}

type DataObjectStorageRelationshipId struct {
	U64
}

type ProposalStatus struct {
	Enum
}

func (p *ProposalStatus) Init() {
	p.ValueList = []string{"Active", "Cancelled", "Expired", "Approved", "Rejected", "Slashed"}
}

type VoteKind struct {
	Enum
}

func (v *VoteKind) Init() {
	v.ValueList = []string{"Abstain", "Approve", "Reject", "Slash"}
}

type TallyResult struct {
	Struct
}

func (t *TallyResult) Init() {
	t.TypeString = "TallyResult<BlockNumber>"
	t.TypeMapping = map[string]interface{}{
		"proposal_id":  "u32",
		"abstentions":  "u32",
		"approvals":    "u32",
		"rejections":   "u32",
		"slashes":      "u32",
		"status":       "ProposalStatus",
		"finalized_at": "BlockNumber",
	}
}

type StorageHasher struct {
	Enum
}

func (s *StorageHasher) Init(data ScaleBytes, valueList []string) {
	valueList = []string{"Blake2_128", "Blake2_256", "Twox128", "Twox256", "Twox128Concat"}
	s.Enum.Init(data, valueList)
}

func (s *StorageHasher) is_blake2_128() bool {
	return s.Index == 0
}
func (s *StorageHasher) is_blake2_256() bool {
	return s.Index == 1
}
func (s *StorageHasher) is_twox128() bool {
	return s.Index == 2
}
func (s *StorageHasher) is_twox256() bool {
	return s.Index == 3
}
func (s *StorageHasher) is_twox128_concat() bool {
	return s.Index == 4
}

type Other struct {
	Bytes
}

type TokenBalance struct {
	U128
}

type Currency struct {
	U128
}

type CurrencyOf struct {
	U128
}
