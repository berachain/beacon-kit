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
	"sort"
	"time"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/p2p"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	rp2p "github.com/berachain/beacon-kit/mod/runtime/pkg/p2p"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	v1 "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	BlobSidecarsT ssz.Marshallable,
	StorageBackendT any,
] struct {
	// chainSpec is the chain specification.
	chainSpec primitives.ChainSpec
	// validatorService is the service responsible for building beacon blocks.
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobSidecarsT,
	]

	chainService BlockchainService[BeaconBlockT, BlobSidecarsT]

	// TODO: we will eventually gossip the blobs separately from
	// CometBFT, but for now, these are no-op gossipers.
	blobGossiper p2p.PublisherReceiver[
		BlobSidecarsT,
		[]byte,
		encoding.ABCIRequest,
		BlobSidecarsT,
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
	BlobSidecarsT ssz.Marshallable,
	StorageBackendT StorageBackend[BeaconStateT],
](
	chainSpec primitives.ChainSpec,
	validatorService ValidatorService[
		BeaconBlockT,
		BeaconStateT,
		BlobSidecarsT,
	],
	chainService BlockchainService[BeaconBlockT, BlobSidecarsT],
	telemetrySink TelemetrySink,
	storageBackend StorageBackendT,
) *ValidatorMiddleware[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
	BeaconStateT, BlobSidecarsT, StorageBackendT,
] {
	return &ValidatorMiddleware[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT,
		BeaconStateT, BlobSidecarsT, StorageBackendT,
	]{
		chainSpec:        chainSpec,
		validatorService: validatorService,
		chainService:     chainService,
		blobGossiper: rp2p.NewNoopBlobHandler[
			BlobSidecarsT, encoding.ABCIRequest](),
		beaconBlockGossiper: rp2p.
			NewNoopBlockGossipHandler[BeaconBlockT, encoding.ABCIRequest](
			chainSpec,
		),
		metrics:        newValidatorMiddlewareMetrics(telemetrySink),
		storageBackend: storageBackend,
	}
}

// PrepareProposalHandler is a wrapper around the prepare proposal handler
// that injects the beacon block into the proposal.
func (h *ValidatorMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) PrepareProposalHandler(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	var (
		startTime     = time.Now()
		sidecarsBz    []byte
		beaconBlockBz []byte
		logger        = ctx.Logger().With(
			"service", "prepare-proposal",
		)
	)
	defer h.metrics.measurePrepareProposalDuration(startTime)

	// Get the best block and blobs.
	blk, blobs, err := h.validatorService.RequestBestBlock(
		ctx, math.Slot(req.GetHeight()))
	if err != nil || blk.IsNil() {
		logger.Error(
			"failed to assemble proposal", "error", err, "block", blk)
		return &cmtabci.PrepareProposalResponse{}, err
	}

	st := h.storageBackend.StateFromContext(ctx)

	// Get the previous slot's hash tree root.
	root := blk.GetParentBlockRoot()

	// Get the attestations from the votes.
	attestations, err := h.attestationDataFromVotes(
		st,
		root,
		req.LocalLastCommit.Votes,
		//#nosec:G701 // safe.
		uint64(req.Height-1),
	)
	if err != nil {
		logger.Error("failed to get attestations from votes", "error", err)
		return &cmtabci.PrepareProposalResponse{}, err
	}
	blk.GetBody().SetAttestations(attestations)

	// Get the slashing info from the misbehaviors.
	slashingInfo, err := h.slashingInfoFromMisbehaviors(st, req.Misbehavior)
	blk.GetBody().SetSlashingInfo(slashingInfo)

	// "Publish" the blobs and the beacon block.
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

// attestationDataFromVotes returns a list of attestation data from the
// comet vote info. This is used to build the attestations for the block.
func (h *ValidatorMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobsSidecarsT,
	DepositStoreT,
]) attestationDataFromVotes(
	st BeaconStateT,
	root primitives.Root,
	votes []v1.ExtendedVoteInfo,
	slot uint64,
) ([]*types.AttestationData, error) {
	var err error
	var index math.U64
	attestations := make([]*types.AttestationData, len(votes))
	for i, vote := range votes {
		index, err = st.ValidatorIndexByCometBFTAddress(vote.Validator.Address)
		if err != nil {
			return nil, err
		}
		attestations[i] = &types.AttestationData{
			Slot:            slot,
			Index:           index.Unwrap(),
			BeaconBlockRoot: root,
		}
	}
	// Attestations are sorted by index.
	sort.Slice(attestations, func(i, j int) bool {
		return attestations[i].Index < attestations[j].Index
	})
	return attestations, nil
}

// slashingInfoFromMisbehaviors returns a list of slashing info from the
// comet misbehaviors.
func (h *ValidatorMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobsSidecarsT,
	DepositStoreT,
]) slashingInfoFromMisbehaviors(
	st BeaconStateT,
	misbehaviors []v1.Misbehavior,
) ([]*types.SlashingInfo, error) {
	var err error
	var index math.U64
	slashingInfo := make([]*types.SlashingInfo, len(misbehaviors))
	for i, misbehavior := range misbehaviors {
		index, err = st.ValidatorIndexByCometBFTAddress(
			misbehavior.Validator.Address,
		)
		if err != nil {
			return nil, err
		}
		slashingInfo[i] = &types.SlashingInfo{
			//#nosec:G701 // safe.
			Slot:  uint64(misbehavior.GetHeight()),
			Index: index.Unwrap(),
		}
	}
	return slashingInfo, nil
}

// ProcessProposalHandler is a wrapper around the process proposal handler
// that extracts the beacon block from the proposal and processes it.
func (h *ValidatorMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
	DepositStoreT,
]) ProcessProposalHandler(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	var (
		startTime = time.Now()
		logger    = ctx.Logger().With(
			"service", "prepare-proposal",
		)
	)
	defer h.metrics.measureProcessProposalDuration(startTime)

	args := []any{"beacon_block", true, "blob_sidecars", true}
	blk, err := h.beaconBlockGossiper.Request(ctx, req)
	if err != nil {
		args[1] = false
	}

	sidecars, err := h.blobGossiper.Request(ctx, req)
	if err != nil {
		args[3] = false
	}

	logger.Info("received proposal with", args...)
	if err = h.chainService.ReceiveBlockAndBlobs(
		ctx, blk, sidecars,
	); errors.IsFatal(err) {
		return &cmtabci.ProcessProposalResponse{
			Status: cmtabci.PROCESS_PROPOSAL_STATUS_REJECT,
		}, err
	}

	return &cmtabci.ProcessProposalResponse{
		Status: cmtabci.PROCESS_PROPOSAL_STATUS_ACCEPT,
	}, nil
}
