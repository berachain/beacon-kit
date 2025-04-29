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

	"github.com/berachain/beacon-kit/log/phuslu"
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

// TestSimulatedCometComponent runs the test suite.
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
	submitTxNonce := uint64(0)
	for currentHeight := startHeight; currentHeight < startHeight+iterations; currentHeight++ {
		submitTxNonce = s.submitTransactions(submitTxNonce, 100)
		// We set the consensus time to currentHeight * 2 to mimick a functional 2s block time.
		// This assumption is not always true but used in a valid test.
		consensusTime := time.Unix(currentHeight*2, 0)
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
	}
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

		txBytes, err := transaction.MarshalBinary()
		s.Require().NoError(err)

		var result interface{}
		err = s.Geth.TestNode.EngineClient.Call(s.Geth.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
	}
	return startNonce + numTransactions
}

// TestExcessBalanceAtFork this test will demonstrate the systems behaviour when there is withdrawal is
// initiated before the fork and set to complete after the fork. The withdrawal is due to excess validators.
func (s *PectraForkSuite) TestTODO_0() {
	s.T().Skip("TODO: Implement this test")
}

// This test will have a proposer propose a valid post-fork block, but one that is never finalized. The round will increase.
// We should observe that the EL correctly handles this and no execution requests are included in the block.
func (s *PectraForkSuite) TestTODO_1() {
	s.T().Skip("TODO: Implement this test")
}

// A valid proposer will attempt to propose a pre-fork block but with a post-fork consensus timestamp.
// We expect this to be rejected.
func (s *PectraForkSuite) TestTODO_2() {
	s.T().Skip("TODO: Implement this test")
}

func processFinalizeCommit(
	t *testing.T,
	node simulated.SharedAccessors,
	processRequest *types.ProcessProposalRequest,
	finalizeRequest *types.FinalizeBlockRequest,
	expectedMessage string,
) {
	// Process the proposal
	processResp, err := node.SimComet.Comet.ProcessProposal(node.CtxComet, processRequest)
	require.NoError(t, err)
	require.Equal(t, types.PROCESS_PROPOSAL_STATUS_ACCEPT, processResp.Status)

	// Finalize the block
	node.LogBuffer.Reset()
	finalizeResp, err := node.SimComet.Comet.FinalizeBlock(node.CtxComet, finalizeRequest)
	require.NoError(t, err)
	require.Contains(t, node.LogBuffer.String(), expectedMessage)
	require.NotEmpty(t, finalizeResp)

	// Commit the block.
	_, err = node.SimComet.Comet.Commit(node.CtxComet, &types.CommitRequest{})
	require.NoError(t, err)
}
