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

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/genesis"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/encoding"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
)

// InitGenesis is called by the base app to initialize the state of the.
func (h *ABCIMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
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
func (h *ABCIMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) PreBlock(
	ctx sdk.Context, req *cometabci.FinalizeBlockRequest,
) error {
	go h.preBlock(ctx, req)
	return nil
}

// handlePreBlock is called by the base app before the block is finalized. It
// is responsible for aggregating oracle data from each validator and writing
// the oracle data to the store.
func (h *ABCIMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) preBlock(
	ctx sdk.Context, req *cometabci.FinalizeBlockRequest,
) {
	blk, blobs, err := encoding.
		ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](req,
		BeaconBlockTxIndex,
		BlobSidecarsTxIndex,
		h.chainSpec.ActiveForkVersionForSlot(
			math.Slot(req.Height),
		))

	if err != nil {
		h.errChannel <- errors.Join(err, ErrBadExtractBlockAndBlocks)
		return
	}

	result, err := h.chainService.ProcessBlockAndBlobs(ctx, blk, blobs)
	if err != nil {
		h.errChannel <- err
	} else {
		h.valUpdatesChannel <- result
	}
}

// EndBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	return h.endBlock(ctx)
}

// endBlock returns the validator set updates from the beacon state.
func (h *ABCIMiddleware[
	AvailabilityStoreT,
	BeaconBlockT,
	BeaconBlockBodyT,
	BeaconStateT,
	BlobSidecarsT,
]) endBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-h.errChannel:
		if errors.Is(err, ErrBadExtractBlockAndBlocks) {
			err = nil
		}
		return nil, err
	case result := <-h.valUpdatesChannel:
		return iter.MapErr(
			result.RemoveDuplicates().Sort(), convertValidatorUpdate,
		)
	}
}
