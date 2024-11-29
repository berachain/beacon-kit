// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package blob

import (
	"time"

	constypes "github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"golang.org/x/sync/errgroup"
)

// SidecarFactory is a factory for sidecars.
type SidecarFactory[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody,
	BeaconBlockHeaderT any,
	ForkDataT any,
] struct {
	// chainSpec defines the specifications of the blockchain.
	chainSpec ChainSpec
	// kzgPosition is the position of the KZG commitment in the block.
	//
	// TODO: This needs to be made configurable / modular.
	kzgPosition uint64
	// metrics is used to collect and report factory metrics.
	metrics *factoryMetrics
}

// NewSidecarFactory creates a new sidecar factory.
func NewSidecarFactory[
	BeaconBlockT BeaconBlock[BeaconBlockBodyT, BeaconBlockHeaderT],
	BeaconBlockBodyT BeaconBlockBody,
	BeaconBlockHeaderT any,
	ForkDataT any,
](
	chainSpec ChainSpec,
	// todo: calculate from config.
	kzgPosition uint64,
	telemetrySink TelemetrySink,
) *SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, ForkDataT,
] {
	return &SidecarFactory[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, ForkDataT,
	]{
		chainSpec: chainSpec,
		// TODO: This should be configurable / modular.
		kzgPosition: kzgPosition,
		metrics:     newFactoryMetrics(telemetrySink),
	}
}

// BuildSidecars builds a sidecar.
func (f *SidecarFactory[BeaconBlockT, _, BeaconBlockHeaderT, ForkDataT]) BuildSidecars(
	blk BeaconBlockT,
	bundle engineprimitives.BlobsBundle,
	signer crypto.BLSSigner,
	forkData ForkDataT,
) (*types.BlobSidecars, error) {
	var (
		blobs       = bundle.GetBlobs()
		commitments = bundle.GetCommitments()
		proofs      = bundle.GetProofs()
		numBlobs    = uint64(len(blobs))
		sidecars    = make([]*types.BlobSidecar, numBlobs)
		body        = blk.GetBody()
		g           = errgroup.Group{}
		header      = any(blk.GetHeader()).(*constypes.BeaconBlockHeader)
	)

	startTime := time.Now()
	defer f.metrics.measureBuildSidecarsDuration(
		startTime, math.U64(numBlobs),
	)

	// Contrary to the spec, we do not need to sign the full
	// block, because the header embeds the body's hash tree root
	// already. We just need a bond between the block signer (already
	// tied to CometBFT's ProposerAddress) and the sidecars.

	domain := any(forkData).(*constypes.ForkData).ComputeDomain(f.chainSpec.DomainTypeProposer())
	signingRoot := constypes.ComputeSigningRoot(
		header,
		domain,
	)
	signature, err := signer.Sign(signingRoot[:])
	if err != nil {
		return nil, err
	}
	sigHeader := constypes.NewSignedBeaconBlockHeader(header, signature)

	for i := range numBlobs {
		g.Go(func() error {
			inclusionProof, err := f.BuildKZGInclusionProof(
				body, math.U64(i),
			)
			if err != nil {
				return err
			}
			sidecars[i] = types.BuildBlobSidecar(
				math.U64(i),
				sigHeader,
				blobs[i],
				commitments[i],
				proofs[i],
				inclusionProof,
			)
			return nil
		})
	}

	return &types.BlobSidecars{Sidecars: sidecars}, g.Wait()
}

// BuildKZGInclusionProof builds a KZG inclusion proof.
func (f *SidecarFactory[_, BeaconBlockBodyT, _, _]) BuildKZGInclusionProof(
	body BeaconBlockBodyT,
	index math.U64,
) ([]common.Root, error) {
	startTime := time.Now()
	defer f.metrics.measureBuildKZGInclusionProofDuration(startTime)

	// Build the merkle proof to the commitment within the
	// list of commitments.
	commitmentsProof, err := f.BuildCommitmentProof(body, index)
	if err != nil {
		return nil, err
	}

	// Build the merkle proof for the body root.
	bodyProof, err := f.BuildBlockBodyProof(body)
	if err != nil {
		return nil, err
	}

	// By property of the merkle tree, we can concatenate the
	// two proofs to get the final proof.
	return append(commitmentsProof, bodyProof...), nil
}

// BuildBlockBodyProof builds a block body proof.
func (f *SidecarFactory[_, BeaconBlockBodyT, _, _]) BuildBlockBodyProof(
	body BeaconBlockBodyT,
) ([]common.Root, error) {
	startTime := time.Now()
	defer f.metrics.measureBuildBlockBodyProofDuration(startTime)
	tree, err := merkle.NewTreeWithMaxLeaves[common.Root](
		body.GetTopLevelRoots(),
		body.Length()-1,
	)
	if err != nil {
		return nil, err
	}

	return tree.MerkleProof(f.kzgPosition)
}

// BuildCommitmentProof builds a commitment proof.
func (f *SidecarFactory[_, BeaconBlockBodyT, _, _]) BuildCommitmentProof(
	body BeaconBlockBodyT,
	index math.U64,
) ([]common.Root, error) {
	startTime := time.Now()
	defer f.metrics.measureBuildCommitmentProofDuration(startTime)
	bodyTree, err := merkle.NewTreeWithMaxLeaves[common.Root](
		body.GetBlobKzgCommitments().Leafify(),
		f.chainSpec.MaxBlobCommitmentsPerBlock(),
	)
	if err != nil {
		return nil, err
	}

	return bodyTree.MerkleProofWithMixin(index.Unwrap())
}
