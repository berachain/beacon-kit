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
	"strings"
	"testing"
	"time"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/itsdevbear/bolaris/e2e/suite"
	"github.com/kurtosis-tech/kurtosis/api/golang/core/lib/services"
	"github.com/sourcegraph/conc"
)

// TestBeaconKitE2ESuite runs the test suite.
func TestBeaconKitE2ESuite(t *testing.T) {
	suite.Run(t, new(BeaconKitE2ESuite))
}

// BeaconE2ESuite is a suite of tests simluating a fully function beacon-kit
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

	var jsonRPCPorts = make([]string, 0)
	for k, v := range svrcs {
		s.Logger().Info("Service started", "service", k, "uuid", v)
		var serviceCtx *services.ServiceContext
		serviceCtx, err = s.Enclave().GetServiceContext(string(v))
		s.Require().NoError(err, "Error getting service context")

		// Get the public ports representing eth JSON-RPC endpoints.
		jsonRPC, ok := serviceCtx.GetPublicPorts()["rpc"]
		if ok {
			jsonRPCPorts = append(
				jsonRPCPorts,
				strings.Split(jsonRPC.String(), "/")[0],
			)
		}
	}

	wg := conc.NewWaitGroup()
	for port := range jsonRPCPorts {
		wg.Go(func() {
			var ethClient *ethclient.Client
			ethClient, err = ethclient.Dial(
				"http://0.0.0.0:" + jsonRPCPorts[port],
			)
			s.Require().NoError(err, "Error creating eth client")
			for {
				ticker := time.NewTicker(3 * time.Second)
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
						s.Require().
							NoError(err, "Error getting current block number during wait")
						s.Logger().Info(
							"Finalized block",
							"block",
							latestBlock.Number().Uint64(),
							"port",
							jsonRPCPorts[port],
						)
						if latestBlock.Number().Uint64() >= targetBlock {
							return
						}
					}
				}
			}
		})
	}

	done := make(chan bool, 1)
	go func() {
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
