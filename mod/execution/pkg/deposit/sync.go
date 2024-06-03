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

package deposit

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// defaultRetryInterval processes a deposit event.
const defaultRetryInterval = 20 * time.Second

// depositFetcher processes a deposit event.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) depositFetcher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case blk := <-s.newBlock:
			querierBlockNum := blk.
				GetBody().GetExecutionPayload().GetNumber() - s.eth1FollowDistance
			s.fetchAndStoreDeposits(ctx, querierBlockNum)
		}
	}
}

// depositCatchupFetcher fetches deposits for blocks that failed to be
// processed.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) depositCatchupFetcher(ctx context.Context) {
	ticker := time.NewTicker(defaultRetryInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Fetch deposits for blocks that failed to be processed.
			for blockNum := range s.failedBlocks {
				s.fetchAndStoreDeposits(ctx, blockNum)
			}
		}
	}
}
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) fetchAndStoreDeposits(ctx context.Context, blockNum math.U64) {
	deposits, err := s.dc.ReadDeposits(ctx, blockNum)
	if err != nil {
		s.logger.Error("Failed to read deposits", "error", err)
		s.failedBlocks[blockNum] = struct{}{}
		return
	}

	if len(deposits) > 0 {
		s.logger.Info(
			"found deposits on execution layer",
			"block", blockNum, "deposits", len(deposits),
		)
	}

	if err = s.ds.EnqueueDeposits(deposits); err != nil {
		s.logger.Error("Failed to store deposits", "error", err)
		s.failedBlocks[blockNum] = struct{}{}
		return
	}

	if s.failedBlocks[blockNum] != struct{}{} {
		delete(s.failedBlocks, blockNum)
	}
}
