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

	"github.com/berachain/beacon-kit/beacon/core"
	beacontypes "github.com/berachain/beacon-kit/beacon/core/types"
	"github.com/berachain/beacon-kit/crypto/kzg"
)

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) ProcessBeaconBlock(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
	blockHash [32]byte,
) error {
	// TODO:
	// expectedProposer, err := epc.GetBeaconProposer(benv.Slot)

	// This go rountine validates the execution level aspects of the block.
	// i.e: does newPayload return VALID?
	isValidPayload, err := s.validateExecutionOnBlock(
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
	// Create a new state processor.
	sp := core.NewStateProcessor(
		s.BeaconCfg(),
		s.rp,
	)
	bv := core.NewBlockValidator(
		s.BeaconCfg(),
	)

	// Validate the block
	if err := bv.ValidatorBlock(
		s.BeaconState(ctx), blk,
	); err != nil {
		return err
	}

	return sp.ProcessBlock(
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
		pv      = core.NewPayloadValidator(
			s.BeaconCfg(),
		)
	)

	// Validate the payload.
	if err := pv.ValidatePayload(
		s.BeaconState(ctx),
		s.ForkchoiceStore(ctx),
		payload,
	); err != nil {
		return false, err
	}

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

// postBlockProcess is called after a block has been processed.
func (s *Service) postBlockProcess(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
	blockHash [32]byte,
	_ bool,
) error {
	// If the block is nil or empty, we return an error.
	if blk == nil || blk.IsNil() {
		return beacontypes.ErrNilBlk
	}

	// If the block does not have a payload, we return an error.
	payload := blk.GetBody().GetExecutionPayload()
	if payload.IsNil() {
		return ErrInvalidPayload
	}
	payloadBlockHash := payload.GetBlockHash()

	// If the builder is enabled attempt to build a block locally.
	// If we are in the sync state, we skip building blocks optimistically.
	if s.BuilderCfg().LocalBuilderEnabled && !s.ss.IsInitSync() {
		// We have to do this in order to update it before FCU.
		// TODO: In general we need to improve the control flow for
		// Preblocker vs ProcessProposal.
		err := s.sendFCUWithAttributes(
			ctx, payloadBlockHash, blk.GetSlot(), blockHash,
		)
		if err == nil {
			return nil
		}
		s.Logger().
			Error("failed to send forkchoice update in postBlockProcess", "error", err)
	}

	// Otherwise we send a forkchoice update to the execution client.
	return s.sendFCU(ctx, payloadBlockHash)
}
