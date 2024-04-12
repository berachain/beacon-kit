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
	"github.com/berachain/beacon-kit/mod/config/params"
	"github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives/engine"
	"github.com/berachain/beacon-kit/mod/primitives/kzg"
	"golang.org/x/sync/errgroup"
)

// SidecarFactory is a factory for sidecars.
type SidecarFactory struct {
	cfg         *params.BeaconChainConfig
	kzgPosition uint64
}

// NewSidecarFactory creates a new sidecar factory.
func NewSidecarFactory(
	cfg *params.BeaconChainConfig,
) *SidecarFactory {
	return &SidecarFactory{
		cfg: cfg,
		// TODO: This should be configurable / modular.
		kzgPosition: types.KZGPositionDeneb,
	}
}

// BuildSidecar builds a sidecar.
func (f *SidecarFactory) BuildSidecars(
	blk types.BeaconBlock,
	blobs *engine.BlobsBundleV1,
) (*datypes.BlobSidecars, error) {
	numBlobs := uint64(len(blobs.Blobs))
	sidecars := make([]*datypes.BlobSidecar, numBlobs)
	g := errgroup.Group{}
	for i := range numBlobs {
		g.Go(func() error {
			var err error
			blob := kzg.Blob(blobs.Blobs[i])
			sidecars[i], err = types.BuildBlobSidecar(
				i,
				blk,
				f.kzgPosition,
				&blob,
				kzg.Commitment(blobs.Commitments[i]),
				kzg.Proof(blobs.Proofs[i]),
			)
			return err
		})
	}

	return &datypes.BlobSidecars{Sidecars: sidecars}, g.Wait()
}

// BuildInclusionProof builds an inclusion proof.
func (f *SidecarFactory) BuildInclusionProof(
	body types.BeaconBlockBody,
	index uint64,
) ([][32]byte, error) {
	return types.MerkleProofKZGCommitment(
		body, types.KZGPositionDeneb, index)
}
