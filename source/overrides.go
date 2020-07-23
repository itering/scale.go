package source

import (
	"fmt"
	"strings"
)

var overrides = map[string]map[string]string{
	"babe": {
		"equivocationProof": "BabeEquivocationProof",
	},
	"grandpa": {
		"equivocation":      "GrandpaEquivocation",
		"equivocationProof": "GrandpaEquivocationProof",
	},
	"balances": {
		"Status": "BalanceStatus",
	},
	"contract": {
		"AccountInfo": "ContractAccountInfo",
	},
	"contracts": {
		"StorageKey": "ContractStorageKey",
	},
	"identity": {
		"Judgement": "IdentityJudgement",
	},
	"parachains": {
		"Id": "ParaId",
	},
	"society": {
		"Judgement": "SocietyJudgement",
		"Vote":      "SocietyVote",
	},
	"staking": {
		"Compact": "CompactAssignments",
	},
	"treasury": {
		"Proposal": "TreasuryProposal",
	},
}

func OverrideModuleType(moduleName, typeString string) string {
	fmt.Println(moduleName, typeString)
	if module, ok := overrides[strings.ToLower(moduleName)]; ok {
		if overrideType, ok := module[typeString]; ok {
			return overrideType
		}
	}
	return typeString
}
