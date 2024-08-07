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

package runtime

import (
	"context"
	"fmt"
	"time"

	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/json"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	"golang.org/x/sync/errgroup"
)

/* -------------------------------------------------------------------------- */
/*                                 InitGenesis                                */
/* -------------------------------------------------------------------------- */

// InitGenesis is called by the base app to initialize the state of the.
func (a *App[
	_, _, _, _, _, _, _, GenesisT, _, _, _,
]) InitChain(
	ctx context.Context,
	bz []byte,
) (transition.ValidatorUpdates, []byte, error) {
	var (
		g          errgroup.Group
		valUpdates transition.ValidatorUpdates
		genesisErr error
	)
	data := new(GenesisT)
	if err := json.Unmarshal(bz, data); err != nil {
		return nil, nil, err
	}
	// Send a request to the chain service to process the genesis data.
	if err := a.genesisBroker.Publish(ctx, asynctypes.NewEvent(
		ctx, events.GenesisDataProcessRequest, *data,
	)); err != nil {
		return nil, nil, err
	}

	// Wait for the genesis data to be processed.
	g.Go(func() error {
		valUpdates, genesisErr = a.waitForGenesisData(ctx)
		return genesisErr
	})

	if err := g.Wait(); err != nil {
		return nil, nil, err
	}

	// stateHash, err := a.sb.StateFromContext(ctx).HashTreeRoot()
	// if err != nil {
	// 	return nil, nil, err
	// }

	stateHash, err := a.sb.StateFromContext(ctx).LatestCommitHash()
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("STATE HASH FROM APP GENESIS", hex.FromBytes(stateHash))
	return valUpdates, stateHash, nil
}

// waitForGenesisData waits for the genesis data to be processed and returns
// the validator updates.
func (a *App[
	_, _, _, _, _, _, _, _, _, _, _,
]) waitForGenesisData(
	ctx context.Context,
) (transition.ValidatorUpdates, error) {
	select {
	case msg := <-a.valUpdateSub:
		if msg.Type() != events.ValidatorSetUpdated {
			return nil, errors.Wrapf(
				ErrUnexpectedEvent,
				"unexpected event type: %s", msg.Type(),
			)
		}
		return msg.Data(), msg.Error()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

/* -------------------------------------------------------------------------- */
/*                               PrepareProposal                              */
/* -------------------------------------------------------------------------- */

// prepareProposal is the internal handler for preparing proposals.
func (a *App[
	_, _, _, _, _, _, _, _, _, SlotDataT, _,
]) PrepareProposal(
	ctx context.Context,
	req *types.PrepareRequest,
) ([][]byte, error) {
	var (
		g                           errgroup.Group
		startTime                   = time.Now()
		beaconBlockErr, sidecarsErr error
		beaconBlockBz, sidecarsBz   []byte
	)
	defer a.metrics.measurePrepareProposalDuration(startTime)
	slotData, err := a.convertPrepareProposalToSlotData(ctx, req)
	if err != nil {
		return nil, err
	}

	// Send a request to the validator service to give us a beacon block
	// and blob sidecards to pass to ABCI.
	if err := a.slotBroker.Publish(ctx, asynctypes.NewEvent(
		ctx, events.NewSlot, slotData,
	)); err != nil {
		return nil, err
	}

	// Wait for the beacon block to be built.
	g.Go(func() error {
		beaconBlockBz, beaconBlockErr = a.waitforBeaconBlk(ctx)
		return beaconBlockErr
	})

	// Wait for the sidecars to be built.
	g.Go(func() error {
		sidecarsBz, sidecarsErr = a.waitForSidecars(ctx)
		return sidecarsErr
	})

	// Wait for both processes to complete and then
	// return the appropriate response.
	return [][]byte{beaconBlockBz, sidecarsBz}, g.Wait()
}

// waitForSidecars waits for the sidecars to be built and returns them.
func (a *App[
	_, _, _, _, _, _, _, _, _, _, _,
]) waitForSidecars(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case msg := <-a.sidecarsCh:
		if msg.Error() != nil {
			return nil, msg.Error()
		}
		return a.blobGossiper.Publish(ctx, msg.Data())
	}
}

// waitforBeaconBlk waits for the beacon block to be built and returns it.
func (a *App[
	_, _, _, _, _, _, _, _, _, _, _,
]) waitforBeaconBlk(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case beaconBlock := <-a.blkCh:
		if beaconBlock.Error() != nil {
			return nil, beaconBlock.Error()
		}
		return a.beaconBlockGossiper.Publish(
			ctx,
			beaconBlock.Data(),
		)
	}
}

/* -------------------------------------------------------------------------- */
/*                               ProcessProposal                              */
/* -------------------------------------------------------------------------- */

// ProcessProposal processes the proposal for the ABCI middleware.
// It handles both the beacon block and blob sidecars concurrently.
func (a *App[
	_, _, BeaconBlockT, _, BlobSidecarsT, _, _, _, _, _, _,
]) ProcessProposal(
	ctx context.Context,
	req *types.ProcessRequest,
) error {
	var (
		blk       BeaconBlockT
		sidecars  BlobSidecarsT
		err       error
		g, _      = errgroup.WithContext(ctx)
		startTime = time.Now()
	)

	defer a.metrics.measureProcessProposalDuration(startTime)

	// Request the beacon block.
	if blk, err = a.beaconBlockGossiper.Request(ctx, req); err != nil {
		a.logger.Error("error requesting beacon block", "error", err)
		return nil
	}

	// Begin processing the beacon block.
	g.Go(func() error {
		return a.verifyBeaconBlock(ctx, blk)
	})

	// Request the blob sidecars.
	if sidecars, err = a.blobGossiper.Request(ctx, req); err != nil {
		a.logger.Error("error requesting blob sidecars", "error", err)
		return nil
	}

	// Begin processing the blob sidecars.
	g.Go(func() error {
		return a.verifyBlobSidecars(ctx, sidecars)
	})

	// Wait for both processes to complete and then
	// return the appropriate response.s
	return g.Wait()
}

