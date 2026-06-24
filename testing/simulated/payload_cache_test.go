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
	"math/big"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	svcencoding "github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/net/url"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	cmtcfg "github.com/cometbft/cometbft/config"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/holiman/uint256"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// PayloadCacheSuite defines our test suite for Pectra related work using simulated Comet component.
type PayloadCacheSuite struct {
	suite.Suite
	Reth  simulated.SharedAccessors
	Reth2 simulated.SharedAccessors
	Reth3 simulated.SharedAccessors
}

// TestPayloadCacheSuite runs the test suite.
func TestPayloadCacheSuite(t *testing.T) {
	suite.Run(t, new(PayloadCacheSuite))
}

// SetupTest initializes the test environment.
func (s *PayloadCacheSuite) SetupTest() {
	testName := s.T().Name()
	useThirdValidatorSetup := strings.Contains(
		testName,
		"TestReth_MisorderedBlobSidecarsCachedEnvelope_IsSuccessful",
	)

	// Create a cancellable context for the duration of the test.
	s.Reth.CtxApp, s.Reth.CtxAppCancelFn = context.WithCancel(context.Background())
	s.Reth2.CtxApp, s.Reth2.CtxAppCancelFn = context.WithCancel(context.Background())
	if useThirdValidatorSetup {
		s.Reth3.CtxApp, s.Reth3.CtxAppCancelFn = context.WithCancel(context.Background())
	}

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.Reth.CtxComet = context.TODO()
	s.Reth.HomeDir = s.T().TempDir()

	s.Reth2.CtxComet = context.TODO()
	s.Reth2.HomeDir = s.T().TempDir()

	if useThirdValidatorSetup {
		s.Reth3.CtxComet = context.TODO()
		s.Reth3.HomeDir = s.T().TempDir()
	}

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-fork-genesis.json"
	chainSpecFunc := simulated.ProvidePectraForkTestChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	var (
		rethCmtCfg  *cmtcfg.Config
		reth2CmtCfg *cmtcfg.Config
		reth3CmtCfg *cmtcfg.Config
	)
	if useThirdValidatorSetup {
		configs, root := simulated.InitializeHomeDirs(
			s.T(), chainSpec, elGenesisPath, s.Reth.HomeDir, s.Reth2.HomeDir, s.Reth3.HomeDir,
		)
		rethCmtCfg, reth2CmtCfg, reth3CmtCfg = configs[0], configs[1], configs[2]
		s.Reth.GenesisValidatorsRoot = root
		s.Reth2.GenesisValidatorsRoot = root
		s.Reth3.GenesisValidatorsRoot = root
	} else {
		configs, root := simulated.InitializeHomeDirs(
			s.T(), chainSpec, elGenesisPath, s.Reth.HomeDir, s.Reth2.HomeDir,
		)
		rethCmtCfg, reth2CmtCfg = configs[0], configs[1]
		s.Reth.GenesisValidatorsRoot = root
		s.Reth2.GenesisValidatorsRoot = root
	}

	// Start the EL (execution layer) Reth node.
	rethNode := execution.NewRethNode(s.Reth.HomeDir, execution.ValidRethImage())
	elHandle, authRPC, elRPC := rethNode.Start(s.T(), path.Base(elGenesisPath))
	s.Reth.ElHandle = elHandle

	// Choose the reth node to run. 2 specific tests require the engine api override flag.
	var (
		reth2Node *execution.ExecNode
		reth3Node *execution.ExecNode
	)
	if strings.Contains(testName, "TestReth_MustRebuildPostForkPayload_IsSuccessful") ||
		strings.Contains(testName, "TestReth_MustRebuildPreForkPayload_IsSuccessful") {
		reth2Node = execution.NewRethNodeWithEngineOverride(s.Reth2.HomeDir, execution.ValidRethImage())
		if useThirdValidatorSetup {
			reth3Node = execution.NewRethNodeWithEngineOverride(s.Reth3.HomeDir, execution.ValidRethImage())
		}
	} else {
		reth2Node = execution.NewRethNode(s.Reth2.HomeDir, execution.ValidRethImage())
		if useThirdValidatorSetup {
			reth3Node = execution.NewRethNode(s.Reth3.HomeDir, execution.ValidRethImage())
		}
	}
	rethHandle, rethAuthRPC, rethRPC := reth2Node.Start(s.T(), path.Base(elGenesisPath))
	s.Reth2.ElHandle = rethHandle
	var (
		reth2Handle  *execution.Resource
		reth2AuthRPC *url.ConnectionURL
		reth2RPC     *url.ConnectionURL
	)
	if useThirdValidatorSetup {
		reth2Handle, reth2AuthRPC, reth2RPC = reth3Node.Start(s.T(), path.Base(elGenesisPath))
		s.Reth3.ElHandle = reth2Handle
	}

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.Reth.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.Reth.LogBuffer, nil)

	s.Reth2.LogBuffer = &simulated.SyncBuffer{}
	rethLogger := phuslu.NewLogger(s.Reth2.LogBuffer, nil)

	var reth2Logger *phuslu.Logger
	if useThirdValidatorSetup {
		s.Reth3.LogBuffer = &simulated.SyncBuffer{}
		reth2Logger = phuslu.NewLogger(s.Reth3.LogBuffer, nil)
	}

	// Build the Beacon node with the simulated Comet component and electra genesis chain spec
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.Reth.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Reth.HomeDir,
		CometConfig: rethCmtCfg,
		AuthRPC:     authRPC,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Reth.SimComet = s.Reth.TestNode.SimComet
	nodeAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth.SimComet.Comet.SetNodeAddress(nodeAddress)

	s.Reth2.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Reth2.HomeDir,
		CometConfig: reth2CmtCfg,
		AuthRPC:     rethAuthRPC,
		ClientRPC:   rethRPC,
		Logger:      rethLogger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Reth2.SimComet = s.Reth2.TestNode.SimComet
	nodeAddress, err = s.Reth2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth2.SimComet.Comet.SetNodeAddress(nodeAddress)

	if useThirdValidatorSetup {
		s.Reth3.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
			TempHomeDir: s.Reth3.HomeDir,
			CometConfig: reth3CmtCfg,
			AuthRPC:     reth2AuthRPC,
			ClientRPC:   reth2RPC,
			Logger:      reth2Logger,
			AppOpts:     viper.New(),
			Components:  components,
		})
		s.Reth3.SimComet = s.Reth3.TestNode.SimComet
		nodeAddress, err = s.Reth3.SimComet.GetNodeAddress()
		s.Require().NoError(err)
		s.Reth3.SimComet.Comet.SetNodeAddress(nodeAddress)
	}

	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Reth.TestNode.Start(s.Reth.CtxApp)
	}()
	// Start the Beacon node in a separate goroutine.
	go func() {
		_ = s.Reth2.TestNode.Start(s.Reth2.CtxApp)
	}()
	if useThirdValidatorSetup {
		go func() {
			_ = s.Reth3.TestNode.Start(s.Reth3.CtxApp)
		}()
	}

	s.Reth.SimulationClient = execution.NewSimulationClient(s.Reth.TestNode.ContractBackend)
	timeOut := 10 * time.Second
	interval := 50 * time.Millisecond
	err = simulated.WaitTillServicesStarted(s.Reth.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
	err = simulated.WaitTillServicesStarted(s.Reth2.LogBuffer, timeOut, interval)
	s.Require().NoError(err)
	if useThirdValidatorSetup {
		err = simulated.WaitTillServicesStarted(s.Reth3.LogBuffer, timeOut, interval)
		s.Require().NoError(err)
	}
}

