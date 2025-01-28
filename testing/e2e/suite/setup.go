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
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// cleanupExistingEnclave attempts to destroy any existing enclave.
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

// setupEnclave creates and initializes the enclave.
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

// setupSingleConsensusClient sets up a single consensus client.
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

// setupConsensusClients initializes and connects consensus clients.
func (s *KurtosisE2ESuite) setupConsensusClients(network *NetworkInstance) error {
	s.Logger().Info("Setting up validator clients", "clients", network.Config.NetworkConfiguration.Validators.Nodes)
	for i := range network.Config.NetworkConfiguration.Validators.Nodes {
		if err := s.setupSingleConsensusClient(network, i); err != nil {
			return err
		}
	}
	s.consensusClients = network.consensusClients
	s.Logger().Info("Set up consensus clients", "clients", s.consensusClients)
	return nil
}

// WaitForRPCReady waits for the RPC endpoint to be ready.
func (s *KurtosisE2ESuite) WaitForRPCReady(network *NetworkInstance) error {
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

// setupLoadBalancer initializes the load balancer.
func (s *KurtosisE2ESuite) setupLoadBalancer(network *NetworkInstance) error {
	balancerType := network.Config.EthJSONRPCEndpoints[0].Type
	sCtx, err := network.enclave.GetServiceContext(balancerType)
	if err != nil {
		return fmt.Errorf("failed to get balancer service context: %w", err)
	}

	network.loadBalancer, err = types.NewLoadBalancer(sCtx)
	if err != nil {
		return fmt.Errorf("failed to create load balancer: %w", err)
	}
	s.loadBalancer = network.loadBalancer

	return s.WaitForRPCReady(network)
}

// generateTestAccounts creates new test accounts with fresh keys.
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
	s.testAccounts = network.testAccounts

	return nil
}

// setupAccounts initializes and funds the accounts.
func (s *KurtosisE2ESuite) setupAccounts(network *NetworkInstance) error {
	// Initialize genesis account
	network.genesisAccount = types.NewEthAccountFromHex(
		"genesisAccount", "fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
	)
	s.genesisAccount = network.genesisAccount

	// Wait for a few blocks to ensure the genesis account has funds
	//nolint:mnd // 5 blocks
	if err := s.WaitForNBlockNumbers(5); err != nil {
		return fmt.Errorf("failed waiting for blocks: %w", err)
	}

	// Verify genesis account balance
	balance, err := s.JSONRPCBalancer().BalanceAt(s.ctx, s.genesisAccount.Address(), nil)
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
		if err = s.FundAccount(account.Address(), amount); err != nil {
			return fmt.Errorf("failed to fund test accounts: %w", err)
		}
	}
	return nil
}

// WaitForFinalizedBlockNumber waits for the finalized block number
// to reach the target block number across all execution clients.
func (s *KurtosisE2ESuite) WaitForFinalizedBlockNumber(
	target uint64,
) error {
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
		finalBlockNum, err = s.JSONRPCBalancer().BlockNumber(cctx)
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

// WaitForNBlockNumbers waits for a specified amount of blocks into the future
// from now.
func (s *KurtosisE2ESuite) WaitForNBlockNumbers(
	n uint64,
) error {
	current, err := s.JSONRPCBalancer().BlockNumber(s.ctx)
	if err != nil {
		return err
	}
	return s.WaitForFinalizedBlockNumber(current + n)
}

// TearDownSuite cleans up resources after all tests have been executed.
// this function executes after all tests executed.
func (s *KurtosisE2ESuite) TearDownSuite() {
	s.Logger().Info("Destroying enclave...")
	for _, client := range s.consensusClients {
		res, err := client.Stop(s.ctx)
		s.Require().NoError(err, "Error stopping consensus client")
		s.Require().Nil(res.ExecutionError, "Error stopping consensus client")
		s.Require().Empty(res.ValidationErrors, "Error stopping consensus client")
	}
	s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave"))
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

// stopSingleConsensusClient stops a single consensus client.
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
func (s *KurtosisE2ESuite) stopConsensusClients(network *NetworkInstance) error {
	s.Logger().Info("Stopping consensus clients in cleanupNetwork", "count", len(network.consensusClients))
	for name, client := range network.consensusClients {
		if err := s.stopSingleConsensusClient(name, client); err != nil {
			return err
		}
	}
	return nil
}

// checkEnclaveExists verifies if the enclave still exists.
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

// JSONRPCBalancer returns the JSON-RPC balancer for the test suite.
func (s *KurtosisE2ESuite) JSONRPCBalancer() *types.LoadBalancer {
	return s.loadBalancer
}

// JSONRPCBalancerType returns the type of the JSON-RPC balancer
// for the test suite.
func (s *KurtosisE2ESuite) JSONRPCBalancerType() string {
	return s.cfg.EthJSONRPCEndpoints[0].Type
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

func (s *KurtosisE2ESuite) SetKurtosisCtx(ctx *kurtosis_context.KurtosisContext) {
	s.kCtx = ctx
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
	if receipt.Status != ethtypes.ReceiptStatusSuccessful {
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
func (s *KurtosisE2ESuite) WaitForTransactionReceipt(hash common.Hash) (*ethtypes.Receipt, error) {
	for range 30 {
		receipt, err := s.JSONRPCBalancer().TransactionReceipt(s.ctx, hash)
		if err == nil {
			return receipt, nil
		}
		time.Sleep(time.Second)
	}
	return nil, fmt.Errorf("transaction not mined within timeout: %s", hash.Hex())
}

// SetContext sets the context for the test suite.
func (s *KurtosisE2ESuite) SetContext(ctx context.Context) {
	s.ctx = ctx
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

// SetLogger sets the logger for the test suite.
func (s *KurtosisE2ESuite) SetLogger(l log.Logger) {
	s.logger = l
}
