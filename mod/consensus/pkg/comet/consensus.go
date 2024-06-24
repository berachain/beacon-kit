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

package comet

import (
	"context"

	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/transition"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
)

// TODO: We must rid this of comet bft types.
type Middleware interface {
	InitGenesis(
		ctx context.Context,
		bz []byte,
	) (transition.ValidatorUpdates, error)

	PrepareProposal(
		context.Context,
		math.Slot,
	) ([]byte, []byte, error)

	ProcessProposal(
		ctx context.Context,
		req *cmtabci.ProcessProposalRequest,
	) (*cmtabci.ProcessProposalResponse, error)

	PreBlock(
		_ context.Context, req *cmtabci.FinalizeBlockRequest,
	) error

	EndBlock(
		ctx context.Context,
	) (transition.ValidatorUpdates, error)
}

// NewConsensus returns a new consensus middleware.
func NewConsensus(m Middleware) *Consensus {
	return &Consensus{
		Middleware: m,
	}
}

// Consensus is used to decouple the Comet consensus engine from the Cosmos SDK.
// Right now, it is very coupled to the sdk base app and we will
// eventually fully decouple this.
type Consensus struct {
	Middleware
}

func (c *Consensus) InitGenesis(
	ctx context.Context,
	bz []byte,
) ([]appmodulev2.ValidatorUpdate, error) {
	updates, err := c.Middleware.InitGenesis(ctx, bz)
	if err != nil {
		return nil, err
	}
	// Convert updates into the Cosmos SDK format.
	return iter.MapErr(updates, convertValidatorUpdate)
}

// TODO: Decouple Comet Types
func (c *Consensus) PrepareProposal(
	ctx sdk.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	slot := math.Slot(req.Height)
	blkBz, sidecarsBz, err := c.Middleware.PrepareProposal(ctx, slot)
	if err != nil {
		return nil, err
	}
	return &cmtabci.PrepareProposalResponse{
		Txs: [][]byte{blkBz, sidecarsBz},
	}, nil
}

// TODO: Decouple Comet Types
func (c *Consensus) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	return c.Middleware.ProcessProposal(ctx, req)
}

// TODO: Decouple Comet Types
func (c *Consensus) PreBlock(
	ctx sdk.Context, req *cmtabci.FinalizeBlockRequest,
) error {
	return c.Middleware.PreBlock(ctx, req)
}

func (c *Consensus) EndBlock(
	ctx context.Context,
) ([]appmodulev2.ValidatorUpdate, error) {
	updates, err := c.Middleware.EndBlock(ctx)
	if err != nil {
		return nil, err
	}
	return iter.MapErr(updates, convertValidatorUpdate)
}
