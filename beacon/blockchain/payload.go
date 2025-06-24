// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"fmt"

	payloadtime "github.com/berachain/beacon-kit/beacon/payload-time"
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/payload/builder"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// forceSyncUponProcess sends a force head FCU to the execution client.
func (s *Service) forceSyncUponProcess(
	ctx context.Context,
	st *statedb.StateDB,
) {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		s.logger.Error(
			"failed to get latest execution payload header",
			"error", err,
		)
		return
	}

	s.logger.Info(
		"Sending startup forkchoice update to execution client",
		"head_eth1_hash", lph.GetBlockHash(),
		"safe_eth1_hash", lph.GetParentHash(),
		"finalized_eth1_hash", lph.GetParentHash(),
		"for_slot", lph.GetNumber(),
	)

	// Submit the forkchoice update to the execution client.
	req := ctypes.BuildForkchoiceUpdateRequestNoAttrs(
		&engineprimitives.ForkchoiceStateV1{
			HeadBlockHash:      lph.GetBlockHash(),
			SafeBlockHash:      lph.GetParentHash(),
			FinalizedBlockHash: lph.GetParentHash(),
		},
		s.chainSpec.ActiveForkVersionForTimestamp(lph.GetTimestamp()),
	)
	if _, err = s.executionEngine.NotifyForkchoiceUpdate(ctx, req); err != nil {
		s.logger.Error(
			"failed to send force head FCU",
			"error", err,
		)
	}
}

// forceSyncUponFinalize sends a new payload and force startup FCU to the Execution
// Layer client. This informs the EL client of the new head and forces a SYNC
// if blocks are missing. This function should only be run once at startup.
func (s *Service) forceSyncUponFinalize(
	ctx context.Context,
	beaconBlock *ctypes.BeaconBlock,
) error {
	// NewPayload call first to load payload into EL client.
	executionPayload := beaconBlock.GetBody().GetExecutionPayload()
	payloadReq, err := ctypes.BuildNewPayloadRequestFromFork(beaconBlock)
	if err != nil {
		return err
	}

	if err = payloadReq.HasValidVersionedAndBlockHashes(); err != nil {
		return err
	}

	// We set retryOnSyncingStatus to false here. We can ignore SYNCING status and proceed
	// to the FCU.
	err = s.executionEngine.NotifyNewPayload(ctx, payloadReq, false)
	if err != nil {
		return fmt.Errorf("startSyncUponFinalize NotifyNewPayload failed: %w", err)
	}

	// Submit the forkchoice update to the EL client. This will ensure that it is either synced or
	// starts up a sync.
	req := ctypes.BuildForkchoiceUpdateRequestNoAttrs(
		&engineprimitives.ForkchoiceStateV1{
			HeadBlockHash:      executionPayload.GetBlockHash(),
			SafeBlockHash:      executionPayload.GetParentHash(),
			FinalizedBlockHash: executionPayload.GetParentHash(),
		},
		s.chainSpec.ActiveForkVersionForTimestamp(executionPayload.GetTimestamp()),
	)

	switch _, err = s.executionEngine.NotifyForkchoiceUpdate(ctx, req); {
	case err == nil:
		return nil

	case errors.IsAny(err,
		engineerrors.ErrSyncingPayloadStatus,
		engineerrors.ErrAcceptedPayloadStatus):
		s.logger.Warn(
			//nolint:lll // long message on one line for readability.
			`Your execution client is syncing. It should be downloading eth blocks from its peers. Restart the beacon node once the execution client is caught up.`,
		)
		return err

	default:
		return fmt.Errorf("force startup NotifyForkchoiceUpdate failed: %w", err)
	}
}

func (s *Service) preFetchBuildDataForRejection(
	st *statedb.StateDB,
	currentTime math.U64,
) (
	*builder.RequestPayloadData,
	error,
) {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving latest execution payload header: %w", err)
	}
	nextPayloadTimestamp := payloadtime.Next(
		currentTime,
		lph.GetTimestamp(),
		true, // buildOptimistically
	)

	// In order to rebuild a payload for the current slot, we need to know the
	// previous block root, since we know that this is an unmodified state.
	// We can safely get the latest block header and then rebuild the
	// previous block and it's root.
	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return nil, err
	}

	stateSlot, err := st.GetSlot()
	if err != nil {
		return nil, err
	}

	// Set the previous state root on the header.
	latestHeader.SetStateRoot(st.HashTreeRoot())

	// We must prepare the state for the fork version of the new block being built to handle
	// the case where the new block is on a new fork version. Although we do not have the
	// confirmed timestamp by the EL, we will assume it to be `nextPayloadTimestamp` to decide
	// the new block's fork version.
	err = s.stateProcessor.ProcessFork(st, nextPayloadTimestamp, false)
	if err != nil {
		return nil, err
	}

	// Expected payloadWithdrawals to include in this payload.
	payloadWithdrawals, _, err := st.ExpectedWithdrawals(nextPayloadTimestamp)
	if err != nil {
		s.logger.Error(
			"Could not get expected withdrawals to get payload attribute",
			"error",
			err,
		)
		return nil, err
	}
	// Get the previous randao mix.
	epoch := s.chainSpec.SlotToEpoch(stateSlot)
	prevRandao, err := st.GetRandaoMixAtIndex(
		epoch.Unwrap() % s.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return nil, err
	}

	return &builder.RequestPayloadData{
		Slot:               stateSlot,
		Timestamp:          nextPayloadTimestamp,
		PayloadWithdrawals: payloadWithdrawals,
		PrevRandao:         prevRandao,
		ParentBlockRoot:    latestHeader.HashTreeRoot(),
		HeadEth1BlockHash:  lph.GetBlockHash(),
		FinalEth1BlockHash: lph.GetParentHash(),
	}, nil
}

