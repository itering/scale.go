package types

import "encoding/json"

type PortableType struct {
	Id   int     `json:"id"`
	Type *SiType `json:"type"`
}

type SiType struct {
	Path   []string          `json:"path"`
	Params []SiTypeParameter `json:"params"`
	Def    SiTypeDef         `json:"def"`
	Docs   []string          `json:"docs"`
}

type SiTypeParameter struct {
	Name string `json:"name"`
	Type int    `json:"type,omitempty"`
}

type SiTypeDef struct {
	Composite   *SiTypeDefComposite   `json:"Composite,omitempty"`
	Variant     *SiTypeDefVariant     `json:"Variant,omitempty"`
	Sequence    *SiTypeDefSequence    `json:"SiTypeDefSequence,omitempty"`
	Array       *SiTypeDefArray       `json:"Array,omitempty"`
	Tuple       *SiTypeDefTuple       `json:"Tuple,omitempty"`
	Primitive   *SiTypeDefPrimitive   `json:"Primitive,omitempty"`
	Compact     *SiTypeDefCompact     `json:"Compact,omitempty"`
	BitSequence *SiTypeDefBitSequence `json:"BitSequence,omitempty"`
	// HistoricMetaCompat compatibility
	HistoricMetaCompat string `json:"HistoricMetaCompat,omitempty"`
}

type SiTypeDefComposite struct {
	Fields []SiField `json:"fields"`
}

type SiField struct {
	Name     string   `json:"name,omitempty"`
	Type     int      `json:"type"`
	TypeName string   `json:"type_name,omitempty"`
	Docs     []string `json:"docs"`
}

type SiTypeDefVariant struct {
	Variants []SiVariant `json:"variants"`
}

type SiVariant struct {
	Name   string    `json:"name"`
	Fields []SiField `json:"fields"`
	Index  int       `json:"index"`
	Docs   []string  `json:"docs"`
}

type SiTypeDefSequence struct {
	Type int `json:"type"`
}

type SiTypeDefCompact struct {
	Type int `json:"type"`
}

type SiTypeDefArray struct {
	Len  int `json:"len"`
	Type int `json:"type"`
}

type SiTypeDefTuple []int

type SiTypeDefPrimitive string

type SiTypeDefBitSequence struct {
	BitStoreType int `json:"bitStoreType"`
	BitOrderType int `json:"bitOrderType"`
}

func InitPortableRaw(raw []interface{}) []PortableType {
	var t []PortableType
	bm, _ := json.Marshal(raw)
	if err := json.Unmarshal(bm, &t); err != nil {
		panic(err)
	}
	return t
}
