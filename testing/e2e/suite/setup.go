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
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// Setup related functions

// setupEnclave creates and initializes the enclave for the network.
// It ensures any existing enclave with the same name is cleaned up first.
func (s *KurtosisE2ESuite) setupEnclave(network *NetworkInstance) error {
	// Create unique enclave name for this chain spec
	s.Logger().Info("Creating enclave", "chainSpec", network.Config.NetworkConfiguration.ChainSpec)
	enclaveName := "e2e-test-enclave-" + network.Config.NetworkConfiguration.ChainSpec

	// Try to destroy any existing enclave with the same name
	if err := s.cleanupExistingEnclave(enclaveName); err != nil {
		return err
	}

	var err error
	network.enclave, err = s.kCtx.CreateEnclave(s.ctx, enclaveName)
	if err != nil {
		return fmt.Errorf("failed to create enclave: %w", err)
	}

	// Run Starlark package
	result, err := network.enclave.RunStarlarkPackageBlocking(
		s.ctx,
		"../../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithSerializedParams(network.Config.MustMarshalJSON()),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to run starlark package: %w", err)
	}
	if result.ExecutionError != nil {
		return fmt.Errorf("starlark execution error: %s", result.ExecutionError)
	}

	return nil
}

// setupConsensusClients initializes and connects all consensus clients for the network.
// It creates clients based on the network's validator configuration.
func (s *KurtosisE2ESuite) setupConsensusClients(network *NetworkInstance) error {
	s.Logger().Info("Setting up validator clients", "clients", network.Config.NetworkConfiguration.Validators.Nodes)
	for i := range network.Config.NetworkConfiguration.Validators.Nodes {
		if err := s.setupSingleConsensusClient(network, i); err != nil {
			return err
		}
	}
	s.Logger().Info("Set up consensus clients", "clients", network.consensusClients)
	return nil
}

// setupLoadBalancer initializes the network's load balancer and waits for RPC readiness.
func (s *KurtosisE2ESuite) setupLoadBalancer(network *NetworkInstance) error {
	balancerType := network.Config.EthJSONRPCEndpoints[0].Type
	sCtx, err := network.enclave.GetServiceContext(balancerType)
	if err != nil {
		return fmt.Errorf("failed to get balancer service context: %w", err)
	}

	loadBalancer, err := types.NewLoadBalancer(sCtx)
	if err != nil {
		return fmt.Errorf("failed to create load balancer: %w", err)
	}

	// Verify the load balancer was created successfully
	if loadBalancer == nil {
		return errors.New("load balancer is nil after creation")
	}
	network.loadBalancer = loadBalancer

	return s.waitForRPCReady(network)
}

// setupAccounts initializes and funds the genesis and test accounts.
// It ensures the genesis account has sufficient funds before creating test accounts.
func (s *KurtosisE2ESuite) setupAccounts(network *NetworkInstance) error {
	// Initialize genesis account
	network.genesisAccount = types.NewEthAccountFromHex(
		"genesisAccount", "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
	)

	// Wait for a few blocks to ensure the genesis account has funds
	//nolint:mnd // 5 blocks
	if err := s.WaitForNBlockNumbers(network, 5); err != nil {
		return fmt.Errorf("failed waiting for blocks: %w", err)
	}

	// Verify genesis account balance
	balance, err := network.JSONRPCBalancer().BalanceAt(s.ctx, network.genesisAccount.Address(), nil)
	if err != nil {
		return fmt.Errorf("failed to get genesis balance: %w", err)
	}
	s.Logger().Info("Genesis account balance", "balance", balance)
	if balance.Cmp(big.NewInt(0)) == 0 {
		return errors.New("genesis account has no funds")
	}

	if err = s.generateTestAccounts(network); err != nil {
		return fmt.Errorf("failed to generate test accounts: %w", err)
	}

	// Fund test accounts using the genesis account
	for _, account := range network.testAccounts {
		//nolint:mnd // 60000 ETH
		amount, ok := new(big.Int).SetString("60000000000000000000000", 10)
		if !ok {
			return errors.New("failed to parse amount")
		}
		if err = s.fundAccount(network, account.Address(), amount); err != nil {
			return fmt.Errorf("failed to fund test accounts: %w", err)
		}
	}
	return nil
}

// Helper functions

