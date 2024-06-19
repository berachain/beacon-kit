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
)

/* -------------------------------------------------------------------------- */
/*                                 InitGenesis                                */
/* -------------------------------------------------------------------------- */

// InitGenesis is called by the base app to initialize the state of the.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) InitGenesis(
	ctx context.Context,
	bz []byte,
) ([]appmodulev2.ValidatorUpdate, error) {
	return h.initGenesis(ctx, bz)
}

// initGenesis is called by the base app to initialize the state of the.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
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
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) PrepareProposal(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	return h.prepareProposal(ctx, req)
}

// prepareProposal is the internal handler for preparing proposals.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
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
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) waitForSidecars(gCtx context.Context) ([]byte, error) {
	select {
	case <-gCtx.Done():
		return nil, gCtx.Err()
	case err := <-h.prepareProposalErrCh:
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
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) waitforBeaconBlk(gCtx context.Context) ([]byte, error) {
	select {
	case <-gCtx.Done():
		return nil, gCtx.Err()
	case err := <-h.prepareProposalErrCh:
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
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	return h.processProposal(ctx, req)
}

// processProposal is the internal handler for processing proposals.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) processProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	var (
		startTime = time.Now()
	)
	defer h.metrics.measureProcessProposalDuration(startTime)

	args := []any{"beacon_block", true, "blob_sidecars", true}
	blk, err := h.beaconBlockGossiper.Request(ctx, req)
	if err != nil {
		args[1] = false
	}

	sidecars, err := h.blobGossiper.Request(ctx, req)
	if err != nil {
		args[3] = false
	}

	h.logger.Info("Received proposal with", args...)
	if err = h.chainService.ReceiveBlockAndBlobs(
		ctx, blk, sidecars,
	); errors.IsFatal(err) {
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, err
	}

	return &cmtabci.ProcessProposalResponse{
		Status: cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT,
	}, nil
}

/* -------------------------------------------------------------------------- */
/*                                FinalizeBlock                               */
/* -------------------------------------------------------------------------- */

// PreBlock is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) PreBlock(
	ctx sdk.Context, req *cmtabci.FinalizeBlockRequest,
) error {
	go h.preBlock(ctx, req)
	return nil
}

// handlePreBlock is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) preBlock(
	ctx sdk.Context, req *cmtabci.FinalizeBlockRequest,
) {
	blk, blobs, err := encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		h.chainSpec.ActiveForkVersionForSlot(
			math.Slot(req.Height),
		))

	if err != nil {
		h.finalizeBlockErrCh <- errors.Join(err, ErrBadExtractBlockAndBlocks)
		return
	}

	result, err := h.chainService.ProcessBlockAndBlobs(ctx, blk, blobs)
	if err != nil {
		h.finalizeBlockErrCh <- err
	} else {
		h.valUpdatesCh <- result
	}
}

// EndBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	return h.endBlock(ctx)
}

// endBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconStateT,
	BlobSidecarsT, DepositT, ExecutionPayloadT, GenesisT,
]) endBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-h.finalizeBlockErrCh:
		if errors.Is(err, ErrBadExtractBlockAndBlocks) {
			err = nil
		}
		return nil, err
	case result := <-h.valUpdatesCh:
		return iter.MapErr(
			result.RemoveDuplicates().Sort(), convertValidatorUpdate,
		)
	}
}
