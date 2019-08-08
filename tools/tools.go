package tools

import "strings"

func ConvertType(name string) string {
	name = strings.ReplaceAll(name, "T::", "")
	name = strings.ReplaceAll(name, "<T>", "")
	name = strings.ReplaceAll(name, "<T as Trait>::", "")
	if name == "()" {
		return "Null"
	}
	if name == "Vec<u8>" {
		return "Bytes"
	}
	if name == "<Lookup as StaticLookup>::Source" {
		return "Address"
	}
	if name == "<Balance as HasCompact>::Type" {
		return "Compact<Balance>"
	}
	if name == "<BlockNumber as HasCompact>::Type" {
		return "Compact<BlockNumber>"
	}
	if name == "<Balance as HasCompact>::Type" {
		return "Compact<Balance>"
	}
	if name == "<Moment as HasCompact>::Type" {
		return "Compact<Moment>"
	}
	if name == "<InherentOfflineReport as InherentOfflineReport>::Inherent" {
		return "Inherent"
	}
	return name
}
