// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package middleware

import (
	"context"

	"github.com/berachain/beacon-kit/consensus/cometbft/service/encoding"
	"github.com/berachain/beacon-kit/consensus/types"
	"github.com/berachain/beacon-kit/primitives/async"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/* -------------------------------------------------------------------------- */
/*                                FinalizeBlock                               */
/* -------------------------------------------------------------------------- */

// EndBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	BeaconBlockT, _, BlobSidecarsT, _, _,
]) FinalizeBlock(
	ctx sdk.Context,
	req *cmtabci.FinalizeBlockRequest,
) (transition.ValidatorUpdates, error) {
	awaitCtx, cancel := context.WithTimeout(ctx, AwaitTimeout)
	defer cancel()
	// flush the channel to ensure that we are not handling old data.
	if numMsgs := async.ClearChan(h.subFinalValidatorUpdates); numMsgs > 0 {
		h.logger.Error(
			"WARNING: messages remaining in final validator updates channel",
			"num_msgs", numMsgs)
	}

	blk, blobs, err := encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](
		req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		h.chainSpec.ActiveForkVersionForSlot(
			math.Slot(req.Height),
		))
	if err != nil {
		// If we don't have a block, we can't do anything.
		return nil, nil
	}

	// notify that the final beacon block has been received.
	var consensusBlk *types.ConsensusBlock[BeaconBlockT]
	consensusBlk = consensusBlk.New(
		blk,
		req.GetProposerAddress(),
		req.GetTime(),
	)
	blkEvent := async.NewEvent(
		ctx,
		async.FinalBeaconBlockReceived,
		consensusBlk,
	)
	if err = h.dispatcher.Publish(blkEvent); err != nil {
		return nil, err
	}

	// notify that the final blob sidecars have been received.
	if err = h.dispatcher.Publish(
		async.NewEvent(ctx, async.FinalSidecarsReceived, blobs),
	); err != nil {
		return nil, err
	}

	// wait for the final validator updates.
	return h.waitForFinalValidatorUpdates(awaitCtx)
}

// waitForFinalValidatorUpdates waits for the final validator updates to be
// received.
func (h *ABCIMiddleware[
	_, _, _, _, _,
]) waitForFinalValidatorUpdates(
	ctx context.Context,
) (transition.ValidatorUpdates, error) {
	select {
	case <-ctx.Done():
		return nil, ErrFinalValidatorUpdatesTimeout(ctx.Err())
	case event := <-h.subFinalValidatorUpdates:
		return event.Data(), event.Error()
	}
}
