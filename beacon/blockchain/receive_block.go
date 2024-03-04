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
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	beacontypes "github.com/itsdevbear/bolaris/beacon/core/types"
	"github.com/itsdevbear/bolaris/crypto/kzg"
	"golang.org/x/sync/errgroup"
)

// ReceiveBeaconBlock receives an incoming beacon block, it first validates
// and then processes the block.
func (s *Service) ReceiveBeaconBlock(
	ctx context.Context,
	blockHash [32]byte,
	buoy beacontypes.ReadOnlyBeaconBuoy,
) error {
	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	var (
		eg, groupCtx   = errgroup.WithContext(ctx)
		isValidPayload bool
		forkChoicer    = s.ForkchoiceStore(ctx)
	)

	// If we have already seen this block, we can skip processing it.
	// TODO: should we store some historical data here?
	if forkChoicer.GetLastSeenBeaconBlock() == blockHash {
		s.Logger().Info(
			"ignoring already processed beacon block",
			// todo: don't use common for beacontypes
			"hash", common.Hash(blockHash).Hex(),
		)
		return nil
	}
	forkChoicer.SetLastSeenBeaconBlock(blockHash)

	// This go routine validates the consensus level aspects of the block.
	// i.e: does it have a valid ancestor?
	eg.Go(func() error {
		err := s.validateStateTransition(groupCtx, buoy)
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
			groupCtx, buoy,
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
		ctx, buoy, blockHash, isValidPayload,
	)
}

// validateStateTransition checks a block's state transition.
// TODO: Expand rules, consider modularity. Current implementation
// is hardcoded for single slot finality, which works but lacks flexibility.
func (s *Service) validateStateTransition(
	ctx context.Context, buoy beacontypes.ReadOnlyBeaconBuoy,
) error {
	if err := beacontypes.BeaconBuoyIsNil(buoy); err != nil {
		return err
	}

	parentBlockRoot := s.BeaconState(ctx).GetParentBlockRoot()
	if !bytes.Equal(parentBlockRoot[:], buoy.GetParentBlockRoot()) {
		return fmt.Errorf(
			"parent root does not match, expected: %x, got: %x",
			parentBlockRoot,
			buoy.GetParentBlockRoot(),
		)
	}

	// TODO: Probably add RANDAO and Staking stuff here?

	// TODO: how do we handle hard fork boundaries?

	return nil
}

// validateExecutionOnBlock checks the validity of a the execution payload
// on the beacon block.
func (s *Service) validateExecutionOnBlock(
	// todo: parentRoot hashs should be on blk.
	ctx context.Context,
	buoy beacontypes.ReadOnlyBeaconBuoy,
) (bool, error) {
	if err := beacontypes.BeaconBuoyIsNil(buoy); err != nil {
		return false, err
	}

	payload, err := buoy.ExecutionPayload()
	if err != nil {
		return false, err
	}

	if payload == nil || payload.IsEmpty() {
		return false, errors.New("no payload in beacon block")
	}

	// In BeaconKit, since we are currently operating on SingleSlot Finality
	// we purposefully reject any block that is not a child of the last
	// finalized block.
	safeHash := s.ForkchoiceStore(ctx).GetSafeEth1BlockHash()
	if !bytes.Equal(safeHash[:], payload.GetParentHash()) {
		return false, fmt.Errorf(
			"parent block with hash %x is not finalized, expected finalized hash %x",
			payload.GetParentHash(),
			safeHash,
		)
	}

	// TODO: add some more safety checks here.
	return s.es.NotifyNewPayload(
		ctx,
		buoy.GetSlot(),
		payload,
		kzg.ConvertCommitmentsToVersionedHashes(buoy.GetBlobKzgCommitments()),
		common.Hash(buoy.GetParentBlockRoot()),
	)
}
