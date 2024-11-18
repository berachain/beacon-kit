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

	"github.com/berachain/beacon-kit/mod/config/pkg/spec"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
)

// sendPostBlockFCU sends a forkchoice update to the execution client.
func (s *Service[
	_, ConsensusBlockT, _, _, _, BeaconStateT, _, _, _, _, _,
]) sendPostBlockFCU(
	ctx context.Context,
	st BeaconStateT,
	blk ConsensusBlockT,
) {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		s.logger.Error(
			"failed to get latest execution payload in postBlockProcess",
			"error", err,
		)
		return
	}

	if !s.shouldBuildOptimisticPayloads() && s.localBuilder.Enabled() {
		s.sendNextFCUWithAttributes(ctx, st, blk, lph)
	} else {
		s.sendNextFCUWithoutAttributes(ctx, blk, lph)
	}
}

// sendNextFCUWithAttributes sends a forkchoice update to the execution
// client with attributes.
func (s *Service[
	_, ConsensusBlockT, _, _, _, BeaconStateT,
	_, _, ExecutionPayloadHeaderT, _, _,
]) sendNextFCUWithAttributes(
	ctx context.Context,
	st BeaconStateT,
	blk ConsensusBlockT,
	lph ExecutionPayloadHeaderT,
) {
	beaconBlk := blk.GetBeaconBlock()

	stCopy := st.Copy()
	if _, err := s.stateProcessor.ProcessSlots(
		stCopy, beaconBlk.GetSlot()+1,
	); err != nil {
		s.logger.Error(
			"failed to process slots in non-optimistic payload",
			"error", err,
		)
		return
	}

	nextPayloadTime := max(
		blk.GetConsensusTime()+1,
		lph.GetTimestamp()+1,
	).Unwrap()

	// We set timestamp check on Bartio for backward compatibility reasons
	// TODO: drop this we drop other Bartio special cases.
	if s.chainSpec.DepositEth1ChainID() == spec.BartioChainID {
		nextPayloadTime = max(
			//#nosec:G701
			uint64(time.Now().Unix()+1),
			uint64((lph.GetTimestamp() + 1)),
		)
	}

	prevBlockRoot := beaconBlk.HashTreeRoot()
	if _, err := s.localBuilder.RequestPayloadAsync(
		ctx,
		stCopy,
		beaconBlk.GetSlot()+1,
		nextPayloadTime,
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
	_, ConsensusBlockT, _, _, _, _, _, _,
	ExecutionPayloadHeaderT, _, PayloadAttributesT,
]) sendNextFCUWithoutAttributes(
	ctx context.Context,
	blk ConsensusBlockT,
	lph ExecutionPayloadHeaderT,
) {
	beaconBlk := blk.GetBeaconBlock()

	if _, _, err := s.executionEngine.NotifyForkchoiceUpdate(
		ctx,
		// TODO: Switch to New().
		engineprimitives.
			BuildForkchoiceUpdateRequestNoAttrs[PayloadAttributesT](
			&engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      lph.GetBlockHash(),
				SafeBlockHash:      lph.GetParentHash(),
				FinalizedBlockHash: lph.GetParentHash(),
			},
			s.chainSpec.ActiveForkVersionForSlot(beaconBlk.GetSlot()),
		),
	); err != nil {
		s.logger.Error(
			"failed to send forkchoice update without attributes",
			"error", err,
		)
	}
}
