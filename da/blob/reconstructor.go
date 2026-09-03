// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"context"
	"fmt"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/kzg"
	datypes "github.com/berachain/beacon-kit/da/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	gethengine "github.com/ethereum/go-ethereum/beacon/engine"
)

// ErrBlobsNotInELPool is returned when the execution client does not have all blobs of a block in its blob pool.
var ErrBlobsNotInELPool = errors.New("execution client does not have the block's blobs")

// ELBlobsFetcher fetches blobs from the local execution client's blob pool. Satisfied by *client.EngineClient.
type ELBlobsFetcher interface {
	GetBlobsV2(ctx context.Context, versionedHashes []common.ExecutionHash) ([]*gethengine.BlobAndProofV2, error)
}

// Reconstructor rebuilds a block's canonical blob sidecars from blobs fetched off the local execution client. Blob transactions gossip
// through the EL mempool with their blobs attached, so at the tip of the chain this is usually a local hit that avoids any
// consensus-layer p2p round trip.
//
// engine_getBlobsV2 returns cell proofs (Osaka onwards) while the sidecar format carries the single blob KZG proof, so the proof is
// recomputed locally; being deterministic, the resulting sidecars are byte-identical to the proposer's.
type Reconstructor struct {
	engineClient ELBlobsFetcher
	factory      *SidecarFactory
	prover       kzg.BlobProofProver
	logger       log.Logger
}

// NewReconstructor creates a Reconstructor.
func NewReconstructor(
	engineClient ELBlobsFetcher,
	factory *SidecarFactory,
	prover kzg.BlobProofProver,
	logger log.Logger,
) *Reconstructor {
	return &Reconstructor{
		engineClient: engineClient,
		factory:      factory,
		prover:       prover,
		logger:       logger,
	}
}

// ReconstructSidecars fetches the blobs committed to by the signed block from the execution client and rebuilds the block's sidecars
// (inclusion proofs and all). Returns ErrBlobsNotInELPool if the execution client is missing any blob (the engine_getBlobsV2 response is
// all-or-nothing).
func (r *Reconstructor) ReconstructSidecars(ctx context.Context, signedBlk *ctypes.SignedBeaconBlock) (datypes.BlobSidecars, error) {
	blk := signedBlk.GetBeaconBlock()
	commitments := blk.GetBody().GetBlobKzgCommitments()
	if len(commitments) == 0 {
		return datypes.BlobSidecars{}, nil
	}

	blobsAndProofs, err := r.engineClient.GetBlobsV2(ctx, commitments.ToVersionedHashes())
	if err != nil {
		return nil, fmt.Errorf("engine_getBlobsV2 failed: %w", err)
	}
	// A nil/empty response (the EL does not hold every requested blob) fails the count check below.
	if len(blobsAndProofs) != len(commitments) {
		return nil, fmt.Errorf(
			"engine_getBlobsV2 returned %d blobs, expected %d: %w",
			len(blobsAndProofs), len(commitments), ErrBlobsNotInELPool,
		)
	}

	bundle := &engineprimitives.BlobsBundleV1{
		Commitments: make([]eip4844.KZGCommitment, len(commitments)),
		Proofs:      make([]eip4844.KZGProof, len(commitments)),
		Blobs:       make([]*eip4844.Blob, len(commitments)),
	}
	for i, bp := range blobsAndProofs {
		if bp == nil {
			return nil, fmt.Errorf("missing blob at index %d: %w", i, ErrBlobsNotInELPool)
		}
		if len(bp.Blob) != len(eip4844.Blob{}) {
			return nil, fmt.Errorf("blob at index %d has invalid length %d", i, len(bp.Blob))
		}

		blob := new(eip4844.Blob)
		copy(blob[:], bp.Blob)

		// The EL returns cell proofs; recompute the canonical blob proof.
		proof, proofErr := r.prover.ComputeBlobProof(blob, commitments[i])
		if proofErr != nil {
			return nil, fmt.Errorf("failed computing blob proof for index %d: %w", i, proofErr)
		}

		bundle.Blobs[i] = blob
		bundle.Commitments[i] = commitments[i]
		bundle.Proofs[i] = proof
	}

	sidecars, err := r.factory.BuildSidecars(signedBlk, bundle)
	if err != nil {
		return nil, fmt.Errorf("failed rebuilding sidecars from EL blobs: %w", err)
	}

	r.logger.Info("Reconstructed blob sidecars from execution client",
		"slot", blk.GetSlot().Unwrap(), "count", len(sidecars))
	return sidecars, nil
}
