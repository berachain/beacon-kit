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

	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/execution"
	enginetypes "github.com/berachain/beacon-kit/mod/execution/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"golang.org/x/sync/errgroup"
)

// ProcessSlot processes the incoming beacon slot.
func (s *Service) ProcessSlot(
	ctx context.Context,
) error {
	// Process the slot.
	return s.sp.ProcessSlot(
		s.BeaconState(ctx),
	)
}

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) ProcessBeaconBlock(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
	blobs *datypes.BlobSidecars,
) error {
	var (
		avs         = s.AvailabilityStore(ctx)
		g, groupCtx = errgroup.WithContext(ctx)
		st          = s.BeaconState(groupCtx)
		err         error
	)

	// Validate payload in Parallel.
	g.Go(func() error {
		return s.pv.ValidatePayload(st, blk.GetBody())
	})

	// Validate block in Parallel.
	g.Go(func() error {
		return s.bv.ValidateBlock(st, blk)
	})

	// Wait for the errgroup to finish, the error will be non-nil if any
	// of the goroutines returned an error.
	if err = g.Wait(); err != nil {
		// If we fail any checks we process the slot and move on.
		return err
	}

	// Then we notify the engine of the new payload.
	body := blk.GetBody()
	parentBeaconBlockRoot := blk.GetParentBlockRoot()
	versionedHashes := primitives.KzgCommitmentsToVersionedHashes(
		body.GetBlobKzgCommitments())
	if _, err = s.ee.VerifyAndNotifyNewPayload(
		ctx,
		execution.BuildNewPayloadRequest(
			body.GetExecutionPayload(),
			versionedHashes,
			&parentBeaconBlockRoot,
		),
	); err != nil {
		s.Logger().
			Error("failed to notify engine of new payload", "error", err)
		return err
	}

	// We want to get a headstart on blob processing since it
	// is a relatively expensive operation.
	g.Go(func() error {
		return s.sp.ProcessBlobs(
			avs,
			blk,
			blobs,
		)
	})

	g.Go(func() error {
		return s.sp.ProcessBlock(
			s.BeaconState(ctx),
			blk,
		)
	})

	// Wait for the errgroup to finish, the error will be non-nil if any
	// of the goroutines returned an error.
	if err = g.Wait(); err != nil {
		// If we fail any checks we process the slot and move on.
		return err
	}

	// TODO: Validate the data availability as well as check for the
	// minimum DA required time.
	// daStartTime := time.Now()
	// if avs != nil {
	// avs.IsDataAvailable(ctx, s.CurrentSlot(), rob); err != nil {
	// 		return errors.Wrap(err, "could not validate blob data availability
	// (AvailabilityStore.IsDataAvailable)")
	// 	}
	// } else {
	// s.isDataAvailable(ctx, blockRoot, blockCopy); err != nil {
	// 		return errors.Wrap(err, "could not validate blob data availability")
	// 	}
	// }

	return nil
}

// PostBlockProcess is called after a block has been processed.
// It is responsible for processing logs and other post block tasks.
func (s *Service) PostBlockProcess(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
) error {
	var (
		payload enginetypes.ExecutionPayload
		st      = s.BeaconState(ctx)
	)

	// No matter what happens we always want to forkchoice at the end of post
	// block processing.
	defer func(payloadPtr *enginetypes.ExecutionPayload) {
		s.sendPostBlockFCU(ctx, *payloadPtr)
	}(&payload)

	// If the block is nil, exit early.
	if blk.IsNil() {
		return nil
	}

	body := blk.GetBody()
	if body.IsNil() {
		return nil
	}

	// Update the forkchoice.
	payload = blk.GetBody().GetExecutionPayload()
	if payload.IsNil() {
		return nil
	}

	prevEth1Block, err := st.GetEth1BlockHash()
	if err != nil {
		return err
	}

	// Process the logs in the block.
	if err = s.sks.ProcessLogsInETH1Block(
		ctx,
		prevEth1Block,
	); err != nil {
		s.Logger().Error("failed to process logs", "error", err)
		return err
	}

	payloadBlockHash := payload.GetBlockHash()
	if err = st.UpdateEth1BlockHash(payloadBlockHash); err != nil {
		return err
	}

	return nil
}
