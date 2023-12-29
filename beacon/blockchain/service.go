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

package blockchain

import (
	"context"
	"crypto/rand"
	"errors"
	"time"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/interfaces"
	payloadattribute "github.com/prysmaticlabs/prysm/v4/consensus-types/payload-attribute"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"

	"cosmossdk.io/core/header"
	"cosmossdk.io/log"

	"github.com/ethereum/go-ethereum/common"

	"github.com/itsdevbear/bolaris/beacon/execution/engine"
	"github.com/itsdevbear/bolaris/beacon/execution/finalizer"
	"github.com/itsdevbear/bolaris/types"
	"github.com/itsdevbear/bolaris/types/config"
)

// var defaultLatestValidHash = bytesutil.PadTo([]byte{0xff}, 32).
const forkchoiceBackoff = 500 * time.Millisecond

type ForkChoiceStoreProvider interface {
	ForkChoiceStore(ctx context.Context) types.ForkChoiceStore
}

type Service struct {
	beaconCfg *config.Beacon
	etherbase common.Address
	logger    log.Logger
	fcsp      ForkChoiceStoreProvider
	engine    engine.Caller
	finalizer *finalizer.Finalizer
}

func NewService(opts ...Option) *Service {
	s := &Service{}
	for _, opt := range opts {
		if err := opt(s); err != nil {
			s.logger.Error("Failed to apply option", "error", err)
		}
	}
	s.finalizer = finalizer.New(s)
	s.finalizer.Start(context.Background())
	return s
}

func (s *Service) ValidateProposedBeaconBlock(ctx context.Context,
	block header.Info, header interfaces.ExecutionData,
) (*enginev1.PayloadIDBytes, error) {
	isValidPayload, err := s.validateExecutionOnBlock(ctx, 0, header /*, nil, [32]byte{}*/)
	if err != nil {
		if !errors.Is(err, execution.ErrAcceptedSyncingPayloadStatus) {
			s.logger.Error("failed to validate execution on block", "err", err)
			return nil, err
		}
	} else if !isValidPayload {
		return nil, execution.ErrInvalidPayloadStatus
	}

	// Forkchoice our execution client's head to be the block that we validated as correct
	// above. NOTE: we do not touch our safe or finalized hashes here.
	finalized := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	safe := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	return s.notifyForkchoiceUpdateWithSyncingRetry(
		ctx, primitives.Slot(block.Height), &notifyForkchoiceUpdateArg{
			headHash:  header.BlockHash(),
			safeHash:  safe[:],
			finalHash: finalized[:],
		}, true)
}

func (s *Service) notifyForkchoiceUpdateWithSyncingRetry(
	ctx context.Context, slot primitives.Slot, arg *notifyForkchoiceUpdateArg, withAttrs bool,
) (*enginev1.PayloadIDBytes, error) {
retry:
	payloadID, err := s.notifyForkchoiceUpdate(ctx, slot, arg, withAttrs)
	if err != nil {
		if errors.Is(err, execution.ErrAcceptedSyncingPayloadStatus) {
			s.logger.Info("retrying forkchoice update", "reason", err)
			time.Sleep(forkchoiceBackoff)
			goto retry
		}
		s.logger.Error("failed to notify forkchoice update", "err", err)
	}
	return payloadID, err
}