// cleanupExistingEnclave attempts to destroy any existing enclave with the given name.
// It logs any errors but continues execution to allow setup to proceed.
func (s *KurtosisE2ESuite) cleanupExistingEnclave(enclaveName string) error {
	enclaves, err := s.kCtx.GetEnclaves(s.ctx)
	if err != nil {
		s.Logger().Error("Failed to get enclaves", "error", err)
		return nil // Continue with setup even if we can't check enclaves
	}

	for _, e := range enclaves.GetEnclavesByUuid() {
		if e.GetName() == enclaveName {
			s.Logger().Info("Destroying existing enclave", "name", enclaveName)
			if err = s.kCtx.DestroyEnclave(s.ctx, e.GetEnclaveUuid()); err != nil {
				s.Logger().Error("Failed to destroy existing enclave", "error", err)
			}
		}
	}
	return nil
}

// setupSingleConsensusClient initializes a single consensus client with the given index.
// It creates the client, establishes connection, and adds it to the network's client map.
func (s *KurtosisE2ESuite) setupSingleConsensusClient(network *NetworkInstance, i int) error {
	clientName := fmt.Sprintf("cl-validator-beaconkit-%d", i)
	sCtx, err := network.enclave.GetServiceContext(clientName)
	if err != nil {
		return fmt.Errorf("failed to get service context: %w", err)
	}

	client := types.NewConsensusClient(
		types.NewWrappedServiceContext(sCtx, network.enclave.RunStarlarkScriptBlocking),
	)
	if err = client.Connect(s.ctx); err != nil {
		return fmt.Errorf("failed to connect consensus client %s: %w", clientName, err)
	}
	network.consensusClients[clientName] = client
	s.Logger().Info("Created consensus client", "name", clientName)
	return nil
}

// waitForRPCReady polls the RPC endpoint until it responds successfully or times out.
// It attempts up to maxRetries times with a 2-second delay between attempts.
func (s *KurtosisE2ESuite) waitForRPCReady(network *NetworkInstance) error {
	s.Logger().Info("Waiting for RPC to be ready", "url", network.loadBalancer.URL())
	maxRetries := 30
	for attempt := range maxRetries {
		blockNum, err := network.loadBalancer.BlockNumber(s.ctx)
		if err == nil {
			s.Logger().Info("RPC is ready", "blockNum", blockNum)
			return nil
		}
		s.Logger().Info("RPC not ready yet",
			"attempt", attempt+1,
			"url", network.loadBalancer.URL(),
			"error", err,
		)
		//nolint:mnd // 2 seconds
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("RPC not ready after %d retries", maxRetries)
}

// generateTestAccounts creates three test accounts with freshly generated keys.
// It returns an error if key generation fails for any account.
func (s *KurtosisE2ESuite) generateTestAccounts(network *NetworkInstance) error {
	// Generate keys for test accounts
	//nolint:mnd // 3 accounts
	keys := make([]*ecdsa.PrivateKey, 3)
	for i := range keys {
		var err error
		keys[i], err = crypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("error generating key%d: %w", i, err)
		}
	}

	// Initialize test accounts with generated keys
	network.testAccounts = []*types.EthAccount{
		types.NewEthAccount("testAccount0", keys[0]),
		types.NewEthAccount("testAccount1", keys[1]),
		types.NewEthAccount("testAccount2", keys[2]),
	}

	return nil
}

// fundAccount sends ETH to the given address.
func (s *KurtosisE2ESuite) fundAccount(network *NetworkInstance, to common.Address, amount *big.Int) error {
	// Get initial balance
	initialBalance, err := network.JSONRPCBalancer().BalanceAt(s.ctx, to, nil)
	if err != nil {
		return fmt.Errorf("failed to get initial balance: %w", err)
	}

	nonce, err := network.JSONRPCBalancer().PendingNonceAt(s.ctx, network.genesisAccount.Address())
	if err != nil {
		return err
	}

	chainID, err := network.JSONRPCBalancer().ChainID(s.ctx)
	if err != nil {
		return err
	}

	// Get the latest block for fee estimation
	header, err := network.JSONRPCBalancer().HeaderByNumber(s.ctx, nil)
	if err != nil {
		return err
	}
	//nolint:mnd // 1 Gwei
	gasFeeCap := new(big.Int).Add(header.BaseFee, big.NewInt(1e9))

	//nolint:mnd // 21000 gas
	tx := ethtypes.NewTx(&ethtypes.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		To:        &to,
		Value:     amount,
		Gas:       21000,
		GasFeeCap: gasFeeCap,
		GasTipCap: big.NewInt(1e9),
	})

	signedTx, err := network.genesisAccount.SignTx(chainID, tx)
	if err != nil {
		return err
	}

	if err = network.JSONRPCBalancer().SendTransaction(s.ctx, signedTx); err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for transaction to be mined
	receipt, err := s.WaitForTransactionReceipt(network, signedTx.Hash())
	if err != nil {
		return fmt.Errorf("failed waiting for transaction: %w", err)
	}

	// Verify transaction success
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	// Verify balance increase
	newBalance, err := network.JSONRPCBalancer().BalanceAt(s.ctx, to, nil)
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

