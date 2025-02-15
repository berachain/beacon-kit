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
	"context"
	"os"
	"testing"
	"time"

	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/testing/simulated"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

type DefaultConfiguration struct {
	suite.Suite
	simulated.TestSuiteHandle
	CometService *cometbft.Service
}

// TestDefaultConfiguration is a test suite with the default configurations as we would normally build beacond
// It serves as a mechanism
//
//nolint:paralleltest // cannot be run in parallel due to use of environment variables.
func TestDefaultConfiguration(t *testing.T) {
	suite.Run(t, new(DefaultConfiguration))
}

func (s *DefaultConfiguration) SetupTest() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.Ctx = ctx
	s.CancelFunc = cancelFunc

	tempHomeDir := s.T().TempDir()
	// Initialize home directory
	cometConfig := simulated.InitializeHomeDir(s.T(), tempHomeDir)

	// Start the Geth node - needs to be done first as we need the auth rpc as input for the beacon node.
	elNode := simulated.NewGethNode(tempHomeDir, simulated.ValidGethImage())
	gethHandle, authRPC := elNode.Start(s.T())
	s.ElHandle = gethHandle

	// Build the Beacon node once we have the auth rpc url
	logger := phuslu.NewLogger(os.Stdout, nil)
	testNode := simulated.NewTestNode(s.T(),
		simulated.TestNodeInput{
			TempHomeDir: tempHomeDir,
			CometConfig: cometConfig,
			AuthRPC:     authRPC,
			Logger:      logger,
			AppOpts:     viper.New(),
			Components:  simulated.DefaultComponents(s.T()),
		})
	s.TestNode = testNode

	// Fetch services we will want to query and interact with so they are easily accessible in testing
	var cometService *cometbft.Service
	err := testNode.FetchService(&cometService)
	s.Require().NoError(err)
	s.NotNil(cometService)
	s.CometService = cometService
}

func (s *DefaultConfiguration) TearDownTest() {
	err := s.ElHandle.Close()
	if err != nil {
		s.T().Error("Error closing geth handle")
	}
	s.CancelFunc()
}

func (s *DefaultConfiguration) TestDriverWorks() {
	go func() {
		// Node blocks on Start and hence we have to run in separate routine
		if err := s.TestNode.Start(s.Ctx); err != nil {
			s.T().Error(err)
		}
	}()
	<-time.After(10 * time.Second)
	minimumBlockHeight := int64(2)
	s.Greater(s.CometService.LastBlockHeight(), minimumBlockHeight)
}
