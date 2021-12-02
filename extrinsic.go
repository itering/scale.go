package scalecodec

import (
	"fmt"

	scaleType "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"golang.org/x/crypto/blake2b"
)

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	TypeName string      `json:"type_name"`
	Value    interface{} `json:"value"`
}

type ExtrinsicDecoder struct {
	scaleType.ScaleDecoder
	ExtrinsicLength     int              `json:"extrinsic_length"`
	ExtrinsicHash       string           `json:"extrinsic_hash"`
	VersionInfo         string           `json:"version_info"`
	ContainsTransaction bool             `json:"contains_transaction"`
	Address             interface{}      `json:"address"`
	Signature           string           `json:"signature"`
	Nonce               int              `json:"nonce"`
	Era                 string           `json:"era"`
	CallIndex           string           `json:"call_index"`
	Params              []ExtrinsicParam `json:"params"`
	Metadata            *scaleType.MetadataStruct
	SignedExtensions    []scaleType.SignedExtension `json:"signed_extensions"`
}

func (e *ExtrinsicDecoder) Init(data scaleType.ScaleBytes, option *scaleType.ScaleDecoderOption) {
	if option == nil || option.Metadata == nil {
		panic("ExtrinsicDecoder option metadata required")
	}
	e.Params = []ExtrinsicParam{}
	e.Metadata = option.Metadata
	e.SignedExtensions = option.SignedExtensions
	e.ScaleDecoder.Init(data, option)
}

func (e *ExtrinsicDecoder) generateHash() string {
	if !e.ContainsTransaction {
		return ""
	}
	var extrinsicData []byte
	if e.ExtrinsicLength > 0 {
		extrinsicData = e.Data.Data
	} else {
		extrinsicLengthType := scaleType.CompactU32{}
		extrinsicLengthType.Encode(len(e.Data.Data))
		extrinsicData = append(extrinsicLengthType.Data.Data[:], e.Data.Data[:]...)
	}
	checksum, _ := blake2b.New(32, []byte{})
	_, _ = checksum.Write(extrinsicData)
	h := checksum.Sum(nil)
	return utiles.BytesToHex(h)
}

func (e *ExtrinsicDecoder) Process() {
	e.ExtrinsicLength = e.ProcessAndUpdateData("Compact<u32>").(int)
	if e.ExtrinsicLength != e.Data.GetRemainingLength() {
		e.ExtrinsicLength = 0
		e.Data.Reset()
	}

	e.VersionInfo = utiles.BytesToHex(e.NextBytes(1))

	e.ContainsTransaction = utiles.U256(e.VersionInfo).Int64() >= 80

	result := map[string]interface{}{
		"extrinsic_length": e.ExtrinsicLength,
		"version_info":     e.VersionInfo,
	}
	if e.VersionInfo == "04" || e.VersionInfo == "84" {
		if e.ContainsTransaction {
			// Address
			address := e.ProcessAndUpdateData("Address")
			switch v := address.(type) {
			case string:
				e.Address = v
				result["address_type"] = "AccountId"
			case map[string]interface{}:
				for name, value := range v {
					result["address_type"] = name
					e.Address = value
				}
			}
			// ExtrinsicSignature
			signature := e.ProcessAndUpdateData("ExtrinsicSignature")
			switch v := signature.(type) {
			case string:
				e.Signature = v
			case map[string]interface{}:
				for _, value := range v {
					e.Signature = value.(string)
				}
			}
			e.Era = e.ProcessAndUpdateData("EraExtrinsic").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			if e.Metadata.Extrinsic != nil {
				if utiles.SliceIndex("ChargeTransactionPayment", e.Metadata.Extrinsic.SignedIdentifier) != -1 {
					result["tip"] = e.ProcessAndUpdateData("Compact<Balance>")
				}
			} else {
				result["tip"] = e.ProcessAndUpdateData("Compact<Balance>")
			}

			// SignedExtensions
			// https://github.com/polkadot-js/api/blob/9ae87bed782a5e3e345e20f6a9b64687d399a257/packages/types/src/extrinsic/signedExtensions/index.ts
			for _, extension := range e.SignedExtensions {
				if utiles.SliceIndex(extension.Name, e.Metadata.Extrinsic.SignedIdentifier) != -1 {
					for _, v := range extension.AdditionalSigned {
						result[v.Name] = e.ProcessAndUpdateData(v.Type)
					}
				}
			}
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))
	} else {
		panic(fmt.Sprintf("Extrinsics version %s is not support", e.VersionInfo))
	}
	if e.CallIndex == "" {
		panic("Not find Extrinsic Lookup, please check type registry")
	}

	call, ok := e.Metadata.CallIndex[e.CallIndex]
	if !ok {
		panic(fmt.Sprintf("Not find Extrinsic Lookup %s, please check metadata info", e.CallIndex))
	}
	e.Module = call.Module.Name

	for _, arg := range call.Call.Args {
		value := e.ProcessAndUpdateData(arg.Type)
		param := ExtrinsicParam{
			Name:  arg.Name,
			Type:  arg.Type,
			Value: value,
		}
		if param.TypeName == "" {
			param.TypeName = arg.TypeName
		}
		e.Params = append(e.Params, param)
	}

	if e.ContainsTransaction {
		result["account_id"] = e.Address
		result["signature"] = e.Signature
		result["nonce"] = e.Nonce
		result["era"] = e.Era
		result["extrinsic_hash"] = e.ExtrinsicHash
	}

	if e.CallIndex != "" {
		result["call_code"] = e.CallIndex
		result["call_module_function"] = call.Call.Name
		result["call_module"] = call.Module.Name
	}

	result["nonce"] = e.Nonce
	result["era"] = e.Era
	result["params"] = e.Params
	e.Value = result
}
