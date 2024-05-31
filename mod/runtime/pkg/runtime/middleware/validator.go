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

package middleware

import (
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"golang.org/x/sync/errgroup"
)

// ValidatorMiddleware is a middleware between ABCI and the validator logic.
type ValidatorMiddleware[
	BeaconBlockT interface {
		types.RawBeaconBlock
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
	},
	BeaconStateT any, BlobsSidecarsT ssz.Marshallable,
] struct {
	// chainSpec is the chain specification.
	chainSpec primitives.ChainSpec
	// validatorService is the service responsible for building beacon blocks.
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobsSidecarsT,
	]
	// TODO: we will eventually gossip the blobs separately from
	// CometBFT, but for now, these are no-op gossipers.
	blobGossiper p2p.Publisher[
		BlobsSidecarsT, []byte,
	]
	// TODO: we will eventually gossip the blocks separately from
	// CometBFT, but for now, these are no-op gossipers.
	beaconBlockGossiper p2p.PublisherReceiver[
		BeaconBlockT,
		[]byte,
		encoding.ABCIRequest,
		BeaconBlockT,
	]
	// metrics is the metrics emitter.
	metrics *validatorMiddlewareMetrics
}

// NewValidatorMiddleware creates a new instance of the Handler struct.
func NewValidatorMiddleware[
	BeaconBlockT interface {
		types.RawBeaconBlock
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
	},
	BeaconStateT any,
	BlobsSidecarsT ssz.Marshallable,
](
	chainSpec primitives.ChainSpec,
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobsSidecarsT,
	],
	telemetrySink TelemetrySink,
) *ValidatorMiddleware[BeaconBlockT, BeaconStateT, BlobsSidecarsT] {
	return &ValidatorMiddleware[BeaconBlockT, BeaconStateT, BlobsSidecarsT]{
		chainSpec:        chainSpec,
		validatorService: validatorService,
		blobGossiper: rp2p.
			NoopGossipHandler[BlobsSidecarsT, []byte]{},
		beaconBlockGossiper: rp2p.
			NewNoopBlockGossipHandler[BeaconBlockT, encoding.ABCIRequest](
			chainSpec,
		),
		metrics: newValidatorMiddlewareMetrics(telemetrySink),
	}
}

// PrepareProposalHandler is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *ValidatorMiddleware[
	BeaconBlockT, BeaconStateT, BlobsSidecarsT,
]) PrepareProposalHandler(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	var (
		logger    = ctx.Logger().With("service", "prepare-proposal")
		startTime = time.Now()
	)

	defer h.metrics.measurePrepareProposalDuration(startTime)

	// Get the best block and blobs.
	blk, blobs, err := h.validatorService.RequestBestBlock(
		ctx, math.Slot(req.GetHeight()))
	if err != nil || blk.IsNil() {
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
func (h *ValidatorMiddleware[
	BeaconBlockT, BeaconStateT, BlobsSidecarsT,
]) ProcessProposalHandler(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	var (
		logger    = ctx.Logger().With("service", "process-proposal")
		startTime = time.Now()
	)

	defer h.metrics.measureProcessProposalDuration(startTime)
	if blk, err := h.beaconBlockGossiper.Request(ctx, req); err != nil {
		logger.Error(
			"failed to retrieve beacon block from request",
			"error",
			err,
		)
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, err
	} else if err = h.validatorService.
		VerifyIncomingBlock(ctx, blk); err != nil {
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, err
	}

	return &cmtabci.ProcessProposalResponse{
		Status: cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT,
	}, nil
}
