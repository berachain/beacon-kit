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
	"fmt"
	"time"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	engineerrors "github.com/berachain/beacon-kit/engine-primitives/errors"
	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
)

type RequestPayloadData struct {
	Slot                 math.Slot
	Timestamp            math.U64
	PayloadWithdrawals   engineprimitives.Withdrawals
	PrevRandao           common.Bytes32
	ParentBlockRoot      common.Root
	HeadEth1BlockHash    common.ExecutionHash
	FinalEth1BlockHash   common.ExecutionHash
	ParentProposerPubkey *crypto.BLSPubkey // nil for fork versions before Electra1
}

// RequestPayloadAsync builds a payload for the given slot and
// returns the payload ID.
func (pb *PayloadBuilder) RequestPayloadAsync(
	ctx context.Context,
	r *RequestPayloadData,
) (*engineprimitives.PayloadID, common.Version, error) {
	if !pb.Enabled() {
		return nil, common.Version{}, ErrPayloadBuilderDisabled
	}

	if payloadID, found := pb.pc.Get(r.Slot, r.ParentBlockRoot); found {
		pb.logger.Info(
			"aborting payload build; payload already exists in cache",
			"for_slot", r.Slot.Base10(),
			"parent_block_root", r.ParentBlockRoot,
		)
		return &payloadID.PayloadID, payloadID.ForkVersion, nil
	}

	// Assemble the payload attributes.
	attrs, err := pb.attributesFactory.BuildPayloadAttributes(
		r.Timestamp,
		r.PayloadWithdrawals,
		r.PrevRandao,
		r.ParentBlockRoot,
		r.ParentProposerPubkey,
	)
	if err != nil {
		return nil, common.Version{}, err
	}

	forkVersion := pb.chainSpec.ActiveForkVersionForTimestamp(r.Timestamp)
	// Submit the forkchoice update to the execution client.
	req := ctypes.BuildForkchoiceUpdateRequest(
		&engineprimitives.ForkchoiceStateV1{
			HeadBlockHash:      r.HeadEth1BlockHash,
			SafeBlockHash:      r.FinalEth1BlockHash,
			FinalizedBlockHash: r.FinalEth1BlockHash,
		},
		attrs,
		forkVersion,
	)
	payloadID, err := pb.ee.NotifyForkchoiceUpdate(ctx, req)
	if err != nil {
		return nil, common.Version{}, fmt.Errorf("RequestPayloadAsync failed sending forkchoice update: %w", err)
	}

	// Only add to cache if we received back a payload ID.
	if payloadID != nil {
		pb.pc.Set(r.Slot, r.ParentBlockRoot, *payloadID, forkVersion)
	}

	return payloadID, forkVersion, nil
}

// RequestPayloadSync request a payload for the given slot and
// blocks until the payload is delivered.
func (pb *PayloadBuilder) RequestPayloadSync(
	ctx context.Context,
	r *RequestPayloadData,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	if !pb.Enabled() {
		return nil, ErrPayloadBuilderDisabled
	}

	// Build the payload and wait for the execution client to
	// return the payload ID.
	payloadID, forkVersion, err := pb.RequestPayloadAsync(ctx, r)
	if err != nil {
		return nil, err
	}
	if payloadID == nil {
		return nil, ErrNilPayloadID
	}

	// Wait for the payload to be delivered to the execution client.
	pb.logger.Info(
		"Waiting for local payload to be delivered to execution client",
		"for_slot", r.Slot.Base10(), "timeout", pb.cfg.PayloadTimeout.String(),
	)
	select {
	case <-time.After(pb.cfg.PayloadTimeout):
		// We want to trigger delivery of the payload to the execution client
		// before the timestamp expires.
		break
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	payload, err := pb.getPayload(ctx, *payloadID, forkVersion)
	if err != nil {
		if errors.Is(err, engineerrors.ErrUnknownPayload) {
			// We may have cached the payloadID, but the payload have become stale
			// in the EVM, or there could have been other issues forcing the EVM
			// to provide no payload. In both cases we can purge the payloadID as it
			// is not usable anymore.
			pb.pc.Delete(r.Slot, r.ParentBlockRoot)
		}
		return nil, fmt.Errorf("failed retrieving payload for ID %x: %w", *payloadID, err)
	}
	return payload, nil
}

// RetrievePayload attempts to pull a previously built payload
// by reading a payloadID from the builder's cache. If it fails to
// retrieve a payload, it will build a new payload and wait for the
// execution client to return the payload.
func (pb *PayloadBuilder) RetrievePayload(
	ctx context.Context,
	slot math.Slot,
	parentBlockRoot common.Root,
	expectedForkVersion common.Version,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	if !pb.Enabled() {
		return nil, ErrPayloadBuilderDisabled
	}

	// Attempt to see if we previously fired off a payload built for
	// this particular slot and parent block root.
	payloadRes, found := pb.pc.Get(slot, parentBlockRoot)
	if !found {
		return nil, ErrPayloadIDNotFound
	}
	if !version.Equals(payloadRes.ForkVersion, expectedForkVersion) {
		pb.pc.Delete(slot, parentBlockRoot)
		return nil, ErrPayloadIDNotFound // force payload rebuild with the right fork
	}

	// Get the payload from the execution client.
	envelope, err := pb.getPayload(ctx, payloadRes.PayloadID, payloadRes.ForkVersion)
	if err != nil {
		if errors.Is(err, engineerrors.ErrUnknownPayload) {
			// We may have cached the payloadID, but the payload have become stale
			// in the EVM, or there could have been other issues. Block builder will
			// try and build again the payload just in time
			pb.pc.Delete(slot, parentBlockRoot)
		}
		return nil, err
	}

	// If the payload was built by a different builder, something is
	// wrong the EL<>CL setup.
	payload := envelope.GetExecutionPayload()
	if payload.GetFeeRecipient() != pb.cfg.SuggestedFeeRecipient {
		pb.logger.Warn(
			"Payload fee recipient does not match suggested fee recipient - "+
				"please check both your CL and EL configuration",
			"payload_fee_recipient", payload.GetFeeRecipient(),
			"suggested_fee_recipient", pb.cfg.SuggestedFeeRecipient,
		)
	}

	// log some data
	args := []any{
		"for_slot", slot.Base10(),
		"override_builder", envelope.ShouldOverrideBuilder(),
		"payload_block_hash", payload.GetBlockHash(),
		"parent_hash", payload.GetParentHash(),
	}
	if blobsBundle := envelope.GetBlobsBundle(); blobsBundle != nil {
		args = append(args, "num_blobs", len(blobsBundle.GetBlobs()))
	}
	pb.logger.Info("Payload retrieved from local builder", args...)

	return envelope, err
}

func (pb *PayloadBuilder) getPayload(
	ctx context.Context,
	payloadID engineprimitives.PayloadID,
	forkVersion common.Version,
) (ctypes.BuiltExecutionPayloadEnv, error) {
	envelope, err := pb.ee.GetPayload(
		ctx,
		&ctypes.GetPayloadRequest{
			PayloadID:   payloadID,
			ForkVersion: forkVersion,
		},
	)
	if err != nil {
		return nil, err
	}
	if envelope == nil {
		return nil, ErrNilPayloadEnvelope // appease linter. This is checked already
	}
	if envelope.GetExecutionPayload().Withdrawals == nil {
		return nil, ErrNilWithdrawals
	}
	return envelope, nil
}
