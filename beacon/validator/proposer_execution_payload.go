// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package validator

import (
	"context"
	"fmt"
	"time"

	"github.com/itsdevbear/bolaris/types/consensus/v1/interfaces"
	"github.com/itsdevbear/bolaris/types/state"
	"github.com/pkg/errors"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
	"github.com/prysmaticlabs/prysm/v4/runtime/version"
)

//nolint:funlen,gocognit // TODO FIX
func (s *Service) getLocalPayload(
	ctx context.Context,
	blk interfaces.ReadOnlyBeaconKitBlock, st state.BeaconState,
) (interfaces.ExecutionData, bool, error) {
	if blk.Version() < version.Bellatrix {
		return nil, false, nil
	}

	slot := blk.GetSlot()
	// vIdx := blk.ProposerIndex()
	// headRoot := blk.ParentRoot()
	// logFields := logrus.Fields{
	// 	"validatorIndex": vIdx,
	// 	"slot":           slot,
	// 	"headRoot":       fmt.Sprintf("%#x", headRoot),
	// }
	// p, err := blk.Execution()
	// if err != nil {
	// 	return nil, false, err
	// }

	// parentHash := p.ParentHash()
	// WILL ALWAYS MISS
	payloadID, ok := s.PayloadIDCache.PayloadID(slot, [32]byte{})

	// val, tracked := s.TrackedValidatorsCache.Validator(vIdx)
	// if !tracked {
	// 	logrus.WithFields(logFields).Warn("could not find tracked proposer index")
	// }

	// If we have a payload ID in the cache, we can return the payload from the cache.
	if ok && payloadID != [8]byte{} {
		// Payload ID is cache hit. Return the cached payload ID.
		var pid primitives.PayloadID
		copy(pid[:], payloadID[:])
		// payloadIDCacheHit.Inc()
		var payload interfaces.ExecutionData
		var overrideBuilder bool
		payload, _, overrideBuilder, err := s.en.GetPayload(ctx, pid, slot)
		switch {
		case err == nil:
			// bundleCache.add(slot, bundle)
			// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)
			return payload, overrideBuilder, nil
		case errors.Is(err, context.DeadlineExceeded):
		default:
			return nil, false, errors.Wrap(err, "could not get cached payload from execution client")
		}
	}

	// Otherwise we did not have a payload in the cache and we must build a new payload.

	// log.WithFields(logFields).Debug("payload ID cache miss")vs
	var parentHash []byte
	var err error
	parentHash, err = s.getParentBlockHash(ctx)
	if err != nil {
		return nil, false, err
	}

	// payloadIDCacheMiss.Inc()

	// random, err := helpers.RandaoMix(st, time.CurrentEpoch(st))
	// if err != nil {
	// 	return nil, false, err
	// }
	random := make([]byte, 32)   //nolint:gomnd // todo: randao
	headRoot := make([]byte, 32) //nolint:gomnd // todo: cancaun
	justifiedBlockHash := s.BeaconState(ctx).GetSafeEth1BlockHash()
	finalizedBlockHash := s.BeaconState(ctx).GetFinalizedEth1BlockHash()

	f := &enginev1.ForkchoiceState{
		HeadBlockHash:      parentHash,
		SafeBlockHash:      justifiedBlockHash[:],
		FinalizedBlockHash: finalizedBlockHash[:],
	}

	t := time.Now()
	var (
		attr        payloadattribute.Attributer
		withdrawals []*enginev1.Withdrawal
	)
	switch st.Version() {
	case version.Deneb:
		withdrawals, err = st.ExpectedWithdrawals()
		if err != nil {
			return nil, false, err
		}
		attr, err = payloadattribute.New(&enginev1.PayloadAttributesV3{
			Timestamp:             uint64(t.Unix()),
			PrevRandao:            random,
			SuggestedFeeRecipient: s.BeaconCfg().Validator.SuggestedFeeRecipient[:],
			Withdrawals:           withdrawals,
			ParentBeaconBlockRoot: headRoot,
		})
		if err != nil {
			return nil, false, err
		}
	case version.Capella:
		withdrawals, err = st.ExpectedWithdrawals()
		if err != nil {
			return nil, false, err
		}
		attr, err = payloadattribute.New(&enginev1.PayloadAttributesV2{
			Timestamp:             uint64(t.Unix()),
			PrevRandao:            random,
			SuggestedFeeRecipient: s.BeaconCfg().Validator.SuggestedFeeRecipient[:],
			Withdrawals:           withdrawals,
		})
		if err != nil {
			return nil, false, err
		}
	case version.Bellatrix:
		attr, err = payloadattribute.New(&enginev1.PayloadAttributes{
			Timestamp:             uint64(t.Unix()),
			PrevRandao:            random,
			SuggestedFeeRecipient: s.BeaconCfg().Validator.SuggestedFeeRecipient[:],
		})
		if err != nil {
			return nil, false, err
		}
	default:
		return nil, false, errors.New("unknown beacon state version")
	}

	var payloadIDBytes *enginev1.PayloadIDBytes
	payloadIDBytes, _, err = s.en.ForkchoiceUpdated(ctx, f, attr)
	if err != nil {
		return nil, false, errors.Wrap(err, "could not prepare payload")
	}
	if payloadIDBytes == nil {
		return nil, false, fmt.Errorf("nil payload with block hash: %#x", parentHash)
	}

	payload, _, overrideBuilder, err := s.en.GetPayload(ctx, *payloadIDBytes, slot)
	if err != nil {
		return nil, false, err
	}
	// bundleCache.add(slot, bundle)
	// warnIfFeeRecipientDiffers(payload, val.FeeRecipient)
	localValueGwei, err := payload.ValueInGwei()
	if err == nil {
		s.Logger().Debug("received execution payload from local engine", "value", localValueGwei)
	}
	return payload, overrideBuilder, nil
}

// getParentBlockHash retrieves the parent block hash for the given slot.
//
//nolint:unparam // todo: review this later.
func (s *Service) getParentBlockHash(ctx context.Context) ([]byte, error) {
	// The first slot should be proposed with the genesis block as parent.
	st := s.BeaconState(ctx)
	if st.Slot() == 1 {
		return st.GenesisEth1Hash().Bytes(), nil
	}

	// We always want the parent block to be the last finalized block.
	return st.GetFinalizedEth1BlockHash().Bytes(), nil
}
