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

package da

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// BlobProcessor is the blob processor.
type BlobProcessor[
	BeaconBlockBodyT any,
	BlobSidecarsT interface{ Len() int },
] struct {
	logger    log.Logger[any]
	chainSpec primitives.ChainSpec
	bv        BlobVerifier[BlobSidecarsT]
}

// NewBlobProcessor creates a new blob processor.
func NewBlobProcessor[
	BeaconBlockBodyT any, BlobSidecarsT interface{ Len() int },
](
	logger log.Logger[any],
	chainSpec primitives.ChainSpec,
	bv BlobVerifier[BlobSidecarsT],
) *BlobProcessor[BeaconBlockBodyT, BlobSidecarsT] {
	return &BlobProcessor[BeaconBlockBodyT, BlobSidecarsT]{
		logger:    logger,
		chainSpec: chainSpec,
		bv:        bv,
	}
}

// ProcessBlobs processes the blobs and ensures they match the local state.
func (sp *BlobProcessor[BeaconBlockBodyT, BlobSidecarsT]) ProcessBlobs(
	slot math.Slot,
	// TODO: decouple from core.
	avs core.AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	sidecars BlobSidecarsT,
) error {
	// If there are no blobs to verify, return early.
	numSidecars := sidecars.Len()
	if numSidecars == 0 {
		sp.logger.Info(
			"no blob sidecars to verify, skipping verifier 🧢",
			"slot",
			slot,
		)
		return nil
	}

	// Otherwise, we run the verification checks on the blobs.
	if err := sp.bv.VerifyBlobs(
		sidecars,
		types.BlockBodyKZGOffset(slot, sp.chainSpec),
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
