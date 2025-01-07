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

package suite

import (
	"context"
	"fmt"

	"cosmossdk.io/log"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// SetupSuite executes before the test suite begins execution.
func (s *KurtosisE2ESuite) SetupSuite() {
	// Initialize maps if needed
	if s.networks == nil {
		s.networks = make(map[string]*NetworkInstance)
	}
	if s.testSpecs == nil {
		s.testSpecs = make(map[string]string)
	}
	// Setup basic suite configuration
	s.ctx = context.Background()
	s.logger = log.NewTestLogger(s.T())

	var err error
	s.kCtx, err = kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)

	// Initialize only one network for all tests (using devnet)
	chainID, err := s.GetChainIDFromSpec("devnet")
	s.Require().NoError(err)

	s.Logger().Info("Chain ID", "chainID", chainID)
	// Initialize the network once
	err = s.initializeNetwork("devnet", int(chainID))
	s.Require().NoError(err)

}

// Define chain spec to chain ID mapping
var ChainSpecToID = map[string]uint64{
	"devnet":  80087,
	"testnet": 80084,
}

// Add a method to get chain ID from spec
func (s *KurtosisE2ESuite) GetChainIDFromSpec(spec string) (uint64, error) {
	chainID, exists := ChainSpecToID[spec]
	if !exists {
		return 0, fmt.Errorf("unknown chain spec: %s", spec)
	}
	return chainID, nil
}

// SetupTest runs before each test
func (s *KurtosisE2ESuite) SetupTest() {
	testName := s.T().Name()
	s.Logger().Info("Setting up test", "testName", testName)
}

// TearDownSuite cleans up resources after all tests have been executed.
func (s *KurtosisE2ESuite) TearDownSuite() {
	s.Logger().Info("Destroying enclaves...")

	// Clean up all networks
	for chainSpec, network := range s.networks {
		for _, client := range network.consensusClients {
			res, err := client.Stop(s.ctx)
			s.Require().NoError(err, "Error stopping consensus client")
			s.Require().Nil(res.ExecutionError, "Error stopping consensus client")
			s.Require().Empty(res.ValidationErrors, "Error stopping consensus client")
		}

		enclaveName := fmt.Sprintf("e2e-test-enclave-%s", chainSpec)
		s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, enclaveName))
	}
}

// CheckForSuccessfulTx returns true if the transaction was successful.
func (s *KurtosisE2ESuite) CheckForSuccessfulTx(
	tx common.Hash,
) bool {
	ctx, cancel := context.WithTimeout(s.Ctx(), DefaultE2ETestTimeout)
	defer cancel()
	receipt, err := s.JSONRPCBalancer().TransactionReceipt(ctx, tx)
	if err != nil {
		s.Logger().Error("Error getting transaction receipt", "error", err)
		return false
	}
	return receipt.Status == ethtypes.ReceiptStatusSuccessful
}
