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

	depositcli "github.com/berachain/beacon-kit/cli/commands/deposit"
	consensustypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/execution/requests/eip7002"
	"github.com/berachain/beacon-kit/geth-primitives/deposit"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	beaconmath "github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethcore "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/validator"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// testPkey corresponds to address 0x56898d1aFb10cad584961eb96AcD476C6826e41E which is prefunded in genesis
const testPkey2 = "9b9bc88a144fff869ae2f4ea8e252f2494d9b52ea1008d0b3537dad27ab489d5"

// PectraWithdrawalSuite defines our test suite for Pectra related work using simulated Comet component.
type PectraWithdrawalSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors
}

// TestSimulatedCometComponent runs the test suite.
func TestPectraWithdrawalSuite(t *testing.T) {
	suite.Run(t, new(PectraWithdrawalSuite))
}

// SetupTest initializes the test environment.
func (s *PectraWithdrawalSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/pectra-fork-genesis.json"
	chainSpecFunc := simulated.ProvidePectraWithdrawalTestChainSpec
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
func (s *PectraWithdrawalSuite) TearDownTest() {
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

// TestExcessValidatorBeforeFork_CorrectlyEvicted verifies that when a validator’s deposit
// exceeds the validator set cap before the Electra fork, it is evicted correctly. The set
// cap is 1.
//
// Scenario timeline:
//   Epoch 1: Move chain by 1 block to include the deposit (deposit store len == 1).
//   Epoch 2: Move chain by 1 block to enqueue the deposit (deposit store len == 2).
//   Epoch 3: Move chain by 1 block → validator status becomes PendingInitialized.
//   Epoch 4: Move chain by 1 block → validator status becomes PendingQueued.
//   Epoch 5: Move chain by 1 block → validator status becomes ExitedUnslashed; withdrawableEpoch is set to 6.
//   —— Electra fork ——
//   Epoch 6: Move chain by 1 block → status is WithdrawalPossible and EL balance is returned immediately.
//   Epoch 7: Move chain by 1 block → status is WithdrawalDone (effective balance = 0).

func (s *PectraWithdrawalSuite) TestExcessValidatorBeforeFork_CorrectlyEvicted() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Send the Deposit
	var senderAddress gethcommon.Address
	depositAmount := beaconmath.Gwei(500_000 * 1e9) // 500K Bera
	{
		depositContractAddress := gethcommon.Address(s.TestNode.ChainSpec.DepositContractAddress())
		depositClient, err := deposit.NewDepositContract(depositContractAddress, s.TestNode.ContractBackend)
		s.Require().NoError(err)
		depositCount, err := depositClient.DepositCount(&bind.CallOpts{
			BlockNumber: big.NewInt(0),
		})
		s.Require().NoError(err)
		s.Require().Equal(uint64(1), depositCount)

		credAddress := common.NewExecutionAddressFromHex("0x56898d1aFb10cad584961eb96AcD476C6826e41E")
		creds := consensustypes.NewCredentialsFromExecutionAddress(credAddress)
		newDepositor := &signer.BLSSigner{PrivValidator: types.NewMockPVWithKeyType(bls12381.KeyType)}
		depositMsg, blsSig, err := depositcli.CreateDepositMessage(
			s.TestNode.ChainSpec,
			newDepositor,
			s.GenesisValidatorsRoot,
			creds,
			depositAmount,
		)
		s.Require().NoError(err)
		err = depositcli.ValidateDeposit(
			s.TestNode.ChainSpec,
			depositMsg.Pubkey,
			depositMsg.Credentials,
			depositMsg.Amount,
			s.GenesisValidatorsRoot,
			blsSig,
		)
		s.Require().NoError(err)

		elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		senderKey, err := crypto.HexToECDSA(testPkey2)
		senderAddress = gethcommon.HexToAddress(credAddress.String())
		s.Require().NoError(err)
		_, err = depositClient.Deposit(&bind.TransactOpts{
			From: senderAddress,
			Signer: func(_ gethcommon.Address, tx *gethcore.Transaction) (*gethcore.Transaction, error) {
				return gethcore.SignTx(
					tx, gethcore.LatestSignerForChainID(elChainID), senderKey,
				)
			},
			Value: big.NewInt(0).Mul(big.NewInt(int64(depositAmount)), big.NewInt(1e9)),
		}, depositMsg.Pubkey[:], depositMsg.Credentials[:], blsSig[:], senderAddress)
		s.Require().NoError(err)
	}

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)
	// Hard fork occurs at t=10, so we start at t=0
	nextBlockTime := time.Unix(0, 0)
	nextBlockHeight := int64(1)

	// [Slot/Epoch 1] Move the chain by 1 block to include the deposit
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		ds := s.TestNode.StorageBackend.DepositStore()
		deposits, err := ds.GetDepositsByIndex(s.CtxApp, 0, uint64(nextBlockHeight)*s.TestNode.ChainSpec.MaxDepositsPerBlock())
		s.Require().NoError(err)
		// There should only be 1 deposit in the deposit store from genesis
		s.Require().Len(deposits, 1)
		nextBlockHeight++
	}
	// [Slot/Epoch 2] Move the chain by 1 block to Enqueue the deposit
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		ds := s.TestNode.StorageBackend.DepositStore()
		deposits, err := ds.GetDepositsByIndex(s.CtxApp, 0, uint64(nextBlockHeight)*s.TestNode.ChainSpec.MaxDepositsPerBlock())
		s.Require().NoError(err)
		// There should be 2 deposits in the deposit store
		s.Require().Len(deposits, 2)
		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1)
		nextBlockHeight++
	}
	// [Slot/Epoch 3] Move the chain by 1 block make the validator pending initialized
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 2)
		s.Require().Equal(validator.PendingInitialized.String(), validators[1].Status)
		nextBlockHeight++
	}
	// [Slot/Epoch 4] Move the chain by 1 block make the validator pending queued
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 2)
		s.Require().Equal(validator.PendingQueued.String(), validators[1].Status)
		nextBlockHeight++
	}

	// [Slot/Epoch 5] Move the chain by 1 block mark the validator as exited
	{
		s.LogBuffer.Reset()
		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 2)
		s.Require().Equal(validator.ExitedUnslashed.String(), validators[1].Status)

		// The validator should withdrawable at Epoch 6 since hard fork has not occurred.
		s.Require().Equal("6", validators[1].Validator.WithdrawableEpoch)
		nextBlockHeight++
	}
	// [Slot/Epoch 6] Move the chain by 1. This block activates the hard fork. No withdrawal delay is expected.
	{
		s.LogBuffer.Reset()

		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 2)
		s.Require().Equal(validator.WithdrawalPossible.String(), validators[1].Status)

		// The validator should withdrawable at Epoch 6 since hard fork has not occurred.
		s.Require().Equal("6", validators[1].Validator.WithdrawableEpoch)

		// Confirm the fork was activated
		s.Require().Contains(s.LogBuffer.String(), "✅  welcome to the electra (0x05000000) fork! 🎉")

		// Confirm the balance change on EL
		previousBlockBalance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		currentBalance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight))
		s.Require().NoError(err)
		depositAmountWei := big.NewInt(0).Mul(big.NewInt(int64(depositAmount)), big.NewInt(1e9))
		expectedBalance := big.NewInt(0).Add(previousBlockBalance, depositAmountWei)
		s.Require().Equal(expectedBalance, currentBalance)
		nextBlockHeight++
	}
	// [Slot/Epoch 7] Move the chain by 1. The effective balance is now 0 and the withdrawal is complete
	{
		s.LogBuffer.Reset()

		_, _, nextBlockTime = s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, nextBlockTime)
		s.Require().Equal(int64(2)*nextBlockHeight, nextBlockTime.Unix())

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 2)
		s.Require().Equal(validator.WithdrawalDone.String(), validators[1].Status)
		nextBlockHeight++
	}
}

