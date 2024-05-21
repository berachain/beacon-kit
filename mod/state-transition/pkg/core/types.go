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
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT any, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool

	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// Context defines an interface for managing state transition context.
type Context interface {
	context.Context

	// GetValidateResult returns whether to validate the result of the state
	// transition.
	GetValidateResult() bool

	// GetSkipPayloadIfExists returns whether to skip verifying the payload if
	// it already exists on the execution client.
	GetSkipPayloadIfExists() bool

	// GetOptimisticEngine returns whether to optimistically assume the
	// execution client has the correct state when certain errors are returned
	// by the execution engine.
	GetOptimisticEngine() bool

	// Unwrap returns the underlying golang standard library context.
	Unwrap() context.Context
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// VerifyAndNotifyNewPayload verifies the new payload and notifies the
	// execution client.
	VerifyAndNotifyNewPayload(
		ctx context.Context,
		req *engineprimitives.NewPayloadRequest[types.ExecutionPayload],
	) error
}

// RandaoProcessor is the interface for the randao processor.
type RandaoProcessor[BeaconBlockT, BeaconStateT any] interface {
	// ProcessRandao processes the RANDAO reveal and ensures it
	// matches the local state.
	ProcessRandao(
		BeaconStateT,
		BeaconBlockT,
	) error

	// ProcessRandaoMixesReset resets the RANDAO mixes as defined
	// in the Ethereum 2.0 specification.
	ProcessRandaoMixesReset(
		BeaconStateT,
	) error
}
