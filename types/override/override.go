package override

import (
	"strings"
)

var typesModules = map[string]map[string]string{
	"parasInclusion": {
		"validatorindex": "ParaValidatorIndex",
	},
	"inclusion": {
		"validatorindex": "ParaValidatorIndex",
	},
	"parasScheduler": {
		"validatorindex": "ParaValidatorIndex",
	},
	"parasShared": {
		"validatorindex": "ParaValidatorIndex",
	},
	"scheduler": {
		"validatorindex": "ParaValidatorIndex",
	},
	"paraDisputes": {
		"validatorindex": "ParaValidatorIndex",
	},
	"parasDisputes": {
		"validatorindex": "ParaValidatorIndex",
	},
	"shared": {
		"validatorindex": "ParaValidatorIndex",
	},
	"assets": {
		"approval":       "AssetApproval",
		"approvalkey":    "AssetApprovalKey",
		"balance":        "TAssetBalance",
		"destroywitness": "AssetDestroyWitness",
	},
	"beefy": {
		"authorityid": "BeefyId",
	},
	"ethereum": {
		"block":             "EthBlock",
		"header":            "EthHeader",
		"receipt":           "EthReceipt",
		"transaction":       "EthTransaction",
		"transactionstatus": "EthTransactionStatus",
	},
	// sora
	"bridgemultisig": {
		"timepoint": "BridgeTimepoint",
	},
}

func ModuleType(t, module string) string {
	if module == "" {
		return t
	}
	moduleTypes, ok := typesModules[strings.ToLower(module)]
	if !ok {
		return t
	}
	if overrideType, ok := moduleTypes[t]; ok {
		return overrideType
	}
	return t
}
