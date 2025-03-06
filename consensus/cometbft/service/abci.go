// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
//

package cometbft

import (
	"context"
	"errors"
	"fmt"

	cmtabci "github.com/cometbft/cometbft/abci/types"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	sdkversion "github.com/cosmos/cosmos-sdk/version"
)

var (
	errInvalidHeight         = errors.New("invalid height")
	errNilFinalizeBlockState = errors.New("finalizeBlockState is nil")
)

func (s *Service) InitChain(
	_ context.Context,
	req *cmtabci.InitChainRequest,
) (*cmtabci.InitChainResponse, error) {
	// Check if ctx is still good. CometBFT does not check this.
	if s.ctx.Err() != nil {
		// If the context is getting cancelled, we are shutting down.
		return &cmtabci.InitChainResponse{}, s.ctx.Err()
	}
	//nolint:contextcheck // see s.ctx comment for more details
	return s.initChain(req) // internally this uses s.ctx
}

// PrepareProposal implements the PrepareProposal ABCI method and returns a
// ResponsePrepareProposal object to the client.
func (s *Service) PrepareProposal(
	_ context.Context,
	req *cmtabci.PrepareProposalRequest,
) (*cmtabci.PrepareProposalResponse, error) {
	// Check if ctx is still good. CometBFT does not check this.
	if s.ctx.Err() != nil {
		// If the context is getting cancelled, we are shutting down.
		// It is ok returning an empty proposal.
		//nolint:nilerr // explicitly allowing this case
		return &cmtabci.PrepareProposalResponse{Txs: req.Txs}, nil
	}
	//nolint:contextcheck // see s.ctx comment for more details
	return s.prepareProposal(s.ctx, req)
}

func (s *Service) Info(context.Context,
	*cmtabci.InfoRequest,
) (*cmtabci.InfoResponse, error) {
	lastCommitID := s.sm.GetCommitMultiStore().LastCommitID()
	appVersion := initialAppVersion
	if lastCommitID.Version > 0 {
		var err error
		appVersion, err = s.appVersion()
		if err != nil {
			return nil, fmt.Errorf("failed getting app version: %w", err)
		}
	}

	return &cmtabci.InfoResponse{
		Data:             AppName,
		Version:          sdkversion.Version,
		AppVersion:       appVersion,
		LastBlockHeight:  lastCommitID.Version,
		LastBlockAppHash: lastCommitID.Hash,
	}, nil
}

// ProcessProposal implements the ProcessProposal ABCI method and returns a
// ResponseProcessProposal object to the client.
func (s *Service) ProcessProposal(
	_ context.Context,
	req *cmtabci.ProcessProposalRequest,
) (*cmtabci.ProcessProposalResponse, error) {
	// Check if ctx is still good. CometBFT does not check this.
	if s.ctx.Err() != nil {
		// Node will panic on context cancel with "CONSENSUS FAILURE!!!" due to
		// returning an error. This is expected. We do not want to accept or
		// reject a proposal based on incomplete data.
		return nil, s.ctx.Err()
	}
	//nolint:contextcheck // see s.ctx comment for more details
	return s.processProposal(req) // internally this uses s.ctx
}

func (s *Service) FinalizeBlock(
	_ context.Context,
	req *cmtabci.FinalizeBlockRequest,
) (*cmtabci.FinalizeBlockResponse, error) {
	// Check if ctx is still good. CometBFT does not check this.
	if s.ctx.Err() != nil {
		// Node will panic on context cancel with "CONSENSUS FAILURE!!!" due to error.
		// We expect this to happen and do not want to finalize any incomplete or invalid state.
		return nil, s.ctx.Err()
	}
	return s.finalizeBlock(req) // internally this uses s.ctx
}

// Commit implements the ABCI interface. It will commit all state that exists in
// the deliver state's multi-store and includes the resulting commit ID in the
// returned cmtabci.ResponseCommit. Commit will set the check state based on the
// latest header and reset the deliver state. Also, if a non-zero halt height is
// defined in config, Commit will execute a deferred function call to check
// against that height and gracefully halt if it matches the latest committed
// height.
func (s *Service) Commit(
	_ context.Context, req *cmtabci.CommitRequest,
) (*cmtabci.CommitResponse, error) {
	// Check if ctx is still good. CometBFT does not check this.
	if s.ctx.Err() != nil {
		// Node will panic on context cancel with "CONSENSUS FAILURE!!!" due to error.
		// We expect this to happen and do not want to commit any incomplete or invalid state.
		return nil, s.ctx.Err()
	}

	return s.commit(req)
}

//
// NOOP methods
//

func (Service) Query(
	context.Context,
	*abci.QueryRequest,
) (*abci.QueryResponse, error) {
	return &abci.QueryResponse{}, nil
}

func (Service) ListSnapshots(
	context.Context,
	*abci.ListSnapshotsRequest,
) (*abci.ListSnapshotsResponse, error) {
	return &abci.ListSnapshotsResponse{}, nil
}

func (Service) LoadSnapshotChunk(
	context.Context,
	*abci.LoadSnapshotChunkRequest,
) (*abci.LoadSnapshotChunkResponse, error) {
	return &abci.LoadSnapshotChunkResponse{}, nil
}

func (Service) OfferSnapshot(
	context.Context,
	*abci.OfferSnapshotRequest,
) (*abci.OfferSnapshotResponse, error) {
	return &abci.OfferSnapshotResponse{}, nil
}

func (Service) ApplySnapshotChunk(
	context.Context,
	*abci.ApplySnapshotChunkRequest,
) (*abci.ApplySnapshotChunkResponse, error) {
	return &abci.ApplySnapshotChunkResponse{}, nil
}

func (Service) ExtendVote(
	context.Context,
	*abci.ExtendVoteRequest,
) (*abci.ExtendVoteResponse, error) {
	return &abci.ExtendVoteResponse{}, nil
}

func (Service) VerifyVoteExtension(
	context.Context,
	*abci.VerifyVoteExtensionRequest,
) (*abci.VerifyVoteExtensionResponse, error) {
	return &abci.VerifyVoteExtensionResponse{}, nil
}

func (*Service) CheckTx(
	context.Context,
	*abci.CheckTxRequest,
) (*abci.CheckTxResponse, error) {
	return &abci.CheckTxResponse{}, nil
}