// TearDownTest cleans up the test environment.
func (s *PayloadCacheSuite) TearDownTest() {
	s.Reth.CleanupTestWithLabel(s.T(), "RETH")
	s.Reth2.CleanupTestWithLabel(s.T(), "RETH2")
	s.Reth3.CleanupTestWithLabel(s.T(), "RETH3")
}

// This tests a reth validator proposing a block. It then accepts the proposal in
// process proposal. But the block is not finalized by consensus. Then this
// validator is chosen to propose at a subsequent round. It should just get the old
// payload from its cache.
func (s *PayloadCacheSuite) TestReth_ReusePayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain(s.T(), 2) // 1 reth validator
	nodeAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare the proposal.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
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
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}
}

// This tests a reth validator proposing a invalid block. The proposal is rejected. Then this
// validator is chosen to propose at a subsequent round. It should now be forced to
// rebuild a new payload (and not reuse the old one from its cache).
func (s *PayloadCacheSuite) TestReth_RebuildPayload_IsSuccessful() {
	// Initialize the chain state.
	s.Reth.InitializeChain(s.T(), 2) // 1 reth validator
	nodeAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth.SimComet.Comet.SetNodeAddress(nodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

	{
		// Prepare an invalid proposal.
		faultyConsensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
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
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Now the block is finalized and committed.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
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
	s.Reth.InitializeChain(s.T(), 2) // 1 reth validator
	rethNodeAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth.SimComet.Comet.SetNodeAddress(rethNodeAddress)
	s.Reth2.InitializeChain(s.T(), 2) // 1 reth validator
	reth2NodeAddress, err := s.Reth2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth2.SimComet.Comet.SetNodeAddress(reth2NodeAddress)

	// Next block is height 1.
	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth2.TestNode.ChainSpec.ElectraForkTime()-1), 0)
	{
		// Prepare the proposal.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: rethNodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Reth also prepares proposal.
		proposal, prepareErr = s.Reth2.SimComet.Comet.PrepareProposal(s.Reth2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: reth2NodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal, with no payload eviction from bkit cache.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     reth2NodeAddress,
			Time:                consensusTime,
			NextProposerAddress: reth2NodeAddress,
		}
		// This will trigger a optimistic payload build for block height 2.
		processResp, respErr = s.Reth2.SimComet.Comet.ProcessProposal(s.Reth2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	}

	// For some reason, the supermajority does not finalize the block.
	// We are now crossing over into post-fork time.
	// Next round is height 1, but simulating consensus time is 1 second after previous round.
	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new (pre-fork) payload from reth EL.
		// NOTE: this will fail because reth (without --engine.always-process-payload-attributes-on-canonical-head
		// flag) does not allow re-building a payload for a height that has already been marked safe/finalized.
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Reth returns an empty proposal.
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	{
		// Try to build a new post-fork payload from reth EL. This works because the reth flag
		// allows us to rebuild a payload that has already been marked safe/finalized.
		consensusTime := time.Unix(int64(s.Reth2.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Reth2.SimComet.Comet.PrepareProposal(s.Reth2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: reth2NodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     reth2NodeAddress,
			Time:                consensusTime,
			NextProposerAddress: reth2NodeAddress,
		}
		processResp, respErr := s.Reth2.SimComet.Comet.ProcessProposal(s.Reth2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		processResp, respErr = s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
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
	s.Reth.InitializeChain(s.T(), 2)
	rethNodeAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth2.InitializeChain(s.T(), 2)
	reth2NodeAddress, err := s.Reth2.SimComet.GetNodeAddress()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// Both nodes prepare and propose a post-fork block without finalizing.
	{
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)

		// Reth builds.
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Reth processes the proposal. No bkit payload eviction here.
		// Optimistically build the next height's payload.
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     rethNodeAddress,
			Time:                consensusTime,
			NextProposerAddress: rethNodeAddress,
		}
		s.Reth.LogBuffer.Reset()
		processResp, respErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Reth also builds.
		proposal, prepareErr = s.Reth2.SimComet.Comet.PrepareProposal(s.Reth2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: reth2NodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// Reth processes the proposal. No bkit payload eviction here.
		// Optimistically build the next height's payload.
		processRequest = &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     reth2NodeAddress,
			Time:                consensusTime,
			NextProposerAddress: reth2NodeAddress,
		}
		s.Reth2.LogBuffer.Reset()
		processResp, respErr = s.Reth2.SimComet.Comet.ProcessProposal(s.Reth2.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	time.Sleep(100 * time.Millisecond) // Next round.
	// The previous payload in cache has been evicted. The optimistic builds for next height
	// should have completed by now.
	{
		// Try to build a new (pre-fork) payload from reth EL.
		// NOTE: this will fail because reth does not allow re-building a payload for a height
		// that has already been marked safe/finalized
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		}
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 0) // Reth returns an empty proposal.
	}

	time.Sleep(10 * time.Millisecond) // Next round.
	// The next block the proposer proposes with a pre-fork timestamp will actually have a pre-fork time
	// Since the previous payload in cache has been evicted, a new payload is built and retrieved.
	{
		// Force build a new (pre-fork) payload from reth EL.
		// NOTE: this requires --engine.always-process-payload-attributes-on-canonical-head.
		consensusTime := time.Unix(int64(s.Reth2.TestNode.ChainSpec.ElectraForkTime())-1, 0)
		proposal, prepareErr := s.Reth2.SimComet.Comet.PrepareProposal(s.Reth2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: reth2NodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     reth2NodeAddress,
			Time:                consensusTime,
			NextProposerAddress: rethNodeAddress,
		}

		// Process the proposal. No bkit payload eviction here from cache. Also trigger an optimistic
		// build for next height.
		s.Reth.LogBuffer.Reset()
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Reth also process proposal and does not evict payload from bkit cache.
		s.Reth2.LogBuffer.Reset()
		processResp, processErr = s.Reth2.SimComet.Comet.ProcessProposal(s.Reth2.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

		// Finalize the block. Evict bkit payload here because finalize is accepted.
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: reth2NodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := s.Reth.SimComet.Comet.FinalizeBlock(s.Reth.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, finalizeErr = s.Reth2.SimComet.Comet.FinalizeBlock(s.Reth2.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)

		// Commit the block.
		_, err := s.Reth.SimComet.Comet.Commit(s.Reth.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Reth.LogBuffer.Reset()
		_, err = s.Reth2.SimComet.Comet.Commit(s.Reth2.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)
		s.Reth2.LogBuffer.Reset()
	}

	// Finally, we cross the fork and show no issues. Reth uses the optimistic build which has the
	// correct payload time and consequently is built correctly for post-fork.
	nextBlockHeight++
	time.Sleep(100 * time.Millisecond) // The optimistic build for next height should have completed by now.
	{
		consensusTime := time.Unix(int64(s.Reth.TestNode.ChainSpec.ElectraForkTime()), 0)
		proposal, prepareErr := s.Reth.SimComet.Comet.PrepareProposal(s.Reth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: rethNodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: rethNodeAddress,
			Time:            consensusTime,
		}
		// Process the proposal.
		processResp, processErr := s.Reth.SimComet.Comet.ProcessProposal(s.Reth.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
		s.Require().Contains(s.Reth.LogBuffer.String(), "Processing execution requests")
		processResp, processErr = s.Reth2.SimComet.Comet.ProcessProposal(s.Reth2.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
		s.Require().Contains(s.Reth2.LogBuffer.String(), "Processing execution requests")
	}
}

// This tests a malicious proposer that reorders blob sidecars in a proposal.
// Another validator accepts and caches the payload envelope without finalizing.
// In a subsequent round at the same height, the caching validator proposes again
// from cache and a third validator accepts that proposal.
func (s *PayloadCacheSuite) TestReth_MisorderedBlobSidecarsCachedEnvelope_IsSuccessful() {
	// Initialize chain state with 3 validators.
	s.Reth.InitializeChain(s.T(), 3)
	s.Reth2.InitializeChain(s.T(), 3)
	s.Reth3.InitializeChain(s.T(), 3)

	// validator A: malicious proposer.
	maliciousProposerAddress, err := s.Reth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth.SimComet.Comet.SetNodeAddress(maliciousProposerAddress)

	// validator B: verifies malicious proposal and later proposes from cache.
	cachingValidatorAddress, err := s.Reth2.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth2.SimComet.Comet.SetNodeAddress(cachingValidatorAddress)

	// validator C: verifies validator B's subsequent proposal.
	verifyingValidatorAddress, err := s.Reth3.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Reth3.SimComet.Comet.SetNodeAddress(verifyingValidatorAddress)

	nextBlockHeight := int64(1)
	consensusTime := time.Unix(int64(s.Reth2.TestNode.ChainSpec.ElectraForkTime()), 0)

	// Seed proposer A's EL with blob txs so the proposal has reorderable sidecars.
	s.submitBlobTransactions(s.Reth, 2)

	var (
		proposal    *types.PrepareProposalResponse
		reorderedTx bool
	)
	for i := 0; i < 10; i++ {
		proposal, err = s.Reth.SimComet.Comet.PrepareProposal(
			s.Reth.CtxComet, &types.PrepareProposalRequest{
				Height:          nextBlockHeight,
				Time:            consensusTime,
				ProposerAddress: maliciousProposerAddress,
			},
		)
		s.Require().NoError(err)
		s.Require().Len(proposal.Txs, 2)

		sidecars, scErr := svcencoding.UnmarshalBlobSidecarsFromABCIRequest(
			proposal.Txs,
			blockchain.BlobSidecarsTxIndex,
		)
		s.Require().NoError(scErr)
		if len(sidecars) < 2 {
			time.Sleep(200 * time.Millisecond)
			continue // NOTE: we retry at most 10 times, with 200ms delay between retries.
		}

		// Maliciously reorder two sidecars in the proposal.
		sidecars[0], sidecars[1] = sidecars[1], sidecars[0]
		sidecarBytes, marshalErr := sidecars.MarshalSSZ()
		s.Require().NoError(marshalErr)
		proposal.Txs[blockchain.BlobSidecarsTxIndex] = sidecarBytes
		reorderedTx = true
		break
	}
	s.Require().True(
		reorderedTx,
		"expected at least 2 blob sidecars so the malicious proposer can reorder them",
	)

	// validator B verifies and accepts the maliciously-ordered proposal. This caches the
	// payload envelope, but the block never finalizes.
	processResp, err := s.Reth2.SimComet.Comet.ProcessProposal(
		s.Reth2.CtxComet, &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     maliciousProposerAddress,
			Time:                consensusTime,
			NextProposerAddress: cachingValidatorAddress,
		},
	)
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Next round at same height; validator B now proposes from cached envelope.
	time.Sleep(200 * time.Millisecond)
	consensusTime = consensusTime.Add(time.Second)
	cachedProposal, err := s.Reth2.SimComet.Comet.PrepareProposal(
		s.Reth2.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: cachingValidatorAddress,
		},
	)
	s.Require().NoError(err)
	s.Require().Len(cachedProposal.Txs, 2)

	// Ensure the proposal from cache carries sidecars in canonical index order.
	cachedSidecars, err := svcencoding.UnmarshalBlobSidecarsFromABCIRequest(
		cachedProposal.Txs,
		blockchain.BlobSidecarsTxIndex,
	)
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(len(cachedSidecars), 2)
	for i := 1; i < len(cachedSidecars); i++ {
		s.Require().GreaterOrEqual(
			cachedSidecars[i].GetIndex(),
			cachedSidecars[i-1].GetIndex(),
		)
	}

	// validator C verifies the cached proposal and accepts it.
	processResp, err = s.Reth3.SimComet.Comet.ProcessProposal(
		s.Reth3.CtxComet, &types.ProcessProposalRequest{
			Txs:                 cachedProposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     cachingValidatorAddress,
			Time:                consensusTime,
			NextProposerAddress: verifyingValidatorAddress,
		},
	)
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
}

func (s *PayloadCacheSuite) submitBlobTransactions(node simulated.SharedAccessors, numBlobs int) {
	s.T().Helper()
	s.Require().GreaterOrEqual(numBlobs, 2)

	blobs := make([]*eip4844.Blob, numBlobs)
	for i := range blobs {
		blob := &eip4844.Blob{}
		blob[0] = byte(i + 1)
		blob[1] = byte(i + 2)
		blobs[i] = blob
	}

	proofs, commitments := simulated.GetProofAndCommitmentsForBlobs(
		require.New(s.T()),
		blobs,
		node.TestNode.KZGVerifier,
	)

	testKey := simulated.GetTestKey(s.T())
	chainID := node.TestNode.ChainSpec.DepositEth1ChainID()
	signer := gethtypes.NewCancunSigner(big.NewInt(int64(chainID)))
	senderAddress := crypto.PubkeyToAddress(testKey.PublicKey)
	recipientAddress := gethcommon.HexToAddress(simulated.WithdrawalExecutionAddress)

	nextNonce, err := node.TestNode.ContractBackend.PendingNonceAt(node.CtxApp, senderAddress)
	s.Require().NoError(err)

	for i := 0; i < numBlobs; i++ {
		txSidecar := &gethtypes.BlobTxSidecar{
			Blobs:       []kzg4844.Blob{kzg4844.Blob(blobs[i][:])},
			Commitments: []kzg4844.Commitment{kzg4844.Commitment(commitments[i])},
			Proofs:      []kzg4844.Proof{kzg4844.Proof(proofs[i])},
		}
		blobHash := commitments[i].ToVersionedHash()

		blobTx, signErr := gethtypes.SignNewTx(
			testKey,
			signer,
			&gethtypes.BlobTx{
				ChainID:    uint256.NewInt(chainID),
				Nonce:      nextNonce + uint64(i),
				GasTipCap:  uint256.NewInt(10_000_000_000),
				GasFeeCap:  uint256.NewInt(10_000_000_000),
				Gas:        210000,
				To:         recipientAddress,
				Value:      uint256.NewInt(0),
				Data:       []byte{},
				AccessList: nil,
				BlobFeeCap: uint256.NewInt(10_000_000_000),
				BlobHashes: []gethcommon.Hash{blobHash},
				// Sidecar must be nil on signing; we attach it immediately after.
				Sidecar: nil,
			},
		)
		s.Require().NoError(signErr)

		blobTx = blobTx.WithBlobTxSidecar(txSidecar)
		sendErr := node.TestNode.ContractBackend.SendTransaction(node.CtxApp, blobTx)
		s.Require().NoError(sendErr)
	}

	for i := 0; i < 10; i++ {
		pendingNonce, nonceErr := node.TestNode.ContractBackend.PendingNonceAt(node.CtxApp, senderAddress)
		s.Require().NoError(nonceErr)
		if pendingNonce >= nextNonce+uint64(numBlobs) {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	s.T().Fatalf("blob transactions did not enter txpool in time")
}
