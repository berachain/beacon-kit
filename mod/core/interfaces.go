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

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, primitives.Slot, types.ReadOnlyBeaconBlock,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(primitives.Slot, *datypes.BlobSidecars) error
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

// BlobsProcessor is the interface for the blobs processor.
type BlobsProcessor interface {
	ProcessBlobs(
		primitives.Slot,
		AvailabilityStore,
		*datypes.BlobSidecars,
	) error
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor interface {
	ProcessRandao(
		state.BeaconState,
		types.BeaconBlock,
	) error
	ProcessRandaoMixesReset(
		state.BeaconState,
	) error
}
