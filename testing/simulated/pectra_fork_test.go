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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package simulated_test

import (
	"bytes"
	"context"
	"path"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// PectraForkSuite defines our test suite for Pectra related work using simulated Comet component.
type PectraForkSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors
}

// TestSimulatedCometComponent runs the test suite.
func TestPectraForkSuite(t *testing.T) {
	suite.Run(t, new(PectraForkSuite))
}

// SetupTest initializes the test environment.
func (s *PectraForkSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-fork-genesis.json"
	chainSpecFunc := simulated.ProvidePectraForkTestChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), chainSpec, s.HomeDir, elGenesisPath)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Geth node.
	elNode := execution.NewRethNode(s.HomeDir, execution.ValidRethImage())
	elHandle, authRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
	s.ElHandle = elHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.LogBuffer = new(bytes.Buffer)
	logger := phuslu.NewLogger(s.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component and electra genesis chain spec
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     authRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})

	s.SimComet = s.TestNode.SimComet

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.TestNode.Start(s.CtxApp)
	}()

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.EngineClient)
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

// TearDownTest cleans up the test environment.
func (s *PectraForkSuite) TearDownTest() {
	// If the test has failed, log additional information.
	if s.T().Failed() {
		s.T().Log(s.LogBuffer.String())
	}
	if err := s.ElHandle.Close(); err != nil {
		s.T().Error("Error closing EL handle:", err)
	}
	// mimics the behaviour of shutdown func
	s.CtxAppCancelFn()
	s.TestNode.ServiceRegistry.StopAll()
}

// TestTimestampFork_ELAndCLInSync_IsSuccessful tests that we can fork successfully if EL and CL have synced timestamps
// The forks timestamp at Unix 0, as the genesis at Unix 0, Cancun is at 10 and Prague is at 20.
func (s *PectraForkSuite) TestTimestampFork_ELAndCLInSync_IsSuccessful() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// Prepare a block proposal.
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	expectedMessages := []string{
		"Finalizing block with fork version service=blockchain\u001B[0m block=1\u001B[0m fork=0x04010000\u001B[0m",
		"Finalizing block with fork version service=blockchain\u001B[0m block=2\u001B[0m fork=0x04010000\u001B[0m",
		"Finalizing block with fork version service=blockchain\u001B[0m block=3\u001B[0m fork=0x04010000\u001B[0m",
		"Finalizing block with fork version service=blockchain\u001B[0m block=4\u001B[0m fork=0x04010000\u001B[0m",
		"Finalizing block with fork version service=blockchain\u001B[0m block=5\u001B[0m fork=0x05000000\u001B[0m",
		"Finalizing block with fork version service=blockchain\u001B[0m block=6\u001B[0m fork=0x05000000\u001B[0m",
		"Finalizing block with fork version service=blockchain\u001B[0m block=7\u001B[0m fork=0x05000000\u001B[0m",
	}

	startHeight := int64(1)
	iterations := int64(len(expectedMessages))
	expectedMessagesIdx := 0
	for currentHeight := startHeight; currentHeight < startHeight+iterations; currentHeight++ {
		// We set the consensus time to currentHeight * 2 to mimick a functional 2s block time.
		// This assumption is not always true but used in a valid test.
		consensusTime := time.Unix(currentHeight*2, 0)
		proposal, err := s.SimComet.Comet.PrepareProposal(s.CtxComet, &types.PrepareProposalRequest{
			Height:          currentHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		})
		s.Require().NoError(err)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block.
		s.LogBuffer.Reset()
		finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.CtxComet, &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		})
		s.Require().Contains(s.LogBuffer.String(), expectedMessages[expectedMessagesIdx])
		s.Require().NoError(err)
		s.Require().NotEmpty(finalizeResp)
		// Commit the block.
		_, err = s.SimComet.Comet.Commit(s.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		expectedMessagesIdx++
	}
}
