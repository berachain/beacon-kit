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
	"fmt"

	"github.com/prysmaticlabs/prysm/v4/beacon-chain/execution"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v4/encoding/bytesutil"
	enginev1 "github.com/prysmaticlabs/prysm/v4/proto/engine/v1"
)

func (s *Service) notifyForkchoiceUpdate(
	ctx context.Context, fcuConfig *FCUConfig,
) error {
	var (
		payloadID   *primitives.PayloadID
		err         error
		beaconState = s.bsp.BeaconState(ctx)
		fc          = &enginev1.ForkchoiceState{
			HeadBlockHash:      fcuConfig.HeadEth1Hash[:],
			SafeBlockHash:      beaconState.GetSafeEth1BlockHash().Bytes(),
			FinalizedBlockHash: beaconState.GetFinalizedEth1BlockHash().Bytes(),
		}
	)

	// TODO: remember and figure out what the middle param is.
	payloadID, _, err = s.engine.ForkchoiceUpdated(ctx, fc, fcuConfig.Attributes)
	if err != nil {
		switch err { //nolint:errorlint // okay for now.
		case execution.ErrAcceptedSyncingPayloadStatus:
			return nil
		case execution.ErrInvalidPayloadStatus:
			s.Logger().Error("invalid payload status", "error", err)
			// In Prysm, this code recursively calls back until we find a valid hash we can
			// insert. In BeaconKit, we don't have the nice ability to do this, *but* in
			// theory we should never need it, since we have single block finality thanks
			// to CometBFT. Essentially, if we get an invalid payload status here, something
			// higher up must've gone wrong and thus we don't really need the retry here.
			return errors.New("invalid payload")
		default:
			s.Logger().Error("undefined execution engine error", "error", err)
			return err
		}
	}

	// We can mark this Eth1Block as the latest valid block.
	// TODO: maybe move to blockchain for IsCanonical and Head checks.
	// TODO: the whole getting the execution payload off the block /
	// the whole LastestExecutionPayload Premine thing "PremineGenesisConfig".
	s.bsp.BeaconState(ctx).SetLastValidHead(fcuConfig.HeadEth1Hash)

	// If the forkchoice update call has an attribute, update the payload ID cache.
	hasAttr := fcuConfig.Attributes != nil && !fcuConfig.Attributes.IsEmpty()
	if hasAttr && payloadID != nil {
		var pID [8]byte
		copy(pID[:], payloadID[:])
		s.Logger().Info("forkchoice updated with payload attributes for proposal",
			"head_eth1_hash", fcuConfig.HeadEth1Hash,
			"proposing_slot", fcuConfig.ProposingSlot,
			"payloadID", fmt.Sprintf("%#x", bytesutil.Trunc(payloadID[:])),
		)
		s.payloadCache.Set(fcuConfig.ProposingSlot, fcuConfig.HeadEth1Hash, pID)
	} else if hasAttr && payloadID == nil {
		/*TODO: introduce this feature && !s.cfg.Features.Get().PrepareAllPayloads*/
		s.Logger().Error("received nil payload ID on VALID engine response",
			"head_eth1_hash", fmt.Sprintf("%#x", fcuConfig.HeadEth1Hash),
			"proposing_slot", fcuConfig.ProposingSlot,
		)
	}

	return nil
}
