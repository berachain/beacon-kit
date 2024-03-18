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
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"golang.org/x/sync/errgroup"
)

// SideCars is a slice of blob side cars to be included in the block.
type BlobSidecars struct {
	Sidecars []*BlobSidecar `ssz-max:"6"`
}

// BlobSidecar is a struct that contains blobs and their associated information.
type BlobSidecar struct {
	Index          uint64
	Blob           []byte   `ssz-size:"131072"`
	KzgCommitment  []byte   `ssz-size:"48"`
	KzgProof       []byte   `ssz-size:"48"`
	InclusionProof [][]byte `ssz-size:"8,32"`
}

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	blk BeaconBlock,
	blobs *enginetypes.BlobsBundleV1,
) (*BlobSidecars, error) {
	numBlobs := uint64(len(blobs.Blobs))
	sidecars := make([]*BlobSidecar, numBlobs)
	g := errgroup.Group{}
	for i := uint64(0); i < numBlobs; i++ {
		i := i // capture range variable
		g.Go(func() error {
			// Create Inclusion Proof
			inclusionProof, err := MerkleProofKZGCommitment(blk, i)
			if err != nil {
				return err
			}

			sidecars[i] = &BlobSidecar{
				Index:          i,
				Blob:           blobs.Blobs[i],
				KzgCommitment:  blobs.Commitments[i],
				KzgProof:       blobs.Proofs[i],
				InclusionProof: inclusionProof,
			}
			return nil
		})
	}

	return &BlobSidecars{Sidecars: sidecars}, g.Wait()
}
