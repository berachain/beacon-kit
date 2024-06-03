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
)

// depositFetcher processes a deposit event.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT, DepositT,
]) depositFetcher(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case blk := <-s.newBlock:
			querierBlockNum := blk.GetBody().GetExecutionPayload().
				GetNumber() - s.eth1FollowDistance

			deposits, err := s.dc.ReadDeposits(ctx, querierBlockNum)
			if err != nil {
				s.logger.Error("Failed to read deposits", "error", err)
				continue
			}

			if len(deposits) == 0 {
				s.logger.Info(
					"waiting for deposits from execution layer",
					"block",
					querierBlockNum,
				)
			} else {
				s.logger.Info(
					"found deposits on execution layer",
					"block", querierBlockNum, "deposits", len(deposits),
				)
			}

			if err = s.ds.EnqueueDeposits(deposits); err != nil {
				s.logger.Error("Failed to enqueue deposits", "error", err)
				continue
			}
		}
	}
}