// Verifies that excess balance is not withdrawn accidentally if a validator has multiple sources of withdrawals.
// Both the partial withdrawal and the excess balance withdrawal will occur simultaneously in a block.
// Multiple deposits are sent so that Consensus Layer has a higher balance when the withdrawal request is processed.
// This also served as a PoC for a now patched bug (see https://github.com/berachain/beacon-kit/pull/2723).
func (s *PectraWithdrawalSuite) TestWithdrawalFromExcessStake_WithPartialWithdrawal_CorrectAmountWithdrawn() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	credAddress := common.NewExecutionAddressFromHex(simulated.WithdrawalExecutionAddress)
	creds := consensustypes.NewCredentialsFromExecutionAddress(credAddress)
	senderAddress := gethcommon.HexToAddress(credAddress.String())
	// Hard fork occurs at t=10, so we move passed the pectra hard fork
	nextBlockHeight := int64(1)
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		nextBlockHeight++
	}
	// 10 million bera on EL at the start.
	expectedStartBalance, isValid := big.NewInt(0).SetString("10000000000000000000000000", 10)
	s.Require().True(isValid)
	// Confirm the validator's expected start balance in Wei
	{
		startBalance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		s.Require().Equal(expectedStartBalance, startBalance)
		s.T().Logf("balance at start: %s wei", startBalance.String())

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight-1), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1)
		s.T().Logf("staked validator balance at start: %v gwei", validators[0].Validator.EffectiveBalance)
		// Starts with 10000000000000000 gwei / 10 million BERA staked.
		s.Require().Equal("10000000000000000", validators[0].Validator.EffectiveBalance)
	}

	// Send the Deposit and progress 1 block so that the deposit is included in the next block
	depositAmount := beaconmath.Gwei(1_000_000 * 1e9) // 1 million Bera
	{
		// Send Deposit Request
		iterations := int64(2)
		s.defaultDeposit(blsSigner, creds, depositAmount, true)
		s.MoveChainToHeight(s.T(), nextBlockHeight, iterations, blsSigner, time.Now())
		nextBlockHeight += iterations
	}

	// Create the Partial Withdrawal Request for a large amount above the MaxEffectiveBalance.
	// It will be reduced to the maximum possible amount above the MinActivationBalance as part of the processing logic.
	totalWithdrawalAmount := beaconmath.Gwei(15_000_000 * 1e9)
	{
		// corresponds with the funded address in genesis `simulated.WithdrawalExecutionAddress`
		senderKey := simulated.GetTestKey(s.T())

		elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
		pragueSigner := gethcore.NewPragueSigner(elChainID)

		fee, err := eip7002.GetWithdrawalFee(s.CtxApp, s.TestNode.EngineClient)
		s.Require().NoError(err)

		withdrawalTxData, err := eip7002.CreateWithdrawalRequestData(blsSigner.PublicKey(), totalWithdrawalAmount)
		s.Require().NoError(err)

		withdrawalTx := gethcore.MustSignNewTx(senderKey, pragueSigner, &gethcore.DynamicFeeTx{
			ChainID:   elChainID,
			Nonce:     1,
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
	// Move forward two blocks to include in the chain
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 2, blsSigner, time.Now())
		nextBlockHeight += 2
	}

	// Send another deposit
	{
		s.defaultDeposit(blsSigner, creds, depositAmount, false)
	}

	// Move the chain by 1 block to include the deposit
	var balanceAfterDepositTxIncluded *big.Int
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())

		ds := s.TestNode.StorageBackend.DepositStore()
		deposits, err := ds.GetDepositsByIndex(s.CtxApp, 0, uint64(nextBlockHeight)*s.TestNode.ChainSpec.MaxDepositsPerBlock())
		s.Require().NoError(err)
		// There should be 2 deposits in the store
		s.Require().Len(deposits, 2)

		balanceAfterDepositTxIncluded, err = s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight))
		s.Require().NoError(err)
		nextBlockHeight++
	}
	// Move the chain by 1 block to Enqueue the deposit
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())

		ds := s.TestNode.StorageBackend.DepositStore()
		deposits, err := ds.GetDepositsByIndex(s.CtxApp, 0, uint64(nextBlockHeight)*s.TestNode.ChainSpec.MaxDepositsPerBlock())
		s.Require().NoError(err)
		// There should be 3 deposits in the deposit store
		s.Require().Len(deposits, 3)
		// Only 1 active validator
		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1)
		s.T().Logf("staked validator balance: %v gwei", validators[0].Validator.EffectiveBalance)
		nextBlockHeight++
	}
	// Move the chain by 1 block trigger the withdrawal.
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())

		balance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight))
		s.Require().NoError(err)
		s.Require().Equal(balanceAfterDepositTxIncluded, balance)
		// The validator's balance should not have changed yet
		nextBlockHeight++
	}
	// The next block will have the partial withdrawal, but not the excess balance withdrawal and increase the validator's EL balance
	// Before the fix, it would also have the excess balance withdrawal.
	{
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		nextBlockHeight++
	}
	{
		iterations := int64(4)
		s.MoveChainToHeight(s.T(), nextBlockHeight, iterations, blsSigner, time.Now())
		nextBlockHeight += iterations

		finalBalance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		s.T().Logf("balance at end: %s wei", finalBalance.String())
		finalBalanceGwei, convertErr := beaconmath.GweiFromWei(finalBalance)
		s.Require().NoError(convertErr)
		expectedStartBalanceGwei, convertErr := beaconmath.GweiFromWei(expectedStartBalance)
		s.Require().NoError(convertErr)
		s.Require().InDelta(
			uint64(expectedStartBalanceGwei)+
				uint64(
					s.TestNode.ChainSpec.MaxEffectiveBalance()-
						s.TestNode.ChainSpec.MinActivationBalance(),
				),
			uint64(finalBalanceGwei),
			500_000, // maximum 0.0005 BERA or 500000 Gwei delta
		)

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight-1), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1)
		s.T().Logf("staked validator balance at end: %v gwei", validators[0].Validator.EffectiveBalance)
		s.Require().Equal(s.TestNode.ChainSpec.MinActivationBalance().Base10(), validators[0].Validator.EffectiveBalance)
	}
}

