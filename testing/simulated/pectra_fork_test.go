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
	"math/big"
	"path"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/execution/requests/eip7251"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// PectraForkSuite defines our test suite for Pectra related work using simulated Comet component.
type PectraForkSuite struct {
	suite.Suite
	Geth simulated.SharedAccessors
	Reth simulated.SharedAccessors
}

// TestPectraForkSuite runs the test suite.
func TestPectraForkSuite(t *testing.T) {
	suite.Run(t, new(PectraForkSuite))
}

// SetupTest initializes the test environment.
func (s *PectraForkSuite) SetupTest() {
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
func (s *PectraForkSuite) TearDownTest() {
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

// TestTimestampFork_ELAndCLInSync_IsSuccessful tests that we can fork successfully if EL and CL have synced timestamps
// The forks timestamp at Unix 0, as the genesis at Unix 0, Cancun is at 10 and Prague is at 20.
// The Geth Node will be the block producer but the Reth node is treated as a full node, i.e. doesn't produce blocks.
func (s *PectraForkSuite) TestTimestampFork_ELAndCLInSync_IsSuccessful() {
	// Initialize the geth chain state.
	s.Geth.InitializeChain(s.T())
	// Initialize the reth chain state.
	s.Reth.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)

	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	expectedMessages := []string{
		"Processing block with fork version service=blockchain\u001B[0m block=1\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=2\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=3\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=4\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=5\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=6\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=7\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=8\u001B[0m fork=0x04010000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=9\u001B[0m fork=0x05000000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=10\u001B[0m fork=0x05000000\u001B[0m",
		"Processing block with fork version service=blockchain\u001B[0m block=11\u001B[0m fork=0x05000000\u001B[0m",
	}

	var (
		startHeight         = int64(1)
		iterations          = int64(len(expectedMessages))
		expectedMessagesIdx = 0
		submitTxNonce       = uint64(0)
		consensusTime       = time.Unix(startHeight*2, 0)
	)

	for currentHeight := startHeight; currentHeight < startHeight+iterations; currentHeight++ {
		submitTxNonce = s.submitTransactions(submitTxNonce, 100)
		proposal, err := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          currentHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}

		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: pubkey.Address(),
			Time:            consensusTime,
		}
		expectedMessage := expectedMessages[expectedMessagesIdx]
		processFinalizeCommit(s.T(), s.Geth, processRequest, finalizeRequest, expectedMessage)
		processFinalizeCommit(s.T(), s.Reth, processRequest, finalizeRequest, expectedMessage)

		expectedMessagesIdx++

		// set consensus time for the next block to match
		// the timestamp of the payload built optimistically.
		forkVersion := s.Geth.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(consensusTime.Unix())) //#nosec: G115
		blk, _, err := encoding.ExtractBlobsAndBlockFromRequest(
			processRequest,
			blockchain.BeaconBlockTxIndex,
			blockchain.BlobSidecarsTxIndex,
			forkVersion,
		)
		s.Require().NoError(err)
		consensusTime = time.Unix(
			int64(payloadtime.Next(blk.GetTimestamp(), blk.GetTimestamp(), true)),
			0,
		)
	}
}

// A user makes a consolidation request on our chain which isn't supported.
func (s *PectraForkSuite) TestMaliciousUser_MakesConsolidationRequest_IsIgnored() {
	// Initialize the chain state.
	s.Geth.InitializeChain(s.T())
	s.Reth.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// We must first move the chain to the fork height, then an extra block
	// such that the consolidation contract has an updated `EXCESS_INHIBITOR`.
	// We set the timestamp such that the fork has occurred (i.e., time.Now)
	{
		for i := 0; i < 2; i++ {
			consensusTime := time.Now()
			proposal, err := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
				Height:          nextBlockHeight,
				Time:            consensusTime,
				ProposerAddress: pubkey.Address(),
			})
			s.Require().NoError(err)
			s.Require().Len(proposal.Txs, 2)

			processRequest := &types.ProcessProposalRequest{
				Txs:             proposal.Txs,
				Height:          nextBlockHeight,
				ProposerAddress: pubkey.Address(),
				Time:            consensusTime,
			}

			finalizeRequest := &types.FinalizeBlockRequest{
				Txs:             proposal.Txs,
				Height:          nextBlockHeight,
				ProposerAddress: pubkey.Address(),
				Time:            consensusTime,
			}
			expectedMsg := "fork=0x05000000"
			processFinalizeCommit(s.T(), s.Geth, processRequest, finalizeRequest, expectedMsg)
			processFinalizeCommit(s.T(), s.Reth, processRequest, finalizeRequest, expectedMsg)
			nextBlockHeight++
		}
	}
	// Next we submit the Consolidation request transaction
	{
		// corresponds with the funded address in genesis
		senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
		s.Require().NoError(err)

		elChainID := big.NewInt(int64(s.Geth.TestNode.ChainSpec.DepositEth1ChainID()))
		signer := gethtypes.NewPragueSigner(elChainID)

		fee, feeErr := eip7251.GetConsolidationFee(s.Geth.CtxApp, s.Geth.TestNode.EngineClient)
		s.Require().NoError(feeErr)

		// The inputs to the request do not necessarily matter, as long as they pass EL validation
		consolidationTxData, requestErr := eip7251.CreateConsolidationRequestData(blsSigner.PublicKey(), blsSigner.PublicKey())
		s.Require().NoError(requestErr)

		consolidationTx := gethtypes.MustSignNewTx(senderKey, signer, &gethtypes.DynamicFeeTx{
			ChainID:   elChainID,
			Nonce:     0,
			To:        &params.ConsolidationQueueAddress,
			Gas:       500_000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
			Value:     fee,
			Data:      consolidationTxData,
		})
		txBytes, marshalErr := consolidationTx.MarshalBinary()
		s.Require().NoError(marshalErr)
		var result interface{}
		err = s.Geth.TestNode.EngineClient.Call(s.Geth.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
		time.Sleep(time.Second) // give it time to allow the tx to be included in the next block
	}
	// Move the chain so that tx is included and progresses correctly afterward.
	{
		for i := 0; i < 5; i++ {
			consensusTime := time.Now()
			proposal, err := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
				Height:          nextBlockHeight,
				Time:            consensusTime,
				ProposerAddress: pubkey.Address(),
			})
			s.Require().NoError(err)
			s.Require().Len(proposal.Txs, 2)

			processRequest := &types.ProcessProposalRequest{
				Txs:             proposal.Txs,
				Height:          nextBlockHeight,
				ProposerAddress: pubkey.Address(),
				Time:            consensusTime,
			}

			finalizeRequest := &types.FinalizeBlockRequest{
				Txs:             proposal.Txs,
				Height:          nextBlockHeight,
				ProposerAddress: pubkey.Address(),
				Time:            consensusTime,
			}
			var expectedMsg string
			if i == 0 {
				// The first block since tx submission will have the consolidation propagated to the CL.
				expectedMsg = "consolidations=1"
			}
			processFinalizeCommit(s.T(), s.Geth, processRequest, finalizeRequest, expectedMsg)
			processFinalizeCommit(s.T(), s.Reth, processRequest, finalizeRequest, expectedMsg)
			nextBlockHeight++
		}
	}
}

