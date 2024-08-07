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

	storetypes "cosmossdk.io/store/types"
	"github.com/berachain/beacon-kit/mod/consensus/pkg/engine"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	abci "github.com/cometbft/cometbft/abci/types"
)

var _ abci.Application = (*Application[engine.Client])(nil)

type Application[ClientT engine.Client] struct {
	// TODO: remove this once we implement the whole interface
	// actually most of these methods are noops, could be useful
	// to keep it instead of implementing the whole interface
	abci.BaseApplication

	// Logger
	logger log.Logger[any]

	// TODO: remove this once we implement the rpc engine
	client ClientT

	// Consensus Params
	ConsensusParamsStore *ConsensusParamsStore
}

func NewApplication[ClientT engine.Client](
	logger log.Logger[any],
	client ClientT,
	chainSpec common.ChainSpec,
) *Application[ClientT] {
	return &Application[ClientT]{
		BaseApplication:      abci.BaseApplication{},
		logger:               logger,
		client:               client,
		ConsensusParamsStore: NewConsensusParamsStore(chainSpec),
	}
}

func (app *Application[ClientT]) InitChain(
	ctx context.Context,
	req *abci.InitChainRequest,
) (*abci.InitChainResponse, error) {
	app.logger.Info(
		"Initializing chain",
		"chain id", req.ChainId,
		"height", req.InitialHeight,
	)
	valUpdates, _, err := app.client.InitChain(ctx, req.AppStateBytes)
	if err != nil {
		return nil, err
	}
	params, err := app.ConsensusParamsStore.Get(uint64(req.InitialHeight))
	if err != nil {
		return nil, err
	}

	// annoying from sdk.
	appHash := (&storetypes.CommitInfo{}).CommitID().Hash

	return &abci.InitChainResponse{
		ConsensusParams: params,
		Validators:      convertValidatorUpdates(valUpdates),
		AppHash:         appHash,
	}, nil
}

func (app *Application[ClientT]) PrepareProposal(
	ctx context.Context,
	req *abci.PrepareProposalRequest,
) (*abci.PrepareProposalResponse, error) {
	app.logger.Info("Starting PrepareProposal")
	prepareReq := prepareRequestFromABCIRequest(req)
	txs, err := app.client.PrepareProposal(ctx, prepareReq)
	if err != nil {
		return nil, err
	}
	return &abci.PrepareProposalResponse{
		Txs: txs,
	}, nil
}

func (app *Application[ClientT]) ProcessProposal(
	ctx context.Context,
	req *abci.ProcessProposalRequest,
) (*abci.ProcessProposalResponse, error) {
	app.logger.Info("Starting ProcessProposal")
	var err error
	processReq := processRequestFromABCIRequest(req)
	status := abci.PROCESS_PROPOSAL_STATUS_ACCEPT
	if err = app.client.ProcessProposal(ctx, processReq); err != nil {
		status = abci.PROCESS_PROPOSAL_STATUS_REJECT
	}
	return &abci.ProcessProposalResponse{
		Status: status,
	}, err
}

func (app *Application[ClientT]) FinalizeBlock(
	ctx context.Context,
	req *abci.FinalizeBlockRequest,
) (*abci.FinalizeBlockResponse, error) {
	app.logger.Info("Starting FinalizeBlock")
	finalizeReq := finalizeRequestFromABCIRequest(req)
	valUpdates, appHash, err := app.client.FinalizeBlock(ctx, finalizeReq)
	if err != nil {
		return nil, err
	}
	params, err := app.ConsensusParamsStore.Get(finalizeReq.Slot)
	if err != nil {
		return nil, err
	}
	return &abci.FinalizeBlockResponse{
		ValidatorUpdates:      convertValidatorUpdates(valUpdates),
		ConsensusParamUpdates: params,
		AppHash:               appHash,
		TxResults:             execTxResultsFromTxs(req.Txs),
	}, nil
}

func (app *Application[ClientT]) Commit(
	ctx context.Context,
	req *abci.CommitRequest,
) (*abci.CommitResponse, error) {
	app.logger.Info("Starting Commit")
	return &abci.CommitResponse{
		RetainHeight: 0, // TODO: implement this
	}, nil
}
