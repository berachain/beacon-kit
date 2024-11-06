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

package types

import (
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
)

type ValidatorResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                any  `json:"data"`
}

// BlockHeader contains the block header details
// that resides in BlockHeaderResponse.
type BlockHeader[BlockHeaderT BeaconBlockHeader] struct {
	Message   BlockHeaderT        `json:"message"`
	Signature crypto.BLSSignature `json:"signature"`
}

// BlockHeaderResponse contains the block header response for the Beacon API.
type BlockHeaderResponse[BlockHeaderT BeaconBlockHeader] struct {
	Root      common.Root                `json:"root"`
	Canonical bool                       `json:"canonical"`
	Header    *BlockHeader[BlockHeaderT] `json:"header"`
}

type messageJSON struct {
	Slot          string      `json:"slot"`
	ProposerIndex string      `json:"proposer_index"`
	ParentRoot    common.Root `json:"parent_root"`
	StateRoot     common.Root `json:"state_root"`
	BodyRoot      common.Root `json:"body_root"`
}

type blockHeaderResponseJSON struct {
	Message   messageJSON         `json:"message"`
	Signature crypto.BLSSignature `json:"signature"`
}

// MarshalJSON implements custom JSON marshaling for BlockHeader.
func (bh *BlockHeader[BlockHeaderT]) MarshalJSON() ([]byte, error) {
	return json.Marshal(&blockHeaderResponseJSON{
		Message: messageJSON{
			Slot: strconv.FormatUint(
				bh.Message.GetSlot().Unwrap(), 10,
			),
			ProposerIndex: strconv.FormatUint(
				bh.Message.GetProposerIndex().Unwrap(), 10,
			),
			ParentRoot: bh.Message.GetParentBlockRoot(),
			StateRoot:  bh.Message.GetStateRoot(),
			BodyRoot:   bh.Message.GetBodyRoot(),
		},
		Signature: bh.Signature,
	})
}

type GenesisData struct {
	GenesisTime           string      `json:"genesis_time"`
	GenesisValidatorsRoot common.Root `json:"genesis_validators_root"`
	GenesisForkVersion    string      `json:"genesis_fork_version"`
}

type RootData struct {
	Root common.Root `json:"root"`
}

type ValidatorBalanceData struct {
	Index   uint64 `json:"index"`
	Balance uint64 `json:"balance"`
}

func (vbd ValidatorBalanceData) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Index   string `json:"index"`
		Balance string `json:"balance"`
	}{
		Index:   strconv.FormatUint(vbd.Index, 10),
		Balance: strconv.FormatUint(vbd.Balance, 10),
	})
}

type ValidatorData[
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	ValidatorBalanceData
	Status    string     `json:"status"`
	Validator ValidatorT `json:"validator"`
}

type validatorJSON struct {
	PublicKey                  string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

type responseJSON struct {
	Index     string        `json:"index"`
	Balance   string        `json:"balance"`
	Status    string        `json:"status"`
	Validator validatorJSON `json:"validator"`
}

func (vd ValidatorData[
	ValidatorT, WithdrawalCredentialsT,
]) MarshalJSON() ([]byte, error) {
	withdrawalCredentials := vd.Validator.GetWithdrawalCredentials()
	withdrawalCredentialsBytes := withdrawalCredentials.Bytes()

	return json.Marshal(responseJSON{
		Index:   strconv.FormatUint(vd.Index, 10),
		Balance: strconv.FormatUint(vd.Balance, 10),
		Status:  vd.Status,
		Validator: validatorJSON{
			PublicKey: vd.Validator.GetPubkey().String(),
			WithdrawalCredentials: "0x" + hex.EncodeToString(
				withdrawalCredentialsBytes,
			),
			EffectiveBalance: strconv.FormatUint(
				vd.Validator.GetEffectiveBalance().Unwrap(), 10,
			),
			Slashed: vd.Validator.IsSlashed(),
			ActivationEligibilityEpoch: strconv.FormatUint(
				vd.Validator.GetActivationEligibilityEpoch().Unwrap(), 10,
			),
			ActivationEpoch: strconv.FormatUint(
				vd.Validator.GetActivationEpoch().Unwrap(), 10,
			),
			ExitEpoch: strconv.FormatUint(
				vd.Validator.GetExitEpoch().Unwrap(), 10,
			),
			WithdrawableEpoch: strconv.FormatUint(
				vd.Validator.GetWithdrawableEpoch().Unwrap(), 10,
			),
		},
	})
}

type BlockRewardsData struct {
	ProposerIndex     uint64 `json:"proposer_index,string"`
	Total             uint64 `json:"total,string"`
	Attestations      uint64 `json:"attestations,string"`
	SyncAggregate     uint64 `json:"sync_aggregate,string"`
	ProposerSlashings uint64 `json:"proposer_slashings,string"`
	AttesterSlashings uint64 `json:"attester_slashings,string"`
}

type RandaoData struct {
	Randao common.Bytes32 `json:"randao"`
}

type ForkData struct {
	Fork
}

type forkJSON struct {
	PreviousVersion string `json:"previous_version"`
	CurrentVersion  string `json:"current_version"`
	Epoch           string `json:"epoch"`
}

func (fr ForkData) MarshalJSON() ([]byte, error) {
	return json.Marshal(forkJSON{
		PreviousVersion: fr.GetPreviousVersion().String(),
		CurrentVersion:  fr.GetCurrentVersion().String(),
		Epoch:           strconv.FormatUint(fr.GetEpoch().Unwrap(), 10),
	})
}

type FinalityCheckpointsData struct {
	PreviousJustified common.Checkpoint `json:"previous_justified"`
	CurrentJustified  common.Checkpoint `json:"current_justified"`
	Finalized         common.Checkpoint `json:"finalized"`
}
