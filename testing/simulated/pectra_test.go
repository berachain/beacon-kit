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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// PectraSuite defines our test suite for Pectra related work using simulated Comet component.
type PectraSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors
}

// TestSimulatedCometComponent runs the test suite.
func TestPectraSuite(t *testing.T) {
	suite.Run(t, new(PectraSuite))
}

// SetupTest initializes the test environment.
func (s *PectraSuite) SetupTest() {
	s.T().Parallel()
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./pectra-eth-genesis.json"
	chainSpecFunc := simulated.ProvideElectraGenesisChainSpec
	// Create the chainSpec.
	chainSpec, err := chainSpecFunc()
	s.Require().NoError(err)
	cometConfig, genesisValidatorsRoot := simulated.InitializeHomeDir(s.T(), chainSpec, s.HomeDir, elGenesisPath)
	s.GenesisValidatorsRoot = genesisValidatorsRoot

	// Start the EL (execution layer) Geth node.
	elNode := execution.NewGethNode(s.HomeDir, execution.ValidGethImageWithSimulate())
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
func (s *PectraSuite) TearDownTest() {
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

// TODO(pectra): Get this test case passing. It currently fails in ProcessProposal with error:
// `block hash in payload does not match assembled block`.
func (s *PectraSuite) TestFullLifecycle_WithoutRequests_IsSuccessful() {
	//s.T().Skip("TODO(pectra): Get this test case passing. It currently fails in ProcessProposal with error")
	const blockHeight = 1
	const coreLoopIterations = 10

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// Go through iterations of the core loop.
	proposals := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)
}

// TODO(pectra): Get this test case passing. It currently fails in ProcessProposal with error:
// `block hash in payload does not match assembled block`.
func (s *PectraSuite) TestFullLifecycle_WithRequests_IsSuccessful() {
	s.T().Skip("TODO(pectra): Get this test case passing. It currently fails in ProcessProposal with error")
	const blockHeight = 1
	const coreLoopIterations = 10

	// Initialize the chain state.
	s.InitializeChain(s.T())

	// Retrieve the BLS signer and proposer address.
	blsSigner := simulated.GetBlsSigner(s.HomeDir)

	// create withdrawal request
	// corresponds with funded address in genesis 0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4
	senderKey, err := crypto.HexToECDSA("fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306")
	s.Require().NoError(err)

	elChainID := big.NewInt(int64(s.TestNode.ChainSpec.DepositEth1ChainID()))
	signer := types.NewPragueSigner(elChainID)
	// Field values roughly copies from Geth
	// https://github.com/ethereum/go-ethereum/blob/39638c81c56db2b2dfe6f51999ffd3029ee212cb/core/blockchain_test.go#L4131-L4130
	withdrawalTx := types.MustSignNewTx(senderKey, signer, &types.DynamicFeeTx{
		ChainID:   elChainID,
		Nonce:     1,
		To:        &params.WithdrawalQueueAddress,
		Gas:       500_000,
		GasFeeCap: big.NewInt(1000000000),
		GasTipCap: big.NewInt(1000000000),
		Value:     big.NewInt(1),
		Data:      common.FromHex("b917cfdc0d25b72d55cf94db328e1629b7f4fde2c30cdacf873b664416f76a0c7f7cc50c9f72a3cb84be88144cde91250000000000000d80"),
	})

	txBytes, err := withdrawalTx.MarshalBinary()
	s.Require().NoError(err)

	var result interface{}
	err = s.TestNode.EngineClient.Call(s.CtxApp, &result, "eth_sendRawTransaction", hexutil.Encode(txBytes))
	s.Require().NoError(err)

	// Go through iterations of the core loop.
	proposals := s.MoveChainToHeight(s.T(), blockHeight, coreLoopIterations, blsSigner)
	s.Require().Len(proposals, coreLoopIterations)
}
