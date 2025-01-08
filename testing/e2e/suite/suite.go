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
	"reflect"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	e2etypes "github.com/berachain/beacon-kit/testing/e2e/types"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
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
	cfg     *config.E2ETestConfig
	logger  log.Logger
	ctx     context.Context
	kCtx    *kurtosis_context.KurtosisContext
	enclave *enclaves.EnclaveContext

	// TODO: Figure out what these may be useful for.
	consensusClients map[string]*types.ConsensusClient
	// executionClients map[string]*types.ExecutionClient
	loadBalancer *types.LoadBalancer

	genesisAccount *types.EthAccount
	testAccounts   []*types.EthAccount

	// Network management
	networks  map[string]*NetworkInstance   // maps chainSpec to network
	testSpecs map[string]e2etypes.ChainSpec // maps testName to chainSpec
	mu        sync.RWMutex
}

// NetworkInstance represents a single network configuration
type NetworkInstance struct {
	Config           *config.E2ETestConfig
	consensusClients map[string]*types.ConsensusClient
	loadBalancer     *types.LoadBalancer
	genesisAccount   *types.EthAccount
	testAccounts     []*types.EthAccount
	enclave          *enclaves.EnclaveContext
}

// NewNetworkInstance creates a new network instance
func NewNetworkInstance(cfg *config.E2ETestConfig) *NetworkInstance {
	return &NetworkInstance{
		Config:           cfg,
		consensusClients: make(map[string]*types.ConsensusClient),
	}
}

// GetCurrentNetwork returns the network for the current test
func (s *KurtosisE2ESuite) GetCurrentNetwork() *NetworkInstance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	testName := s.T().Name()
	spec := s.testSpecs[testName]
	chainKey := fmt.Sprintf("%d-%s", spec.ChainID, spec.Network)
	return s.networks[chainKey]
}

// ConsensusClients returns the consensus clients associated with the
// KurtosisE2ESuite.
func (
	s *KurtosisE2ESuite,
) ConsensusClients() map[string]*types.ConsensusClient {
	return s.consensusClients
}

// Ctx returns the context associated with the KurtosisE2ESuite.
// This context is used throughout the suite to control the flow of operations,
// including timeouts and cancellations.
func (s *KurtosisE2ESuite) Ctx() context.Context {
	return s.ctx
}

// Enclave returns the enclave running the beacon-kit network.
func (s *KurtosisE2ESuite) Enclave() *enclaves.EnclaveContext {
	return s.enclave
}

// Config returns the E2ETestConfig associated with the KurtosisE2ESuite.
func (s *KurtosisE2ESuite) Config() *config.E2ETestConfig {
	return s.cfg
}

// KurtosisCtx returns the KurtosisContext associated with the KurtosisE2ESuite.
// The KurtosisContext is a critical component that facilitates interaction with
// the Kurtosis testnet, including creating and managing enclaves.
func (s *KurtosisE2ESuite) KurtosisCtx() *kurtosis_context.KurtosisContext {
	return s.kCtx
}

// ExecutionClients returns the execution clients associated with the
// KurtosisE2ESuite.
func (
	s *KurtosisE2ESuite,
) ExecutionClients() map[string]*types.ExecutionClient {
	return nil
}

// JSONRPCBalancer returns the JSON-RPC balancer for the test suite.
func (s *KurtosisE2ESuite) JSONRPCBalancer() *types.LoadBalancer {
	s.Logger().Info("JSONRPCBalancer", "balancer", s.loadBalancer)
	return s.loadBalancer
}

// JSONRPCBalancerType returns the type of the JSON-RPC balancer
// for the test suite.
func (s *KurtosisE2ESuite) JSONRPCBalancerType() string {
	return s.cfg.EthJSONRPCEndpoints[0].Type
}

// Logger returns the logger for the test suite.
func (s *KurtosisE2ESuite) Logger() log.Logger {
	return s.logger
}

// GenesisAccount returns the genesis account for the test suite.
func (s *KurtosisE2ESuite) GenesisAccount() *types.EthAccount {
	return s.genesisAccount
}

// TestAccounts returns the test accounts for the test suite.
func (s *KurtosisE2ESuite) TestAccounts() []*types.EthAccount {
	return s.testAccounts
}

// RegisterTest associates a test with its chain specification
func (s *KurtosisE2ESuite) RegisterTest(testName string, spec e2etypes.ChainSpec) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.testSpecs[testName] = spec
}

// SetLogger sets the logger for the test suite.
func (s *KurtosisE2ESuite) SetLogger(l log.Logger) {
	s.logger = l
}

// SetContext sets the context for the test suite.
func (s *KurtosisE2ESuite) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// SetNetworks sets the networks for the test suite.
func (s *KurtosisE2ESuite) SetNetworks(networks map[string]*NetworkInstance) {
	s.networks = networks
}

// SetTestSpecs sets the test specs for the test suite.
func (s *KurtosisE2ESuite) SetTestSpecs(specs map[string]e2etypes.ChainSpec) {
	s.testSpecs = specs
}

// Networks returns the networks for the test suite.
func (s *KurtosisE2ESuite) Networks() map[string]*NetworkInstance {
	return s.networks
}

// TestSpecs returns the test specs for the test suite.
func (s *KurtosisE2ESuite) TestSpecs() map[string]e2etypes.ChainSpec {
	return s.testSpecs
}

