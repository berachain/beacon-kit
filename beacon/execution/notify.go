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

	"github.com/berachain/beacon-kit/engine/client"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// notifyNewPayload notifies the execution client of a new payload.
func (s *Service) notifyNewPayload(
	ctx context.Context,
	slot primitives.Slot,
	payload enginetypes.ExecutionPayload,
	versionedHashes []primitives.ExecutionHash,
	parentBlockRoot [32]byte,
) (bool, error) {
	s.Logger().Info("notifying new payload",
		"payload_block_hash", (payload.GetBlockHash()),
		"parent_hash", (payload.GetParentHash()),
		"for_slot", slot,
	)

	lastValidHash, err := s.engine.NewPayload(
		ctx, payload, versionedHashes, &parentBlockRoot,
	)
	switch {
	case errors.Is(err, client.ErrAcceptedSyncingPayloadStatus):
		s.Logger().Info("new payload called with optimistic block",
			"payload_block_hash", (payload.GetBlockHash()),
			"parent_hash", (payload.GetParentHash()),
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
) (*enginetypes.PayloadID, error) {
	forkChoicer := s.ForkchoiceStore(ctx)

	fcs := &enginetypes.ForkchoiceState{
		HeadBlockHash:      fcuConfig.HeadEth1Hash,
		SafeBlockHash:      forkChoicer.JustifiedPayloadBlockHash(),
		FinalizedBlockHash: forkChoicer.FinalizedPayloadBlockHash(),
	}

	s.Logger().Info("notifying forkchoice update",
		"head_eth1_hash", fcuConfig.HeadEth1Hash,
		"safe_eth1_hash", forkChoicer.JustifiedPayloadBlockHash(),
		"finalized_eth1_hash", forkChoicer.FinalizedPayloadBlockHash(),
		"for_slot", fcuConfig.ProposingSlot,
		"has_attributes", fcuConfig.Attributes != nil,
	)

	// Notify the execution engine of the forkchoice update.
	payloadID, _, err := s.engine.ForkchoiceUpdated(
		ctx,
		fcs,
		fcuConfig.Attributes,
		s.ActiveForkVersionForSlot(fcuConfig.ProposingSlot),
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
			HeadEth1Hash:  forkChoicer.JustifiedPayloadBlockHash(),
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
