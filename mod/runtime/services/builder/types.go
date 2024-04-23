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

package builder

import (
	"context"

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/da"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
)

// BlobFactory is the interface for building blobs.
type BlobFactory[BBB da.BeaconBlockBody] interface {
	BuildSidecars(
		blk da.BeaconBlock[BBB],
		blobs engineprimitives.BlobsBundle,
	) (*datypes.BlobSidecars, error)
}

// BLSSigner defines an interface for cryptographic signing operations.
// It uses generic type parameters Signature and Pubkey, both of which are
// slices of bytes.
type BLSSigner interface {
	// PublicKey returns the public key of the signer.
	PublicKey() primitives.BLSPubkey

	// Sign takes a message as a slice of bytes and returns a signature as a
	// slice of bytes and an error.
	Sign([]byte) (primitives.BLSSignature, error)
}

// RandaoProcessor defines the interface for processing RANDAO reveals.
type RandaoProcessor interface {
	// BuildReveal generates a RANDAO reveal based on the given beacon state.
	// It returns a Reveal object and any error encountered during the process.
	BuildReveal(st state.BeaconState) (primitives.BLSSignature, error)
}

// PayloadBuilder represents a service that is responsible for
// building eth1 blocks.
type PayloadBuilder interface {
	RetrieveBuiltPayload(
		ctx context.Context,
		st state.BeaconState,
		slot math.Slot,
		parentBlockRoot primitives.Root,
		parentEth1Hash primitives.ExecutionHash,
	) (engineprimitives.BuiltExecutionPayloadEnv, error)
}
