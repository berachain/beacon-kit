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
	"math/big"
	"sync/atomic"
	"time"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/testing/e2e/config"
	"github.com/berachain/beacon-kit/testing/e2e/suite/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/sourcegraph/conc/iter"
)

// SetupSuite executes before the test suite begins execution.
func (s *KurtosisE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions()
}

// SetupSuiteWithOptions sets up the test suite with the provided options.
func (s *KurtosisE2ESuite) SetupSuiteWithOptions(opts ...Option) {
	var (
		key1, key2, key3 *ecdsa.PrivateKey
		err              error
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
	key2, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")
	key3, err = crypto.GenerateKey()
	s.Require().NoError(err, "Error generating key")

	s.testAccounts = append(s.testAccounts, types.NewEthAccount(
		"testAccount1",
		key1,
	))
	s.testAccounts = append(s.testAccounts, types.NewEthAccount(
		"testAccount2",
		key2,
	))
	s.testAccounts = append(s.testAccounts, types.NewEthAccount(
		"testAccount1",
		key3,
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
	_ = s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave")

	s.logger.Info("Creating enclave...")
	s.enclave, err = s.kCtx.CreateEnclave(s.ctx, "e2e-test-enclave")
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
	err = s.WaitForFinalizedBlockNumber(3)
	s.Require().NoError(err, "Error waiting for finalized block number")

	// Fund any requested accounts.
	s.FundAccounts()
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

// FundAccounts funds the accounts for the test suite.
func (s *KurtosisE2ESuite) FundAccounts() {
	ctx := context.Background()
	nonce := atomic.Uint64{}
	pendingNonce, err := s.JSONRPCBalancer().PendingNonceAt(
		ctx, s.genesisAccount.Address())
	s.Require().NoError(err, "Failed to get nonce for genesis account")
	nonce.Store(pendingNonce)

	var chainID *big.Int
	chainID, err = s.JSONRPCBalancer().ChainID(ctx)
	s.Require().NoError(err, "failed to get chain ID")
	s.logger.Info("Chain-id is", "chain_id", chainID)
	_, err = iter.MapErr(
		s.testAccounts,
		func(acc **types.EthAccount) (*ethtypes.Receipt, error) {
			account := *acc
			var gasTipCap *big.Int

			if gasTipCap, err = s.JSONRPCBalancer().SuggestGasTipCap(ctx); err != nil {
				var rpcErr rpc.Error
				if errors.As(err, &rpcErr) && rpcErr.ErrorCode() == -32601 {
					// Besu does not support eth_maxPriorityFeePerGas
					// so we use a default value of 10 Gwei.
					gasTipCap = big.NewInt(0).SetUint64(TenGwei)
				} else {
					return nil, err
				}
			}

			gasFeeCap := new(big.Int).Add(
				gasTipCap, big.NewInt(0).SetUint64(TenGwei))
			nonceToSubmit := nonce.Add(1) - 1
			//nolint:mnd // 20000 Ether
			value := new(big.Int).Mul(big.NewInt(20000), big.NewInt(Ether))
			dest := account.Address()
			var signedTx *ethtypes.Transaction
			if signedTx, err = s.genesisAccount.SignTx(
				chainID, ethtypes.NewTx(&ethtypes.DynamicFeeTx{
					ChainID: chainID, Nonce: nonceToSubmit,
					GasTipCap: gasTipCap, GasFeeCap: gasFeeCap,
					Gas: EtherTransferGasLimit, To: &dest,
					Value: value, Data: nil,
				}),
			); err != nil {
				return nil, err
			}

			cctx, cancel := context.WithTimeout(ctx, DefaultE2ETestTimeout)
			defer cancel()
			if err = s.JSONRPCBalancer().SendTransaction(cctx, signedTx); err != nil {
				s.logger.Error(
					"error submitting funding transaction",
					"error",
					err,
				)
				return nil, err
			}

			s.logger.Info(
				"Funding transaction submitted, waiting for confirmation...",
				"tx_hash", signedTx.Hash().Hex(), "nonce", nonceToSubmit,
				"account", account.Name(), "value", value,
			)

			var receipt *ethtypes.Receipt

			if receipt, err = bind.WaitMined(
				cctx, s.JSONRPCBalancer(), signedTx,
			); err != nil {
				return nil, err
			}

			s.logger.Info(
				"Funding transaction confirmed",
				"tx_hash", signedTx.Hash().Hex(),
				"account", account.Name(),
			)

			// Verify the receipt status.
			if receipt.Status != ethtypes.ReceiptStatusSuccessful {
				return nil, err
			}

			// Wait an extra block to ensure all clients are in sync.
			//nolint:mnd,contextcheck // its okay.
			if err = s.WaitForFinalizedBlockNumber(
				receipt.BlockNumber.Uint64() + 2,
			); err != nil {
				return nil, err
			}

			// Verify the balance of the account
			var balance *big.Int
			if balance, err = s.JSONRPCBalancer().BalanceAt(
				ctx, account.Address(), nil); err != nil {
				return nil, err
			} else if balance.Cmp(value) != 0 {
				return nil, errors.Wrap(
					ErrUnexpectedBalance,
					fmt.Sprintf("expected: %v, got: %v", value, balance),
				)
			}
			return receipt, nil
		},
	)
	s.Require().NoError(err, "Error funding accounts")
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
