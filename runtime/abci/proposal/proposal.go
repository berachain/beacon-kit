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

package proposal

import (
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/prysmaticlabs/prysm/v4/consensus-types/primitives"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	consensusv1 "github.com/itsdevbear/bolaris/types/consensus/v1"
)

// TODO: Need to have the wait for syncing phase at the start to allow the Execution Client
// to sync up and the consensus client shouldn't join the validator set yet.
// TODO: also need to make payload position a config variable or something.
const PayloadPosition = 0

// Handler is a struct that encapsulates the necessary components to handle the proposal processes.
type Handler struct {
	prepareProposal sdk.PrepareProposalHandler
	processProposal sdk.ProcessProposalHandler
	beaconChain     *blockchain.Service
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(
	beaconChain *blockchain.Service,
	prepareProposal sdk.PrepareProposalHandler,
	processProposal sdk.ProcessProposalHandler,
) *Handler {
	return &Handler{
		beaconChain:     beaconChain,
		prepareProposal: prepareProposal,
		processProposal: processProposal,
	}
}

func (h *Handler) PrepareProposalHandler(
	ctx sdk.Context, req *abci.RequestPrepareProposal,
) (*abci.ResponsePrepareProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyPrepareProposalTime, "ms")
	logger := ctx.Logger().With("module", "prepare-proposal")

	// We start by requesting a block from the execution client. This may be from pulling
	// a previously built payload from the local cache via `getPayload()` or it may be
	// by asking for a forkchoice from the execution client, depending on timing.
	block, err := h.beaconChain.GetOrBuildBlock(
		ctx, primitives.Slot(req.Height),
	)

	if err != nil {
		logger.Error("failed to build block", "error", err)
		return nil, err
	}

	// Marshal the block into bytes.
	bz, err := block.Marshal()
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

func (h *Handler) ProcessProposalHandler(
	ctx sdk.Context, req *abci.RequestProcessProposal,
) (*abci.ResponseProcessProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyProcessProposalTime, "ms")
	logger := ctx.Logger().With("module", "process-proposal")

	// Extract the beacon kit block from the proposal and unmarshal it.
	block, err := consensusv1.ReadOnlyBeaconKitBlockFromABCIRequest(
		req, PayloadPosition,
	)
	if err != nil {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// If we get any sort of error from the execution client, we bubble
	// it up and reject the proposal, as we do not want to write a block
	// finalization to the consensus layer that is invalid.
	if err = h.beaconChain.ProcessReceivedBlock(
		ctx, block,
	); err != nil {
		logger.Error("failed to validate block", "error", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// Run the remainder of the proposal. We remove the beacon block from the proposal
	// before passing it to the next handler.
	return h.processProposal(ctx, h.removeBeaconBlockFromTxs(req))
}

// removeBeaconBlockFromTxs removes the beacon block from the proposal.
// TODO: optimize this function to avoid the giga memory copy.
func (h *Handler) removeBeaconBlockFromTxs(
	req *abci.RequestProcessProposal,
) *abci.RequestProcessProposal {
	// Extract and remove the PayloadPosition'th tx from the proposal.
	txsLen := len(req.Txs)
	switch PayloadPosition {
	case 0: // Remove the first element
		req.Txs = req.Txs[1:]
	case txsLen - 1: // Remove the last element
		req.Txs = req.Txs[:txsLen-1]
	default: // Remove an element from the middle
		// Shift elements to the left to overwrite the element at PayloadPosition
		copy(req.Txs[PayloadPosition:], req.Txs[PayloadPosition+1:])
		req.Txs = req.Txs[:txsLen-1] // Slice off the last element which is now duplicated
	}
	return req
}
