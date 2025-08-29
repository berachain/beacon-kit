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

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// sendPostBlockFCU sends a forkchoice update to the execution client after a
// block is finalized.This function should only be used to notify
// the EL client of the new head and should not request optimistic builds, as:
// Optimistic clients already request builds in handleOptimisticPayloadBuild()
// Non-optimistic clients should never request optimistic builds.
func (s *Service) sendPostBlockFCU(
	ctx context.Context,
	st *statedb.StateDB,
) error {
	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return fmt.Errorf("failed getting latest payload: %w", err)
	}

	// Send a forkchoice update without payload attributes to notify EL of the new head.
	req := ctypes.BuildForkchoiceUpdateRequestNoAttrs(
		&engineprimitives.ForkchoiceStateV1{
			HeadBlockHash:      lph.GetBlockHash(),
			SafeBlockHash:      lph.GetParentHash(),
			FinalizedBlockHash: lph.GetParentHash(),
		},
		s.chainSpec.ActiveForkVersionForTimestamp(lph.GetTimestamp()),
	)
	if _, err = s.executionEngine.NotifyForkchoiceUpdate(ctx, req); err != nil {
		return fmt.Errorf("failed forkchoice update, head %s: %w",
			lph.GetBlockHash().String(),
			err,
		)
	}
	return nil
}
