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
	"fmt"
	"time"

	async "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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
	var (
		err      error
		gdpEvent async.Event[transition.ValidatorUpdates]
	)
	// in theory this channel should already be empty, but we clear it anyways
	h.subGenDataProcessed.Clear()

	data := new(GenesisT)
	if err = json.Unmarshal(bz, data); err != nil {
		h.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	if err = h.dispatcher.PublishEvent(
		async.NewEvent(ctx, events.GenesisDataReceived, *data),
	); err != nil {
		return nil, err
	}

	gdpEvent, err = h.subGenDataProcessed.Await(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("GOT IT BACK")
	return gdpEvent.Data(), gdpEvent.Error()
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
	var (
		err          error
		builtBBEvent async.Event[BeaconBlockT]
		builtSCEvent async.Event[BlobSidecarsT]
		startTime    = time.Now()
	)
	defer h.metrics.measurePrepareProposalDuration(startTime)
	// in theory these subs should already be empty, but we clear them anyways
	h.subBuiltBeaconBlock.Clear()
	h.subBuiltSidecars.Clear()

	if err = h.dispatcher.PublishEvent(
		async.NewEvent(
			ctx, events.NewSlot, slotData,
		),
	); err != nil {
		return nil, nil, err
	}

	// wait for built beacon block
	builtBBEvent, err = h.subBuiltBeaconBlock.Await(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err = builtBBEvent.Error(); err != nil {
		return nil, nil, err
	}

	// wait for built sidecars
	builtSCEvent, err = h.subBuiltSidecars.Await(ctx)
	if err != nil {
		return nil, nil, err
	}
	if err = builtSCEvent.Error(); err != nil {
		return nil, nil, err
	}

	// gossip the built beacon block and blob sidecars
	return h.handleBuiltBeaconBlockAndSidecars(
		ctx, builtBBEvent.Data(), builtSCEvent.Data(),
	)
}

// handleBeaconBlockBundleResponse gossips the built beacon block and blob
// sidecars to the network.
func (h *ABCIMiddleware[
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _, _,
]) handleBuiltBeaconBlockAndSidecars(
	ctx context.Context,
	bb BeaconBlockT,
	sc BlobSidecarsT,
) ([]byte, []byte, error) {
	// gossip beacon block
	bbBz, bbErr := h.beaconBlockGossiper.Publish(
		ctx, bb,
	)
	if bbErr != nil {
		return nil, nil, bbErr
	}
	// gossip blob sidecars
	scBz, scErr := h.blobGossiper.Publish(ctx, sc)
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
		err       error
		startTime = time.Now()
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
	)
	// in theory these subs should already be empty, but we clear them anyways
	h.subBBVerified.Clear()
	h.subSCVerified.Clear()
	abciReq, ok := req.(*cmtabci.ProcessProposalRequest)
	if !ok {
		return nil, ErrInvalidProcessProposalRequestType
	}

	defer h.metrics.measureProcessProposalDuration(startTime)

	// Request the beacon block.
	if blk, err = h.beaconBlockGossiper.Request(ctx, abciReq); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// TODO: implement service
	// notify that the beacon block has been received.
	if err = h.dispatcher.PublishEvent(
		async.NewEvent(ctx, events.BeaconBlockReceived, blk),
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// Request the blob sidecars.
	if sidecars, err = h.blobGossiper.Request(ctx, abciReq); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// notify that the sidecars have been received.
	if err = h.dispatcher.PublishEvent(
		async.NewEvent(ctx, events.SidecarsReceived, sidecars),
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// err if the built beacon block or sidecars failed verification.
	_, err = h.subBBVerified.Await(ctx)
	if err != nil {
		return h.createProcessProposalResponse(err)
	}
	_, err = h.subSCVerified.Await(ctx)
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
		err                  error
		blk                  BeaconBlockT
		blobs                BlobSidecarsT
		finalValUpdatesEvent async.Event[transition.ValidatorUpdates]
	)
	// in theory this sub should already be empty, but we clear them anyways
	h.subFinalValidatorUpdates.Clear()
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

	// notify that the final beacon block has been received.
	if err = h.dispatcher.PublishEvent(
		async.NewEvent(ctx, events.FinalBeaconBlockReceived, blk),
	); err != nil {
		return nil, err
	}

	// notify that the final blob sidecars have been received.
	if err = h.dispatcher.PublishEvent(
		async.NewEvent(ctx, events.FinalSidecarsReceived, blobs),
	); err != nil {
		return nil, err
	}

	// wait for the final validator updates.
	finalValUpdatesEvent, err = h.subFinalValidatorUpdates.Await(ctx)
	if err != nil {
		return nil, err
	}

	return finalValUpdatesEvent.Data(), finalValUpdatesEvent.Error()
}
