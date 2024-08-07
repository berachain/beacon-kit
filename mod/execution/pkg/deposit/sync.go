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

package deposit

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
)

// defaultRetryInterval processes a deposit event.
const defaultRetryInterval = 20 * time.Second

// depositFetcher processes a deposit event.
func (s *Service[
	_, _, _, _, _, _,
]) depositFetcher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.feed:
			if msg.Is(messages.BeaconBlockFinalized) {
				blockNum := msg.Data().
					GetBody().GetExecutionPayload().GetNumber()
				s.fetchAndStoreDeposits(ctx, blockNum-s.eth1FollowDistance)
			}
		}
	}
}

// depositCatchupFetcher fetches deposits for blocks that failed to be
// processed.
func (s *Service[
	_, _, _, _, _, _,
]) depositCatchupFetcher(ctx context.Context) {
	ticker := time.NewTicker(defaultRetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if len(s.failedBlocks) == 0 {
				continue
			}
			s.logger.Warn(
				"Failed to get deposits from block(s), retrying...",
				"num_blocks",
				s.failedBlocks,
			)

			// Fetch deposits for blocks that failed to be processed.
			for blockNum := range s.failedBlocks {
				s.fetchAndStoreDeposits(ctx, blockNum)
			}
		}
	}
}

func (s *Service[
	_, _, _, _, _, _,
]) fetchAndStoreDeposits(ctx context.Context, blockNum math.U64) {
	deposits, err := s.dc.ReadDeposits(ctx, blockNum)
	if err != nil {
		s.metrics.markFailedToGetBlockLogs(blockNum)
		s.failedBlocks[blockNum] = struct{}{}
		return
	}

	if len(deposits) > 0 {
		s.logger.Info(
			"Found deposits on execution layer",
			"block", blockNum, "deposits", len(deposits),
		)
	}

	if err = s.ds.EnqueueDeposits(deposits); err != nil {
		s.logger.Error("Failed to store deposits", "error", err)
		s.failedBlocks[blockNum] = struct{}{}
		return
	}

	delete(s.failedBlocks, blockNum)
}