// This test will have a proposer propose a valid post-fork block, but one that is not finalized.
// The next round will propose a valid pre-fork block that gets finalized due to deviance in the consensus timestamp.
// The proposer will then propose a valid post-fork block that is correctly finalized.
//
// NOTE: this test requires reth with the --engine.always-process-payload-attributes-on-canonical-head flag
// to propose the valid pre-fork block.
func (s *PectraForkSuite) TestValidProposer_ProposesPostForkBlockIsNotFinalized_IsSuccessful() {
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

// The proposer prepares a proposal with a pre-fork timestamp, but a post-fork process proposal consensus time.
// This will be rejected and is expected to occur around the fork for 1 or 2 rounds.
func (s *PectraForkSuite) TestValidProposer_ProposesPreForkBlockWithPostForkConsensusTimestamp_IsRejected() {
	// Initialize the chain state.
	s.Geth.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.Geth.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// The proposer prepares a proposal with a pre-fork timestamp, but a post-fork process proposal consensus time.
	{
		// a pre-fork time.
		consensusTime := time.Unix(2, 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Now(),
		}

		// Process the proposal, expect a reject
		{
			s.Geth.LogBuffer.Reset()
			processResp, err := s.Geth.SimComet.Comet.ProcessProposal(s.Geth.CtxComet, processRequest)
			s.Require().NoError(err)
			s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)
			s.Require().Contains(
				s.Geth.LogBuffer.String(),
				"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous: decoded 392, previous was 396",
			)
		}
	}
}

func processFinalizeCommit(
	t *testing.T,
	node simulated.SharedAccessors,
	processRequest *types.ProcessProposalRequest,
	finalizeRequest *types.FinalizeBlockRequest,
	expectedMessage string,
) {
	// Process the proposal
	node.LogBuffer.Reset()
	processResp, err := node.SimComet.Comet.ProcessProposal(node.CtxComet, processRequest)
	require.NoError(t, err)
	require.Equal(t, types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)
	require.Contains(t, node.LogBuffer.String(), expectedMessage)

	// Finalize the block
	finalizeResp, err := node.SimComet.Comet.FinalizeBlock(node.CtxComet, finalizeRequest)
	require.NoError(t, err)
	require.NotEmpty(t, finalizeResp)

	// Commit the block.
	_, err = node.SimComet.Comet.Commit(node.CtxComet, &types.CommitRequest{})
	require.NoError(t, err)
}

func (s *PectraForkSuite) submitTransactions(startNonce uint64, numTransactions uint64) uint64 {
	// corresponds with funded address in genesis 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4
	senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)
	elChainID := big.NewInt(int64(s.Geth.TestNode.ChainSpec.DepositEth1ChainID()))
	signer := gethtypes.NewPragueSigner(elChainID)

	for i := startNonce; i < startNonce+numTransactions; i++ {
		transaction := gethtypes.MustSignNewTx(senderKey, signer, &gethtypes.DynamicFeeTx{
			ChainID:   elChainID,
			Nonce:     i,
			To:        &params.BeaconRootsAddress, // any address
			Gas:       500_000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
			Value:     big.NewInt(0),
			Data:      nil,
		})

		txBytes, marshalErr := transaction.MarshalBinary()
		s.Require().NoError(marshalErr)

		var result interface{}
		err = s.Geth.TestNode.EngineClient.Call(s.Geth.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
	}
	return startNonce + numTransactions
}
