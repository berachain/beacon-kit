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

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

// BuildLocalPayload builds a payload for the given slot and
// returns the payload ID.
func (s *Service) BuildLocalPayload(
	ctx context.Context,
	parentEth1Hash common.Hash,
	slot primitives.Slot,
	timestamp uint64,
	parentBeaconBlockRoot []byte,
) (*enginev1.PayloadIDBytes, error) {
	var (
		st = s.BeaconState(ctx)
		// TODO: RANDAO
		prevRandao = make([]byte, 32) //nolint:gomnd // TODO: later
		// prevRandao, err := helpers.RandaoMix(st, time.CurrentEpoch(st))
	)

	// If the local builder is disabled, we skip the payload building.
	if !s.cfg.LocalBuilderEnabled {
		s.Logger().
			Info("local builder is disabled, skipping payload building...")
		return nil, ErrLocalBuildingDisabled
	}

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		s.Logger().Error(
			"Could not get expected withdrawals to get payload attribute", "error", err)
		return nil, err
	}

	// Build the payload attributes.
	attrs, err := engine.NewPayloadAttributesContainer(
		s.ActiveForkVersionForSlot(slot),
		timestamp,
		prevRandao,
		s.BeaconCfg().Validator.SuggestedFeeRecipient[:],
		withdrawals,
		parentBeaconBlockRoot,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create payload attributes")
	}

	fcuConfig := &execution.FCUConfig{
		HeadEth1Hash:  parentEth1Hash,
		ProposingSlot: slot,
		Attributes:    attrs,
	}

	// Notify the execution client of the forkchoice update.
	var payloadID *enginev1.PayloadIDBytes
	s.Logger().Info(
		"bob the builder; can we fix it; bob the builder; yes we can 🚧",
		"slot", slot, "parent_eth1_hash", parentEth1Hash,
	)
	payloadID, err = s.es.NotifyForkchoiceUpdate(
		ctx, fcuConfig,
	)
	if err != nil {
		return nil, errors.Wrap(err, "error when notifying forkchoice update")
	} else if payloadID == nil {
		s.Logger().Error("received nil payload ID on VALID engine response",
			"head_eth1_hash", fmt.Sprintf("%#x", fcuConfig.HeadEth1Hash),
			"slot", fcuConfig.ProposingSlot,
		)
		return nil, ErrNilPayloadOnValidResponse
	}

	s.Logger().Info("forkchoice updated with payload attributes",
		"head_eth1_hash", fcuConfig.HeadEth1Hash,
		"slot", fcuConfig.ProposingSlot,
		"payload_id", fmt.Sprintf("%#x", *payloadID),
	)

	s.payloadCache.Set(
		fcuConfig.ProposingSlot,
		[32]byte(parentBeaconBlockRoot),
		primitives.PayloadID(payloadID[:]),
	)

	return payloadID, nil
}

// GetLocalPayload attempts to pull a previously built payload
// by reading a payloadID from the builder's cache. If it fails to
// retrieve a payload, it will build a new payload and wait for the
// execution client to return the payload.
func (s *Service) GetLocalPayload(
	ctx context.Context,
	slot primitives.Slot,
	parentBeaconBlockRoot []byte,
	parentEth1Hash common.Hash,
) (engine.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	// vIdx := blk.ProposerIndex()
	// headRoot := blk.ParentRoot()
	// logFields := logrus.Fields{
	// 	"validatorIndex": vIdx,
	// 	"slot":           slot,
	// 	"headRoot":       fmt.Sprintf("%#x", headRoot),
	// }

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
		parentBeaconBlockRoot,
	)
	if err == nil {
		return payload, blobsBundle, overrideBuilder, nil
	}

	s.Logger().Warn(
		"lookup failed during local payload retrieval",
		"error", err,
	)

	// Otherwise we have to trigger a new payload to be built, wait for it to be
	// resolved and then return the data. This case should very rarely be hit
	// if your consensus client is operating well.
	//
	//#nosec:G701 // won't overflow, time cannot be negative.
	return s.buildAndWaitForLocalPayload(
		ctx,
		parentEth1Hash,
		slot,
		uint64(time.Now().Unix()),
		parentBeaconBlockRoot,
	)
}

// getPayloadFromCachedPayloadIDs attempts to retrieve a payload from the
// execution client via a payload ID that is stored in the builder's cache.
func (s *Service) getPayloadFromCachedPayloadIDs(
	ctx context.Context,
	slot primitives.Slot,
	parentBeaconBlockRoot []byte,
) (engine.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	// If we have a payload ID in the cache, we can return the payload from the
	// cache.
	payloadID, found := s.payloadCache.Get(
		slot,
		[32]byte(parentBeaconBlockRoot),
	)
	if found && (payloadID != primitives.PayloadID{}) {
		// Payload ID is cache hit.
		telemetry.IncrCounter(1, MetricsPayloadIDCacheHit)
		payload, blobsBundle, overrideBuilder, err :=
			s.getPayloadFromExecutionClient(
				ctx, primitives.PayloadID(payloadID[:]), slot,
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
	parentEth1Hash common.Hash,
	slot primitives.Slot,
	timestamp uint64,
	parentBeaconBlockRoot []byte,
) (engine.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	// Build the payload and wait for the execution client to return the payload
	// ID.
	payloadID, err := s.BuildLocalPayload(
		ctx, parentEth1Hash, slot, timestamp, parentBeaconBlockRoot,
	)
	if err != nil {
		return nil, nil, false, err
	}

	// Wait for the payload to be delivered to the execution client.
	s.Logger().Info(
		"waiting for local payload to be delivered to execution client",
		"slot", slot, "timeout", s.cfg.LocalBuildPayloadTimeout.String(),
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
			ctx, primitives.PayloadID(*payloadID), slot,
		)
	if err != nil {
		return nil, nil, false, err
	}

	// TODO: Dencun
	_ = blobsBundle
	// bundleCache.add(slot, bundle)
	// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)

	s.Logger().Debug(
		"received execution payload from local engine", "value", payload.GetValue(),
	)
	return payload, blobsBundle, overrideBuilder, nil
}

// getPayloadFromExecutionClient retrieves the payload and blobs bundle for the
// given slot.
func (s *Service) getPayloadFromExecutionClient(
	ctx context.Context,
	payloadID primitives.PayloadID,
	slot primitives.Slot,
) (engine.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	payload, blobsBundle, overrideBuilder, err := s.es.GetPayload(
		ctx,
		payloadID,
		slot,
	)
	if err != nil {
		return nil, nil, false, err
	}

	s.Logger().Info("payload retrieved from local builder 🏗️ ",
		"slot", slot,
		"block_hash", common.BytesToHash(payload.GetBlockHash()),
		"parent_hash", common.BytesToHash(payload.GetParentHash()),
		"value", payload.GetValue().ToEther(),
		"override_builder", overrideBuilder,
		"num_blobs", len(blobsBundle.GetBlobs()),
	)
	return payload, blobsBundle, overrideBuilder, err
}
