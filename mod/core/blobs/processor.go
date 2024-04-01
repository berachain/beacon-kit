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

package blobs

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/sourcegraph/conc/iter"
)

// Processor is the processor for blobs.
type Processor struct{}

// NewProcessor creates a new processor.
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessBlob processes a blob.
func (sp *Processor) ProcessBlobs(
	avs state.AvailabilityStore,
	blk types.BeaconBlock,
	sidecars *datypes.BlobSidecars,
) error {
	// Verify the KZG inclusion proofs.
	bodyRoot, err := blk.GetBody().HashTreeRoot()
	if err != nil {
		return err
	}

	// Ensure the blobs are available.
	if err = errors.Join(iter.Map(
		sidecars.Sidecars,
		func(sidecar **datypes.BlobSidecar) error {
			if *sidecar == nil {
				return ErrAttemptedToVerifyNilSidecar
			}
			// Store the blobs under a single height.
			return types.VerifyKZGInclusionProof(
				bodyRoot[:], *sidecar, (*sidecar).Index,
			)
		},
	)...); err != nil {
		return err
	}

	return avs.Persist(blk.GetSlot(), sidecars.Sidecars...)
}
