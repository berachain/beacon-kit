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

package validator

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// forceStartupHead sends a force head FCU to the execution client.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
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
	if err = s.localPayloadBuilder.SendForceHeadFCU(ctx, st, slot+1); err != nil {
		s.logger.Error(
			"failed to send force head FCU",
			"error", err,
		)
	}
}

// retrieveExecutionPayload retrieves the execution payload for the block.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) retrieveExecutionPayload(
	ctx context.Context, st BeaconStateT, blk BeaconBlockT,
) (engineprimitives.BuiltExecutionPayloadEnv[*types.ExecutionPayload], error) {
	// Get the payload for the block.
	envelope, err := s.localPayloadBuilder.
		RetrievePayload(
			ctx,
			blk.GetSlot(),
			blk.GetParentBlockRoot(),
		)
	if err != nil {
		s.metrics.failedToRetrievePayload(
			blk.GetSlot(),
			err,
		)

		// The latest execution payload header will be from the previous block
		// during the block building phase.
		var lph *types.ExecutionPayloadHeader
		lph, err = st.GetLatestExecutionPayloadHeader()
		if err != nil {
			return nil, err
		}

		// If we failed to retrieve the payload, request a synchrnous payload.
		//
		// NOTE: The state here is properly configured by the
		// prepareStateForBuilding
		//
		// call that needs to be called before requesting the Payload.
		// TODO: We should decouple the PayloadBuilder from BeaconState to make
		// this less confusing.
		return s.localPayloadBuilder.RequestPayloadSync(
			ctx,
			st,
			blk.GetSlot(),
			// TODO: this is hood.
			max(
				//#nosec:G701
				uint64(time.Now().Unix()+1),
				uint64((lph.GetTimestamp()+1)),
			),
			blk.GetParentBlockRoot(),
			lph.GetBlockHash(),
			lph.GetParentHash(),
		)
	}
	return envelope, nil
}

// handleRebuildPayloadForRejectedBlock handles the case where the incoming
// block was rejected and we need to rebuild the payload for the current slot.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
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
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) rebuildPayloadForRejectedBlock(
	ctx context.Context,
	st BeaconStateT,
) error {
	var (
		previousBlockRoot primitives.Root
		lph               engineprimitives.ExecutionPayloadHeader
		slot              math.Slot
	)

	s.logger.Info("rebuilding payload for rejected block ⏳ ")

	// In order to rebuild a payload for the current slot, we need to know the
	// previous block root, since we know that this is unmodified state.
	// We can safely get the latest block header and then rebuild the
	// previous block and it's root.
	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return err
	}

	stateRoot, err := st.HashTreeRoot()
	if err != nil {
		return err
	}

	latestHeader.StateRoot = stateRoot
	previousBlockRoot, err = latestHeader.HashTreeRoot()
	if err != nil {
		return err
	}

	// We need to get the *last* finalized execution payload, thus
	// the BeaconState that was passed in must be `unmodified`.
	lph, err = st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	slot, err = st.GetSlot()
	if err != nil {
		return err
	}

	// Submit a request for a new payload.
	if _, err = s.localPayloadBuilder.RequestPayloadAsync(
		ctx,
		st,
		// We are rebuilding for the current slot.
		slot,
		// TODO: this is hood as fuck.
		max(
			//#nosec:G701
			uint64(time.Now().Unix()+1),
			uint64((lph.GetTimestamp()+1)),
		),
		// We set the parent root to the previous block root.
		previousBlockRoot,
		// We set the head of our chain to previous finalized block.
		lph.GetBlockHash(),
		// We can say that the payload from the previous block is *finalized*,
		// TODO: This is making an assumption about the consensus rules
		// and possibly should be made more explicit later on.
		lph.GetBlockHash(),
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
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) handleOptimisticPayloadBuild(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) {
	if err := s.optimisticPayloadBuild(ctx, st, blk); err != nil {
		s.logger.Error(
			"failed to build optimistic payload",
			"for_slot", blk.GetSlot()+1,
			"error", err,
		)
	}
}

// optimisticPayloadBuild builds a payload for the next slot.
func (s *Service[
	BeaconBlockT, BeaconBlockBodyT, BeaconStateT,
	BlobSidecarsT, DepositStoreT, ForkDataT,
]) optimisticPayloadBuild(
	ctx context.Context,
	st BeaconStateT,
	blk BeaconBlockT,
) error {
	// We are building for the next slot, so we increment the slot relative
	// to the block we just processed.
	slot := blk.GetSlot() + 1

	s.logger.Info(
		"optimistically triggering payload build for next slot 🛩️ ",
		"next_slot", slot,
	)

	// We know that this block was properly formed so we can
	// calculate the block hash easily.
	blkRoot, err := blk.HashTreeRoot()
	if err != nil {
		return err
	}

	// We process the slot to update any RANDAO values.
	if _, err = s.stateProcessor.ProcessSlot(
		st,
	); err != nil {
		return err
	}

	// We then trigger a request for the next payload.
	payload := blk.GetBody().GetExecutionPayload()
	if _, err = s.localPayloadBuilder.RequestPayloadAsync(
		ctx, st,
		slot,
		// TODO: this is hood as fuck.
		max(
			//#nosec:G701
			uint64(time.Now().Unix()+int64(s.chainSpec.TargetSecondsPerEth1Block())),
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
