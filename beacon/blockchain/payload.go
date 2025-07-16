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

// Once you provide the right state, we really need to carry out the very same operations
// to extract the data necessary to build the next block, whether current block is
// being rejected or accepted. This is way there can be (and so should be)
// a single function doing these ops. preFetchBuildData is that function.
func (s *Service) preFetchBuildData(st *statedb.StateDB, currentTime math.U64) (
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

	// Carry out on the support state st all the operations needed to
	// process a new payload, namely ProcessSlots and ProcessFork
	if _, err = s.stateProcessor.ProcessSlots(st, blkSlot); err != nil {
		return nil, fmt.Errorf("failed processing block slot: %w", err)
	}
	if err = s.stateProcessor.ProcessFork(st, nextPayloadTimestamp, false); err != nil {
		return nil, fmt.Errorf("failed processing fork: %w", err)
	}

	// Once the state is ready, extract relevant data to build next payload
	payloadWithdrawals, _, err := st.ExpectedWithdrawals(nextPayloadTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed computing expected withdrawals: %w", err)
	}
	epoch := s.chainSpec.SlotToEpoch(blkSlot)
	prevRandao, err := st.GetRandaoMixAtIndex(
		epoch.Unwrap() % s.chainSpec.EpochsPerHistoricalVector(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed retrieving randao: %w", err)
	}

	latestHeader, err := st.GetLatestBlockHeader()
	if err != nil {
		return nil, err
	}

	return &builder.RequestPayloadData{
		Slot:               blkSlot,
		Timestamp:          nextPayloadTimestamp,
		PayloadWithdrawals: payloadWithdrawals,
		PrevRandao:         prevRandao,
		ParentBlockRoot:    latestHeader.HashTreeRoot(),

		// We set the head of our chain to the latest verified block (whether it is final or not)
		HeadEth1BlockHash: lph.GetBlockHash(),

		// Assumuming consensus guarantees single slot finality, the parent
		// of the latest block we verified must be final already.
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
