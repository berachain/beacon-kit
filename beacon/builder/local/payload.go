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
	"errors"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// BuildLocalPayload builds a payload for the given slot and
// returns the payload ID.
func (s *Service) BuildLocalPayload(
	ctx context.Context,
	parentEth1Hash common.Hash,
	slot primitives.Slot,
	timestamp uint64,
	parentBlockRoot [32]byte,
) (*enginev1.PayloadIDBytes, error) {
	// Assemble the payload attributes.
	attrs, err := s.getPayloadAttribute(
		ctx, slot, timestamp, parentBlockRoot,
	)
	if err != nil {
		return nil, fmt.Errorf("%w error when getting payload attributes", err)
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
		"for_slot", slot,
		"parent_eth1_hash", parentEth1Hash,
		"parent_block_root", common.Hash(parentBlockRoot),
	)
	payloadID, err = s.es.NotifyForkchoiceUpdate(
		ctx, fcuConfig,
	)
	if err != nil {
		return nil, err
	} else if payloadID == nil {
		s.Logger().Warn("received nil payload ID on VALID engine response",
			"head_eth1_hash", fmt.Sprintf("%#x", fcuConfig.HeadEth1Hash),
			"for_slot", fcuConfig.ProposingSlot,
		)

		s.SetStatus(ErrNilPayloadOnValidResponse)
		return payloadID, ErrNilPayloadOnValidResponse
	}

	s.Logger().Info("forkchoice updated with payload attributes",
		"head_eth1_hash", fcuConfig.HeadEth1Hash,
		"for_slot", fcuConfig.ProposingSlot,
		"payload_id", fmt.Sprintf("%#x", *payloadID),
	)

	s.payloadCache.Set(
		fcuConfig.ProposingSlot,
		parentBlockRoot,
		primitives.PayloadID(payloadID[:]),
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
	parentBlockRoot [32]byte,
	parentEth1Hash common.Hash,
) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
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
		s.Logger().Warn(fmt.Sprintf(
			"%s, notifying execution client to construct a new payload ...",
			err.Error(),
		))

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
	parentBlockRoot [32]byte,
) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	// If we have a payload ID in the cache, we can return the payload from the
	// cache.
	payloadID, found := s.payloadCache.Get(slot, parentBlockRoot)
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
	parentBlockRoot [32]byte,
) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
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
		st = s.BeaconState(ctx)
		// TODO: RANDAO
		prevRandao = make([]byte, 32) //nolint:gomnd // TODO: later
		// prevRandao, err := helpers.RandaoMix(st, time.CurrentEpoch(st))
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		s.Logger().Error(
			"Could not get expected withdrawals to get payload attribute", "error", err)
		return nil, err
	}

	// Build the payload attributes.
	attrs, err := enginetypes.NewPayloadAttributesContainer(
		s.ActiveForkVersionForSlot(slot),
		timestamp,
		prevRandao,
		s.BeaconCfg().Validator.SuggestedFeeRecipient[:],
		withdrawals,
		prevHeadRoot,
	)
	if err != nil {
		return nil, errors.New("could not create payload attributes")
	}

	return attrs, nil
}

// getPayloadFromExecutionClient retrieves the payload and blobs bundle for the
// given slot.
func (s *Service) getPayloadFromExecutionClient(
	ctx context.Context,
	payloadID primitives.PayloadID,
	slot primitives.Slot,
) (enginetypes.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	payload, blobsBundle, overrideBuilder, err := s.es.GetPayload(
		ctx,
		payloadID,
		slot,
	)
	if err != nil {
		return nil, nil, false, err
	}

	s.Logger().Info("payload retrieved from local builder 🏗️ ",
		"for_slot", slot,
		"block_hash", common.BytesToHash(payload.GetBlockHash()),
		"parent_hash", common.BytesToHash(payload.GetParentHash()),
		"value", payload.GetValue().ToEther(),
		"override_builder", overrideBuilder,
		"num_blobs", len(blobsBundle.GetBlobs()),
	)
	return payload, blobsBundle, overrideBuilder, err
}
