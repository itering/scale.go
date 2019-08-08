package types

type MetadataCallAndEvent struct {
	CallIndex  map[string]interface{} `json:"call_index"`
	EventIndex map[string]interface{} `json:"event_index"`
}

type MetadataModules struct {
	Name    string            `json:"name"`
	Prefix  string            `json:"prefix"`
	Storage []MetadataStorage `json:"storage"`
	Calls   []MetadataCalls   `json:"calls"`
	Events  []MetadataEvents  `json:"events"`
}

type MetadataStorage struct {
	Name     string                 `json:"name"`
	Modifier string                 `json:"modifier"`
	Type     map[string]interface{} `json:"type"`
	Fallback string                 `json:"fallback"`
	Docs     []string               `json:"docs"`
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
	MagicNumber int                    `json:"magicNumber"`
	Metadata    MetadataTag            `json:"metadata"`
	CallIndex   map[string]interface{} `json:"call_index"`
	EventIndex  map[string]interface{} `json:"event_index"`
}

type MetadataTag struct {
	Modules []MetadataModules `json:"modules"`
}
