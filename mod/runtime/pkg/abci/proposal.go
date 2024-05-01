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
	"errors"
	"time"

	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ssz "github.com/ferranbt/fastssz"
	"github.com/sourcegraph/conc/iter"
)

// Handler is a struct that encapsulates the necessary components to handle
// the proposal processes.
type Handler struct {
	builderService BuilderService
	chainService   BlockchainService
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(
	builderService BuilderService,
	chainService BlockchainService,
) *Handler {
	return &Handler{
		builderService: builderService,
		chainService:   chainService,
	}
}

// PrepareProposalHandler is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *Handler) PrepareProposalHandler(
	ctx sdk.Context, req *cmtabci.RequestPrepareProposal,
) (*cmtabci.ResponsePrepareProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyPrepareProposalTime, "ms")
	logger := ctx.Logger().With("module", "prepare-proposal")
	st := h.chainService.BeaconState(ctx)

	// Process the Slot to set the state root for the block.
	if err := h.chainService.ProcessSlot(st); err != nil {
		return &cmtabci.ResponsePrepareProposal{}, err
	}

	blk, blobs, err := h.builderService.RequestBestBlock(
		ctx, st, math.Slot(req.Height))
	if err != nil || blk == nil || blk.IsNil() {
		logger.Error("failed to build block", "error", err, "block", blk)
		return &cmtabci.ResponsePrepareProposal{}, err
	}

	// Serialize the block and blobs.
	txs, err := iter.MapErr[ssz.Marshaler, []byte](
		[]ssz.Marshaler{blk, blobs},
		func(m *ssz.Marshaler) ([]byte, error) {
			return (*m).MarshalSSZ()
		})

	return &cmtabci.ResponsePrepareProposal{
		Txs: txs,
	}, err
}

// ProcessProposalHandler is a wrapper around the process proposal handler
// that extracts the beacon block from the proposal and processes it.
func (h *Handler) ProcessProposalHandler(
	ctx sdk.Context, req *cmtabci.RequestProcessProposal,
) (*cmtabci.ResponseProcessProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyProcessProposalTime, "ms")
	logger := ctx.Logger().With("module", "process-proposal")

	// Unmarshal the beacon block from the abci request.
	blk, err := encoding.UnmarshalBeaconBlockFromABCIRequest(
		req, BeaconBlockTxIndex,
		h.chainService.ChainSpec().
			ActiveForkVersionForSlot(math.Slot(req.Height)))
	if err != nil {
		logger.Error(
			"failed to retrieve beacon block from request",
			"error",
			err,
		)

		return &cmtabci.ResponseProcessProposal{
			Status: cmtabci.ResponseProcessProposal_REJECT,
		}, nil
	}

	// If the block is syncing, we reject the proposal. This is guard against a
	// potential attack under the unlikely scenario in which a supermajority of
	// validators have their EL's syncing. If nodes were to accept this proposal
	// optmistically when they are syncing, it could potentially allow for a
	// malicious validator to push a bad block through.
	//
	// TODO: figure out a way to prevent newPayload from being called twiced as
	// it will be called again
	// in PreBlocker.
	if err = h.chainService.VerifyPayloadOnBlk(ctx, blk); errors.Is(
		err,
		engineclient.ErrSyncingPayloadStatus,
	) {
		return &cmtabci.ResponseProcessProposal{
			Status: cmtabci.ResponseProcessProposal_REJECT,
		}, err
	}

	return &cmtabci.ResponseProcessProposal{
		Status: cmtabci.ResponseProcessProposal_ACCEPT,
	}, nil
}
