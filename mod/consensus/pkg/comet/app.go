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

	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	abci "github.com/cometbft/cometbft/abci/types"
)

var _ abci.Application = (*Application[engine.Client])(nil)

type Application[NodeT engine.Client] struct {
	// TODO: remove this once we implement the whole interface
	// actually most of these methods are noops, could be useful
	// to keep it instead of implementing the whole interface
	abci.BaseApplication

	// Logger
	logger log.Logger[any]

	// TODO: remove this once we implement the rpc engine
	node NodeT

	// CometBFT Params
	lastFinalizedHeight  int64
	consensusParamsStore *comet.ConsensusParamsStore
}

func NewApplication[NodeT engine.Client](
	logger log.Logger[any],
	node NodeT,
	chainSpec common.ChainSpec,
) *Application[NodeT] {
	return &Application[NodeT]{
		BaseApplication:      abci.BaseApplication{},
		logger:               logger,
		node:                 node,
		consensusParamsStore: comet.NewConsensusParamsStore(chainSpec),
	}
}

func (app *Application[NodeT]) InitChain(
	ctx context.Context,
	req *abci.InitChainRequest,
) (*abci.InitChainResponse, error) {
	app.logger.Info(
		"Initializing chain",
		"chain id", req.ChainId,
		"height", req.InitialHeight,
	)
	valUpdates, err := app.node.InitChain(ctx, req.AppStateBytes)
	if err != nil {
		return nil, err
	}

	return &abci.InitChainResponse{
		ConsensusParams: req.ConsensusParams, // TODO: do we need to override this??? how can we abstract this???
		Validators:      convertValidatorUpdates(valUpdates),
		AppHash:         nil,
	}, nil
}

func (app *Application[NodeT]) PrepareProposal(
	ctx context.Context,
	req *abci.PrepareProposalRequest,
) (*abci.PrepareProposalResponse, error) {
	app.logger.Info("PrepareProposal", "req", req)
	txs, err := app.node.PrepareProposal(ctx, req)
	if err != nil {
		return nil, err
	}
	return &abci.PrepareProposalResponse{
		Txs: txs,
	}, nil
}

func (app *Application[NodeT]) ProcessProposal(
	ctx context.Context,
	req *abci.ProcessProposalRequest,
) (*abci.ProcessProposalResponse, error) {
	app.logger.Info("ProcessProposal", "req", req)
	var err error
	status := abci.PROCESS_PROPOSAL_STATUS_ACCEPT
	if err = app.node.ProcessProposal(ctx, req); err != nil {
		status = abci.PROCESS_PROPOSAL_STATUS_REJECT
	}
	return &abci.ProcessProposalResponse{
		Status: status,
	}, err
}

func (app *Application[NodeT]) FinalizeBlock(
	ctx context.Context,
	req *abci.FinalizeBlockRequest,
) (*abci.FinalizeBlockResponse, error) {
	app.logger.Info("FinalizeBlock", "req", req)
	valUpdates, err := app.node.FinalizeBlock(ctx, req)
	if err != nil {
		return nil, err
	}
	app.lastFinalizedHeight = req.Height
	params, err := app.consensusParamsStore.Get(ctx)
	if err != nil {
		return nil, err
	}
	return &abci.FinalizeBlockResponse{
		ValidatorUpdates:      convertValidatorUpdates(valUpdates),
		ConsensusParamUpdates: &params,
		AppHash:               nil,
	}, nil
}

func (app *Application[NodeT]) Commit(
	ctx context.Context,
	req *abci.CommitRequest,
) (*abci.CommitResponse, error) {
	app.logger.Info("Commit", "req", req)
	return &abci.CommitResponse{
		RetainHeight: 0, // TODO: implement this
	}, nil
}
