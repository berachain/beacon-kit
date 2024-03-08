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
	"math/big"
	"time"

	"github.com/berachain/beacon-kit/e2e/suite"
	"github.com/ethereum/go-ethereum/rpc"
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
	wg := conc.NewWaitGroup()
	for name, executionClient := range s.ExecutionClients() {
		wg.Go(func() {
			s.Logger().
				Info("Waiting for connection...", "name", name)
			for {
				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-s.Ctx().Done():
						return
					case <-ticker.C:
						latestBlock, err := executionClient.BlockByNumber(
							s.Ctx(),
							nil,
						)
						finalBlock, _ := executionClient.BlockByNumber(
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
							"chain info",
							"latest_block_num",
							latestBlockNum,
							"finalized_block_num",
							finalizedBlockNum,
							"name",
							name,
						)

						// If the finalized block number is greater than or
						// equal to the target block, we can stop waiting and
						// the test is successful for this node.
						if finalizedBlockNum >= targetBlock {
							s.Logger().Info(
								"Target block reached 🎉",
								"block", targetBlock,
								"name", name,
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
