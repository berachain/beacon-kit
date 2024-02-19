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

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
	"github.com/pkg/errors"
)

func (s *Service) GetOrBuildLocalPayload(
	ctx context.Context,
	slot primitives.Slot,
) (engine.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	// vIdx := blk.ProposerIndex()
	// headRoot := blk.ParentRoot()
	// logFields := logrus.Fields{
	// 	"validatorIndex": vIdx,
	// 	"slot":           slot,
	// 	"headRoot":       fmt.Sprintf("%#x", headRoot),
	// }

	// TODO: Proposer-Builder Seperation Improvements Later.
	// val, tracked := s.TrackedValidatorsCache.Validator(vIdx)
	// if !tracked {
	// 	logrus.WithFields(logFields).Warn("could not find tracked proposer index")
	// }

	parentEth1Hash, err := s.getParentEth1Hash(ctx)
	if err != nil {
		return nil, nil, false, err
	}

	// If we have a payload ID in the cache, we can return the payload from the cache.
	payloadID, ok := s.payloadCache.Get(slot, parentEth1Hash)
	if ok && (payloadID != primitives.PayloadID{}) {
		var (
			pidCpy          primitives.PayloadID
			payload         engine.ExecutionPayload
			overrideBuilder bool
			blobsBundle     *enginev1.BlobsBundle
		)

		// Payload ID is cache hit.
		telemetry.IncrCounter(1, MetricsPayloadIDCacheHit)
		copy(pidCpy[:], payloadID[:])
		if payload, blobsBundle, overrideBuilder, err = s.en.GetPayload(ctx, pidCpy, slot); err == nil {
			// bundleCache.add(slot, bundle)
			// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)
			//  Return the cached payload ID.
			return payload, blobsBundle, overrideBuilder, nil
		}
		s.Logger().Warn("could not get cached payload from execution client", "error", err)
		telemetry.IncrCounter(1, MetricsPayloadIDCacheError)
	}

	// If we reach this point, we have a cache miss and must build a new payload.
	telemetry.IncrCounter(1, MetricsPayloadIDCacheMiss)
	s.Logger().Warn(
		"could not find payload in cache, building new payload",
		"slot", slot, "parent_eth1-hash", parentEth1Hash.Hex(),
	)

	//#nosec:G701 // won't overflow, time cannot be negative.
	return s.BuildAndWaitForLocalPayload(ctx, parentEth1Hash, slot, uint64(time.Now().Unix()))
}

func (s *Service) BuildLocalPayload(
	ctx context.Context,
	parentEth1Hash common.Hash,
	_ primitives.Slot,
	timestamp uint64,
) (*enginev1.PayloadIDBytes, error) {
	var (
		st = s.BeaconState(ctx)
		// TODO: RANDAO
		prevRandao = make([]byte, 32) //nolint:gomnd // TODO: later
		// prevRandao, err := helpers.RandaoMix(st, time.CurrentEpoch(st))
		// TODO: Cancun
		headRoot = make([]byte, 32) //nolint:gomnd // TODO: Cancun
	)

	// Get the expected withdrawals to include in this payload.
	withdrawals, err := st.ExpectedWithdrawals()
	if err != nil {
		s.Logger().Error(
			"Could not get expected withdrawals to get payload attribute", "error", err)
		return nil, err
	}

	// Build the payload attributes.
	attrs, err := engine.NewPayloadAttributesContainer(
		st.Version(),
		timestamp,
		prevRandao,
		s.BeaconCfg().Validator.SuggestedFeeRecipient[:],
		withdrawals,
		headRoot,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not create payload attributes")
	}

	// Notify the execution client of the forkchoice update.
	var payloadID *enginev1.PayloadIDBytes
	payloadID, _, err = s.en.ForkchoiceUpdated(
		ctx,
		&enginev1.ForkchoiceState{
			HeadBlockHash:      parentEth1Hash.Bytes(),
			SafeBlockHash:      s.BeaconState(ctx).GetSafeEth1BlockHash().Bytes(),
			FinalizedBlockHash: s.BeaconState(ctx).GetFinalizedEth1BlockHash().Bytes(),
		},
		attrs,
	)
	if err != nil {
		return nil, errors.Wrap(err, "could not prepare payload")
	} else if payloadID == nil {
		s.Logger().Warn(
			"local block builder received nil payload ID on VALID engine response",
		)
		return nil, fmt.Errorf("nil payload with block hash: %#x", parentEth1Hash)
	}
	return payloadID, nil
}

// GetExecutionPayload retrieves the execution payload for the given slot.
func (s *Service) BuildAndWaitForLocalPayload(
	ctx context.Context,
	parentEth1Hash common.Hash,
	slot primitives.Slot,
	timestamp uint64,
) (engine.ExecutionPayload, *enginev1.BlobsBundle, bool, error) {
	// Build the payload and wait for the execution client to return the payload ID.
	payloadID, err := s.BuildLocalPayload(ctx, parentEth1Hash, slot, timestamp)
	if err != nil {
		return nil, nil, false, err
	}

	// Calculate the duration to wait for the payload to be delivered.
	var duration time.Duration
	nowUnix := uint64(time.Now().Unix())
	if timestamp <= nowUnix {
		duration = 500 * time.Millisecond //nolint:gomnd // for now.
	} else {
		duration = time.Duration(timestamp-nowUnix) * time.Second
	}

	select {
	case <-time.After(duration):
		// We want to trigger delivery of the payload to the execution client
		// before the timestamp expires.
		break
	case <-ctx.Done():
		return nil, nil, false, ctx.Err()
	}

	// Get the payload from the execution client.
	payload, blobsBundle, overrideBuilder, err := s.en.GetPayload(
		ctx, primitives.PayloadID(*payloadID), slot,
	)
	if err != nil {
		return nil, nil, false, err
	}

	// TODO: Dencun
	_ = blobsBundle
	// bundleCache.add(slot, bundle)
	// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)

	s.Logger().Debug("received execution payload from local engine", "value", payload.GetValue())
	return payload, blobsBundle, overrideBuilder, nil
}

// getParentEth1Hash retrieves the parent block hash for the given slot.
//
//nolint:unparam // todo: review this later.
func (s *Service) getParentEth1Hash(ctx context.Context) (common.Hash, error) {
	// The first slot should be proposed with the genesis block as parent.
	st := s.BeaconState(ctx)
	if st.Slot() == 1 {
		return st.GenesisEth1Hash(), nil
	}

	// We always want the parent block to be the last finalized block.
	return st.GetFinalizedEth1BlockHash(), nil
}
