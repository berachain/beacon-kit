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
	"testing"
	"time"

	injectedconsensus "github.com/berachain/beacon-kit/testing/injected-consensus"
	comettypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/suite"
)

type InjectedConsensus struct {
	suite.Suite
	testNode *injectedconsensus.TestNode
}

func (s *InjectedConsensus) SetupTest() {
	s.testNode = injectedconsensus.NewTestNode(s.T())
}

func (s *InjectedConsensus) TearDownTest() {
	// Ensure teardown runs no matter what
	s.testNode.CancelFunc()
}

func (s *InjectedConsensus) TestInitChainRequestsInvalidChainID() {
	request := &comettypes.InitChainRequest{
		ChainId: "80090",
	}
	_, err := s.testNode.CometService.InitChain(s.testNode.Context, request)
	s.Require().Error(err, "invalid chain-id on InitChain; expected: beacond-2061, got: 80090")
}

// TestProcessProposalRequestInvalidBlock tests the scenario where a peer sends us a block with an invalid timestamp.
func (s *InjectedConsensus) TestProcessProposalRequestInvalidBlock() {
	go func() {
		if err := s.testNode.Node.Start(s.testNode.Context); err != nil {
			s.T().Error(err)
		}
	}()

	<-time.After(30 * time.Second)

	// genesis := genesisFromFile(t, testNode.cometConfig.Genesis)

	// genesisFile := testNode.cometConfig.GenesisFile()

	// request := &comettypes.InitChainRequest{
	//	ChainId:       "beacond-2061",
	//	AppStateBytes: genesis.AppState,
	//}
	// fmt.Println(genesis)
	// fmt.Println(genesisFile)
	// response, err := testNode.cometService.InitChain(ctx, request)
	// require.NoError(t, err)
	minimumBlockHeight := int64(2)
	s.Greater(s.testNode.CometService.LastBlockHeight(), minimumBlockHeight)
	// We expect one deposit given the genesis file in 'config/genesis.json'
	// require.Len(t, response.GetValidators(), 1)
}

func TestInjectedConsensus(t *testing.T) {
	suite.Run(t, new(InjectedConsensus))
}