// RunTestsByChainSpec runs all tests for each chain spec
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

		// Wait for RPC to be ready
		if err := s.WaitForRPCReady(network); err != nil {
			s.T().Fatalf("Failed waiting for RPC: %v", err)
		}

		// Run all tests for this chain spec
		for _, testName := range tests {
			s.Logger().Info("Running test", "test", testName)
			method := reflect.ValueOf(interface{}(s)).MethodByName(testName)
			s.Logger().Info("Method", "method", method)
			if !method.IsValid() {
				s.T().Errorf("Test method %s not found", testName)
				continue
			}
			s.Run(testName, func() {
				method.Call(nil)
			})
		}
	}
}

// InitializeNetwork sets up a network using the provided configuration
func (s *KurtosisE2ESuite) InitializeNetwork(network *NetworkInstance) error {
	// Create unique enclave name for this chain spec
	s.Logger().Info("Creating enclave", "chainSpec", network.Config.NetworkConfiguration.ChainSpec)
	enclaveName := fmt.Sprintf("e2e-test-enclave-%s", network.Config.NetworkConfiguration.ChainSpec)

	// Try to destroy any existing enclave with the same name
	enclaves, err := s.kCtx.GetEnclaves(s.ctx)
	if err != nil {
		s.Logger().Error("Failed to get enclaves", "error", err)
	} else {
		for _, e := range enclaves.GetEnclavesByUuid() {
			if e.GetName() == enclaveName {
				s.Logger().Info("Destroying existing enclave", "name", enclaveName)
				if err := s.kCtx.DestroyEnclave(s.ctx, string(e.GetEnclaveUuid())); err != nil {
					s.Logger().Error("Failed to destroy existing enclave", "error", err)
				}
			}
		}
	}

	network.enclave, err = s.kCtx.CreateEnclave(s.ctx, enclaveName)
	if err != nil {
		return fmt.Errorf("failed to create enclave: %w", err)
	}
	s.Logger().Info("Created enclave", "enclave", network.enclave)

	// Run Starlark package
	result, err := network.enclave.RunStarlarkPackageBlocking(
		s.ctx,
		"../../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithSerializedParams(network.Config.MustMarshalJSON()),
		),
	)
	s.Logger().Info("Ran starlark package")
	if err != nil {
		return fmt.Errorf("failed to run starlark package: %w", err)
	}
	if result.ExecutionError != nil {
		return fmt.Errorf("starlark execution error: %s", result.ExecutionError)
	}

	// Setup consensus clients
	s.Logger().Info("Setting up validator clients", "clients", network.Config.NetworkConfiguration.Validators.Nodes)
	for i := range network.Config.NetworkConfiguration.Validators.Nodes {
		clientName := fmt.Sprintf("cl-validator-beaconkit-%d", i)
		sCtx, err := network.enclave.GetServiceContext(clientName)
		if err != nil {
			return fmt.Errorf("failed to get service context: %w", err)
		}

		client := types.NewConsensusClient(
			types.NewWrappedServiceContext(sCtx, network.enclave.RunStarlarkScriptBlocking),
		)
		network.consensusClients[clientName] = client
	}
	s.Logger().Info("Set up validator clients", "clients", network.consensusClients)
	// Setup JSON-RPC balancer
	balancerType := network.Config.EthJSONRPCEndpoints[0].Type
	sCtx, err := network.enclave.GetServiceContext(balancerType)
	if err != nil {
		return fmt.Errorf("failed to get balancer service context: %w", err)
	}
	network.loadBalancer, err = types.NewLoadBalancer(sCtx)
	if err != nil {
		return fmt.Errorf("failed to create load balancer: %w", err)
	}
	s.Logger().Info("Verifying load balancer",
		"type", balancerType,
		"ports", sCtx.GetPublicPorts(),
		"balancer", network.loadBalancer,
	)
	// Set the suite's load balancer to match the network's
	s.loadBalancer = network.loadBalancer
	s.Logger().Info("Set suite's load balancer", "balancer", s.loadBalancer)
	return nil
}

// cleanupNetwork cleans up the network resources
func (s *KurtosisE2ESuite) CleanupNetwork(network *NetworkInstance) error {
	// Stop consensus clients
	s.Logger().Info("Stopping consensus clients in cleanupNetwork", "clients", len(network.consensusClients))
	for name, client := range network.consensusClients {
		s.Logger().Info("Stopping consensus client", "name", name)
		if client != nil && client.Client != nil {
			if res, err := client.Stop(s.ctx); err != nil {
				s.Logger().Error("Failed to stop consensus client", "error", err)
			} else if res != nil && res.ExecutionError != nil {
				s.Logger().Error("Client stop returned error", "error", res.ExecutionError)
			}
		}
	}

	// Destroy enclave
	if network.enclave != nil {
		s.Logger().Info("Destroying enclave in cleanupNetwork", "enclave", network.enclave)
		if err := s.kCtx.DestroyEnclave(s.ctx, string(network.enclave.GetEnclaveUuid())); err != nil {
			return fmt.Errorf("failed to destroy enclave: %w", err)
		}
	}

	return nil
}

func (s *KurtosisE2ESuite) SetKurtosisCtx(ctx *kurtosis_context.KurtosisContext) {
	s.kCtx = ctx
}

// WaitForRPCReady waits for the RPC endpoint to be ready
func (s *KurtosisE2ESuite) WaitForRPCReady(network *NetworkInstance) error {
	s.Logger().Info("Waiting for RPC to be ready")
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		blockNum, err := network.loadBalancer.BlockNumber(s.ctx)
		if err == nil {
			s.Logger().Info("RPC is ready", "blockNum", blockNum)
			return nil
		}
		s.Logger().Info("RPC not ready yet", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("RPC not ready after %d retries", maxRetries)
}
