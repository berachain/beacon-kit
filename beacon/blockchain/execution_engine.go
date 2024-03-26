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

	"github.com/berachain/beacon-kit/beacon/execution"
	enginetypes "github.com/berachain/beacon-kit/engine/types"
	"github.com/berachain/beacon-kit/primitives"
)

// sendFCU sends a forkchoice update to the execution client.
func (s *Service) sendFCU(
	ctx context.Context,
	headEth1Hash primitives.ExecutionHash,
) error {
	_, err := s.es.NotifyForkchoiceUpdate(
		ctx, &execution.FCUConfig{
			HeadEth1Hash: headEth1Hash,
		})
	return err
}

// sendFCUWithAttributes sends a forkchoice update to the
// execution client with payload attributes. It does
// so via the local builder service.
func (s *Service) sendFCUWithAttributes(
	ctx context.Context,
	headEth1Hash primitives.ExecutionHash,
	forSlot primitives.Slot,
	parentBlockRoot [32]byte,
) error {
	_, err := s.lb.BuildLocalPayload(
		ctx,
		headEth1Hash,
		forSlot,
		//#nosec:G701 // won't realistically overflow.
		uint64(time.Now().Unix()),
		parentBlockRoot,
	)
	return err
}

// sendPostBlockFCU sends a forkchoice update to the execution client.
func (s *Service) sendPostBlockFCU(
	ctx context.Context,
	payload enginetypes.ExecutionPayload,
) {
	var (
		headHash primitives.ExecutionHash
		st       = s.BeaconState(ctx)
		fcs      = s.ForkchoiceStore(ctx)
	)

	// If we have a payload we want to set our head to it's block hash.
	// Otherwise we are going to use the justified payload block hash.
	if payload != nil {
		headHash = payload.GetBlockHash()
	} else {
		headHash = fcs.JustifiedPayloadBlockHash()
	}

	// If we are the local builder and we are not in init sync
	// forkchoice update with attributes.
	//nolint:nestif // todo:cleanup
	if s.BuilderCfg().LocalBuilderEnabled && !s.ss.IsInitSync() {
		// TODO: This BlockRoot calculation is sound, but very confusing
		// and hard to explain to someone who is not familiar with the
		// nuance of our implementation. We should refactor this.
		h, err := st.GetLatestBlockHeader()
		if err != nil {
			s.Logger().
				Error("failed to get latest block header in postBlockProcess", "error", err)
			return
		}

		stateRoot, err := st.HashTreeRoot()
		if err != nil {
			s.Logger().
				Error("failed to get state root in postBlockProcess", "error", err)
			return
		}

		h.StateRoot = stateRoot
		root, err := h.HashTreeRoot()
		if err != nil {
			s.Logger().
				Error("failed to get block header root in postBlockProcess", "error", err)
			return
		}

		slot, err := st.GetSlot()
		if err != nil {
			s.Logger().
				Error("failed to get slot in postBlockProcess", "error", err)
		}

		if err = s.sendFCUWithAttributes(
			ctx,
			headHash,
			slot+1,
			root,
		); err == nil {
			return
		}

		// If we error we log and continue, we try again without building a
		// block
		// just incase this can help get our execution client back on track.
		s.Logger().
			Error("failed to send forkchoice update with attributes", "error", err)
	}

	// Otherwise we send a forkchoice update to the execution client.
	if err := s.sendFCU(ctx, headHash); err != nil {
		s.Logger().
			Error("failed to send forkchoice update in postBlockProcess", "error", err)
	}
}
