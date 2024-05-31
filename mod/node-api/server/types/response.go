// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
)

type ErrorResponse struct {
	Code    int `json:"code"`
	Message any `json:"message"`
}

type DataResponse struct {
	Data any `json:"data"`
}

type MessageResponse struct {
	Message any `json:"message"`
}

//nolint:lll
type GenesisData struct {
	GenesisTime           uint64             `json:"genesis_time,string"`
	GenesisValidatorsRoot primitives.Bytes32 `json:"genesis_validators_root,string"`
	GenesisForkVersion    primitives.Version `json:"genesis_fork_version,string"`
}

type RootData struct {
	Root primitives.Root `json:"root"`
}

type ValidatorResponse struct {
	ExecutionOptimistic bool `json:"execution_optimistic"`
	Finalized           bool `json:"finalized"`
	Data                any  `json:"data"`
}

type ValidatorData struct {
	Index     uint64           `json:"index,string"`
	Balance   uint64           `json:"balance,string"`
	Status    string           `json:"status"`
	Validator *types.Validator `json:"validator"`
}

type ValidatorBalanceData struct {
	Index   uint64 `json:"index,string"`
	Balance uint64 `json:"balance,string"`
}

type CommitteeData struct {
	Index      uint64   `json:"index,string"`
	Slot       uint64   `json:"slot,string"`
	Validators []uint64 `json:"validators,string"`
}

type SyncCommitteeData struct {
	Validators          []uint64   `json:"validators,string,string"`
	ValidatorAggregates [][]uint64 `json:"validator_aggregates,string,string"`
}

type BlockRewardsData struct {
	ProposerIndex     uint64 `json:"proposer_index,string"`
	Total             uint64 `json:"total,string"`
	Attestations      uint64 `json:"attestations,string"`
	SyncAggregate     uint64 `json:"sync_aggregate,string"`
	ProposerSlashings uint64 `json:"proposer_slashings,string"`
	AttesterSlashings uint64 `json:"attester_slashings,string"`
}

//nolint:lll
type SpecParamsResponse struct {
	DepositContractAddress          string `json:"DEPOSIT_CONTRACT_ADDRESS,string"`
	DepositNetworkID                uint64 `json:"DEPOSIT_NETWORK_ID,string"`
	DomainAggregateAndProof         string `json:"DOMAIN_AGGREGATE_AND_PROOF,string"`
	InactivityPenaltyQuotient       uint64 `json:"INACTIVITY_PENALTY_QUOTIENT,string"`
	InactivityPenaltyQuotientAltair uint64 `json:"INACTIVITY_PENALTY_QUOTIENT_ALTAIR,string"`
}

type VoluntaryExitData struct {
	Epoch          uint64 `json:"epoch,string"`
	ValidatorIndex uint64 `json:"validator_index,string"`
}

type MessageSignature struct {
	Message   any                 `json:"message"`
	Signature crypto.BLSSignature `json:"signature"`
}

type BtsToExecutionChangeData struct {
	ValidatorIndex     uint64                  `json:"validator_index,string"`
	FromBlsPubkey      crypto.BLSPubkey        `json:"from_bls_pubkey,string"`
	ToExecutionAddress common.ExecutionAddress `json:"to_execution_address,string"`
}

type ProposerDutiesData struct {
	Pubkey         crypto.BLSPubkey `json:"pubkey,string"`
	ValidatorIndex uint64           `json:"validator_index,string"`
	Slot           uint64           `json:"slot,string"`
}

type BlockProposerDutiesResponse struct {
	DependentRoot       primitives.Root       `json:"dependent_root,string"`
	ExecutionOptimistic bool                  `json:"execution_optimistic"`
	Data                []*ProposerDutiesData `json:"data"`
}

type BlockHeaderData struct {
	Root      primitives.Root     `json:"root"`
	Canonical bool                `json:"canonical"`
	Header    MessageResponse     `json:"header"`
	Signature crypto.BLSSignature `json:"signature"`
}

type BlockResponse struct {
	Version             string            `json:"version"`
	ExecutionOptimistic bool              `json:"execution_optimistic"`
	Finalized           bool              `json:"finalized"`
	Data                *MessageSignature `json:"data"`
}
