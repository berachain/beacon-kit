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
	"github.com/ethereum/go-ethereum/common"
	beacontypesv1 "github.com/itsdevbear/bolaris/beacon/core/types/v1"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
)

type ForkchoiceStore interface {
	WriteOnlyForkChoice
	ReadOnlyForkChoice

	SetLastSeenBeaconBlock([32]byte)
	GetLastSeenBeaconBlock() [32]byte
}

// Write Only Fork Choice.
type WriteOnlyForkChoice interface {
	// Write Only Fork Choice.
	SetSafeEth1BlockHash(safeBlockHash common.Hash)
	SetFinalizedEth1BlockHash(finalizedBlockHash common.Hash)
	SetGenesisEth1Hash(genesisEth1Hash common.Hash)
}

// ReadOnlyForkChoice.
type ReadOnlyForkChoice interface {
	GetSafeEth1BlockHash() common.Hash
	GetFinalizedEth1BlockHash() common.Hash
	GenesisEth1Hash() common.Hash
}

// BeaconState is the interface for the beacon state. It
// is a combination of the read-only and write-only beacon state consensus.
type BeaconState interface {
	ReadOnlyBeaconState
	WriteOnlyBeaconState
}

// ReadOnlyBeaconState is the interface for a read-only beacon state.
type ReadOnlyBeaconState interface {
	// TODO: fill these in as we develop impl
	ReadWriteDepositQueue
	ReadOnlyWithdrawals

	GetParentBlockRoot() [32]byte

	// TODO: Actually decouple epocha nd slot
	// GetEpochBySlot(primitives.Slot) primitives.Epoch
}

// WriteOnlyBeaconState is the interface for a write-only beacon state.
type WriteOnlyBeaconState interface {
	SetParentBlockRoot([32]byte)
}

// ReadWriteDepositQueue has read and write access to deposit queue.
type ReadWriteDepositQueue interface {
	EnqueueDeposits([]*beacontypesv1.Deposit) error
	DequeueDeposits(n uint64) ([]*beacontypesv1.Deposit, error)
}

// ReadOnlyWithdrawals only has read access to withdrawal methods.
type ReadOnlyWithdrawals interface {
	ExpectedWithdrawals() ([]*enginev1.Withdrawal, error)
}
