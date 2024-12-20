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
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/testing/play_e2e/config"
	"github.com/berachain/beacon-kit/testing/play_e2e/suite/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
)

// SetupSuite executes before the test suite begins execution.
func (s *KurtosisE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions()
}

// SetupSuiteWithOptions sets up the test suite with the provided options.
func (s *KurtosisE2ESuite) SetupSuiteWithOptions(opts ...Option) {
	var (
		key1 *ecdsa.PrivateKey
		err  error
	)

	// Setup some sane defaults.
	s.cfg = config.DefaultE2ETestConfig()
	s.ctx = context.Background()
	s.logger = log.NewTestLogger(s.T())
	s.Require().NoError(err, "Error loading starlark helper file")
	s.testAccounts = make([]*types.EthAccount, 0)

	s.genesisAccount = types.NewEthAccountFromHex(
		"genesisAccount",
		"fffdbb37105441e14b0ee6330d855d8504ff39e705c3afa8f859ac9865f99306",
	)
	key1, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")

	s.testAccounts = append(s.testAccounts, types.NewEthAccount(
		"testAccount1",
		key1,
	))

	// Apply all the provided options, this allows
	// the test suite to be configured in a flexible manner by
	// the caller (i.e. overriding defaults).
	for _, opt := range opts {
		if err = opt(s); err != nil {
			s.Require().NoError(err, "Error applying option")
		}
	}

	s.kCtx, err = kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)
	s.logger.Info("Destroying any existing enclave...")
	//#nosec:G703 // It's okay if this errors out. It will error out
	// if there is no enclave to destroy.
	_ = s.kCtx.DestroyEnclave(s.ctx, "play-test-enclave")

	s.logger.Info("Creating enclave...")
	s.enclave, err = s.kCtx.CreateEnclave(s.ctx, "play-test-enclave")
	s.Require().NoError(err)

	s.logger.Info(
		"Spinning up enclave...",
		"num_validators",
		len(s.cfg.NetworkConfiguration.Validators.Nodes),
		"num_full_nodes",
		len(s.cfg.NetworkConfiguration.FullNodes.Nodes),
	)

	result, err := s.enclave.RunStarlarkPackageBlocking(
		s.ctx,
		"../../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithSerializedParams(
				string(s.cfg.MustMarshalJSON()),
			),
		),
	)
	s.Require().NoError(err, "Error running Starlark package")
	s.Require().Nil(result.ExecutionError, "Error running Starlark package")
	s.Require().Empty(result.ValidationErrors)
	s.logger.Info("Enclave spun up successfully")

	s.logger.Info("Setting up consensus clients")
	err = s.SetupConsensusClients()
	s.Require().NoError(err, "Error setting up consensus clients")

	// Setup the JSON-RPC balancer.
	s.logger.Info("Setting up JSON-RPC balancer")
	err = s.SetupJSONRPCBalancer()
	s.Require().NoError(err, "Error setting up JSON-RPC balancer")

	s.logger.Info("Waiting for nodes to get ready...")
	//nolint:mnd // its okay.
	time.Sleep(5 * time.Second)
	// Wait for the finalized block number to increase a bit to
	// ensure all clients are in sync.
	//nolint:mnd // 3 blocks
	err = s.WaitForFinalizedBlockNumber(2)
	s.Require().NoError(err, "Error waiting for finalized block number")

}

func (s *KurtosisE2ESuite) SetupConsensusClients() error {
	s.consensusClients = make(map[string]*types.ConsensusClient)
	sCtx, err := s.Enclave().GetServiceContext("cl-validator-beaconkit-0")
	if err != nil {
		return err
	}

	s.consensusClients["cl-validator-beaconkit-0"] = types.NewConsensusClient(
		types.NewWrappedServiceContext(
			sCtx,
			s.Enclave().RunStarlarkScriptBlocking,
		),
	)

	sCtx, err = s.Enclave().GetServiceContext("cl-validator-beaconkit-1")
	if err != nil {
		return err
	}
	s.consensusClients["cl-validator-beaconkit-1"] = types.NewConsensusClient(
		types.NewWrappedServiceContext(
			sCtx,
			s.Enclave().RunStarlarkScriptBlocking,
		),
	)

	sCtx, err = s.Enclave().GetServiceContext("cl-validator-beaconkit-2")
	if err != nil {
		return err
	}
	s.consensusClients["cl-validator-beaconkit-2"] = types.NewConsensusClient(
		types.NewWrappedServiceContext(
			sCtx,
			s.Enclave().RunStarlarkScriptBlocking,
		),
	)

	sCtx, err = s.Enclave().GetServiceContext("cl-validator-beaconkit-3")
	if err != nil {
		return err
	}
	s.consensusClients["cl-validator-beaconkit-3"] = types.NewConsensusClient(
		types.NewWrappedServiceContext(
			sCtx,
			s.Enclave().RunStarlarkScriptBlocking,
		),
	)
	return nil
}

// SetupJSONRPCBalancer sets up the load balancer for the test suite.
func (s *KurtosisE2ESuite) SetupJSONRPCBalancer() error {
	// get the type for EthJSONRPCEndpoint
	typeRPCEndpoint := s.JSONRPCBalancerType()

	s.logger.Info("Setting up JSON-RPC balancer:", "type", typeRPCEndpoint)

	sCtx, err := s.Enclave().GetServiceContext(typeRPCEndpoint)
	if err != nil {
		return err
	}

	if s.loadBalancer, err = types.NewLoadBalancer(
		sCtx,
	); err != nil {
		return err
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
//func (s *KurtosisE2ESuite) TearDownSuite() {
//	s.Logger().Info("Destroying enclave...")
//	s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave"))
//}

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
