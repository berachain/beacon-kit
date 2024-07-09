// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

// Package ssz contains fixture from fast ssz
//
//nolint:tagliatelle,lll // imported test
package ssz

type BeaconBlockHeader struct {
	Slot          uint64 `json:"slot"`
	ProposerIndex uint64 `json:"proposer_index"`
	ParentRoot    []byte `json:"parent_root"    ssz-size:"32"`
	StateRoot     []byte `json:"state_root"     ssz-size:"32"`
	BodyRoot      []byte `json:"body_root"      ssz-size:"32"`
}
type Checkpoint struct {
	Epoch uint64 `json:"epoch"`
	Root  []byte `json:"root"  ssz-size:"32"`
}
type Fork struct {
	PreviousVersion []byte `json:"previous_version" ssz-size:"4"`
	CurrentVersion  []byte `json:"current_version"  ssz-size:"4"`
	Epoch           uint64 `json:"epoch"`
}

type Validator struct {
	Pubkey                     []byte `json:"pubkey"                       ssz-size:"48"`
	WithdrawalCredentials      []byte `json:"withdrawal_credentials"       ssz-size:"32"`
	EffectiveBalance           uint64 `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch uint64 `json:"activation_eligibility_epoch"`
	ActivationEpoch            uint64 `json:"activation_epoch"`
	ExitEpoch                  uint64 `json:"exit_epoch"`
	WithdrawableEpoch          uint64 `json:"withdrawable_epoch"`
}
type SyncCommittee struct {
	PubKeys         [][]byte `json:"pubkeys"          ssz-size:"512,48"`
	AggregatePubKey [48]byte `json:"aggregate_pubkey" ssz-size:"48"`
}
type Eth1Data struct {
	DepositRoot  []byte `json:"deposit_root"  ssz-size:"32"`
	DepositCount uint64 `json:"deposit_count"`
	BlockHash    []byte `json:"block_hash"    ssz-size:"32"`
}

type ExecutionPayloadHeader struct {
	ParentHash       []byte `json:"parent_hash"       ssz-size:"32"`
	FeeRecipient     []byte `json:"fee_recipient"     ssz-size:"20"`
	StateRoot        []byte `json:"state_root"        ssz-size:"32"`
	ReceiptsRoot     []byte `json:"receipts_root"     ssz-size:"32"`
	LogsBloom        []byte `json:"logs_bloom"        ssz-size:"256"`
	PrevRandao       []byte `json:"prev_randao"       ssz-size:"32"`
	BlockNumber      uint64 `json:"block_number"`
	GasLimit         uint64 `json:"gas_limit"`
	GasUsed          uint64 `json:"gas_used"`
	Timestamp        uint64 `json:"timestamp"`
	ExtraData        []byte `json:"extra_data"                       ssz-max:"32"`
	BaseFeePerGas    []byte `json:"base_fee_per_gas"  ssz-size:"32"`
	BlockHash        []byte `json:"block_hash"        ssz-size:"32"`
	TransactionsRoot []byte `json:"transactions_root" ssz-size:"32"`
}

//go:generate go run github.com/ferranbt/fastssz/sszgen -path bellatrix.go -objs BeaconStateBellatrix -include ../../../../pkg/crypto,../../../../pkg/common,../../../../pkg/bytes,../../../constants,../../../../pkg,../../../../pkg/math -output bellatrix.ssz.go
type BeaconStateBellatrix struct {
	GenesisTime                  uint64                  `json:"genesis_time"`
	GenesisValidatorsRoot        []byte                  `json:"genesis_validators_root"         ssz-size:"32"`
	Slot                         uint64                  `json:"slot"`
	Fork                         *Fork                   `json:"fork"`
	LatestBlockHeader            *BeaconBlockHeader      `json:"latest_block_header"`
	BlockRoots                   [][]byte                `json:"block_roots"                     ssz-size:"8192,32"`
	StateRoots                   [][]byte                `json:"state_roots"                     ssz-size:"8192,32"`
	HistoricalRoots              [][]byte                `json:"historical_roots"                ssz-size:"?,32"     ssz-max:"16777216"`
	Eth1Data                     *Eth1Data               `json:"eth1_data"`
	Eth1DataVotes                []*Eth1Data             `json:"eth1_data_votes"                                     ssz-max:"2048"`
	Eth1DepositIndex             uint64                  `json:"eth1_deposit_index"`
	Validators                   []*Validator            `json:"validators"                                          ssz-max:"1099511627776"`
	Balances                     []uint64                `json:"balances"                                            ssz-max:"1099511627776"`
	RandaoMixes                  [][]byte                `json:"randao_mixes"                    ssz-size:"65536,32"`
	Slashings                    []uint64                `json:"slashings"                       ssz-size:"8192"`
	PreviousEpochParticipation   []byte                  `json:"previous_epoch_participation"                        ssz-max:"1099511627776"`
	CurrentEpochParticipation    []byte                  `json:"current_epoch_participation"                         ssz-max:"1099511627776"`
	JustificationBits            []byte                  `json:"justification_bits"              ssz-size:"1"                                cast-type:"github.com/prysmaticlabs/go-bitfield.Bitvector4"`
	PreviousJustifiedCheckpoint  *Checkpoint             `json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint   *Checkpoint             `json:"current_justified_checkpoint"`
	FinalizedCheckpoint          *Checkpoint             `json:"finalized_checkpoint"`
	InactivityScores             []uint64                `json:"inactivity_scores"                                   ssz-max:"1099511627776"`
	CurrentSyncCommittee         *SyncCommittee          `json:"current_sync_committee"`
	NextSyncCommittee            *SyncCommittee          `json:"next_sync_committee"`
	LatestExecutionPayloadHeader *ExecutionPayloadHeader `json:"latest_execution_payload_header"`
}
