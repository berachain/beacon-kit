// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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

package blockchain

import (
	"context"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
)

// sendPostBlockFCU sends a forkchoice update to the execution client.
func (s *Service[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositT,
	DepositStoreT,
]) sendPostBlockFCU(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		s.logger.Error(
			"failed to get latest execution payload in postBlockProcess",
			"error", err,
		)
		return
	}

	// This is technically not an optimistic payload
	// TODO: This needs a refactor, big hood energy.
	//nolint:nestif // todo fix.5
	if !s.shouldBuildOptimisticPayloads() && s.lb.Enabled() {
		stCopy := st.Copy()
		if _, err = s.sp.ProcessSlots(
			stCopy, blk.GetSlot()+1,
		); err != nil {
			return
		}

		var prevBlockRoot primitives.Root
		prevBlockRoot, err = blk.HashTreeRoot()
		if err != nil {
			s.logger.Error(
				"failed to get block root in postBlockProcess",
				"error", err,
			)
			return
		}

		// Ask the builder to send a forkchoice update with attributes.
		// This will trigger a new payload to be built.
		if _, err = s.lb.RequestPayloadAsync(
			ctx,
			stCopy,
			blk.GetSlot()+1,
			//#nosec:G701 // won't realistically overflow.
			// TODO: clock time properly.
			(max(
				uint64(time.Now().Unix()+int64(s.cs.TargetSecondsPerEth1Block())),
				uint64(blk.GetBody().GetExecutionPayload().GetTimestamp()+1),
			)),
			prevBlockRoot,
			lph.GetBlockHash(),
			lph.GetParentHash(),
		); err == nil {
			return
		}

		// If we error we log and continue, we try again without building a
		// block
		// just incase this can help get our execution client back on track.
		s.logger.
			Error(
				"failed to send forkchoice update with attributes",
				"error", err,
			)
	} else {
		// If we are not building blocks, or we failed to build a block
		// we can just send the forkchoice update without attributes.
		_, _, err = s.ee.NotifyForkchoiceUpdate(
			ctx,
			engineprimitives.BuildForkchoiceUpdateRequest(
				&engineprimitives.ForkchoiceStateV1{
					HeadBlockHash:      lph.GetBlockHash(),
					SafeBlockHash:      lph.GetParentHash(),
					FinalizedBlockHash: lph.GetParentHash(),
				},
				nil,
				s.cs.ActiveForkVersionForSlot(blk.GetSlot()),
			),
		)
		if err != nil {
			s.logger.Error(
				"failed to send forkchoice update without attributes",
				"error", err,
			)
		}
	}
}
