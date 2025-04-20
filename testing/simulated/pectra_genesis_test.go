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
	"github.com/berachain/beacon-kit/primitives/eip7002"
	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
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
	elHandle, authRPC := elNode.Start(s.T(), path.Base(elGenesisPath))
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

func (s *PectraGenesisSuite) TestFullLifecycle_WithPartialWithdrawalRequest_IsSuccessful() {
	const blockHeight = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// create withdrawal request
	// corresponds with funded address in genesis `simulated.WithdrawalExecutionAddress`
	senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)

	elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
	signer := types.NewPragueSigner(elChainID)

	fee, err := eip7002.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
	s.Require().NoError(err)

	withdrawalAmount := 3456
	withdrawalTxData, err := eip7002.CreateWithdrawalRequestData(blsSigner.PublicKey(), math.Gwei(withdrawalAmount))
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

	// Go through 1 iteration of the core loop so that the withdrawal tx is included
	s.LogBuffer.Reset()
	proposals, _, _ := s.MoveChainToHeight(s.T(), blockHeight, 1, blsSigner, time.Now())
	s.Require().Len(proposals, 1)
	// Log contains 1 withdrawal
	s.Require().Contains(s.LogBuffer.String(), "Processing execution requests service=state-processor\u001B[0m deposits=0\u001B[0m withdrawals=1\u001B[0m consolidations=0\u001B[0m")

	s.LogBuffer.Reset()
	var afterRequestBalance hexutil.Big
	err = s.TestNode.EngineClient.Call(s.CtxApp, &afterRequestBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
	s.T().Logf("Balance after withdrawal request included in block: %s", afterRequestBalance.ToInt().String())

	// We must progress to Epoch `nextEpoch + MinValidatorWithdrawabilityDelay` before the balance will be removed.
	// As such, we move the chain to the height `nextEpoch + MinValidatorWithdrawabilityDelay - 1`.
	lastBlockHeight := blockHeight + 1
	iterations := s.TestNode.ChainSpec.SlotsPerEpoch() * s.TestNode.ChainSpec.MinValidatorWithdrawabilityDelay()
	proposals, _, _ = s.MoveChainToHeight(s.T(), int64(lastBlockHeight), int64(iterations), blsSigner, time.Now())

	s.LogBuffer.Reset()
	var beforeWithdrawalBalance hexutil.Big
	err = s.TestNode.EngineClient.Call(s.CtxApp, &beforeWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
	s.T().Logf("Balance before withdrawal processed: %s", beforeWithdrawalBalance.ToInt().String())

	// Balance should not have changed yet
	s.Require().Equal(afterRequestBalance.ToInt().String(), beforeWithdrawalBalance.ToInt().String())

	// The next block will be the turn of the Epoch, and balance will change
	lastBlockHeight = int(iterations) + lastBlockHeight
	proposals, _, _ = s.MoveChainToHeight(s.T(), int64(lastBlockHeight), 1, blsSigner, time.Now())

	var afterWithdrawalBalance hexutil.Big
	err = s.TestNode.EngineClient.Call(s.CtxApp, &afterWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
	s.T().Logf("Balance after withdrawal processed: %s", afterWithdrawalBalance.ToInt().String())

	withdrawalAmountWei := new(big.Int).Mul(big.NewInt(int64(withdrawalAmount)), big.NewInt(params.GWei))

	// Expected balance is balance before withdrawal + withdrawalAmount
	expectedBalance := new(big.Int).Add(beforeWithdrawalBalance.ToInt(), withdrawalAmountWei)

	// The new balance of the validator is updated
	s.Require().Equal(expectedBalance.String(), afterWithdrawalBalance.ToInt().String())
}

func (s *PectraGenesisSuite) TestFullLifecycle_WithFullWithdrawalRequest_IsSuccessful() {
	const blockHeight = 1

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// create withdrawal request
	// corresponds with funded address in genesis `simulated.WithdrawalExecutionAddress`
	senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)

	elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
	signer := types.NewPragueSigner(elChainID)

	fee, err := eip7002.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
	s.Require().NoError(err)

	// 0 amount will correspond with a full withdrawal request.
	withdrawalAmount := 0
	withdrawalTxData, err := eip7002.CreateWithdrawalRequestData(blsSigner.PublicKey(), math.Gwei(withdrawalAmount))
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

	// Go through 1 iteration of the core loop so that the withdrawal tx is included
	s.LogBuffer.Reset()
	proposals, finalizeBlockResponses, _ := s.MoveChainToHeight(s.T(), blockHeight, 1, blsSigner, time.Now())
	s.Require().Len(proposals, 1)
	// Log contains 1 withdrawal
	s.Require().Contains(s.LogBuffer.String(), "Processing execution requests service=state-processor\u001B[0m deposits=0\u001B[0m withdrawals=1\u001B[0m consolidations=0\u001B[0m")
	s.Require().Len(finalizeBlockResponses, 1)
	// No validator updates yet
	s.Require().Len(finalizeBlockResponses[0].GetValidatorUpdates(), 0)

	s.LogBuffer.Reset()
	var afterRequestBalance hexutil.Big
	err = s.TestNode.EngineClient.Call(s.CtxApp, &afterRequestBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
	s.T().Logf("Balance after withdrawal request included in block: %s", afterRequestBalance.ToInt().String())

	// Once a validator's full withdrawal request has been included in a block, it's exit epoch will be set to the next epoch.
	// We enforce that it is exited by checking that FinalizeBlock returns the updated validator set without the validator.
	nextBlockHeight := blockHeight + 1
	proposals, finalizeBlockResponses, _ = s.MoveChainToHeight(s.T(), int64(nextBlockHeight), int64(1), blsSigner, time.Now())
	s.Require().Len(finalizeBlockResponses, 1)
	// We expect the validator to be kicked out now, with power 0
	s.Require().Len(finalizeBlockResponses[0].GetValidatorUpdates(), 1)
	ejectedValidator := finalizeBlockResponses[0].GetValidatorUpdates()[0]
	s.Require().Equal(int64(0), ejectedValidator.GetPower())
	s.Require().Equal(blsSigner.PublicKey().String(), hex.EncodeBytes(ejectedValidator.GetPubKeyBytes()))

	// We also expect the validators balance to have updated.

	// We must progress to Epoch `nextEpoch + MinValidatorWithdrawabilityDelay` before the balance will be removed.
	// As such, we move the chain to the height `nextEpoch + MinValidatorWithdrawabilityDelay - 1`.
	nextBlockHeight = nextBlockHeight + 1
	iterations := s.TestNode.ChainSpec.SlotsPerEpoch()*s.TestNode.ChainSpec.MinValidatorWithdrawabilityDelay() - 1
	proposals, _, _ = s.MoveChainToHeight(s.T(), int64(nextBlockHeight), int64(iterations), blsSigner, time.Now())

	s.LogBuffer.Reset()
	var beforeWithdrawalBalance hexutil.Big
	err = s.TestNode.EngineClient.Call(s.CtxApp, &beforeWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
	s.T().Logf("Balance before withdrawal processed: %s", beforeWithdrawalBalance.ToInt().String())

	// Balance should not have changed yet
	s.Require().Equal(afterRequestBalance.ToInt().String(), beforeWithdrawalBalance.ToInt().String())

	// The next block will be the turn of the Epoch, and balance will change
	nextBlockHeight = int(iterations) + nextBlockHeight
	proposals, _, _ = s.MoveChainToHeight(s.T(), int64(nextBlockHeight), 1, blsSigner, time.Now())

	var afterWithdrawalBalance hexutil.Big
	err = s.TestNode.EngineClient.Call(s.CtxApp, &afterWithdrawalBalance, "eth_getBalance", simulated.WithdrawalExecutionAddress, "latest")
	s.T().Logf("Balance after withdrawal processed: %s", afterWithdrawalBalance.ToInt().String())

	// Since this is a full withdrawal, the full balance will be withdrawn.
	// The validator started with a balance equal to math.Gwei(chainSpec.MaxEffectiveBalance())
	withdrawalAmountWei := new(big.Int).Mul(big.NewInt(int64(s.TestNode.ChainSpec.MaxEffectiveBalance())), big.NewInt(params.GWei))

	// Expected balance is balance before withdrawal + withdrawalAmount
	expectedBalance := new(big.Int).Add(beforeWithdrawalBalance.ToInt(), withdrawalAmountWei)

	// The new balance of the validator is updated
	s.Require().Equal(expectedBalance.String(), afterWithdrawalBalance.ToInt().String())

}
