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
	"github.com/berachain/beacon-kit/geth-primitives/deposit"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/node-core/components/signer"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/cometbft/cometbft/crypto/bls12381"
	"github.com/cometbft/cometbft/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
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

// TestExcessValidatorBeforeFork_CorrectlyEvicted checks that if a validator deposits and hits the validator set cap
// before the fork, it will correctly be evicted.
func (s *PectraWithdrawalSuite) TestExcessValidatorBeforeFork_CorrectlyEvicted() {
	const blockHeight = 1
	const coreLoopIterations = 10

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// Hard fork occurs at t=10, so we start at t=0
	startTime := time.Unix(0, 0)

	// Go through iterations of the core loop.
	proposals, _, _ := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner, startTime)
	s.Require().Len(proposals, coreLoopIterations)

	depositContractAddress := gethcommon.Address(s.TestNode.ChainSpec.DepositContractAddress())
	depositClient, err := deposit.NewDepositContract(depositContractAddress, s.TestNode.ContractBackend)
	s.Require().NoError(err)
	depositCount, err := depositClient.DepositCount(&bind.CallOpts{
		BlockNumber: big.NewInt(0),
	})
	s.Require().NoError(err)
	s.Require().Equal(uint64(1), depositCount)

	creds := consensustypes.NewCredentialsFromExecutionAddress(
		common.NewExecutionAddressFromHex("0x56898d1aFb10cad584961eb96AcD476C6826e41E"),
	)
	newDepositor := &signer.BLSSigner{PrivValidator: types.NewMockPVWithKeyType(bls12381.KeyType)}
	depositAmount := math.Gwei(500_000 * 1e9) // 500K Bera
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

	//tx, err = depositClient.Deposit(&bind.TransactOpts{
	//	From:     from,
	//	Value:    value,
	//	Signer:   signer,
	//	Nonce:    new(big.Int).SetUint64(nonce + i),
	//	GasLimit: 1000000,
	//	Context:  s.Ctx(),
	//}, pubkey, curVal.WithdrawalCredentials[:], signature[:], operator)
}
