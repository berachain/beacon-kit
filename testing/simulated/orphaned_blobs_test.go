//go:build simulated

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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package simulated_test

import (
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/kzg"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/stretchr/testify/require"
)

// TestOrphanedBlobCleanup tests that orphaned blob sidecars are properly cleaned up on node restart.
// This simulates the scenario where sidecars are saved to disk but the block finalization fails.
func (s *SimulatedSuite) TestOrphanedBlobCleanup() {
	// Initialize chain and move forward two blocks.
	s.InitializeChain(s.T())
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	_, _, proposalTime := s.MoveChainToHeight(s.T(), 1, 2, nodeAddress, time.Now())

	// Get the last committed block height.
	lastBlockHeight := s.SimComet.Comet.CommitMultiStore().LastCommitID().Version
	orphanedSlot := math.Slot(lastBlockHeight + 1)

	// Create and persist orphaned blob sidecars.
	// This simulates FinalizeSidecars succeeding but finalizeBeaconBlock failing.
	orphanedSidecars := createOrphanedSidecars(s.T(), orphanedSlot, s.TestNode.KZGVerifier)
	err = s.TestNode.StorageBackend.AvailabilityStore().Persist(orphanedSidecars)
	s.Require().NoError(err)

	// Verify orphaned blobs exist.
	sidecars, err := s.TestNode.StorageBackend.AvailabilityStore().GetBlobSidecars(orphanedSlot)
	s.Require().NoError(err)
	s.Require().Len(sidecars, 1)

	// Simulate node restart by calling PruneOrphanedBlobs.
	err = s.TestNode.Blockchain.PruneOrphanedBlobs(lastBlockHeight)
	s.Require().NoError(err)

	// Verify orphaned blobs were cleaned up.
	sidecars, err = s.TestNode.StorageBackend.AvailabilityStore().GetBlobSidecars(orphanedSlot)
	s.Require().NoError(err)
	s.Require().Empty(sidecars)

	// Verify chain continues normally.
	proposals, _, _ := s.MoveChainToHeight(s.T(), 3, 1, nodeAddress, proposalTime)
	s.Require().Len(proposals, 1)
}

// createOrphanedSidecars creates fake blob sidecars for testing orphaned blob cleanup.
func createOrphanedSidecars(
	t require.TestingT,
	slot math.Slot,
	verifier kzg.BlobProofVerifier,
) datypes.BlobSidecars {
	blobs := []*eip4844.Blob{{1, 2, 3}}
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(require.New(t), blobs, verifier)

	sidecars := make(datypes.BlobSidecars, len(blobs))
	for i := range blobs {
		sidecars[i] = datypes.BuildBlobSidecar(
			math.U64(i),
			&ctypes.SignedBeaconBlockHeader{
				Header:    &ctypes.BeaconBlockHeader{Slot: slot},
				Signature: crypto.BLSSignature{},
			},
			blobs[i],
			commitments[i],
			proofs[i],
			make([]common.Root, ctypes.KZGInclusionProofDepth),
		)
	}
	return sidecars
}
