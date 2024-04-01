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

package localbuilder

import (
	"context"
	"fmt"
	"time"

	"github.com/berachain/beacon-kit/mod/execution"
	enginetypes "github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/engine"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// BuildLocalPayload builds a payload for the given slot and
// returns the payload ID.
func (s *Service) BuildLocalPayload(
	ctx context.Context,
	parentEth1Hash primitives.ExecutionHash,
	slot primitives.Slot,
	timestamp uint64,
	parentBlockRoot primitives.Root,
) (*engine.PayloadID, error) {
	// Assemble the payload attributes.
	attrs, err := s.getPayloadAttribute(
		ctx, slot, timestamp, parentBlockRoot,
	)
	if err != nil {
		return nil, fmt.Errorf("%w error when getting payload attributes", err)
	}

	// Notify the execution client of the forkchoice update.
	var payloadID *engine.PayloadID
	s.Logger().Info(
		"bob the builder; can we fix it; bob the builder; yes we can 🚧",
		"for_slot", slot,
		"parent_eth1_hash", parentEth1Hash,
		"parent_block_root", parentBlockRoot,
	)

	parentEth1BlockHash, err := s.BeaconState(ctx).GetEth1BlockHash()
	if err != nil {
		return nil, err
	}
	payloadID, _, err = s.ee.NotifyForkchoiceUpdate(
		ctx, &execution.ForkchoiceUpdateRequest{
			State: &engine.ForkchoiceState{
				HeadBlockHash:      parentEth1Hash,
				SafeBlockHash:      parentEth1BlockHash,
				FinalizedBlockHash: parentEth1BlockHash,
			},
			PayloadAttributes: attrs,
			ForkVersion:       s.ActiveForkVersionForSlot(slot),
		},
	)
	if err != nil {
		return nil, err
	} else if payloadID == nil {
		s.Logger().Warn("received nil payload ID on VALID engine response",
			"head_eth1_hash", parentEth1Hash,
			"for_slot", slot,
		)

		s.SetStatus(ErrNilPayloadOnValidResponse)
		return payloadID, ErrNilPayloadOnValidResponse
	}

	s.Logger().Info("forkchoice updated with payload attributes",
		"head_eth1_hash", parentEth1Hash,
		"for_slot", slot,
		"payload_id", payloadID,
	)

	s.pc.Set(
		slot,
		parentBlockRoot,
		*payloadID,
	)

	s.SetStatus(nil)
	return payloadID, nil
}

// GetBestPayload attempts to pull a previously built payload
// by reading a payloadID from the builder's cache. If it fails to
// retrieve a payload, it will build a new payload and wait for the
// execution client to return the payload.
func (s *Service) GetBestPayload(
	ctx context.Context,
	slot primitives.Slot,
	parentBlockRoot primitives.Root,
	parentEth1Hash primitives.ExecutionHash,
) (enginetypes.ExecutionPayload, *engine.BlobsBundleV1, bool, error) {
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
	payload, blobsBundle, overrideBuilder, err := s.getPayloadFromCachedPayloadIDs(
		ctx,
		slot,
		parentBlockRoot,
	)
	if err != nil {
		// If we see an error we have to trigger a new payload to be built, wait
		// for it to be resolved and then return the data. This case should very
		// rarely be hit
		// if your consensus and execution clients are operating well.
		s.Logger().Warn(
			err.Error() +
				": notifying execution client to construct a new payload ...",
		)

		//#nosec:G701 // won't overflow, time cannot be negative.
		payload, blobsBundle, overrideBuilder, err = s.buildAndWaitForLocalPayload(
			ctx,
			parentEth1Hash,
			slot,
			uint64(time.Now().Unix()),
			parentBlockRoot,
		)
	}

	return payload, blobsBundle, overrideBuilder, err
}

// getPayloadFromCachedPayloadIDs attempts to retrieve a payload from the
// execution client via a payload ID that is stored in the builder's cache.
func (s *Service) getPayloadFromCachedPayloadIDs(
	ctx context.Context,
	slot primitives.Slot,
	parentBlockRoot primitives.Root,
) (enginetypes.ExecutionPayload, *engine.BlobsBundleV1, bool, error) {
	// If we have a payload ID in the cache, we can return the payload from the
	// cache.
	payloadID, found := s.pc.Get(slot, parentBlockRoot)
	if found && (payloadID != engine.PayloadID{}) {
		// Payload ID is cache hit.
		telemetry.IncrCounter(1, MetricsPayloadIDCacheHit)
		payload, blobsBundle, overrideBuilder, err :=
			s.getPayloadFromExecutionClient(
				ctx, &payloadID, slot,
			)
		if err == nil {
			// bundleCache.add(slot, bundle)
			// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)
			//  Return the cached payload ID.
			return payload, blobsBundle, overrideBuilder, nil
		}

		telemetry.IncrCounter(1, MetricsPayloadIDCacheError)
		return nil, nil, false, ErrCachedPayloadNotFoundOnExecutionClient
	}
	return nil, nil, false, ErrPayloadIDNotFound
}

// buildAndWaitForLocalPayload, triggers a payload build process, waits
// for a configuration specified period, and then retrieves the built
// payload from the execution client.
func (s *Service) buildAndWaitForLocalPayload(
	ctx context.Context,
	parentEth1Hash primitives.ExecutionHash,
	slot primitives.Slot,
	timestamp uint64,
	parentBlockRoot primitives.Root,
) (enginetypes.ExecutionPayload, *engine.BlobsBundleV1, bool, error) {
	// Build the payload and wait for the execution client to return the payload
	// ID.
	payloadID, err := s.BuildLocalPayload(
		ctx, parentEth1Hash, slot, timestamp, parentBlockRoot,
	)
	if err != nil {
		return nil, nil, false, err
	}

	// Wait for the payload to be delivered to the execution client.
	s.Logger().Info(
		"waiting for local payload to be delivered to execution client",
		"for_slot", slot, "timeout", s.cfg.LocalBuildPayloadTimeout.String(),
	)
	select {
	case <-time.After(s.cfg.LocalBuildPayloadTimeout):
		// We want to trigger delivery of the payload to the execution client
		// before the timestamp expires.
		break
	case <-ctx.Done():
		return nil, nil, false, ctx.Err()
	}

	// Get the payload from the execution client.
	payload, blobsBundle, overrideBuilder, err :=
		s.getPayloadFromExecutionClient(
			ctx, payloadID, slot,
		)
	if err != nil {
		return nil, nil, false, err
	}

	// TODO: Dencun
	_ = blobsBundle
	// bundleCache.add(slot, bundle)
	// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)

	// s.Logger().Debug(
	// 	"received execution payload from local engine", "value",
	// payload.GetValue(),
	// )
	return payload, blobsBundle, overrideBuilder, nil
}

// getPayloadAttributes returns the payload attributes for the given state and
// slot. The attribute is required to initiate a payload build process in the
// context of an `engine_forkchoiceUpdated` call.
func (s *Service) getPayloadAttribute(
	ctx context.Context,
	slot primitives.Slot,
	timestamp uint64,
	prevHeadRoot [32]byte,
) (enginetypes.PayloadAttributer, error) {
	var (
		prevRandao [32]byte
		st         = s.BeaconState(ctx)
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals(
		s.BeaconCfg().MaxWithdrawalsPerPayload,
	)
	if err != nil {
		s.Logger().Error(
			"Could not get expected withdrawals to get payload attribute", "error", err)
		return nil, err
	}

	epoch := s.BeaconCfg().SlotToEpoch(slot)

	// Get the previous randao mix.
	prevRandao, err = st.GetRandaoMixAtIndex(
		uint64(epoch) % s.BeaconCfg().EpochsPerHistoricalVector,
	)
	if err != nil {
		return nil, err
	}

	return enginetypes.NewPayloadAttributes(
		s.BeaconCfg().ActiveForkVersionByEpoch(epoch),
		timestamp,
		prevRandao,
		s.cfg.SuggestedFeeRecipient,
		withdrawals,
		prevHeadRoot,
	)
}

// getPayloadFromExecutionClient retrieves the payload and blobs bundle for the
// given slot.
func (s *Service) getPayloadFromExecutionClient(
	ctx context.Context,
	payloadID *engine.PayloadID,
	slot primitives.Slot,
) (enginetypes.ExecutionPayload, *engine.BlobsBundleV1, bool, error) {
	if payloadID == nil {
		return nil, nil, false, ErrNilPayloadID
	}

	payload, blobsBundle, overrideBuilder, err := s.ee.GetPayload(
		ctx,
		&execution.GetPayloadRequest{
			PayloadID:   *payloadID,
			ForkVersion: s.BeaconCfg().ActiveForkVersion(slot),
		},
	)
	if err != nil {
		return nil, nil, false, err
	}

	args := []any{
		"for_slot", slot,
		"override_builder", overrideBuilder,
	}

	if payload != nil && !payload.IsNil() {
		args = append(args,
			"payload_block_hash", payload.GetBlockHash(),
			"parent_hash", payload.GetParentHash(),
		)
	}

	if blobsBundle != nil {
		args = append(args, "num_blobs", len(blobsBundle.Blobs))
	}

	s.Logger().Info("payload retrieved from local builder 🏗️ ", args...)
	return payload, blobsBundle, overrideBuilder, err
}
