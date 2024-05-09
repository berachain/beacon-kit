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

	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

// sendFCU sends a forkchoice update to the execution client.
// It sets the head and finalizes the latest.
func (s *Service[BlobSidecarsT]) sendFCU(
	ctx context.Context,
	st state.BeaconState,
	headEth1Hash common.ExecutionHash,
) error {
	latestExecutionPayloadHeader, err := st.GetLatestExecutionPayloadHeader()
	if err != nil {
		return err
	}
	eth1BlockHash := latestExecutionPayloadHeader.GetBlockHash()

	_, _, err = s.ee.NotifyForkchoiceUpdate(
		ctx,
		&engineprimitives.ForkchoiceUpdateRequest{
			State: &engineprimitives.ForkchoiceState{
				HeadBlockHash:      headEth1Hash,
				SafeBlockHash:      eth1BlockHash,
				FinalizedBlockHash: eth1BlockHash,
			},
		},
	)
	return err
}

// sendPostBlockFCU sends a forkchoice update to the execution client.
func (s *Service[BlobSidecarsT]) sendPostBlockFCU(
	ctx context.Context,
	st state.BeaconState,
	payload engineprimitives.ExecutionPayload,
) {
	var (
		headHash common.ExecutionHash
	)

	// If we have a payload we want to set our head to it's block hash.
	// Otherwise we are going to use the justified payload block hash.
	if payload != nil {
		headHash = payload.GetBlockHash()
	} else {
		latestExecutionPayloadHeader, err := st.GetLatestExecutionPayloadHeader()
		if err != nil {
			s.logger.Error(
				"failed to get latest execution payload in postBlockProcess",
				"error", err,
			)
			return
		}
		headHash = latestExecutionPayloadHeader.GetBlockHash()
	}

	// If we are the local builder and we are not in init sync
	// forkchoice update with attributes.
	//nolint:nestif // todo:cleanup
	// TODO: re-enable this flag.
	if true /*s.BuilderCfg().LocalBuilderEnabled */ /*&& !s.ss.IsInitSync()*/ {
		// TODO: This BlockRoot calculation is sound, but very confusing
		// and hard to explain to someone who is not familiar with the
		// nuance of our implementation. We should refactor this.
		h, err := st.GetLatestBlockHeader()
		if err != nil {
			s.logger.
				Error(
					"failed to get latest block header in postBlockProcess",
					"error",
					err,
				)
			return
		}

		stateRoot, err := st.HashTreeRoot()
		if err != nil {
			s.logger.
				Error(
					"failed to get state root in postBlockProcess",
					"error",
					err,
				)
			return
		}

		h.StateRoot = stateRoot
		root, err := h.HashTreeRoot()
		if err != nil {
			s.logger.
				Error(
					"failed to get block header root in postBlockProcess",
					"error",
					err,
				)
			return
		}

		slot, err := st.GetSlot()
		if err != nil {
			s.logger.
				Error("failed to get slot in postBlockProcess", "error", err)
		}

		stCopy := st.Copy()
		if err = s.sp.ProcessSlot(stCopy); err != nil {
			return
		}

		// Ask the builder to send a forkchoice update with attributes.
		// This will trigger a new payload to be built.
		if _, err = s.lb.RequestPayload(
			ctx,
			stCopy,
			slot+1,
			//#nosec:G701 // won't realistically overflow.
			// TODO: clock time properly.
			uint64(time.Now().Unix()+1),
			root,
			headHash,
		); err == nil {
			return
		}

		// If we error we log and continue, we try again without building a
		// block
		// just incase this can help get our execution client back on track.
		s.logger.
			Error(
				"failed to send forkchoice update with attributes",
				"error",
				err,
			)
	}

	// Otherwise we send a forkchoice update to the execution client.
	if err := s.sendFCU(ctx, st, headHash); err != nil {
		s.logger.
			Error(
				"failed to send forkchoice update in postBlockProcess",
				"error",
				err,
			)
	}
}
