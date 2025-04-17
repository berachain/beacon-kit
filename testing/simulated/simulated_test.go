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
	"path"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/berachain/beacon-kit/testing/simulated/execution"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

// SimulatedSuite defines our test suite for the simulated Comet component.
type SimulatedSuite struct {
	suite.Suite
	// Embedded shared accessors for convenience.
	simulated.SharedAccessors
}

// TestSimulatedCometComponent runs the test suite.
func TestSimulatedCometComponent(t *testing.T) {
	suite.Run(t, new(SimulatedSuite))
}

// SetupTest initializes the test environment.
func (s *SimulatedSuite) SetupTest() {
	// Create a cancellable context for the duration of the test.
	s.CtxApp, s.CtxAppCancelFn = context.WithCancel(context.Background())

	// CometBFT uses context.TODO() for all ABCI calls, so we replicate that.
	s.CtxComet = context.TODO()

	s.HomeDir = s.T().TempDir()

	// Initialize the home directory, Comet configuration, and genesis info.
	const elGenesisPath = "./el-genesis-files/eth-genesis.json"
	chainSpecFunc := simulated.ProvideSimulationChainSpec
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

	// Build the Beacon node with the simulated Comet component.
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
func (s *SimulatedSuite) TearDownTest() {
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
