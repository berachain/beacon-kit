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

	"github.com/berachain/beacon-kit/mod/core"
	"github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/sourcegraph/conc/iter"
)

type BlobVerifier interface{}

// Processor is the processor for blobs.
type Processor struct {
	bv BlobVerifier
}

// NewProcessor creates a new processor.
func NewProcessor(bv BlobVerifier) *Processor {
	return &Processor{
		bv: bv,
	}
}

// ProcessBlob processes a blob.
func (sp *Processor) ProcessBlobs(
	slot primitives.Slot,
	avs core.AvailabilityStore,
	sidecars *datypes.BlobSidecars,
) error {
	// Ensure the blobs are available.
	if err := errors.Join(iter.Map(
		sidecars.Sidecars,
		func(sidecar **datypes.BlobSidecar) error {
			if *sidecar == nil {
				return ErrAttemptedToVerifyNilSidecar
			}
			// Store the blobs under a single height.
			return types.VerifyKZGInclusionProof(
				*sidecar,
			)
		},
	)...); err != nil {
		return err
	}

	return avs.Persist(slot, sidecars.Sidecars...)
}
