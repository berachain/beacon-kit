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
	"time"

	"github.com/berachain/beacon-kit/e2e/suite"
	"golang.org/x/sync/errgroup"
)

// BeaconE2ESuite is a suite of tests simulating a fully function beacon-kit
// network.
type BeaconKitE2ESuite struct {
	suite.KurtosisE2ESuite
}

// TestBasicStartup tests the basic startup of the beacon-kit network.
func (s *BeaconKitE2ESuite) TestBasicStartup() {
	targetBlock := uint64(5)
	eg, groupCtx := errgroup.WithContext(context.Background())
	groupCctx, cancel := context.WithTimeout(groupCtx, 90*time.Second)
	defer cancel()

	for _, executionClient := range s.ExecutionClients() {
		eg.Go(
			func() error {
				return executionClient.WaitForLatestBlockNumber(
					groupCctx,
					targetBlock,
				)
			},
		)
	}

	if err := eg.Wait(); err != nil {
		s.T().Fatal(err)
	}
}
