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

package blockchain

import (
	"context"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/events"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"golang.org/x/sync/errgroup"
)

// ProcessGenesisData processes the genesis state and initializes the beacon
// state.
func (s *Service[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) ProcessGenesisData(
	ctx context.Context,
	genesisData *genesis.Genesis[
		*types.Deposit, *types.ExecutionPayloadHeaderDeneb,
	],
) ([]*transition.ValidatorUpdate, error) {
	return s.sp.InitializePreminedBeaconStateFromEth1(
		s.sb.StateFromContext(ctx),
		genesisData.Deposits,
		genesisData.ExecutionPayloadHeader,
		genesisData.ForkVersion,
	)
}

// ProcessBlockAndBlobs receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) ProcessBlockAndBlobs(
	ctx context.Context,
	blk types.BeaconBlock,
	sidecars BlobSidecarsT,
) ([]*transition.ValidatorUpdate, error) {
	var (
		g, gCtx    = errgroup.WithContext(ctx)
		st         = s.sb.StateFromContext(ctx)
		valUpdates []*transition.ValidatorUpdate
	)

	// If the block is nil, exit early.
	if blk == nil || blk.IsNil() {
		return nil, ErrNilBlk
	}

	// Launch a goroutine to process the incoming beacon block.
	g.Go(func() error {
		var err error
		// We set `OptimisticEngine` to true since this is called during
		// FinalizeBlock. We want to assume the payload is valid. If it
		// ends up not being valid later, the node will simply AppHash,
		// which is completely fine. This means we were syncing from a
		// bad peer, and we would likely AppHash anyways.
		valUpdates, err = s.processBeaconBlock(gCtx, st, blk, true)
		return err
	})

	// Launch a goroutine to process the blob sidecars.
	g.Go(func() error {
		return s.processBlobSidecars(gCtx, blk.GetSlot(), sidecars)
	})

	// Wait for the goroutines to finish.
	if err := g.Wait(); err != nil {
		return nil, err
	}

	// If the blobs needed to process the block are not available, we
	// return an error. It is safe to use the slot off of the beacon block
	// since it has been verified as correct already.
	if !s.sb.AvailabilityStore(ctx).IsDataAvailable(
		ctx, blk.GetSlot(), blk.GetBody(),
	) {
		return nil, ErrDataNotAvailable
	}

	// emit new block event
	s.blockFeed.Send(events.NewBlock(ctx, blk))

	// No matter what happens we always want to forkchoice at the end of post
	// block processing.
	// TODO: this is hood as fuck.
	// We won't send a fcu if the block is bad, should be addressed
	// via ticker later.
	go s.sendPostBlockFCU(ctx, st, blk)
	go s.postBlockProcessTasks(ctx, st)

	return valUpdates, nil
}

// postBlockProcessTasks performs post block processing tasks.
//
// TODO: Deprecate this function and move it's usage outside of the main block
// processing thread.
func (s *Service[
	AvailabilityStoreT, BeaconStateT,
	BlobSidecarsT, DepositStoreT,
]) postBlockProcessTasks(
	ctx context.Context,
	st BeaconStateT,
) {
	// Prune deposits.
	// TODO: This should be moved into a go-routine in the background.
	// Watching for logs should be completely decoupled as well.
	idx, err := st.GetEth1DepositIndex()
	if err != nil {
		s.logger.Error(
			"failed to get eth1 deposit index in postBlockProcessTasks",
			"error", err)
		return
	}

	// TODO: pruner shouldn't be in main block processing thread.
	if err = s.PruneDepositEvents(ctx, idx); err != nil {
		s.logger.Error(
			"failed to prune deposit events in postBlockProcessTasks",
			"error", err)
		return
	}
}

// ProcessBeaconBlock processes the beacon block.
func (s *Service[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) ProcessBeaconBlock(
	ctx context.Context,
	blk types.BeaconBlock,
) ([]*transition.ValidatorUpdate, error) {
	st := s.sb.StateFromContext(ctx)
	return s.processBeaconBlock(ctx, st, blk, false)
}

// ProcessBeaconBlock processes the beacon block.
func (s *Service[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) processBeaconBlock(
	ctx context.Context,
	st BeaconStateT,
	blk types.BeaconBlock,
	optimisticEngine bool,
) ([]*transition.ValidatorUpdate, error) {
	startTime := time.Now()
	defer s.metrics.measureStateTransitionDuration(startTime)
	valUpdates, err := s.sp.Transition(
		&transition.Context{
			Context:          ctx,
			OptimisticEngine: optimisticEngine,
		},
		st,
		blk,
	)
	return valUpdates, err
}

// ProcessBlobSidecars processes the blob sidecars.
func (s *Service[
	AvailabilityStoreT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) processBlobSidecars(
	ctx context.Context,
	slot math.Slot,
	sidecars BlobSidecarsT,
) error {
	startTime := time.Now()
	defer s.metrics.measureBlobProcessingDuration(startTime)
	return s.bp.ProcessBlobs(
		slot,
		s.sb.AvailabilityStore(ctx),
		sidecars,
	)
}
