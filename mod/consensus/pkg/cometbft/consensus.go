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

package cometbft

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sourcegraph/conc/iter"
)

// ConsensusEngine is used to decouple the Comet consensus engine from
// the Cosmos SDK.
// Right now, it is very coupled to the sdk base app and we will
// eventually fully decouple this.
type ConsensusEngine[ValidatorUpdateT any] struct {
	Middleware
}

// NewConsensusEngine returns a new consensus middleware.
func NewConsensusEngine[ValidatorUpdateT any](
	m Middleware,
) *ConsensusEngine[ValidatorUpdateT] {
	return &ConsensusEngine[ValidatorUpdateT]{
		Middleware: m,
	}
}

func (c *ConsensusEngine[ValidatorUpdateT]) InitGenesis(
	ctx context.Context,
	genesisBz []byte,
) ([]ValidatorUpdateT, error) {
	updates, err := c.Middleware.InitGenesis(ctx, genesisBz)
	if err != nil {
		return nil, err
	}
	// Convert updates into the Cosmos SDK format.
	return iter.MapErr(updates, convertValidatorUpdate[ValidatorUpdateT])
}

// TODO: Decouple Comet Types
func (c *ConsensusEngine[ValidatorUpdateT]) PrepareProposal(
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
func (c *ConsensusEngine[ValidatorUpdateT]) ProcessProposal(
	ctx sdk.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	resp, err := c.Middleware.ProcessProposal(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*cmtabci.ProcessProposalResponse), nil
}

// TODO: Decouple Comet Types
func (c *ConsensusEngine[ValidatorUpdateT]) PreBlock(
	ctx sdk.Context,
	req *cmtabci.FinalizeBlockRequest,
) error {
	return c.Middleware.PreBlock(ctx, req)
}

func (c *ConsensusEngine[ValidatorUpdateT]) EndBlock(
	ctx context.Context,
) ([]ValidatorUpdateT, error) {
	updates, err := c.Middleware.EndBlock(ctx)
	if err != nil {
		return nil, err
	}
	return iter.MapErr(updates, convertValidatorUpdate[ValidatorUpdateT])
}
