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

package builder

import (
	"context"
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// RequestPayloadAsync builds a payload for the given slot and
// returns the payload ID.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	PayloadAttributesT, PayloadIDT, WithdrawalT,
]) RequestPayloadAsync(
	ctx context.Context,
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot common.Root,
	headEth1BlockHash common.ExecutionHash,
	finalEth1BlockHash common.ExecutionHash,
) (*PayloadIDT, error) {
	if !pb.Enabled() {
		return nil, ErrPayloadBuilderDisabled
	}

	if payloadID, found := pb.pc.Get(slot, parentBlockRoot); found {
		pb.logger.Warn(
			"aborting payload build; payload already exists in cache",
			"for_slot",
			slot.Base10(),
			"parent_block_root",
			parentBlockRoot,
		)
		return &payloadID, nil
	}

	// Assemble the payload attributes.
	attrs, err := pb.attributesFactory.
		BuildPayloadAttributes(st, slot, timestamp, parentBlockRoot)
	if err != nil {
		return nil, err
	}

	// Submit the forkchoice update to the execution client.
	var payloadID *PayloadIDT
	payloadID, _, err = pb.ee.NotifyForkchoiceUpdate(
		ctx, &engineprimitives.ForkchoiceUpdateRequest[PayloadAttributesT]{
			State: &engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      headEth1BlockHash,
				SafeBlockHash:      finalEth1BlockHash,
				FinalizedBlockHash: finalEth1BlockHash,
			},
			PayloadAttributes: attrs,
			ForkVersion:       pb.chainSpec.GetActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	}

	// Only add to cache if we received back a payload ID.
	if payloadID != nil {
		pb.pc.Set(slot, parentBlockRoot, *payloadID)
	}

	return payloadID, nil
}

// RequestPayloadSync request a payload for the given slot and
// blocks until the payload is delivered.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	PayloadAttributesT, PayloadIDT, WithdrawalT,
]) RequestPayloadSync(
	ctx context.Context,
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot common.Root,
	parentEth1Hash common.ExecutionHash,
	finalBlockHash common.ExecutionHash,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	if !pb.Enabled() {
		return nil, ErrPayloadBuilderDisabled
	}

	// Build the payload and wait for the execution client to
	// return the payload ID.
	payloadID, err := pb.RequestPayloadAsync(
		ctx,
		st,
		slot,
		timestamp,
		parentBlockRoot,
		parentEth1Hash,
		finalBlockHash,
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

// RetrievePayload attempts to pull a previously built payload
// by reading a payloadID from the builder's cache. If it fails to
// retrieve a payload, it will build a new payload and wait for the
// execution client to return the payload.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	PayloadAttributesT, PayloadIDT, WithdrawalT,
]) RetrievePayload(
	ctx context.Context,
	slot math.Slot,
	parentBlockRoot common.Root,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	if !pb.Enabled() {
		return nil, ErrPayloadBuilderDisabled
	}

	// Attempt to see if we previously fired off a payload built for
	// this particular slot and parent block root.
	payloadID, found := pb.pc.Get(slot, parentBlockRoot)
	if !found {
		return nil, ErrPayloadIDNotFound
	}

	// Get the payload from the execution client.
	envelope, err := pb.getPayload(ctx, payloadID, slot)
	if err != nil {
		return nil, err
	}

	overrideBuilder := envelope.ShouldOverrideBuilder()
	args := []any{
		"for_slot", slot.Base10(),
		"override_builder", overrideBuilder,
	}

	payload := envelope.GetExecutionPayload()
	if !payload.IsNil() {
		args = append(args,
			"payload_block_hash", payload.GetBlockHash(),
			"parent_hash", payload.GetParentHash(),
		)
	}

	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle != nil {
		args = append(args, "num_blobs", len(blobsBundle.GetBlobs()))
	}

	pb.logger.Info("Payload retrieved from local builder", args...)

	// If the payload was built by a different builder, something is
	// wrong the EL<>CL setup.
	if payload.GetFeeRecipient() != pb.cfg.SuggestedFeeRecipient {
		pb.logger.Warn(
			"Payload fee recipient does not match suggested fee recipient - "+
				"please check both your CL and EL configuration",
			"payload_fee_recipient", payload.GetFeeRecipient(),
			"suggested_fee_recipient", pb.cfg.SuggestedFeeRecipient,
		)
	}
	return envelope, err
}

// SendForceHeadFCU builds a payload for the given slot and
// returns the payload ID.
//
// TODO: This should be moved onto a "sync service"
// of some kind.
func (pb *PayloadBuilder[
	BeaconStateT, ExecutionPayloadT, ExecutionPayloadHeaderT,
	PayloadAttributesT, PayloadIDT, WithdrawalT,
]) SendForceHeadFCU(
	ctx context.Context,
	st BeaconStateT,
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
	var attrs PayloadAttributesT
	_, _, err = pb.ee.NotifyForkchoiceUpdate(
		ctx, &engineprimitives.ForkchoiceUpdateRequest[PayloadAttributesT]{
			State: &engineprimitives.ForkchoiceStateV1{
				HeadBlockHash:      lph.GetBlockHash(),
				SafeBlockHash:      lph.GetParentHash(),
				FinalizedBlockHash: lph.GetParentHash(),
			},
			PayloadAttributes: attrs,
			ForkVersion:       pb.chainSpec.GetActiveForkVersionForSlot(slot),
		},
	)
	return err
}

func (pb *PayloadBuilder[
	_, ExecutionPayloadT, _,
	_, PayloadIDT, _,
]) getPayload(
	ctx context.Context,
	payloadID PayloadIDT,
	slot math.U64,
) (engineprimitives.BuiltExecutionPayloadEnv[ExecutionPayloadT], error) {
	envelope, err := pb.ee.GetPayload(
		ctx,
		&engineprimitives.GetPayloadRequest[PayloadIDT]{
			PayloadID:   payloadID,
			ForkVersion: pb.chainSpec.GetActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	}
	if envelope == nil {
		return nil, ErrNilPayloadEnvelope
	}
	return envelope, nil
}
