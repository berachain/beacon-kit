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

package injectedconsensus_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/log/phuslu"
	injectedconsensus "github.com/berachain/beacon-kit/testing/injected-consensus"
	comettypes "github.com/cometbft/cometbft/abci/types"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type InjectedConsensus struct {
	suite.Suite
	ctx        context.Context
	cancelFunc context.CancelFunc
	testNode   *injectedconsensus.TestNode

	// Geth dockertest handles for closing
	gethHandle *dockertest.Resource
}

func (s *InjectedConsensus) SetupTest() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancelFunc = cancelFunc

	tempHomeDir := s.T().TempDir()
	// Initialize home directory
	cometConfig := injectedconsensus.InitializeHomeDir(s.T(), tempHomeDir)

	// Start the Geth node - needs to be done first as we need the auth rpc as input for the beacon node.
	elNode := injectedconsensus.NewGethNode(tempHomeDir, injectedconsensus.ValidGethImage())
	gethHandle, authRPC := elNode.Start(s.T())
	s.gethHandle = gethHandle

	// Build the Beacon node once we have the auth rpc url
	logger := phuslu.NewLogger(os.Stdout, nil)
	testNode := injectedconsensus.NewTestNode(s.T(),
		injectedconsensus.TestNodeInput{
			TempHomeDir: tempHomeDir,
			CometConfig: cometConfig,
			AuthRPC:     authRPC,
			Logger:      logger,
		})
	s.testNode = testNode
}

func (s *InjectedConsensus) TearDownTest() {
	err := s.gethHandle.Close()
	if err != nil {
		s.T().Error("Error closing geth handle")
	}
	s.cancelFunc()
}

func (s *InjectedConsensus) TestInitChainRequestsInvalidChainID() {
	request := &comettypes.InitChainRequest{
		ChainId: "80090",
	}
	_, err := s.testNode.CometService.InitChain(s.ctx, request)
	s.Require().ErrorContains(err, "invalid chain-id on InitChain; expected: test-mainnet-chain, got: 80090")
}

func (s *InjectedConsensus) TestInjectedConsensusWorks() {
	go func() {
		if err := s.testNode.Node.Start(s.ctx); err != nil {
			s.T().Error(err)
		}
	}()
	<-time.After(30 * time.Second)
	minimumBlockHeight := int64(2)
	s.Greater(s.testNode.CometService.LastBlockHeight(), minimumBlockHeight)
}

//nolint:paralleltest // cannot be run in parallel due to use of environment variables.
func TestInjectedConsensus(t *testing.T) {
	suite.Run(t, new(InjectedConsensus))
}
