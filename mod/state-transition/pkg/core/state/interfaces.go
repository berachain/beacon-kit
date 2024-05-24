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

package state

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state types.
type BeaconState interface {
	Copy() BeaconState
	Save()
	Context() context.Context
	HashTreeRoot() ([32]byte, error)
	ReadOnlyBeaconState
	WriteOnlyBeaconState
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState interface {
	ReadOnlyEth1Data
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators
	ReadOnlyWithdrawals

	GetBalance(math.ValidatorIndex) (math.Gwei, error)
	GetSlot() (math.Slot, error)
	GetGenesisValidatorsRoot() (primitives.Root, error)
	GetBlockRootAtIndex(uint64) (primitives.Root, error)
	GetLatestBlockHeader() (*types.BeaconBlockHeader, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	GetValidators() ([]*types.Validator, error)
	GetTotalSlashing() (math.Gwei, error)
	GetNextWithdrawalIndex() (uint64, error)
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	GetTotalValidators() (uint64, error)
	GetValidatorsByEffectiveBalance() ([]*types.Validator, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	WriteOnlyEth1Data
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators

	SetGenesisValidatorsRoot(root primitives.Root) error
	SetFork(*types.Fork) error
	SetSlot(math.Slot) error
	UpdateBlockRootAtIndex(uint64, primitives.Root) error
	SetLatestBlockHeader(*types.BeaconBlockHeader) error
	IncreaseBalance(math.ValidatorIndex, math.Gwei) error
	DecreaseBalance(math.ValidatorIndex, math.Gwei) error
	UpdateSlashingAtIndex(uint64, math.Gwei) error
	SetNextWithdrawalIndex(uint64) error
	SetNextWithdrawalValidatorIndex(math.ValidatorIndex) error
	RemoveValidatorAtIndex(math.ValidatorIndex) error
	SetTotalSlashing(math.Gwei) error
}

// WriteOnlyStateRoots defines a struct which only has write access to state
// roots methods.
type WriteOnlyStateRoots interface {
	UpdateStateRootAtIndex(uint64, primitives.Root) error
}

// ReadOnlyStateRoots defines a struct which only has read access to state roots
// methods.
type ReadOnlyStateRoots interface {
	StateRootAtIndex(uint64) (primitives.Root, error)
}

// WriteOnlyRandaoMixes defines a struct which only has write access to randao
// mixes methods.
type WriteOnlyRandaoMixes interface {
	UpdateRandaoMixAtIndex(uint64, primitives.Bytes32) error
}

// ReadOnlyRandaoMixes defines a struct which only has read access to randao
// mixes methods.
type ReadOnlyRandaoMixes interface {
	GetRandaoMixAtIndex(uint64) (primitives.Bytes32, error)
}

// WriteOnlyValidators has write access to validator methods.
type WriteOnlyValidators interface {
	UpdateValidatorAtIndex(
		math.ValidatorIndex,
		*types.Validator,
	) error

	AddValidator(*types.Validator) error
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators interface {
	ValidatorIndexByPubkey(
		crypto.BLSPubkey,
	) (math.ValidatorIndex, error)

	ValidatorByIndex(
		math.ValidatorIndex,
	) (*types.Validator, error)
}

// WriteOnlyEth1Data has write access to eth1 data.
type WriteOnlyEth1Data interface {
	SetEth1Data(*types.Eth1Data) error
	SetEth1DepositIndex(uint64) error
	SetLatestExecutionPayloadHeader(
		engineprimitives.ExecutionPayloadHeader,
	) error
}

// ReadOnlyEth1Data has read access to eth1 data.
type ReadOnlyEth1Data interface {
	GetEth1Data() (*types.Eth1Data, error)
	GetEth1DepositIndex() (uint64, error)
	GetLatestExecutionPayloadHeader() (
		engineprimitives.ExecutionPayloadHeader, error,
	)
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals interface {
	ExpectedWithdrawals() ([]*engineprimitives.Withdrawal, error)
}

type KVStore[KVStoreT any] interface {
	Context() context.Context
	WithContext(
		ctx context.Context,
	) KVStoreT
	Save()
	GetLatestExecutionPayloadHeader() (
		engineprimitives.ExecutionPayloadHeader, error,
	)
	SetLatestExecutionPayloadHeader(
		payloadHeader engineprimitives.ExecutionPayloadHeader,
	) error
	GetEth1DepositIndex() (uint64, error)
	SetEth1DepositIndex(
		index uint64,
	) error
	GetBalance(idx math.ValidatorIndex) (math.Gwei, error)
	SetBalance(idx math.ValidatorIndex, balance math.Gwei) error
	Copy() KVStoreT
	GetSlot() (math.Slot, error)
	SetSlot(slot math.Slot) error
	GetFork() (*types.Fork, error)
	SetFork(fork *types.Fork) error
	GetGenesisValidatorsRoot() (common.Root, error)
	SetGenesisValidatorsRoot(root common.Root) error
	GetLatestBlockHeader() (*types.BeaconBlockHeader, error)
	SetLatestBlockHeader(header *types.BeaconBlockHeader) error
	GetBlockRootAtIndex(index uint64) (primitives.Root, error)
	StateRootAtIndex(index uint64) (primitives.Root, error)
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
	GetRandaoMixAtIndex(index uint64) (primitives.Bytes32, error)
	GetSlashings() ([]uint64, error)
	SetSlashingAtIndex(index uint64, amount math.Gwei) error
	GetSlashingAtIndex(index uint64) (math.Gwei, error)
	GetTotalValidators() (uint64, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	ValidatorByIndex(index math.ValidatorIndex) (*types.Validator, error)
	UpdateBlockRootAtIndex(index uint64, root primitives.Root) error
	UpdateStateRootAtIndex(index uint64, root primitives.Root) error
	UpdateRandaoMixAtIndex(index uint64, mix primitives.Bytes32) error
	UpdateValidatorAtIndex(
		index math.ValidatorIndex,
		validator *types.Validator,
	) error
	ValidatorIndexByPubkey(pubkey crypto.BLSPubkey) (math.ValidatorIndex, error)
	AddValidator(
		val *types.Validator,
	) error
	GetValidatorsByEffectiveBalance() ([]*types.Validator, error)
	RemoveValidatorAtIndex(idx math.ValidatorIndex) error
}
