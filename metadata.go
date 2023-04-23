package scalecodec

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
)

type MetadataDecoder struct {
	types.ScaleDecoder
	Version  string               `json:"version"`
	Metadata types.MetadataStruct `json:"metadata"`
}

func (m *MetadataDecoder) Init(data []byte) {
	sData := scaleBytes.ScaleBytes{Data: data}
	m.ScaleDecoder.Init(sData, nil)
}

func (m *MetadataDecoder) Process() error {
	magicBytes := m.NextBytes(4)
	if string(magicBytes) == "meta" {
		metadataVersion := utiles.U256(utiles.BytesToHex(m.Data.Data[m.Data.Offset : m.Data.Offset+1]))
		m.Version = m.ProcessAndUpdateData("MetadataVersion").(string)
		m.Metadata = m.ProcessAndUpdateData(m.Version).(types.MetadataStruct)
		m.Metadata.MetadataVersion = int(metadataVersion.Int64())
		return nil
	}
	return errors.New("not metadata")

}

// CheckRegistry
// Deprecated
func (m *MetadataDecoder) CheckRegistry() (notReg []string) { return }

// Compare compare metadata
// dest is older metadata

type MetadataCompareResult struct {
	NewModules    []string                 `json:"new_modules,omitempty"`
	ModuleChanges map[string]ModuleChanges `json:"module_changes"`
}

type ModuleChanges struct {
	Calls  *ChangesList `json:"calls,omitempty"`
	Events *ChangesList `json:"events,omitempty"`
}

type ChangesList struct {
	New     []string         `json:"new,omitempty"`
	Changes []CompareChanges `json:"changes,omitempty"`
}

type CompareChanges struct {
	Prev    string `json:"prev"`
	Current string `json:"next"`
}

func (m *MetadataDecoder) Compare(dest *types.MetadataTag) *MetadataCompareResult {
	var getModuleByName = func(name string) *types.MetadataModules {
		for _, module := range dest.Modules {
			if module.Name == name {
				return &module
			}
		}
		return nil
	}

	var getCallByName = func(name string, calls []types.MetadataCalls) *types.MetadataCalls {
		for _, call := range calls {
			if call.Name == name {
				return &call
			}
		}
		return nil
	}

	var getEventByName = func(name string, events []types.MetadataEvents) *types.MetadataEvents {
		for _, event := range events {
			if event.Name == name {
				return &event
			}
		}
		return nil
	}

	var buildCallArgs = func(call types.MetadataCalls) string {
		var args []string
		for _, arg := range call.Args {
			args = append(args, arg.Type)
		}
		return fmt.Sprintf("%s(%s)", call.Name, strings.Join(args, ","))
	}

	var buildEventArgs = func(event types.MetadataEvents) string {
		return fmt.Sprintf("%s(%s)", event.Name, strings.Join(event.Args, ","))
	}

	var result MetadataCompareResult
	result.ModuleChanges = make(map[string]ModuleChanges)
	for _, module := range m.Metadata.Metadata.Modules {
		destModule := getModuleByName(module.Name)
		if destModule == nil {
			// new modules
			result.NewModules = append(result.NewModules, module.Name)
			continue
		}

		// check call
		moduleChanges := ModuleChanges{}
		calls := ChangesList{}
		for _, call := range module.Calls {
			destCall := getCallByName(call.Name, destModule.Calls)
			if destCall == nil {
				calls.New = append(calls.New, buildCallArgs(call))
				continue
			}
			// check call args
			if len(call.Args) != len(destCall.Args) {
				calls.Changes = append(calls.Changes, CompareChanges{
					Prev:    buildCallArgs(*destCall),
					Current: buildCallArgs(call),
				})
				continue
			}
			for index, arg := range call.Args {
				// check every type
				destType := destCall.Args[index].Type
				if !types.Eq(arg.Type, destType) {
					calls.Changes = append(calls.Changes, CompareChanges{
						Prev:    buildCallArgs(*destCall),
						Current: buildCallArgs(call),
					})
					continue
				}
			}
		}
		if calls.New != nil || calls.Changes != nil {
			moduleChanges.Calls = &calls
		}

		// check event
		Events := ChangesList{}
		for _, event := range module.Events {
			destEvent := getEventByName(event.Name, destModule.Events)
			if destEvent == nil {
				Events.New = append(Events.New, buildEventArgs(event))
				continue
			}
			// check call args
			if len(event.Args) != len(destEvent.Args) {
				Events.Changes = append(Events.Changes, CompareChanges{
					Prev:    buildEventArgs(*destEvent),
					Current: buildEventArgs(event),
				})
				continue
			}
			for index, arg := range event.Args {
				// check every type
				if !types.Eq(arg, destEvent.Args[index]) {
					Events.Changes = append(Events.Changes, CompareChanges{
						Prev:    buildEventArgs(*destEvent),
						Current: buildEventArgs(event),
					})
					continue
				}
			}
		}
		if Events.New != nil || Events.Changes != nil {
			moduleChanges.Events = &Events
		}
		if (moduleChanges.Calls != nil && (len(moduleChanges.Calls.Changes) > 0 || len(moduleChanges.Calls.New) > 0)) || (moduleChanges.Events != nil && (len(moduleChanges.Events.Changes) > 0 || len(moduleChanges.Events.New) > 0)) {
			result.ModuleChanges[module.Name] = moduleChanges
		}

	}
	return &result
}
