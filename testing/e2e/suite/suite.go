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
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/crypto"
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

// KurtosisE2ESuite manages the test suite and network instances
type KurtosisE2ESuite struct {
	suite.Suite
	logger log.Logger
	ctx    context.Context
	kCtx   *kurtosis_context.KurtosisContext

	// Network management
	networks  map[string]*NetworkInstance // maps chainSpec to network
	testSpecs map[string]string           // maps testName to chainSpec
	mu        sync.RWMutex
}

// NetworkInstance manages a single network configuration
type NetworkInstance struct {
	enclave          *enclaves.EnclaveContext
	cfg              *config.E2ETestConfig
	consensusClients map[string]*types.ConsensusClient
	loadBalancer     *types.LoadBalancer
	genesisAccount   *types.EthAccount
	testAccounts     []*types.EthAccount
}

// NewKurtosisE2ESuite creates and initializes a new KurtosisE2ESuite
func NewKurtosisE2ESuite() *KurtosisE2ESuite {
	return &KurtosisE2ESuite{
		networks:  make(map[string]*NetworkInstance),
		testSpecs: make(map[string]string),
	}
}

// Update accessor methods to use the current network
func (s *KurtosisE2ESuite) JSONRPCBalancer() *types.LoadBalancer {
	return s.GetNetworkForTest().loadBalancer
}

func (s *KurtosisE2ESuite) ConsensusClients() map[string]*types.ConsensusClient {
	return s.GetNetworkForTest().consensusClients
}

// Ctx returns the context associated with the KurtosisE2ESuite.
// This context is used throughout the suite to control the flow of operations,
// including timeouts and cancellations.
func (s *KurtosisE2ESuite) Ctx() context.Context {
	return s.ctx
}

// Enclave returns the enclave running the beacon-kit network.
func (s *KurtosisE2ESuite) Enclave() *enclaves.EnclaveContext {
	return s.GetNetworkForTest().enclave
}

// // Config returns the E2ETestConfig associated with the KurtosisE2ESuite.
// func (s *KurtosisE2ESuite) Config() *config.E2ETestConfig {
// 	return s.cfg
// }

// KurtosisCtx returns the KurtosisContext associated with the KurtosisE2ESuite.
// The KurtosisContext is a critical component that facilitates interaction with
// the Kurtosis testnet, including creating and managing enclaves.
func (s *KurtosisE2ESuite) KurtosisCtx() *kurtosis_context.KurtosisContext {
	return s.kCtx
}

// ExecutionClients returns the execution clients associated with the
// KurtosisE2ESuite.
func (s *KurtosisE2ESuite) ExecutionClients() map[string]*types.ExecutionClient {
	return nil
}

// // JSONRPCBalancerType returns the type of the JSON-RPC balancer
// // for the test suite.
func (s *KurtosisE2ESuite) JSONRPCBalancerType() string {
	return s.GetNetworkForTest().cfg.EthJSONRPCEndpoints[0].Type
}

// Logger returns the logger for the test suite.
func (s *KurtosisE2ESuite) Logger() log.Logger {
	return s.logger
}

// // GenesisAccount returns the genesis account for the test suite.
// func (s *KurtosisE2ESuite) GenesisAccount() *types.EthAccount {
// 	return s.genesisAccount
// }

// TestAccounts returns the test accounts for the current network
func (s *KurtosisE2ESuite) TestAccounts() []*types.EthAccount {
	return s.GetNetworkForTest().testAccounts
}

// RegisterTest associates a test with its chain specification
func (s *KurtosisE2ESuite) RegisterTest(testName, chainSpec string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.testSpecs[testName] = chainSpec
}

// GetNetworkForTest returns the network instance for the current test
func (s *KurtosisE2ESuite) GetNetworkForTest() *NetworkInstance {
	s.mu.RLock()
	defer s.mu.RUnlock()
	testName := s.T().Name()
	spec := s.testSpecs[testName]
	return s.networks[spec]
}

