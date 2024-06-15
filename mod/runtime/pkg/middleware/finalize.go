// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
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
	"encoding/json"
	"sort"
	"time"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
)

// FinalizeBlockMiddleware is a struct that encapsulates the necessary
// components to handle
// the proposal processes.
type FinalizeBlockMiddleware[
	BeaconBlockT interface {
		ssz.Marshallable
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
	},
	BeaconStateT any,
	BlobSidecarsT ssz.Marshallable,
] struct {
	// chainSpec is the chain specification.
	chainSpec primitives.ChainSpec
	// chainService represents the blockchain service.
	chainService BlockchainService[BeaconBlockT, BlobSidecarsT]
	// metrics is the metrics for the middleware.
	metrics *finalizeMiddlewareMetrics
	// valUpdates caches the validator updates as they are produced.
	valUpdates []*transition.ValidatorUpdate
}

// NewFinalizeBlockMiddleware creates a new instance of the Handler struct.
func NewFinalizeBlockMiddleware[
	BeaconBlockT interface {
		ssz.Marshallable
		NewFromSSZ([]byte, uint32) (BeaconBlockT, error)
	},
	BeaconStateT any, BlobSidecarsT ssz.Marshallable,
](
	chainSpec primitives.ChainSpec,
	chainService BlockchainService[BeaconBlockT, BlobSidecarsT],
	telemetrySink TelemetrySink,
) *FinalizeBlockMiddleware[BeaconBlockT, BeaconStateT, BlobSidecarsT] {
	// This is just for nilaway, TODO: remove later.
	if chainService == nil {
		panic("chain service is nil")
	}

	return &FinalizeBlockMiddleware[BeaconBlockT, BeaconStateT, BlobSidecarsT]{
		chainSpec:    chainSpec,
		chainService: chainService,
		metrics:      newFinalizeMiddlewareMetrics(telemetrySink),
	}
}

// InitGenesis is called by the base app to initialize the state of the.
func (h *FinalizeBlockMiddleware[
	BeaconBlockT, BeaconStateT, BlobSidecarsT,
]) InitGenesis(
	ctx context.Context,
	bz []byte,
) ([]appmodulev2.ValidatorUpdate, error) {
	data := new(
		genesis.Genesis[*types.Deposit, *types.ExecutionPayloadHeaderDeneb],
	)
	if err := json.Unmarshal(bz, data); err != nil {
		return nil, err
	}
	updates, err := h.chainService.ProcessGenesisData(
		ctx,
		data,
	)
	if err != nil {
		return nil, err
	}

	// Convert updates into the Cosmos SDK format.
	return iter.MapErr(updates, convertValidatorUpdate)
}

// PreBlock is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *FinalizeBlockMiddleware[
	BeaconBlockT, BeaconStateT, BlobSidecarsT,
]) PreBlock(
	ctx sdk.Context, req *cometabci.FinalizeBlockRequest,
) error {
	startTime := time.Now()
	defer h.metrics.measureEndBlockDuration(startTime)

	blk, blobs, err := encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		h.chainSpec.ActiveForkVersionForSlot(
			math.Slot(req.Height),
		))

	if err != nil {
		// We want to return nil here to prevent the
		// middleware from triggering a panic.
		return nil
	}

	// Process the state transition and produce the required delta from
	// the sync committee.
	h.valUpdates, err = h.chainService.ProcessBlockAndBlobs(
		ctx, blk, blobs,
		// TODO: Speak with @melekes about this, doesn't seem to
		// work reliably.
		/*req.SyncingToHeight == req.Height*/
	)
	return err
}

// EndBlock returns the validator set updates from the beacon state.
func (h FinalizeBlockMiddleware[
	BeaconBlockT, BeaconStateT, BlobSidecarsT,
]) EndBlock(
	context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	// Deduplicate h.valUpdates by pubkey, keeping the later element over any
	// earlier ones
	valUpdatesMap := make(map[string]*transition.ValidatorUpdate)
	for _, update := range h.valUpdates {
		pubKey := string(update.Pubkey[:])
		valUpdatesMap[pubKey] = update
	}

	// Convert map back to slice and sort by pubkey
	dedupedValUpdates := make(
		[]*transition.ValidatorUpdate,
		0,
		len(valUpdatesMap),
	)
	for _, update := range valUpdatesMap {
		dedupedValUpdates = append(dedupedValUpdates, update)
	}
	sort.Slice(dedupedValUpdates, func(i, j int) bool {
		return string(
			dedupedValUpdates[i].Pubkey[:],
		) < string(
			dedupedValUpdates[j].Pubkey[:],
		)
	})
	h.valUpdates = dedupedValUpdates
	return iter.MapErr(h.valUpdates, convertValidatorUpdate)
}
