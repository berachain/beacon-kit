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

package transition

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state types.
type BeaconState[
	T any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	ExecutionPayloadHeaderT,
	ForkT,
	KVStoreT any,
	ValidatorT Validator[ValidatorT, WithdrawalCredentialsT],
	ValidatorsT,
	WithdrawalT any,
	WithdrawalCredentialsT interface {
		~[32]byte
		ToExecutionAddress() (common.ExecutionAddress, error)
	},
] interface {
	NewFromDB(
		bdb KVStoreT,
		cs common.ChainSpec,
	) T
	Copy() T
	Save()
	Context() context.Context
	HashTreeRoot() common.Root
	ReadOnlyBeaconState[
		BeaconBlockHeaderT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT, WithdrawalT,
	]
	WriteOnlyBeaconState[
		BeaconBlockHeaderT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	]
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState[
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	ExecutionPayloadHeaderT, ForkT,
	ValidatorT, ValidatorsT, WithdrawalT any,
] interface {
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators[ValidatorT]

	GetBlockRootAtIndex(uint64) (common.Root, error)
	GetEth1DepositIndex() (uint64, error)
	GetFork() (ForkT, error)
	GetGenesisValidatorsRoot() (common.Root, error)
	GetLatestBlockHeader() (BeaconBlockHeaderT, error)
	GetLatestExecutionPayloadHeader() (
		ExecutionPayloadHeaderT, error,
	)
	GetSlot() (math.Slot, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	GetTotalValidators() (uint64, error)
	GetValidators() (ValidatorsT, error)
	GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState[
	BeaconBlockHeaderT, ExecutionPayloadHeaderT,
	ForkT, ValidatorT any,
] interface {
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators[ValidatorT]

	SetEth1DepositIndex(uint64) error
	SetFork(ForkT) error
	SetGenesisValidatorsRoot(root common.Root) error
	SetLatestBlockHeader(BeaconBlockHeaderT) error
	SetLatestExecutionPayloadHeader(
		ExecutionPayloadHeaderT,
	) error
	SetSlot(math.Slot) error
	UpdateBlockRootAtIndex(uint64, common.Root) error
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
type WriteOnlyValidators[ValidatorT any] interface {
	UpdateValidatorAtIndex(
		math.ValidatorIndex,
		ValidatorT,
	) error

	AddValidator(ValidatorT) error
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators[ValidatorT any] interface {
	ValidatorIndexByPubkey(
		crypto.BLSPubkey,
	) (math.ValidatorIndex, error)

	ValidatorByIndex(
		math.ValidatorIndex,
	) (ValidatorT, error)
}
