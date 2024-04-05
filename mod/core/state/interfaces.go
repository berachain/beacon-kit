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

	"github.com/berachain/beacon-kit/mod/core/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state consensus.
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
	ReadOnlyDeposits
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators
	ReadOnlyWithdrawals

	GetSlot() (primitives.Slot, error)
	GetGenesisValidatorsRoot() (primitives.Root, error)
	GetBlockRootAtIndex(uint64) (primitives.Root, error)
	GetLatestBlockHeader() (*primitives.BeaconBlockHeader, error)
	GetTotalActiveBalances(uint64) (primitives.Gwei, error)
	GetValidators() ([]*types.Validator, error)
	GetEth1BlockHash() (primitives.ExecutionHash, error)
	GetTotalSlashing() (primitives.Gwei, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	WriteOnlyDeposits
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators
	WriteOnlyWithdrawals
	SetSlot(primitives.Slot) error
	UpdateBlockRootAtIndex(uint64, primitives.Root) error
	SetLatestBlockHeader(*primitives.BeaconBlockHeader) error
	IncreaseBalance(primitives.ValidatorIndex, primitives.Gwei) error
	DecreaseBalance(primitives.ValidatorIndex, primitives.Gwei) error
	UpdateEth1BlockHash(primitives.ExecutionHash) error
	UpdateSlashingAtIndex(uint64, primitives.Gwei) error
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
	// Add methods here
	UpdateValidatorAtIndex(
		primitives.ValidatorIndex,
		*types.Validator,
	) error

	AddValidator(*types.Validator) error
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators interface {
	ValidatorIndexByPubkey(
		primitives.BLSPubkey,
	) (primitives.ValidatorIndex, error)

	ValidatorByIndex(
		primitives.ValidatorIndex,
	) (*types.Validator, error)
}

// WriteOnlyDeposits has write access to deposit queue.
type WriteOnlyDeposits interface {
	SetEth1DepositIndex(uint64) error
	EnqueueDeposits([]*types.Deposit) error
	DequeueDeposits(uint64) ([]*types.Deposit, error)
}

// ReadOnlyDeposits has read access to deposit queue.
type ReadOnlyDeposits interface {
	GetEth1DepositIndex() (uint64, error)
	ExpectedDeposits(uint64) ([]*types.Deposit, error)
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals interface {
	ExpectedWithdrawals(uint64) ([]*primitives.Withdrawal, error)
}

// WriteOnlyWithdrawals only has write access to withdrawal methods.
type WriteOnlyWithdrawals interface {
	EnqueueWithdrawals([]*primitives.Withdrawal) error
	DequeueWithdrawals(uint64) ([]*primitives.Withdrawal, error)
}
