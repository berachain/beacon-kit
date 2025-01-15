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

package core

import (
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
)

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState interface {
	ReadOnlyEth1Data
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators
	ReadOnlyWithdrawals

	GetBalance(math.ValidatorIndex) (math.Gwei, error)
	GetSlot() (math.Slot, error)
	GetFork() (*ctypes.Fork, error)
	GetGenesisValidatorsRoot() (common.Root, error)
	GetBlockRootAtIndex(uint64) (common.Root, error)
	GetLatestBlockHeader() (*ctypes.BeaconBlockHeader, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	GetValidators() (ctypes.Validators, error)
	GetSlashingAtIndex(uint64) (math.Gwei, error)
	GetTotalSlashing() (math.Gwei, error)
	GetNextWithdrawalIndex() (uint64, error)
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	GetTotalValidators() (uint64, error)
	GetValidatorsByEffectiveBalance() ([]*ctypes.Validator, error)
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	WriteOnlyEth1Data
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators

	SetGenesisValidatorsRoot(root common.Root) error
	SetFork(*ctypes.Fork) error
	SetSlot(math.Slot) error
	UpdateBlockRootAtIndex(uint64, common.Root) error
	SetLatestBlockHeader(*ctypes.BeaconBlockHeader) error
	IncreaseBalance(math.ValidatorIndex, math.Gwei) error
	DecreaseBalance(math.ValidatorIndex, math.Gwei) error
	UpdateSlashingAtIndex(uint64, math.Gwei) error
	SetNextWithdrawalIndex(uint64) error
	SetNextWithdrawalValidatorIndex(math.ValidatorIndex) error
	SetTotalSlashing(math.Gwei) error
}

// WriteOnlyStateRoots defines a struct which only has write access to state
// roots methods.
type WriteOnlyStateRoots interface {
	UpdateStateRootAtIndex(uint64, common.Root) error
}

// ReadOnlyStateRoots defines a struct which only has read access to state roots
// methods.
type ReadOnlyStateRoots interface {
	StateRootAtIndex(uint64) (common.Root, error)
}

// WriteOnlyRandaoMixes defines a struct which only has write access to randao
// mixes methods.
type WriteOnlyRandaoMixes interface {
	UpdateRandaoMixAtIndex(uint64, common.Bytes32) error
}

// ReadOnlyRandaoMixes defines a struct which only has read access to randao
// mixes methods.
type ReadOnlyRandaoMixes interface {
	GetRandaoMixAtIndex(uint64) (common.Bytes32, error)
}

// WriteOnlyValidators has write access to validator methods.
type WriteOnlyValidators interface {
	UpdateValidatorAtIndex(
		math.ValidatorIndex,
		*ctypes.Validator,
	) error

	AddValidator(*ctypes.Validator) error
	AddValidatorBartio(*ctypes.Validator) error
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators interface {
	ValidatorIndexByPubkey(
		crypto.BLSPubkey,
	) (math.ValidatorIndex, error)

	ValidatorByIndex(
		math.ValidatorIndex,
	) (*ctypes.Validator, error)
}

// WriteOnlyEth1Data has write access to eth1 data.
type WriteOnlyEth1Data interface {
	SetEth1Data(*ctypes.Eth1Data) error
	SetEth1DepositIndex(uint64) error
	SetLatestExecutionPayloadHeader(*ctypes.ExecutionPayloadHeader) error
}

// ReadOnlyEth1Data has read access to eth1 data.
type ReadOnlyEth1Data interface {
	GetEth1Data() (*ctypes.Eth1Data, error)
	GetEth1DepositIndex() (uint64, error)
	GetLatestExecutionPayloadHeader() (*ctypes.ExecutionPayloadHeader, error)
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals interface {
	EVMInflationWithdrawal() *engineprimitives.Withdrawal
	ExpectedWithdrawals() (engineprimitives.Withdrawals, error)
}
