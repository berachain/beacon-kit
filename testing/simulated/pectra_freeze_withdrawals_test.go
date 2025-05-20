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

	"github.com/berachain/beacon-kit/execution/requests/eip7002"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// PectraFreezeWithdrawalsSuite defines our test suite for Pectra with withdrawals frozen, e.g. in emergency scenario.
type PectraFreezeWithdrawalsSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors
}

// TestSimulatedCometComponent runs the test suite.
func TestPectraFreezeWithdrawalsSuite(t *testing.T) {
	suite.Run(t, new(PectraFreezeWithdrawalsSuite))
}

// SetupTest initializes the test environment.
func (s *PectraFreezeWithdrawalsSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-eth-genesis.json"
	chainSpecFunc := simulated.ProvideFreezeWithdrawalsChainSpec
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
func (s *PectraFreezeWithdrawalsSuite) TearDownTest() {
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

func (s *PectraFreezeWithdrawalsSuite) TestWithdrawals_FullWithdrawal_WhileWithdrawalsFrozen() {
	// Initialize the chain and BLS signer
	s.InitializeChain(s.T())
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// Derive the execution address for the withdrawal
	execAddr := simulated.WithdrawalExecutionAddress
	senderAddress := gethcommon.HexToAddress(common.NewExecutionAddressFromHex(execAddr).String())

	var nextBlockHeight int64 = 1

	// 1. Advance one block to update the withdrawal contract's EXCESS_INHIBITOR
	proposals, _, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Unix(nextBlockHeight*2, 0))
	s.Require().Len(proposals, 1)
	nextBlockHeight++

	// 2. Send a full withdrawal request, i.e. 0 value
	totalWithdrawalGwei := beaconmath.Gwei(0)
	{
		senderKey := simulated.GetTestKey(s.T())
		chainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		signer := gethcore.NewPragueSigner(chainID)

		// Fetch the protocol fee
		fee, err := eip7002.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
		s.Require().NoError(err)

		// Build and sign the withdrawal transaction
		txData, err := eip7002.CreateWithdrawalRequestData(blsSigner.PublicKey(), totalWithdrawalGwei)
		s.Require().NoError(err)

		tx := gethcore.MustSignNewTx(senderKey, signer, &gethcore.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			To:        &params.WithdrawalQueueAddress,
			Gas:       500_000,
			GasFeeCap: big.NewInt(1e9),
			GasTipCap: big.NewInt(1e9),
			Value:     fee,
			Data:      txData,
		})

		// Submit the raw transaction
		txBytes, err := tx.MarshalBinary()
		s.Require().NoError(err)

		var result interface{}
		err = s.TestNode.EngineClient.Call(s.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
	}

	// 3. Mine 2 blocks to include the withdrawal in the chain
	{
		iterations := int64(2)
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, iterations, blsSigner, time.Unix(nextBlockHeight*2, 0))
		nextBlockHeight += iterations
	}

	// 4. Advance passed MinValidatorWithdrawabilityDelay so validator can be withdrawn but before withdrawals are re-activated.
	{
		delay := int64(s.TestNode.ChainSpec.MinValidatorWithdrawabilityDelay().Unwrap()) + 2
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, delay, blsSigner, time.Unix(nextBlockHeight*2, 0))
		nextBlockHeight += delay

		validators, apiErr := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight-1), nil, nil)
		s.Require().NoError(apiErr)
		s.Require().Len(validators, 1)
		s.Require().Equal(constants.ValidatorStatusWithdrawalPossible, validators[0].Status)
	}

	// 5. Capture balance before withdrawals re-activate
	var beforeWithdrawBalance *big.Int
	{
		var err error
		beforeWithdrawBalance, err = s.TestNode.ContractBackend.BalanceAt(
			s.CtxApp,
			senderAddress,
			big.NewInt(nextBlockHeight-1),
		)
		s.Require().NoError(err)
		s.T().Logf("Balance before re-activation: %s wei", beforeWithdrawBalance)
	}

	// 6. Mine one block at timestamp=30s to re-activate withdrawals
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Unix(30, 0))
		nextBlockHeight++
	}

	// 7. Capture balance after withdrawal executes
	var afterWithdrawBalance *big.Int
	{
		var err error
		afterWithdrawBalance, err = s.TestNode.ContractBackend.BalanceAt(
			s.CtxApp,
			senderAddress,
			big.NewInt(nextBlockHeight-1),
		)
		s.Require().NoError(err)
		s.T().Logf("Balance after withdrawal: %s wei", afterWithdrawBalance)
	}

	// 8. Ensure the withdrawal amount equals 10_000_000 BERA
	delta := new(big.Int).Sub(afterWithdrawBalance, beforeWithdrawBalance)
	expected := big.NewInt(0).Mul(big.NewInt(10_000_000), big.NewInt(1e18)) // 10 million BERA
	s.Require().Equal(expected, delta, "expected withdrawal of 10_000_000 BERA, got %s", delta)
}

