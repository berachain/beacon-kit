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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// forceStartupHead sends a force head FCU to the execution client.
func (s *Service[
	_, _, _, _, _, _,
]) forceStartupHead(
	ctx context.Context,
	st *statedb.StateDB,
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
	if err = s.localBuilder.SendForceHeadFCU(ctx, st, slot+1); err != nil {
		s.logger.Error(
			"failed to send force head FCU",
			"error", err,
		)
	}
}

// handleRebuildPayloadForRejectedBlock handles the case where the incoming
// block was rejected and we need to rebuild the payload for the current slot.
func (s *Service[
	_, _, _, _, _, _,
]) handleRebuildPayloadForRejectedBlock(
	ctx context.Context,
	st *statedb.StateDB,
	nextPayloadTimestamp math.U64,
) {
	if err := s.rebuildPayloadForRejectedBlock(
		ctx,
		st,
		nextPayloadTimestamp,
	); err != nil {
		s.logger.Error(
			"failed to rebuild payload for nil block",
			"error", err,
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
	_, _, _, _,
	_, _,
]) rebuildPayloadForRejectedBlock(
	ctx context.Context,
	st *statedb.StateDB,
	nextPayloadTimestamp math.U64,
) error {
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

	// Set the previous state root on the header.
	latestHeader.SetStateRoot(st.HashTreeRoot())

	// We need to get the *last* finalized execution payload, thus
	// the BeaconState that was passed in must be `unmodified`.
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	// Submit a request for a new payload.
	if _, err = s.localBuilder.RequestPayloadAsync(
		ctx,
		st,
		// We are rebuilding for the current slot.
		stateSlot,
		nextPayloadTimestamp.Unwrap(),
		// We set the parent root to the previous block root.
		latestHeader.HashTreeRoot(),
		// We set the head of our chain to the previous finalized block.
		lph.GetBlockHash(),
		// We can say that the payload from the previous block is *finalized*,
		// TODO: This is making an assumption about the consensus rules
		// and possibly should be made more explicit later on.
		lph.GetParentHash(),
	); err != nil {
		s.metrics.markRebuildPayloadForRejectedBlockFailure(stateSlot, err)
		return err
	}
	s.metrics.markRebuildPayloadForRejectedBlockSuccess(stateSlot)
	return nil
}

// handleOptimisticPayloadBuild handles optimistically
// building for the next slot.
func (s *Service[
	_, _, _, _, _, _,
]) handleOptimisticPayloadBuild(
	ctx context.Context,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	nextPayloadTimestamp math.U64,
) {
	if err := s.optimisticPayloadBuild(
		ctx,
		st,
		blk,
		nextPayloadTimestamp,
	); err != nil {
		s.logger.Error(
			"Failed to build optimistic payload",
			"for_slot", (blk.GetSlot() + 1).Base10(),
			"error", err,
		)
	}
}

// optimisticPayloadBuild builds a payload for the next slot.
func (s *Service[
	_, _, _, _, _, _,
]) optimisticPayloadBuild(
	ctx context.Context,
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	nextPayloadTimestamp math.U64,
) error {
	// We are building for the next slot, so we increment the slot relative
	// to the block we just processed.
	slot := blk.GetSlot() + 1

	s.logger.Info(
		"Optimistically triggering payload build for next slot 🛩️ ",
		"next_slot", slot.Base10(),
	)

	// We process the slot to update any RANDAO values.
	if _, err := s.stateProcessor.ProcessSlots(st, slot); err != nil {
		return err
	}

	// We then trigger a request for the next payload.
	payload := blk.GetBody().GetExecutionPayload()
	if _, err := s.localBuilder.RequestPayloadAsync(
		ctx, st,
		slot,
		nextPayloadTimestamp.Unwrap(),
		// The previous block root is simply the root of the block we just
		// processed.
		blk.HashTreeRoot(),
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
