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

package blob

import (
	"time"

	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Processor is the blob processor that handles the processing and verification
// of blob sidecars.
type Processor[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, *types.BlobSidecars,
	],
	BeaconBlockBodyT any,
] struct {
	// logger is used to log information and errors.
	logger log.Logger[any]
	// chainSpec defines the specifications of the blockchain.
	chainSpec primitives.ChainSpec
	// verifier is responsible for verifying the blobs.
	verifier *Verifier
	// blockBodyOffsetFn is a function that calculates the block body offset
	// based on the slot and chain specifications.
	blockBodyOffsetFn func(math.Slot, primitives.ChainSpec) uint64
	// metrics is used to collect and report processor metrics.
	metrics *processorMetrics
}

// NewProcessor creates a new blob processor.
func NewProcessor[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, *types.BlobSidecars,
	],
	BeaconBlockBodyT any,
](
	logger log.Logger[any],
	chainSpec primitives.ChainSpec,
	verifier *Verifier,
	blockBodyOffsetFn func(math.Slot, primitives.ChainSpec) uint64,
	telemetrySink TelemetrySink,
) *Processor[AvailabilityStoreT, BeaconBlockBodyT] {
	return &Processor[AvailabilityStoreT, BeaconBlockBodyT]{
		logger:            logger,
		chainSpec:         chainSpec,
		verifier:          verifier,
		blockBodyOffsetFn: blockBodyOffsetFn,
		metrics:           newProcessorMetrics(telemetrySink),
	}
}

// ProcessBlobs processes the blobs and ensures they match the local state.
func (sp *Processor[AvailabilityStoreT, BeaconBlockBodyT]) ProcessBlobs(
	slot math.Slot,
	avs AvailabilityStoreT,
	sidecars *types.BlobSidecars,
) error {
	var (
		numSidecars = math.U64(sidecars.Len())
		startTime   = time.Now()
	)

	defer sp.metrics.measureProcessBlobsDuration(startTime, numSidecars)

	// If there are no blobs to verify, return early.
	if numSidecars == 0 {
		sp.logger.Info(
			"no blob sidecars to verify, skipping verifier 🧢",
			"slot",
			slot,
		)
		return nil
	}

	// Otherwise, we run the verification checks on the blobs.
	if err := sp.verifier.VerifyBlobs(
		sidecars,
		sp.blockBodyOffsetFn(slot, sp.chainSpec),
	); err != nil {
		return err
	}

	sp.logger.Info(
		"successfully verified all blob sidecars 💦",
		"num-sidecars",
		numSidecars,
		"slot",
		slot,
	)

	// Lastly, we store the blobs in the availability store.
	return avs.Persist(slot, sidecars)
}
