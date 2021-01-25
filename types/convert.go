package types

import "strings"

func ConvertType(name string, option ...interface{}) string {
	defer func() {
		if len(option) == 0 {
			RuntimeCodecType = append(RuntimeCodecType, name)
		}
	}()
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "T::", "")
	name = strings.ReplaceAll(name, "VecDeque<", "Vec<")
	name = strings.ReplaceAll(name, "<T>", "")
	name = strings.ReplaceAll(name, "<T as Trait>::", "")
	name = strings.ReplaceAll(name, "<T, I>", "")
	name = strings.ReplaceAll(name, "\n", " ")
	switch name {
	case "()", "<InherentOfflineReport as InherentOfflineReport>::Inherent":
		name = "Null"
	case "Vec<u8>":
		name = "Bytes"
	case "<Lookup as StaticLookup>::Source":
		name = "Address"
	case "<Balance as HasCompact>::Type":
		name = "Compact<Balance>"
	case "<BlockNumber as HasCompact>::Type":
		name = "Compact<BlockNumber>"
	case "<Moment as HasCompact>::Type":
		name = "Compact<Moment>"
	case "<T as Trait<I>>::Proposal":
		name = "Proposal"
	case "wasm::PrefabWasmModule":
		name = "PrefabWasmModule"
	}
	return name
}
