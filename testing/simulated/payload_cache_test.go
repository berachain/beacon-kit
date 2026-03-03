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
	Primary   simulated.SharedAccessors
	Secondary simulated.SharedAccessors
}

// TestPayloadCacheSuite runs the test suite.
func TestPayloadCacheSuite(t *testing.T) {
	suite.Run(t, new(PayloadCacheSuite))
}

// SetupTest initializes the test environment.
func (s *PayloadCacheSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.Primary.CtxApp, s.Primary.CtxAppCancelFn = context.WithCancel(context.Background())
	s.Secondary.CtxApp, s.Secondary.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.Primary.CtxComet = context.TODO()
	s.Primary.HomeDir = s.T().TempDir()

	s.Secondary.CtxComet = context.TODO()
	s.Secondary.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-fork-genesis.json"
	chainSpecFunc := simulated.ProvidePectraForkTestChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	primaryCmtCfg, secondaryCmtCfg, genesisValidatorsRoot := simulated.Initialize2HomeDirs(
		s.T(), chainSpec, s.Primary.HomeDir, s.Secondary.HomeDir, elGenesisPath,
	)
	s.Primary.GenesisValidatorsRoot = genesisValidatorsRoot
	s.Secondary.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the primary EL (execution layer) Reth node.
	primaryNode := execution.NewRethNode(s.Primary.HomeDir, execution.ValidRethImage())
	elHandle, authRPC, elRPC := primaryNode.Start(s.T(), path.Base(elGenesisPath))
	s.Primary.ElHandle = elHandle

	// Choose the secondary reth node to run. 2 specific tests require the engine api override flag.
	var secondaryNode *execution.ExecNode
	testName := s.T().Name()
	if strings.Contains(testName, "TestReth_MustRebuildPostForkPayload_IsSuccessful") ||
		strings.Contains(testName, "TestReth_MustRebuildPreForkPayload_IsSuccessful") {
		secondaryNode = execution.NewRethNodeWithEngineOverride(s.Secondary.HomeDir, execution.ValidRethImage())
	} else {
		secondaryNode = execution.NewRethNode(s.Secondary.HomeDir, execution.ValidRethImage())
	}
	secondaryHandle, secondaryAuthRPC, elRPC := secondaryNode.Start(s.T(), path.Base(elGenesisPath))
	s.Secondary.ElHandle = secondaryHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.Primary.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.Primary.LogBuffer, nil)

	s.Secondary.LogBuffer = &simulated.SyncBuffer{}
	secondaryLogger := phuslu.NewLogger(s.Secondary.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component and electra genesis chain spec
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.Primary.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Primary.HomeDir,
		CometConfig: primaryCmtCfg,
		AuthRPC:     authRPC,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Primary.SimComet = s.Primary.TestNode.SimComet
	nodeAddress, err := s.Primary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Primary.SimComet.Comet.SetNodeAddress(nodeAddress)

	s.Secondary.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Secondary.HomeDir,
		CometConfig: secondaryCmtCfg,
		AuthRPC:     secondaryAuthRPC,
		ClientRPC:   elRPC,
		Logger:      secondaryLogger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Secondary.SimComet = s.Secondary.TestNode.SimComet
	nodeAddress, err = s.Secondary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Secondary.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Primary.TestNode.Start(s.Primary.CtxApp)
	}()
	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Secondary.TestNode.Start(s.Secondary.CtxApp)
	}()

	s.Primary.SimulationClient = execution.NewSimulationClient(s.Primary.TestNode.EngineClient)
	// Secondary node does not have a simulation API
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.Primary.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
	err = simulated.WaitTillServicesStarted(s.Secondary.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

// TearDownTest cleans up the test environment.
func (s *PayloadCacheSuite) TearDownTest() {
	s.Primary.CleanupTestWithLabel(s.T(), "PRIMARY")
	s.Secondary.CleanupTestWithLabel(s.T(), "SECONDARY")
}

// This tests a reth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PayloadCacheSuite) TestReth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Secondary.InitializeChain2Validators(s.T()) // 1 reth validator
	nodeAddress, err := s.Secondary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Secondary.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
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
		processResp, respErr := s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from reth.
		proposal, prepareErr := s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
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
		processResp, processErr := s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Secondary.SimComet.Comet.FinalizeBlock(s.Secondary.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Secondary.SimComet.Comet.Commit(s.Secondary.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a reth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PayloadCacheSuite) TestReth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Secondary.InitializeChain2Validators(s.T()) // 1 reth validator
	nodeAddress, err := s.Secondary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Secondary.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
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
		processResp, respErr := s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Secondary.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
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
		processResp, processErr := s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Secondary.SimComet.Comet.FinalizeBlock(s.Secondary.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Secondary.SimComet.Comet.Commit(s.Secondary.CtxComet, &types.CommitRequest{})
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
	s.Primary.InitializeChain2Validators(s.T()) // primary node
	primaryNodeAddress, err := s.Primary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Primary.SimComet.Comet.SetNodeAddress(primaryNodeAddress)
	s.Secondary.InitializeChain2Validators(s.T()) // secondary node (with engine override flag)
	secondaryNodeAddress, err := s.Secondary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Secondary.SimComet.Comet.SetNodeAddress(secondaryNodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime()-1), 0)
	{
		// Prepare the proposal.
		proposal, prepareErr := s.Primary.SimComet.Comet.PrepareProposal(s.Primary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: primaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     primaryNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: primaryNodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Primary.SimComet.Comet.ProcessProposal(s.Primary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Secondary node also prepares proposal.
		proposal, prepareErr = s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: secondaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     secondaryNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: secondaryNodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr = s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// We are now crossing over into post-fork time.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new payload from the primary node EL.
		// NOTE: this will fail because the primary node does not allow re-building a payload for a height
		// that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Primary.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Primary.SimComet.Comet.PrepareProposal(s.Primary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: primaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Primary node returns an empty proposal (no engine override flag).
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new post-fork payload from secondary node (with engine override flag) EL.
		// This works because the engine override flag allows us to rebuild a payload that has
		// already been marked safe/finalized.
		consensusTime := time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: secondaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     secondaryNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: secondaryNodeAddress,
		}
		processResp, respErr := s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		processResp, respErr = s.Primary.SimComet.Comet.ProcessProposal(s.Primary.CtxComet, processRequest)
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
	s.Primary.InitializeChain2Validators(s.T())
	primaryNodeAddress, err := s.Primary.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Secondary.InitializeChain2Validators(s.T())
	secondaryNodeAddress, err := s.Secondary.SimComet.GetNodeAddress()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// Both primary and secondary nodes prepare and propose a post-fork block without finalizing.
	{
		consensusTime := time.Unix(int64(s.Primary.TestNode.ChainSpec.ElectraForkTime()), 0)

		// Primary node builds.
		proposal, prepareErr := s.Primary.SimComet.Comet.PrepareProposal(s.Primary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: primaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Primary node processes the proposal. No bkit payload eviction here.
		// Optimistically build the next height's payload.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     primaryNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: primaryNodeAddress,
		}
		s.Primary.LogBuffer.Reset()
		processResp, respErr := s.Primary.SimComet.Comet.ProcessProposal(s.Primary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Secondary node also builds.
		proposal, prepareErr = s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: secondaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Secondary node processes the proposal. No bkit payload eviction here.
		// Optimistically build the next height's payload.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     secondaryNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: secondaryNodeAddress,
		}
		s.Secondary.LogBuffer.Reset()
		processResp, respErr = s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	time.Sleep(100 * time.Millisecond) // Next round.
	// The previous payload in cache has been evicted. The optimistic builds for next height
	// should have completed by now.
	{
		// Try to build a new (pre-fork) payload from the primary node EL.
		// NOTE: this will fail because the primary node (without engine override flag) does not allow
		// re-building a payload for a height that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Primary.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: primaryNodeAddress,
		}
		proposal, prepareErr := s.Primary.SimComet.Comet.PrepareProposal(s.Primary.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Primary node returns an empty proposal (no engine override flag).
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	// The next block the proposer proposes with a pre-fork timestamp will actually have a pre-fork time
	// Since the previous payload in cache has been evicted, a new payload is built and retrieved.
	{
		// Force build a new (pre-fork) payload from secondary node (with engine override flag) EL.
		// NOTE: this requires --engine.always-process-payload-attributes-on-canonical-head.
		consensusTime := time.Unix(int64(s.Secondary.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Secondary.SimComet.Comet.PrepareProposal(s.Secondary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: secondaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     secondaryNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: primaryNodeAddress,
		}

		// Process the proposal. No bkit payload eviction here from cache. Also trigger an optimistic
		// build for next height.
		s.Primary.LogBuffer.Reset()
		processResp, processErr := s.Primary.SimComet.Comet.ProcessProposal(s.Primary.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Secondary node also process proposal and does not evict payload from bkit cache.
		s.Secondary.LogBuffer.Reset()
		processResp, processErr = s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block. Evict bkit payload here because finalize is accepted.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: secondaryNodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Primary.SimComet.Comet.FinalizeBlock(s.Primary.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, finalizeErr = s.Secondary.SimComet.Comet.FinalizeBlock(s.Secondary.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)

		// Commit the block.
		_, err := s.Primary.SimComet.Comet.Commit(s.Primary.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Primary.LogBuffer.Reset()
		_, err = s.Secondary.SimComet.Comet.Commit(s.Secondary.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Secondary.LogBuffer.Reset()
	}

	// Finally, we cross the fork and show no issues. Primary node uses the optimistic build which has the
	// correct payload time and consequently is built correctly for post-fork.
	nextBlockHeight++
	time.Sleep(100 * time.Millisecond) // The optimistic build for next height should have completed by now.
	{
		consensusTime := time.Unix(int64(s.Primary.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Primary.SimComet.Comet.PrepareProposal(s.Primary.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: primaryNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: primaryNodeAddress,
			Time:            consensusTime,
		}
		// Process the proposal.
		processResp, processErr := s.Primary.SimComet.Comet.ProcessProposal(s.Primary.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
		s.Require().Contains(s.Primary.LogBuffer.String(), "Processing execution requests")
		processResp, processErr = s.Secondary.SimComet.Comet.ProcessProposal(s.Secondary.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		s.Require().Contains(s.Secondary.LogBuffer.String(), "Processing execution requests")
	}
}

// Test a scenario where reth must rebuild a payload for a failed state transition.
func (s *PectraForkSuite) TestReth_MustRebuildForFailedStateTransition_IsSuccessful() {
	// Initialize the chain state.
	testEL := s.Secondary
	helpBuilder := s.Primary
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