// stopSingleConsensusClient stops a single consensus client safely.
// It handles nil clients and logs any errors that occur during shutdown.
func (s *KurtosisE2ESuite) stopSingleConsensusClient(name string, client *types.ConsensusClient) error {
	if client == nil || client.Client == nil {
		s.Logger().Info("Client is nil, skipping", "name", name)
		return nil
	}

	s.Logger().Info("Stopping consensus client", "name", name)
	res, err := client.Stop(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to stop client %s: %w", name, err)
	}

	if res != nil && res.ExecutionError != nil {
		return fmt.Errorf("client %s stop returned error: %s", name, res.ExecutionError)
	}

	return nil
}

// stopConsensusClients stops all consensus clients in the network.
// It attempts to stop each client individually and returns on first error.
func (s *KurtosisE2ESuite) stopConsensusClients(network *NetworkInstance) error {
	s.Logger().Info("Stopping consensus clients in cleanupNetwork", "count", len(network.consensusClients))
	for name, client := range network.consensusClients {
		if err := s.stopSingleConsensusClient(name, client); err != nil {
			return err
		}
	}
	return nil
}

// checkEnclaveExists verifies if the enclave still exists in Kurtosis.
// It returns true if the enclave UUID is found in the list of active enclaves.
func (s *KurtosisE2ESuite) checkEnclaveExists(enclave *enclaves.EnclaveContext) (bool, error) {
	enclaves, err := s.kCtx.GetEnclaves(s.ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get enclaves: %w", err)
	}

	for _, e := range enclaves.GetEnclavesByUuid() {
		if e.GetEnclaveUuid() == string(enclave.GetEnclaveUuid()) {
			return true, nil
		}
	}

	return false, nil
}

// destroyEnclave destroys the network's enclave if it exists.
// It checks for existence first to avoid errors on already destroyed enclaves.
func (s *KurtosisE2ESuite) destroyEnclave(network *NetworkInstance) error {
	if network.enclave == nil {
		s.Logger().Info("Enclave is nil, skipping destruction")
		return nil
	}

	exists, err := s.checkEnclaveExists(network.enclave)
	if err != nil {
		s.Logger().Error("Failed to check enclave existence", "error", err)
		return nil // Continue with cleanup even if we can't check existence
	}

	if !exists {
		s.Logger().Info("Enclave already destroyed", "uuid", network.enclave.GetEnclaveUuid())
		return nil
	}

	s.Logger().Info("Destroying enclave", "uuid", network.enclave.GetEnclaveUuid())
	return s.kCtx.DestroyEnclave(s.ctx, string(network.enclave.GetEnclaveUuid()))
}

// WaitForFinalizedBlockNumber waits until the chain reaches the target block number.
// It polls the chain state and returns an error if the context deadline is exceeded.
func (s *KurtosisE2ESuite) WaitForFinalizedBlockNumber(network *NetworkInstance, target uint64) error {
	cctx, cancel := context.WithTimeout(s.ctx, DefaultE2ETestTimeout)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var finalBlockNum uint64
	for finalBlockNum < target {
		// check if cctx deadline is exceeded to prevent endless loop
		select {
		case <-cctx.Done():
			return cctx.Err()
		default:
		}

		var err error
		finalBlockNum, err = network.JSONRPCBalancer().BlockNumber(cctx)
		if err != nil {
			s.logger.Error("error getting finalized block number", "error", err)
			continue
		}
		s.logger.Info(
			"Waiting for finalized block number to reach target",
			"target",
			target,
			"finalized",
			finalBlockNum,
		)

		select {
		case <-s.ctx.Done():
			return s.ctx.Err()
		case <-ticker.C:
			continue
		}
	}

	s.logger.Info(
		"Finalized block number reached target 🎉",
		"target",
		target,
		"finalized",
		finalBlockNum,
	)

	return nil
}