func (s *Service) notifyForkchoiceUpdate(ctx context.Context,
	slot primitives.Slot, arg *notifyForkchoiceUpdateArg, withAttrs bool,
) (*enginev1.PayloadIDBytes, error) {
	// currSafeBlk := s.fcsp.ForkChoiceStore(ctx).GetSafeBlockHash()
	// currFinalizedBlk := s.fcsp.ForkChoiceStore(ctx).GetFinalizedBlockHash()
	// TODO FIX, rn we are just blindly finalizing whatever the proposer has sent us.
	// The blind finalization is "sota safe" cause we will get an STATUS_INVALID From the
	// forkchoice update
	fc := &enginev1.ForkchoiceState{
		HeadBlockHash:      arg.headHash,
		SafeBlockHash:      arg.safeHash,
		FinalizedBlockHash: arg.finalHash,
	}

	// We want to start building the next block as part of this forkchoice update.
	// nextSlot := slot + 1 // Cache payload ID for next slot proposer.
	nextSlot := slot
	var attrs payloadattribute.Attributer
	var err error
	if withAttrs {
		attrs, err = s.getPayloadAttributes(ctx, nextSlot, uint64(time.Now().Unix()))
		if err != nil {
			s.logger.Error("failed to get payload attributes in notifyForkchoiceUpdated", "err", err)
			return nil, err
		}
	} else {
		attrs = payloadattribute.EmptyWithVersion(s.beaconCfg.ActiveForkVersion(primitives.Epoch(slot)))
	}
	payloadID, _, err := s.engine.ForkchoiceUpdated(ctx, fc, attrs)
	if err != nil {
		switch err { //nolint:errorlint // okay for now.
		case execution.ErrAcceptedSyncingPayloadStatus:
			return payloadID, err
		case execution.ErrInvalidPayloadStatus:
			s.logger.Error("invalid payload status", "err", err)
			previousHead := s.fcsp.ForkChoiceStore(ctx).GetLastValidHead()
			var pid *enginev1.PayloadIDBytes
			pid, err = s.notifyForkchoiceUpdate(ctx, slot, &notifyForkchoiceUpdateArg{
				headHash: previousHead[:],
			}, withAttrs)

			if err != nil {
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
			return pid, errors.New("invalid payload")
			// return pid, invalidBlock{error: ErrInvalidPayload,
			//root: arg.headRoot, invalidAncestorRoots: invalidRoots}

		default:
			// log.WithError(err).Error(ErrUndefinedExecutionEngineError)
			s.logger.Error("undefined execution engine error", "err", err)
			return nil, errors.New("undefined execution engine error")
		}
	}
	// forkchoiceUpdatedValidNodeCount.Inc()
	// if err := s.cfg.ForkChoiceStore.SetOptimisticToValid(ctx, arg.headRoot); err != nil {
	// 	log.WithError(err).Error("Could not set head root to valid")
	// 	return nil, nil
	// }
	// If the forkchoice update call has an attribute, update the proposer payload ID cache.
	// if hasAttr && payloadID != nil {
	// 	var pId [8]byte
	// 	copy(pId[:], payloadID[:])
	// 	log.WithFields(logrus.Fields{
	// 		"blockRoot": fmt.Sprintf("%#x", bytesutil.Trunc(arg.headRoot[:])),
	// 		"headSlot":  headBlk.Slot(),
	// 		"payloadID": fmt.Sprintf("%#x", bytesutil.Trunc(payloadID[:])),
	// 	}).Info("Forkchoice updated with payload attributes for proposal")
	// 	s.cfg.ProposerSlotIndexCache.SetProposerAndPayloadIDs(nextSlot, proposerId, pId, arg.headRoot)
	// } else if hasAttr && payloadID == nil && !features.Get().PrepareAllPayloads {
	// 	log.WithFields(logrus.Fields{
	// 		"blockHash": fmt.Sprintf("%#x", headPayload.BlockHash()),
	// 		"slot":      headBlk.Slot(),
	// 	}).Error("Received nil payload ID on VALID engine response")
	// }
	return payloadID, nil
}

// It returns true if the EL has returned VALID for the block.
func (s *Service) notifyNewPayload(ctx context.Context /*preStateVersion*/, _ int,
	preStateHeader interfaces.ExecutionData, /*, blk interfaces.ReadOnlySignedBeaconBlock*/
) (bool, error) {
	lastValidHash, err := s.engine.NewPayload(ctx, preStateHeader,
		[]common.Hash{}, &common.Hash{} /*empty version hashes and root before Deneb*/)
	return lastValidHash != nil, err
}

// validateExecutionOnBlock notifies the engine of the incoming block execution payload and
// returns true if the payload is valid.
func (s *Service) validateExecutionOnBlock(
	ctx context.Context, ver int, header interfaces.ExecutionData,
	/*,signed interfaces.ReadOnlySignedBeaconBlock, blockRoot [32]byte*/) (bool, error) {
	isValidPayload, err := s.notifyNewPayload(ctx, ver, header /*, signed*/)
	if err != nil {
		return false, err
		// return false, s.handleInvalidExecutionError(ctx, err, blockRoot,
		// signed.Block().ParentRoot())
	}
	// if signed.Version() < version.Capella && isValidPayload {
	// 	if err := s.validateMergeTransitionBlock(ctx, ver, header, signed); err != nil {
	// 		return isValidPayload, err
	// 	}
	// }
	return isValidPayload, nil
}

func (s *Service) getPayloadAttributes(_ context.Context,
	_ /*slot*/ primitives.Slot, timestamp uint64) (payloadattribute.Attributer, error) {
	// TODO: modularize andn make better.
	var random [32]byte
	if _, err := rand.Read(random[:]); err != nil {
		return nil, err
	}

	return payloadattribute.New(&enginev1.PayloadAttributesV2{
		Timestamp:             timestamp,
		SuggestedFeeRecipient: s.etherbase.Bytes(),
		Withdrawals:           nil,
		PrevRandao:            append([]byte{}, random[:]...),
	})
}
