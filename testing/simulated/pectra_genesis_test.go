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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/eip7685"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	v1 "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// PectraGenesisSuite defines our test suite for Pectra related work using simulated Comet component.
type PectraGenesisSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors
}

// TestSimulatedCometComponent runs the test suite.
func TestPectraSuite(t *testing.T) {
	suite.Run(t, new(PectraGenesisSuite))
}

// SetupTest initializes the test environment.
func (s *PectraGenesisSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-eth-genesis.json"
	chainSpecFunc := simulated.ProvideElectraGenesisChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), chainSpec, s.HomeDir, elGenesisPath)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Geth node.
	elNode := execution.NewGethNode(s.HomeDir, execution.ValidGethImage())
	elHandle, authRPC, elRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
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
		ClientRPC:   elRPC,
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
func (s *PectraGenesisSuite) TearDownTest() {
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

func (s *PectraGenesisSuite) TestFullLifecycle_WithoutRequests_IsSuccessful() {
	const blockHeight = 1
	const coreLoopIterations = 10

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// Test happens post Electra fork.
	startTime := time.Now()

	// Go through iterations of the core loop.
	proposals, _, _ := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner, startTime)
	s.Require().Len(proposals, coreLoopIterations)
}

func (s *PectraGenesisSuite) TestFullLifecycle_WithPartialWithdrawalRequests_IsSuccessful() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	nextBlockHeight := int64(1)
	// We must first move the chain by 1 height such that the withdrawal contract has an updated `EXCESS_INHIBITOR`.
	{
		proposals, _, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		s.Require().Len(proposals, 1)
		nextBlockHeight++
	}

	// create and submit the withdrawal request
	totalWithdrawalAmount := 3456
	{
		// corresponds with the funded address in genesis `simulated.WithdrawalExecutionAddress`
		senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
		s.Require().NoError(err)

		elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		signer := types.NewPragueSigner(elChainID)

		fee, err := eip7685.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
		s.Require().NoError(err)

		totalTxs := 2
		amountPerTx := totalWithdrawalAmount / totalTxs
		withdrawalTxData, err := eip7685.CreateWithdrawalRequestData(blsSigner.PublicKey(), math.Gwei(amountPerTx))
		s.Require().NoError(err)

		// submit 2 txs
		for i := 0; i < totalTxs; i++ {
			withdrawalTx := types.MustSignNewTx(senderKey, signer, &types.DynamicFeeTx{
				ChainID:   elChainID,
				Nonce:     uint64(i),
				To:        &params.WithdrawalQueueAddress,
				Gas:       500_000,
				GasFeeCap: big.NewInt(1000000000),
				GasTipCap: big.NewInt(1000000000),
				Value:     fee,
				Data:      withdrawalTxData,
			})

			var balance hexutil.Big
			err = s.TestNode.EngineClient.Call(s.CtxApp, &balance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
			s.T().Logf("Balance before withdrawal request sent: %s", balance.ToInt().String())

			var txBytes []byte
			txBytes, err = withdrawalTx.MarshalBinary()
			s.Require().NoError(err)

			var result interface{}
			err = s.TestNode.EngineClient.Call(s.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
			s.Require().NoError(err)
		}
	}

	// Go through 1 iteration of the core loop so that both withdrawal txs is included
	var afterRequestBalance hexutil.Big
	{
		s.LogBuffer.Reset()
		proposals, _, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		s.Require().Len(proposals, 1)
		// Log contains 2 withdrawals
		s.Require().Contains(s.LogBuffer.String(), "Processing execution requests service=state-processor\u001B[0m deposits=0\u001B[0m withdrawals=2\u001B[0m consolidations=0\u001B[0m")

		s.LogBuffer.Reset()
		err := s.TestNode.EngineClient.Call(s.CtxApp, &afterRequestBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.Require().NoError(err)
		s.T().Logf("Balance after withdrawal request included in block: %s", afterRequestBalance.ToInt().String())
		nextBlockHeight++
	}

	// We must progress to Epoch `nextEpoch + MinValidatorWithdrawabilityDelay` before the balance will be removed.
	// IterationsToTurn will get us to the slot before the turn of the target
	var beforeWithdrawalBalance hexutil.Big
	{
		prevBlockHeight := nextBlockHeight - 1
		epochOfWithdrawalRequest := s.TestNode.ChainSpec.SlotToEpoch(math.Slot(prevBlockHeight))
		nextEpoch := epochOfWithdrawalRequest + 1
		targetEpoch := nextEpoch + s.TestNode.ChainSpec.MinValidatorWithdrawabilityDelay()
		iterationsToTurn := (s.TestNode.ChainSpec.SlotsPerEpoch() * uint64(targetEpoch)) - uint64(prevBlockHeight) - 1

		_, _, _ = s.MoveChainToHeight(s.T(), nextBlockHeight, int64(iterationsToTurn), blsSigner, time.Now())

		s.LogBuffer.Reset()
		err := s.TestNode.EngineClient.Call(s.CtxApp, &beforeWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.Require().NoError(err)
		s.T().Logf("Balance before withdrawal processed: %s", beforeWithdrawalBalance.ToInt().String())

		// Balance should not have changed yet
		s.Require().Equal(afterRequestBalance.ToInt().String(), beforeWithdrawalBalance.ToInt().String())
		nextBlockHeight = nextBlockHeight + int64(iterationsToTurn)
	}

	// The next block will be the turn of the Epoch, and the balance will change
	{
		_, _, _ = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())

		var afterWithdrawalBalance hexutil.Big
		err := s.TestNode.EngineClient.Call(s.CtxApp, &afterWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.Require().NoError(err)
		s.T().Logf("Balance after withdrawal processed: %s", afterWithdrawalBalance.ToInt().String())

		withdrawalAmountWei := new(big.Int).Mul(big.NewInt(int64(totalWithdrawalAmount)), big.NewInt(params.GWei))

		// Expected balance is balance before withdrawal + totalWithdrawalAmount
		expectedBalance := new(big.Int).Add(beforeWithdrawalBalance.ToInt(), withdrawalAmountWei)

		// The new balance of the validator is updated
		s.Require().Equal(expectedBalance.String(), afterWithdrawalBalance.ToInt().String())
	}
}

func (s *PectraGenesisSuite) TestFullLifecycle_WithFullWithdrawalRequest_IsSuccessful() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	nextBlockHeight := int64(1)
	// We must first move the chain by 1 height such that the withdrawal contract has an updated `EXCESS_INHIBITOR`.
	{
		proposals, _, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		s.Require().Len(proposals, 1)
		nextBlockHeight = nextBlockHeight + 1
	}

	// create a withdrawal request and submit
	{
		// corresponds with the funded address in genesis `simulated.WithdrawalExecutionAddress`
		senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
		s.Require().NoError(err)

		elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		signer := types.NewPragueSigner(elChainID)

		fee, err := eip7685.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
		s.Require().NoError(err)

		// 0 amount will correspond with a full withdrawal request.
		withdrawalAmount := 0
		withdrawalTxData, err := eip7685.CreateWithdrawalRequestData(blsSigner.PublicKey(), math.Gwei(withdrawalAmount))
		s.Require().NoError(err)

		withdrawalTx := types.MustSignNewTx(senderKey, signer, &types.DynamicFeeTx{
			ChainID:   elChainID,
			Nonce:     0,
			To:        &params.WithdrawalQueueAddress,
			Gas:       500_000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
			Value:     fee,
			Data:      withdrawalTxData,
		})

		var balance hexutil.Big
		err = s.TestNode.EngineClient.Call(s.CtxApp, &balance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.T().Logf("Balance before withdrawal request sent: %s", balance.ToInt().String())

		txBytes, err := withdrawalTx.MarshalBinary()
		s.Require().NoError(err)

		var result interface{}
		err = s.TestNode.EngineClient.Call(s.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
	}

	// Go through 1 iteration of the core loop so that the withdrawal tx is included
	var afterRequestBalance hexutil.Big
	{
		proposals, finalizeBlockResponses, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		s.Require().Len(proposals, 1)
		// Log contains 1 withdrawal
		s.Require().Contains(s.LogBuffer.String(), "Processing execution requests service=state-processor\u001B[0m deposits=0\u001B[0m withdrawals=1\u001B[0m consolidations=0\u001B[0m")
		s.Require().Len(finalizeBlockResponses, 1)
		// No validator updates yet
		s.Require().Len(finalizeBlockResponses[0].GetValidatorUpdates(), 0)

		err := s.TestNode.EngineClient.Call(s.CtxApp, &afterRequestBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.Require().NoError(err)
		s.T().Logf("Balance after withdrawal request included in block: %s", afterRequestBalance.ToInt().String())
		nextBlockHeight++
	}

	// Once a validator's full withdrawal request has been included in a block, it's exit epoch will be set to the next epoch.
	// We enforce that it is exited by checking that FinalizeBlock returns the updated validator set without the validator.
	var exitEpoch math.Epoch
	{
		s.LogBuffer.Reset()
		prevBlockHeight := nextBlockHeight - 1
		epochOfWithdrawalRequest := s.TestNode.ChainSpec.SlotToEpoch(math.Slot(prevBlockHeight))
		nextEpoch := epochOfWithdrawalRequest + 1
		exitEpoch = nextEpoch
		iterationsToExitEpoch := (s.TestNode.ChainSpec.SlotsPerEpoch() * uint64(exitEpoch)) - uint64(prevBlockHeight)

		_, finalizeBlockResponses, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, int64(iterationsToExitEpoch), blsSigner, time.Now())
		s.Require().Len(finalizeBlockResponses, int(iterationsToExitEpoch))
		lastBlockIdx := len(finalizeBlockResponses) - 1
		// We expect the validator to be kicked out now, with power 0
		s.Require().Len(finalizeBlockResponses[lastBlockIdx].GetValidatorUpdates(), 1)
		ejectedValidator := finalizeBlockResponses[lastBlockIdx].GetValidatorUpdates()[0]
		s.Require().Equal(int64(0), ejectedValidator.GetPower())
		s.Require().Equal(blsSigner.PublicKey().String(), hex.EncodeBytes(ejectedValidator.GetPubKeyBytes()))

		nextBlockHeight = nextBlockHeight + int64(iterationsToExitEpoch)
	}

	// We must progress to Epoch `exitEpoch + MinValidatorWithdrawabilityDelay` before the balance will be removed.
	// We progress to the slot before the turn of the target epoch to enforce the balance has not changed.
	var beforeWithdrawalBalance hexutil.Big
	{
		// IterationsToTurn will get us to the slot before the turn of the target
		targetEpoch := exitEpoch + s.TestNode.ChainSpec.MinValidatorWithdrawabilityDelay()
		iterationsToTurn := (s.TestNode.ChainSpec.SlotsPerEpoch() * uint64(targetEpoch)) - uint64(nextBlockHeight)
		_, _, _ = s.MoveChainToHeight(s.T(), nextBlockHeight, int64(iterationsToTurn), blsSigner, time.Now())

		s.LogBuffer.Reset()
		err := s.TestNode.EngineClient.Call(s.CtxApp, &beforeWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.Require().NoError(err)
		s.T().Logf("Balance before withdrawal processed: %s", beforeWithdrawalBalance.ToInt().String())

		// Balance should not have changed yet
		s.Require().Equal(afterRequestBalance.ToInt().String(), beforeWithdrawalBalance.ToInt().String())
		nextBlockHeight = nextBlockHeight + int64(iterationsToTurn)
	}

	// The next block will be the turn of the Epoch, and the balance will change
	{
		_, _, _ = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		var afterWithdrawalBalance hexutil.Big
		err := s.TestNode.EngineClient.Call(s.CtxApp, &afterWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
		s.Require().NoError(err)
		s.T().Logf("Balance after withdrawal processed: %s", afterWithdrawalBalance.ToInt().String())

		// Since this is a full withdrawal, the full balance will be withdrawn.
		// The validator started with a balance equal to math.Gwei(chainSpec.MaxEffectiveBalance())
		withdrawalAmountWei := new(big.Int).Mul(big.NewInt(int64(s.TestNode.ChainSpec.MaxEffectiveBalance())), big.NewInt(params.GWei))

		// Expected balance is balance before withdrawal + withdrawalAmount
		expectedBalance := new(big.Int).Add(beforeWithdrawalBalance.ToInt(), withdrawalAmountWei)

		// The new balance of the validator is updated
		s.Require().Equal(expectedBalance.String(), afterWithdrawalBalance.ToInt().String())
	}
}

// TestMaliciousProposer_AddInvalidExecutionRequests_IsRejected a malicious proposer adds execution requests
// that were not actually requested.
func (s *PectraGenesisSuite) TestMaliciousProposer_AddInvalidExecutionRequests_IsRejected() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	pubkey, err := blsSigner.GetPubKey()
	s.Require().NoError(err)

	nextBlockHeight := int64(1)
	// We must first move the chain by 1 height such that the withdrawal contract has an updated `EXCESS_INHIBITOR`.
	{
		proposals, _, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		s.Require().Len(proposals, 1)
		nextBlockHeight++
	}

	// Create a signed block with invalid execution requests.
	var maliciousSignedBlock *ctypes.SignedBeaconBlock
	var proposal *v1.PrepareProposalResponse
	proposalTime := time.Now()
	{
		s.LogBuffer.Reset()
		proposal, err = s.SimComet.Comet.PrepareProposal(s.CtxComet, &v1.PrepareProposalRequest{
			Height:          nextBlockHeight,
			Time:            time.Now(),
			ProposerAddress: pubkey.Address(),
		})
		s.Require().NoError(err)
		s.Require().Len(proposal.Txs, 2)
		// Unmarshal the proposal block.
		proposedBlock, unmarshalErr := encoding.UnmarshalBeaconBlockFromABCIRequest(
			proposal.Txs,
			blockchain.BeaconBlockTxIndex,
			s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(math.U64(proposalTime.Unix())),
		)
		s.Require().NoError(unmarshalErr)

		// Invalid Execution Request
		invalidExecutionRequests := &ctypes.ExecutionRequests{
			Deposits: []*ctypes.DepositRequest{
				{
					Pubkey:      [48]byte{0, 1, 2},
					Credentials: [32]byte{0, 3, 2},
					Amount:      10000000,
					Signature:   [96]byte{5, 6, 7},
					Index:       5,
				},
			},
			Withdrawals:    nil,
			Consolidations: nil,
		}

		// Create a malicious block by injecting an invalid Execution Request.
		maliciousBlock := simulated.ComputeAndSetInvalidExecutionBlock(
			s.T(), proposedBlock.GetBeaconBlock(), s.TestNode.ChainSpec, nil, invalidExecutionRequests,
		)
		// Re-sign the block
		maliciousSignedBlock, err = ctypes.NewSignedBeaconBlock(
			maliciousBlock,
			&ctypes.ForkData{
				CurrentVersion:        s.TestNode.ChainSpec.ActiveForkVersionForTimestamp(maliciousBlock.GetTimestamp()),
				GenesisValidatorsRoot: s.GenesisValidatorsRoot,
			},
			s.TestNode.ChainSpec,
			blsSigner,
		)
		s.Require().NoError(err)

		// Check that the block contains the invalid execution request.
		requests, getErr := maliciousSignedBlock.GetBeaconBlock().GetBody().GetExecutionRequests()
		s.Require().NoError(getErr)
		s.Require().Len(requests.Deposits, 1)

	}
	// Propose the invalid block
	{
		maliciousBlockBytes, sszErr := maliciousSignedBlock.MarshalSSZ()
		s.Require().NoError(sszErr)

		// Replace the valid block with the malicious block in the proposal.
		proposal.Txs[0] = maliciousBlockBytes

		// Reset the log buffer to discard old logs we don't care about
		s.LogBuffer.Reset()
		// Process the proposal containing the malicious block.
		processResp, err := s.SimComet.Comet.ProcessProposal(s.CtxComet, &v1.ProcessProposalRequest{
			Txs:             proposal.Txs,
			Height:          nextBlockHeight,
			ProposerAddress: pubkey.Address(),
			Time:            proposalTime,
		})
		s.Require().NoError(err)
		s.Require().Equal(v1.PROCESS_PROPOSAL_STATUS_REJECT, processResp.Status)

		// Verify that the log contains the expected error message.
		s.Require().Contains(s.LogBuffer.String(), errors.ErrInvalidPayloadStatus.Error())
		s.Require().Contains(s.LogBuffer.String(), "invalid requests hash (remote: 33ba74e937423115e3abf4250db02588388b4b3a7918950ed44a28e4bf3428d2 local: e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855)")
	}
}

// TestMaliciousUser_MakesConsolidationRequest_IsIgnored a user makes a consolidation request on our chain
// which isn't supported.
func (s *PectraGenesisSuite) TestMaliciousUser_MakesConsolidationRequest_IsIgnored() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	nextBlockHeight := int64(1)
	// We must first move the chain by 1 height such that the consolidation contract has an updated `EXCESS_INHIBITOR`.
	{
		_, _, _ = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		nextBlockHeight++
	}
	// Next we submit the Consolidation request transaction
	{
		// corresponds with funded address in genesis
		senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
		s.Require().NoError(err)

		elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		signer := types.NewPragueSigner(elChainID)

		//fee, feeErr := eip7685.GetConsolidationFee(s.CtxApp, s.TestNode.EngineClient)
		//s.Require().NoError(feeErr)

		// The inputs to the request do not necessarily matter, as long as they pass EL validation
		consolidationTxData, requestErr := eip7685.CreateConsolidationRequestData(blsSigner.PublicKey(), blsSigner.PublicKey())
		s.Require().NoError(requestErr)

		consolidationTx := types.MustSignNewTx(senderKey, signer, &types.DynamicFeeTx{
			ChainID:   elChainID,
			Nonce:     0,
			To:        &params.WithdrawalQueueAddress,
			Gas:       500_000,
			GasFeeCap: big.NewInt(1000000000),
			GasTipCap: big.NewInt(1000000000),
			Value:     big.NewInt(1000),
			Data:      consolidationTxData,
		})
		txBytes, marshalErr := consolidationTx.MarshalBinary()
		s.Require().NoError(marshalErr)
		var result interface{}
		err = s.TestNode.EngineClient.Call(s.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
	}
	// Move the chain so that tx is included
	{
		_, _, _ = s.MoveChainToHeight(s.T(), nextBlockHeight, 5, blsSigner, time.Now())
		nextBlockHeight++
	}
}
