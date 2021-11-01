package types

type MetadataModuleCall struct {
	ScaleDecoder
	Name string                       `json:"name"`
	Args []MetadataModuleCallArgument `json:"args"`
	Docs []string                     `json:"docs"`
}

func (m *MetadataModuleCall) Process() {
	m.Name = m.ProcessAndUpdateData("String").(string)
	argsValue := m.ProcessAndUpdateData("Vec<MetadataModuleCallArgument>").([]interface{})
	var args []MetadataModuleCallArgument
	for _, v := range argsValue {
		arg := v.(map[string]string)
		args = append(args, MetadataModuleCallArgument{Name: arg["name"], Type: arg["type"]})
	}
	m.Args = args
	docs := m.ProcessAndUpdateData("Vec<String>").([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(string))
	}
	m.Value = MetadataModuleCall{
		Name: m.Name,
		Args: m.Args,
		Docs: m.Docs,
	}
}

type MetadataModuleCallArgument struct {
	ScaleDecoder
	Name     string `json:"name"`
	Type     string `json:"type"`
	TypeName string `json:"type_name,omitempty"`
}

func (m *MetadataModuleCallArgument) Process() {
	m.Name = m.ProcessAndUpdateData("String").(string)
	m.Type = ConvertType(m.ProcessAndUpdateData("String").(string))
	m.Value = map[string]string{
		"name": m.Name,
		"type": m.Type,
	}
}

type MetadataModuleEvent struct {
	ScaleDecoder
	Name string   `json:"name"`
	Docs []string `json:"docs"`
	Args []string `json:"args"`
}

func (m *MetadataModuleEvent) Process() {
	m.Name = m.ProcessAndUpdateData("String").(string)
	args := m.ProcessAndUpdateData("Vec<String>").([]interface{})
	for _, v := range args {
		m.Args = append(m.Args, v.(string))
	}
	docs := m.ProcessAndUpdateData("Vec<String>").([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(string))
	}
	m.Value = MetadataEvents{
		Name: m.Name,
		Args: m.Args,
		Docs: m.Docs,
	}
}
