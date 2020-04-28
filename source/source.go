package source

import (
	"encoding/json"
)

type TypeStruct struct {
	Type        string     `json:"type"`
	TypeString  string     `json:"type_string"`
	TypeMapping [][]string `json:"type_mapping,omitempty"`
	ValueList   []string   `json:"value_list,omitempty"`
}

func LoadTypeRegistry(source []byte) map[string]TypeStruct {
	var original map[string]interface{}
	err := json.Unmarshal(source, &original)
	if err != nil {
		panic(err)
	}
	ts := make(map[string]TypeStruct)
	for key, value := range original {
		switch value.(type) {
		case string:
			ts[key] = TypeStruct{Type: "string", TypeString: value.(string)}
		default:
			b, _ := json.Marshal(value)
			var t TypeStruct
			_ = json.Unmarshal(b, &t)
			ts[key] = t
		}
	}
	return ts
}
