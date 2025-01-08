// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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

package e2e_test

import (
	"context"
	"fmt"
	"os"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite"
	e2etypes "github.com/berachain/beacon-kit/testing/e2e/types"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// BeaconKitE2ESuite is a suite of tests simulating a fully functional beacon-kit network.
type BeaconKitE2ESuite struct {
	suite.KurtosisE2ESuite
	opts []suite.Option
}

func (s *BeaconKitE2ESuite) SetupSuiteWithOptions(opts ...suite.Option) {
	s.opts = opts
}

// SetupSuite executes before the test suite begins execution.
func (s *BeaconKitE2ESuite) SetupSuite() {
	// Initialize basic configuration
	s.SetContext(context.Background())
	logger := log.NewLogger(os.Stdout)
	s.SetLogger(logger)

	var err error
	kCtx, err := kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)
	s.SetKurtosisCtx(kCtx)

	s.SetNetworks(make(map[string]*suite.NetworkInstance))
	s.SetTestSpecs(make(map[string]e2etypes.ChainSpec))

	// Apply all chain options
	for _, opt := range s.opts {
		if err := opt(&s.KurtosisE2ESuite); err != nil {
			s.Require().NoError(err)
		}
	}

	// Initialize networks
	s.initializeNetworks()
}

// TearDownSuite cleans up after all tests have run
func (s *BeaconKitE2ESuite) TearDownSuite() {
	for _, network := range s.Networks() {
		if err := s.CleanupNetwork(network); err != nil {
			s.Logger().Error("Failed to cleanup network", "error", err)
		}
	}
}

// initializeNetworks sets up networks for each unique chain spec
func (s *BeaconKitE2ESuite) initializeNetworks() {
	for _, spec := range s.GetTestSpecs() {
		chainKey := fmt.Sprintf("%d-%s", spec.ChainID, spec.Network)
		if networks := s.Networks(); networks[chainKey] == nil {
			network := suite.NewNetworkInstance(config.DefaultE2ETestConfig())
			network.Config.NetworkConfiguration.ChainID = int(spec.ChainID)
			network.Config.NetworkConfiguration.ChainSpec = spec.Network
			networks[chainKey] = network
			s.SetNetworks(networks)
		}
	}
}

func (s *BeaconKitE2ESuite) TestRunE2E() {
	s.RunTestsByChainSpec()
}

func (s *BeaconKitE2ESuite) runBasicStartup() {
	err := s.WaitForFinalizedBlockNumber(10)
	s.Require().NoError(err)
}
