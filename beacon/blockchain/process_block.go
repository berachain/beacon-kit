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
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/beacon/execution"
	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// TODO: calculate this based off of all the comet timeouts.
const approximateBlkTime = 3 * time.Second

// postBlockProcess(.
func (s *Service) postBlockProcess(
	ctx context.Context, blk consensus.ReadOnlyBeaconKitBlock, isValidPayload bool,
) error {
	if !isValidPayload {
		telemetry.IncrCounter(1, MetricReceivedInvalidPayload)
		// If the incoming payload for this block is not valid, we submit a forkchoice to bring us back
		// to the last valid one.
		// TODO: Is doing this potentially the cause of the weird Geth SnapSync issue?
		// TODO: Should introduce the concept of missed slots?
		if err := s.sendFCU(ctx, s.BeaconState(ctx).GetLastValidHead(), blk.GetSlot()+1); err != nil {
			s.Logger().Error("failed to send forkchoice update", "error", err)
		}
		return ErrInvalidPayload
	}

	executionPayload, err := blk.ExecutionPayload()
	if err != nil {
		return err
	}

	// We notify the execution client of the new block, and wait for it to return
	// a payload ID. If the payload ID is nil, we return an error. One thing to notice here however
	// is that we pass in `slot+1` to the execution client. We do this so that we can begin building
	// the next block in the background while we are finalizing this block.
	// We are okay pushing this asynchonous work to the execution client, as it is designed for it.
	//
	// TODO: we should probably just have a validator job in the background that is
	// constantly building new payloads and then not worry about anything here triggering
	// payload builds.
	return s.sendFCU(ctx, common.Hash(executionPayload.GetBlockHash()), blk.GetSlot()+1)
}

// sendFCU sends a forkchoice update to the execution client.
func (s *Service) sendFCU(
	ctx context.Context, headEth1Hash common.Hash, proposingSlot primitives.Slot,
) error {
	// If we are preparing all payloads and we are a validator, we delegate the responsibility
	// of submitting our forkchoice update to the builder service in order to have it prepare
	// a new execution payload for us. We discard the response, since we don't really
	// care about it in the context of processing the current block. As long as it doesn't error.
	if s.FeatureFlags().PrepareAllPayloads {
		//#nosec:G701 // won't overflow, time cannot be negative.
		_, err := s.bs.BuildLocalPayload(
			ctx, headEth1Hash, proposingSlot, uint64((time.Now().Add(approximateBlkTime)).Unix()),
		)
		if err == nil {
			return nil
		}

		// If we see an error here, we fallback to submitting a forkchoice without building a payload.
		// Incase there is a block building issue, but the payload was still totally valid, we don't
		// want this node to reject the block.
		s.Logger().Warn(
			"failed to build local payload via builder service prepare all payloads",
			"error", err)
		telemetry.IncrCounter(1, MetricFailedToBuildLocalPayload)
	}

	// If we did not return above we just send the forkchoice update to the execution client
	// with no attributes.
	// By not passing any attributes, the execution client will NOT build a new execution payload.
	_, err := s.en.NotifyForkchoiceUpdate(
		ctx, &execution.FCUConfig{
			HeadEth1Hash:  headEth1Hash,
			ProposingSlot: proposingSlot,
		})
	return err
}
