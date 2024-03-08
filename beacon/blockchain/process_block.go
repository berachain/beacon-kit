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
	"errors"

	beacontypes "github.com/itsdevbear/bolaris/beacon/core/types"
	"github.com/itsdevbear/bolaris/beacon/sync"
)

// postBlockProcess(.
func (s *Service) postBlockProcess(
	ctx context.Context,
	blk beacontypes.ReadOnlyBeaconBlock,
	blockHash [32]byte,
	_ bool,
) error {
	// If the block does not have a payload, we return an error.
	payload := blk.GetBody().GetExecutionPayload()
	if payload.IsNil() {
		return ErrInvalidPayload
	}
	payloadBlockHash := payload.GetBlockHash()

	// If the consensus client is still syncing we are going to skip forkchoice
	// updates. This means that if the consensus client is still syncing, we
	// will not be able to build a block locally.
	//
	// NOTE: Status() will return nil during the initial syncing phase.
	if errors.Is(s.ss.Status(), sync.ErrConsensusClientIsSyncing) {
		return nil
	}

	// If the builder is enabled attempt to build a block locally.
	// If we are in the sync state, we skip building blocks optimistically.
	if s.BuilderCfg().LocalBuilderEnabled && !s.ss.IsInitSync() {
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
