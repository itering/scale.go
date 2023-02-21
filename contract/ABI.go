package contract

import (
	"encoding/json"

	"github.com/itering/scale.go/types"
)

type MetadataAbi struct {
	RegisteredSiType map[int]string `json:"-"`
	Spec             Spec           `json:"spec"`
	Types            []AbiType      `json:"types"`
	Version          string         `json:"version,omitempty"`
}

type Spec struct {
	Constructors []Message           `json:"constructors"`
	Docs         []interface{}       `json:"docs"`
	LangError    TypeWithDisplayName `json:"lang_error"`
	Events       []SpecEvent         `json:"events"`
	Messages     []interface{}       `json:"messages"`
}

type SpecEvent struct {
	Args  []EventArg `json:"args"`
	Docs  []string   `json:"docs"`
	Label string     `json:"label"`
}
type EventArg struct {
	Docs    []string            `json:"docs"`
	Indexed bool                `json:"indexed"`
	Payable bool                `json:"payable"`
	Type    TypeWithDisplayName `json:"type"`
}

type TypeWithDisplayName struct {
	DisplayName []string `json:"displayName"`
	Type        int      `json:"type"`
}

type Message struct {
	Args       []MessageArg        `json:"args"`
	Docs       []string            `json:"docs"`
	Label      string              `json:"label"`
	Payable    bool                `json:"payable"`
	ReturnType TypeWithDisplayName `json:"returnType"`
	Selector   string              `json:"selector"`
}

type MessageArg struct {
	Label string              `json:"label"`
	Type  TypeWithDisplayName `json:"type"`
}

type AbiType struct {
	Id   int       `json:"id"`
	Type AbiSiType `json:"type"`
}

type AbiSiType struct {
	Path   []string                `json:"path,omitempty"`
	Params []types.SiTypeParameter `json:"params,omitempty"`
	Def    SiTypeDef               `json:"def"`
	Docs   []string                `json:"docs,omitempty"`
}

type SiTypeDef struct {
	Composite      *types.SiTypeDefComposite   `json:"composite,omitempty"`
	Variant        *types.SiTypeDefVariant     `json:"variant,omitempty"`
	Sequence       *types.SiTypeDefSequence    `json:"sequence,omitempty"`
	Array          *types.SiTypeDefArray       `json:"array,omitempty"`
	Tuple          *types.SiTypeDefTuple       `json:"tuple,omitempty"`
	Primitive      *types.SiTypeDefPrimitive   `json:"primitive,omitempty"`
	Compact        *types.SiTypeDefCompact     `json:"compact,omitempty"`
	BitSequence    *types.SiTypeDefBitSequence `json:"bitSequence,omitempty"`
	SiTypeDefRange *types.SiTypeDefRange       `json:"range,omitempty"`
}

func InitAbi(raw []byte) (*MetadataAbi, error) {
	var abi MetadataAbi
	err := json.Unmarshal(raw, &abi)
	if err != nil {
		return nil, err
	}
	return &abi, nil
}

func (a *AbiType) ToScalePortableType() *types.PortableType {
	pt := types.PortableType{Id: a.Id}

	siDef := types.SiTypeDef{}
	def := a.Type.Def
	switch {
	case def.Array != nil:
		siDef.Array = def.Array
	case def.Composite != nil:
		siDef.Composite = def.Composite
	case def.Variant != nil:
		siDef.Variant = def.Variant
	case def.Sequence != nil:
		siDef.Sequence = def.Sequence
	case def.Tuple != nil:
		siDef.Tuple = def.Tuple
	case def.Primitive != nil:
		siDef.Primitive = def.Primitive
	case def.Compact != nil:
		siDef.Compact = def.Compact
	case def.BitSequence != nil:
		siDef.BitSequence = def.BitSequence
	case def.SiTypeDefRange != nil:
		siDef.SiTypeDefRange = def.SiTypeDefRange
	}

	pt.Type = types.SiType{
		Path:   a.Type.Path,
		Params: a.Type.Params,
		Docs:   a.Type.Docs,
		Def:    siDef,
	}
	return &pt
}

func (m *MetadataAbi) Register(s *types.ScaleDecoder, prefix string) {
	id2Portable := make(map[int]types.SiType)
	for _, v := range m.Types {
		id2Portable[v.Id] = v.ToScalePortableType().Type
	}
	scaleInfo := types.ScaleInfo{
		ScaleDecoder: s,
		PathPrefix:   prefix,
	}
	scaleInfo.ProcessSiType(id2Portable)
	m.RegisteredSiType = s.RegisteredSiType
}

func (m *MetadataAbi) GetTypeNames(id int) string {
	return m.RegisteredSiType[id]
}
