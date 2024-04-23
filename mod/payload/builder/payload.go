// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package builder

import (
	"context"
	"fmt"
	"time"

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/execution"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/math"
)

// RequestPayload builds a payload for the given slot and
// returns the payload ID.
func (pb *PayloadBuilder) RequestPayload(
	ctx context.Context,
	st state.BeaconState,
	parentEth1Hash primitives.ExecutionHash,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot primitives.Root,
) (*engineprimitives.PayloadID, error) {
	// Assemble the payload attributes.
	attrs, err := pb.getPayloadAttribute(st, slot, timestamp, parentBlockRoot)
	if err != nil {
		return nil, fmt.Errorf("%w error when getting payload attributes", err)
	}

	// Notify the execution client of the forkchoice update.
	var payloadID *engineprimitives.PayloadID
	pb.logger.Info(
		"bob the builder; can we fix it; bob the builder; yes we can 🚧",
		"for_slot", slot,
		"parent_eth1_hash", parentEth1Hash,
		"parent_block_root", parentBlockRoot,
	)

	latestExecutionPayload, err := st.GetLatestExecutionPayload()
	if err != nil {
		return nil, err
	}
	parentEth1BlockHash := latestExecutionPayload.GetBlockHash()

	payloadID, _, err = pb.ee.NotifyForkchoiceUpdate(
		ctx, &execution.ForkchoiceUpdateRequest{
			State: &engineprimitives.ForkchoiceState{
				HeadBlockHash:      parentEth1Hash,
				SafeBlockHash:      parentEth1BlockHash,
				FinalizedBlockHash: parentEth1BlockHash,
			},
			PayloadAttributes: attrs,
			ForkVersion:       pb.chainSpec.ActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	} else if payloadID == nil {
		pb.logger.Warn("received nil payload ID on VALID engine response",
			"head_eth1_hash", parentEth1Hash,
			"for_slot", slot,
		)

		return payloadID, ErrNilPayloadOnValidResponse
	}

	pb.logger.Info("forkchoice updated with payload attributes",
		"head_eth1_hash", parentEth1Hash,
		"for_slot", slot,
		"payload_id", payloadID,
	)

	pb.pc.Set(
		slot,
		parentBlockRoot,
		*payloadID,
	)

	return payloadID, nil
}

// RetrieveBuiltPayload attempts to pull a previously built payload
// by reading a payloadID from the builder's cache. If it fails to
// retrieve a payload, it will build a new payload and wait for the
// execution client to return the payload.
func (pb *PayloadBuilder) RetrieveBuiltPayload(
	ctx context.Context,
	st state.BeaconState,
	slot math.Slot,
	parentBlockRoot primitives.Root,
	parentEth1Hash primitives.ExecutionHash,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// TODO: Proposer-Builder Separation Improvements Later.
	// val, tracked := s.TrackedValidatorsCache.Validator(vIdx)
	// if !tracked {
	// 	logrus.WithFields(logFields).Warn("could not find tracked proposer
	// index")
	// }

	// We first attempt to see if we previously fired off a payload built for
	// this particular slot and parent block root. If we have, and we are able
	// to
	// retrieve it from our execution client, we can return it immediately.
	envelope, err := pb.
		retrieveBuiltPayload(
			ctx,
			slot,
			parentBlockRoot,
		)

	// If there was no error we can simply return the payload that we
	// just retrieved.
	if err == nil {
		return envelope, nil
	}

	// Otherwise we will fall back to triggering a payload build.
	return pb.buildAndWaitForLocalPayload(
		ctx,
		st,
		parentEth1Hash,
		slot,
		// TODO: we need to do the proper timestamp math here for EIP4788.
		//#nosec:G701 // won't realistically overflow.
		uint64(time.Now().Unix()),
		parentBlockRoot,
	)
}

// retrieveBuiltPayload retrieves the payload and blobs bundle
// from the execution client.
func (pb *PayloadBuilder) retrieveBuiltPayload(
	ctx context.Context,
	slot math.Slot,
	parentBlockRoot primitives.Root,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// See if we have a payload ID for this slot and parent block root.
	payloadID, found := pb.pc.Get(slot, parentBlockRoot)
	if !found || (payloadID == engineprimitives.PayloadID{}) {
		// If we don't have a payload ID, we can't retrieve the payload.
		return nil, ErrPayloadIDNotFound
	}

	// Request the payload from the execution client.
	envelope, err := pb.getPayloadFromExecutionClient(
		ctx, &payloadID, slot,
	)
	if err != nil {
		return nil, err
	} else if envelope == nil {
		return nil, ErrNilPayloadEnvelope
	}

	// Cache the payload and return.
	payload := envelope.GetExecutionPayload()
	if payload == nil || payload.IsNil() {
		return nil, ErrNilPayload
	}
	pb.pc.Set(slot, payload.GetParentHash(), payloadID)
	return envelope, nil
}

// buildAndWaitForLocalPayload, triggers a payload build process, waits
// for a configuration specified period, and then retrieves the built
// payload from the execution client.
func (pb *PayloadBuilder) buildAndWaitForLocalPayload(
	ctx context.Context,
	st state.BeaconState,
	parentEth1Hash primitives.ExecutionHash,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot primitives.Root,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// Build the payload and wait for the execution client to return the payload
	// ID.
	payloadID, err := pb.RequestPayload(
		ctx, st, parentEth1Hash, slot, timestamp, parentBlockRoot,
	)
	if err != nil {
		return nil, err
	}

	// Wait for the payload to be delivered to the execution client.
	pb.logger.Info(
		"waiting for local payload to be delivered to execution client",
		"for_slot", slot, "timeout", pb.cfg.PayloadTimeout.String(),
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
	return pb.getPayloadFromExecutionClient(
		ctx, payloadID, slot)
}

// getPayloadFromExecutionClient retrieves the payload and blobs bundle for the
// given slot.
func (pb *PayloadBuilder) getPayloadFromExecutionClient(
	ctx context.Context,
	payloadID *engineprimitives.PayloadID,
	slot math.Slot,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	if payloadID == nil {
		return nil, ErrNilPayloadID
	}

	envelope, err := pb.ee.GetPayload(
		ctx,
		&execution.GetPayloadRequest{
			PayloadID:   *payloadID,
			ForkVersion: pb.chainSpec.ActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	} else if envelope == nil {
		return nil, ErrNilPayloadEnvelope
	}

	overrideBuilder := envelope.ShouldOverrideBuilder()
	args := []any{
		"for_slot", slot,
		"override_builder", overrideBuilder,
	}

	payload := envelope.GetExecutionPayload()
	if payload != nil && !payload.IsNil() {
		args = append(args,
			"payload_block_hash", payload.GetBlockHash(),
			"parent_hash", payload.GetParentHash(),
		)
	}

	blobsBundle := envelope.GetBlobsBundle()
	if blobsBundle != nil {
		args = append(args, "num_blobs", len(blobsBundle.GetBlobs()))
	}

	pb.logger.Info("payload retrieved from local builder 🏗️ ", args...)
	return envelope, err
}
