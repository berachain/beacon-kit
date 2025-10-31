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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/blockchain/testhelpers"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dablob "github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/da/blobreactor"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/observability/metrics/discard"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/storage/filedb"
	"github.com/berachain/beacon-kit/testing/simulated"
	cmtconfig "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/p2p"
	"github.com/stretchr/testify/require"
)

// TestBlobFetcher_MultiNodeFetch tests Node1 fetching blobs from Node2 via P2P blob reactor
func (s *SimulatedSuite) TestBlobFetcher_MultiNodeFetch() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Move chain forward one block
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	startTime := time.Now()
	_, _, _ = s.MoveChainToHeight(s.T(), 1, 1, nodeAddress, startTime)

	// Create test blobs
	testSlot := math.Slot(100)
	blobs := []*eip4844.Blob{{1, 2, 3}, {4, 5, 6}}

	// Create sidecars, block, and commitments
	sidecars, block, commitments := createTestSidecars(s.T(), s, blobs, testSlot)
	s.Require().Len(sidecars, 2)
	s.Require().Len(commitments, 2)

	// Setup two nodes: Node1 (empty) and Node2 (has blobs)
	node1HomeDir := filepath.Join(os.TempDir(), "node1_multinode_test")
	node2HomeDir := filepath.Join(os.TempDir(), "node2_multinode_test")
	defer os.RemoveAll(node1HomeDir)
	defer os.RemoveAll(node2HomeDir)

	node1Store := createBlobStore(node1HomeDir)
	node2Store := createBlobStore(node2HomeDir)
	s.Require().NoError(node2Store.Persist(sidecars))

	node1Reactor := blobreactor.NewBlobReactor(
		node1Store,
		log.NewNopLogger(),
		blobreactor.Config{RequestTimeout: 5 * time.Second},
		blobreactor.NewMetrics(discard.NewFactory()),
	)
	node2Reactor := blobreactor.NewBlobReactor(
		node2Store,
		log.NewNopLogger(),
		blobreactor.Config{RequestTimeout: 5 * time.Second},
		blobreactor.NewMetrics(discard.NewFactory()),
	)

	// Connect via P2P
	switches := setupP2PReactors([]*blobreactor.BlobReactor{node1Reactor, node2Reactor})
	defer func() {
		for _, sw := range switches {
			_ = sw.Stop()
		}
	}()

	// Create and start Node1's blob fetcher
	node1Fetcher, err := blockchain.NewBlobFetcher(
		filepath.Join(node1HomeDir, "data"),
		log.NewNopLogger(),
		s.TestNode.BlobProcessor,
		node1Reactor,
		testhelpers.NewSimpleStorageBackend(node1Store),
		s.TestNode.ChainSpec,
		blockchain.BlobFetcherConfig{
			CheckInterval: 100 * time.Millisecond,
			RetryInterval: 200 * time.Millisecond,
			MaxRetries:    3,
		},
		blockchain.NewBlobFetcherMetrics(discard.NewFactory()),
	)
	s.Require().NoError(err)
	node1Fetcher.Start(s.CtxApp)

	// Set head slots on both nodes (within DA period)
	node1Reactor.SetHeadSlot(testSlot + 10)
	node2Reactor.SetHeadSlot(testSlot + 10)
	node1Fetcher.SetHeadSlot(testSlot + 10)

	// Queue blob request, wait for it to be downloaded and validate
	s.Require().NoError(node1Fetcher.QueueBlobRequest(block))
	time.Sleep(500 * time.Millisecond)
	storedSidecars, err := node1Store.GetBlobSidecars(testSlot)
	s.Require().NoError(err)
	s.Require().Len(storedSidecars, 2)
	verifyBlobCommitments(s.T(), storedSidecars, commitments)
	s.Require().True(assertQueueLength(node1HomeDir, 0), "queue should be empty after successful fetch")

	node1Fetcher.Stop()
}

// Helper to check queue state
func assertQueueLength(homeDir string, expected int) bool {
	queueDir := filepath.Join(homeDir, "data", "blobs", "download_queue")
	files, err := os.ReadDir(queueDir)
	if err != nil {
		return expected == 0
	}

	jsonFiles := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			jsonFiles++
		}
	}
	return jsonFiles == expected
}

// Helper to create test sidecars with minimal viable data for testing
func createTestSidecars(t *testing.T, s *SimulatedSuite, blobs []*eip4844.Blob, slot math.Slot) (
	datypes.BlobSidecars, *ctypes.BeaconBlock, []eip4844.KZGCommitment,
) {
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(s.Require(), blobs, s.TestNode.KZGVerifier)

	block, err := ctypes.NewBeaconBlockWithVersion(slot, 0, common.Root{}, version.Deneb())
	s.Require().NoError(err)
	block.Body.SetBlobKzgCommitments(
		eip4844.KZGCommitments[common.ExecutionHash](commitments),
	)
	signedHeader := ctypes.NewSignedBeaconBlockHeader(block.GetHeader(), crypto.BLSSignature{})

	sidecarFactory := dablob.NewSidecarFactory(dablob.NewFactoryMetrics(discard.NewFactory()))
	sidecars := make(datypes.BlobSidecars, len(blobs))
	for i := range blobs {
		inclusionProof, err := sidecarFactory.BuildKZGInclusionProof(block.Body, math.U64(i))
		s.Require().NoError(err)

		sidecars[i] = &datypes.BlobSidecar{
			Index:                   uint64(i),
			Blob:                    *blobs[i],
			KzgCommitment:           commitments[i],
			KzgProof:                proofs[i],
			SignedBeaconBlockHeader: signedHeader,
			InclusionProof:          inclusionProof,
		}
	}

	return sidecars, block, commitments
}

// Helper to create a blob availability store for testing
func createBlobStore(homeDir string) *dastore.Store {
	return dastore.New(
		filedb.NewRangeDB(
			filedb.NewDB(
				filedb.WithRootDirectory(filepath.Join(homeDir, "data", "blobs")),
				filedb.WithFileExtension("ssz"),
				filedb.WithDirectoryPermissions(os.ModePerm),
				filedb.WithLogger(log.NewNopLogger()),
			),
		),
		log.NewNopLogger(),
	)
}

// Helper to setup P2P connected reactors
func setupP2PReactors(reactors []*blobreactor.BlobReactor) []*p2p.Switch {
	p2pConfig := cmtconfig.DefaultP2PConfig()
	p2pConfig.ListenAddress = "tcp://127.0.0.1:0"

	initSwitch := func(i int, sw *p2p.Switch) *p2p.Switch {
		sw.AddReactor(blobreactor.ReactorName, reactors[i])
		return sw
	}
	return p2p.MakeConnectedSwitches(p2pConfig, len(reactors), initSwitch, p2p.Connect2Switches)
}

// Helper to verify blobs match expected commitments
func verifyBlobCommitments(t *testing.T, sidecars datypes.BlobSidecars, expectedCommitments []eip4844.KZGCommitment) {
	t.Helper()
	require := require.New(t)

	sidecarsByIndex := make(map[uint64]*datypes.BlobSidecar)
	for _, sidecar := range sidecars {
		sidecarsByIndex[sidecar.Index] = sidecar
	}

	for i, expectedCommitment := range expectedCommitments {
		sidecar, exists := sidecarsByIndex[uint64(i)]
		require.True(exists, "sidecar index %d should exist", i)
		require.Equal(expectedCommitment, sidecar.KzgCommitment, "commitment mismatch at index %d", i)
	}
}
