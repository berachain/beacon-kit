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

	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	datypes "github.com/berachain/beacon-kit/mod/interfaces/pkg/da/types"
	engineprimitivesI "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/telemetry"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"golang.org/x/sync/errgroup"
)

// SidecarFactory is a factory for sidecars.
type SidecarFactory[
	BeaconBlockT types.BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BlobsBundleT engineprimitivesI.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	BlobSidecarT datypes.BlobSidecar[
		BlobSidecarT, BeaconBlockHeaderT,
	],
	BlobSidecarsT datypes.BlobSidecars[BlobSidecarsT, BlobSidecarT],
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
] struct {
	// chainSpec defines the specifications of the blockchain.
	chainSpec common.ChainSpec
	// kzgPosition is the position of the KZG commitment in the block.
	//
	// TODO: This needs to be made configurable / modular.
	kzgPosition uint64
	// metrics is used to collect and report factory metrics.
	metrics *factoryMetrics
}

// NewSidecarFactory creates a new sidecar factory.
func NewSidecarFactory[
	BeaconBlockT types.BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BlobsBundleT engineprimitivesI.BlobsBundle[
		eip4844.KZGCommitment, eip4844.KZGProof, eip4844.Blob,
	],
	BlobSidecarT datypes.BlobSidecar[
		BlobSidecarT, BeaconBlockHeaderT,
	],
	BlobSidecarsT datypes.BlobSidecars[BlobSidecarsT, BlobSidecarT],
	DepositT any,
	Eth1DataT any,
	ExecutionPayloadT any,
](
	chainSpec common.ChainSpec,
	// todo: calculate from config.
	kzgPosition uint64,
	telemetrySink telemetry.Sink,
) *SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BlobsBundleT,
	BlobSidecarT, BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
] {
	return &SidecarFactory[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BlobsBundleT,
		BlobSidecarT, BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
	]{
		chainSpec: chainSpec,
		// TODO: This should be configurable / modular.
		kzgPosition: kzgPosition,
		metrics:     newFactoryMetrics(telemetrySink),
	}
}

// BuildSidecars builds a sidecar.
func (f *SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BlobsBundleT,
	BlobSidecarT, BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
]) BuildSidecars(
	blk BeaconBlockT,
	bundle BlobsBundleT,
) (BlobSidecarsT, error) {
	var (
		out         BlobSidecarsT
		blobs       = bundle.GetBlobs()
		commitments = bundle.GetCommitments()
		proofs      = bundle.GetProofs()
		numBlobs    = uint64(len(blobs))
		sidecars    = make([]BlobSidecarT, numBlobs)
		body        = blk.GetBody()
		g           = errgroup.Group{}
	)

	startTime := time.Now()
	defer f.metrics.measureBuildSidecarsDuration(
		startTime, math.U64(numBlobs),
	)
	var sidecar BlobSidecarT
	for i := range numBlobs {
		g.Go(func() error {
			inclusionProof, err := f.BuildKZGInclusionProof(
				body, math.U64(i),
			)
			if err != nil {
				return err
			}
			sidecars[i] = sidecar.New(
				math.U64(i), blk.GetHeader(),
				blobs[i],
				commitments[i],
				proofs[i],
				inclusionProof,
			)
			return nil
		})
	}
	return out.NewFromSidecars(sidecars), g.Wait()
}

// BuildKZGInclusionProof builds a KZG inclusion proof.
func (f *SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BlobsBundleT,
	BlobSidecarT, BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
]) BuildKZGInclusionProof(
	body BeaconBlockBodyT,
	index math.U64,
) ([][32]byte, error) {
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
func (f *SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BlobsBundleT,
	BlobSidecarT, BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
]) BuildBlockBodyProof(
	body BeaconBlockBodyT,
) ([][32]byte, error) {
	startTime := time.Now()
	defer f.metrics.measureBuildBlockBodyProofDuration(startTime)
	membersRoots, err := body.GetTopLevelRoots()
	if err != nil {
		return nil, err
	}

	tree, err := merkle.NewTreeWithMaxLeaves[[32]byte](
		membersRoots,
		body.Length()-1,
	)
	if err != nil {
		return nil, err
	}

	return tree.MerkleProof(f.kzgPosition)
}

// BuildCommitmentProof builds a commitment proof.
func (f *SidecarFactory[
	BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT, BlobsBundleT,
	BlobSidecarT, BlobSidecarsT, DepositT, Eth1DataT, ExecutionPayloadT,
]) BuildCommitmentProof(
	body BeaconBlockBodyT,
	index math.U64,
) ([][32]byte, error) {
	startTime := time.Now()
	defer f.metrics.measureBuildCommitmentProofDuration(startTime)

	bodyTree, err := merkle.NewTreeWithMaxLeaves[[32]byte](
		body.GetBlobKzgCommitments().Leafify(),
		f.chainSpec.MaxBlobCommitmentsPerBlock(),
	)
	if err != nil {
		return nil, err
	}

	return bodyTree.MerkleProofWithMixin(index.Unwrap())
}
