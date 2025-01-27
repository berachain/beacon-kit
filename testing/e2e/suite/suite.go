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
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/log"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	types "github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
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
	genesisAccount   *types.EthAccount
	testAccounts     []*types.EthAccount
	enclave          *enclaves.EnclaveContext
}

// NewNetworkInstance creates a new network instance.
func NewNetworkInstance(cfg *config.E2ETestConfig) *NetworkInstance {
	return &NetworkInstance{
		Config:           cfg,
		consensusClients: make(map[string]*types.ConsensusClient),
	}
}

// GetCurrentNetwork returns the network for the current test.
func (s *KurtosisE2ESuite) GetCurrentNetwork() *NetworkInstance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	testName := s.T().Name()
	// Extract the actual test name from the full path
	if idx := strings.LastIndex(testName, "/"); idx != -1 {
		testName = testName[idx+1:]
	}

	// s.Logger().Info("Getting network for test",
	// 	"testName", testName,
	// 	"networks", s.networks)

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

// GetAccounts returns the test accounts for the test suite.
func (s *KurtosisE2ESuite) GetAccounts() []*types.EthAccount {
	return s.testAccounts
}

// RegisterTest associates a test with its chain specification.
func (s *KurtosisE2ESuite) RegisterTest(testName string, spec ChainSpec) {
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
func (s *KurtosisE2ESuite) SetTestSpecs(specs map[string]ChainSpec) {
	s.testSpecs = specs
}

// Networks returns the networks for the test suite.
func (s *KurtosisE2ESuite) Networks() map[string]*NetworkInstance {
	return s.networks
}

// TestSpecs returns the test specs for the test suite.
func (s *KurtosisE2ESuite) GetTestSpecs() map[string]ChainSpec {
	return s.testSpecs
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
		return fmt.Errorf("network instance cannot be nil")
	}
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
				if err := s.kCtx.DestroyEnclave(s.ctx, e.GetEnclaveUuid()); err != nil {
					s.Logger().Error("Failed to destroy existing enclave", "error", err)
				}
			}
		}
	}

	network.enclave, err = s.kCtx.CreateEnclave(s.ctx, enclaveName)
	if err != nil {
		return fmt.Errorf("failed to create enclave: %w", err)
	}
	s.Logger().Info("Created enclave")

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
		var sCtx *services.ServiceContext
		clientName := fmt.Sprintf("cl-validator-beaconkit-%d", i)
		sCtx, err = network.enclave.GetServiceContext(clientName)
		if err != nil {
			return fmt.Errorf("failed to get service context: %w", err)
		}

		client := types.NewConsensusClient(
			types.NewWrappedServiceContext(sCtx, network.enclave.RunStarlarkScriptBlocking),
		)
		// Connect the client
		if err = client.Connect(s.ctx); err != nil {
			return fmt.Errorf("failed to connect consensus client %s: %w", clientName, err)
		}
		network.consensusClients[clientName] = client
		s.Logger().Info("Created consensus client", "name", clientName)
	}

	// Add this line to update the suite's consensus clients
	s.consensusClients = network.consensusClients
	s.Logger().Info("Set up consensus clients", "clients", s.consensusClients)

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

	// Set the suite's load balancer to match the network's
	s.loadBalancer = network.loadBalancer
	// Wait for RPC to be ready before funding accounts
	if err = s.WaitForRPCReady(network); err != nil {
		return fmt.Errorf("failed waiting for RPC: %w", err)
	}

	// Initialize genesis account
	network.genesisAccount = types.NewEthAccountFromHex(
		"genesisAccount", "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
	)
	s.genesisAccount = network.genesisAccount

	// Wait for a few blocks to ensure the genesis account has funds
	if err = s.WaitForNBlockNumbers(5); err != nil {
		return fmt.Errorf("failed waiting for blocks: %w", err)
	}

	// Verify genesis account balance
	balance, err := s.JSONRPCBalancer().BalanceAt(s.ctx, s.genesisAccount.Address(), nil)
	if err != nil {
		return fmt.Errorf("failed to get genesis balance: %w", err)
	}
	s.Logger().Info("Genesis account balance", "balance", balance)
	if balance.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("genesis account has no funds")
	}

	// Wait for RPC to be ready before funding accounts
	if err = s.WaitForRPCReady(network); err != nil {
		return fmt.Errorf("failed waiting for RPC: %w", err)
	}

	var (
		key0, key1, key2 *ecdsa.PrivateKey
	)
	key0, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")
	key1, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")
	key2, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")

	// Initialize test accounts
	network.testAccounts = []*types.EthAccount{
		types.NewEthAccount("testAccount0", key0),
		types.NewEthAccount("testAccount1", key1),
		types.NewEthAccount("testAccount2", key2),
	}
	s.testAccounts = network.testAccounts

	// Fund test accounts using the genesis account
	for _, account := range network.testAccounts {
		amount, ok := new(big.Int).SetString("20000000000000000000000", 10) // 20000 ETH
		if !ok {
			return fmt.Errorf("failed to parse amount")
		}
		if err = s.FundAccount(account.Address(), amount); err != nil {
			return fmt.Errorf("failed to fund test accounts: %w", err)
		}
	}

	return nil
}

