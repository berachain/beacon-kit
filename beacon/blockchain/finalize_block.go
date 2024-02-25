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

	"github.com/ethereum/go-ethereum/common"
	"github.com/itsdevbear/bolaris/types/consensus"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// FinalizeBeaconBlock finalizes a beacon block by processing the logs,
// deposits,
// and voluntary exits. It also updates the finalized and safe eth1 block hashes
// on the beacon state.
func (s *Service) FinalizeBeaconBlock(
	ctx context.Context,
	blk consensus.ReadOnlyBeaconKitBlock,
	blockRoot [32]byte,
) error {
	payload, err := blk.ExecutionPayload()
	if err != nil {
		return err
	}

	// TODO: PROCESS LOGS HERE
	// TODO: PROCESS DEPOSITS HERE
	// TODO: PROCESS VOLUNTARY EXITS HERE

	eth1BlockHash := common.Hash(payload.GetBlockHash())
	state := s.BeaconState(ctx)
	state.SetFinalizedEth1BlockHash(eth1BlockHash)
	state.SetSafeEth1BlockHash(eth1BlockHash)
	state.SetLastValidHead(eth1BlockHash)
	state.SetParentBlockRoot(blockRoot)

	return nil
}

// PostFinalizeBeaconBlock is called after a beacon block has been finalized.
func (s *Service) PostFinalizeBeaconBlock(
	ctx context.Context,
	slot primitives.Slot,
) error {
	var err error
	state := s.BeaconState(ctx)
	if s.BuilderCfg().LocalBuilderEnabled {
		err = s.sendFCUWithAttributes(
			ctx,
			state.GetSafeEth1BlockHash(),
			slot,
			state.GetParentBlockRoot(),
		)
	} else {
		err = s.sendFCU(
			ctx, s.BeaconState(ctx).GetSafeEth1BlockHash())
	}

	return err
}