func (s *PectraWithdrawalSuite) TestWithdrawalFromExcessStake_HasCorrectWithdrawalAmount() {
	// Initialize the chain state.
	s.InitializeChain(s.T())

	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	credAddress := common.NewExecutionAddressFromHex(simulated.WithdrawalExecutionAddress)
	creds := consensustypes.NewCredentialsFromExecutionAddress(credAddress)
	senderAddress := gethcommon.HexToAddress(credAddress.String())
	// Hard fork occurs at t=10, so we move passed the pectra hard fork
	nextBlockHeight := int64(1)
	{
		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, 1, blsSigner, time.Now())
		nextBlockHeight++
	}
	// 10 million bera on EL at the start.
	expectedStartBalance, isValid := big.NewInt(0).SetString("10000000000000000000000000", 10)
	s.Require().True(isValid)
	// Confirm the validator's expected start balance in Wei
	{
		startBalance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		s.Require().Equal(expectedStartBalance, startBalance)

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight-1), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1)
		// Starts with 10000000000000000 gwei / 10 million BERA staked.
		s.Require().Equal("10000000000000000", validators[0].Validator.EffectiveBalance)
	}

	// Send the Deposit and progress 1 block so that the deposit is included in the next block
	depositAmount := beaconmath.Gwei(10_000 * 1e9) // 10K BERA
	{
		// Send Deposit Requests
		iterations := int64(10)
		s.defaultDepositWithNonce(blsSigner, creds, depositAmount, true, big.NewInt(0))
		for i := 1; i < 30; i++ {
			s.defaultDepositWithNonce(blsSigner, creds, depositAmount, false, big.NewInt(int64(i)))
		}

		s.LogBuffer.Reset()
		s.MoveChainToHeight(s.T(), nextBlockHeight, iterations, blsSigner, time.Now())
		nextBlockHeight += iterations
		// We expect that withdrawals due to excess balance were created
		s.Require().Contains(s.LogBuffer.String(), "expectedWithdrawals: validator withdrawal due to excess balance")

		deposits, err := s.TestNode.StorageBackend.DepositStore().GetDepositsByIndex(
			s.CtxApp, 0, uint64(nextBlockHeight)*s.TestNode.ChainSpec.MaxDepositsPerBlock(),
		)
		s.Require().Len(deposits, 31)
		s.Require().NoError(err)
	}
	// Confirm the validator's end balance on EL is similar to the start balance on EL
	// Confirm the validator's end balance on CL is still at the cap, i.e., 10 Mil BERA.
	{
		endBalance, err := s.TestNode.ContractBackend.BalanceAt(s.CtxApp, senderAddress, big.NewInt(nextBlockHeight-1))
		s.Require().NoError(err)
		finalBalanceGwei, convertErr := beaconmath.GweiFromWei(endBalance)
		s.Require().NoError(convertErr)
		expectedStartBalanceGwei, convertErr := beaconmath.GweiFromWei(expectedStartBalance)
		s.Require().NoError(convertErr)
		s.Require().InDelta(finalBalanceGwei.Unwrap(), expectedStartBalanceGwei.Unwrap(), 2_000_000) // maximum 2_000_000 Gwei delta

		validators, err := s.TestNode.APIBackend.FilteredValidators(beaconmath.Slot(nextBlockHeight-1), nil, nil)
		s.Require().NoError(err)
		s.Require().Len(validators, 1)
		// Ends with 10000000000000000 gwei / 10 million BERA staked.
		s.Require().Equal("10000000000000000", validators[0].Validator.EffectiveBalance)
	}
}

