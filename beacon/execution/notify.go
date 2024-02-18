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
	eth "github.com/itsdevbear/bolaris/execution/engine/ethclient"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
	"github.com/itsdevbear/bolaris/types/consensus/version"
	"github.com/itsdevbear/bolaris/types/engine"
	enginev1 "github.com/itsdevbear/bolaris/types/engine/v1"
)

// notifyNewPayload notifies the execution client of a new payload.
func (s *Service) notifyNewPayload(
	ctx context.Context, payload engine.ExecutionPayload,
) (bool, error) {
	var (
		lastValidHash []byte
		err           error
		beaconState   = s.BeaconState(ctx)
	)

	//nolint:revive // okay for now.
	if beaconState.Version() >= version.Deneb {
		// TODO: Deneb
		// var versionedHashes []common.Hash
		// versionedHashes, err = kzgCommitmentsToVersionedHashes(blk.Block().Body())
		// if err != nil {
		// 	return false, errors.Wrap(err, "could not get versioned hashes to feed the engine")
		// }
		// pr := common.Hash(blk.Block().ParentRoot())
		// lastValidHash, err = s.engine.NewPayload(ctx, payload, versionedHashes, &pr)
	} else {
		lastValidHash, err = s.engine.NewPayload(
			ctx, payload, []common.Hash{}, &common.Hash{}, /*empty version hashes and root before Deneb*/
		)
	}
	switch {
	case err == nil:
		return true, nil
	case errors.Is(err, eth.ErrAcceptedSyncingPayloadStatus):
		s.Logger().Info("new payload called with optimistic block",
			"head_eth1_hash", common.Bytes2Hex(payload.GetBlockHash()),
			"proposing_slot", beaconState.Slot,
		)
		return false, nil
	case errors.Is(err, eth.ErrInvalidPayloadStatus):
		s.Logger().Error("invalid payload status", "last_valid_hash", fmt.Sprintf("%#x", lastValidHash))
		err = ErrBadBlockProduced
	}

	return false, err
}
func (s *Service) notifyForkchoiceUpdate(
	ctx context.Context, fcuConfig *FCUConfig,
) (*enginev1.PayloadIDBytes, error) {
	beaconState := s.BeaconState(ctx)

	// Ensure we don't pass a nil attribute to the execution engine.
	if fcuConfig.Attributes == nil {
		fcuConfig.Attributes = engine.EmptyPayloadAttributesWithVersion(beaconState.Version())
	}

	// Notify the execution engine of the forkchoice update.
	payloadID, _, err := s.engine.ForkchoiceUpdated(
		ctx,
		&enginev1.ForkchoiceState{
			HeadBlockHash:      fcuConfig.HeadEth1Hash[:],
			SafeBlockHash:      beaconState.GetSafeEth1BlockHash().Bytes(),
			FinalizedBlockHash: beaconState.GetFinalizedEth1BlockHash().Bytes(),
		},
		fcuConfig.Attributes,
	)
	switch {
	case errors.Is(err, eth.ErrAcceptedSyncingPayloadStatus):
		s.Logger().Info("forkchoice updated with optimistic block",
			"head_eth1_hash", fcuConfig.HeadEth1Hash,
			"proposing_slot", fcuConfig.ProposingSlot,
		)
		telemetry.IncrCounter(1, MetricsKeyAcceptedSyncingPayloadStatus)
		return payloadID, nil
	case errors.Is(err, eth.ErrInvalidPayloadStatus):
		s.Logger().Error("invalid payload status", "error", err)
		telemetry.IncrCounter(1, MetricsKeyInvalidPayloadStatus)
		// Attempt to get the chain back into a valid state.
		payloadID, err = s.notifyForkchoiceUpdate(ctx, &FCUConfig{
			HeadEth1Hash:  beaconState.GetLastValidHead(),
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

	// We can mark this Eth1Block as the latest valid block.
	// TODO: maybe move to blockchain for IsCanonical and Head checks.
	// TODO: the whole getting the execution payload off the block /
	// the whole LastestExecutionPayload Premine thing "PremineGenesisConfig".
	beaconState.SetLastValidHead(fcuConfig.HeadEth1Hash)

	// If the forkchoice update call has an attribute, update the payload ID cache.
	hasAttr := fcuConfig.Attributes != nil && !fcuConfig.Attributes.IsEmpty()
	if hasAttr && payloadID != nil {
		s.Logger().Info("forkchoice updated with payload attributes for proposal",
			"head_eth1_hash", fcuConfig.HeadEth1Hash,
			"proposing_slot", fcuConfig.ProposingSlot,
			"payload_id", fmt.Sprintf("%#x", payloadID),
		)
		s.payloadCache.Set(
			fcuConfig.ProposingSlot, fcuConfig.HeadEth1Hash, primitives.PayloadID(payloadID[:]))
	} else if hasAttr && payloadID == nil {
		/*TODO: introduce this feature && !s.cfg.Features.Get().PrepareAllPayloads*/
		s.Logger().Error("received nil payload ID on VALID engine response",
			"head_eth1_hash", fmt.Sprintf("%#x", fcuConfig.HeadEth1Hash),
			"proposing_slot", fcuConfig.ProposingSlot,
		)
	}

	return payloadID, nil
}
