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
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/itsdevbear/bolaris/types/consensus/primitives"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	sync "github.com/itsdevbear/bolaris/beacon/sync"
	"github.com/itsdevbear/bolaris/config"
	abcitypes "github.com/itsdevbear/bolaris/runtime/abci/types"
	"github.com/itsdevbear/bolaris/validator"
)

// Handler is a struct that encapsulates the necessary components to handle
// the proposal processes.
type Handler struct {
	cfg             *config.ABCI
	validator       *validator.Service
	beaconChain     *blockchain.Service
	syncService     *sync.Service
	prepareProposal sdk.PrepareProposalHandler
	processProposal sdk.ProcessProposalHandler
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(
	cfg *config.ABCI,
	validator *validator.Service,
	syncService *sync.Service,
	beaconChain *blockchain.Service,
	prepareProposal sdk.PrepareProposalHandler,
	processProposal sdk.ProcessProposalHandler,
) *Handler {
	return &Handler{
		cfg:             cfg,
		validator:       validator,
		syncService:     syncService,
		beaconChain:     beaconChain,
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
	if bsp := h.syncService.CheckSyncStatus(ctx); bsp.Status == sync.StatusExecutionAhead {
		return nil, fmt.Errorf("err: %w, status: %d", ErrValidatorClientNotSynced, bsp.Status)
	}

	// We start by requesting the validator service to build us a block. This may
	// be from pulling a previously built payload from the local cache or it may be
	// by asking for a forkchoice from the execution client, depending on timing.
	block, err := h.validator.BuildBeaconBlock(
		ctx, primitives.Slot(req.Height),
	)

	if err != nil {
		logger.Error("failed to build block", "error", err)
		return nil, err
	}

	// Marshal the block into bytes.
	bz, err := block.MarshalSSZ()
	if err != nil {
		logger.Error("failed to marshal block", "error", err)
	}

	// Run the remainder of the prepare proposal handler.
	resp, err := h.prepareProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	// Inject the beacon kit block into the proposal.
	resp.Txs = append([][]byte{bz}, resp.Txs...)
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
	// TODO: I don't love how we have to use the BeaconConfig here.
	// TODO: Block factory struct?
	// TODO: Use protobuf and .(type)?
	block, err := abcitypes.ReadOnlyBeaconKitBlockFromABCIRequest(
		req, h.cfg.BeaconBlockPosition,
		h.beaconChain.BeaconCfg().ActiveForkVersion(primitives.Epoch(req.Height)),
	)
	if err != nil {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	if err = h.beaconChain.ReceiveBeaconBlock(
		ctx, block,
	); err != nil {
		logger.Error("failed to validate block", "error", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// Run the remainder of the proposal. We remove the beacon block from the proposal
	// before passing it to the next handler.
	return h.processProposal(ctx, h.RemoveBeaconBlockFromTxs(req))
}

// removeBeaconBlockFromTxs removes the beacon block from the proposal.
// TODO: optimize this function to avoid the giga memory copy.
func (h *Handler) RemoveBeaconBlockFromTxs(
	req *abci.RequestProcessProposal,
) *abci.RequestProcessProposal {
	req.Txs = removeAtIndex(req.Txs, h.cfg.BeaconBlockPosition)
	return req
}

// removeAtIndex removes an element at a given index from a slice.
func removeAtIndex[T any](s []T, index uint) []T {
	return append(s[:index], s[index+1:]...)
}