// Test partial withdrawal while withdrawals are frozen
func (s *PectraFreezeWithdrawalsSuite) TestWithdrawals_PartialWithdrawal_WhileWithdrawalsFrozen() {
	// Initialize chain and signer
	s.InitializeChain(s.T())
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// Execution address and sender
	execAddr := simulated.WithdrawalExecutionAddress
	senderAddress := gethcommon.HexToAddress(common.NewExecutionAddressFromHex(execAddr).String())

	var nextBlockHeight int64 = 1

	// 1. Advance one block to bump EXCESS_INHIBITOR
	proposals, _, _ := s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Unix(nextBlockHeight*2, 0))
	s.Require().Len(proposals, 1)
	nextBlockHeight++

	// 2. Submit a partial withdrawal request (e.g. 5 M BERA)
	totalWithdrawalGwei := beaconmath.Gwei(5_000_000 * 1e9)
	{
		senderKey := simulated.GetTestKey(s.T())
		chainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		signer := gethcore.NewPragueSigner(chainID)

		// Fetch withdrawal fee
		fee, err := eip7002.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
		s.Require().NoError(err)

		// Build request data
		txData, err := eip7002.CreateWithdrawalRequestData(blsSigner.PublicKey(), totalWithdrawalGwei)
		s.Require().NoError(err)

		tx := gethcore.MustSignNewTx(senderKey, signer, &gethcore.DynamicFeeTx{
			ChainID:   chainID,
			Nonce:     0,
			To:        &params.WithdrawalQueueAddress,
			Gas:       500_000,
			GasFeeCap: big.NewInt(1e9),
			GasTipCap: big.NewInt(1e9),
			Value:     fee,
			Data:      txData,
		})

		// Send tx
		txBytes, err := tx.MarshalBinary()
		s.Require().NoError(err)

		var result interface{}
		err = s.TestNode.EngineClient.Call(s.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
		s.Require().NoError(err)
	}

	// 3. Mine 2 blocks to include the request
	{
		iters := int64(2)
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, iters, blsSigner, time.Unix(nextBlockHeight*2, 0))
		nextBlockHeight += iters
	}

	// 4. Advance past withdrawability delay but before re-activation
	{
		delay := int64(s.TestNode.ChainSpec.MinValidatorWithdrawabilityDelay().Unwrap()) + 2
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, delay, blsSigner, time.Unix(nextBlockHeight*2, 0))
		nextBlockHeight += delay

		// Validator should be withdrawable
		validators, apiErr := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight-1), nil, nil)
		s.Require().NoError(apiErr)
		s.Require().Len(validators, 1)
		//s.Require().Equal(constants.ValidatorStatusWithdrawalPossible, validators[0].Status)
	}

	// 5. Get balance before re-activation
	var beforeBalance *big.Int
	{
		var err error
		beforeBalance, err = s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		s.T().Logf("Balance before re-activation: %s wei", beforeBalance)
	}

	// 6. Activate withdrawals at 30s
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Unix(30, 0))
		nextBlockHeight++
	}

	// 7. Get balance after partial withdrawal
	var afterBalance *big.Int
	{
		var err error
		afterBalance, err = s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		s.T().Logf("Balance after withdrawal: %s wei", afterBalance)
	}

	// 8. Verify withdrawal equals 5M BERA
	delta := new(big.Int).Sub(afterBalance, beforeBalance)
	expected := new(big.Int).Mul(big.NewInt(5_000_000), big.NewInt(1e18))
	s.Require().Equal(expected.String(), delta.String(), "expected withdrawal of %d wei, got %s", expected.String(), delta)
}
