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

package preblock

import (
	"context"

	"github.com/prysmaticlabs/prysm/v4/math"

	"cosmossdk.io/log"

	cometabci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	consensustypes "github.com/itsdevbear/bolaris/beacon/consensus-types"
	v1 "github.com/itsdevbear/bolaris/types/v1"
)

type BeaconKeeper interface {
	ForkChoiceStore(ctx context.Context) v1.ForkChoiceStore
}

// BeaconPreBlockHandler is responsible for aggregating oracle data from each
// validator and writing the oracle data into the store before any transactions
// are executed/finalized for a given block.
type BeaconPreBlockHandler struct {
	logger log.Logger

	// keeper is the keeper for the oracle module. This is utilized to write
	// oracle data to state.
	keeper BeaconKeeper
}

// NewBeaconPreBlockHandler returns a new BeaconPreBlockHandler. The handler
// is responsible for writing oracle data included in vote extensions to state.
func NewBeaconPreBlockHandler(
	logger log.Logger,
	beaconKeeper BeaconKeeper,
) *BeaconPreBlockHandler {
	return &BeaconPreBlockHandler{
		logger: logger,
		keeper: beaconKeeper,
	}
}

// PreBlocker is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *BeaconPreBlockHandler) PreBlocker() sdk.PreBlocker {
	return func(ctx sdk.Context, req *cometabci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
		h.logger.Info(
			"executing the pre-finalize block hook",
			"height", req.Height,
		)

		beaconBlockData := req.Txs[0] // todo modularize.

		// todo handle hardforks without needing codechange.
		data, err := consensustypes.BytesToExecutionData(
			beaconBlockData, math.Gwei(0), 3) //nolint:gomnd // fix later.
		if err != nil {
			h.logger.Error("payload in beacon block could not be unmarshalled", "err", err)
			return nil, err
		}

		// Finalize the block that is being proposed.
		store := h.keeper.ForkChoiceStore(ctx)
		store.SetFinalizedBlockHash([32]byte(data.BlockHash()))
		store.SetSafeBlockHash([32]byte(data.BlockHash()))
		store.SetLastValidHead([32]byte(data.BlockHash()))
		return &sdk.ResponsePreBlock{}, nil
	}
}
