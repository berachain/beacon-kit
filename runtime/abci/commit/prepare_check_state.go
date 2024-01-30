// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package commit

import (
	"context"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/logs"
	"github.com/itsdevbear/bolaris/types"
)

type BeaconKeeper interface {
	ForkChoiceStore(ctx context.Context) types.ForkChoiceStore
}

type BeaconPrepareCheckStateHandler struct {
	logger       log.Logger
	beaconKeeper BeaconKeeper
	beaconChain  *blockchain.Service
	childHandler sdk.PrepareCheckStater
	logProcessor *logs.Processor
}

func NewBeaconPrepareCheckStateHandler(
	logger log.Logger,
	beaconKeeper BeaconKeeper,
	beaconChain *blockchain.Service,
	childHandler sdk.PrepareCheckStater,
	logProcessor *logs.Processor,
) *BeaconPrepareCheckStateHandler {
	return &BeaconPrepareCheckStateHandler{
		logger:       logger,
		beaconKeeper: beaconKeeper,
		beaconChain:  beaconChain,
		childHandler: childHandler,
		logProcessor: logProcessor,
	}
}

// This should maybe be moved to Precommit for the sole reason of if for some
// reason the finalization on eth1 fails, we should be given the opportunity to
// abort committing the beacon block before we write the block as final on the execution layer.
// TODO: later.

// Also lets explore instead of finalizing the block async, we finalize it synchonrously
// but allow for some sort of --fast-sync flag, which will call some sort of forkchoice
// update at node start, to begin syncing' the execution client to the  head block, specified
// as part of the --fast-sync command?
func (h *BeaconPrepareCheckStateHandler) PrepareCheckStater() sdk.PrepareCheckStater {
	return func(ctx sdk.Context) {
		// fcs := h.beaconKeeper.ForkChoiceStore(ctx)
		// finalHash := fcs.GetFinalizedEth1BlockHash()
		// if err := h.beaconChain.FinalizeBlockAsync(ctx, ctx.HeaderInfo(), finalHash[:]); err != nil {
		// 	h.logger.Error("failed to finalize block", "error", err)
		// 	panic(err)
		// }

		// TODO THIS IS HACK and needs to be moved to either preblock of block n+1 or precommit of
		// this block, cause we can't perform db writes in prepare check state, since the block
		// // was just committed.
		// if err := h.logProcessor.ProcessETH1Block(ctx, big.NewInt((ctx.BlockHeight()))); err != nil {
		// 	panic(err)
		// }
	}
}
