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
	"context"
	"path"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// BlobConsensusSuite runs the simulated Comet component against a chain spec that enables the blob-consensus
// transition at height 2, exercising the enable-height boundary (two txs before, one tx after) through the
// real PrepareProposal/ProcessProposal/FinalizeBlock/Commit code paths in CI.
type BlobConsensusSuite struct {
	suite.Suite
	simulated.SharedAccessors
}

// TestBlobConsensusSuite runs the test suite.
func TestBlobConsensusSuite(t *testing.T) {
	suite.Run(t, new(BlobConsensusSuite))
}

// SetupTest initializes the test environment.
func (s *BlobConsensusSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/eth-genesis.json"
	chainSpecFunc := simulated.ProvideBlobConsensusTestChainSpec
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	configs, genesisValidatorsRoot := simulated.InitializeHomeDirs(s.T(), chainSpec, elGenesisPath, s.HomeDir)
	cometConfig := configs[0]
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Reth node.
	elNode := execution.NewRethNode(s.HomeDir, execution.ValidRethImage())
	elHandle, authRPC, elRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
	s.ElHandle = elHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component.
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)
	s.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     authRPC,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})

	s.SimComet = s.TestNode.SimComet

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.TestNode.Start(s.CtxApp)
	}()

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.ContractBackend)
	err = simulated.WaitTillServicesStarted(s.LogBuffer, 10*time.Second, 50*time.Millisecond)
	s.Require().NoError(err)
}

// TearDownTest cleans up the test environment.
func (s *BlobConsensusSuite) TearDownTest() {
	s.CleanupTest(s.T())
}

// TestBlobConsensusTransition_TxLayoutAcrossEnableHeight drives the chain through the blob-consensus enable
// height. Height 1 must carry two consensus txs (legacy layout), and every height from 2 on exactly one, while
// blocks keep proposing, verifying, and finalizing through the full ABCI flow (including the state-cache
// finalization path, active from height 2 via stable block time). With no blob transactions in the payloads
// the sidecar sets are empty; the blob-carrying flow is covered by the unit tests and the devnet harnesses.
func (s *BlobConsensusSuite) TestBlobConsensusTransition_TxLayoutAcrossEnableHeight() {
	const blocksToRun = 5

	s.InitializeChain(s.T(), 1)
	nodeAddress, err := s.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.SimComet.Comet.SetNodeAddress(nodeAddress)

	// MoveChainToHeight asserts the per-height expected tx count internally (two below the enable height, one
	// at and above it) and fails if any proposal is rejected or fails to finalize.
	proposals, _, _ := s.MoveChainToHeight(s.T(), 1, blocksToRun, nodeAddress, time.Now())
	s.Require().Len(proposals, blocksToRun)

	// Belt and suspenders: re-assert the layout boundary on the recorded proposals.
	s.Require().Len(proposals[0].Txs, 2, "height 1 is below the enable height and must carry block+sidecars")
	for i := 1; i < blocksToRun; i++ {
		s.Require().Len(proposals[i].Txs, 1, "height %d must carry only the block", i+1)
	}

	// Nothing is pending in the background blob fetcher: every finalized block's data was available.
	s.Require().Zero(s.TestNode.Blockchain.PendingBlobRequests())
}
