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

package execution

import (
	"context"
	"errors"
	"time"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

// TODO: this function is dog and retries need to be managed better in general.
func (s *Service) notifyForkchoiceUpdateWithSyncingRetry(
	ctx context.Context, slot primitives.Slot, arg *NotifyForkchoiceUpdateArg, withAttrs bool,
) (*primitives.PayloadID, error) {
retry:
	payloadID, err := s.notifyForkchoiceUpdate(ctx, slot, arg, withAttrs)
	if err != nil {
		if errors.Is(err, execution.ErrAcceptedSyncingPayloadStatus) {
			s.logger.Info("retrying forkchoice update", "reason", err)
			time.Sleep(forkchoiceBackoff)
			goto retry
		}
		s.logger.Error("failed to notify forkchoice update", "error", err)
	}
	return payloadID, err
}

func (s *Service) notifyForkchoiceUpdate(ctx context.Context,
	slot primitives.Slot, arg *NotifyForkchoiceUpdateArg, withAttrs bool,
) (*primitives.PayloadID, error) {
	fc := &enginev1.ForkchoiceState{
		HeadBlockHash:      arg.headHash.Bytes(),
		SafeBlockHash:      arg.safeHash.Bytes(),
		FinalizedBlockHash: arg.finalHash.Bytes(),
	}

	// Cache payloads if we get a payloadID in our response.
	var payloadID *primitives.PayloadID
	defer func() {
		if payloadID != nil {
			s.payloadCache.Set(slot, arg.headHash, *payloadID)
		}
	}()

	// We want to start building the next block as part of this forkchoice update.
	// nextSlot := slot + 1 // Cache payload ID for next slot proposer.
	nextSlot := slot
	var attrs payloadattribute.Attributer
	var err error
	if withAttrs {
		// todo: handle version.
		attrs, err = s.getPayloadAttributes(ctx, nextSlot, uint64(time.Now().Unix()))
		if err != nil {
			s.logger.Error("failed to get payload attributes in notifyForkchoiceUpdated", "error", err)
			return nil, err
		}
	} else {
		attrs = payloadattribute.EmptyWithVersion(s.beaconCfg.ActiveForkVersion(primitives.Epoch(slot)))
	}
	// TODO: remember and figure out what the middle param is.
	payloadID, _, err = s.engine.ForkchoiceUpdated(ctx, fc, attrs)
	if err != nil {
		// TODO: ensure this switch statement isn't fucked.
		switch err { //nolint:errorlint // okay for now.
		case execution.ErrAcceptedSyncingPayloadStatus:
			return payloadID, err
		case execution.ErrInvalidPayloadStatus:
			s.logger.Error("invalid payload status", "error", err)
			// TODO: Get last valid is kinda hood, its just a ptr in mem rn.
			previousHead := s.fcsp.ForkChoiceStore(ctx).GetLastValidHead()
			payloadID, err = s.notifyForkchoiceUpdate(ctx, slot, &NotifyForkchoiceUpdateArg{
				headHash: previousHead,
			}, withAttrs)

			if err != nil {
				// TODO: if u hit here, you're cooked.
				return nil, err // Returning err because it's recursive here.
			}

			// if err := s.saveHead(ctx, r, b, st); err != nil {
			// 	log.WithError(err).Error("could not save head after pruning invalid blocks")
			// }

			// log.WithFields(logrus.Fields{
			// 	"slot":                 headBlk.Slot(),
			// 	"blockRoot":            fmt.Sprintf("%#x", bytesutil.Trunc(headRoot[:])),
			// 	"invalidChildrenCount": len(invalidRoots),
			// 	"newHeadRoot":          fmt.Sprintf("%#x", bytesutil.Trunc(r[:])),
			// }).Warn("Pruned invalid blocks")
			return payloadID, errors.New("invalid payload")
			// return pid, invalidBlock{error: ErrInvalidPayload,
			//root: arg.headRoot, invalidAncestorRoots: invalidRoots}

		default:
			s.logger.Error("undefined execution engine error", "error", err)
			return nil, err
		}
	}
	// forkchoiceUpdatedValidNodeCount.Inc()
	//
	//	if err := s.cfg.ForkChoiceStore.SetOptimisticToValid(ctx, arg.headRoot); err != nil {
	//		log.WithError(err).Error("Could not set head root to valid")
	//		return nil, nil
	//	}
	//
	// If the forkchoice update call has an attribute, update the proposer payload ID cache.
	//
	//	if hasAttr && payloadID != nil {
	//		var pId [8]byte
	//		copy(pId[:], payloadID[:])
	//		log.WithFields(logrus.Fields{
	//			"blockRoot": fmt.Sprintf("%#x", bytesutil.Trunc(arg.headRoot[:])),
	//			"headSlot":  headBlk.Slot(),
	//			"payloadID": fmt.Sprintf("%#x", bytesutil.Trunc(payloadID[:])),
	//		}).Info("Forkchoice updated with payload attributes for proposal")
	//		s.cfg.ProposerSlotIndexCache.SetProposerAndPayloadIDs(nextSlot, proposerId, pId, arg.headRoot)
	//	} else if hasAttr && payloadID == nil && !features.Get().PrepareAllPayloads {
	//
	//		log.WithFields(logrus.Fields{
	//			"blockHash": fmt.Sprintf("%#x", headPayload.BlockHash()),
	//			"slot":      headBlk.Slot(),
	//		}).Error("Received nil payload ID on VALID engine response")
	//	}
	return payloadID, nil
}
