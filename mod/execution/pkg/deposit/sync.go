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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// defaultRetryInterval processes a deposit event.
const defaultRetryInterval = 20 * time.Second

// depositFetcher returns a function that retrieves the block number from the
// event and fetches and stores the deposits for that block.
func (s *Service[
	BeaconBlockT, _, _, _, _,
]) depositFetcher(ctx context.Context, event async.Event[BeaconBlockT]) {
	blockNum := event.Data().GetBody().GetExecutionPayload().GetNumber()
	s.fetchAndStoreDeposits(ctx, blockNum-s.eth1FollowDistance)
}

// depositCatchupFetcher fetches deposits for blocks that failed to be
// processed.
func (s *Service[
	_, _, _, _, _,
]) depositCatchupFetcher(ctx context.Context) {
	ticker := time.NewTicker(defaultRetryInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Fetch deposits for blocks that failed to be processed.
			s.failedBlocks.Range(func(key, _ any) bool {
				blockNum, ok := key.(math.U64)
				if !ok {
					s.logger.Error("Failed to parse key to math.U64", "key", key)
				} else {
					s.logger.Info("Retry fetch and store deposits", "block_num", blockNum)
					s.fetchAndStoreDeposits(ctx, blockNum)
				}
				return true
			})
		}
	}
}

func (s *Service[
	_, _, _, _, _,
]) fetchAndStoreDeposits(ctx context.Context, blockNum math.U64) {
	deposits, err := s.dc.ReadDeposits(ctx, blockNum)
	if err != nil {
		s.metrics.markFailedToGetBlockLogs(blockNum)
		s.failedBlocks.Store(blockNum, true)
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
		s.failedBlocks.Store(blockNum, true)
		return
	}

	s.failedBlocks.Delete(blockNum)
}
