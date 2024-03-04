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

package execution

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/engine/client"
	enginetypes "github.com/itsdevbear/bolaris/engine/types"
	enginev1 "github.com/itsdevbear/bolaris/engine/types/v1"
	"github.com/itsdevbear/bolaris/primitives"
)

// notifyNewPayload notifies the execution client of a new payload.
func (s *Service) notifyNewPayload(
	ctx context.Context,
	slot primitives.Slot,
	payload enginetypes.ExecutionPayload,
	versionedHashes []common.Hash,
	parentBlockRoot [32]byte,
) (bool, error) {
	s.Logger().Info("notifying new payload",
		"payload_block_hash", common.BytesToHash(payload.GetBlockHash()),
		"parent_hash", common.BytesToHash(payload.GetParentHash()),
		"for_slot", slot,
	)

	lastValidHash, err := s.engine.NewPayload(
		ctx, payload, versionedHashes, &parentBlockRoot,
	)
	switch {
	case errors.Is(err, client.ErrAcceptedSyncingPayloadStatus):
		s.Logger().Info("new payload called with optimistic block",
			"block_hash", common.BytesToHash(payload.GetBlockHash()),
			"parent_hash", common.BytesToHash(payload.GetParentHash()),
			"for_slot", slot,
		)
		return false, nil
	case errors.Is(err, client.ErrInvalidPayloadStatus):
		s.Logger().Error(
			"invalid payload status",
			"last_valid_hash", fmt.Sprintf("%#x", lastValidHash),
		)
		return false, ErrBadBlockProduced
	case err != nil:
		return false, err
	}
	return true, nil
}

// notifyForkchoiceUpdate notifies the execution client of a forkchoice update.
func (s *Service) notifyForkchoiceUpdate(
	ctx context.Context, fcuConfig *FCUConfig,
) (*enginev1.PayloadIDBytes, error) {
	forkChoicer := s.ForkchoiceStore(ctx)

	// TODO: intercept here and ask builder service for payload attributes.
	// if isValidator && PrepareAllPayloads {
	// Ensure we don't pass a nil attribute to the execution engine.
	if fcuConfig.Attributes == nil {
		fcuConfig.Attributes = enginetypes.EmptyPayloadAttributesWithVersion(
			s.ActiveForkVersionForSlot(fcuConfig.ProposingSlot))
	}

	fcs := &enginev1.ForkchoiceState{
		HeadBlockHash:      fcuConfig.HeadEth1Hash[:],
		SafeBlockHash:      forkChoicer.GetSafeEth1BlockHash().Bytes(),
		FinalizedBlockHash: forkChoicer.GetFinalizedEth1BlockHash().Bytes(),
	}

	s.Logger().Info("notifying forkchoice update",
		"head_eth1_hash", fcuConfig.HeadEth1Hash,
		"safe_eth1_hash", forkChoicer.GetSafeEth1BlockHash(),
		"finalized_eth1_hash", forkChoicer.GetFinalizedEth1BlockHash(),
		"for_slot", fcuConfig.ProposingSlot,
		"has_attributes", fcuConfig.Attributes != nil &&
			!fcuConfig.Attributes.IsEmpty(),
	)

	// Notify the execution engine of the forkchoice update.
	payloadID, _, err := s.engine.ForkchoiceUpdated(
		ctx,
		fcs,
		fcuConfig.Attributes,
	)
	switch {
	case errors.Is(err, client.ErrAcceptedSyncingPayloadStatus):
		s.Logger().Info("forkchoice updated with optimistic block",
			"head_eth1_hash", fcuConfig.HeadEth1Hash,
			"for_slot", fcuConfig.ProposingSlot,
		)
		telemetry.IncrCounter(1, MetricsKeyAcceptedSyncingPayloadStatus)
		return payloadID, nil
	case errors.Is(err, client.ErrInvalidPayloadStatus):
		s.Logger().Error("invalid payload status", "error", err)
		telemetry.IncrCounter(1, MetricsKeyInvalidPayloadStatus)
		// Attempt to get the chain back into a valid state, by
		// getting finding an ancestor block with a valid payload and
		// forcing a recovery.
		payloadID, err = s.notifyForkchoiceUpdate(ctx, &FCUConfig{
			// TODO: in the case of CometBFT BeaconKit, this could in theory
			// just be the last finalized block, bc we are always inserting
			// ontop of that, however making that assumption here feels
			// a little coupled.
			HeadEth1Hash:  forkChoicer.GetSafeEth1BlockHash(),
			ProposingSlot: fcuConfig.ProposingSlot,
			Attributes:    fcuConfig.Attributes,
		})
		if err != nil {
			// We have to return the error here since this function
			// is recursive.
			return nil, err
		}
		return payloadID, ErrBadBlockProduced
	case err != nil:
		s.Logger().Error("undefined execution engine error", "error", err)
		return nil, err
	}

	return payloadID, nil
}
