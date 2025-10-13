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

// PayloadCacheSuite defines our test suite for Pectra related work using simulated Comet component.
type PayloadCacheSuite struct {
	suite.Suite
	Geth simulated.SharedAccessors
	Reth simulated.SharedAccessors
}

// TestPayloadCacheSuite runs the test suite.
func TestPayloadCacheSuite(t *testing.T) {
	suite.Run(t, new(PayloadCacheSuite))
}

// SetupTest initializes the test environment.
func (s *PayloadCacheSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.Geth.CtxApp, s.Geth.CtxAppCancelFn = context.WithCancel(context.Background())
	s.Reth.CtxApp, s.Reth.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.Geth.CtxComet = context.TODO()
	s.Geth.HomeDir = s.T().TempDir()

	s.Reth.CtxComet = context.TODO()
	s.Reth.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-fork-genesis.json"
	chainSpecFunc := simulated.ProvidePectraForkTestChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	gethCmtCfg, rethCmtCfg, genesisValidatorsRoot := simulated.Initialize2HomeDirs(
		s.T(), chainSpec, s.Geth.HomeDir, s.Reth.HomeDir, elGenesisPath,
	)
	s.Geth.GenesisValidatorsRoot = genesisValidatorsRoot
	s.Reth.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Geth node.
	gethNode := execution.NewGethNode(s.Geth.HomeDir, execution.ValidGethImage())
	elHandle, authRPC, elRPC := gethNode.Start(s.T(), path.Base(elGenesisPath))
	s.Geth.ElHandle = elHandle

	rethNode := execution.NewRethNode(s.Reth.HomeDir, execution.ValidRethImage())
	rethHandle, rethAuthRPC, elRPC := rethNode.Start(s.T(), path.Base(elGenesisPath))
	s.Reth.ElHandle = rethHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.Geth.LogBuffer = new(bytes.Buffer)
	logger := phuslu.NewLogger(s.Geth.LogBuffer, nil)

	s.Reth.LogBuffer = new(bytes.Buffer)
	rethLogger := phuslu.NewLogger(s.Reth.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component and electra genesis chain spec
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.Geth.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Geth.HomeDir,
		CometConfig: gethCmtCfg,
		AuthRPC:     authRPC,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Geth.SimComet = s.Geth.TestNode.SimComet

	s.Reth.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Reth.HomeDir,
		CometConfig: rethCmtCfg,
		AuthRPC:     rethAuthRPC,
		ClientRPC:   elRPC,
		Logger:      rethLogger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Reth.SimComet = s.Reth.TestNode.SimComet

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Geth.TestNode.Start(s.Geth.CtxApp)
	}()
	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Reth.TestNode.Start(s.Reth.CtxApp)
	}()

	s.Geth.SimulationClient = execution.NewSimulationClient(s.Geth.TestNode.EngineClient)
	// Reth does not have a simulation API
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.Geth.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
	err = simulated.WaitTillServicesStarted(s.Reth.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
}

// TearDownTest cleans up the test environment.
func (s *PayloadCacheSuite) TearDownTest() {
	// If the test has failed, log additional information.
	if s.T().Failed() {
		s.T().Log("GETH CL LOGS:")
		s.T().Log(s.Geth.LogBuffer.String())
		s.T().Log("RETH CL LOGS:")
		s.T().Log(s.Reth.LogBuffer.String())
	}
	if err := s.Geth.ElHandle.Close(); err != nil {
		s.T().Error("Error closing Geth EL handle:", err)
	}
	if err := s.Reth.ElHandle.Close(); err != nil {
		s.T().Error("Error closing Reth EL handle:", err)
	}
	// mimics the behaviour of shutdown func
	s.Geth.CtxAppCancelFn()
	s.Geth.TestNode.ServiceRegistry.StopAll()
	s.Reth.CtxAppCancelFn()
	s.Reth.TestNode.ServiceRegistry.StopAll()
}

