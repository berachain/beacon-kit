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

package proposal

import (
	"bytes"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/itsdevbear/bolaris/beacon/blockchain"
	builder "github.com/itsdevbear/bolaris/beacon/builder/local"
	sync "github.com/itsdevbear/bolaris/beacon/sync"
	"github.com/itsdevbear/bolaris/config"
	abcitypes "github.com/itsdevbear/bolaris/runtime/abci/types"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"
)

// Handler is a struct that encapsulates the necessary components to handle
// the proposal processes.
type Handler struct {
	cfg             *config.ABCI
	builderService  *builder.Service
	chainService    *blockchain.Service
	syncService     *sync.Service
	prepareProposal sdk.PrepareProposalHandler
	processProposal sdk.ProcessProposalHandler
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(
	cfg *config.ABCI,
	builderService *builder.Service,
	syncService *sync.Service,
	chainService *blockchain.Service,
	prepareProposal sdk.PrepareProposalHandler,
	processProposal sdk.ProcessProposalHandler,
) *Handler {
	return &Handler{
		cfg:             cfg,
		builderService:  builderService,
		syncService:     syncService,
		chainService:    chainService,
		prepareProposal: prepareProposal,
		processProposal: processProposal,
	}
}

// PrepareProposalHandler is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *Handler) PrepareProposalHandler(
	ctx sdk.Context, req *abci.RequestPrepareProposal,
) (*abci.ResponsePrepareProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyPrepareProposalTime, "ms")
	logger := ctx.Logger().With("module", "prepare-proposal")

	// TODO: Make this more sophisticated.
	//nolint:lll // couldnt fix.
	if bsp := h.syncService.CheckSyncStatus(ctx); bsp.Status == sync.StatusExecutionAhead {
		return nil, fmt.Errorf(
			"err: %w, status: %d", ErrValidatorClientNotSynced, bsp.Status,
		)
	}

	// TODO abstract this into BeaconState()
	parentRoot := ctx.BlockHeader().AppHash
	if req.Height == 1 {
		parentRoot = make([]byte, 32) //nolint:gomnd //temp
	}

	// We start by requesting the validator service to build us a block. This
	// may be from pulling a previously built payload from the local cache or it
	// may be by asking for a forkchoice from the execution client, depending on
	// timing.
	block, err := h.builderService.RequestBestBlock(
		ctx, primitives.Slot(req.Height), parentRoot,
	)

	if err != nil {
		logger.Error("failed to build block", "error", err)
		return nil, err
	}

	// Marshal the block into bytes.
	beaconBz, err := block.MarshalSSZ()
	if err != nil {
		logger.Error("failed to marshal block", "error", err)
	}

	// Run the remainder of the prepare proposal handler.
	resp, err := h.prepareProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	// Inject the beacon kit block into the proposal.
	resp.Txs = append([][]byte{beaconBz}, resp.Txs...)
	return resp, nil
}

// ProcessProposalHandler is a wrapper around the process proposal handler
// that extracts the beacon block from the proposal and processes it.
func (h *Handler) ProcessProposalHandler(
	ctx sdk.Context, req *abci.RequestProcessProposal,
) (*abci.ResponseProcessProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyProcessProposalTime, "ms")
	logger := ctx.Logger().With("module", "process-proposal")

	// Extract the beacon block from the ABCI request.
	//
	// TODO: Block factory struct?
	// TODO: Use protobuf and .(type)?
	block, err := abcitypes.ReadOnlyBeaconKitBlockFromABCIRequest(
		req, h.cfg.BeaconBlockPosition,
		h.chainService.ActiveForkVersionForSlot(primitives.Slot(req.Height)),
	)
	if err != nil {
		return &abci.ResponseProcessProposal{
			Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// TODO abstract this into BeaconState()
	parentRoot := ctx.BlockHeader().AppHash
	if req.Height == 1 {
		parentRoot = make([]byte, 32) //nolint:gomnd //temp
	}

	// TODO: move this to a better spot.
	if !bytes.Equal(parentRoot, block.GetParentRoot()) {
		return &abci.ResponseProcessProposal{
				Status: abci.ResponseProcessProposal_REJECT}, fmt.Errorf(
				"parent root does not match, expected: %x, got: %x",
				ctx.BlockHeader().AppHash, block.GetParentRoot(),
			)
	}

	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	if err = h.chainService.ReceiveBeaconBlock(
		ctx, block,
	); err != nil {
		logger.Error("failed to validate block", "error", err)
		return &abci.ResponseProcessProposal{
			Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// We have to keep a copy of beaconBz to re-inject it into the proposal
	// after the underlying process proposal handler has run. This is to avoid
	// making a
	// copy of the entire request.
	//
	// TODO: there has to be a more friendly way to handle this, but hey it
	// works.
	pos := h.cfg.BeaconBlockPosition
	beaconBz := req.Txs[pos]
	defer func() {
		req.Txs = append([][]byte{beaconBz}, req.Txs...)
	}()
	req.Txs = append(
		req.Txs[:pos], req.Txs[pos+1:]...,
	)

	return h.processProposal(ctx, req)
}
