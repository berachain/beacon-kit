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
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// RequestPayload request a payload for the given slot and
// blocks until the payload is delivered.
func (pb *PayloadBuilder) RequestPayload(
	ctx context.Context,
	st *statedb.StateDB,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot common.Root,
	parentEth1Hash common.ExecutionHash,
	finalBlockHash common.ExecutionHash,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	if !pb.Enabled() {
		return nil, ErrPayloadBuilderDisabled
	}

	// Build the payload and wait for the execution client to
	// return the payload ID.
	// Assemble the payload attributes.
	attrs, err := pb.attributesFactory.BuildPayloadAttributes(
		st,
		slot,
		timestamp,
		parentBlockRoot,
	)
	if err != nil {
		return nil, err
	}

	// Submit the forkchoice update to the execution client.
	payloadID, _, err := pb.ee.NotifyForkchoiceUpdate(
		ctx, &ctypes.ForkchoiceUpdateRequest{
			State: &engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      parentEth1Hash,
				SafeBlockHash:      finalBlockHash,
				FinalizedBlockHash: finalBlockHash,
			},
			PayloadAttributes: attrs,
			ForkVersion:       pb.chainSpec.ActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	}
	if payloadID == nil {
		return nil, ErrNilPayloadID
	}

	// Wait for the payload to be delivered to the execution client.
	pb.logger.Info(
		"Waiting for local payload to be delivered to execution client",
		"for_slot", slot.Base10(), "timeout", pb.cfg.PayloadTimeout.String(),
	)
	select {
	case <-time.After(pb.cfg.PayloadTimeout):
		// We want to trigger delivery of the payload to the execution client
		// before the timestamp expires.
		break
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Get the payload from the execution client.
	return pb.getPayload(ctx, *payloadID, slot)
}

// SendForceHeadFCU builds a payload for the given slot and
// returns the payload ID.
//
// TODO: This should be moved onto a "sync service"
// of some kind.
func (pb *PayloadBuilder) SendForceHeadFCU(
	ctx context.Context,
	st *statedb.StateDB,
	slot math.Slot,
) error {
	if !pb.Enabled() {
		return ErrPayloadBuilderDisabled
	}

	lph, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	pb.logger.Info(
		"Sending startup forkchoice update to execution client",
		"head_eth1_hash", lph.GetBlockHash(),
		"safe_eth1_hash", lph.GetParentHash(),
		"finalized_eth1_hash", lph.GetParentHash(),
		"for_slot", slot.Base10(),
	)

	// Submit the forkchoice update to the execution client.
	var attrs *engineprimitives.PayloadAttributes
	_, _, err = pb.ee.NotifyForkchoiceUpdate(
		ctx, &ctypes.ForkchoiceUpdateRequest{
			State: &engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      lph.GetBlockHash(),
				SafeBlockHash:      lph.GetParentHash(),
				FinalizedBlockHash: lph.GetParentHash(),
			},
			PayloadAttributes: attrs,
			ForkVersion:       pb.chainSpec.ActiveForkVersionForSlot(slot),
		},
	)
	return err
}

func (pb *PayloadBuilder) getPayload(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	slot math.U64,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	envelope, err := pb.ee.GetPayload(
		ctx,
		&ctypes.GetPayloadRequest{
			PayloadID:   payloadID,
			ForkVersion: pb.chainSpec.ActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	}
	if envelope == nil {
		return nil, ErrNilPayloadEnvelope
	}
	if envelope.GetExecutionPayload().Withdrawals == nil {
		return nil, ErrNilWithdrawals
	}
	return envelope, nil
}
