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

package injected_test

import (
	"context"
	"os"
	"testing"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/testing/injected"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type CustomCometComponent struct {
	suite.Suite
	injected.TestSuiteHandle
	NoopComet         *injected.NoopCometService
	BlockchainService *blockchain.Service
}

// TestCustomCometComponent is a test suite with a custom comet driver can can control ourselves
//
//nolint:paralleltest // cannot be run in parallel due to use of environment variables.
func TestCustomCometComponent(t *testing.T) {
	suite.Run(t, new(CustomCometComponent))
}

func (s *CustomCometComponent) SetupTest() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.Ctx = ctx
	s.CancelFunc = cancelFunc

	tempHomeDir := s.T().TempDir()
	// Initialize home directory
	cometConfig := injected.InitializeHomeDir(s.T(), tempHomeDir)

	// Start the Geth node - needs to be done first as we need the auth rpc as input for the beacon node.
	elNode := injected.NewGethNode(tempHomeDir, injected.ValidGethImage())
	gethHandle, authRPC := elNode.Start(s.T())
	s.ElHandle = gethHandle

	// Build the Beacon node once we have the auth rpc url
	logger := phuslu.NewLogger(os.Stdout, &phuslu.Config{
		LogLevel: "DEBUG",
	})

	components := injected.FixedComponents(s.T())
	components = append(components, injected.ProvideNoopCometService)

	testNode := injected.NewTestNode(s.T(),
		injected.TestNodeInput{
			TempHomeDir: tempHomeDir,
			CometConfig: cometConfig,
			AuthRPC:     authRPC,
			Logger:      logger,
			AppOpts:     viper.New(),
			Components:  components,
		})
	s.TestNode = testNode

	// Fetch services we will want to query and interact with so they are easily accessible in testing
	var noopCometService *injected.NoopCometService
	err := testNode.FetchService(&noopCometService)
	s.Require().NoError(err)
	s.NotNil(noopCometService)
	s.NoopComet = noopCometService

	var blockchainService *blockchain.Service
	err = testNode.FetchService(&blockchainService)
	s.Require().NoError(err)
	s.NotNil(blockchainService)
	s.BlockchainService = blockchainService
}

func (s *CustomCometComponent) TearDownTest() {
	err := s.ElHandle.Close()
	if err != nil {
		s.T().Error("Error closing geth handle")
	}
	s.CancelFunc()
}

func (s *CustomCometComponent) TestInitChain_InvalidChainID_MustError() {
	_, err := s.NoopComet.Comet.InitChain(s.Ctx, &types.InitChainRequest{
		ChainId: "henlo-im-invalid",
	})
	s.Require().ErrorContains(err, "invalid chain-id on InitChain; expected: test-mainnet-chain, got: henlo-im-invalid")
}

func (s *CustomCometComponent) TestProcessProposal_HeightZero_MustError() {
	_, err := s.NoopComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Height: 0,
	})
	s.Require().ErrorContains(err, "processProposal at height 0: invalid height")
}

func (s *CustomCometComponent) TestProcessProposal_InvalidBlock_MustError() {
	res, err := s.NoopComet.Comet.ProcessProposal(s.Ctx, &types.ProcessProposalRequest{
		Txs:    make([][]byte, 0),
		Height: 1,
	})
	s.Require().NoError(err)
	s.Require().Equal(types.PROCESS_PROPOSAL_STATUS_REJECT, res.Status)
}
