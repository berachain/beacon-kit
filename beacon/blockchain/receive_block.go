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
	"bytes"
	"context"
	"fmt"

	"github.com/itsdevbear/bolaris/types/consensus"
	"golang.org/x/sync/errgroup"
)

// ReceiveBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) ReceiveBeaconBlock(
	ctx context.Context,
	blk consensus.ReadOnlyBeaconKitBlock,
) error {
	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	var (
		eg, groupCtx   = errgroup.WithContext(ctx)
		isValidPayload bool
	)

	// TODO: Do we need to wait for the forkchoice to update?
	// TODO: move the error group to use GCD.

	// This go routine validators the consensus level aspects of the block.
	// i.e: does it have a valid ancesor?
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

	// Wait for the goroutines to finish.
	if err := eg.Wait(); err != nil {
		return err
	}

	// If the block is valid, we can process it.
	return s.postBlockProcess(
		ctx, blk, isValidPayload,
	)
}

// validateStateTransition checks a block's state transition.
// TODO: Expand rules, consider modularity. Current implementation
// is hardcoded for single slot finality, which works but lacks flexibility.
func (s *Service) validateStateTransition(
	ctx context.Context, blk consensus.ReadOnlyBeaconKitBlock,
) error {
	executionData, err := blk.ExecutionPayload()
	if err != nil {
		return err
	}

	finalizedHash := s.BeaconState(ctx).GetFinalizedEth1BlockHash()
	if !bytes.Equal(finalizedHash[:], executionData.GetParentHash()) {
		return fmt.Errorf(
			"parent block with hash %x is not finalized, expected finalized hash %x",
			executionData.GetParentHash(),
			finalizedHash,
		)
	}

	// TODO: Probably add RANDAO and Staking stuff here?

	// TODO: how do we handle hard fork boundaries?

	return nil
}

// validateExecutionOnBlock checks the validity of a proposed beacon block.
func (s *Service) validateExecutionOnBlock(
	ctx context.Context, blk consensus.ReadOnlyBeaconKitBlock,
) (bool, error) {
	payload, err := blk.ExecutionPayload()
	if err != nil {
		return false, err
	}

	// TODO: add some more safety checks here.

	return s.es.NotifyNewPayload(ctx, payload, blk.GetSlot())
}
