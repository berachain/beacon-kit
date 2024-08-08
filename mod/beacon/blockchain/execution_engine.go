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

package blockchain

import (
	"context"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
)

// sendPostBlockFCU sends a forkchoice update to the execution client.
func (s *Service[
	_, BeaconBlockT, _, _, BeaconStateT, _, _, _, _, _, _,
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

	if !s.shouldBuildOptimisticPayloads() && s.lb.Enabled() {
		s.sendNextFCUWithAttributes(ctx, st, blk, lph)
	} else {
		s.sendNextFCUWithoutAttributes(ctx, blk, lph)
	}
}

// sendNextFCUWithAttributes sends a forkchoice update to the execution
// client with attributes.
func (s *Service[
	_, BeaconBlockT, _, _, BeaconStateT,
	_, _, ExecutionPayloadHeaderT, _, _, _,
]) sendNextFCUWithAttributes(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
	lph ExecutionPayloadHeaderT,
) {
	stCopy := st.Copy()
	if _, err := s.sp.ProcessSlots(stCopy, blk.GetSlot()+1); err != nil {
		s.logger.Error(
			"failed to process slots in non-optimistic payload",
			"error", err,
		)
		return
	}

	prevBlockRoot := blk.HashTreeRoot()
	if _, err := s.lb.RequestPayloadAsync(
		ctx,
		stCopy,
		blk.GetSlot()+1,
		s.calculateNextTimestamp(blk),
		prevBlockRoot,
		lph.GetBlockHash(),
		lph.GetParentHash(),
	); err != nil {
		s.logger.Error(
			"failed to send forkchoice update with attributes in non-optimistic payload",
			"error",
			err,
		)
	}
}

// sendNextFCUWithoutAttributes sends a forkchoice update to the
// execution client without attributes.
func (s *Service[
	_, BeaconBlockT, _, _, _, _, _,
	ExecutionPayloadHeaderT, _, PayloadAttributesT, _,
]) sendNextFCUWithoutAttributes(
	ctx context.Context,
	blk BeaconBlockT,
	lph ExecutionPayloadHeaderT,
) {
	if _, _, err := s.ee.NotifyForkchoiceUpdate(
		ctx,
		// TODO: Switch to New().
		engineprimitives.
			BuildForkchoiceUpdateRequestNoAttrs[PayloadAttributesT](
			&engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      lph.GetBlockHash(),
				SafeBlockHash:      lph.GetParentHash(),
				FinalizedBlockHash: lph.GetParentHash(),
			},
			s.cs.ActiveForkVersionForSlot(blk.GetSlot()),
		),
	); err != nil {
		s.logger.Error(
			"failed to send forkchoice update without attributes",
			"error", err,
		)
	}
}

// calculateNextTimestamp calculates the next timestamp for an execution
// payload.
//
// TODO: This is hood and needs to be improved.
func (s *Service[
	_, BeaconBlockT, _, _, _, _, _, _, _, _, _,
]) calculateNextTimestamp(blk BeaconBlockT) uint64 {
	//#nosec:G701 // not an issue in practice.
	return max(
		uint64(time.Now().Unix()+int64(s.cs.TargetSecondsPerEth1Block())),
		uint64(blk.GetBody().GetExecutionPayload().GetTimestamp()+1),
	)
}
