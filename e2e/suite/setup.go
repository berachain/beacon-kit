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
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync/atomic"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/e2e/suite/types"
	"github.com/berachain/beacon-kit/kurtosis"
	"github.com/cockroachdb/errors"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/sourcegraph/conc/iter"
	"golang.org/x/sync/errgroup"
)

// SetupSuite executes before the test suite begins execution.
func (s *KurtosisE2ESuite) SetupSuite() {
	s.SetupSuiteWithOptions()
}

// Option is a function that sets a field on the KurtosisE2ESuite.
func (s *KurtosisE2ESuite) SetupSuiteWithOptions(opts ...Option) {
	var (
		key1, key2, key3 *ecdsa.PrivateKey
		err              error
	)

	// Setup some sane defaults.
	s.cfg = kurtosis.DefaultE2ETestConfig()
	s.ctx = context.Background()
	s.logger = log.NewTestLogger(s.T())
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
	// the caller (i.e overriding defaults).
	for _, opt := range opts {
		if err = opt(s); err != nil {
			s.Require().NoError(err, "Error applying option")
		}
	}

	s.kCtx, err = kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)
	s.logger.Info("destroying any existing enclave...")
	//#nosec:G703 // its okay if this errors out. It will error out
	// if there is no enclave to destroy.
	_ = s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave")

	s.logger.Info("creating enclave...")
	s.enclave, err = s.kCtx.CreateEnclave(s.ctx, "e2e-test-enclave")
	s.Require().NoError(err)

	s.logger.Info(
		"spinning up enclave...",
		"num_participants",
		len(s.cfg.Participants),
	)
	result, err := s.enclave.RunStarlarkPackageBlocking(
		s.ctx,
		"../kurtosis",
		starlark_run_config.NewRunStarlarkConfig(
			starlark_run_config.WithSerializedParams(
				string(s.cfg.MustMarshalJSON()),
			),
		),
	)
	s.Require().NoError(err, "Error running Starlark package")
	s.Require().Nil(result.ExecutionError, "Error running Starlark package")
	s.Require().Empty(result.ValidationErrors)

	// Setup the clients and connect.
	s.SetupExecutionClients()

	// Wait for the finalized block number to reach 1.
	err = s.WaitForFinalizedBlockNumber(1)
	s.Require().NoError(err, "Error waiting for finalized block number")

	// Fund any requested accounts.
	s.FundAccounts()
}

// SetupExecutionClients sets up the execution clients for the test suite.
func (s *KurtosisE2ESuite) SetupExecutionClients() {
	s.executionClients = make(map[string]*types.ExecutionClient)
	svrcs, err := s.Enclave().GetServices()
	s.Require().NoError(err, "Error getting services")
	for name, v := range svrcs {
		var serviceCtx *services.ServiceContext
		serviceCtx, err = s.Enclave().GetServiceContext(string(v))
		s.Require().NoError(err, "Error getting service context")
		if strings.HasPrefix(string(name), "el-") {
			if s.executionClients[string(name)],
				err = types.NewExecutionClientFromServiceCtx(
				serviceCtx,
				s.logger,
			); err != nil {
				// TODO: Figoure out how to handle clients that purposefully
				// don't expose JSON-RPC.
				s.Require().NoError(err, "Error creating execution client")
			}
		}
	}
}

// FundAccounts funds the accounts for the test suite.
func (s *KurtosisE2ESuite) FundAccounts() {
	ctx := context.Background()
	nonce := atomic.Uint64{}
	ecKeys := make([]string, 0, len(s.executionClients))
	for key := range s.executionClients {
		ecKeys = append(ecKeys, key)
	}

	// Send ether from the genesis account to the test account
	randomIndex, err := rand.Int(
		rand.Reader, big.NewInt(int64(len(ecKeys))))
	s.Require().NoError(err, "Error generating random index")
	el := s.executionClients[ecKeys[randomIndex.Int64()]]
	pendingNonce, err := el.PendingNonceAt(
		ctx, s.genesisAccount.Address(),
	)
	nonce.Store(pendingNonce)
	s.Require().NoError(err, "Failed to get nonce for genesis account")

	var chainID *big.Int
	chainID, err = el.NetworkID(ctx)
	s.Require().NoError(err, "Failed to get network ID")

	_, err = iter.MapErr(
		s.testAccounts,
		func(acc **types.EthAccount) (*ethtypes.Receipt, error) {
			account := *acc
			// Select a random execution client to send the transaction to.
			// TODO: Filter by RPC support.
			var i *big.Int
			i, err = rand.Int(rand.Reader, big.NewInt(int64(len(ecKeys))))
			if err != nil {
				return nil, err
			}
			executionClient := s.executionClients[ecKeys[i.Int64()]]

			var gasTipCap *big.Int
			if gasTipCap, err = executionClient.SuggestGasTipCap(ctx); err != nil {
				return nil, err
			}

			gasFeeCap := new(big.Int).Add(gasTipCap, big.NewInt(TenGwei))
			nonceToSubmit := nonce.Add(1) - 1
			value := big.NewInt(Ether)
			dest := account.Address()
			var signedTx *ethtypes.Transaction
			if signedTx, err = s.genesisAccount.SignTx(
				chainID, ethtypes.NewTx(&ethtypes.DynamicFeeTx{
					ChainID:   chainID,
					Nonce:     nonceToSubmit,
					GasTipCap: gasTipCap,
					GasFeeCap: gasFeeCap,
					Gas:       EtherTransferGasLimit,
					To:        &dest,
					Value:     value,
					Data:      nil,
				}),
			); err != nil {
				return nil, err
			}

			cctx, cancel := context.WithTimeout(ctx, DefaultE2ETestTimeout)
			defer cancel()

			if err = executionClient.SendTransaction(cctx, signedTx); err != nil {
				return nil, err
			}

			s.logger.Info(
				"funding transaction submitted, waiting for confirmation...",
				"tx_hash", signedTx.Hash().Hex(), "nonce", nonceToSubmit,
				"account", account.Name(), "value", value,
			)

			var receipt *ethtypes.Receipt
			receipt, err = bind.WaitMined(cctx, executionClient, signedTx)
			if err != nil {
				return nil, err
			}

			s.logger.Info(
				"funding transaction confirmed",
				"tx_hash", signedTx.Hash().Hex(),
				"account", account.Name(),
			)

			// Verify the receipt status.
			if receipt.Status != ethtypes.ReceiptStatusSuccessful {
				return nil, err
			}

			// Verify the balance of the account
			var balance *big.Int
			if balance, err = executionClient.BalanceAt(
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
	eg, groupCtx := errgroup.WithContext(context.Background())
	groupCctx, cancel := context.WithTimeout(
		groupCtx, DefaultE2ETestTimeout)
	defer cancel()
	for _, executionClient := range s.ExecutionClients() {
		eg.Go(
			func() error {
				return executionClient.WaitForFinalizedBlockNumber(
					groupCctx,
					target,
				)
			},
		)
	}

	return eg.Wait()
}

// TearDownSuite cleans up resources after all tests have been executed.
// this function executes after all tests executed.
func (s *KurtosisE2ESuite) TearDownSuite() {
	s.Logger().Info("destroying enclave...")
	s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave"))
}