// WaitForNBlockNumbers waits for a specified amount of blocks into the future from now.
// It gets the current block number and waits until target = current + n blocks.
func (s *KurtosisE2ESuite) WaitForNBlockNumbers(network *NetworkInstance, n uint64) error {
	current, err := network.JSONRPCBalancer().BlockNumber(s.ctx)
	if err != nil {
		return err
	}
	return s.WaitForFinalizedBlockNumber(network, current+n)
}

// GetAccounts returns all test accounts created for the test suite.
func (s *KurtosisE2ESuite) GetAccounts() []*types.EthAccount {
	network := s.GetCurrentNetwork()
	if network == nil {
		s.Logger().Error("No network found for current test")
		return nil
	}
	return network.testAccounts
}

// RegisterTest associates a test with its chain specification.
func (s *KurtosisE2ESuite) RegisterTest(testName string, spec ChainSpec) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.testSpecs[testName] = spec
}

// SetNetworks sets the networks for the test suite.
func (s *KurtosisE2ESuite) SetNetworks(networks map[string]*NetworkInstance) {
	s.networks = networks
}

// Networks returns the networks for the test suite.
func (s *KurtosisE2ESuite) GetNetworks() map[string]*NetworkInstance {
	return s.networks
}

// SetTestSpecs sets the test specs for the test suite.
func (s *KurtosisE2ESuite) SetTestSpecs(specs map[string]ChainSpec) {
	s.testSpecs = specs
}

// TestSpecs returns the test specs for the test suite.
func (s *KurtosisE2ESuite) GetTestSpecs() map[string]ChainSpec {
	return s.testSpecs
}

func (s *KurtosisE2ESuite) SetKurtosisCtx(ctx *kurtosis_context.KurtosisContext) {
	s.kCtx = ctx
}

// SetContext sets the main context for the test suite.
func (s *KurtosisE2ESuite) SetContext(ctx context.Context) {
	s.ctx = ctx
}

// RegisterTestFunc registers a test function with a name.
func (s *KurtosisE2ESuite) RegisterTestFunc(name string, fn func()) {
	if s.testFuncs == nil {
		s.testFuncs = make(map[string]func())
	}
	s.testFuncs[name] = fn
}

// ConsensusClients returns all consensus clients associated with the test suite.
func (s *KurtosisE2ESuite) ConsensusClients() map[string]*types.ConsensusClient {
	network := s.GetCurrentNetwork()
	if network == nil {
		s.Logger().Error("No network found for current test")
		return nil
	}
	return network.consensusClients
}

// WaitForTransactionReceipt waits for a transaction to be mined and returns the receipt.
// It attempts to get the receipt for up to 30 seconds before timing out.
func (s *KurtosisE2ESuite) WaitForTransactionReceipt(network *NetworkInstance, tx common.Hash) (*ethtypes.Receipt, error) {
	for range 30 {
		receipt, err := network.JSONRPCBalancer().TransactionReceipt(s.ctx, tx)
		if err == nil {
			return receipt, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("transaction not mined within timeout: %s", tx.Hex())
}

// GetCurrentNetwork returns the network instance for the current running test.
// It extracts the test name from the full path and looks up the corresponding network.
func (s *KurtosisE2ESuite) GetCurrentNetwork() *NetworkInstance {
	s.mu.RLock()
	defer s.mu.RUnlock()

	testName := s.T().Name()
	// Extract the actual test name from the full path
	if idx := strings.LastIndex(testName, "/"); idx != -1 {
		testName = testName[idx+1:]
	}

	spec := s.testSpecs[testName]
	chainKey := fmt.Sprintf("%d-%s", spec.ChainID, spec.Network)
	return s.networks[chainKey]
}

// Ctx returns the context used throughout the test suite.
func (s *KurtosisE2ESuite) Ctx() context.Context {
	return s.ctx
}

// SetLogger sets the logger instance for the test suite.
func (s *KurtosisE2ESuite) SetLogger(l log.Logger) {
	s.logger = l
}

// Network Instance Methods

// ConsensusClients returns the consensus clients for a specific network instance.
func (n *NetworkInstance) ConsensusClients() map[string]*types.ConsensusClient {
	return n.consensusClients
}

// JSONRPCBalancer returns the JSON-RPC load balancer for the network instance.
func (n *NetworkInstance) JSONRPCBalancer() *types.LoadBalancer {
	return n.loadBalancer
}
