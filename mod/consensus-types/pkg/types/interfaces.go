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
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	eip4844 "github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
)

// BeaconBlockBody is the interface for a beacon block body.
type BeaconBlockBody interface {
	WriteOnlyBeaconBlockBody
	ReadOnlyBeaconBlockBody
	Length() uint64
}

// WriteOnlyBeaconBlockBody is the interface for a write-only beacon block body.
type WriteOnlyBeaconBlockBody interface {
	SetDeposits([]*Deposit)
	SetEth1Data(*Eth1Data)
	SetExecutionData(engineprimitives.ExecutionPayload) error
	SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash])
	SetRandaoReveal(crypto.BLSSignature)
}

// ReadOnlyBeaconBlockBody is the interface for
// a read-only beacon block body.
type ReadOnlyBeaconBlockBody interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool

	// Execution returns the execution data of the block.
	GetDeposits() []*Deposit
	GetEth1Data() *Eth1Data
	GetGraffiti() bytes.B32
	GetRandaoReveal() crypto.BLSSignature
	GetExecutionPayload() engineprimitives.ExecutionPayload
	GetBlobKzgCommitments() eip4844.KZGCommitments[common.ExecutionHash]
	GetTopLevelRoots() ([][32]byte, error)
}

// BeaconBlock is the interface for a beacon block.
type BeaconBlock interface {
	SetStateRoot(common.Root)
	ReadOnlyBeaconBlock[BeaconBlockBody]
}

type BeaconBlockG[BodyT any] struct {
	ReadOnlyBeaconBlock[BodyT]
}

// ReadOnlyBeaconBlock is the interface for a read-only beacon block.
type ReadOnlyBeaconBlock[BodyT any] interface {
	ssz.Marshaler
	ssz.Unmarshaler
	ssz.HashRoot
	IsNil() bool
	Version() uint32
	GetSlot() math.Slot
	GetProposerIndex() math.ValidatorIndex
	GetParentBlockRoot() common.Root
	GetStateRoot() common.Root
	GetBody() BodyT
	GetHeader() *BeaconBlockHeader
}
