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

package preblock

import (
	"context"

	"cosmossdk.io/log"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	"github.com/itsdevbear/bolaris/beacon/core/state"
	sync "github.com/itsdevbear/bolaris/beacon/sync"
	"github.com/itsdevbear/bolaris/config"
	abcitypes "github.com/itsdevbear/bolaris/runtime/abci/types"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

type BeaconKeeper interface {
	BeaconState(ctx context.Context) state.BeaconState
}

// BeaconPreBlockHandler is responsible for aggregating oracle data from each
// validator and writing the oracle data into the store before any transactions
// are executed/finalized for a given block.
type BeaconPreBlockHandler struct {
	// cfg is the configuration for block proposals and finalization.
	cfg *config.ABCI

	// logger is the logger used by the handler.
	logger log.Logger

	// chainService is the service that is responsible for interacting with
	// the beacon chain.
	chainService *blockchain.Service

	// syncStatus is the service that is responsible for determining if the
	// node is currently syncing.
	syncStatus *sync.Service

	// nextHandler is the next pre-block handler in the chain. This is always
	// nesting of the next pre-block handler into this handler.
	nextHandler sdk.PreBlocker
}

// NewBeaconPreBlockHandler returns a new BeaconPreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewBeaconPreBlockHandler(
	cfg *config.ABCI,
	logger log.Logger,
	chainService *blockchain.Service,
	syncService *sync.Service,
	nextHandler sdk.PreBlocker,
) *BeaconPreBlockHandler {
	return &BeaconPreBlockHandler{
		cfg:          cfg,
		logger:       logger,
		chainService: chainService,
		syncStatus:   syncService,
		nextHandler:  nextHandler,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *BeaconPreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(
		ctx sdk.Context, req *cometabci.RequestFinalizeBlock,
	) (*sdk.ResponsePreBlock, error) {
		// Extract the beacon block from the ABCI request.
		//
		// TODO: Block factory struct?
		// TODO: Use protobuf and .(type)?
		beaconBlock, err := abcitypes.ReadOnlyBeaconKitBlockFromABCIRequest(
			req,
			h.cfg.BeaconBlockPosition,
			h.chainService.ActiveForkVersionForSlot(
				primitives.Slot(req.Height),
			),
		)
		if err != nil {
			return nil, err
		}

		// Process the finalization of the beacon block.
		if err = h.chainService.FinalizeBeaconBlock(
			ctx, beaconBlock, [32]byte(req.Hash),
		); err != nil {
			return nil, err
		}

		// If there is no child handler, we are done, this preblocker
		// does not modify any consensus params so we return an empty
		// response.
		if h.nextHandler == nil {
			return &sdk.ResponsePreBlock{}, nil
		}

		// Call the nested child handler.
		return h.nextHandler(ctx, req)
	}
}