// verifyBeaconBlock handles the processing of the beacon block.
// It requests the block, publishes a received event, and waits for
// verification.
func (a *App[
	_, _, BeaconBlockT, _, _, _, _, _, _, _, _,
]) verifyBeaconBlock(
	ctx context.Context,
	blk BeaconBlockT,
) error {
	// Publish the received event.
	if err := a.blkBroker.Publish(
		ctx,
		asynctypes.NewEvent(ctx, events.BeaconBlockReceived, blk, nil),
	); err != nil {
		return err
	}

	// Wait for a response.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-a.blkCh:
		if msg.Type() != events.BeaconBlockVerified {
			return errors.Wrapf(
				ErrUnexpectedEvent, "unexpected event type: %s", msg.Type(),
			)
		}
		return msg.Error()
	}
}

// processBlobSidecars handles the processing of blob sidecars.
// It requests the sidecars, publishes a received event, and waits for
// processing.
func (a *App[
	_, _, _, _, BlobSidecarsT, _, _, _, _, _, _,
]) verifyBlobSidecars(
	ctx context.Context,
	sidecars BlobSidecarsT,
) error {
	// Publish the received event.
	if err := a.sidecarsBroker.Publish(
		ctx,
		asynctypes.NewEvent(ctx, events.BlobSidecarsReceived, sidecars),
	); err != nil {
		return err
	}

	// Wait for a response.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-a.sidecarsCh:
		if msg.Type() != events.BlobSidecarsProcessed {
			return errors.Wrapf(
				ErrUnexpectedEvent, "unexpected event type: %s", msg.Type(),
			)
		}
		return msg.Error()
	}
}

/* -------------------------------------------------------------------------- */
/*                                FinalizeBlock                               */
/* -------------------------------------------------------------------------- */

// EndBlock returns the validator set updates from the beacon state.
func (a *App[
	_, _, BeaconBlockT, _, BlobSidecarsT, _, _, _, _, _, _,
]) FinalizeBlock(
	ctx context.Context,
	req *types.FinalizeRequest,
) (transition.ValidatorUpdates, []byte, error) {
	blk, blobs, err := encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](
		req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		a.chainSpec.ActiveForkVersionForSlot(
			math.Slot(req.Slot),
		))
	if err != nil {
		// If we don't have a block, we can't do anything.
		//nolint:nilerr // by design.

		// stateHash, err := a.sb.StateFromContext(ctx).LatestCommitHash()
		// if err != nil {
		// 	return nil, nil, err
		// }
		stateHash, err := a.sb.StateFromContext(ctx).Commit()
		if err != nil {
			return nil, nil, err
		}
		return nil, stateHash, nil
	}

	// Send the sidecars to the sidecars feed and wait for a response
	if err = a.processSidecars(ctx, blobs); err != nil {
		return nil, nil, err
	}

	// Process the beacon block and return the validator updates.
	valUpdates, err := a.processBeaconBlock(
		ctx, blk,
	)
	if err != nil {
		return nil, nil, err
	}

	// fmt.Println("STATE HASH FROM APP FINALIZE", a.sb.StateFromContext(ctx).HashTreeRoot())

	stateHash, err := a.sb.StateFromContext(ctx).LatestCommitHash()
	if err != nil {
		return nil, nil, err
	}
	return valUpdates, stateHash, err
}

// processSidecars publishes the sidecars and waits for a response.
func (a *App[
	_, _, _, _, BlobSidecarsT, _, _, _, _, _, _,
]) processSidecars(ctx context.Context, blobs BlobSidecarsT) error {
	// Publish the sidecars.
	if err := a.sidecarsBroker.Publish(ctx, asynctypes.NewEvent(
		ctx, events.BlobSidecarsProcessRequest, blobs,
	)); err != nil {
		return err
	}

	// Wait for the sidecars to be processed.
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-a.sidecarsCh:
		if msg.Type() != events.BlobSidecarsProcessed {
			return errors.Wrapf(
				ErrUnexpectedEvent,
				"unexpected event type: %s", msg.Type(),
			)
		}
		return msg.Error()
	}
}

// processBeaconBlock processes the beacon block and returns validator updates.
func (a *App[
	_, _, BeaconBlockT, _, _, _, _, _, _, _, _,
]) processBeaconBlock(
	ctx context.Context, blk BeaconBlockT,
) (transition.ValidatorUpdates, error) {
	// Publish the verified block event.
	if err := a.blkBroker.Publish(
		ctx, asynctypes.NewEvent(
			ctx, events.BeaconBlockFinalizedRequest, blk,
		)); err != nil {
		return nil, err
	}

	// Wait for the block to be processed.
	select {
	case msg := <-a.valUpdateSub:
		if msg.Type() != events.ValidatorSetUpdated {
			return nil, errors.Wrapf(
				ErrUnexpectedEvent,
				"unexpected event type: %s", msg.Type(),
			)
		}
		return msg.Data(), msg.Error()
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
