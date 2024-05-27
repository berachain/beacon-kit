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
	"time"

	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// RetrieveOrBuildPayload attempts to pull a previously built payload
// by reading a payloadID from the builder's cache. If it fails to
// retrieve a payload, it will build a new payload and wait for the
// execution client to return the payload.
func (pb *PayloadBuilder[BeaconStateT]) RetrieveOrBuildPayload(
	ctx context.Context,
	st BeaconStateT,
	slot math.Slot,
	parentBlockRoot primitives.Root,
	parentEth1Hash common.ExecutionHash,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// We first attempt to see if we previously fired off a payload built for
	// this particular slot and parent block root. If we have, and we are able
	// to retrieve it from our execution client, we can return it immediately.
	if payloadID, found := pb.pc.Get(slot, parentBlockRoot); found {
		if envelope, err := pb.getPayload(
			ctx,
			slot,
			payloadID,
		); err == nil {
			// If there was no error we can simply return the payload that we
			// just retrieved.
			return envelope, nil
		}
	}

	// Otherwise we will fall back to triggering a payload build.
	return pb.RequestPayloadAndWait(
		ctx,
		st,
		slot,
		// TODO: we need to do the proper timestamp math here for EIP4788.
		//#nosec:G701 // won't realistically overflow.
		uint64(time.Now().Unix()),
		parentBlockRoot,
		parentEth1Hash,
	)
}

// RequestPayload builds a payload for the given slot and
// returns the payload ID.
func (pb *PayloadBuilder[BeaconStateT]) RequestPayload(
	ctx context.Context,
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot primitives.Root,
	parentEth1Hash common.ExecutionHash,
) (*engineprimitives.PayloadID, error) {
	pb.logger.Info(
		"bob the builder; can we fix it; bob the builder; yes we can 🚧",
		"for_slot", slot,
		"parent_eth1_hash", parentEth1Hash,
		"parent_block_root", parentBlockRoot,
	)

	// Assemble the payload attributes.
	attrs, err := pb.getPayloadAttribute(st, slot, timestamp, parentBlockRoot)
	if err != nil {
		return nil, errors.Newf("%w error when getting payload attributes", err)
	}

	// Submit the forkchoice update to the execution client.
	var payloadID *engineprimitives.PayloadID
	payloadID, _, err = pb.submitForkchoiceUpdate(
		ctx,
		st,
		slot,
		attrs,
		parentEth1Hash,
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

	pb.pc.Set(slot, parentBlockRoot, *payloadID)
	return payloadID, nil
}

// RequestPayload request a payload for the given slot and
// blocks until the payload is delivered.
func (pb *PayloadBuilder[BeaconStateT]) RequestPayloadAndWait(
	ctx context.Context,
	st BeaconStateT,
	slot math.Slot,
	timestamp uint64,
	parentBlockRoot primitives.Root,
	parentEth1Hash common.ExecutionHash,
) (engineprimitives.BuiltExecutionPayloadEnv, error) {
	// Build the payload and wait for the execution client to return the payload
	// ID.
	payloadID, err := pb.RequestPayload(
		ctx, st, slot, timestamp, parentBlockRoot, parentEth1Hash,
	)
	if err != nil {
		return nil, err
	} else if payloadID == nil {
		return nil, ErrNilPayloadID
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
	return pb.ee.GetPayload(
		ctx,
		&engineprimitives.GetPayloadRequest{
			PayloadID:   *payloadID,
			ForkVersion: pb.chainSpec.ActiveForkVersionForSlot(slot),
		},
	)
}
