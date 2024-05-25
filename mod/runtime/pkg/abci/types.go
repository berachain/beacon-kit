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

package abci

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// BuilderService is responsible for building beacon blocks.
type BuilderService[
	BeaconBlockT types.BeaconBlock,
	BeaconStateT state.BeaconState,
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
}

// BlockchainService defines the interface for interacting with the blockchain
// state and processing blocks.
type BlockchainService[BlobsSidecarsT ssz.Marshallable] interface {
	// ProcessBlockAndBlobs processes the given beacon block and associated
	// blobs
	// sidecars.
	ProcessBlockAndBlobs(
		context.Context,
		types.BeaconBlock,
		BlobsSidecarsT,
	) ([]*transition.ValidatorUpdate, error)

	// VerifyPayloadOnBlk verifies the payload on the given beacon block.
	VerifyPayloadOnBlk(
		context.Context, types.BeaconBlock,
	) error
}
