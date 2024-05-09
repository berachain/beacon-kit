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
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	engineclient "github.com/berachain/beacon-kit/mod/execution/pkg/client"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/sync/errgroup"
)

// Handler is a struct that encapsulates the necessary components to handle
// the proposal processes.
type Handler[BlobsSidecarsT ssz.Marshallable] struct {
	builderService BuilderService[BlobsSidecarsT]
	chainService   BlockchainService[BlobsSidecarsT]

	// TODO: we will eventually gossip the blobs separately from
	// CometBFT, but for now, these are no-op gossipers.
	blobGossiper        p2p.Publisher[BlobsSidecarsT, []byte]
	beaconBlockGossiper p2p.PublisherReceiver[
		types.BeaconBlock, []byte, encoding.ABCIRequest, types.BeaconBlock,
	]
}

// NewHandler creates a new instance of the Handler struct.
func NewHandler[BlobsSidecarsT ssz.Marshallable](
	builderService BuilderService[BlobsSidecarsT],
	chainService BlockchainService[BlobsSidecarsT],
) *Handler[BlobsSidecarsT] {
	// This is just for nilaway, TODO: remove later.
	if chainService == nil {
		panic("chain service is nil")
	}
	return &Handler[BlobsSidecarsT]{
		builderService: builderService,
		chainService:   chainService,
		// TODO: we will eventually gossipt the blobs separately from
		// CometBFT.
		blobGossiper: rp2p.NoopGossipHandler[BlobsSidecarsT, []byte]{},
		beaconBlockGossiper: rp2p.NewNoopBlockGossipHandler[encoding.ABCIRequest](
			chainService.ChainSpec(),
		),
	}
}

// PrepareProposalHandler is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *Handler[BlobsSidecarsT]) PrepareProposalHandler(
	ctx sdk.Context, req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyPrepareProposalTime, "ms")

	var (
		logger = ctx.Logger().With("module", "prepare-proposal")
		st     = h.chainService.BeaconState(ctx)
	)

	// Process the Slot to set the state root for the block.
	if err := h.chainService.ProcessSlot(st); err != nil {
		return &cmtabci.PrepareProposalResponse{}, err
	}

	// Get the best block and blobs.
	blk, blobs, err := h.builderService.RequestBestBlock(
		ctx, st, math.Slot(req.Height))
	if err != nil || blk == nil || blk.IsNil() {
		logger.Error("failed to build block", "error", err, "block", blk)
		return &cmtabci.PrepareProposalResponse{}, err
	}

	// "Publish" the blobs and the beacon block.
	var sidecarsBz, beaconBlockBz []byte
	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		var localErr error
		sidecarsBz, localErr = h.blobGossiper.Publish(gCtx, blobs)
		if localErr != nil {
			logger.Error("failed to publish blobs", "error", err)
		}
		return err
	})

	g.Go(func() error {
		var localErr error
		beaconBlockBz, localErr = h.beaconBlockGossiper.Publish(gCtx, blk)
		if localErr != nil {
			logger.Error("failed to publish beacon block", "error", err)
		}
		return err
	})

	return &cmtabci.PrepareProposalResponse{
		Txs: [][]byte{beaconBlockBz, sidecarsBz},
	}, g.Wait()
}

// ProcessProposalHandler is a wrapper around the process proposal handler
// that extracts the beacon block from the proposal and processes it.
func (h *Handler[BlobsSidecarsT]) ProcessProposalHandler(
	ctx sdk.Context, req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	defer telemetry.MeasureSince(time.Now(), MetricKeyProcessProposalTime, "ms")

	var (
		logger = ctx.Logger().With("module", "process-proposal")
		blk    types.BeaconBlock
		err    error
	)

	if blk, err = h.beaconBlockGossiper.Request(ctx, req); err != nil {
		logger.Error(
			"failed to retrieve beacon block from request",
			"error",
			err,
		)
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, err
	}

	// If the block is syncing, we reject the proposal. This is guard against a
	// potential attack under the unlikely scenario in which a supermajority of
	// validators have their EL's syncing. If nodes were to accept this proposal
	// optmistically when they are syncing, it could potentially allow for a
	// malicious validator to push a bad block through.
	if err = h.chainService.VerifyPayloadOnBlk(ctx, blk); errors.Is(
		err,
		engineclient.ErrSyncingPayloadStatus,
	) {
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, err
	}

	return &cmtabci.ProcessProposalResponse{
		Status: cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT,
	}, nil
}
