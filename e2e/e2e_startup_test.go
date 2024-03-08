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

//nolint:cyclop // todo:fix.
package e2e_test

import (
	"math/big"
	"strings"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/itsdevbear/bolaris/e2e/suite"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/sourcegraph/conc"
)

// BeaconE2ESuite is a suite of tests simulating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.KurtosisE2ESuite
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
//
//nolint:gocognit // todo break this function into smaller functions.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	targetBlock := uint64(5)
	svrcs, err := s.Enclave().GetServices()
	s.Require().NoError(err, "Error getting services")

	var jsonRPCPorts = make(map[string]string)
	for k, v := range svrcs {
		s.Logger().Info("Service started", "service", k, "uuid", v)
		var serviceCtx *services.ServiceContext
		serviceCtx, err = s.Enclave().GetServiceContext(string(v))
		s.Require().NoError(err, "Error getting service context")

		// Get the public ports representing eth JSON-RPC endpoints.
		jsonRPC, ok := serviceCtx.GetPublicPorts()["rpc"]
		if ok {
			str := strings.Split(jsonRPC.String(), "/")
			s.Require().NotNil(str, "Error getting public ports")
			jsonRPCPorts[string(k)] = str[0]
		}
	}

	wg := conc.NewWaitGroup()
	for service, port := range jsonRPCPorts {
		wg.Go(func() {
			s.Logger().Info("Waiting for connection...", "service", service)
			var ethClient *ethclient.Client
			ethClient, err = ethclient.Dial(
				"http://0.0.0.0:" + port,
			)
			s.Require().NoError(err, "Error creating eth client")
			for {
				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-s.Ctx().Done():
						return
					case <-ticker.C:
						var latestBlock *ethtypes.Block
						latestBlock, err = ethClient.BlockByNumber(
							s.Ctx(),
							nil,
						)
						finalBlock, _ := ethClient.BlockByNumber(
							s.Ctx(),
							big.NewInt(int64(rpc.FinalizedBlockNumber)),
						)

						latestBlockNum := latestBlock.Number().Uint64()
						finalizedBlockNum := uint64(0)
						if finalBlock != nil {
							finalizedBlockNum = finalBlock.Number().Uint64()
						}
						s.Require().
							NoError(err, "Error getting current block number during wait")
						s.Logger().Info(
							"block info",
							"latest",
							latestBlockNum,
							"finalized",
							finalizedBlockNum,
							"service",
							service,
						)

						// If the finalized block number is greater than or
						// equal to the target block, we can stop waiting and
						// the test is successful for this node.
						if finalizedBlockNum >= targetBlock {
							s.Logger().Info(
								"Target block reached 🎉",
								"service", service, "block", targetBlock,
							)
							return
						}
					}
				}
			}
		})
	}
	done := make(chan bool, 1)
	go func() {
		// We wait for all the goroutines to finish before we can complete the
		// test.
		wg.WaitAndRecover()
		done <- true
	}()
	select {
	case <-done:
		// Completed without timeout
	case <-time.After(90 * time.Second):
		s.T().Fatal("Timeout waiting for goroutines to finish")
	}
}