func (s *PectraWithdrawalSuite) defaultDepositWithNonce(
	blsSigner *signer.BLSSigner, creds consensustypes.WithdrawalCredentials, depositAmount beaconmath.Gwei, setOperator bool, nonce *big.Int) {
	depositContractAddress := gethcommon.Address(s.TestNode.ChainSpec.DepositContractAddress())
	depositClient, err := deposit.NewDepositContract(depositContractAddress, s.TestNode.ContractBackend)
	s.Require().NoError(err)

	depositMsg, blsSig, err := depositcli.CreateDepositMessage(
		s.TestNode.ChainSpec,
		blsSigner,
		s.GenesisValidatorsRoot,
		creds,
		depositAmount,
	)
	s.Require().NoError(err)
	err = depositcli.ValidateDeposit(
		s.TestNode.ChainSpec,
		depositMsg.Pubkey,
		depositMsg.Credentials,
		depositMsg.Amount,
		s.GenesisValidatorsRoot,
		blsSig,
	)
	s.Require().NoError(err)

	elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
	senderKey := simulated.GetTestKey(s.T())
	senderAddress := gethcommon.HexToAddress(creds.String())
	s.Require().NoError(err)
	operator := senderAddress
	if !setOperator {
		operator = gethcommon.HexToAddress("0x0000000000000000000000000000000000000000")
	}

	txOpts := &bind.TransactOpts{
		From: senderAddress,
		Signer: func(_ gethcommon.Address, tx *gethcore.Transaction) (*gethcore.Transaction, error) {
			return gethcore.SignTx(
				tx, gethcore.LatestSignerForChainID(elChainID), senderKey,
			)
		},
		GasLimit: 200_000,
		Value:    big.NewInt(0).Mul(big.NewInt(int64(depositAmount)), big.NewInt(1e9)),
	}

	if nonce != nil {
		txOpts.Nonce = nonce
	}

	_, err = depositClient.Deposit(txOpts, depositMsg.Pubkey[:], depositMsg.Credentials[:], blsSig[:], operator)
	s.Require().NoError(err)
}

func (s *PectraWithdrawalSuite) defaultDeposit(blsSigner *signer.BLSSigner, creds consensustypes.WithdrawalCredentials, depositAmount beaconmath.Gwei, setOperator bool) {
	s.defaultDepositWithNonce(blsSigner, creds, depositAmount, setOperator, nil)
}