// initializeNetwork creates a new network instance for a chain spec
func (s *KurtosisE2ESuite) initializeNetwork(chainSpec string) error {
	s.logger.Info("Initializing network", "chainSpec", chainSpec)
	network := &NetworkInstance{
		cfg:              config.DefaultE2ETestConfig(),
		consensusClients: make(map[string]*types.ConsensusClient),
	}

	// Apply chain spec configuration
	network.cfg.NetworkConfiguration.ChainSpec = chainSpec

	// Create unique enclave name for this chain spec
	enclaveName := fmt.Sprintf("e2e-test-enclave-%s", chainSpec)

	var err error
	network.enclave, err = s.kCtx.CreateEnclave(s.ctx, enclaveName)
	if err != nil {
		return err
	}

	// Initialize the network components
	if err := s.setupNetwork(network); err != nil {
		return err
	}

	s.networks[chainSpec] = network
	return nil
}

func (s *KurtosisE2ESuite) setupConsensusClientsForNetwork(network *NetworkInstance) error {
	for i := range network.cfg.NetworkConfiguration.Validators.Nodes {
		clientName := fmt.Sprintf("cl-validator-beaconkit-%d", i)
		sCtx, err := network.enclave.GetServiceContext(clientName)
		if err != nil {
			return err
		}

		client := types.NewConsensusClient(
			types.NewWrappedServiceContext(sCtx, network.enclave.RunStarlarkScriptBlocking),
		)
		network.consensusClients[clientName] = client
	}
	return nil
}
func (s *KurtosisE2ESuite) setupJSONRPCBalancerForNetwork(network *NetworkInstance) error {
	// Get the balancer type from config
	balancerType := network.cfg.EthJSONRPCEndpoints[0].Type

	// Get service context for the balancer
	sCtx, err := network.enclave.GetServiceContext(balancerType)
	if err != nil {
		return err
	}

	// Create new load balancer
	network.loadBalancer, err = types.NewLoadBalancer(sCtx)
	if err != nil {
		return err
	}

	return nil
}

func (s *KurtosisE2ESuite) setupNetwork(network *NetworkInstance) error {
	// Run Starlark package
	result, err := network.enclave.RunStarlarkPackageBlocking(
		s.ctx,
		"../../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithSerializedParams(network.cfg.MustMarshalJSON()),
		),
	)
	if err != nil {
		return err
	}
	if result.ExecutionError != nil {
		return fmt.Errorf("starlark execution error: %s", result.ExecutionError)
	}

	// Setup consensus clients
	if err := s.setupConsensusClientsForNetwork(network); err != nil {
		return err
	}

	// Setup JSON-RPC balancer
	if err := s.setupJSONRPCBalancerForNetwork(network); err != nil {
		return err
	}

	// Initialize accounts
	if err := s.initializeAccountsForNetwork(network); err != nil {
		return err
	}

	return nil
}

func (s *KurtosisE2ESuite) initializeAccountsForNetwork(network *NetworkInstance) error {
	var (
		key0 *ecdsa.PrivateKey
		err  error
	)

	network.genesisAccount = types.NewEthAccountFromHex(
		"genesisAccount", "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
	)
	key0, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")

	network.testAccounts = make([]*types.EthAccount, 1)
	network.testAccounts[0] = types.NewEthAccount("testAccount0", key0)

	return nil
}

// WaitForFinalizedBlockNumber waits for the finalized block number to reach the target.
func (s *KurtosisE2ESuite) WaitForFinalizedBlockNumber(target uint64) error {
	network := s.GetNetworkForTest()
	finalBlockNum, err := network.loadBalancer.BlockNumber(s.ctx)
	if err != nil {
		return err
	}
	if finalBlockNum < target {
		return fmt.Errorf("block number %d not reached (current: %d)", target, finalBlockNum)
	}
	return nil
}

// WaitForNBlockNumbers waits for a specified amount of blocks into the future from now.
func (s *KurtosisE2ESuite) WaitForNBlockNumbers(n uint64) error {
	network := s.GetNetworkForTest()
	current, err := network.loadBalancer.BlockNumber(s.ctx)
	if err != nil {
		return err
	}
	return s.WaitForFinalizedBlockNumber(current + n)
}
