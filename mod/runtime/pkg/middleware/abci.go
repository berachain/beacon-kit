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
	"time"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/messages"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
)

/* -------------------------------------------------------------------------- */
/*                                 InitGenesis                                */
/* -------------------------------------------------------------------------- */

// InitGenesis is called by the base app to initialize the state of the.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, GenesisT, _,
]) InitGenesis(
	ctx context.Context,
	bz []byte,
) (transition.ValidatorUpdates, error) {
	var err error
	data := new(GenesisT)
	if err = json.Unmarshal(bz, data); err != nil {
		h.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	// request for validator updates
	valUpdates := asynctypes.NewFuture[transition.ValidatorUpdates]()
	if err = h.dispatcher.SendRequest(
		asynctypes.NewMessage(
			ctx, messages.ProcessGenesisData, *data,
		), valUpdates,
	); err != nil {
		return nil, err
	}

	return valUpdates.Resolve()
}

/* -------------------------------------------------------------------------- */
/*                               PrepareProposal                              */
/* -------------------------------------------------------------------------- */

// prepareProposal is the internal handler for preparing proposals.
func (h *ABCIMiddleware[
	_, BeaconBlockT, BeaconBlockBundleT, BlobSidecarsT, _, _, _, SlotDataT,
]) PrepareProposal(
	ctx context.Context,
	slotData SlotDataT,
) ([]byte, []byte, error) {
	var err error
	startTime := time.Now()
	defer h.metrics.measurePrepareProposalDuration(startTime)

	// request a built beacon block for the given slot
	beaconBlkBundleFuture := asynctypes.NewFuture[BeaconBlockBundleT]()
	if err = h.dispatcher.SendRequest(
		asynctypes.NewMessage(
			ctx, messages.BuildBeaconBlockAndSidecars, slotData,
		), beaconBlkBundleFuture,
	); err != nil {
		return nil, nil, err
	}

	// resolve the beacon block bundle from the future
	beaconBlkBundle, err := beaconBlkBundleFuture.Resolve()
	if err != nil {
		return nil, nil, err
	}

	// gossip the built beacon block and blob sidecars
	return h.handleBeaconBlockBundleResponse(ctx, beaconBlkBundle)
}

// handleBeaconBlockBundleResponse gossips the built beacon block and blob
// sidecars to the network.
func (h *ABCIMiddleware[
	_, BeaconBlockT, BeaconBlockBundleT, BlobSidecarsT, _, _, _, _,
]) handleBeaconBlockBundleResponse(
	ctx context.Context,
	bbb BeaconBlockBundleT,
) ([]byte, []byte, error) {
	// gossip beacon block
	bbBz, bbErr := h.beaconBlockGossiper.Publish(
		ctx, bbb.GetBeaconBlock(),
	)
	if bbErr != nil {
		return nil, nil, bbErr
	}
	// gossip blob sidecars
	scBz, scErr := h.blobGossiper.Publish(ctx, bbb.GetSidecars())
	if scErr != nil {
		return nil, nil, scErr
	}
	return bbBz, scBz, nil
}

/* -------------------------------------------------------------------------- */
/*                               ProcessProposal                              */
/* -------------------------------------------------------------------------- */

// ProcessProposal processes the proposal for the ABCI middleware.
// It handles both the beacon block and blob sidecars concurrently.
func (h *ABCIMiddleware[
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _, _,
]) ProcessProposal(
	ctx context.Context,
	req proto.Message,
) (proto.Message, error) {
	var (
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
		err       error
		startTime = time.Now()
	)
	abciReq, ok := req.(*cmtabci.ProcessProposalRequest)
	if !ok {
		return nil, ErrInvalidProcessProposalRequestType
	}

	defer h.metrics.measureProcessProposalDuration(startTime)

	// Request the beacon block.
	if blk, err = h.beaconBlockGossiper.Request(ctx, abciReq); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// verify the beacon block
	beaconBlockFuture := asynctypes.NewFuture[BeaconBlockT]()
	if err = h.dispatcher.SendRequest(
		asynctypes.NewMessage(
			ctx, messages.VerifyBeaconBlock, blk,
		), beaconBlockFuture,
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// Request the blob sidecars.
	if sidecars, err = h.blobGossiper.Request(ctx, abciReq); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// verify the blob sidecars
	sidecarsFuture := asynctypes.NewFuture[BlobSidecarsT]()
	if err = h.dispatcher.SendRequest(
		asynctypes.NewMessage(
			ctx, messages.VerifySidecars, sidecars,
		), sidecarsFuture,
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// error if the beacon block or sidecars are invalid
	_, err = beaconBlockFuture.Resolve()
	if err != nil {
		return h.createProcessProposalResponse(err)
	}

	_, err = sidecarsFuture.Resolve()
	if err != nil {
		return h.createProcessProposalResponse(err)
	}

	return h.createProcessProposalResponse(nil)
}

// createResponse generates the appropriate ProcessProposalResponse based on the
// error.
func (*ABCIMiddleware[
	_, BeaconBlockT, _, _, BlobSidecarsT, _, _, _,
]) createProcessProposalResponse(err error) (proto.Message, error) {
	status := cmtabci.PROCESS_PROPOSAL_STATUS_REJECT
	if !errors.IsFatal(err) {
		status = cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT
		err = nil
	}
	return &cmtabci.ProcessProposalResponse{Status: status}, err
}

/* -------------------------------------------------------------------------- */
/*                                FinalizeBlock                               */
/* -------------------------------------------------------------------------- */

// EndBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _, _,
]) FinalizeBlock(
	ctx context.Context, req proto.Message,
) (transition.ValidatorUpdates, error) {
	var (
		blk   BeaconBlockT
		blobs BlobSidecarsT
		err   error
	)
	abciReq, ok := req.(*cmtabci.FinalizeBlockRequest)
	if !ok {
		return nil, ErrInvalidFinalizeBlockRequestType
	}
	blk, blobs, err = encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](
		abciReq,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		h.chainSpec.ActiveForkVersionForSlot(
			math.Slot(abciReq.Height),
		))
	if err != nil {
		// If we don't have a block, we can't do anything.
		//nolint:nilerr // by design.
		return nil, nil
	}

	// process the sidecars.
	sidecarsFuture := asynctypes.NewFuture[BlobSidecarsT]()
	if err = h.dispatcher.SendRequest(
		asynctypes.NewMessage(
			ctx, messages.ProcessSidecars, blobs,
		), sidecarsFuture,
	); err != nil {
		return nil, err
	}

	// finalize the beacon block.
	valUpdatesFuture := asynctypes.NewFuture[transition.ValidatorUpdates]()
	if err = h.dispatcher.SendRequest(
		asynctypes.NewMessage(
			ctx, messages.FinalizeBeaconBlock, blk,
		), valUpdatesFuture,
	); err != nil {
		return nil, err
	}

	return valUpdatesFuture.Resolve()
}
