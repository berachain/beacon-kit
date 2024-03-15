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
	randaotypes "github.com/berachain/beacon-kit/beacon/core/randao/types"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	ssz "github.com/prysmaticlabs/fastssz"
)

// BeaconBlock is the interface for a beacon block.
type BeaconBlock interface {
	ReadOnlyBeaconBlock
	WriteOnlyBeaconBlock
}

// WriteOnlyBeaconBlock is the interface for a write-only beacon block.
type WriteOnlyBeaconBlock interface {
}

// ReadOnlyBeaconBlock is the interface for a read-only beacon block.
type ReadOnlyBeaconBlock interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	Version() uint32
	GetSlot() primitives.Slot
	GetParentBlockRoot() [32]byte
	GetProposerIndex() primitives.ValidatorIndex
	GetBody() BeaconBlockBody
}

type BeaconBlockBody interface {
	WriteOnlyBeaconBlockBody
	ReadOnlyBeaconBlockBody
}

type WriteOnlyBeaconBlockBody interface {
	SetDeposits([]*Deposit)
	SetExecutionData(enginetypes.ExecutionPayload) error
}

type ReadOnlyBeaconBlockBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool

	// Execution returns the execution data of the block.
	GetDeposits() []*Deposit
	GetRandaoReveal() randaotypes.Reveal
	GetExecutionPayload() enginetypes.ExecutionPayload
	GetBlobKzgCommitments() [][48]byte
}
