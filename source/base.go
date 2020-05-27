package source

var BaseType = `{
  "MetadataVersion": {
    "type": "enum",
    "value_list": [
      "MetadataV0Decoder",
      "MetadataV1Decoder",
      "MetadataV2Decoder",
      "MetadataV3Decoder",
      "MetadataV4Decoder",
      "MetadataV5Decoder",
      "MetadataV6Decoder",
      "MetadataV7Decoder",
      "MetadataV8Decoder",
      "MetadataV9Decoder",
      "MetadataV10Decoder",
      "MetadataV11Decoder"
    ]
  },
  "StorageModify": {
    "type": "enum",
    "value_list": [
      "Optional",
      "Default"
    ]
  },
  "StorageFunctionType": {
    "type": "enum",
    "value_list": [
      "PlainType",
      "MapType",
      "DoubleMapType"
    ]
  },
  "SetId": "U64",
  "RoundNumber": "U64",
  "SessionIndex": "U32",
  "AuctionIndex": "U32",
  "AuthIndex": "U32",
  "AuthorityIndex": "u64",
  "AuthorityWeight": "u64",
  "CollatorId": "H256",
  "EventIndex": "u32",
  "LeasePeriod": "BlockNumber",
  "LeasePeriodOf": "BlockNumber",
  "MaybeVrf": "[u8; 32]",
  "MemberCount": "u32",
  "MomentOf": "Moment",
  "MoreAttestations": "Null",
  "Multiplier": "u64",
  "PhantomData": "Null",
  "Reporter": "AccountId",
  "RegistrarIndex": "u32",
  "ReportIdOf": "Hash",
  "SubId": "u32",
  "Weight": "u32",
  "WeightMultiplier": "u64",
  "Index": "U32",
  "Kind": "[u8; 16]",
  "OpaqueTimeSlot": "Bytes",
  "Box<Call>": "BoxProposal",
  "Box<<T as Trait<I>>::Proposal>": "BoxProposal",
  "<AuthorityId as RuntimeAppPublic>::Signature": "Signature",
  "&[u8]": "Bytes",
  "ConsensusEngineId": "[u8; 4]",
  "SpanIndex": "u32",
  "StrikeCount": "u32",
  "RefCount": "u8",
  "ValidatorId": "AccountId",
  "CallHash": "H256",
  "PhragmenScore": "[u128; 3]",
  "schnorrkel::Randomness": "Hash",
  "schnorrkel::RawVRFOutput": "[u8; 32]",
  "NominatorIndex": "u32",
  "PerU16": "u16",
  "OffchainAccuracy": "u16",
  "AuthorityList": "Vec<NextAuthority>",
  "BalanceOf": "Balance",
  "ProposalIndex": "u32",
  "OpaquePeerId": "Bytes",
  "OpaqueMultiaddr": "Bytes",
  "ProposalTitle": "Bytes",
  "ProposalContents": "Bytes",
  "StorageKey": "Bytes",
  "DoNotConstruct": "Null",
  "Origin": "Null",
  "Randomness": "Hash",
  "MaybeRandomness": "Option<Randomness>",
  "Perbill": "u32",
  "Proposal": "BoxProposal",
  "AuthoritySignature": "Signature",
  "VrfData": "[u8; 32]",
  "VrfProof": "[u8; 64]",
  "KeyValue": {
    "type": "struct",
    "type_mapping": [
      [
        "key",
        "Vec<u8>"
      ],
      [
        "value",
        "Vec<u8>"
      ]
    ]
  },
  "UnlockChunk": {
    "type": "struct",
    "type_mapping": [
      [
        "value",
        "Compact<Balance>"
      ],
      [
        "era",
        "Compact<EraIndex>"
      ]
    ]
  },
  "IndividualExposure": {
    "type": "struct",
    "type_mapping": [
      [
        "who",
        "AccountId"
      ],
      [
        "value",
        "Compact<Balance>"
      ]
    ]
  },
  "OpaqueNetworkState": {
    "type": "struct",
    "type_mapping": [
      [
        "peerId",
        "OpaquePeerId"
      ],
      [
        "externalAddresses",
        "Vec<OpaqueMultiaddr>"
      ]
    ]
  },
  "(AccountId, Balance)": {
    "type": "struct",
    "type_mapping": [
      [
        "account",
        "AccountId"
      ],
      [
        "balance",
        "Balance"
      ]
    ]
  },
  "AccountData": {
    "type": "struct",
    "type_mapping": [
      [
        "free",
        "Balance"
      ],
      [
        "reserved",
        "Balance"
      ],
      [
        "miscFrozen",
        "Balance"
      ],
      [
        "feeFrozen",
        "Balance"
      ]
    ]
  },
  "CandidateReceipt": {
    "type": "struct",
    "type_mapping": [
      [
        "parachainIndex",
        "ParaId"
      ],
      [
        "collator",
        "AccountId"
      ],
      [
        "signature",
        "CollatorSignature"
      ],
      [
        "headData",
        "HeadData"
      ],
      [
        "balanceUploads",
        "Vec<BalanceUpload>"
      ],
      [
        "egressQueueRoots",
        "Vec<EgressQueueRoot>"
      ],
      [
        "fees",
        "u64"
      ],
      [
        "blockDataHash",
        "Hash"
      ]
    ]
  },
  "AttestedCandidate": {
    "type": "struct",
    "type_mapping": [
      [
        "candidate",
        "CandidateReceipt"
      ],
      [
        "validityVotes",
        "Vec<ValidityVote>"
      ]
    ]
  },
  "FullIdentification": {
    "type": "struct",
    "type_mapping": [
      [
        "total",
        "Compact<Balance>"
      ],
      [
        "own",
        "Compact<Balance>"
      ],
      [
        "others",
        "Vec<IndividualExposure>"
      ]
    ]
  },
  "IdentificationTuple": {
    "type": "struct",
    "type_mapping": [
      [
        "validatorId",
        "ValidatorId"
      ],
      [
        "exposure",
        "FullIdentification"
      ]
    ]
  },
  "Reasons": {
    "type": "enum",
    "value_list": [
      "Fee",
      "Misc",
      "All"
    ]
  },
  "NextAuthority": {
    "type": "struct",
    "type_mapping": [
      [
        "AuthorityId",
        "AuthorityId"
      ],
      [
        "weight",
        "AuthorityWeight"
      ]
    ]
  },
  "BalanceUpload": {
    "type": "struct",
    "type_mapping": [
      [
        "col1",
        "AccountId"
      ],
      [
        "col2",
        "u64"
      ]
    ]
  },
  "DispatchClass": {
    "type": "enum",
    "value_list": [
      "Normal",
      "Operational"
    ]
  },
  "DispatchInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "weight",
        "Weight"
      ],
      [
        "class",
        "DispatchClass"
      ],
      [
        "paysFee",
        "bool"
      ]
    ]
  },
  "EgressQueueRoot": {
    "type": "struct",
    "type_mapping": [
      [
        "col1",
        "ParaId"
      ],
      [
        "col2",
        "Hash"
      ]
    ]
  },
  "IdentityFields": {
    "type": "set",
    "value_list": [
      "Display",
      "Legal",
      "Web",
      "Riot",
      "Email",
      "PgpFingerprint",
      "Image",
      "Twitter"
    ]
  },
  "IdentityInfoAdditional": {
    "type": "struct",
    "type_mapping": [
      [
        "field",
        "Data"
      ],
      [
        "value",
        "Data"
      ]
    ]
  },
  "IdentityInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "additional",
        "Vec<IdentityInfoAdditional>"
      ],
      [
        "display",
        "Data"
      ],
      [
        "legal",
        "Data"
      ],
      [
        "web",
        "Data"
      ],
      [
        "riot",
        "Data"
      ],
      [
        "email",
        "Data"
      ],
      [
        "pgpFingerprint",
        "Option<H160>"
      ],
      [
        "image",
        "Data"
      ],
      [
        "twitter",
        "Data"
      ]
    ]
  },
  "Judgement": {
    "type": "enum",
    "type_mapping": [
      [
        "Unknown",
        "Null"
      ],
      [
        "FeePaid",
        "Balance"
      ],
      [
        "Reasonable",
        "Null"
      ],
      [
        "KnownGood",
        "Null"
      ],
      [
        "OutOfDate",
        "Null"
      ],
      [
        "LowQuality",
        "Null"
      ],
      [
        "Erroneous",
        "Null"
      ]
    ]
  },
  "Offender": {
    "type": "struct",
    "type_mapping": [
      [
        "col1",
        "ValidatorId"
      ],
      [
        "col2",
        "Exposure"
      ]
    ]
  },
  "OffenceDetails<AccountId, IdentificationTuple>": {
    "type": "struct",
    "type_mapping": [
      [
        "offender",
        "Offender"
      ],
      [
        "reporters",
        "Vec<Reporter>"
      ]
    ]
  },
  "OpenTip<AccountId, BalanceOf, BlockNumber, Hash>": {
    "type": "struct",
    "type_mapping": [
      [
        "reason",
        "Hash"
      ],
      [
        "who",
        "AccountId"
      ],
      [
        "finder",
        "Option<OpenTipFinder>"
      ],
      [
        "closes",
        "Option<BlockNumber>"
      ],
      [
        "tips",
        "Vec<OpenTipTip>"
      ]
    ]
  },
  "ParaScheduling": {
    "type": "enum",
    "value_list": [
      "Always",
      "Dynamic"
    ]
  },
  "ParaInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "scheduling",
        "ParaScheduling"
      ]
    ]
  },
  "ReferendumInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "end",
        "BlockNumber"
      ],
      [
        "proposal",
        "Proposal"
      ],
      [
        "threshold",
        "VoteThreshold"
      ],
      [
        "delay",
        "BlockNumber"
      ]
    ]
  },
  "ReferendumInfo<BlockNumber, Proposal>": {
    "type": "struct",
    "type_mapping": [
      [
        "end",
        "BlockNumber"
      ],
      [
        "proposal",
        "Proposal"
      ],
      [
        "threshold",
        "VoteThreshold"
      ],
      [
        "delay",
        "BlockNumber"
      ]
    ]
  },
  "ReferendumInfo<BlockNumber, Hash>": {
    "type": "struct",
    "type_mapping": [
      [
        "end",
        "BlockNumber"
      ],
      [
        "proposalHash",
        "Hash"
      ],
      [
        "threshold",
        "VoteThreshold"
      ],
      [
        "delay",
        "BlockNumber"
      ]
    ]
  },
  "RegistrarInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "account",
        "AccountId"
      ],
      [
        "fee",
        "Balance"
      ],
      [
        "fields",
        "IdentityFields"
      ]
    ]
  },
  "Registration": {
    "type": "struct",
    "type_mapping": [
      [
        "judgements",
        "Vec<RegistrationJudgement>"
      ],
      [
        "deposit",
        "Balance"
      ],
      [
        "info",
        "IdentityInfo"
      ]
    ]
  },
  "RegistrationJudgement": {
    "type": "struct",
    "type_mapping": [
      [
        "col1",
        "RegistrarIndex"
      ],
      [
        "col2",
        "Judgement"
      ]
    ]
  },
  "WinningDataEntry": {
    "type": "struct",
    "type_mapping": [
      [
        "col1",
        "AccountId"
      ],
      [
        "col2",
        "ParaIdOf"
      ],
      [
        "col3",
        "BalanceOf"
      ]
    ]
  },
  "WithdrawReasons": {
    "type": "set",
    "value_list": [
      "TransactionPayment",
      "Transfer",
      "Reserve",
      "Fee",
      "Tip"
    ]
  },
  "Nominations": {
    "type": "struct",
    "type_mapping": [
      [
        "targets",
        "Vec<AccountId>"
      ],
      [
        "submittedIn",
        "EraIndex"
      ],
      [
        "suppressed",
        "bool"
      ]
    ]
  },
  "Forcing": {
    "type": "enum",
    "value_list": [
      "NotForcing",
      "ForceNew",
      "ForceNone"
    ]
  },
  "VoteStage": {
    "type": "enum",
    "value_list": [
      "PreVoting",
      "Commit",
      "Voting",
      "Completed"
    ]
  },
  "ProposalStage": {
    "type": "enum",
    "value_list": [
      "PreVoting",
      "Voting",
      "Completed"
    ]
  },
  "ProposalCategory": {
    "type": "enum",
    "value_list": [
      "Signaling"
    ]
  },
  "Heartbeat": {
    "type": "struct",
    "type_mapping": [
      [
        "blockNumber",
        "BlockNumber"
      ],
      [
        "networkState",
        "OpaqueNetworkState"
      ],
      [
        "sessionIndex",
        "SessionIndex"
      ],
      [
        "authorityIndex",
        "AuthIndex"
      ]
    ]
  },
  "RewardDestination": {
    "type": "enum",
    "value_list": [
      "Staked",
      "Stash",
      "Controller"
    ]
  },
  "DigestItem": {
    "type": "enum",
    "type_mapping": [
      [
        "Other",
        "Vec<u8>"
      ],
      [
        "AuthoritiesChange",
        "Vec<AuthorityId>"
      ],
      [
        "ChangesTrieRoot",
        "Hash"
      ],
      [
        "SealV0",
        "SealV0"
      ],
      [
        "Consensus",
        "Consensus"
      ],
      [
        "Seal",
        "Seal"
      ],
      [
        "PreRuntime",
        "PreRuntime"
      ]
    ]
  },
  "Digest": {
    "type": "struct",
    "type_mapping": [
      [
        "logs",
        "Vec<DigestItem<Hash>>"
      ]
    ]
  },
  "SpanRecord": {
    "type": "struct",
    "type_mapping": [
      [
        "slashed",
        "Balance"
      ],
      [
        "paidOut",
        "Balance"
      ]
    ]
  },
  "UnappliedSlashOther": {
    "type": "struct",
    "type_mapping": [
      [
        "account",
        "AccountId"
      ],
      [
        "amount",
        "Balance"
      ]
    ]
  },
  "UnappliedSlash": {
    "type": "struct",
    "type_mapping": [
      [
        "validator",
        "AccountId"
      ],
      [
        "own",
        "AccountId"
      ],
      [
        "others",
        "Vec<UnappliedSlashOther>"
      ],
      [
        "reporters",
        "Vec<AccountId>"
      ],
      [
        "payout",
        "Balance"
      ]
    ]
  },
  "Header": {
    "type": "struct",
    "type_mapping": [
      [
        "parent_hash",
        "H256"
      ],
      [
        "number",
        "Compact<BlockNumber>"
      ],
      [
        "state_root",
        "H256"
      ],
      [
        "extrinsics_root",
        "H256"
      ],
      [
        "digest",
        "Digest"
      ]
    ]
  },
  "DispatchErrorModule": {
    "type": "struct",
    "type_mapping": [
      [
        "index",
        "u8"
      ],
      [
        "error",
        "u8"
      ]
    ]
  },
  "DispatchError": {
    "type": "enum",
    "type_mapping": [
      [
        "Other",
        "Null"
      ],
      [
        "CannotLookup",
        "Null"
      ],
      [
        "BadOrigin",
        "Null"
      ],
      [
        "Module",
        "DispatchErrorModule"
      ]
    ]
  },
  "DispatchResult": {
    "type": "enum",
    "type_mapping": [
      [
        "Error",
        "DispatchError"
      ],
      [
        "Ok",
        "Null"
      ]
    ]
  },
  "ActiveRecovery": {
    "type": "struct",
    "type_mapping": [
      [
        "created",
        "BlockNumber"
      ],
      [
        "deposit",
        "Balance"
      ],
      [
        "friends",
        "Vec<AccountId>"
      ]
    ]
  },
  "RecoveryConfig": {
    "type": "struct",
    "type_mapping": [
      [
        "delayPeriod",
        "BlockNumber"
      ],
      [
        "deposit",
        "Balance"
      ],
      [
        "friends",
        "Vec<AccountId>"
      ],
      [
        "threshold",
        "u16"
      ]
    ]
  },
  "BidKindVouch": {
    "type": "struct",
    "type_mapping": [
      [
        "account",
        "AccountId"
      ],
      [
        "amount",
        "Balance"
      ]
    ]
  },
  "BidKind": {
    "type": "enum",
    "type_mapping": [
      [
        "Deposit",
        "Balance"
      ],
      [
        "Vouch",
        "BidKindVouch"
      ]
    ]
  },
  "Bid": {
    "type": "struct",
    "type_mapping": [
      [
        "who",
        "AccountId"
      ],
      [
        "kind",
        "BidKind"
      ],
      [
        "value",
        "Balance"
      ]
    ]
  },
  "VouchingStatus": {
    "type": "enum",
    "value_list": [
      "Vouching",
      "Banned"
    ]
  },
  "ExtrinsicMetadata": {
    "type": "struct",
    "type_mapping": [
      [
        "version",
        "u8"
      ],
      [
        "signedExtensions",
        "Vec<Bytes>"
      ]
    ]
  },
  "AccountInfo<Index, AccountData>": {
    "type": "struct",
    "type_mapping": [
      [
        "nonce",
        "Index"
      ],
      [
        "refcount",
        "RefCount"
      ],
      [
        "data",
        "AccountData"
      ]
    ]
  },
  "ActiveEraInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "index",
        "EraIndex"
      ],
      [
        "start",
        "MomentOf"
      ]
    ]
  },
  "(BalanceOf<T, I>, BidKind<AccountId, BalanceOf<T, I>>)": {
    "type": "struct",
    "type_mapping": [
      [
        "balance",
        "Balance"
      ],
      [
        "bidkind",
        "BidKind"
      ]
    ]
  },
  "StakingLedger<AccountId, BalanceOf>": {
    "type": "struct",
    "type_mapping": [
      [
        "stash",
        "AccountId"
      ],
      [
        "total",
        "Compact<Balance>"
      ],
      [
        "active",
        "Compact<Balance>"
      ],
      [
        "unlocking",
        "Vec<UnlockChunk<Balance>>"
      ],
      [
        "lastReward",
        "Option<EraIndex>"
      ]
    ]
  },
  "LastRuntimeUpgradeInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "specVersion",
        "Compact<u32>"
      ],
      [
        "specName",
        "Bytes"
      ]
    ]
  },
  "AccountVoteSplit": {
    "type": "struct",
    "type_mapping": [
      [
        "aye",
        "Balance"
      ],
      [
        "nay",
        "Balance"
      ]
    ]
  },
  "AccountVoteStandard": {
    "type": "struct",
    "type_mapping": [
      [
        "vote",
        "Vote"
      ],
      [
        "balance",
        "Balance"
      ]
    ]
  },
  "AccountVote": {
    "type": "enum",
    "type_mapping": [
      [
        "Standard",
        "AccountVoteStandard"
      ],
      [
        "Split",
        "AccountVoteSplit"
      ]
    ]
  },
  "Delegations": {
    "type": "struct",
    "type_mapping": [
      [
        "votes",
        "Balance"
      ],
      [
        "capital",
        "Balance"
      ]
    ]
  },
  "PriorLock": {
    "type": "struct",
    "type_mapping": [
      [
        "blockNumber",
        "BlockNumber"
      ],
      [
        "balance",
        "Balance"
      ]
    ]
  },
  "ReferendumInfo<BlockNumber, Hash, Balance>": {
    "type": "enum",
    "type_mapping": [
      [
        "Ongoing",
        "ReferendumStatus"
      ],
      [
        "Finished",
        "ReferendumInfoFinished"
      ]
    ]
  },
  "ReferendumInfoFinished": {
    "type": "struct",
    "type_mapping": [
      [
        "approved",
        "bool"
      ],
      [
        "end",
        "BlockNumber"
      ]
    ]
  },
  "Tally": {
    "type": "struct",
    "type_mapping": [
      [
        "ayes",
        "Balance"
      ],
      [
        "nays",
        "Balance"
      ],
      [
        "turnout",
        "Balance"
      ]
    ]
  },
  "ReferendumStatus": {
    "type": "struct",
    "type_mapping": [
      [
        "end",
        "BlockNumber"
      ],
      [
        "proposalHash",
        "Hash"
      ],
      [
        "threshold",
        "VoteThreshold"
      ],
      [
        "delay",
        "BlockNumber"
      ],
      [
        "tally",
        "Tally"
      ]
    ]
  },
  "VotingDirectVote": {
    "type": "struct",
    "type_mapping": [
      [
        "referendumIndex",
        "ReferendumIndex"
      ],
      [
        "accountVote",
        "AccountVote"
      ]
    ]
  },
  "VotingDirect": {
    "type": "struct",
    "type_mapping": [
      [
        "votes",
        "Vec<VotingDirectVote>"
      ],
      [
        "delegations",
        "Delegations"
      ],
      [
        "prior",
        "PriorLock"
      ]
    ]
  },
  "VotingDelegating": {
    "type": "struct",
    "type_mapping": [
      [
        "balance",
        "Balance"
      ],
      [
        "target",
        "AccountId"
      ],
      [
        "conviction",
        "Conviction"
      ],
      [
        "delegations",
        "Delegations"
      ],
      [
        "prior",
        "PriorLock"
      ]
    ]
  },
  "Voting": {
    "type": "enum",
    "type_mapping": [
      [
        "Direct",
        "VotingDirect"
      ],
      [
        "Delegating",
        "VotingDelegating"
      ]
    ]
  },
  "Timepoint": {
    "type": "struct",
    "type_mapping": [
      [
        "height",
        "BlockNumber"
      ],
      [
        "index",
        "u32"
      ]
    ]
  },
  "Multisig": {
    "type": "struct",
    "type_mapping": [
      [
        "when",
        "Timepoint"
      ],
      [
        "deposit",
        "Balance"
      ],
      [
        "depositor",
        "AccountId"
      ],
      [
        "approvals",
        "Vec<AccountId>"
      ]
    ]
  },
  "SessionKeysPolkadot": {
    "type": "struct",
    "type_mapping": [
      [
        "grandpa",
        "AccountId"
      ],
      [
        "babe",
        "AccountId"
      ],
      [
        "im_online",
        "AccountId"
      ],
      [
        "authority_discovery",
        "AccountId"
      ],
      [
        "parachains",
        "AccountId"
      ]
    ]
  },
  "BalanceLock": {
    "type": "struct",
    "type_mapping": [
      [
        "id",
        "LockIdentifier"
      ],
      [
        "amount",
        "Balance"
      ],
      [
        "until",
        "u32"
      ],
      [
        "reasons",
        "WithdrawReasons"
      ]
    ]
  },
  "BalanceLock<Balance>": {
    "type": "struct",
    "type_mapping": [
      [
        "id",
        "LockIdentifier"
      ],
      [
        "amount",
        "Balance"
      ],
      [
        "reasons",
        "Reasons"
      ]
    ]
  },
  "VoteThreshold": {
    "type": "enum",
    "value_list": [
      "SuperMajorityApprove",
      "SuperMajorityAgainst",
      "SimpleMajority"
    ]
  },
  "ModuleId": "LockIdentifier",
  "RuntimeDbWeight": {
    "type": "struct",
    "type_mapping": [
      [
        "read",
        "Weight"
      ],
      [
        "write",
        "Weight"
      ]
    ]
  },
  "ElectionCompute": {
    "type": "enum",
    "value_list": [
      "OnChain",
      "Signed",
      "Authority"
    ]
  },
  "ElectionResult": {
    "type": "struct",
    "type_mapping": [
      [
        "compute",
        "ElectionCompute"
      ],
      [
        "slotStake",
        "Balance"
      ],
      [
        "electedStashes",
        "Vec<AccountId>"
      ],
      [
        "exposures",
        "Vec<(AccountId, Exposure)>"
      ]
    ]
  },
  "ElectionStatus": {
    "type": "enum",
    "type_mapping": [
      [
        "Close",
        "Null"
      ],
      [
        "Open",
        "BlockNumber"
      ]
    ]
  },
  "Phase": {
    "type": "enum",
    "type_mapping": [
      [
        "ApplyExtrinsic",
        "u32"
      ],
      [
        "Finalization",
        "Null"
      ],
      [
        "Initialization",
        "Null"
      ]
    ]
  },
  "PreimageStatusAvailable": {
    "type": "struct",
    "type_mapping": [
      [
        "data",
        "Bytes"
      ],
      [
        "provider",
        "AccountId"
      ],
      [
        "deposit",
        "Balance"
      ],
      [
        "since",
        "BlockNumber"
      ],
      [
        "expiry",
        "Option<BlockNumber>"
      ]
    ]
  },
  "PreimageStatus": {
    "type": "enum",
    "type_mapping": [
      [
        "Missing",
        "BlockNumber"
      ],
      [
        "Available",
        "PreimageStatusAvailable"
      ]
    ]
  },
  "TaskAddress": {
    "type": "struct",
    "type_mapping": [
      [
        "blockNumber",
        "BlockNumber"
      ],
      [
        "index",
        "u32"
      ]
    ]
  },
  "ValidationCode": "Bytes",
  "ValidatorIndex": "u16",
  "ParaPastCodeMeta": {
    "type": "struct",
    "type_mapping": [
      [
        "upgrade_times",
        "Vec<BlockNumber>"
      ],
      [
        "last_pruned",
        "Option<BlockNumber>"
      ]
    ]
  },
  "CompactScore": {
    "type": "struct",
    "type_mapping": [
      [
        "validatorIndex",
        "ValidatorIndex"
      ],
      [
        "offchainAccuracy",
        "OffchainAccuracy"
      ]
    ]
  },
  "CompactAssignmentsVote": {
    "type": "struct",
    "type_mapping": [
      [
        "accountId1",
        "AccountId"
      ],
      [
        "scores",
        "Vec<CompactScore>"
      ],
      [
        "accountId2",
        "AccountId"
      ]
    ]
  },
  "CompactAssignments": {
    "type": "struct",
    "type_mapping": [
      [
        "votes1",
        "Vec<(NominatorIndex, ValidatorIndex)>"
      ],
      [
        "votes2",
        "Vec<(NominatorIndex, [CompactScore; 1], ValidatorIndex)>"
      ],
      [
        "votes3",
        "Vec<(NominatorIndex, [CompactScore; 2], ValidatorIndex)>"
      ],
      [
        "votes4",
        "Vec<(NominatorIndex, [CompactScore; 3], ValidatorIndex)>"
      ],
      [
        "votes5",
        "Vec<(NominatorIndex, [CompactScore; 4], ValidatorIndex)>"
      ],
      [
        "votes6",
        "Vec<(NominatorIndex, [CompactScore; 5], ValidatorIndex)>"
      ],
      [
        "votes7",
        "Vec<(NominatorIndex, [CompactScore; 6], ValidatorIndex)>"
      ],
      [
        "votes8",
        "Vec<(NominatorIndex, [CompactScore; 7], ValidatorIndex)>"
      ],
      [
        "votes9",
        "Vec<(NominatorIndex, [CompactScore; 8], ValidatorIndex)>"
      ],
      [
        "votes10",
        "Vec<(NominatorIndex, [CompactScore; 9], ValidatorIndex)>"
      ],
      [
        "votes11",
        "Vec<(NominatorIndex, [CompactScore; 10], ValidatorIndex)>"
      ],
      [
        "votes12",
        "Vec<(NominatorIndex, [CompactScore; 11], ValidatorIndex)>"
      ],
      [
        "votes13",
        "Vec<(NominatorIndex, [CompactScore; 12], ValidatorIndex)>"
      ],
      [
        "votes14",
        "Vec<(NominatorIndex, [CompactScore; 13], ValidatorIndex)>"
      ],
      [
        "votes15",
        "Vec<(NominatorIndex, [CompactScore; 14], ValidatorIndex)>"
      ],
      [
        "votes16",
        "Vec<(NominatorIndex, [CompactScore; 15], ValidatorIndex)>"
      ]
    ]
  },
  "ReleasesBalances": {
    "type": "enum",
    "value_list": [
      "V1_0_0",
      "V2_0_0"
    ]
  },
  "Releases": "ReleasesBalances",
  "SlotRange": {
    "type": "enum",
    "value_list": [
      "ZeroZero",
      "ZeroOne",
      "ZeroTwo",
      "ZeroThree",
      "OneOne",
      "OneTwo",
      "OneThree",
      "TwoTwo",
      "TwoThree",
      "ThreeThree"
    ]
  },
  "ValidityAttestation": {
    "type": "enum",
    "type_mapping": [
      [
        "None",
        "Null"
      ],
      [
        "Implicit",
        "CollatorSignature"
      ],
      [
        "Explicit",
        "CollatorSignature"
      ]
    ]
  },
  "VestingInfo": {
    "type": "struct",
    "type_mapping": [
      [
        "locked",
        "Balance"
      ],
      [
        "perBlock",
        "Balance"
      ],
      [
        "startingBlock",
        "BlockNumber"
      ]
    ]
  },
  "ProxyState": {
    "type": "struct",
    "type_mapping": [
      [
        "Open",
        "AccountId"
      ],
      [
        "Active",
        "AccountId"
      ]
    ]
  },
  "IncomingParachainFixed": {
    "type": "struct",
    "type_mapping": [
      [
        "codeHash",
        "Hash"
      ],
      [
        "initialHeadData",
        "Bytes"
      ]
    ]
  },
  "IncomingParachain": {
    "type": "enum",
    "type_mapping": [
      [
        "Unset",
        "NewBidder"
      ],
      [
        "Fixed",
        "IncomingParachainFixed"
      ],
      [
        "Deploy",
        "IncomingParachainDeploy"
      ]
    ]
  },
  "Pays": {
    "type": "enum",
    "value_list": [
      "Yes",
      "No"
    ]
  },
  "Renouncing": {
    "type": "enum",
    "type_mapping": [
      [
        "Member",
        "Null"
      ],
      [
        "RunnerUp",
        "Null"
      ],
      [
        "Candidate",
        "Compact<u32>"
      ]
    ]
  },
  "ExtrinsicsWeight": {
    "type": "struct",
    "type_mapping": [
      [
        "normal",
        "Weight"
      ],
      [
        "operational",
        "Weight"
      ]
    ]
  },
  "ValidatorCount": "u32",
  "KeyOwnerProof": {
    "type": "struct",
    "type_mapping": [
      [
        "session",
        "SessionIndex"
      ],
      [
        "trieNodes",
        "Vec<Bytes>"
      ],
      [
        "validatorCount",
        "ValidatorCount"
      ]
    ]
  },
  "DefunctVoter": {
    "type": "struct",
    "type_mapping": [
      [
        "who",
        "AccountId"
      ],
      [
        "voteCount",
        "Compact<u32>"
      ],
      [
        "candidateCount",
        "Compact<u32>"
      ]
    ]
  },
  "ElectionSize": {
    "type": "struct",
    "type_mapping": [
      [
        "validators",
        "ValidatorIndex"
      ],
      [
        "nominators",
        "NominatorIndex"
      ]
    ]
  },
  "AllowedSlots": {
    "type": "enum",
    "value_list": [
      "PrimarySlots",
      "PrimaryAndSecondaryPlainSlots",
      "PrimaryAndSecondaryVRFSlots"
    ]
  },
  "NextConfigDescriptorV1": {
    "type": "struct",
    "type_mapping": [
      [
        "c",
        "(u64, u64)"
      ],
      [
        "allowedSlots",
        "AllowedSlots"
      ]
    ]
  },
  "NextConfigDescriptor": {
    "type": "enum",
    "type_mapping": [
      [
        "V0",
        "Null"
      ],
      [
        "V1",
        "NextConfigDescriptorV1"
      ]
    ]
  },
  "StatementKind": {
    "type": "enum",
    "value_list": [
      "Regular",
      "Saft"
    ]
  },
  "schedule::Priority": "u8",
  "GrandpaEquivocation": {
    "type": "struct",
    "type_mapping": [
      [
        "roundNumber",
        "u64"
      ],
      [
        "identity",
        "AuthorityId"
      ],
      [
        "first",
        "(GrandpaPrevote, AuthoritySignature)"
      ],
      [
        "second",
        "(GrandpaPrevote, AuthoritySignature)"
      ]
    ]
  },
  "GrandpaPrevote": {
    "type": "struct",
    "type_mapping": [
      [
        "targetHash",
        "Hash"
      ],
      [
        "targetNumber",
        "BlockNumber"
      ]
    ]
  },
  "Equivocation": {
    "type": "enum",
    "type_mapping": [
      [
        "Prevote",
        "GrandpaEquivocation"
      ],
      [
        "Precommit",
        "GrandpaEquivocation"
      ]
    ]
  },
  "EquivocationProof": {
    "type": "struct",
    "type_mapping": [
      [
        "setId",
        "SetId"
      ],
      [
        "equivocation",
        "Equivocation"
      ]
    ]
  },
  "ValidatorPrefs": {
    "type": "struct",
    "type_mapping": [
      [
        "Commission",
        "Compact<Balance>"
      ]
    ]
  }
}`
