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

package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

type Backend struct {
	getNewStateDB func(context.Context, string) StateDB
}

// New creates and returns a new Backend instance.
// TODO: need to add state_id resolver; possible values are: "head" (canonical
// head in node's view), "genesis", "finalized", "justified", <slot>, <hex
// encoded stateRoot with 0x prefix>.
func New(
	getNewStateDB func(ctx context.Context, stateId string) StateDB,
) *Backend {
	return &Backend{
		getNewStateDB: getNewStateDB,
	}
}

type StateDB interface {
	GetGenesisValidatorsRoot() (common.Root, error)
	GetSlot() (math.Slot, error)
	GetLatestExecutionPayloadHeader() (
		*types.ExecutionPayloadHeader, error,
	)
	SetLatestExecutionPayloadHeader(
		payloadHeader *types.ExecutionPayloadHeader,
	) error
	GetEth1DepositIndex() (uint64, error)
	SetEth1DepositIndex(
		index uint64,
	) error
	GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
	SetBalance(idx math.ValidatorIndex, balance math.Gwei) error
	SetSlot(slot math.Slot) error
	GetFork() (*types.Fork, error)
	SetFork(fork *types.Fork) error
	GetLatestBlockHeader() (*types.BeaconBlockHeader, error)
	SetLatestBlockHeader(header *types.BeaconBlockHeader) error
	GetBlockRootAtIndex(index uint64) (common.Root, error)
	StateRootAtIndex(index uint64) (common.Root, error)
	GetEth1Data() (*types.Eth1Data, error)
	SetEth1Data(data *types.Eth1Data) error
	GetValidators() ([]*types.Validator, error)
	GetBalances() ([]uint64, error)
	GetNextWithdrawalIndex() (uint64, error)
	SetNextWithdrawalIndex(index uint64) error
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	SetNextWithdrawalValidatorIndex(index math.ValidatorIndex) error
	GetTotalSlashing() (math.Gwei, error)
	SetTotalSlashing(total math.Gwei) error
	GetRandaoMixAtIndex(index uint64) (common.Bytes32, error)
	GetSlashings() ([]uint64, error)
	SetSlashingAtIndex(index uint64, amount math.Gwei) error
	GetSlashingAtIndex(index uint64) (math.Gwei, error)
	GetTotalValidators() (uint64, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	ValidatorByIndex(index math.ValidatorIndex) (*types.Validator, error)
	UpdateBlockRootAtIndex(index uint64, root common.Root) error
	UpdateStateRootAtIndex(index uint64, root common.Root) error
	UpdateRandaoMixAtIndex(index uint64, mix common.Bytes32) error
	UpdateValidatorAtIndex(
		index math.ValidatorIndex,
		validator *types.Validator,
	) error
	ValidatorIndexByPubkey(pubkey crypto.BLSPubkey) (math.ValidatorIndex, error)
	AddValidator(
		val *types.Validator,
	) error
	GetValidatorsByEffectiveBalance() ([]*types.Validator, error)
}
