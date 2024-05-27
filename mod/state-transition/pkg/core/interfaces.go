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

package core

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
)

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state types.
type BeaconState[ValidatorT any] interface {
	Copy() BeaconState[ValidatorT]
	Save()
	Context() context.Context
	HashTreeRoot() ([32]byte, error)
	ReadOnlyBeaconState[ValidatorT]
	WriteOnlyBeaconState[ValidatorT]
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState[ValidatorT any] interface {
	ReadOnlyEth1Data
	ReadOnlyRandaoMixes
	ReadOnlyStateRoots
	ReadOnlyValidators[ValidatorT]
	ReadOnlyWithdrawals

	GetBalance(math.ValidatorIndex) (math.Gwei, error)
	GetSlot() (math.Slot, error)
	GetGenesisValidatorsRoot() (primitives.Root, error)
	GetBlockRootAtIndex(uint64) (primitives.Root, error)
	GetLatestBlockHeader() (*types.BeaconBlockHeader, error)
	GetTotalActiveBalances(uint64) (math.Gwei, error)
	GetValidators() ([]ValidatorT, error)
	GetTotalSlashing() (math.Gwei, error)
	GetNextWithdrawalIndex() (uint64, error)
	GetNextWithdrawalValidatorIndex() (math.ValidatorIndex, error)
	GetTotalValidators() (uint64, error)
	GetValidatorsByEffectiveBalance() ([]ValidatorT, error)
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState[ValidatorT any] interface {
	WriteOnlyEth1Data
	WriteOnlyRandaoMixes
	WriteOnlyStateRoots
	WriteOnlyValidators[ValidatorT]

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

// Validator represents an interface for a validator with generic type
// ValidatorT.
type Validator[
	ValidatorT any,
	WithdrawalCredentialsT ~[32]byte,
] interface {
	ssz.Marshaler
	ssz.HashRoot
	// New creates a new validator with the given parameters.
	New(
		pubkey crypto.BLSPubkey,
		withdrawalCredentials WithdrawalCredentialsT,
		amount math.Gwei,
		effectiveBalanceIncrement math.Gwei,
		maxEffectiveBalance math.Gwei,
	) ValidatorT
	// IsSlashed returns true if the validator is slashed.
	IsSlashed() bool
	// GetPubkey returns the public key of the validator.
	GetPubkey() crypto.BLSPubkey
	// GetEffectiveBalance returns the effective balance of the validator in
	// Gwei.
	GetEffectiveBalance() math.Gwei
	// SetEffectiveBalance sets the effective balance of the validator in Gwei.
	SetEffectiveBalance(math.Gwei)
	// GetWithdrawableEpoch returns the epoch when the validator can withdraw.
	GetWithdrawableEpoch() math.Epoch
}
