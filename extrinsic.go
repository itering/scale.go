package scalecodec

import (
	"errors"
	"fmt"

	scaleType "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/types/scaleBytes"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
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

// https://github.com/polkadot-js/api/blob/master/packages/types/src/extrinsic/signedExtensions/index.ts#L24
var signedExts = map[string]bool{
	"CheckSpecVersion":         false,
	"CheckTxVersion":           false,
	"CheckGenesis":             false,
	"CheckMortality":           false,
	"CheckNonce":               false,
	"CheckWeight":              false,
	"ChargeTransactionPayment": false,
	"CheckBlockGasLimit":       false,
	"ChargeAssetTxPayment":     true,
}

func (e *ExtrinsicDecoder) Init(data scaleBytes.ScaleBytes, option *scaleType.ScaleDecoderOption) {
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

type GenericExtrinsic struct {
	VersionInfo        string
	ExtrinsicLength    int
	AddressType        string
	Tip                decimal.Decimal
	SignedExtensions   map[string]interface{}
	AccountId          interface{}
	Signer             interface{}
	Signature          string
	SignatureRaw       interface{}
	Nonce              int
	Era                string
	ExtrinsicHash      string
	CallModuleFunction string
	CallCode           string
	CallModule         string
	Params             []ExtrinsicParam
}

func (e *ExtrinsicDecoder) Process() {
	e.ExtrinsicLength = e.ProcessAndUpdateData("Compact<u32>").(int)
	if e.ExtrinsicLength != e.Data.GetRemainingLength() {
		e.ExtrinsicLength = 0
		e.Data.Reset()
	}

	e.VersionInfo = utiles.BytesToHex(e.NextBytes(1))

	e.ContainsTransaction = utiles.U256(e.VersionInfo).Int64() >= 80

	result := GenericExtrinsic{
		ExtrinsicLength: e.ExtrinsicLength,
		VersionInfo:     e.VersionInfo,
	}

	if e.VersionInfo == "04" || e.VersionInfo == "84" {
		if e.ContainsTransaction {
			// Address
			result.Signer = e.ProcessAndUpdateData(utiles.TrueOrElse(e.Metadata.MetadataVersion >= 14 && scaleType.HasReg("ExtrinsicSigner"), "ExtrinsicSigner", "Address"))
			switch v := result.Signer.(type) {
			case string:
				e.Address = v
				result.AddressType = "AccountId"
			case map[string]interface{}:
				for name, value := range v {
					result.AddressType = name
					e.Address = value
				}
			}
			// ExtrinsicSignature
			result.SignatureRaw = e.ProcessAndUpdateData("ExtrinsicSignature")
			switch v := result.SignatureRaw.(type) {
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
					result.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
				}
			} else {
				result.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
			}
			// spec SignedExtensions
			if len(e.SignedExtensions) > 0 {
				for _, extension := range e.SignedExtensions {
					if utiles.SliceIndex(extension.Name, e.Metadata.Extrinsic.SignedIdentifier) != -1 {
						for _, v := range extension.AdditionalSigned {
							result.SignedExtensions[v.Name] = e.ProcessAndUpdateData(v.Type)
						}
					}
				}
			} else {
				if e.Metadata.MetadataVersion >= 14 {
					for _, ext := range e.Metadata.Extrinsic.SignedExtensions {
						if enable, ok := signedExts[ext.Identifier]; ok && enable {
							result.SignedExtensions[ext.Identifier] = e.ProcessAndUpdateData(ext.TypeString)
						}
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
		e.Params = append(e.Params, ExtrinsicParam{Name: arg.Name, Type: arg.Type, Value: e.ProcessAndUpdateData(arg.Type), TypeName: arg.TypeName})
	}

	if e.ContainsTransaction {
		result.AccountId = e.Address
		result.Signature = e.Signature
		result.Nonce = e.Nonce
		result.Era = e.Era
		result.ExtrinsicHash = e.ExtrinsicHash
	}
	result.CallCode = e.CallIndex
	result.CallModuleFunction = call.Call.Name
	result.CallModule = call.Module.Name
	result.Params = e.Params
	e.Value = &result
}

// Encode extrinsic with option
func (g *GenericExtrinsic) Encode(opt *scaleType.ScaleDecoderOption) (string, error) {
	if opt.Metadata == nil {
		return "", errors.New("invalid metadata")
	}
	data := g.VersionInfo
	if g.VersionInfo == "84" {
		data = data + scaleType.Encode(utiles.TrueOrElse(opt.Metadata.MetadataVersion >= 14 && scaleType.HasReg("ExtrinsicSigner"), "ExtrinsicSigner", "AccountId"), g.Signer) // accountId
		data = data + scaleType.Encode("ExtrinsicSignature", g.SignatureRaw)                                                                                                   // signature
		data = data + scaleType.Encode("EraExtrinsic", g.Era)                                                                                                                  // era
		data = data + scaleType.Encode("Compact<U64>", g.Nonce)                                                                                                                // nonce
		data = data + scaleType.Encode("Compact<Balance>", g.Tip)                                                                                                              // tip
	}

	data = data + g.CallCode
	call, ok := opt.Metadata.CallIndex[g.CallCode]

	if !ok {
		return "", fmt.Errorf("not find Extrinsic Lookup %s, please check metadata info", g.CallCode)
	}

	if len(g.Params) != len(call.Call.Args) {
		return "", fmt.Errorf("extrinsic params length not match, expect %d, got %d", len(call.Call.Args), len(g.Params))
	}

	for index, arg := range call.Call.Args {
		data = data + utiles.TrimHex(scaleType.EncodeWithOpt(arg.Type, g.Params[index].Value, opt))
	}

	return scaleType.Encode("Compact<u32>", len(utiles.HexToBytes(data))) + data, nil
}
