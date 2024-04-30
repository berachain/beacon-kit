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

package abci

import (
	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/runtime/encoding"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FinalizeBlock is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *Handler) FinalizeBlock(
	ctx sdk.Context, req *cometabci.RequestFinalizeBlock,
) error {
	logger := ctx.Logger().With("module", "pre-block")

	// Extract the beacon block from the ABCI request.
	blk, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		req,
		BeaconBlockTxIndex,
		h.chainService.ChainSpec().ActiveForkVersionForSlot(
			math.Slot(req.Height),
		),
	)
	if err != nil {
		return err
	}

	blobSideCars, err := encoding.UnmarshalBlobSidecarsFromABCIRequest(
		req,
		BlobSidecarsTxIndex,
	)
	if err != nil {
		return err
	}

	st := h.chainService.BeaconState(ctx)

	// Process the Slot.
	if err = h.chainService.ProcessSlot(st); err != nil {
		logger.Error("failed to process slot", "error", err)
		return err
	}

	// Processing the incoming beacon block and blobs.
	stCopy := st.Copy()
	if err = h.chainService.ProcessBeaconBlock(
		ctx,
		stCopy,
		blk,
		blobSideCars,
	); err != nil {
		logger.Warn(
			"failed to receive beacon block",
			"error",
			err,
		)
		// TODO: Emit Evidence so that the validator can be slashed.
	} else {
		// We only want to persist state changes if we successfully
		// processed the block.
		stCopy.Save()
	}

	// Process the finalization of the beacon block.
	return h.chainService.PostBlockProcess(ctx, st, blk)
}
