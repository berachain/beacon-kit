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

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
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
	_, _, GenesisT, _,
]) InitGenesis(
	ctx context.Context,
	bz []byte,
) (transition.ValidatorUpdates, error) {
	var (
		err             error
		waitCtx, cancel = context.WithTimeout(ctx, AwaitTimeout)
	)
	defer cancel()
	// TODO: in theory the GenesisDataReceived channel should be empty, but we
	// should clear it anyways here to ensure that data is valid.

	data := new(GenesisT)
	if err = json.Unmarshal(bz, data); err != nil {
		h.logger.Error("Failed to unmarshal genesis data", "error", err)
		return nil, err
	}

	if err = h.dispatcher.Publish(
		async.NewEvent(ctx, async.GenesisDataReceived, *data),
	); err != nil {
		return nil, err
	}
	return h.waitForGenesisProcessed(waitCtx)
}

// waitForGenesisProcessed waits until the genesis data has been processed and
// returns the validator updates, or err if the context is cancelled.
func (h *ABCIMiddleware[
	_, _, _, _,
]) waitForGenesisProcessed(
	ctx context.Context,
) (transition.ValidatorUpdates, error) {
	select {
	case <-ctx.Done():
		return nil, ErrInitGenesisTimeout(ctx.Err())
	case gdpEvent := <-h.subGenDataProcessed:
		return gdpEvent.Data(), gdpEvent.Error()
	}
}

/* -------------------------------------------------------------------------- */
/*                               PrepareProposal                              */
/* -------------------------------------------------------------------------- */

// prepareProposal is the internal handler for preparing proposals.
func (h *ABCIMiddleware[
	BeaconBlockT, BlobSidecarsT, _, SlotDataT,
]) PrepareProposal(
	ctx context.Context,
	slotData SlotDataT,
) ([]byte, []byte, error) {
	var (
		err              error
		builtBeaconBlock BeaconBlockT
		builtSidecars    BlobSidecarsT
		startTime        = time.Now()
		awaitCtx, cancel = context.WithTimeout(ctx, AwaitTimeout)
	)
	// TODO: clear the built beacon block and sidecars channels, else we may
	// end up handling old data from previous slots.
	defer cancel()
	defer h.metrics.measurePrepareProposalDuration(startTime)

	if err = h.dispatcher.Publish(
		async.NewEvent(
			ctx, async.NewSlot, slotData,
		),
	); err != nil {
		return nil, nil, err
	}

	// wait for built beacon block
	builtBeaconBlock, err = h.waitForBuiltBeaconBlock(awaitCtx)
	if err != nil {
		return nil, nil, err
	}

	// wait for built sidecars
	builtSidecars, err = h.waitForBuiltSidecars(awaitCtx)
	if err != nil {
		return nil, nil, err
	}

	return h.handleBuiltBeaconBlockAndSidecars(builtBeaconBlock, builtSidecars)
}

// waitForBuiltBeaconBlock waits for the built beacon block to be received.
func (h *ABCIMiddleware[
	BeaconBlockT, BlobSidecarsT, _, SlotDataT,
]) waitForBuiltBeaconBlock(
	ctx context.Context,
) (BeaconBlockT, error) {
	select {
	case <-ctx.Done():
		return *new(BeaconBlockT), ErrBuildBeaconBlockTimeout(ctx.Err())
	case bbEvent := <-h.subBuiltBeaconBlock:
		return bbEvent.Data(), bbEvent.Error()
	}
}

// waitForBuiltSidecars waits for the built sidecars to be received.
func (h *ABCIMiddleware[
	_, BlobSidecarsT, _, _,
]) waitForBuiltSidecars(
	ctx context.Context,
) (BlobSidecarsT, error) {
	select {
	case <-ctx.Done():
		return *new(BlobSidecarsT), ErrBuildSidecarsTimeout(ctx.Err())
	case scEvent := <-h.subBuiltSidecars:
		return scEvent.Data(), scEvent.Error()
	}
}

