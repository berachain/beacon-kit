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
	"context"
	"path"
	"strings"
	"testing"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// PayloadCacheSuite defines our test suite for Pectra related work using simulated Comet component.
type PayloadCacheSuite struct {
	suite.Suite
	Node1 simulated.SharedAccessors
	Node2 simulated.SharedAccessors
}

// TestPayloadCacheSuite runs the test suite.
func TestPayloadCacheSuite(t *testing.T) {
	suite.Run(t, new(PayloadCacheSuite))
}

// SetupTest initializes the test environment.
func (s *PayloadCacheSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.Node1.CtxApp, s.Node1.CtxAppCancelFn = context.WithCancel(context.Background())
	s.Node2.CtxApp, s.Node2.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.Node1.CtxComet = context.TODO()
	s.Node1.HomeDir = s.T().TempDir()

	s.Node2.CtxComet = context.TODO()
	s.Node2.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-fork-genesis.json"
	chainSpecFunc := simulated.ProvidePectraForkTestChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	gethCmtCfg, rethCmtCfg, genesisValidatorsRoot := simulated.Initialize2HomeDirs(
		s.T(), chainSpec, s.Node1.HomeDir, s.Node2.HomeDir, elGenesisPath,
	)
	s.Node1.GenesisValidatorsRoot = genesisValidatorsRoot
	s.Node2.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL nodes. 2 specific tests require the engine api override flag.
	var (
		rethNode1 *execution.ExecNode
		rethNode2 *execution.ExecNode
	)
	testName := s.T().Name()
	if strings.Contains(testName, "TestReth_MustRebuildPostForkPayload_IsSuccessful") ||
		strings.Contains(testName, "TestReth_MustRebuildPreForkPayload_IsSuccessful") {
		rethNode1 = execution.NewRethNodeWithEngineOverride(s.Node1.HomeDir, execution.ValidRethImage())
		rethNode2 = execution.NewRethNodeWithEngineOverride(s.Node2.HomeDir, execution.ValidRethImage())
	} else {
		rethNode1 = execution.NewRethNode(s.Node1.HomeDir, execution.ValidRethImage())
		rethNode2 = execution.NewRethNode(s.Node2.HomeDir, execution.ValidRethImage())
	}
	rethHandle1, rethAuthRPC1, rethRPC1 := rethNode1.Start(s.T(), path.Base(elGenesisPath))
	rethHandle2, rethAuthRPC2, rethRPC2 := rethNode2.Start(s.T(), path.Base(elGenesisPath))
	s.Node1.ElHandle = rethHandle1
	s.Node2.ElHandle = rethHandle2

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.Node1.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.Node1.LogBuffer, nil)

	s.Node2.LogBuffer = &simulated.SyncBuffer{}
	rethLogger := phuslu.NewLogger(s.Node2.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component and electra genesis chain spec
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.Node1.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Node1.HomeDir,
		CometConfig: gethCmtCfg,
		AuthRPC:     rethAuthRPC1,
		ClientRPC:   rethRPC1,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Node1.SimComet = s.Node1.TestNode.SimComet
	nodeAddress, err := s.Node1.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node1.SimComet.Comet.SetNodeAddress(nodeAddress)

	s.Node2.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Node2.HomeDir,
		CometConfig: rethCmtCfg,
		AuthRPC:     rethAuthRPC2,
		ClientRPC:   rethRPC2,
		Logger:      rethLogger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Node2.SimComet = s.Node2.TestNode.SimComet
	nodeAddress, err = s.Node2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node2.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Node1.TestNode.Start(s.Node1.CtxApp)
	}()
	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Node2.TestNode.Start(s.Node2.CtxApp)
	}()

	s.Node1.SimulationClient = execution.NewSimulationClient(s.Node1.TestNode.EngineClient)
	// Reth does not have a simulation API
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.Node1.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
	err = simulated.WaitTillServicesStarted(s.Node2.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

// TearDownTest cleans up the test environment.
func (s *PayloadCacheSuite) TearDownTest() {
	s.Node1.CleanupTestWithLabel(s.T(), "Node1")
	s.Node2.CleanupTestWithLabel(s.T(), "Node2")
}

// This tests a reth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PayloadCacheSuite) TestReth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Node2.InitializeChain2Validators(s.T()) // 1 reth validator
	nodeAddress, err := s.Node2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node2.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from reth.
		proposal, prepareErr := s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal.
		processResp, processErr := s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Node2.SimComet.Comet.FinalizeBlock(s.Node2.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Node2.SimComet.Comet.Commit(s.Node2.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a geth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PayloadCacheSuite) TestGeth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Node1.InitializeChain2Validators(s.T()) // 1 geth validator
	nodeAddress, err := s.Node1.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node1.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from geth.
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal.
		processResp, processErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Node1.SimComet.Comet.FinalizeBlock(s.Node1.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Node1.SimComet.Comet.Commit(s.Node1.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a reth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PayloadCacheSuite) TestReth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Node2.InitializeChain2Validators(s.T()) // 1 reth validator
	nodeAddress, err := s.Node2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node2.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            faultyConsensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		// As we reject our own built proposal, bkit should evict this payload from its cache.
		processResp, respErr := s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Node2.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal.
		processResp, processErr := s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Node2.SimComet.Comet.FinalizeBlock(s.Node2.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Node2.SimComet.Comet.Commit(s.Node2.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a geth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PayloadCacheSuite) TestGeth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Node1.InitializeChain2Validators(s.T()) // 1 geth validator
	nodeAddress, err := s.Node1.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node1.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            faultyConsensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		// As we reject our own built proposal, bkit should evict this payload from its cache.
		processResp, respErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Node1.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal.
		processResp, processErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Node1.SimComet.Comet.FinalizeBlock(s.Node1.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Node1.SimComet.Comet.Commit(s.Node1.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// Test a scenario where the first proposed block is pre-fork and accepted initially but never
// finalized by the network (which triggeres optimistic builds). The subsequent rounds are now
// post-fork. Only reth with the flag can force rebuild a payload.
//
// NOTE: this test requires reth with the --engine.always-process-payload-attributes-on-canonical-head flag.
func (s *PayloadCacheSuite) TestReth_MustRebuildPostForkPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Node1.InitializeChain2Validators(s.T()) // 1 geth validator
	gethNodeAddress, err := s.Node1.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node1.SimComet.Comet.SetNodeAddress(gethNodeAddress)
	s.Node2.InitializeChain2Validators(s.T()) // 1 reth validator
	rethNodeAddress, err := s.Node2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node2.SimComet.Comet.SetNodeAddress(rethNodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime()-1), 0)
	{
		// Prepare the proposal.
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     gethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: gethNodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Reth also prepares proposal.
		proposal, prepareErr = s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: rethNodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr = s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// We are now crossing over into post-fork time.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new (pre-fork) payload from geth EL.
		// NOTE: this will fail because geth does not allow re-building a payload for a height
		// that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Geth returns an empty proposal.
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new post-fork payload from reth EL. This works because the reth flag
		// allows us to rebuild a payload that has already been marked safe/finalized.
		consensusTime := time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: rethNodeAddress,
		}
		processResp, respErr := s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		processResp, respErr = s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}
}

// Test a scenario where the first proposed block is post-fork and accepted initially but never
// finalized by the network (which triggeres optimistic builds). The subsequent rounds are actually
// pre-fork. Only reth with the flag can force rebuild a payload.
//
// NOTE: this test requires reth with the --engine.always-process-payload-attributes-on-canonical-head flag
// to propose the valid pre-fork block.
func (s *PayloadCacheSuite) TestReth_MustRebuildPreForkPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Node1.InitializeChain2Validators(s.T())
	gethNodeAddress, err := s.Node1.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Node2.InitializeChain2Validators(s.T())
	rethNodeAddress, err := s.Node2.SimComet.GetNodeAddress()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// Both reth and geth prepare and propose a post-fork block without finalizing.
	{
		consensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime()), 0)

		// Geth builds.
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Geth processes the proposal. No bkit payload eviction here.
		// Optimistically build the next height's payload.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     gethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: gethNodeAddress,
		}
		s.Node1.LogBuffer.Reset()
		processResp, respErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Reth also builds.
		proposal, prepareErr = s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Reth processes the proposal. No bkit payload eviction here.
		// Optimistically build the next height's payload.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: rethNodeAddress,
		}
		s.Node2.LogBuffer.Reset()
		processResp, respErr = s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	time.Sleep(100 * time.Millisecond) // Next round.
	// The previous payload in cache has been evicted. The optimistic builds for next height
	// should have completed by now.
	{
		// Try to build a new (pre-fork) payload from geth EL.
		// NOTE: this will fail because geth does not allow re-building a payload for a height
		// that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		}
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Geth returns an empty proposal.
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	// The next block the proposer proposes with a pre-fork timestamp will actually have a pre-fork time
	// Since the previous payload in cache has been evicted, a new payload is built and retrieved.
	{
		// Force build a new (pre-fork) payload from reth EL.
		// NOTE: this requires --engine.always-process-payload-attributes-on-canonical-head.
		consensusTime := time.Unix(int64(s.Node2.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Node2.SimComet.Comet.PrepareProposal(s.Node2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: gethNodeAddress,
		}

		// Process the proposal. No bkit payload eviction here from cache. Also trigger an optimistic
		// build for next height.
		s.Node1.LogBuffer.Reset()
		processResp, processErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Reth also process proposal and does not evict payload from bkit cache.
		s.Node2.LogBuffer.Reset()
		processResp, processErr = s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block. Evict bkit payload here because finalize is accepted.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: rethNodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Node1.SimComet.Comet.FinalizeBlock(s.Node1.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, finalizeErr = s.Node2.SimComet.Comet.FinalizeBlock(s.Node2.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)

		// Commit the block.
		_, err := s.Node1.SimComet.Comet.Commit(s.Node1.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Node1.LogBuffer.Reset()
		_, err = s.Node2.SimComet.Comet.Commit(s.Node2.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Node2.LogBuffer.Reset()
	}

	// Finally, we cross the fork and show no issues. Geth uses the optimistic build which has the
	// correct payload time and consequently is built correctly for post-fork.
	nextBlockHeight++
	time.Sleep(100 * time.Millisecond) // The optimistic build for next height should have completed by now.
	{
		consensusTime := time.Unix(int64(s.Node1.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Node1.SimComet.Comet.PrepareProposal(s.Node1.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: gethNodeAddress,
			Time:            consensusTime,
		}
		// Process the proposal.
		processResp, processErr := s.Node1.SimComet.Comet.ProcessProposal(s.Node1.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
		s.Require().Contains(s.Node1.LogBuffer.String(), "Processing execution requests")
		processResp, processErr = s.Node2.SimComet.Comet.ProcessProposal(s.Node2.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		s.Require().Contains(s.Node2.LogBuffer.String(), "Processing execution requests")
	}
}

// Test a scenario where reth must rebuild a payload for a failed state transition.
func (s *PectraForkSuite) TestReth_MustRebuildForFailedStateTransition_IsSuccessful() {
	// Initialize the chain state.
	testEL := s.Node2
	helpBuilder := s.Node1
	testEL.InitializeChain(s.T())
	helpBuilder.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(testEL.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	nodeAddress := pubkey.Address()
	testEL.SimComet.Comet.SetNodeAddress(nodeAddress)

	const blkHeight = int64(1)
	var (
		specs         = testEL.TestNode.ChainSpec
		consensusTime = time.Unix(int64(specs.ElectraForkTime()), 0)

		validTxsHeight1 [][]byte
	)

	{
		// 1- Build a valid block at height 1, via the helpBuilder
		prepareRequest := &types.PrepareProposalRequest{
			Height:          blkHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		}
		proposal, prepareErr := helpBuilder.SimComet.Comet.PrepareProposal(helpBuilder.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)
		validTxsHeight1 = proposal.Txs

		// 2- Process the block via testEL node. The proposal is expected
		// to pass and start building payload for height 2, optimistically.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              blkHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		processResp, respErr := testEL.SimComet.Comet.ProcessProposal(testEL.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	// For some reason, the supermajority does not finalize the block.
	// Another block comes, still at height 1, this time *invalid*. This would force
	// the node to rebuild height 1, which the EL cannot do since it has already received
	// an FCU(head == block_at_height_2)
	{
		invalidTxs := testBuildInvalidBlock(
			s.Require(),
			helpBuilder,
			&types.PrepareProposalRequest{
				Txs:             validTxsHeight1,
				Height:          blkHeight,
				Time:            consensusTime,
				ProposerAddress: pubkey.Address(),
			},
			func(sbb *ctypes.SignedBeaconBlock) {
				sbb.Body.RandaoReveal = [96]byte{'t', 'e', 's', 't'} // this makes the block invalid
			},
		)

		// 3- Process the invalid proposal proposal. It will be rejected
		// and attempt to build optimistically a block at height 1.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 invalidTxs,
			Height:              blkHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		processResp, processErr := testEL.SimComet.Comet.ProcessProposal(testEL.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT.String(), processResp.Status.String())
	}

	{
		// 4- Finally let reth node build block at height 1, process and finalize it
		prepareRequest := &types.PrepareProposalRequest{
			Height:          blkHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := testEL.SimComet.Comet.PrepareProposal(testEL.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal via testEL node. The proposal is expected
		// to pass and start building payload for height 2, optimistically.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              blkHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		processResp, respErr := testEL.SimComet.Comet.ProcessProposal(testEL.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          blkHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := testEL.SimComet.Comet.FinalizeBlock(testEL.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := testEL.SimComet.Comet.Commit(testEL.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}
