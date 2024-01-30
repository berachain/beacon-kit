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
	"errors"
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
		ctx, primitives.Slot(req.Height), uint64(req.Time.UTC().Unix()),
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
	block, err := h.extractAndUnmarshalBlock(req)
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

	// Run the remainder of the proposal.
	return h.processProposal(ctx, req)
}

// extractAndUnmarshalBlock extracts the block from the proposal and unmarshals it.
func (h *Handler) extractAndUnmarshalBlock(
	req *abci.RequestProcessProposal,
) (*consensusv1.BaseBeaconKitBlock, error) {
	// Extract the marshalled payload from the proposal
	if len(req.Txs) == 0 {
		return nil, errors.New("no transactions in proposal")
	}
	bz := req.Txs[PayloadPosition]
	req.Txs = req.Txs[1:]

	block := &consensusv1.BaseBeaconKitBlock{}
	if err := block.Unmarshal(bz); err != nil {
		return nil, err
	}
	return block, nil
}
