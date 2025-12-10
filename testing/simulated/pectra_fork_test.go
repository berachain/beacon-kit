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
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), chainSpec, s.Geth.HomeDir, elGenesisPath)
	s.Geth.GenesisValidatorsRoot = genesisValidatorsRoot
	s.Reth.GenesisValidatorsRoot = genesisValidatorsRoot

	// Copy the home dir for the Reth Node
	simulated.CopyHomeDir(s.T(), s.Geth.HomeDir, s.Reth.HomeDir)

	// Start the EL (execution layer) Geth node.
	gethNode := execution.NewGethNode(s.Geth.HomeDir, execution.ValidGethImage())
	elHandle, authRPC, elRPC := gethNode.Start(s.T(), path.Base(elGenesisPath))
	s.Geth.ElHandle = elHandle

	rethNode := execution.NewRethNode(s.Reth.HomeDir, execution.ValidRethImage())
	rethHandle, rethAuthRPC, rethElRPC := rethNode.Start(s.T(), path.Base(elGenesisPath))
	s.Reth.ElHandle = rethHandle

	// Prepare a logger backed by a buffer to capture logs for assertions.
	s.Geth.LogBuffer = &simulated.SyncBuffer{}
	logger := phuslu.NewLogger(s.Geth.LogBuffer, nil)

	s.Reth.LogBuffer = &simulated.SyncBuffer{}
	rethLogger := phuslu.NewLogger(s.Reth.LogBuffer, nil)

	// Build the Beacon node with the simulated Comet component and electra genesis chain spec
	components := simulated.FixedComponents(s.T())
	components = append(components, simulated.ProvideSimComet)
	components = append(components, chainSpecFunc)

	s.Geth.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Geth.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     authRPC,
		ClientRPC:   elRPC,
		Logger:      logger,
		AppOpts:     viper.New(),
		Components:  components,
	})
	s.Geth.SimComet = s.Geth.TestNode.SimComet

	s.Reth.TestNode = simulated.NewTestNode(s.T(), simulated.TestNodeInput{
		TempHomeDir: s.Reth.HomeDir,
		CometConfig: cometConfig,
		AuthRPC:     rethAuthRPC,
		ClientRPC:   rethElRPC,
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
	s.Geth.CleanupTestWithLabel(s.T(), "GETH")
	s.Reth.CleanupTestWithLabel(s.T(), "RETH")
}

// TestTimestampFork_ELAndCLInSync_IsSuccessful tests that we can fork successfully if EL and CL have synced timestamps
// The forks timestamp at Unix 0, as the genesis at Unix 0, Cancun is at 10 and Prague is at 20.
// The Geth Node will be the block producer but the Reth node is treated as a full node, i.e. doesn't produce blocks.
func (s *PectraForkSuite) TestTimestampFork_ELAndCLInSync_IsSuccessful() {
	// Initialize the geth chain state.
	s.Geth.InitializeChain(s.T())
	// Initialize the reth chain state.
	s.Reth.InitializeChain(s.T())

	nodeAddress, err := s.Geth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Geth.SimComet.Comet.SetNodeAddress(nodeAddress)

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
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(err)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              currentHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          currentHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		expectedMessage := expectedMessages[expectedMessagesIdx]
		processFinalizeCommit(s.T(), s.Geth, processRequest, finalizeRequest, expectedMessage)
		processFinalizeCommit(s.T(), s.Reth, processRequest, finalizeRequest, expectedMessage)

		expectedMessagesIdx++

		// set consensus time for the next block to match
		// the timestamp of the payload built optimistically.
		forkVersion := s.Geth.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(consensusTime.Unix())) //#nosec: G115
		blk, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
			processRequest.GetTxs(),
			blockchain.BeaconBlockTxIndex,
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
	nodeAddress := pubkey.Address()
	s.Geth.SimComet.Comet.SetNodeAddress(nodeAddress)

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
				ProposerAddress: nodeAddress,
			})
			s.Require().NoError(err)
			s.Require().Len(proposal.Txs, 2)

			processRequest := &types.ProcessProposalRequest{
				Txs:                 proposal.Txs,
				Height:              nextBlockHeight,
				ProposerAddress:     nodeAddress,
				Time:                consensusTime,
				NextProposerAddress: nodeAddress,
			}

			finalizeRequest := &types.FinalizeBlockRequest{
				Txs:             proposal.Txs,
				Height:          nextBlockHeight,
				ProposerAddress: nodeAddress,
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
				ProposerAddress: nodeAddress,
			})
			s.Require().NoError(err)
			s.Require().Len(proposal.Txs, 2)

			processRequest := &types.ProcessProposalRequest{
				Txs:                 proposal.Txs,
				Height:              nextBlockHeight,
				ProposerAddress:     nodeAddress,
				Time:                consensusTime,
				NextProposerAddress: nodeAddress,
			}

			finalizeRequest := &types.FinalizeBlockRequest{
				Txs:             proposal.Txs,
				Height:          nextBlockHeight,
				ProposerAddress: nodeAddress,
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
func (s *PectraForkSuite) TestValidProposer_ProposesPostForkBlockIsNotFinalized_IsSuccessful() {
	client := s.Geth
	helper := s.Reth

	// Initialize the chain state.
	client.InitializeChain(s.T())
	helper.InitializeChain(s.T()) // helper to build "invalid" blocks

	nodeAddress, err := client.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	client.SimComet.Comet.SetNodeAddress(nodeAddress)

	nextBlockHeight := int64(1)
	var proposal *types.PrepareProposalResponse

	// 1 - Build a block whose consensus and payloadTimestamp are both post-fork.
	//     Check that it verifies, but do not finalize it
	{
		consensusTime := time.Unix(int64(client.TestNode.ChainSpec.ElectraForkTime()), 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		}
		proposal, err = helper.SimComet.Comet.PrepareProposal(helper.CtxComet, prepareReq)
		s.Require().NoError(err)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal
		client.LogBuffer.Reset()
		processResp, respErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	// 2 - Build a block whose consensus timestamp is pre-fork, while the payload is post fork.
	//     Check that it does not verifies.
	// Note: to build the invalid block we reuse the beaconBlock from point 1 and just change CometBFT timestamp
	{
		maliciouConsensusTime := time.Unix(int64(client.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                maliciouConsensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal
		client.LogBuffer.Reset()
		processResp, processErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT.String(), processResp.Status.String())
		s.Require().Contains(
			client.LogBuffer.String(),
			"failed decoding *types.SignedBeaconBlock: ssz: offset smaller than previous",
		)
	}

	// 3 - Build a block whose consensus and payload timestamp are pre-fork.
	//     Check that it does verify (even if we already validated a post fork block).
	// Note: to build the block we reuse the beaconBlock from point 1 and just change CometBFT timestamp
	{
		consensusTime := time.Unix(int64(client.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		}
		proposal, prepareErr := helper.SimComet.Comet.PrepareProposal(helper.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal
		client.LogBuffer.Reset()
		processResp, processErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
	}

	// The next block the proposer proposes with a pre-fork timestamp will actually have a pre-fork time
	// Since the previous payload in cache has been evicted and a new payload is retrieved.
	{
		consensusTime := time.Unix(int64(client.TestNode.ChainSpec.ElectraForkTime())-2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		}
		proposal, prepareErr := helper.SimComet.Comet.PrepareProposal(helper.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}

		// Process the proposal
		client.LogBuffer.Reset()
		processResp, processErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// Finalize the block
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: nodeAddress,
			Time:            consensusTime,
		}
		_, finalizeErr := client.SimComet.Comet.FinalizeBlock(client.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)

		// Commit the block.
		_, err = client.SimComet.Comet.Commit(client.CtxComet, &types.CommitRequest{})
		s.Require().NoError(err)

		nextBlockHeight++
	}
	// Finally, we cross the fork and show no issues
	{
		consensusTime := time.Unix(int64(client.TestNode.ChainSpec.ElectraForkTime())+2, 0)
		prepareReq := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		}
		proposal, prepareErr := client.SimComet.Comet.PrepareProposal(client.CtxComet, prepareReq)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                consensusTime,
			NextProposerAddress: nodeAddress,
		}
		// Process the proposal
		client.LogBuffer.Reset()
		processResp, processErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(processErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())
		s.Require().Contains(client.LogBuffer.String(), "Processing execution requests")
	}
}

// The proposer prepares a proposal with a pre-fork timestamp, but a post-fork process proposal consensus time.
// This will be rejected and is expected to occur around the fork for 1 or 2 rounds.
func (s *PectraForkSuite) TestValidProposer_ProposesPreForkBlockWithPostForkConsensusTimestamp_IsRejected() {
	// Initialize the chain state.
	s.Geth.InitializeChain(s.T())
	nodeAddress, err := s.Geth.SimComet.GetNodeAddress()
	s.Require().NoError(err)
	s.Geth.SimComet.Comet.SetNodeAddress(nodeAddress)

	nextBlockHeight := int64(1)
	// The proposer prepares a proposal with a pre-fork timestamp, but a post-fork process proposal consensus time.
	{
		// a pre-fork time.
		consensusTime := time.Unix(2, 0)
		proposal, prepareErr := s.Geth.SimComet.Comet.PrepareProposal(s.Geth.CtxComet, &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            consensusTime,
			ProposerAddress: nodeAddress,
		})
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		processRequest := &types.ProcessProposalRequest{
			Txs:                 proposal.Txs,
			Height:              nextBlockHeight,
			ProposerAddress:     nodeAddress,
			Time:                time.Now(),
			NextProposerAddress: nodeAddress,
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

// This test will show that an optimistically building a payload across the fork boundary
// correctly invokes `ProcessFork` on the state processor.
func (s *PectraForkSuite) Test_OptimisticBuildAtFork_IsSuccessful() {
	// Initialize the chain state.
	client := s.Geth
	client.InitializeChain(s.T()) // init the validator

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(client.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	// setup first payload and consensus timestamps so that
	// - it would be pre Electra fork
	// - the second block, built optimistically would be at the electra fork.
	specs := client.TestNode.ChainSpec
	firstBlkConsensusTime := specs.ElectraForkTime() - 1 // before fork
	firstBlkPayloadTime := firstBlkConsensusTime
	secondBlkConsensusTime := specs.ElectraForkTime()
	secondBlkPayloadTime := payloadtime.Next(
		math.U64(secondBlkConsensusTime),
		math.U64(firstBlkPayloadTime),
		true, // this is the formula used while setting second block timestamp optimistically
	)
	s.Require().GreaterOrEqual(secondBlkPayloadTime, math.U64(specs.ElectraForkTime())) // post fork

	nextBlockHeight := int64(1)
	{
		// 1- Build pre-fork block
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            time.Unix(int64(firstBlkConsensusTime), 0),
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := client.SimComet.Comet.PrepareProposal(client.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)

		// 2- Process the proposal. This will trigger am optimistic payload build for block height 2.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(firstBlkConsensusTime), 0),
		}
		processResp, respErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// 3- finalize and commit the first block
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(firstBlkConsensusTime), 0),
		}
		_, finalizeErr := client.SimComet.Comet.FinalizeBlock(client.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := client.SimComet.Comet.Commit(client.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
	}

	// Now build the next block
	nextBlockHeight++
	{
		// 4- Build post-fork block. Make sure that the fork transition happens
		prepareRequest := &types.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            time.Unix(int64(secondBlkConsensusTime), 0),
			ProposerAddress: pubkey.Address(),
		}
		proposal, prepareErr := client.SimComet.Comet.PrepareProposal(client.CtxComet, prepareRequest)
		s.Require().NoError(prepareErr)
		s.Require().Len(proposal.Txs, 2)
		s.Require().Contains(client.LogBuffer.String(), "✅  welcome to the")

		// 5- Process the proposal.
		processRequest := &types.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(secondBlkConsensusTime), 0),
		}
		processResp, respErr := client.SimComet.Comet.ProcessProposal(client.CtxComet, processRequest)
		s.Require().NoError(respErr)
		s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_ACCEPT.String(), processResp.Status.String())

		// 6- finalize and commit the second block
		finalizeRequest := &types.FinalizeBlockRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            time.Unix(int64(secondBlkConsensusTime), 0),
		}
		_, finalizeErr := client.SimComet.Comet.FinalizeBlock(client.CtxComet, finalizeRequest)
		s.Require().NoError(finalizeErr)
		_, commitErr := client.SimComet.Comet.Commit(client.CtxComet, &types.CommitRequest{})
		s.Require().NoError(commitErr)
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
