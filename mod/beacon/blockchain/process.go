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

	"github.com/berachain/beacon-kit/mod/core/state"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"golang.org/x/sync/errgroup"
)

// ProcessSlot processes the incoming beacon slot.
func (s *Service[BlobSidecarsT]) ProcessSlot(
	st state.BeaconState,
) error {
	return s.sp.ProcessSlot(st)
}

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service[BlobSidecarsT]) ProcessBeaconBlock(
	ctx context.Context,
	st state.BeaconState,
	blk consensus.ReadOnlyBeaconBlock[consensus.BeaconBlockBody],
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
		)
	})

	// We ask for the slot before waiting as a minor optimization.
	slot, err := st.GetSlot()
	if err != nil {
		return err
	}

	// Wait for the errgroup to finish, the error will be non-nil if any
	// of the goroutines returned an error.
	if err = g.Wait(); err != nil {
		// If we fail any checks we process the slot and move on.
		return err
	}

	// If the blobs needed to process the block are not available, we
	// return an error.
	if !s.bsb.AvailabilityStore(ctx).IsDataAvailable(ctx, slot, body) {
		return ErrDataNotAvailable
	}

	// Prune deposits.
	// TODO: This should be moved into a go-routine in the background.
	// Watching for logs should be completely decoupled as well.
	if err = s.sks.PruneDepositEvents(st); err != nil {
		s.logger.Error("failed to prune deposit events", "error", err)
		return err
	}

	return nil
}

// ValidateBlock validates the incoming beacon block.
func (s *Service[BlobSidecarsT]) ValidateBlock(
	ctx context.Context,
	blk consensus.ReadOnlyBeaconBlock[consensus.BeaconBlockBody],
) error {
	return s.bv.ValidateBlock(
		s.bsb.BeaconState(ctx), blk,
	)
}

// VerifyPayload validates the execution payload on the block.
func (s *Service[BlobSidecarsT]) VerifyPayloadOnBlk(
	ctx context.Context,
	blk consensus.ReadOnlyBeaconBlock[consensus.BeaconBlockBody],
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
		s.bsb.BeaconState(ctx),
		body.GetExecutionPayload(),
	); err != nil {
		return err
	}

	// We notify the engine of the new payload.
	parentBeaconBlockRoot := blk.GetParentBlockRoot()
	return s.ee.VerifyAndNotifyNewPayload(
		ctx,
		engineprimitives.BuildNewPayloadRequest(
			body.GetExecutionPayload(),
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			false,
			// We do not want to optimistically assume truth here.
			false,
		),
	)
}

// PostBlockProcess is called after a block has been processed.
// It is responsible for processing logs and other post block tasks.
func (s *Service[BlobSidecarsT]) PostBlockProcess(
	ctx context.Context,
	st state.BeaconState,
	blk consensus.ReadOnlyBeaconBlock[consensus.BeaconBlockBody],
) error {
	var (
		payload engineprimitives.ExecutionPayload
	)

	// No matter what happens we always want to forkchoice at the end of post
	// block processing.
	defer func(payloadPtr *engineprimitives.ExecutionPayload) {
		s.sendPostBlockFCU(ctx, st, *payloadPtr)
	}(&payload)

	// If the block is nil, exit early.
	if blk == nil || blk.IsNil() {
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

	latestExecutionPayloadHeader, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}
	prevEth1Block := latestExecutionPayloadHeader.GetBlockHash()

	// Process the logs in the block.
	if err = s.sks.ProcessLogsInETH1Block(ctx, prevEth1Block); err != nil {
		s.logger.Error("failed to process logs", "error", err)
		return err
	}

	// Get the merkle roots of transactions and withdrawals in parallel.
	var (
		g, _            = errgroup.WithContext(ctx)
		txsRoot         primitives.Root
		withdrawalsRoot primitives.Root
	)

	g.Go(func() error {
		var txsRootErr error
		txsRoot, txsRootErr = engineprimitives.Transactions(
			payload.GetTransactions(),
		).HashTreeRoot()
		return txsRootErr
	})

	g.Go(func() error {
		var withdrawalsRootErr error
		withdrawalsRoot, withdrawalsRootErr = consensus.Withdrawals(
			payload.GetWithdrawals(),
		).HashTreeRoot()
		return withdrawalsRootErr
	})

	// If deriving either of the roots fails, return the error.
	if err = g.Wait(); err != nil {
		return err
	}

	// Set the latest execution payload header.
	if err = st.SetLatestExecutionPayloadHeader(
		&engineprimitives.ExecutionPayloadHeaderDeneb{
			ParentHash:       payload.GetParentHash(),
			FeeRecipient:     payload.GetFeeRecipient(),
			StateRoot:        payload.GetStateRoot(),
			ReceiptsRoot:     payload.GetReceiptsRoot(),
			LogsBloom:        payload.GetLogsBloom(),
			Random:           payload.GetPrevRandao(),
			Number:           payload.GetNumber(),
			GasLimit:         payload.GetGasLimit(),
			GasUsed:          payload.GetGasUsed(),
			Timestamp:        payload.GetTimestamp(),
			ExtraData:        payload.GetExtraData(),
			BaseFeePerGas:    payload.GetBaseFeePerGas(),
			BlockHash:        payload.GetBlockHash(),
			TransactionsRoot: txsRoot,
			WithdrawalsRoot:  withdrawalsRoot,
			BlobGasUsed:      payload.GetBlobGasUsed(),
			ExcessBlobGas:    payload.GetExcessBlobGas(),
		},
	); err != nil {
		return err
	}

	return nil
}
