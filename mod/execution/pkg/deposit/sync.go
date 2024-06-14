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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// depositFetcher processes a deposit event.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT,
	WithdrawalCredentialsT, DepositT,
]) depositFetcher(ctx context.Context) {
	ch := make(chan BlockEventT)
	sub := s.feed.Subscribe(ch)
	defer sub.Unsubscribe()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-ch:
			if event.Is(events.BeaconBlockFinalized) {
				blockNum := event.Data().
					GetBody().GetExecutionPayload().GetNumber()
				s.fetchAndStoreDeposits(ctx, blockNum-s.eth1FollowDistance)
			}
		}
	}
}

// depositCatchupFetcher fetches deposits for blocks that failed to be
// processed.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT,
	WithdrawalCredentialsT, DepositT,
]) depositCatchupFetcher(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case blockNum := <-s.failedBlocks:
			s.logger.Warn(
				"failed to get deposits from block(s), retrying...",
				"block_num",
				blockNum,
			)
			s.fetchAndStoreDeposits(ctx, blockNum)
		}
	}
}

func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BlockEventT,
	ExecutionPayloadT, SubscriptionT,
	WithdrawalCredentialsT, DepositT,
]) fetchAndStoreDeposits(ctx context.Context, blockNum math.U64) {
	deposits, err := s.dc.ReadDeposits(ctx, blockNum)
	if err != nil {
		s.metrics.markFailedToGetBlockLogs(blockNum)
		s.failedBlocks <- blockNum
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
		s.failedBlocks <- blockNum
		return
	}
}
