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
	"fmt"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// SimulatedSuite defines our test suite for the simulated Comet component.
type SimulatedSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors

	SimComet              *simulated.SimComet
	LogBuffer             *bytes.Buffer
	GenesisValidatorsRoot common.Root
	SimulationClient      *execution.SimulationClient
}

// TestSimulatedCometComponent runs the test suite.
func TestSimulatedCometComponent(t *testing.T) {
	suite.Run(t, new(SimulatedSuite))
}

// SetupTest initializes the test environment.
func (s *SimulatedSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.Ctx, s.CancelFunc = context.WithCancel(context.Background())
	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), s.HomeDir)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Geth node.
	elNode := execution.NewGethNode(s.HomeDir, execution.ValidGethImageWithSimulate())
	elHandle, authRPC := elNode.Start(s.T())
	s.ElHandle = elHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.LogBuffer = new(bytes.Buffer)
	logger := phuslu.NewLogger(s.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component.
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	s.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     authRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})

	// Retrieve the simulated Comet service.
	var cometService *simulated.SimComet
	err := s.TestNode.FetchService(&cometService)
	s.Require().NoError(err)
	s.Require().NotNil(cometService)
	s.SimComet = cometService

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.TestNode.Start(s.Ctx)
	}()

	s.SimulationClient = execution.NewSimulationClient(s.TestNode.EngineClient)
	// Allow a short period for services to fully initialize.
	time.Sleep(2 * time.Second)
}

// TearDownTest cleans up the test environment.
func (s *SimulatedSuite) TearDownTest() {
	fmt.Println(s.LogBuffer.String())
	if err := s.ElHandle.CloseIfOpen(); err != nil {
		s.T().Error("Error closing EL handle:", err)
	}
	s.CancelFunc()
}

// initializeChain sets up the chain using the genesis file.
func (s *SimulatedSuite) initializeChain() {
	// Load the genesis state.
	appGenesis, err := genutiltypes.AppGenesisFromFile(s.HomeDir + "/config/genesis.json")
	s.Require().NoError(err)

	// Initialize the chain.
	initResp, err := s.SimComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId:       simulated.TestnetBeaconChainID,
		AppStateBytes: appGenesis.AppState,
	})
	s.Require().NoError(err)
	s.Require().Len(initResp.Validators, 1, "Expected 1 validator")

	// Verify that the deposit store contains the expected deposits.
	deposits, err := s.TestNode.StorageBackend.DepositStore().GetDepositsByIndex(
		s.Ctx,
		constants.FirstDepositIndex,
		constants.FirstDepositIndex+s.TestNode.ChainSpec.MaxDepositsPerBlock(),
	)
	s.Require().NoError(err)
	s.Require().Len(deposits, 1, "Expected 1 deposit")
}

// moveChainToHeight will iterate through the core loop `iterations` times, i.e. Propose, Process, Finalize and Commit.
// Returns the list of proposed comet blocks.
func (s *SimulatedSuite) moveChainToHeight(startHeight, iterations int64, proposer *signer.BLSSigner) []*types.PrepareProposalResponse {
	// Prepare a block proposal.
	pubkey, err := proposer.GetPubKey()
	s.Require().NoError(err)

	var proposedCometBlocks []*types.PrepareProposalResponse

	for currentHeight := startHeight; currentHeight < startHeight+iterations; currentHeight++ {
		proposalTime := time.Now()
		proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
			Height:          currentHeight,
			Time:            proposalTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(proposal)

		// Process the proposal.
		processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
			Time:            proposalTime,
		})
		s.Require().NoError(err)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block.
		finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.Ctx, &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().NotEmpty(finalizeResp)

		// Commit the block.
		_, err = s.SimComet.Comet.Commit(s.Ctx, &types.CommitRequest{})
		s.Require().NoError(err)

		// Record the Commit Block
		proposedCometBlocks = append(proposedCometBlocks, proposal)
	}
	return proposedCometBlocks
}
