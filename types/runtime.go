package types

import (
	"errors"
	"strings"
)

type RuntimeConfig struct {
}

var regRuntimeStruct map[string]interface{}

func (r *RuntimeConfig) init() {
	regRuntimeStruct = make(map[string]interface{})
	regRuntimeStruct["vec<u8>"] = &Bytes{}
	regRuntimeStruct["enum"] = &Enum{}
	regRuntimeStruct["bytes"] = &Bytes{}
	regRuntimeStruct["vec"] = &Vec{}
	regRuntimeStruct["compact<u32>"] = &CompactU32{}
	regRuntimeStruct["bool"] = &Bool{}
	regRuntimeStruct["storagehasher"] = &StorageHasher{}
	regRuntimeStruct["hexbytes"] = &HexBytes{}
	regRuntimeStruct["moment"] = &Moment{}
	regRuntimeStruct["compact<moment>"] = &Moment{}
	regRuntimeStruct["u32"] = &U32{}
	regRuntimeStruct["blocknumber"] = &BlockNumber{}
	regRuntimeStruct["accountid"] = &AccountId{}
	regRuntimeStruct["sessionindex"] = &SessionIndex{}

	//metadata
	regRuntimeStruct["metadatamoduleevent"] = &MetadataModuleEvent{}
	regRuntimeStruct["metadatamodulecallargument"] = &MetadataModuleCallArgument{}
	regRuntimeStruct["metadatamodulecall"] = &MetadataModuleCall{}
	regRuntimeStruct["metadatav5decoder"] = &MetadataV5Decoder{}
	regRuntimeStruct["metadatav5module"] = &MetadataV5Module{}
	regRuntimeStruct["metadatav5modulestorage"] = &MetadataV5ModuleStorage{}
	regRuntimeStruct["metadatav6decoder"] = &MetadataV6Decoder{}
	regRuntimeStruct["metadatav6module"] = &MetadataV6Module{}
	regRuntimeStruct["metadatav6modulestorage"] = &MetadataV6ModuleStorage{}
	regRuntimeStruct["metadatav6moduleconstants"] = &MetadataV6ModuleConstants{}
}

func (r *RuntimeConfig) GetDecoderClass(typeStr string) (interface{}, error) {
	typeStr = strings.ToLower(typeStr)
	if regRuntimeStruct[typeStr] == nil {
		return nil, errors.New("Scalecodec type nil" + typeStr)
	}
	return regRuntimeStruct[typeStr], nil
}
