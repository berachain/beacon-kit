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
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
)

// BuildBlobSidecar creates a blob sidecar from the given blobs and
// beacon block.
func BuildBlobSidecar(
	index uint64,
	blk BeaconBlock,
	kzgPosition uint64,
	blob *kzg.Blob,
	commitment kzg.Commitment,
	proof kzg.Proof,
) (*datypes.BlobSidecar, error) {
	// Create Inclusion Proof
	inclusionProof, err := MerkleProofKZGCommitment(
		blk.GetBody(), kzgPosition, index,
	)
	if err != nil {
		return nil, err
	}
	return &datypes.BlobSidecar{
		Index:             index,
		Blob:              *blob,
		KzgCommitment:     commitment,
		KzgProof:          proof,
		BeaconBlockHeader: blk.GetHeader(),
		InclusionProof:    inclusionProof,
	}, nil
}
