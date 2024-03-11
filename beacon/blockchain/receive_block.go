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
	"fmt"

	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
	"github.com/berachain/beacon-kit/primitives"
	"golang.org/x/sync/errgroup"
)

// ReceiveBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) ReceiveBeaconBlock(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
	blockHash [32]byte,
) error {
	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	var (
		eg, groupCtx   = errgroup.WithContext(ctx)
		isValidPayload bool
		forkChoicer    = s.ForkchoiceStore(ctx)
	)

	// If the block is nil, We have to abort.
	if blk == nil || blk.IsNil() {
		return beacontypes.ErrNilBlk
	}

	// If we have already seen this block, we can skip processing it.
	// TODO: should we store some historical data here?
	if forkChoicer.HeadBeaconBlock() == blockHash {
		s.Logger().Info(
			"ignoring already processed beacon block",
			// todo: don't use common for beacontypes
			"hash", primitives.ExecutionHash(blockHash).Hex(),
		)
		return nil
	}
	forkChoicer.UpdateHeadBeaconBlock(blockHash)

	// This go routine validates the consensus level aspects of the block.
	// i.e: does it have a valid ancestor?
	eg.Go(func() error {
		err := s.validateStateTransition(groupCtx, blk)
		if err != nil {
			s.Logger().
				Error("failed to validate state transition", "error", err)
			return err
		}
		return nil
	})

	// This go rountine validates the execution level aspects of the block.
	// i.e: does newPayload return VALID?
	eg.Go(func() error {
		var err error
		if isValidPayload, err = s.validateExecutionOnBlock(
			groupCtx, blk,
		); err != nil {
			s.Logger().
				Error("failed to notify engine of new payload", "error", err)
			return err
		}

		return nil
	})

	// daStartTime := time.Now()
	// if avs != nil {
	// 	if err := avs.IsDataAvailable(ctx, s.CurrentSlot(), rob); err != nil {
	// 		return errors.Wrap(err, "could not validate blob data availability
	// (AvailabilityStore.IsDataAvailable)")
	// 	}
	// } else {
	// 	if err := s.isDataAvailable(ctx, blockRoot, blockCopy); err != nil {
	// 		return errors.Wrap(err, "could not validate blob data availability")
	// 	}
	// }

	// Wait for the goroutines to finish.
	if err := eg.Wait(); err != nil {
		return err
	}

	// Perform post block processing.
	return s.postBlockProcess(
		ctx, blk, blockHash, isValidPayload,
	)
}

// validateStateTransition checks a block's state transition.
// TODO: Expand rules, consider modularity. Current implementation
// is hardcoded for single slot finality, which works but lacks flexibility.
func (s *Service) validateStateTransition(
	ctx context.Context, blk beacontypes.ReadOnlyBeaconBlock,
) error {
	if blk.IsNil() {
		return beacontypes.ErrNilBlk
	}

	parentBlockRoot := s.BeaconState(ctx).GetParentBlockRoot()
	if parentBlockRoot != blk.GetParentBlockRoot() {
		return fmt.Errorf(
			"parent root does not match, expected: %x, got: %x",
			parentBlockRoot,
			blk.GetParentBlockRoot(),
		)
	}

	// Ensure Body is non nil.
	body := blk.GetBody()
	if body.IsNil() {
		return beacontypes.ErrNilBlkBody
	}

	// ---------------------///
	// VALIDATE RANDAO HERE ///
	// ---------------------///

	// ---------------------///
	//   VALIDATE KZG HERE  ///
	// ---------------------///

	// Ensure the block deposits are within the limits.
	deposits := body.GetDeposits()
	if uint64(len(deposits)) > s.BeaconCfg().Limits.MaxDepositsPerBlock {
		return fmt.Errorf(
			"too many deposits, expected: %d, got: %d",
			s.BeaconCfg().Limits.MaxDepositsPerBlock, len(deposits),
		)
	}

	// Ensure the deposits match the local state.
	localDeposits, err := s.BeaconState(ctx).
		ExpectedDeposits(uint64(len(deposits)))
	if err != nil {
		return err
	}

	// Ensure the deposits match the local state.
	for i, dep := range deposits {
		if dep == nil {
			return beacontypes.ErrNilDeposit
		}
		if dep.Index != localDeposits[i].Index {
			return fmt.Errorf(
				"deposit index does not match, expected: %d, got: %d",
				localDeposits[i].Index, dep.Index)
		}
	}

	return nil
}

// validateExecutionOnBlock checks the validity of a the execution payload
// on the beacon block.
func (s *Service) validateExecutionOnBlock(
	// todo: parentRoot hashs should be on blk.
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
) (bool, error) {
	if blk.IsNil() {
		return false, beacontypes.ErrNilBlk
	}

	body := blk.GetBody()
	payload := body.GetExecutionPayload()
	if payload.IsNil() {
		return false, beacontypes.ErrNilPayloadInBlk
	}

	// In BeaconKit, since we are currently operating on SingleSlot Finality
	// we purposefully reject any block that is not a child of the last
	// finalized block.
	safeHash := s.ForkchoiceStore(ctx).JustifiedPayloadBlockHash()
	if safeHash != payload.GetParentHash() {
		return false, fmt.Errorf(
			"parent block with hash %x is not finalized, expected finalized hash %x",
			payload.GetParentHash(),
			safeHash,
		)
	}

	// TODO: add some more safety checks here.
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
