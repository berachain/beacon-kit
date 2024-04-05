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

package types

import (
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives/engine"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"golang.org/x/sync/errgroup"
)

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	blk BeaconBlock,
	blobs *engine.BlobsBundleV1,
) (*datypes.BlobSidecars, error) {
	numBlobs := uint64(len(blobs.Blobs))
	sidecars := make([]*datypes.BlobSidecar, numBlobs)
	g := errgroup.Group{}

	blkHeader := blk.GetHeader()
	for i := uint64(0); i < numBlobs; i++ {
		i := i // capture range variable
		g.Go(func() error {
			// Create Inclusion Proof
			inclusionProof, err := MerkleProofKZGCommitment(blk, i)
			if err != nil {
				return err
			}

			blob := kzg.Blob(blobs.Blobs[i])
			sidecars[i] = &datypes.BlobSidecar{
				Index:             i,
				Blob:              &blob,
				KzgCommitment:     kzg.Commitment(blobs.Commitments[i]),
				KzgProof:          kzg.Proof(blobs.Proofs[i]),
				BeaconBlockHeader: blkHeader,
				InclusionProof:    inclusionProof,
			}
			return nil
		})
	}

	return &datypes.BlobSidecars{Sidecars: sidecars}, g.Wait()
}
