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
	"github.com/prysmaticlabs/prysm/v4/math"

	abci "github.com/cometbft/cometbft/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/itsdevbear/bolaris/beacon/blockchain"
	consensus_types "github.com/itsdevbear/bolaris/beacon/consensus-types"
	block_sync "github.com/itsdevbear/bolaris/beacon/execution/block-sync"
)

// TODO: Need to have the wait for syncing phase at the start to allow the Execution Client
// to sync up and the consensus client shouldn't join the validator set yet.
const PayloadPosition = 0

// Handler is a struct that handles the proposal process.
type Handler struct {
	prepareProposal sdk.PrepareProposalHandler
	processProposal sdk.ProcessProposalHandler
	beaconChain     *blockchain.Service
	blockSync       *block_sync.BlockSync
}

func NewHandler(
	beaconChain *blockchain.Service,
	prepareProposal sdk.PrepareProposalHandler,
	processProposal sdk.ProcessProposalHandler,
	blockSync *block_sync.BlockSync,
) *Handler {
	return &Handler{
		beaconChain:     beaconChain,
		prepareProposal: prepareProposal,
		processProposal: processProposal,
		blockSync:       blockSync,
	}
}

func (h *Handler) PrepareProposalHandler(
	ctx sdk.Context, req *abci.RequestPrepareProposal,
) (*abci.ResponsePrepareProposal, error) {
	logger := ctx.Logger().With("module", "prepare-proposal")
	// Build the block on the execution layer.
	// TODO: manage the different type of engine API errors.
	payload, err := h.beaconChain.BuildNextBlock(ctx, ctx.HeaderInfo() /* slot */)
	if err != nil {
		logger.Error("failed to build block", "err", err)
		return nil, err
	}

	// Run the remainder of the prepare proposal handler.
	resp, err := h.prepareProposal(ctx, req)
	if err != nil {
		return nil, err
	}

	// Marshal the payload.
	bz, err := payload.MarshalSSZ()
	if err != nil {
		return nil, err
	}

	// Inject the payload into the proposal.
	resp.Txs = append([][]byte{bz}, resp.Txs...)
	return resp, nil
}

func (h *Handler) ProcessProposalHandler(
	ctx sdk.Context, req *abci.RequestProcessProposal,
) (*abci.ResponseProcessProposal, error) {
	logger := ctx.Logger().With("module", "process-proposal")

	// Extract the marshalled payload from the proposal
	if len(req.Txs) == 0 {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
	}
	bz := req.Txs[PayloadPosition]
	req.Txs = req.Txs[1:]

	// Unmarshal the payload.
	data, err := consensus_types.BytesToExecutionData(bz, math.Gwei(0), 3)
	if err != nil {
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, err
	}

	if _, err = h.beaconChain.ProcessExecutionData(ctx, ctx.HeaderInfo(), data); err != nil {
		logger.Error("failed to validate block", "err", err)
		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, err
	}

	// Run the remainder of the proposal.
	return h.processProposal(ctx, req)
}