// handleRebuildPayloadForRejectedBlock handles the case where the incoming
// block was rejected and we need to rebuild the payload for the current slot.
func (s *Service) handleRebuildPayloadForRejectedBlock(
	ctx context.Context,
	buildData *builder.RequestPayloadData,
) {
	s.logger.Info("Rebuilding payload for rejected block ⏳ ")
	nextBlkSlot := buildData.Slot
	if _, _, err := s.localBuilder.RequestPayloadAsync(ctx, buildData); err != nil {
		s.metrics.markRebuildPayloadForRejectedBlockFailure(nextBlkSlot, err)
		s.logger.Error(
			"failed to rebuild payload for nil block",
			"error", err,
		)
		return
	}

	s.metrics.markRebuildPayloadForRejectedBlockSuccess(nextBlkSlot)
}

func (s *Service) preFetchBuildDataForSuccess(
	st *statedb.StateDB,
	blk *ctypes.BeaconBlock,
	currentTime math.U64,
) (
	*builder.RequestPayloadData,
	error,
) {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving latest execution payload header: %w", err)
	}
	nextPayloadTimestamp := payloadtime.Next(
		currentTime,
		lph.GetTimestamp(),
		true, // buildOptimistically
	)

	stateSlot, err := st.GetSlot()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving slot from state: %w", err)
	}
	blkSlot := stateSlot + 1

	// We process the slot to update any RANDAO values.
	if _, err = s.stateProcessor.ProcessSlots(st, blkSlot); err != nil {
		return nil, fmt.Errorf("failed processing block slot: %w", err)
	}

	// We must prepare the state for the fork version of the new block being built to handle
	// the case where the new block is on a new fork version. Although we do not have the
	// confirmed timestamp by the EL, we will assume it to be `nextPayloadTimestamp` to decide
	// the new block's fork version.
	err = s.stateProcessor.ProcessFork(st, nextPayloadTimestamp, false)
	if err != nil {
		return nil, fmt.Errorf("failed processing fork: %w", err)
	}

	// Expected payloadWithdrawals to include in this payload.
	payloadWithdrawals, _, err := st.ExpectedWithdrawals(nextPayloadTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed computing expected withdrawals: %w", err)
	}
	// Get the previous randao mix.
	epoch := s.chainSpec.SlotToEpoch(blkSlot)
	prevRandao, err := st.GetRandaoMixAtIndex(
		epoch.Unwrap() % s.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving randao: %w", err)
	}

	return &builder.RequestPayloadData{
		Slot:               blkSlot,
		Timestamp:          nextPayloadTimestamp,
		PayloadWithdrawals: payloadWithdrawals,
		PrevRandao:         prevRandao,
		// The previous block root is simply the root of
		// the block we just verified.
		ParentBlockRoot: blk.HashTreeRoot(),

		// We set the head of our chain to the block we just verified (the latest)
		HeadEth1BlockHash: lph.GetBlockHash(),

		// Assumuming consensus guarantees single slot finality, the parent
		// of the latest block we verified must be final already.
		FinalEth1BlockHash: lph.GetParentHash(),
	}, nil
}

// handleOptimisticPayloadBuild handles optimistically
// building for the next slot.
func (s *Service) handleOptimisticPayloadBuild(
	ctx context.Context,
	buildData *builder.RequestPayloadData,
) {
	s.logger.Info(
		"Optimistically triggering payload build for next slot 🛩️ ",
		"next_slot", buildData.Slot.Base10(),
	)
	if _, _, err := s.localBuilder.RequestPayloadAsync(ctx, buildData); err != nil {
		s.metrics.markOptimisticPayloadBuildFailure(buildData.Slot, err)
		s.logger.Error(
			"Failed to build optimistic payload",
			"for_slot", buildData.Slot.Base10(),
			"error", err,
		)
		return
	}

	s.metrics.markOptimisticPayloadBuildSuccess(buildData.Slot)
}
