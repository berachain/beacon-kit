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

package e2e_test

import (
	"context"
	"math/big"
	"strings"
	"time"

	"cosmossdk.io/log"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/kurtosis"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/enclaves"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/starlark_run_config"
	"github.com/kurtosis-tech/kurtosis/api/golang/engine/lib/kurtosis_context"
	"github.com/sourcegraph/conc"
	"github.com/stretchr/testify/suite"
)

// BeaconE2ESuite is a suite of tests simluating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.Suite
	cfg    *kurtosis.E2ETestConfig
	logger log.Logger
	ctx    context.Context
	kCtx   *kurtosis_context.KurtosisContext

	// enclave is the enclave running the beacon-kit network
	// that is currently under test.
	enclave *enclaves.EnclaveContext
}

// this function executes before the test suite begins execution.
func (s *BeaconKitE2ESuite) SetupSuite() {
	s.cfg = kurtosis.DefaultE2ETestConfig()
	s.ctx = context.Background()
	s.logger = log.NewTestLogger(s.T())
	var err error
	s.kCtx, err = kurtosis_context.NewKurtosisContextFromLocalEngine()
	s.Require().NoError(err)
	s.logger.Info("destroying any existing enclave...")
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
}

// this function executes after all tests executed.
func (s *BeaconKitE2ESuite) TearDownSuite() {
	s.logger.Info("destroying enclave...")
	s.Require().NoError(s.kCtx.DestroyEnclave(s.ctx, "e2e-test-enclave"))
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
//
//nolint:gocognit // todo break this function into smaller functions.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	targetBlock := uint64(5)
	svrcs, err := s.enclave.GetServices()
	s.Require().NoError(err, "Error getting services")

	var jsonRPCPorts = make([]string, 0)
	for k, v := range svrcs {
		s.logger.Info("Service started", "service", k, "uuid", v)
		var serviceCtx *services.ServiceContext
		serviceCtx, err = s.enclave.GetServiceContext(string(v))
		s.Require().NoError(err, "Error getting service context")

		// Get the public ports representing eth JSON-RPC endpoints.
		jsonRPC, ok := serviceCtx.GetPublicPorts()["rpc"]
		if !ok {
			continue
		}
		jsonRPCPorts = append(
			jsonRPCPorts,
			strings.Split(jsonRPC.String(), "/")[0],
		)
	}

	wg := conc.NewWaitGroup()
	for port := range jsonRPCPorts {
		wg.Go(
			func() {
				var ethClient *ethclient.Client
				ethClient, err = ethclient.Dial(
					"http://0.0.0.0:" + jsonRPCPorts[port],
				)
				s.Require().NoError(err, "Error creating eth client")
				_, err = ethClient.NetworkID(s.ctx)
				s.Require().NoError(err, "Error getting network id")

				for {
					ticker := time.NewTicker(3 * time.Second)
					defer ticker.Stop()

					for {
						select {
						case <-s.ctx.Done():
							return
						case <-ticker.C:
							var currentBlock *ethtypes.Block
							finalizedBlock, err = ethClient.BlockByNumber(
								s.ctx,
								big.NewInt(rpc.FinalizedBlockNumber),
							)
							s.Require().
								NoError(err, "Error getting current block number during wait")
							s.logger.Info(
								"Current block number",
								"block",
								currentBlock.Number().Uint64(),
								"port",
								jsonRPCPorts[port],
							)
							if currentBlock.Number().Uint64() >= targetBlock {
								return
							}
						}
					}
				}
			},
		)
	}

	done := make(chan bool, 1)
	go func() {
		wg.WaitAndRecover()
		done <- true
	}()
	select {
	case <-done:
		// Completed without timeout
	case <-time.After(60 * time.Second): // Adjust the timeout duration as needed
		s.logger.Error("Timeout waiting for goroutines to finish")
	}
}
