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

package suite

import (
	"context"
	"fmt"
	"sync"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/stretchr/testify/suite"
)

// Run is an alias for suite.Run to help with importing
// in other packages.
//
//nolint:gochecknoglobals // intentionally.
var Run = suite.Run

// KurtosisE2ESuite.
type KurtosisE2ESuite struct {
	suite.Suite
	logger log.Logger
	ctx    context.Context
	kCtx   *kurtosis_context.KurtosisContext

	loadBalancer *types.LoadBalancer

	genesisAccount *types.EthAccount
	testAccounts   []*types.EthAccount

	// Network management
	networks  map[string]*NetworkInstance // maps chainSpec to network
	testSpecs map[string]ChainSpec        // maps testName to chainSpec
	testFuncs map[string]func()           // maps test names to test functions
	mu        sync.RWMutex
}

// NetworkInstance represents a single network configuration.
type NetworkInstance struct {
	Config           *config.E2ETestConfig
	consensusClients map[string]*types.ConsensusClient
	loadBalancer     *types.LoadBalancer
	testAccounts []*types.EthAccount
	enclave      *enclaves.EnclaveContext
}

// NewNetworkInstance creates a new network instance.
func NewNetworkInstance(cfg *config.E2ETestConfig) *NetworkInstance {
	return &NetworkInstance{
		Config:           cfg,
		consensusClients: make(map[string]*types.ConsensusClient),
	}
}

// Logger returns the logger for the test suite.
func (s *KurtosisE2ESuite) Logger() log.Logger {
	return s.logger
}

// RunTestsByChainSpec runs all tests for each chain spec.
func (s *KurtosisE2ESuite) RunTestsByChainSpec() {
	s.Logger().Info("RunTestsByChainSpec", "testSpecs", s.testSpecs)
	// Group tests by chain spec
	testsBySpec := make(map[string][]string)
	for testName, spec := range s.testSpecs {
		chainKey := fmt.Sprintf("%d-%s", spec.ChainID, spec.Network)
		testsBySpec[chainKey] = append(testsBySpec[chainKey], testName)
	}

	// For each chain spec
	for chainKey, tests := range testsBySpec {
		s.Logger().Info("Setting up network for chain spec", "chainKey", chainKey)

		// Initialize network for this chain spec
		network := s.networks[chainKey]
		if err := s.InitializeNetwork(network); err != nil {
			s.T().Fatalf("Failed to initialize network for %s: %v", chainKey, err)
		}

		// Run all tests for this chain spec
		for _, testName := range tests {
			s.Logger().Info("Running test", "test", testName)
			s.Run(testName, func() {
				fn, ok := s.testFuncs[testName]
				if !ok {
					s.T().Errorf("Test method %s not found", testName)
					return
				}
				fn()
			})
		}

		// Clean up network after all tests for this chain spec are done
		if err := s.CleanupNetwork(network); err != nil {
			s.Logger().Error("Failed to cleanup network", "error", err)
		}
	}
}

// InitializeNetwork sets up a network using the provided configuration.
func (s *KurtosisE2ESuite) InitializeNetwork(network *NetworkInstance) error {
	if network == nil {
		return errors.New("network instance cannot be nil")
	}

	if err := s.setupEnclave(network); err != nil {
		return fmt.Errorf("failed to setup enclave: %w", err)
	}

	if err := s.setupConsensusClients(network); err != nil {
		return fmt.Errorf("failed to setup consensus clients: %w", err)
	}

	if err := s.setupLoadBalancer(network); err != nil {
		return fmt.Errorf("failed to setup load balancer: %w", err)
	}

	if err := s.setupAccounts(network); err != nil {
		return fmt.Errorf("failed to setup accounts: %w", err)
	}

	return nil
}

// CleanupNetwork cleans up the network resources.
func (s *KurtosisE2ESuite) CleanupNetwork(network *NetworkInstance) error {
	if network == nil || len(network.consensusClients) == 0 {
		// Network already cleaned up
		s.Logger().Info("Network is nil, skipping cleanup")
		return nil
	}

	if err := s.stopConsensusClients(network); err != nil {
		s.Logger().Error("Failed to stop consensus clients", "error", err)
		// Continue with cleanup even if consensus clients fail to stop
	}

	if err := s.destroyEnclave(network); err != nil {
		return fmt.Errorf("failed to destroy enclave: %w", err)
	}

	return nil
}
