// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package suite

import (
	"context"

	"cosmossdk.io/log"
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

// TestAccounts returns the test accounts for the test suite.
func (s *KurtosisE2ESuite) TestAccounts() []*types.EthAccount {
	return s.testAccounts
}
