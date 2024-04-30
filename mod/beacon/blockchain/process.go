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
	beacontypes "github.com/berachain/beacon-kit/mod/core/types"
	datypes "github.com/berachain/beacon-kit/mod/da/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"golang.org/x/sync/errgroup"
)

// ProcessSlot processes the incoming beacon slot.
func (s *Service) ProcessSlot(
	st state.BeaconState,
) error {
	return s.sp.ProcessSlot(st)
}

// ProcessBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) ProcessBeaconBlock(
	ctx context.Context,
	st state.BeaconState,
	blk primitives.ReadOnlyBeaconBlock,
	blobs *datypes.BlobSidecars,
) error {
	var (
		avs  = s.AvailabilityStore(ctx)
		g, _ = errgroup.WithContext(ctx)
		err  error
	)

	// If the block is nil, exit early.
	if blk == nil || blk.IsNil() {
		return beacontypes.ErrNilBlk
	}

	// Validate payload in Parallel.
	g.Go(func() error {
		return s.pv.VerifyPayload(st, blk.GetBody())
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
	if _, err = s.ee.VerifyAndNotifyNewPayload(
		ctx, engineprimitives.BuildNewPayloadRequest(
			body.GetExecutionPayload(),
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			false,
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
			st,
			avs,
			blobs,
		)
	})

	g.Go(func() error {
		return s.sp.ProcessBlock(
			st,
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

	// Prune deposits
	if err = s.sks.PruneDepositEvents(st); err != nil {
		s.Logger().Error("failed to prune deposit events", "error", err)
		return err
	}

	return nil
}

// ValidateBlock validates the incoming beacon block.
func (s *Service) ValidateBlock(
	ctx context.Context,
	blk primitives.ReadOnlyBeaconBlock,
) error {
	return s.bv.ValidateBlock(
		s.BeaconState(ctx), blk,
	)
}

// VerifyPayload validates the execution payload on the block.
func (s *Service) VerifyPayloadOnBlk(
	ctx context.Context,
	blk primitives.ReadOnlyBeaconBlock,
) error {
	if blk == nil || blk.IsNil() {
		return beacontypes.ErrNilBlk
	}

	body := blk.GetBody()
	if body.IsNil() {
		return beacontypes.ErrNilBlkBody
	}

	// Call the standard payload validator.
	if err := s.pv.VerifyPayload(
		s.BeaconState(ctx),
		body,
	); err != nil {
		return err
	}

	// We notify the engine of the new payload.
	parentBeaconBlockRoot := blk.GetParentBlockRoot()
	if _, err := s.ee.VerifyAndNotifyNewPayload(
		ctx,
		engineprimitives.BuildNewPayloadRequest(
			body.GetExecutionPayload(),
			body.GetBlobKzgCommitments().ToVersionedHashes(),
			&parentBeaconBlockRoot,
			false,
		),
	); err != nil {
		s.Logger().
			Error("failed to notify engine of new payload", "error", err)
		return err
	}
	return nil
}

// PostBlockProcess is called after a block has been processed.
// It is responsible for processing logs and other post block tasks.
func (s *Service) PostBlockProcess(
	ctx context.Context,
	st state.BeaconState,
	blk primitives.ReadOnlyBeaconBlock,
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
		s.Logger().Error("failed to process logs", "error", err)
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
		withdrawalsRoot, withdrawalsRootErr = primitives.Withdrawals(
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
