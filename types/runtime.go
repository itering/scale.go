package types

import (
	"errors"
	"strings"
)

type RuntimeConfig struct {
	Runtime map[string]interface{}
}

func (r *RuntimeConfig) init() {
	r.Runtime = map[string]interface{}{
		"vec<u8>":                      &Bytes{},
		"enum":                         &Enum{},
		"bytes":                        &Bytes{},
		"vec":                          &Vec{},
		"compact<u32>":                 &CompactU32{},
		"bool":                         &Bool{},
		"storagehasher":                &StorageHasher{},
		"hexbytes":                     &HexBytes{},
		"moment":                       &Moment{},
		"compact<moment>":              &Moment{},
		"u32":                          &U32{},
		"blocknumber":                  &BlockNumber{},
		"accountid":                    &AccountId{},
		"sessionindex":                 &SessionIndex{},
		"eraindex":                     &EraIndex{},
		"stakingledger":                &StakingLedgers{},
		"extendedbalance":              &ExtendedBalance{},
		"ringbalanceof":                &RingBalanceOf{},
		"ktonbalanceof":                &KtonBalanceOf{},
		"unlockchunk":                  &UnlockChunk{},
		"compact":                      &Compact{},
		"regularitem":                  &RegularItem{},
		"stakingbalance":               &StakingBalance{},
		"keys":                         &Keys{},
		"metadatamoduleevent":          &MetadataModuleEvent{},
		"metadatamodulecallargument":   &MetadataModuleCallArgument{},
		"metadatamodulecall":           &MetadataModuleCall{},
		"metadatav6decoder":            &MetadataV6Decoder{},
		"metadatav6module":             &MetadataV6Module{},
		"metadatav6modulestorage":      &MetadataV6ModuleStorage{},
		"metadatav6moduleconstants":    &MetadataV6ModuleConstants{},
		"metadatav7decoder":            &MetadataV7Decoder{},
		"metadatav7module":             &MetadataV7Module{},
		"metadatav7modulestorage":      &MetadataV7ModuleStorage{},
		"metadatav7moduleconstants":    &MetadataV7ModuleConstants{},
		"metadatav7modulestorageentry": &MetadataV7ModuleStorageEntry{},
	}
}

func (r *RuntimeConfig) GetDecoderClass(typeStr string) (interface{}, error) {
	typeStr = strings.ToLower(typeStr)
	if r.Runtime[typeStr] == nil {
		return nil, errors.New("Scalecodec type nil" + typeStr)
	}
	return r.Runtime[typeStr], nil
}

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
