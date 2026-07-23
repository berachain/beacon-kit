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

//go:build simulated

package simulated_test

import (
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/blockchain/mocks"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	dablob "github.com/berachain/beacon-kit/da/blob"
	"github.com/berachain/beacon-kit/da/blobreactor"
	dastore "github.com/berachain/beacon-kit/da/store"
	datypes "github.com/berachain/beacon-kit/da/types"
	"github.com/berachain/beacon-kit/node-core/components/metrics"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/storage/filedb"
	"github.com/berachain/beacon-kit/testing/simulated"
	cmtconfig "github.com/cometbft/cometbft/config"
	"github.com/cometbft/cometbft/p2p"
)

// blobTestEnv wires two blob reactors (node1 empty, node2 serving) over real
// connected switches, plus a background fetcher for node1.
type blobTestEnv struct {
	node1Store   *dastore.Store
	node2Store   *dastore.Store
	node1Fetcher blockchain.BlobFetcher
	switches     []*p2p.Switch
}

func (s *SimulatedSuite) setupBlobFetchEnv(headSlot math.Slot) *blobTestEnv {
	node1HomeDir := filepath.Join(s.T().TempDir(), "node1")
	node2HomeDir := filepath.Join(s.T().TempDir(), "node2")

	node1Store := createBlobStore(node1HomeDir)
	node2Store := createBlobStore(node2HomeDir)

	blobProcessor := dablob.NewProcessor(
		log.NewNopLogger(), s.TestNode.KZGVerifier, metrics.NewNoOpTelemetrySink())

	cfg := blobreactor.Config{RequestTimeout: 5 * time.Second, FetchTimeout: 10 * time.Second}
	node1Reactor := blobreactor.NewBlobReactor(
		node1Store, blobProcessor, log.NewNopLogger(), cfg, 6, metrics.NewNoOpTelemetrySink())
	node2Reactor := blobreactor.NewBlobReactor(
		node2Store, blobProcessor, log.NewNopLogger(), cfg, 6, metrics.NewNoOpTelemetrySink())

	switches := setupP2PReactors([]*blobreactor.BlobReactor{node1Reactor, node2Reactor})

	node1Storage := &mocks.StorageBackend{}
	node1Storage.On("AvailabilityStore").Return(node1Store)

	node1Fetcher, err := blockchain.NewBlobFetcher(
		filepath.Join(node1HomeDir, "data"),
		log.NewNopLogger(),
		blobProcessor,
		node1Reactor,
		node1Storage,
		s.TestNode.ChainSpec,
		blockchain.BlobFetcherConfig{
			CheckInterval: 100 * time.Millisecond,
			RetryInterval: 200 * time.Millisecond,
		},
		metrics.NewNoOpTelemetrySink(),
	)
	s.Require().NoError(err)
	node1Fetcher.Start(s.CtxApp)

	node1Reactor.SetHeadSlot(headSlot)
	node2Reactor.SetHeadSlot(headSlot)
	node1Fetcher.SetHeadSlot(headSlot)

	env := &blobTestEnv{
		node1Store:   node1Store,
		node2Store:   node2Store,
		node1Fetcher: node1Fetcher,
		switches:     switches,
	}
	s.T().Cleanup(func() {
		node1Fetcher.Stop()
		for _, sw := range switches {
			_ = sw.Stop()
		}
	})
	return env
}

