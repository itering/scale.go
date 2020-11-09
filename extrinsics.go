package scalecodec

import (
	"fmt"
	scaleType "github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/blake2b"
)

type ExtrinsicParam struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	ValueRaw string      `json:"value_raw"`
}

type ExtrinsicDecoder struct {
	scaleType.ScaleDecoder
	ExtrinsicLength     int                       `json:"extrinsic_length"`
	ExtrinsicHash       string                    `json:"extrinsic_hash"`
	VersionInfo         string                    `json:"version_info"`
	ContainsTransaction bool                      `json:"contains_transaction"`
	Address             string                    `json:"address"`
	Signature           string                    `json:"signature"`
	SignatureVersion    int                       `json:"signature_version"`
	Nonce               int                       `json:"nonce"`
	Era                 string                    `json:"era"`
	CallIndex           string                    `json:"call_index"`
	Tip                 decimal.Decimal           `json:"tip"`
	CallModule          scaleType.MetadataModules `json:"call_module"`
	Call                scaleType.MetadataCalls   `json:"call"`
	Params              []ExtrinsicParam          `json:"params"`
	Metadata            *scaleType.MetadataStruct
}

func (e *ExtrinsicDecoder) Init(data scaleType.ScaleBytes, option *scaleType.ScaleDecoderOption) {
	if option == nil || option.Metadata == nil {
		panic("ExtrinsicDecoder option metadata required")
	}
	e.Params = []ExtrinsicParam{}
	e.Metadata = option.Metadata
	e.ScaleDecoder.Init(data, option)
}

func (e *ExtrinsicDecoder) generateHash() string {
	if e.ContainsTransaction {
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
	return ""
}

func (e *ExtrinsicDecoder) Process() {
	e.ExtrinsicLength = e.ProcessAndUpdateData("Compact<u32>").(int)
	if e.ExtrinsicLength != e.Data.GetRemainingLength() {
		e.ExtrinsicLength = 0
		e.Data.Reset()
	}

	e.VersionInfo = utiles.BytesToHex(e.NextBytes(1))

	e.ContainsTransaction = utiles.U256(e.VersionInfo).Int64() >= 80

	if e.VersionInfo == "01" || e.VersionInfo == "81" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("Address").(string)
			e.Signature = e.ProcessAndUpdateData("Signature").(string)
			e.Nonce = e.ProcessAndUpdateData("Compact<u32>").(int)
			e.Era = e.ProcessAndUpdateData("Era").(string)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))

	} else if e.VersionInfo == "02" || e.VersionInfo == "82" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("Address").(string)
			e.Signature = e.ProcessAndUpdateData("Signature").(string)
			e.Era = e.ProcessAndUpdateData("Era").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			e.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))

	} else if e.VersionInfo == "03" || e.VersionInfo == "83" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("Address").(string)
			e.Signature = e.ProcessAndUpdateData("Signature").(string)
			e.Era = e.ProcessAndUpdateData("Era").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			e.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))

	} else if e.VersionInfo == "04" || e.VersionInfo == "84" {

		if e.ContainsTransaction {
			e.Address = e.ProcessAndUpdateData("GenericAddress").(string)
			e.SignatureVersion = e.ProcessAndUpdateData("U8").(int)
			if e.SignatureVersion == 2 {
				e.Signature = e.ProcessAndUpdateData("EcdsaSignature").(string)
			} else {
				e.Signature = e.ProcessAndUpdateData("Signature").(string)
			}
			e.Era = e.ProcessAndUpdateData("EraExtrinsic").(string)
			e.Nonce = int(e.ProcessAndUpdateData("Compact<U64>").(uint64))
			e.Tip = e.ProcessAndUpdateData("Compact<Balance>").(decimal.Decimal)
			e.ExtrinsicHash = e.generateHash()
		}
		e.CallIndex = utiles.BytesToHex(e.NextBytes(2))
	} else {
		panic(fmt.Sprintf("Extrinsics version %s is not support", e.VersionInfo))
	}

	if e.CallIndex != "" {
		call, ok := e.Metadata.CallIndex[e.CallIndex]
		if !ok {
			panic(fmt.Sprintf("Not find Extrinsic Lookup %s, please check metadata info", e.CallIndex))
		}
		e.Call = call.Call
		e.CallModule = call.Module

		for _, arg := range e.Call.Args {
			e.Params = append(e.Params,
				ExtrinsicParam{
					Name:  arg.Name,
					Type:  arg.Type,
					Value: e.ProcessAndUpdateData(arg.Type),
				})
		}
	}

	result := map[string]interface{}{
		"extrinsic_length": e.ExtrinsicLength,
		"version_info":     e.VersionInfo,
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
		result["call_module_function"] = e.Call.Name
		result["call_module"] = e.CallModule.Name
	}

	result["nonce"] = e.Nonce
	result["era"] = e.Era
	result["tip"] = e.Tip
	result["params"] = e.Params
	e.Value = result
}
