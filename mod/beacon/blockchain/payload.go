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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// forceStartupHead sends a force head FCU to the execution client.
func (s *Service[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) forceStartupHead(
	ctx context.Context,
	st BeaconStateT,
) {
	slot, err := st.GetSlot()
	if err != nil {
		s.logger.Error(
			"failed to get slot for force startup head",
			"error", err,
		)
		return
	}

	// TODO: Verify if the slot number is correct here, I believe in current
	// form
	// it should be +1'd. Not a big deal until hardforks are in play though.
	if err = s.lb.SendForceHeadFCU(ctx, st, slot+1); err != nil {
		s.logger.Error(
			"failed to send force head FCU",
			"error", err,
		)
	}
}

// handleRebuildPayloadForRejectedBlock handles the case where the incoming
// block was rejected and we need to rebuild the payload for the current slot.
func (s *Service[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) handleRebuildPayloadForRejectedBlock(
	ctx context.Context,
	st BeaconStateT,
) {
	if pErr := s.rebuildPayloadForRejectedBlock(
		ctx, st,
	); pErr != nil {
		s.logger.Error(
			"failed to rebuild payload for nil block",
			"error", pErr,
		)
	}
}

// rebuildPayloadForRejectedBlock rebuilds a payload for the current
// slot, if the incoming block was rejected.
//
// NOTE: We cannot use any data off the incoming block and must recompute
// any required information from our local state. We do this since we have
// rejected the incoming block and it would be unsafe to use any
// information from it.
func (s *Service[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _,
	_, ExecutionPayloadHeaderT, _, _, _, _, _, _, _,
]) rebuildPayloadForRejectedBlock(
	ctx context.Context,
	st BeaconStateT,
) error {
	var (
		prevStateRoot common.Root
		prevBlockRoot common.Root
		lph           ExecutionPayloadHeaderT
		slot          math.Slot
	)

	s.logger.Info("Rebuilding payload for rejected block ⏳ ")

	// In order to rebuild a payload for the current slot, we need to know the
	// previous block root, since we know that this is an unmodified state.
	// We can safely get the latest block header and then rebuild the
	// previous block and it's root.
	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return err
	}

	stateSlot, err := st.GetSlot()
	if err != nil {
		return err
	}

	prevStateRoot, err = st.HashTreeRoot()
	if err != nil {
		return err
	}

	latestHeader.SetStateRoot(prevStateRoot)
	prevBlockRoot, err = latestHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	// We need to get the *last* finalized execution payload, thus
	// the BeaconState that was passed in must be `unmodified`.
	lph, err = st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	// Submit a request for a new payload.
	if _, err = s.lb.RequestPayloadAsync(
		ctx,
		st,
		// We are rebuilding for the current slot.
		stateSlot,
		// TODO: this is hood as fuck.
		max(
			//#nosec:G701
			uint64(time.Now().Unix()+1),
			uint64((lph.GetTimestamp()+1)),
		),
		// We set the parent root to the previous block root.
		prevBlockRoot,
		// We set the head of our chain to the previous finalized block.
		lph.GetBlockHash(),
		// We can say that the payload from the previous block is *finalized*,
		// TODO: This is making an assumption about the consensus rules
		// and possibly should be made more explicit later on.
		lph.GetParentHash(),
	); err != nil {
		s.metrics.markRebuildPayloadForRejectedBlockFailure(slot, err)
		return err
	}
	s.metrics.markRebuildPayloadForRejectedBlockSuccess(slot)
	return nil
}

// handleOptimisticPayloadBuild handles optimistically
// building for the next slot.
func (s *Service[
	_, BeaconBlockT, _, _, BeaconStateT, _, _, _,
	_, _, _, _, _, _, _, _, _, _, _, _,
]) handleOptimisticPayloadBuild(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) {
	if err := s.optimisticPayloadBuild(ctx, st, blk); err != nil {
		s.logger.Error(
			"Failed to build optimistic payload",
			"for_slot", (blk.GetSlot() + 1).Base10(),
			"error", err,
		)
	}
}

// optimisticPayloadBuild builds a payload for the next slot.
func (s *Service[
	_, BeaconBlockT, _, _, BeaconStateT, _, _, _,
	_, _, _, _, _, _, _, _, _, _, _, _,
]) optimisticPayloadBuild(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	// We are building for the next slot, so we increment the slot relative
	// to the block we just processed.
	slot := blk.GetSlot() + 1

	s.logger.Info(
		"Optimistically triggering payload build for next slot 🛩️ ",
		"next_slot", slot.Base10(),
	)

	// We know that this block was properly formed so we can
	// calculate the block hash easily.
	blkRoot, err := blk.HashTreeRoot()
	if err != nil {
		return err
	}

	// We process the slot to update any RANDAO values.
	if _, err = s.sp.ProcessSlots(
		st, slot,
	); err != nil {
		return err
	}

	// We then trigger a request for the next payload.
	payload := blk.GetBody().GetExecutionPayload()
	if _, err = s.lb.RequestPayloadAsync(
		ctx, st,
		slot,
		// TODO: this is hood as fuck.
		max(
			//#nosec:G701
			uint64(time.Now().Unix()+int64(s.cs.TargetSecondsPerEth1Block())),
			uint64((payload.GetTimestamp()+1)),
		),
		// The previous block root is simply the root of the block we just
		// processed.
		blkRoot,
		// We set the head of our chain to the block we just processed.
		payload.GetBlockHash(),
		// We can say that the payload from the previous block is *finalized*,
		// This is safe to do since this block was accepted and the thus the
		// parent hash was deemed valid by the state transition function we
		// just processed.
		payload.GetParentHash(),
	); err != nil {
		s.metrics.markOptimisticPayloadBuildFailure(slot, err)
		return err
	}
	s.metrics.markOptimisticPayloadBuildSuccess(slot)
	return nil
}
