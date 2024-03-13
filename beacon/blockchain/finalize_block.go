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
	"github.com/berachain/beacon-kit/primitives"
)

// FinalizeBeaconBlock finalizes a beacon block by processing the logs,
// deposits,
// and voluntary exits. It also updates the finalized and safe eth1 block hashes
// on the beacon state.
func (s *Service) FinalizeBeaconBlock(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
	blockRoot [32]byte,
) error {
	var (
		err         error
		state       = s.BeaconState(ctx)
		forkChoicer = s.ForkchoiceStore(ctx)
	)

	defer func() {
		// Always update the parent block root in the event
		// that the beacon block is not valid.
		state.SetParentBlockRoot(blockRoot)

		// If something bad happens, we defensivelessly send a forkchoice update
		// to bring us back to the last valid head.
		go func() {
			if err != nil {
				s.missedBlockTasks(ctx, blk.GetSlot(), blockRoot)
			}

			s.Logger().Info(
				"finalizing current forkchoice state",
				"safe_hash", forkChoicer.JustifiedPayloadBlockHash().Hex(),
				"finalized_hash", forkChoicer.FinalizedPayloadBlockHash().Hex(),
			)
		}()
	}()

	if blk.IsNil() {
		return beacontypes.ErrNilBlk
	}

	payload := blk.GetBody().GetExecutionPayload()
	if payload.IsNil() {
		// TODO: Slash the proposer for not including a payload.
		return ErrNoPayloadInBeaconBlock
	}

	payloadBlockHash := payload.GetBlockHash()
	if err = forkChoicer.InsertNode(payloadBlockHash); err != nil {
		return err
	}

	if err = s.rp.MixinNewReveal(ctx, blk); err != nil {
		return err
	}

	// TODO: PROCESS LOGS HERE
	// TODO: PROCESS DEPOSITS HERE
	// TODO: PROCESS VOLUNTARY EXITS HERE
	err = s.es.ProcessLogsInETH1Block(
		ctx,
		payloadBlockHash,
	)
	if err != nil {
		s.Logger().Error("failed to process logs", "error", err)
		return err
	}

	// TODO: put into an actual function / flow
	_, err = s.BeaconState(ctx).DequeueDeposits(
		uint64(len(blk.GetBody().GetDeposits())))
	return err
}

// missed block tasks is called when a block is missed. It sends a forkchoice
// update to the execution client to bring the client back to the last valid
// head (safe).
func (s *Service) missedBlockTasks(
	ctx context.Context,
	slot primitives.Slot,
	blockRoot [32]byte,
) {
	forkChoicer := s.ForkchoiceStore(ctx)

	// If we are in the sync state, we skip building blocks
	// optimistically.
	if s.BuilderCfg().LocalBuilderEnabled && !s.ss.IsInitSync() {
		err := s.sendFCUWithAttributes(
			ctx,
			forkChoicer.JustifiedPayloadBlockHash(),
			slot,
			blockRoot,
		)
		if err == nil {
			return
		}
		s.Logger().Error(
			"failed to send recovery forkchoice update w/attributes", "error", err,
		)
	}

	// Otherwise we send a forkchoice update to the execution client.
	err := s.sendFCU(ctx, forkChoicer.JustifiedPayloadBlockHash())
	if err != nil {
		s.Logger().Error(
			"failed to send recovery forkchoice update", "error", err,
		)
	}
}
