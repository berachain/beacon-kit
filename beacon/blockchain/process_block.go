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

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/itsdevbear/bolaris/types/consensus"
)

// postBlockProcess(.
func (s *Service) postBlockProcess(
	ctx context.Context,
	blk consensus.ReadOnlyBeaconKitBlock,
	blockHash [32]byte,
	isValidPayload bool,
) error {
	if !isValidPayload {
		telemetry.IncrCounter(1, MetricReceivedInvalidPayload)

		// If the incoming payload for this block is not valid, we submit a
		// forkchoice
		// to bring us back to the last valid one.
		// TODO: Is doing this potentially the cause of the weird Geth SnapSync
		// issue?
		// TODO: Should introduce the concept of missed slots?
		if err := s.sendFCU(
			ctx, s.BeaconState(ctx).GetLastValidHead(),
		); err != nil {
			s.Logger().Error("failed to send forkchoice update", "error", err)
		}
		return ErrInvalidPayload
	}

	payload, err := blk.ExecutionPayload()
	if err != nil {
		return err
	}

	// If the builder is enabled attempt to build a block locally.
	if s.BuilderCfg().LocalBuilderEnabled {
		if err = s.sendFCUWithAttributes(
			ctx, [32]byte(payload.GetBlockHash()), blk.GetSlot(), blockHash,
		); err == nil {
			return nil
		}
	}

	// If builder is not enabled, or failed to build, fallback to a vanilla
	// fcu.
	return s.sendFCU(ctx, [32]byte(payload.GetBlockHash()))
}