// handleBuiltBeaconBlockAndSidecars gossips the built beacon block and blob
// sidecars to the network.
func (h *ABCIMiddleware[
	BeaconBlockT, BlobSidecarsT, _, _,
]) handleBuiltBeaconBlockAndSidecars(
	bb BeaconBlockT,
	sc BlobSidecarsT,
) ([]byte, []byte, error) {
	bbBz, bbErr := bb.MarshalSSZ()
	if bbErr != nil {
		return nil, nil, bbErr
	}
	scBz, scErr := sc.MarshalSSZ()
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
	BeaconBlockT, BlobSidecarsT, _, _,
]) ProcessProposal(
	ctx context.Context,
	req proto.Message,
) (proto.Message, error) {
	var (
		err              error
		startTime        = time.Now()
		blk              BeaconBlockT
		sidecars         BlobSidecarsT
		awaitCtx, cancel = context.WithTimeout(ctx, AwaitTimeout)
	)
	defer cancel()
	abciReq, ok := req.(*cmtabci.ProcessProposalRequest)
	if !ok {
		return nil, ErrInvalidProcessProposalRequestType
	}

	defer h.metrics.measureProcessProposalDuration(startTime)

	// Request the beacon block.
	if blk, err = encoding.
		UnmarshalBeaconBlockFromABCIRequest[BeaconBlockT](
		abciReq, 0, h.chainSpec.ActiveForkVersionForSlot(math.U64(abciReq.Height)),
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// notify that the beacon block has been received.
	if err = h.dispatcher.Publish(
		async.NewEvent(ctx, async.BeaconBlockReceived, blk),
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// Request the blob sidecars.
	if sidecars, err = encoding.
		UnmarshalBlobSidecarsFromABCIRequest[BlobSidecarsT](
		abciReq, 1,
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// notify that the sidecars have been received.
	if err = h.dispatcher.Publish(
		async.NewEvent(ctx, async.SidecarsReceived, sidecars),
	); err != nil {
		return h.createProcessProposalResponse(errors.WrapNonFatal(err))
	}

	// err if the built beacon block or sidecars failed verification.
	_, err = h.waitForBeaconBlockVerification(awaitCtx)
	if err != nil {
		return h.createProcessProposalResponse(err)
	}
	_, err = h.waitForSidecarVerification(awaitCtx)
	if err != nil {
		return h.createProcessProposalResponse(err)
	}
	return h.createProcessProposalResponse(nil)
}

// waitForBeaconBlockVerification waits for the built beacon block to be
// verified.
func (h *ABCIMiddleware[
	BeaconBlockT, _, _, _,
]) waitForBeaconBlockVerification(
	ctx context.Context,
) (BeaconBlockT, error) {
	select {
	case <-ctx.Done():
		return *new(BeaconBlockT), ErrVerifyBeaconBlockTimeout(ctx.Err())
	case vEvent := <-h.subBBVerified:
		return vEvent.Data(), vEvent.Error()
	}
}

// waitForSidecarVerification waits for the built sidecars to be verified.
func (h *ABCIMiddleware[
	_, BlobSidecarsT, _, _,
]) waitForSidecarVerification(
	ctx context.Context,
) (BlobSidecarsT, error) {
	select {
	case <-ctx.Done():
		return *new(BlobSidecarsT), ErrVerifySidecarsTimeout(ctx.Err())
	case vEvent := <-h.subSCVerified:
		return vEvent.Data(), vEvent.Error()
	}
}

// createResponse generates the appropriate ProcessProposalResponse based on the
// error.
func (*ABCIMiddleware[
	BeaconBlockT, _, BlobSidecarsT, _,
]) createProcessProposalResponse(
	err error,
) (proto.Message, error) {
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
	BeaconBlockT, BlobSidecarsT, _, _,
]) FinalizeBlock(
	ctx context.Context, req proto.Message,
) (transition.ValidatorUpdates, error) {
	var (
		err              error
		blk              BeaconBlockT
		blobs            BlobSidecarsT
		awaitCtx, cancel = context.WithTimeout(ctx, AwaitTimeout)
	)
	defer cancel()
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
		return nil, nil
	}

	// notify that the final beacon block has been received.
	if err = h.dispatcher.Publish(
		async.NewEvent(ctx, async.FinalBeaconBlockReceived, blk),
	); err != nil {
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
	_, _, _, _,
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
