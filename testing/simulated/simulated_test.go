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
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	mathpkg "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const finalizeWaitDuration = 500 * time.Millisecond

// SimulatedSuite defines our test suite for the simulated Comet component.
type SimulatedSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors

	SimComet              *simulated.SimComet
	LogBuffer             *bytes.Buffer
	GenesisValidatorsRoot common.Root
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
	elNode := execution.NewGethNode(s.HomeDir, execution.ValidGethImage())
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

	// Allow a short period for services to fully initialize.
	time.Sleep(2 * time.Second)
}

// TearDownTest cleans up the test environment.
func (s *SimulatedSuite) TearDownTest() {
	if err := s.ElHandle.Close(); err != nil {
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
		ChainId:       "test-mainnet-chain",
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

// TestFullLifecycle_ValidBlock_IsSuccessful tests that a valid block proposal is processed, finalized, and committed.
func (s *SimulatedSuite) TestFullLifecycle_ValidBlock_IsSuccessful() {
	const blockHeight = 1

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Prepare a block proposal.
	proposalTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          blockHeight,
		Time:            proposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Process the proposal.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          blockHeight,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Wait for the block to be finalized.
	time.Sleep(finalizeWaitDuration)

	// Finalize the block.
	finalizeResp, err := s.SimComet.Comet.FinalizeBlock(s.Ctx, &types.FinalizeBlockRequest{
		Txs:             proposal.Txs,
		Height:          blockHeight,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(finalizeResp)

	// Commit the block.
	_, err = s.SimComet.Comet.Commit(s.Ctx, &types.CommitRequest{})
	s.Require().NoError(err)

	// Validate post-commit state.
	queryCtx, err := s.SimComet.CreateQueryContext(blockHeight, false)
	s.Require().NoError(err)

	stateDB := s.TestNode.StorageBackend.StateFromContext(queryCtx)
	slot, err := stateDB.GetSlot()
	s.Require().NoError(err)
	s.Require().Equal(mathpkg.U64(blockHeight), slot)

	stateHeader, err := stateDB.GetLatestBlockHeader()
	s.Require().NoError(err)

	// Unmarshal the beacon block from the ABCI request.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForSlot(blockHeight),
	)
	s.Require().NoError(err)
	s.Require().Equal(proposedBlock.Message.GetHeader().GetBodyRoot(), stateHeader.GetBodyRoot())
}

// TestProcessProposal_BadBlock_IsRejected tests that a block with an invalid tx is rejected
func (s *SimulatedSuite) TestProcessProposal_BadBlock_IsRejected() {
	const blockHeight = 1

	// Initialize the chain state.
	s.initializeChain()

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Prepare a valid block proposal.
	proposalTime := time.Now()
	proposal, err := s.SimComet.Comet.PrepareProposal(s.Ctx, &types.PrepareProposalRequest{
		Height:          blockHeight,
		Time:            proposalTime,
		ProposerAddress: pubkey.Address(),
	})
	s.Require().NoError(err)
	s.Require().NotEmpty(proposal)

	// Unmarshal the proposal block.
	proposedBlock, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		proposal.Txs,
		blockchain.BeaconBlockTxIndex,
		s.TestNode.ChainSpec.ActiveForkVersionForSlot(blockHeight),
	)
	s.Require().NoError(err)

	// Create a malicious block by injecting an invalid transaction.
	maliciousBlock := simulated.CreateInvalidBlock(require.New(s.T()), proposedBlock, blsSigner, s.TestNode.ChainSpec, s.GenesisValidatorsRoot)
	maliciousBlockBytes, err := maliciousBlock.MarshalSSZ()
	s.Require().NoError(err)

	// Replace the valid block with the malicious block in the proposal.
	proposal.Txs[0] = maliciousBlockBytes

	// Process the proposal containing the malicious block.
	processResp, err := s.SimComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:             proposal.Txs,
		Height:          blockHeight,
		ProposerAddress: pubkey.Address(),
		Time:            proposalTime,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)

	// Verify that the log contains the expected error message.
	s.Require().Contains(s.LogBuffer.String(), errors.ErrInvalidPayloadStatus.Error())
	s.Require().Contains(s.LogBuffer.String(), "max fee per gas less than block base fee: address 0x71562b71999873DB5b286dF957af199Ec94617F7, maxFeePerGas: 10000000, baseFee: 875000000")
}
