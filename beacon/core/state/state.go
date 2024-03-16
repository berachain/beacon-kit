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
	"github.com/berachain/beacon-kit/beacon/core/randao/types"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state consensus.
type BeaconState interface {
	ReadOnlyBeaconState
	WriteOnlyBeaconState
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState interface {
	ReadOnlyDeposits
	ReadOnlyRandaoMixes
	ReadOnlyValidators
	ReadOnlyWithdrawals

	GetSlot() primitives.Slot
	SetParentBlockRoot([32]byte)
	GetParentBlockRoot() [32]byte
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	WriteOnlyDeposits
	WriteOnlyRandaoMixes
	WriteOnlyValidators
	WriteOnlyWithdrawals

	SetParentBlockRoot([32]byte)
}

// WriteOnlyRandaoMixes defines a struct which only has write access to randao
// mixes methods.
type WriteOnlyRandaoMixes interface {
	SetRandaoMix(types.Mix) error
}

// ReadOnlyRandaoMixes defines a struct which only has read access to randao
// mixes methods.
type ReadOnlyRandaoMixes interface {
	RandaoMix() (types.Mix, error)
}

// WriteOnlyValidators has write access to validator methods.
type WriteOnlyValidators interface {
	// Add methods here
}

// ReadOnlyValidators has read access to validator methods.
type ReadOnlyValidators interface {
	ValidatorIndexByPubkey(
		[]byte,
	) (primitives.ValidatorIndex, error)

	ValidatorPubKeyByIndex(
		primitives.ValidatorIndex,
	) ([]byte, error)
}

// ReadWriteValidators has read and write access to validator methods.
type ReadWriteDeposits interface {
	ReadOnlyDeposits
	WriteOnlyDeposits
}

// ReadWriteDepositQueue has read and write access to deposit queue.
type WriteOnlyDeposits interface {
	EnqueueDeposits([]*beacontypes.Deposit) error
	DequeueDeposits(uint64) ([]*beacontypes.Deposit, error)
}

// ReadOnlyDeposits has read access to deposit queue.
type ReadOnlyDeposits interface {
	ExpectedDeposits(uint64) ([]*beacontypes.Deposit, error)
}

// ReadWriteWithdrawals has read and write access to withdrawal methods.
type ReadWriteWithdrawals interface {
	ReadOnlyWithdrawals
	WriteOnlyWithdrawals
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals interface {
	ExpectedWithdrawals(uint64) ([]*enginetypes.Withdrawal, error)
}

// WriteOnlyWithdrawals only has write access to withdrawal methods.
type WriteOnlyWithdrawals interface {
	EnqueueWithdrawals([]*enginetypes.Withdrawal) error
	DequeueWithdrawals(uint64) ([]*enginetypes.Withdrawal, error)
}
