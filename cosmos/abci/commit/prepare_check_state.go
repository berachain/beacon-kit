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
	"math/big"

	"cosmossdk.io/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/execution/logs"
	v1 "github.com/itsdevbear/bolaris/types/v1"
)

type BeaconKeeper interface {
	ForkChoiceStore(ctx context.Context) v1.ForkChoiceStore
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

func (h *BeaconPrepareCheckStateHandler) PrepareCheckStater() sdk.PrepareCheckStater {
	return func(ctx sdk.Context) {
		fcs := h.beaconKeeper.ForkChoiceStore(ctx)
		finalHash := fcs.GetFinalizedBlockHash()
		if err := h.beaconChain.FinalizeBlockAsync(ctx, ctx.HeaderInfo(), finalHash[:]); err != nil {
			h.logger.Error("failed to finalize block", "err", err)
			panic(err)
		}

		// TODO THIS IS HACK
		if err := h.logProcessor.ProcessETH1Block(ctx, big.NewInt((ctx.BlockHeight()))); err != nil {
			panic(err)
		}
	}
}
