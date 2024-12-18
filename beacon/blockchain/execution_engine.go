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
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// sendPostBlockFCU sends a forkchoice update to the execution client after a
// block is finalized.This function should only be used to notify
// the EL client of the new head and should not request optimistic builds, as:
// Optimistic clients already request builds in handleOptimisticPayloadBuild()
// Non-optimistic clients should never request optimistic builds.
func (s *Service[
	_, _, ConsensusBlockT, _, _, _,
]) sendPostBlockFCU(
	ctx context.Context,
	st *statedb.StateDB,
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

	// Send a forkchoice update without payload attributes to notify
	// EL of the new head.
	beaconBlk := blk.GetBeaconBlock()
	if _, _, err = s.executionEngine.NotifyForkchoiceUpdate(
		ctx,
		// TODO: Switch to New().
		ctypes.
			BuildForkchoiceUpdateRequestNoAttrs(
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
