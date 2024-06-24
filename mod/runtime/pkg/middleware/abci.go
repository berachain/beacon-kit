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
	"fmt"
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
	_, _, _, _, _, _, GenesisT,
]) InitGenesis(
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

// prepareProposal is the internal handler for preparing proposals.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) PrepareProposal(
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
	if err := h.slotFeed.Publish(asynctypes.NewEvent(
		ctx, events.NewSlot, math.Slot(req.Height),
	)); err != nil {
		return nil, err
	}

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
]) waitForSidecars(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-h.sidecarsCh:
		if msg.Error() != nil {
			return nil, msg.Error()
		}
		return h.blobGossiper.Publish(ctx, msg.Data())
	}
}

// waitforBeaconBlk waits for the beacon block to be built and returns it.
func (h *ABCIMiddleware[
	_, _, _, _, _, _, _,
]) waitforBeaconBlk(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case beaconBlock := <-h.blkCh:
		if beaconBlock.Error() != nil {
			return nil, beaconBlock.Error()
		}
		return h.beaconBlockGossiper.Publish(
			ctx,
			beaconBlock.Data(),
		)
	}
}

/* -------------------------------------------------------------------------- */
/*                               ProcessProposal                              */
/* -------------------------------------------------------------------------- */

// ProcessProposal is a wrapper around the process proposal handler
// that extracts the beacon block from the proposal and processes it.
func (h *ABCIMiddleware[
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _,
]) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	var (
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
		err       error
		g, _      = errgroup.WithContext(ctx)
		startTime = time.Now()
	)
	defer h.metrics.measureProcessProposalDuration(startTime)

	// TODO: Consider exiting early if this node is not a validator to
	// reduce resource usage for full nodes.

	// Decode the beacon block and emit an event.
	blk, err = h.beaconBlockGossiper.Request(ctx, req)
	if err != nil {
		h.logger.Debug("failed to get beacon block", "error", err)
	}

	g.Go(func() error {
		// Emit event to notify the block has been received.
		localErr := h.blkBroker.Publish(asynctypes.NewEvent(
			ctx, events.BeaconBlockReceived, blk, err,
		))
		if localErr != nil {
			return localErr
		}

		if localErr = h.chainService.ReceiveBlock(
			ctx, blk,
		); !errors.IsFatal(localErr) {
			localErr = nil
		}
		return localErr
	})

	g.Go(func() error {
		// We can't notify the sidecars if the block is nil, since
		// we currently rely on the slot from the beacon block.
		if blk.IsNil() {
			return nil
		}

		// Decode the blob sidecars and emit an event.
		var localErr error
		sidecars, localErr = h.blobGossiper.Request(ctx, req)
		if localErr != nil {
			h.logger.Debug("failed to get sidecars", "error", localErr)
		}

		// Emit event to notify the sidecars have been received.
		if localErr = h.sidecarsBroker.Publish(asynctypes.NewEvent(
			ctx, events.BlobSidecarsReceived, sidecars, localErr,
		)); localErr != nil {
			return localErr
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-h.sidecarsCh:
			if msg.Type() != events.BlobSidecarsProcessed {
				return fmt.Errorf(
					"unexpected event type: %s", msg.Type(),
				)
			}
			if msg.Error() != nil {
				return msg.Error()
			}
			sidecars = msg.Data()
		}
		return nil
	})

	resp := &cmtabci.ProcessProposalResponse{
		Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
	}

	// If we see a non fatal error, clear everything.
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
	_, BeaconBlockT, _, BlobSidecarsT, _, _, _,
]) EndBlock(
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

	// Send the sidecars to the sidecars feed, we know at this point
	// That the blobs have been successfully verified in process proposal.
	if err = h.sidecarsBroker.Publish(asynctypes.NewEvent(
		ctx, events.BlobSidecarsVerified, blobs,
	)); err != nil {
		return nil, err
	}

	// Wait for a response from the da service, with the current codepaths
	// we can't parallelize retrieving the DA service response and the
	// validator updates, since we need to check for IsDataAvailable in
	// `ProcessBeaconBlock`, we should improve this though.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case sidecars := <-h.sidecarsCh:
		if sidecars.Type() != events.BlobSidecarsProcessed {
			return nil, fmt.Errorf(
				"unexpected event type: %s", sidecars.Type())
		}
		if sidecars.Error() != nil {
			return nil, sidecars.Error()
		}
	}

	// TODO: Move to Async.
	valUpdates, err := h.chainService.ProcessBeaconBlock(
		ctx, blk,
	)
	if err != nil {
		return nil, err
	}

	return iter.MapErr(
		valUpdates.RemoveDuplicates().Sort(), convertValidatorUpdate,
	)
}
