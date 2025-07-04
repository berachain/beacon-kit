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

package builder

import (
	"context"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/payload/cache"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

type PayloadCache interface {
	GetAndEvict(slot math.Slot, stateRoot common.Root) (cache.PayloadIDCacheResult, bool)
	Set(slot math.Slot, stateRoot common.Root, pid engineprimitives.PayloadID, version common.Version)
}

// AttributesFactory is the interface for the attributes factory.
type AttributesFactory interface {
	BuildPayloadAttributes(
		timestamp math.U64,
		payloadWithdrawals engineprimitives.Withdrawals,
		prevRandao common.Bytes32,
		prevHeadRoot common.Root,
	) (*engineprimitives.PayloadAttributes, error)
}

// ExecutionEngine is the interface for the execution engine.
type ExecutionEngine interface {
	// GetPayload returns the payload and blobs bundle for the given slot.
	GetPayload(
		ctx context.Context,
		req *ctypes.GetPayloadRequest,
	) (ctypes.BuiltExecutionPayloadEnv, error)
	// NotifyForkchoiceUpdate notifies the execution client of a forkchoice
	// update.
	NotifyForkchoiceUpdate(
		ctx context.Context,
		req *ctypes.ForkchoiceUpdateRequest,
	) (*engineprimitives.PayloadID, error)
}

type ChainSpec interface {
	ActiveForkVersionForTimestamp(timestamp math.U64) common.Version
	SlotToEpoch(slot math.Slot) math.Epoch
	EpochsPerHistoricalVector() uint64
}