// CleanupNetwork cleans up the network resources.
func (s *KurtosisE2ESuite) CleanupNetwork(network *NetworkInstance) error {
	// Stop consensus clients
	if network == nil || len(network.consensusClients) == 0 {
		// Network already cleaned up
		return nil
	}

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
		// Check if enclave still exists before trying to destroy it
		enclaves, err := s.kCtx.GetEnclaves(s.ctx)
		if err != nil {
			s.Logger().Error("Failed to get enclaves", "error", err)
			return nil // Continue with cleanup even if we can't check enclaves
		}

		enclaveExists := false
		for _, e := range enclaves.GetEnclavesByUuid() {
			if e.GetEnclaveUuid() == string(network.enclave.GetEnclaveUuid()) {
				enclaveExists = true
				break
			}
		}

		if !enclaveExists {
			s.Logger().Info("Enclave already destroyed", "uuid", network.enclave.GetEnclaveUuid())
			return nil
		}

		s.Logger().Info("Destroying enclave in cleanupNetwork", "enclave", network.enclave)
		if err = s.kCtx.DestroyEnclave(s.ctx, string(network.enclave.GetEnclaveUuid())); err != nil {
			return fmt.Errorf("failed to destroy enclave: %w", err)
		}
	}

	return nil
}

func (s *KurtosisE2ESuite) SetKurtosisCtx(ctx *kurtosis_context.KurtosisContext) {
	s.kCtx = ctx
}

// WaitForRPCReady waits for the RPC endpoint to be ready.
func (s *KurtosisE2ESuite) WaitForRPCReady(network *NetworkInstance) error {
	s.Logger().Info("Waiting for RPC to be ready", "url", network.loadBalancer.URL())
	maxRetries := 30
	for i := 0; i < maxRetries; i++ {
		blockNum, err := network.loadBalancer.BlockNumber(s.ctx)
		if err == nil {
			s.Logger().Info("RPC is ready", "blockNum", blockNum)
			return nil
		}
		s.Logger().Info("RPC not ready yet",
			"attempt", i+1,
			"url", network.loadBalancer.URL(),
			"error", err,
		)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("RPC not ready after %d retries", maxRetries)
}

// RegisterTestFunc registers a test function with a name.
func (s *KurtosisE2ESuite) RegisterTestFunc(name string, fn func()) {
	if s.testFuncs == nil {
		s.testFuncs = make(map[string]func())
	}
	s.testFuncs[name] = fn
}

// ConsensusClients returns the consensus clients for this network.
func (n *NetworkInstance) ConsensusClients() map[string]*types.ConsensusClient {
	return n.consensusClients
}

func (n *NetworkInstance) TestAccounts() []*types.EthAccount {
	return n.testAccounts
}

// FundAccount sends ETH to the given address.
func (s *KurtosisE2ESuite) FundAccount(to common.Address, amount *big.Int) error {
	// Get initial balance
	initialBalance, err := s.JSONRPCBalancer().BalanceAt(s.ctx, to, nil)
	if err != nil {
		return fmt.Errorf("failed to get initial balance: %w", err)
	}

	nonce, err := s.JSONRPCBalancer().PendingNonceAt(s.ctx, s.GenesisAccount().Address())
	if err != nil {
		return err
	}

	chainID, err := s.JSONRPCBalancer().ChainID(s.ctx)
	if err != nil {
		return err
	}

	// Get the latest block for fee estimation
	header, err := s.JSONRPCBalancer().HeaderByNumber(s.ctx, nil)
	if err != nil {
		return err
	}
	gasFeeCap := new(big.Int).Add(header.BaseFee, big.NewInt(1e9))

	tx := coretypes.NewTx(&coretypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       21000,
		GasFeeCap: gasFeeCap,
		GasTipCap: big.NewInt(1e9),
	})

	signedTx, err := s.GenesisAccount().SignTx(chainID, tx)
	if err != nil {
		return err
	}

	if err = s.JSONRPCBalancer().SendTransaction(s.ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := s.WaitForTransactionReceipt(signedTx.Hash())
	if err != nil {
		return fmt.Errorf("failed waiting for transaction: %w", err)
	}

	// Verify transaction success
	if receipt.Status != coretypes.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	// Verify balance increase
	newBalance, err := s.JSONRPCBalancer().BalanceAt(s.ctx, to, nil)
	if err != nil {
		return fmt.Errorf("failed to get new balance: %w", err)
	}

	if newBalance.Cmp(initialBalance) <= 0 {
		return fmt.Errorf("balance did not increase: old=%s new=%s", initialBalance, newBalance)
	}

	s.Logger().Info("Successfully funded account",
		"address", to.Hex(),
		"amount", amount,
		"oldBalance", initialBalance,
		"newBalance", newBalance,
	)
	return nil
}

// WaitForTransactionReceipt waits for a transaction to be mined and returns the receipt.
func (s *KurtosisE2ESuite) WaitForTransactionReceipt(hash common.Hash) (*coretypes.Receipt, error) {
	for range 30 {
		receipt, err := s.JSONRPCBalancer().TransactionReceipt(s.ctx, hash)
		if err == nil {
			return receipt, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("transaction not mined within timeout: %s", hash.Hex())
}
