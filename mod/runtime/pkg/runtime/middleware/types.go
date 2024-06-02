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

package middleware

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
)

// BlockchainService defines the interface for interacting with the blockchain
// state and processing blocks.
type BlockchainService[
	BeaconBlockT any, BlobsSidecarsT ssz.Marshallable,
] interface {
	// ProcessGenesisData processes the genesis data and initializes the beacon
	// state.
	ProcessGenesisData(
		context.Context,
		*genesis.Genesis[
			*types.Deposit, *types.ExecutionPayloadHeaderDeneb,
		],
	) ([]*transition.ValidatorUpdate, error)
	// ProcessBlockAndBlobs processes the given beacon block and associated
	// blobs sidecars.
	ProcessBlockAndBlobs(
		context.Context,
		BeaconBlockT,
		BlobsSidecarsT,
		bool,
	) ([]*transition.ValidatorUpdate, error)
}

// ValidatorService is responsible for building beacon blocks.
type ValidatorService[
	BeaconBlockT any,
	BeaconStateT any,
	BlobsSidecarsT ssz.Marshallable,
] interface {
	// RequestBestBlock requests the best beacon block for a given slot.
	// It returns the beacon block, associated blobs sidecars, and an error if
	// any.
	RequestBestBlock(
		context.Context, // The context for the request.
		math.Slot, // The slot for which the best block is requested.
	) (
		BeaconBlockT, BlobsSidecarsT, error,
	)
	// VerifyIncomingBlock verifies the incoming block and returns an error if
	// the block is invalid.
	VerifyIncomingBlock(
		ctx context.Context,
		blk BeaconBlockT,
	) error
}

// TelemetrySink is an interface for sending metrics to a telemetry backend.
type TelemetrySink interface {
	// MeasureSince measures the time since the given time.
	MeasureSince(key string, start time.Time, args ...string)
}

// StorageBackend is an interface for accessing the storage backend.
type StorageBackend[BeaconStateT any] interface {
	StateFromContext(ctx context.Context) BeaconStateT
}

// BeaconState is an interface for accessing the beacon state.
type BeaconState interface {
	ValidatorIndexByPubkey(
		pubkey crypto.BLSPubkey,
	) (math.ValidatorIndex, error)

	GetBlockRootAtIndex(
		index uint64,
	) (primitives.Root, error)

	ValidatorIndexByCometBFTAddress(
		cometBFTAddress []byte,
	) (math.ValidatorIndex, error)
}
