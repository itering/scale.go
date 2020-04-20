package types

type MetadataModuleCall struct {
	ScaleDecoder
	Name string              `json:"name"`
	Args []map[string]string `json:"args"`
	Docs []string            `json:"docs"`
}

func (m *MetadataModuleCall) Process() {
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	argsValue := m.ProcessAndUpdateData("Vec<MetadataModuleCallArgument>").([]interface{})
	var args []map[string]string
	for _, v := range argsValue {
		args = append(args, v.(map[string]string))
	}
	m.Args = args
	docs := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
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
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m *MetadataModuleCallArgument) Process() {
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	m.Type = ConvertType(m.ProcessAndUpdateData("Bytes").(string))
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
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	args := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range args {
		m.Args = append(m.Args, v.(string))
	}
	docs := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(string))
	}
	m.Value = MetadataEvents{
		Name: m.Name,
		Args: m.Args,
		Docs: m.Docs,
	}
}
