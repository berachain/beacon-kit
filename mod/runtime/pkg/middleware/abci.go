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
	"encoding/json"
	"sync"
	"time"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
	"golang.org/x/sync/errgroup"
)

/* -------------------------------------------------------------------------- */
/*                                 InitGenesis                                */
/* -------------------------------------------------------------------------- */

// InitGenesis is called by the base app to initialize the state of the.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) InitGenesis(
	ctx context.Context,
	bz []byte,
) ([]appmodulev2.ValidatorUpdate, error) {
	return h.initGenesis(ctx, bz)
}

// initGenesis is called by the base app to initialize the state of the.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, GenesisT,
]) initGenesis(
	ctx context.Context,
	bz []byte,
) ([]appmodulev2.ValidatorUpdate, error) {
	data := new(GenesisT)
	if err := json.Unmarshal(bz, data); err != nil {
		return nil, err
	}
	updates, err := h.chainService.ProcessGenesisData(
		ctx,
		*data,
	)
	if err != nil {
		return nil, err
	}

	// Convert updates into the Cosmos SDK format.
	return iter.MapErr(updates, convertValidatorUpdate)
}

/* -------------------------------------------------------------------------- */
/*                               PrepareProposal                              */
/* -------------------------------------------------------------------------- */

// PrepareProposal is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) PrepareProposal(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	return h.prepareProposal(ctx, req)
}

// prepareProposal is the internal handler for preparing proposals.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) prepareProposal(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	var (
		wg                          sync.WaitGroup
		startTime                   = time.Now()
		beaconBlockErr, sidecarsErr error
		beaconBlockBz, sidecarsBz   []byte
	)
	defer h.metrics.measurePrepareProposalDuration(startTime)

	// Send a request to the validator service to give us a beacon block
	// and blob sidecards to pass to ABCI.
	h.slotFeed.Send(asynctypes.NewEvent(
		ctx, events.NewSlot, math.Slot(req.Height),
	))

	// Using a wait group instead of an errgroup to ensure we drain
	// the associated channels for the beacon block and sidecars.
	//nolint:mnd // bet.
	wg.Add(2)
	go func() {
		defer wg.Done()
		beaconBlockBz, beaconBlockErr = h.waitforBeaconBlk(ctx)
	}()

	go func() {
		defer wg.Done()
		sidecarsBz, sidecarsErr = h.waitForSidecars(ctx)
	}()

	wg.Wait()
	if beaconBlockErr != nil {
		return nil, beaconBlockErr
	} else if sidecarsErr != nil {
		return nil, sidecarsErr
	}

	return &cmtabci.PrepareProposalResponse{
		Txs: [][]byte{beaconBlockBz, sidecarsBz},
	}, nil
}

// waitForSidecars waits for the sidecars to be built and returns them.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) waitForSidecars(gCtx context.Context) ([]byte, error) {
	select {
	case <-gCtx.Done():
		return nil, gCtx.Err()
	case err := <-h.errCh:
		return nil, err
	case sidecars := <-h.prepareProposalSidecarsCh:
		if sidecars.Error() != nil {
			return nil, sidecars.Error()
		}

		sidecarsBz, err := h.blobGossiper.Publish(gCtx, sidecars.Data())
		if err != nil {
			h.logger.Error("failed to publish blobs", "error", err)
		}
		return sidecarsBz, err
	}
}

// waitforBeaconBlk waits for the beacon block to be built and returns it.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) waitforBeaconBlk(gCtx context.Context) ([]byte, error) {
	select {
	case <-gCtx.Done():
		return nil, gCtx.Err()
	case err := <-h.errCh:
		return nil, err
	case beaconBlock := <-h.prepareProposalBlkCh:
		if beaconBlock.Error() != nil {
			return nil, beaconBlock.Error()
		}
		beaconBlockBz, err := h.beaconBlockGossiper.Publish(
			gCtx,
			beaconBlock.Data(),
		)
		if err != nil {
			h.logger.Error("failed to publish beacon block", "error", err)
		}
		return beaconBlockBz, err
	}
}

/* -------------------------------------------------------------------------- */
/*                               ProcessProposal                              */
/* -------------------------------------------------------------------------- */

// ProcessProposal is a wrapper around the process proposal handler
// that extracts the beacon block from the proposal and processes it.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	return h.processProposal(ctx, req)
}

// processProposal is the internal handler for processing proposals.
func (h *ABCIMiddleware[
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _,
]) processProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	var (
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
		err       error
		g, gCtx   = errgroup.WithContext(ctx)
		startTime = time.Now()
		args      = []any{"beacon_block", true, "blob_sidecars", true}
	)
	defer h.metrics.measureProcessProposalDuration(startTime)

	blk, err = h.beaconBlockGossiper.Request(gCtx, req)
	if err != nil {
		args[1] = false
	}

	g.Go(func() error {
		if err = h.chainService.ReceiveBlock(
			ctx, blk,
		); !errors.IsFatal(err) {
			err = nil
		}
		return err
	})

	g.Go(func() error {
		sidecars, err = h.blobGossiper.Request(gCtx, req)
		if err != nil {
			args[3] = false
		}

		if blk.IsNil() {
			return nil
		}

		if err = h.daService.ReceiveSidecars(
			gCtx, blk.GetSlot(), sidecars,
		); !errors.IsFatal(err) {
			err = nil
		}
		return err
	})

	resp := &cmtabci.ProcessProposalResponse{
		Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
	}

	// If we see a non fatal error, clear everything.
	defer h.logger.Info("processed proposal", args...)
	if err = g.Wait(); !errors.IsFatal(err) {
		resp.Status = cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT
		err = nil
	}
	return resp, err
}

/* -------------------------------------------------------------------------- */
/*                                FinalizeBlock                               */
/* -------------------------------------------------------------------------- */

// PreBlock is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) PreBlock(
	_ sdk.Context, req *cmtabci.FinalizeBlockRequest,
) error {
	h.req = req
	return nil
}

// EndBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	return h.endBlock(ctx)
}

// endBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _,
]) endBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	blk, blobs, err := encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](
		h.req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		h.chainSpec.ActiveForkVersionForSlot(
			math.Slot(h.req.Height),
		))
	if err != nil {
		// If we don't have a block, we can't do anything.
		//nolint:nilerr // by design.
		return nil, nil
	}

	// TODO: Move to Async.
	if err = h.daService.ProcessSidecars(
		ctx, blk.GetSlot(), blobs,
	); err != nil {
		return nil, err
	}

	// TODO: Move to Async.
	valUpdates, err := h.chainService.ProcessBeaconBlock(ctx, blk)
	if err != nil {
		return nil, err
	}

	return iter.MapErr(
		valUpdates.RemoveDuplicates().Sort(), convertValidatorUpdate,
	)
}
