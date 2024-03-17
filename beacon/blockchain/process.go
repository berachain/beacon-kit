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

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
	"github.com/berachain/beacon-kit/primitives"
	"golang.org/x/sync/errgroup"
)

// ProcessSlot is run when a slot is missed.
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
) error {
	var (
		body        = blk.GetBody()
		root        primitives.HashRoot
		payload     = body.GetExecutionPayload()
		fcs         = s.ForkchoiceStore(ctx)
		g, groupCtx = errgroup.WithContext(ctx)
		st          = s.BeaconState(groupCtx)
		err         error
	)

	// We can use an errgroup to run the validation functions concurrently.
	g.Go(func() error {
		return s.pv.ValidatePayload(st, fcs, payload)
	})

	g.Go(func() error {
		return s.bv.ValidateBlock(st, blk)
	})

	// Wait for the errgroup to finish, the error will be non-nil if any
	// of the goroutines returned an error.
	if err = g.Wait(); err != nil {
		// If we fail any checks we process the slot and move on.
		return err
	}

	// This go rountine validates the execution level aspects of the block.
	// i.e: does newPayload return VALID?
	_, err = s.validateExecutionOnBlock(
		ctx, blk,
	)
	if err != nil {
		s.Logger().
			Error("failed to notify engine of new payload", "error", err)
		return err
	}

	// This go routine validates the consensus level aspects of the block.
	// i.e: does it have a valid ancestor?
	if err = s.validateStateTransition(ctx, blk); err != nil {
		s.Logger().
			Error("failed to validate state transition", "error", err)
		return err
	}

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

	// Perform post block processing.
	// If the builder is enabled attempt to build a block locally.
	// If we are in the sync state, we skip building blocks optimistically.
	if s.BuilderCfg().LocalBuilderEnabled && !s.ss.IsInitSync() {
		if root, err = blk.HashTreeRoot(); err != nil {
			return err
		}

		err = s.sendFCUWithAttributes(
			ctx,
			payload.GetBlockHash(),
			blk.GetSlot(),
			root,
		)
		if err == nil {
			return nil
		}
		s.Logger().
			Error("failed to send forkchoice update in postBlockProcess", "error", err)
	}

	// Otherwise we send a forkchoice update to the execution client.
	return s.sendFCU(ctx, payload.GetBlockHash())
}

// validateStateTransition checks a block's state transition.
// TODO: Expand rules, consider modularity. Current implementation
// is hardcoded for single slot finality, which works but lacks flexibility.
func (s *Service) validateStateTransition(
	ctx context.Context, blk beacontypes.ReadOnlyBeaconBlock,
) error {
	return s.sp.ProcessBlock(
		s.BeaconState(ctx),
		blk,
	)
}

// validateExecutionOnBlock checks the validity of a the execution payload
// on the beacon block.
func (s *Service) validateExecutionOnBlock(
	// todo: parentRoot hashs should be on blk.
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
) (bool, error) {
	var (
		body    = blk.GetBody()
		payload = body.GetExecutionPayload()
	)

	// Then we notify the engine of the new payload.
	return s.es.NotifyNewPayload(
		ctx,
		blk.GetSlot(),
		payload,
		kzg.ConvertCommitmentsToVersionedHashes(
			body.GetBlobKzgCommitments(),
		),
		blk.GetParentBlockRoot(),
	)
}
