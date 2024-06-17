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
	// resChannel is used to communicate the validator updates to the
	// EndBlock method.
	valUpdatesChannel chan []*transition.ValidatorUpdate
	// errChannel is used to communicate errors to the EndBlock method.
	errChannel chan error
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
		chainSpec:         chainSpec,
		chainService:      chainService,
		valUpdatesChannel: make(chan []*transition.ValidatorUpdate),
		errChannel:        make(chan error),
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
	// Call the function asynchronously
	go func() {
		blk, blobs, err := encoding.
			ExtractBlobsAndBlockFromRequest[BeaconBlockT, BlobSidecarsT](req,
			BeaconBlockTxIndex,
			BlobSidecarsTxIndex,
			h.chainSpec.ActiveForkVersionForSlot(
				math.Slot(req.Height),
			))

		if err != nil {
			h.errChannel <- err
			return
		}

		result, err := h.chainService.ProcessBlockAndBlobs(ctx, blk, blobs)
		if err != nil {
			h.errChannel <- err
		} else {
			h.valUpdatesChannel <- result
		}
	}()

	return nil
}

// EndBlock returns the validator set updates from the beacon state.
func (h FinalizeBlockMiddleware[
	BeaconBlockT, BeaconStateT, BlobSidecarsT,
]) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	// Deduplicate h.valUpdates by pubkey, keeping the later element over any
	// earlier ones
	select {
	case <-ctx.Done():
		return nil, nil
	case err := <-h.errChannel:
		return nil, err
	case result := <-h.valUpdatesChannel:
		return handleValUpdateConversion(result)
	}
}
