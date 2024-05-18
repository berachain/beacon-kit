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
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
)

// ProcessSlot processes the incoming beacon slot.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) ProcessSlot(
	st ReadOnlyBeaconStateT,
) error {
	return s.sp.ProcessSlot(st)
}

// ProcessStateTransition receives an incoming beacon block, it first validates
// and then processes the block.
//
//nolint:funlen // todo cleanup.
func (s *Service[
	ReadOnlyBeaconStateT, BlobSidecarsT, DepositStoreT,
]) ProcessStateTransition(
	ctx context.Context,
	st ReadOnlyBeaconStateT,
	blk types.BeaconBlock,
	blobs BlobSidecarsT,
) error {
	// If the block is nil, exit early.
	if blk == nil || blk.IsNil() {
		return ErrNilBlk
	}

	// Perform the state transition.
	if err := s.sp.Transition(
		// We set `OptimisticEngine` to true since this is called during
		// FinalizeBlock. We want to assume the payload is valid. If it
		// ends up not being valid later, the node will simply AppHash,
		// which is completely fine. This means we were syncing from a
		// bad peer, and we would likely AppHash anyways.
		//
		// TODO: Figure out why SkipPayloadIfExists being `true`
		// causes nodes to create gaps in their chain.
		core.NewContext(ctx, true, false, true),
		st,
		s.bsb.AvailabilityStore(ctx),
		blk,
		blobs,
	); err != nil {
		return err
	}

	// If the blobs needed to process the block are not available, we
	// return an error. It is safe to use the slot off of the beacon block
	// since it has been verified as correct already.
	if !s.bsb.AvailabilityStore(ctx).IsDataAvailable(
		ctx, blk.GetSlot(), blk.GetBody(),
	) {
		return ErrDataNotAvailable
	}

	// No matter what happens we always want to forkchoice at the end of post
	// block processing.
	defer func() {
		go s.sendPostBlockFCU(ctx, st, blk)
	}()

	//
	//
	//
	//
	//
	// TODO: EVERYTHING BELOW THIS LINE SHOULD NOT PART OF THE
	//  MAIN BLOCK PROCESSING THREAD.
	//
	//
	//
	//
	//
	//

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

	// We notify the engine of the new payload.
	var (
		parentBeaconBlockRoot = blk.GetParentBlockRoot()
		body                  = blk.GetBody()
		payload               = body.GetExecutionPayload()
	)

	if err := s.ee.VerifyAndNotifyNewPayload(
		ctx,
		engineprimitives.BuildNewPayloadRequest(
			payload,
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			false,
			// We do not want to optimistically assume truth here, since
			// this is being called in process proposal.
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
