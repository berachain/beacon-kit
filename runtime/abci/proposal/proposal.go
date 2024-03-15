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
	"time"

	"github.com/berachain/beacon-kit/beacon/blockchain"
	"github.com/berachain/beacon-kit/beacon/builder"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/health"
	"github.com/berachain/beacon-kit/primitives"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Handler is a struct that encapsulates the necessary components to handle
// the proposal processes.
type Handler struct {
	cfg *config.ABCI
	// stakingKeeper provides access to the staking module. In the handler
	// it is used to convert consAddress to pubkey, before passing it into
	// the core beacon chain logic.
	stakingKeeper StakingKeeper

	builderService *builder.Service
	chainService   *blockchain.Service
	healthService  *health.Service
	nextPrepare    sdk.PrepareProposalHandler
	nextProcess    sdk.ProcessProposalHandler
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler(
	cfg *config.ABCI,
	stakingKeeper StakingKeeper,
	builderService *builder.Service,
	healthService *health.Service,
	chainService *blockchain.Service,
	nextPrepare sdk.PrepareProposalHandler,
	nextProcess sdk.ProcessProposalHandler,

) *Handler {
	return &Handler{
		cfg:            cfg,
		stakingKeeper:  stakingKeeper,
		builderService: builderService,
		healthService:  healthService,
		chainService:   chainService,
		nextPrepare:    nextPrepare,
		nextProcess:    nextProcess,
	}
}

// PrepareProposalHandler is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *Handler) PrepareProposalHandler(
	ctx sdk.Context, req *abci.RequestPrepareProposal,
) (*abci.ResponsePrepareProposal, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyPrepareProposalTime, "ms")
	logger := ctx.Logger().With("module", "prepare-proposal")

	proposerPubkey, err := h.stakingKeeper.GetValidatorPubkeyFromConsAddress(
		ctx, req.ProposerAddress,
	)
	if err != nil {
		return &abci.ResponsePrepareProposal{}, err
	}

	// We start by requesting the validator service to build us a block. This
	// may be from pulling a previously built payload from the local cache or it
	// may be by asking for a forkchoice from the execution client, depending on
	// timing.
	blk, err := h.builderService.RequestBestBlock(
		ctx,
		primitives.Slot(req.Height),
		proposerPubkey,
	)
	if err != nil || blk == nil || blk.IsNil() {
		logger.Error("failed to build block", "error", err, "block", blk)
		return &abci.ResponsePrepareProposal{}, err
	}

	// Marshal the block into bytes.
	beaconBz, err := blk.MarshalSSZ()
	if err != nil {
		logger.Error("failed to marshal block", "error", err)
	}

	// Run the remainder of the prepare proposal handler.
	resp, err := h.nextPrepare(ctx, req)
	if err != nil {
		return nil, err
	}

	// If the response is nil, the implementations of
	// `nextPrepare` is bad.
	if resp == nil {
		return nil, ErrNextPrepareNilResp
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

	return h.nextProcess(ctx, req)
}
