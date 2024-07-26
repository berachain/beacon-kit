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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlockHeader is the interface for a beacon block header.
type BeaconBlockHeader[BeaconBlockHeaderT any] interface {
	constraints.SSZRootable
	New(
		slot math.Slot,
		proposerIndex math.ValidatorIndex,
		parentBlockRoot common.Root,
		stateRoot common.Root,
		bodyRoot common.Root,
	) BeaconBlockHeaderT
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	SetStateRoot(stateRoot common.Root)
	GetBodyRoot() common.Root
}

// BeaconState is the interface for a beacon state.
type BeaconState[
	BeaconBlockHeaderT any,
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ValidatorT any,
] interface {
	constraints.SSZRootable
	// GetLatestExecutionPayloadHeader retrieves the latest execution payload
	// header.
	GetLatestExecutionPayloadHeader() (ExecutionPayloadHeaderT, error)
	// GetEth1DepositIndex retrieves the eth1 deposit index.
	GetEth1DepositIndex() (uint64, error)
	// GetBalance retrieves the balance of a validator.
	GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
	// GetSlot retrieves the current slot.
	GetSlot() (math.Slot, error)
	// GetFork retrieves the fork.
	GetFork() (ForkT, error)
	// GetGenesisValidatorsRoot retrieves the genesis validators root.
	GetGenesisValidatorsRoot() (common.Root, error)
	// GetLatestBlockHeader retrieves the latest block header.
	GetLatestBlockHeader() (BeaconBlockHeaderT, error)
	// GetBlockRootAtIndex retrieves the block root at the given index.
	GetBlockRootAtIndex(index uint64) (common.Root, error)
	// GetEth1Data retrieves the eth1 data.
	GetEth1Data() (Eth1DataT, error)
	// GetValidators retrieves all validators.
	GetValidators() ([]ValidatorT, error)
	// GetBalances retrieves all balances.
	GetBalances() ([]uint64, error)
	// GetNextWithdrawalIndex retrieves the next withdrawal index.
	GetNextWithdrawalIndex() (uint64, error)
	// GetNextWithdrawalValidatorIndex retrieves the next withdrawal validator
	// index.
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	// GetTotalSlashing retrieves the total slashing.
	GetTotalSlashing() (math.Gwei, error)
	// GetRandaoMixAtIndex retrieves the randao mix at the given index.
	GetRandaoMixAtIndex(index uint64) (common.Bytes32, error)
	// GetSlashings retrieves all slashings.
	GetSlashings() ([]uint64, error)
	// GetSlashingAtIndex retrieves the slashing at the given index.
	GetSlashingAtIndex(index uint64) (math.Gwei, error)
	// GetTotalValidators retrieves the total validators.
	GetTotalValidators() (uint64, error)
	// GetTotalActiveBalances retrieves the total active balances.
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	// ValidatorByIndex retrieves the validator at the given index.
	ValidatorByIndex(index math.ValidatorIndex) (ValidatorT, error)
	// ValidatorIndexByPubkey retrieves the validator index by the given pubkey.
	ValidatorIndexByPubkey(pubkey crypto.BLSPubkey) (math.ValidatorIndex, error)
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
	// GetValidatorsByEffectiveBalance retrieves validators by effective
	// balance.
	GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
	// StateRootAtIndex retrieves the state root at the given index.
	StateRootAtIndex(index uint64) (common.Root, error)
}
