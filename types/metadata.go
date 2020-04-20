package types

type MetadataCallAndEvent struct {
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

type MetadataModules struct {
	Name      string                `json:"name"`
	Prefix    string                `json:"prefix"`
	Storage   []MetadataStorage     `json:"storage"`
	Calls     []MetadataCalls       `json:"calls"`
	Events    []MetadataEvents      `json:"events"`
	Constants []MetadataConstants   `json:"constants,omitempty"`
	Errors    []MetadataModuleError `json:"errors"`
}

type MetadataStorage struct {
	Name     string      `json:"name"`
	Modifier string      `json:"modifier"`
	Type     StorageType `json:"type"`
	Fallback string      `json:"fallback"`
	Docs     []string    `json:"docs"`
	Hasher   string      `json:"hasher,omitempty"`
}

type StorageType struct {
	Origin        string   `json:"origin"`
	PlainType     *string  `json:"plain_type,omitempty"`
	MapType       *MapType `json:"map_type,omitempty"`
	DoubleMapType *MapType `json:"double_map_type,omitempty"`
}

type MapType struct {
	Hasher     string `json:"hasher"`
	Key        string `json:"key"`
	Key2       string `json:"key2,omitempty"`
	Key2Hasher string `json:"key2Hasher,omitempty"`
	Value      string `json:"value"`
	IsLinked   bool   `json:"isLinked"`
}

type MetadataCalls struct {
	Lookup string                   `json:"lookup"`
	Name   string                   `json:"name"`
	Docs   []string                 `json:"docs"`
	Args   []map[string]interface{} `json:"args"`
}

type MetadataEvents struct {
	Lookup string   `json:"lookup"`
	Name   string   `json:"name"`
	Docs   []string `json:"docs"`
	Args   []string `json:"args"`
}

type MetadataStruct struct {
	MetadataVersion int         `json:"metadata_version"`
	Metadata        MetadataTag `json:"metadata"`
}

type MetadataTag struct {
	Modules []MetadataModules `json:"modules"`
}

type MetadataConstants struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	ConstantsValue string   `json:"constants_value"`
	Docs           []string `json:"docs"`
}