// This tests a reth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
//
// Note: this test does NOT require --engine.always-process-payload-attributes-on-canonical-head.
func (s *PayloadCacheSuite) TestReth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain2Validators(s.T()) // 1 reth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Reth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from reth.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}

		// Process the proposal.
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a geth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PayloadCacheSuite) TestGeth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Geth.InitializeChain2Validators(s.T()) // 1 geth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(200 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())+1, 0)
	{
		// Prepare the proposal. Bkit cached the payload ID, so we just get the old one from geth.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     pubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: pubkey.Address(),
		}

		// Process the proposal.
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Geth.SimComet.Comet.FinalizeBlock(s.Geth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Geth.SimComet.Comet.Commit(s.Geth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a reth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
//
// Note: this test does NOT require --engine.always-process-payload-attributes-on-canonical-head.
func (s *PayloadCacheSuite) TestReth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain2Validators(s.T()) // 1 reth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Reth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            faultyConsensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		// As we reject our own built proposal, bkit should evict this payload from its cache.
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Reth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}

		// Process the proposal.
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a geth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PayloadCacheSuite) TestGeth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Geth.InitializeChain2Validators(s.T()) // 1 geth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            faultyConsensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		// As we reject our own built proposal, bkit should evict this payload from its cache.
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// Subsequent round where we are selected to propose again.
	{
		// Prepare the valid proposal. This should now request the EL for a new payload.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}

		// Process the proposal.
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		_, finalizeErr := s.Geth.SimComet.Comet.FinalizeBlock(s.Geth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Geth.SimComet.Comet.Commit(s.Geth.CtxComet, &types.CommitRequest{})
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
	s.Geth.InitializeChain2Validators(s.T()) // 1 geth validator
	s.Reth.InitializeChain2Validators(s.T()) // 1 reth validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	gethPubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)
	blsSigner = simulated.GetBlsSigner(s.Reth.HomeDir)
	rethPubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()-4), 0)
	{
		// Prepare the proposal.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethPubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     gethPubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: gethPubkey.Address(),
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Reth also prepares proposal.
		proposal, prepareErr = s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethPubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethPubkey.Address(),
			Time:                consensusTime,
			NextProposerAddress: rethPubkey.Address(),
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// We are now crossing over into post-fork time.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(100 * time.Millisecond) // This lets the optimistic build complete.
	consensusTime = time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)
	{
		// Geth tries to prepare proposal again by retrieving from cache, but post-fork this should fail.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethPubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with payload eviction from bkit cache.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: gethPubkey.Address(),
			Time:            consensusTime,
		}
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
		processRequest.ProposerAddress = rethPubkey.Address()
		processResp, respErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Reth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new (pre-fork) payload from geth EL.
		// NOTE: this will fail because geth does not allow re-building a payload for a height
		// that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethPubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Geth returns an empty proposal.
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new post-fork payload from reth EL. This works because the reth flag
		// allows us to rebuild a payload that has already been marked safe/finalized.
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethPubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: rethPubkey.Address(),
			Time:            consensusTime,
		}
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		processResp, respErr = s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
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
	s.Geth.InitializeChain2Validators(s.T())
	gethNodeAddress, err := s.Geth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth.InitializeChain2Validators(s.T())
	rethNodeAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// Both reth and geth prepare and propose a post-fork block without finalizing.
	{
		consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)

		// Geth builds.
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
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
		s.Geth.LogBuffer.Reset()
		processResp, respErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Reth also builds.
		proposal, prepareErr = s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
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
		s.Reth.LogBuffer.Reset()
		processResp, respErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	time.Sleep(50 * time.Millisecond) // Next round.
	// The proposer prepares a pre-fork block with finalization. The first pre-fork block it proposes will be rejected
	// As it will propose a post-fork block due to retrieving an Execution Payload in the PayloadCache.
	{
		// Get the same already built payload from cache.
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())-4, 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
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
			NextProposerAddress: rethNodeAddress,
		}

		// Process the proposal and rejects but does not bkit evict payload from cache because it didnt build it.
		s.Geth.LogBuffer.Reset()
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT.String(), processResp.Status.String())
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)

		// Reth also process proposal --> Trigger bkit eviction of payload from cache because its rejected.
		s.Reth.LogBuffer.Reset()
		processResp, processErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	time.Sleep(50 * time.Millisecond) // Next round.
	// Geth now proposes a post-fork block from its cache, which should also be invalid.
	{
		// Get the same already built payload from cache.
		consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())-3, 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     gethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: gethNodeAddress,
		}

		// Process the proposal and rejects but does not bkit evict payload from cache because it didnt build it.
		s.Reth.LogBuffer.Reset()
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)

		// Geth also process proposal --> Trigger bkit eviction of payload from cache because its rejected.
		s.Geth.LogBuffer.Reset()
		processResp, processErr = s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
		s.Require().Contains(
			s.Geth.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	// The previous payload in cache has been evicted. The optimistic builds for next height
	// should have completed by now.
	{
		// Try to build a new (pre-fork) payload from geth EL.
		// NOTE: this will fail because geth does not allow re-building a payload for a height
		// that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: gethNodeAddress,
		}
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Geth returns an empty proposal.
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	// The next block the proposer proposes with a pre-fork timestamp will actually have a pre-fork time
	// Since the previous payload in cache has been evicted, a new payload is built and retrieved.
	{
		// Force build a new (pre-fork) payload from reth EL.
		// NOTE: this requires --engine.always-process-payload-attributes-on-canonical-head.
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
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
		s.Geth.LogBuffer.Reset()
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Reth also process proposal and does not evict payload from bkit cache.
		s.Reth.LogBuffer.Reset()
		processResp, processErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block. Evict bkit payload here because finalize is accepted.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: rethNodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Geth.SimComet.Comet.FinalizeBlock(s.Geth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, finalizeErr = s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)

		// Commit the block.
		_, err := s.Geth.SimComet.Comet.Commit(s.Geth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Geth.LogBuffer.Reset()
		_, err = s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Reth.LogBuffer.Reset()
	}

	// Finally, we cross the fork and show no issues. Geth uses the optimistic build which has the
	// correct payload time and consequently is built correctly for post-fork.
	nextBlockHeight++
	time.Sleep(100 * time.Millisecond) // The optimistic build for next height should have completed by now.
	{
		consensusTime := time.Unix(int64(s.Geth.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
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
		processResp, processErr := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
		s.Require().Contains(s.Geth.LogBuffer.String(), "Processing execution requests")
		processResp, processErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		s.Require().Contains(s.Reth.LogBuffer.String(), "Processing execution requests")
	}
}
