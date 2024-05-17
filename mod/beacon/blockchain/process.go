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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"golang.org/x/sync/errgroup"
)

// ProcessSlot processes the incoming beacon slot.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) ProcessSlot(
	st ReadOnlyBeaconStateT,
) error {
	return s.sp.ProcessSlot(st)
}

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
//
//nolint:funlen // todo cleanup.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) ProcessBeaconBlock(
	ctx context.Context,
	st ReadOnlyBeaconStateT,
	blk types.BeaconBlock,
	blobs BlobSidecarsT,
) error {
	var (
		g, _ = errgroup.WithContext(ctx)
		err  error
	)

	// If the block is nil, exit early.
	if blk == nil || blk.IsNil() {
		return ErrNilBlk
	}

	// Validate payload in Parallel.
	g.Go(func() error {
		body := blk.GetBody()
		if body == nil || body.IsNil() {
			return ErrNilBlkBody
		}
		return s.pv.VerifyPayload(st, body.GetExecutionPayload())
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

	// We want to get a headstart on blob processing since it
	// is a relatively expensive operation.
	g.Go(func() error {
		return s.sp.ProcessBlobs(
			st,
			s.bsb.AvailabilityStore(ctx),
			blobs,
		)
	})

	body := blk.GetBody()

	// We can also parallelize the call to the execution layer.
	g.Go(func() error {
		// Then we notify the engine of the new payload.
		parentBeaconBlockRoot := blk.GetParentBlockRoot()
		if err = s.ee.VerifyAndNotifyNewPayload(
			ctx, engineprimitives.BuildNewPayloadRequest(
				body.GetExecutionPayload(),
				body.GetBlobKzgCommitments().ToVersionedHashes(),
				&parentBeaconBlockRoot,
				false,
				// Since this is called during FinalizeBlock, we want to assume
				// the payload is valid, if it ends up not being valid later the
				// node will simply AppHash which is completely fine, since this
				// means we were syncing from a bad peer, and we would likely
				// AppHash anyways.
				true,
			),
		); err != nil {
			s.logger.
				Error("failed to notify engine of new payload", "error", err)
			return err
		}

		// We also want to verify the payload on the block.
		return s.sp.ProcessBlock(
			st,
			blk,
			true,
		)
	})

	// Wait for the errgroup to finish, the error will be non-nil if any
	// of the goroutines returned an error.
	if err = g.Wait(); err != nil {
		// If we fail any checks we process the slot and move on.
		return err
	}

	// If the blobs needed to process the block are not available, we
	// return an error. It is safe to use the slot off of the beacon block
	// since it has been verified as correct already.
	if !s.bsb.AvailabilityStore(ctx).IsDataAvailable(
		ctx, blk.GetSlot(), body,
	) {
		return ErrDataNotAvailable
	}

	// No matter what happens we always want to forkchoice at the end of post
	// block processing.
	defer func() {
		go s.sendPostBlockFCU(ctx, st, blk)
	}()

	// TODO: EVERYTHING BELOW THIS LINE SHOULD NOT PART OF THE
	//  MAIN BLOCK PROCESSING THREAD.

	// Prune deposits.
	// TODO: This should be moved into a go-routine in the background.
	// Watching for logs should be completely decoupled as well.
	idx, err := st.GetEth1DepositIndex()
	if err != nil {
		return err
	}

	// TODO: pruner shouldn't be in main block processing thread.
	if err = s.PruneDepositEvents(ctx, idx); err != nil {
		return err
	}

	var latestExecutionPayloadHeader engineprimitives.ExecutionPayloadHeader
	latestExecutionPayloadHeader, err = st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}

	// Process the logs from the previous blocks execution payload.
	// TODO: This should be moved out of the main block processing flow.
	// TODO: eth1FollowDistance should be done actually proper
	eth1FollowDistance := math.U64(1)
	if err = s.retrieveDepositsFromBlock(
		ctx, latestExecutionPayloadHeader.GetNumber()-eth1FollowDistance,
	); err != nil {
		s.logger.Error("failed to process logs", "error", err)
		return err
	}

	return nil
}

// ValidateBlock validates the incoming beacon block.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) ValidateBlock(
	ctx context.Context,
	blk types.ReadOnlyBeaconBlock[types.BeaconBlockBody],
) error {
	return s.bv.ValidateBlock(
		s.bsb.StateFromContext(ctx), blk,
	)
}

// VerifyPayload validates the execution payload on the block.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) VerifyPayloadOnBlk(
	ctx context.Context,
	blk types.BeaconBlock,
) error {
	if blk == nil || blk.IsNil() {
		return ErrNilBlk
	}

	body := blk.GetBody()
	if body.IsNil() {
		return ErrNilBlkBody
	}

	// Call the standard payload validator.
	if err := s.pv.VerifyPayload(
		s.bsb.StateFromContext(ctx),
		body.GetExecutionPayload(),
	); err != nil {
		return err
	}

	// We notify the engine of the new payload.
	parentBeaconBlockRoot := blk.GetParentBlockRoot()
	payload := body.GetExecutionPayload()
	if err := s.ee.VerifyAndNotifyNewPayload(
		ctx,
		engineprimitives.BuildNewPayloadRequest(
			payload,
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			false,
			// We do not want to optimistically assume truth here.
			false,
		),
	); err != nil {
		return err
	}

	s.logger.Info(
		"successfully verified execution payload 💸",
		"payload-block-number", payload.GetNumber(),
		"num-txs", len(payload.GetTransactions()),
	)
	return nil
}
