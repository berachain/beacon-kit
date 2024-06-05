// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package middleware

import (
	"context"
	"time"

	"cosmossdk.io/server/v2/cometbft/handlers"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/gogoproto/proto"
	"golang.org/x/sync/errgroup"
)

// ValidatorMiddleware is a middleware between ABCI and the validator logic.
type ValidatorMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT interface {
		ValidatorIndexByPubkey(pk crypto.BLSPubkey) (math.ValidatorIndex, error)
		GetBlockRootAtIndex(slot uint64) (primitives.Root, error)
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
	},
	BlobsSidecarsT ssz.Marshallable,
	StorageBackendT any,
	TransactionT types.Tx,
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

	// storageBackend is the storage backend.
	storageBackend StorageBackend[BeaconStateT]
}

// NewValidatorMiddleware creates a new instance of the Handler struct.
func NewValidatorMiddleware[
	AvailabilityStoreT any,
	BeaconBlockT interface {
		types.RawBeaconBlock[BeaconBlockBodyT]
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
		NewWithVersion(
			math.Slot,
			math.ValidatorIndex,
			primitives.Root,
			uint32,
		) (BeaconBlockT, error)
		Empty(uint32) BeaconBlockT
	},
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT interface {
		ValidatorIndexByPubkey(pk crypto.BLSPubkey) (math.ValidatorIndex, error)
		GetBlockRootAtIndex(slot uint64) (primitives.Root, error)
		ValidatorIndexByCometBFTAddress(
			cometBFTAddress []byte,
		) (math.ValidatorIndex, error)
	},
	BlobsSidecarsT ssz.Marshallable,
	StorageBackendT StorageBackend[BeaconStateT],
	TransactionT types.Tx,
](
	chainSpec primitives.ChainSpec,
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobsSidecarsT,
	],
	telemetrySink TelemetrySink,
	storageBackend StorageBackendT,
) *ValidatorMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconStateT, BlobsSidecarsT, StorageBackendT, TransactionT,
] {
	return &ValidatorMiddleware[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobsSidecarsT, StorageBackendT, TransactionT,
	]{
		chainSpec:        chainSpec,
		validatorService: validatorService,
		blobGossiper: rp2p.
			NoopGossipHandler[BlobsSidecarsT, []byte]{},
		beaconBlockGossiper: rp2p.
			NewNoopBlockGossipHandler[BeaconBlockT, encoding.ABCIRequest](
			chainSpec,
		),
		metrics:        newValidatorMiddlewareMetrics(telemetrySink),
		storageBackend: storageBackend,
	}
}

func (h *ValidatorMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobsSidecarsT,
	DepositStoreT,
	TransactionT,
]) PrepareProposalHandler() handlers.PrepareHandler[TransactionT] {
	return func(
		ctx context.Context,
		app handlers.AppManager[TransactionT],
		txs []TransactionT,
		req proto.Message,
	) ([]TransactionT, error) {
		var (
			// logger    = ctx.Logger().With("service", "prepare-proposal")
			startTime = time.Now()
		)

		abciReq, ok := req.(*cmtabci.PrepareProposalRequest)
		if !ok {
			return nil, ErrInvalidRequestType
		}

		defer h.metrics.measurePrepareProposalDuration(startTime)

		// Get the best block and blobs.
		blk, blobs, err := h.validatorService.RequestBestBlock(
			ctx, math.Slot(abciReq.GetHeight()))
		if err != nil || blk.IsNil() {
			// logger.Error("failed to build block", "error", err, "block", blk)
			return nil, err
		}

		// "Publish" the blobs and the beacon block.
		var sidecarsBz, beaconBlockBz []byte
		g, gCtx := errgroup.WithContext(ctx)
		g.Go(func() error {
			sidecarsBz, err = h.blobGossiper.Publish(gCtx, blobs)
			if err != nil {
				// logger.Error("failed to publish blobs", "error", err)
			}
			return err
		})

		g.Go(func() error {
			beaconBlockBz, err = h.beaconBlockGossiper.Publish(gCtx, blk)
			if err != nil {
				// logger.Error("failed to publish beacon block", "error", err)
			}
			return err
		})

		return []TransactionT{
			types.NewTx[TransactionT](beaconBlockBz),
			types.NewTx[TransactionT](sidecarsBz),
		}, g.Wait()
	}
}

func (h *ValidatorMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobsSidecarsT,
	DepositStoreT,
	TransactionT,
]) ProcessProposalHandler() handlers.ProcessHandler[TransactionT] {
	return func(
		ctx context.Context,
		app handlers.AppManager[TransactionT],
		txs []TransactionT,
		req proto.Message,
	) error {
		abciReq, ok := req.(*cmtabci.PrepareProposalRequest)
		if !ok {
			return ErrInvalidRequestType
		}

		startTime := time.Now()
		defer h.metrics.measureProcessProposalDuration(startTime)

		//#nosec:G703
		blk, err := h.beaconBlockGossiper.Request(ctx, abciReq)
		if err != nil {
			// TODO: Handle better.
			blk = blk.Empty(version.Deneb)
		}
		if err = h.validatorService.VerifyIncomingBlock(ctx, blk); err != nil {
			return err
		}
		return nil
	}
}