// TestBlobFetcher_FetchLifecycle drives one request through the whole background-fetch story. While the only
// peer serves sidecars for a different block at the same slot, the request must keep retrying (never silently
// dropped) and nothing may be persisted. Once the correct sidecars appear on the peer, the next retry fetches,
// verifies and persists them, and the queue drains.
func (s *SimulatedSuite) TestBlobFetcher_FetchLifecycle() {
	s.InitializeChain(s.T(), 1)

	testSlot := math.Slot(50)
	blobs := []*eip4844.Blob{{1, 2, 3}, {4, 5, 6}}
	sidecars, signedBlk, commitments := createFetchableSidecars(s, blobs, testSlot)

	// Node2 initially holds sidecars for a DIFFERENT block at the same slot (one blob, other content).
	wrongSidecars, _, _ := createFetchableSidecars(s, []*eip4844.Blob{{7, 7, 7}}, testSlot)

	env := s.setupBlobFetchEnv(testSlot + 10)
	s.Require().NoError(env.node2Store.Persist(wrongSidecars))
	s.Require().NoError(env.node1Fetcher.QueueBlobRequest(signedBlk))

	// Several retry intervals elapse against the junk-serving peer: the request must survive and the
	// mismatched sidecars must never reach the store.
	time.Sleep(1 * time.Second)
	s.Require().Equal(1, env.node1Fetcher.PendingRequests(),
		"request must not be dropped while its slot is within the DA window")
	stored, err := env.node1Store.GetBlobSidecars(testSlot)
	s.Require().NoError(err)
	s.Require().Empty(stored, "nothing may be persisted from a mismatched response")

	// The real sidecars appear on node2; node1's next retry fetches and persists them.
	s.Require().NoError(env.node2Store.DeleteBlobSidecars(testSlot))
	s.Require().NoError(env.node2Store.Persist(sidecars))
	s.Require().Eventually(func() bool {
		return env.node1Fetcher.PendingRequests() == 0
	}, 10*time.Second, 100*time.Millisecond, "queue should drain after successful fetch")

	storedSidecars, err := env.node1Store.GetBlobSidecars(testSlot)
	s.Require().NoError(err)
	s.Require().Len(storedSidecars, 2)
	verifyBlobCommitments(s, storedSidecars, commitments)
}

// createFetchableSidecars builds a signed block committing to the given blobs
// and its canonical sidecars (real KZG proofs and inclusion proofs), so the
// full fetch verification passes.
func createFetchableSidecars(
	s *SimulatedSuite, blobs []*eip4844.Blob, slot math.Slot,
) (datypes.BlobSidecars, *ctypes.SignedBeaconBlock, []eip4844.KZGCommitment) {
	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(s.Require(), blobs, s.TestNode.KZGVerifier)

	block, err := ctypes.NewBeaconBlockWithVersion(slot, 0, common.Root{}, version.Deneb())
	s.Require().NoError(err)
	block.Body.SetBlobKzgCommitments(eip4844.KZGCommitments[common.ExecutionHash](commitments))
	signedBlk := &ctypes.SignedBeaconBlock{BeaconBlock: block, Signature: crypto.BLSSignature{}}
	signedHeader := ctypes.NewSignedBeaconBlockHeader(block.GetHeader(), crypto.BLSSignature{})

	sidecarFactory := dablob.NewSidecarFactory(metrics.NewNoOpTelemetrySink())
	sidecars := make(datypes.BlobSidecars, len(blobs))
	for i := range blobs {
		inclusionProof, proofErr := sidecarFactory.BuildKZGInclusionProof(block.Body, math.U64(i)) //#nosec:G115
		s.Require().NoError(proofErr)

		sidecars[i] = datypes.BuildBlobSidecar(
			math.U64(i), //#nosec:G115
			signedHeader,
			blobs[i],
			commitments[i],
			proofs[i],
			inclusionProof,
		)
	}
	return sidecars, signedBlk, commitments
}

// createBlobStore creates a blob availability store rooted in homeDir.
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

// setupP2PReactors connects the given reactors over in-memory switches.
func setupP2PReactors(reactors []*blobreactor.BlobReactor) []*p2p.Switch {
	p2pConfig := cmtconfig.DefaultP2PConfig()
	p2pConfig.ListenAddress = "tcp://127.0.0.1:0"

	initSwitch := func(i int, sw *p2p.Switch) *p2p.Switch {
		sw.AddReactor(blobreactor.ReactorName, reactors[i])
		return sw
	}
	return p2p.MakeConnectedSwitches(p2pConfig, len(reactors), initSwitch, p2p.Connect2Switches)
}

// verifyBlobCommitments asserts the stored sidecars carry the expected
// commitments by index.
func verifyBlobCommitments(s *SimulatedSuite, sidecars datypes.BlobSidecars, expected []eip4844.KZGCommitment) {
	sidecarsByIndex := make(map[uint64]*datypes.BlobSidecar)
	for _, sidecar := range sidecars {
		sidecarsByIndex[sidecar.Index] = sidecar
	}
	for i, expectedCommitment := range expected {
		sidecar, exists := sidecarsByIndex[uint64(i)] //#nosec:G115
		s.Require().True(exists, "sidecar index %d should exist", i)
		s.Require().Equal(expectedCommitment, sidecar.KzgCommitment, "commitment mismatch at index %d", i)
	}
}
