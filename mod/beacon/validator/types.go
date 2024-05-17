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

package validator

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	ssz "github.com/ferranbt/fastssz"
)

// BeaconState defines the interface for accessing various components of the
// beacon state.
type BeaconState interface {
	// GetSlot retrieves the current slot of the beacon state.
	GetSlot() (math.Slot, error)

	// GetBlockRootAtIndex fetches the block root at a specified index.
	GetBlockRootAtIndex(uint64) (primitives.Root, error)

	// GetLatestExecutionPayloadHeader returns the most recent execution payload
	// header.
	GetLatestExecutionPayloadHeader() (
		engineprimitives.ExecutionPayloadHeader,
		error,
	)

	// ValidatorIndexByPubkey finds the index of a validator based on their
	// public key.
	ValidatorIndexByPubkey(crypto.BLSPubkey) (math.ValidatorIndex, error)

	HashTreeRoot() ([32]byte, error)
}

type BeaconStorageBackend[BeaconStateT BeaconState] interface {
	StateFromContext(context.Context) BeaconStateT
}

// BlobFactory is the interface for building blobs.
type BlobFactory[
	BlobSidecarsT BlobSidecars,
	BeaconBlockBodyT types.ReadOnlyBeaconBlockBody,
] interface {
	// BuildSidecars generates sidecars for a given block and blobs bundle.
	BuildSidecars(
		blk types.ReadOnlyBeaconBlock[BeaconBlockBodyT],
		blobs engineprimitives.BlobsBundle,
	) (BlobSidecarsT, error)
}

// BlobSidecars is the interface for blobs sidecars.
type BlobSidecars interface {
	ssz.Marshaler
	ssz.Unmarshaler
	Len() int
}

// DepositStore defines the interface for deposit storage.
type DepositStore interface {
	// ExpectedDeposits returns `numView` expected deposits.
	ExpectedDeposits(
		numView uint64,
	) ([]*types.Deposit, error)
}

// RandaoProcessor defines the interface for processing RANDAO reveals.
type RandaoProcessor[
	BeaconStateT BeaconState,
] interface {
	// BuildReveal generates a RANDAO reveal based on the given beacon state.
	// It returns a Reveal object and any error encountered during the process.
	BuildReveal(st BeaconStateT) (crypto.BLSSignature, error)
}

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder[BeaconStateT BeaconState] interface {
	// RetrieveOrBuildPayload retrieves or builds the payload for the given
	// slot.
	RetrieveOrBuildPayload(
		ctx context.Context,
		st BeaconStateT,
		slot math.Slot,
		parentBlockRoot primitives.Root,
		parentEth1Hash common.ExecutionHash,
	) (engineprimitives.BuiltExecutionPayloadEnv, error)
}

// StateProcessor defines the interface for processing the state.
type StateProcessor[
	BeaconStateT BeaconState,
] interface {
	// BuildReveal generates a RANDAO reveal based on the given beacon state.
	// It returns a Reveal object and any error encountered during the process.
	ProcessBlock(
		st BeaconStateT,
		blk types.BeaconBlock,
		validateResult bool,
	) error

	// ProcessSlot processes the slot.
	ProcessSlot(
		st BeaconStateT,
	) error
}
