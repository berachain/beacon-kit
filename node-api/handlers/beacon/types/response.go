// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"github.com/berachain/beacon-kit/primitives/common"
)

type GenericResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                any  `json:"data"`
}

// NewResponse creates a new response with CometBFT's finality guarantees.
func NewResponse(data any) GenericResponse {
	return GenericResponse{
		// All data is finalized in CometBFT since we only return data for slots up to head
		Finalized: true,
		// Never optimistic since we only return finalized data
		ExecutionOptimistic: false,
		Data:                data,
	}
}

type BlockResponse struct {
	Version string `json:"version"`
	GenericResponse
}

type StateResponse struct {
	Version             string `json:"version"`
	ExecutionOptimistic bool   `json:"execution_optimistic"`
	Finalized           bool   `json:"finalized"`
	Data                any    `json:"data"`
}

type BlockHeaderResponse struct {
	Root      common.Root              `json:"root"`
	Canonical bool                     `json:"canonical"`
	Header    *SignedBeaconBlockHeader `json:"header"`
}

type BeaconBlockHeader struct {
	Slot          string `json:"slot"`
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

type SignedBeaconBlockHeader struct {
	Message   *BeaconBlockHeader `json:"message"`
	Signature string             `json:"signature"`
}

type GenesisData struct {
	GenesisTime           string      `json:"genesis_time"`
	GenesisValidatorsRoot common.Root `json:"genesis_validators_root"`
	GenesisForkVersion    string      `json:"genesis_fork_version"`
}

// GenesisResponse is handled with this explicit type since "finalized"
// and "execution_optimistic" are not part of the return value.
//
// https://ethereum.github.io/beacon-APIs/#/Beacon/getGenesis
type GenesisResponse struct {
	Data GenesisData `json:"data"`
}

type RootData struct {
	Root common.Root `json:"root"`
}

type ValidatorData struct {
	ValidatorBalanceData
	Status    string     `json:"status"`
	Validator *Validator `json:"validator"`
}

type ValidatorBalanceData struct {
	Index   uint64 `json:"index,string"`
	Balance uint64 `json:"balance,string"`
}

// Validator is the spec representation of the struct.
type Validator struct {
	PublicKey                  string `json:"pubkey"`
	WithdrawalCredentials      string `json:"withdrawal_credentials"`
	EffectiveBalance           string `json:"effective_balance"`
	Slashed                    bool   `json:"slashed"`
	ActivationEligibilityEpoch string `json:"activation_eligibility_epoch"`
	ActivationEpoch            string `json:"activation_epoch"`
	ExitEpoch                  string `json:"exit_epoch"`
	WithdrawableEpoch          string `json:"withdrawable_epoch"`
}

//nolint:staticcheck // todo: figure this out.
type CommitteeData struct {
	Index      uint64   `json:"index,string"`
	Slot       uint64   `json:"slot,string"`
	Validators []uint64 `json:"validators,string"`
}

type BlockRewardsData struct {
	ProposerIndex     uint64 `json:"proposer_index,string"`
	Total             uint64 `json:"total,string"`
	Attestations      uint64 `json:"attestations,string"`
	SyncAggregate     uint64 `json:"sync_aggregate,string"`
	ProposerSlashings uint64 `json:"proposer_slashings,string"`
	AttesterSlashings uint64 `json:"attester_slashings,string"`
}

type Sidecar struct {
	Index                       string                   `json:"index"`
	Blob                        string                   `json:"blob"`
	KZGCommitment               string                   `json:"kzg_commitment"`
	KZGProof                    string                   `json:"kzg_proof"`
	SignedBlockHeader           *SignedBeaconBlockHeader `json:"signed_block_header"`
	KZGCommitmentInclusionProof []string                 `json:"kzg_commitment_inclusion_proof"`
}

type SidecarsResponse struct {
	Data []*Sidecar `json:"data"`
}
