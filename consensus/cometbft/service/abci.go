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
	"strings"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	cmtabci "github.com/cometbft/cometbft/abci/types"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	return s.initChain(s.ctx, req)
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
	return s.processProposal(s.ctx, req)
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
	//nolint:contextcheck // see s.ctx comment for more details
	return s.finalizeBlock(s.ctx, req)
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

func (s *Service) Query(
	_ context.Context,
	req *abci.QueryRequest,
) (*abci.QueryResponse, error) {
	resp := &abci.QueryResponse{}

	// add panic recovery for all queries
	//
	// Ref: https://github.com/cosmos/cosmos-sdk/pull/8039
	defer func() {
		if r := recover(); r != nil {
			resp = queryResult(errorsmod.Wrapf(sdkerrors.ErrPanic, "%v", r))
		}
	}()

	// when a client did not provide a query height, manually inject the latest
	if req.Height == 0 {
		req.Height = s.LastBlockHeight()
	}

	path := SplitABCIQueryPath(req.Path)
	if len(path) == 0 {
		return queryResult(errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "no query path provided")), nil
	}

	switch path[0] {
	case "store":
		// "/store" prefix for store queries
		req.Path = "/" + strings.Join(path[1:], "/")
		cms := s.sm.GetCommitMultiStore()
		queryable, ok := cms.(storetypes.Queryable)
		if !ok {
			panic(ok)
		}

		sdkReq := storetypes.RequestQuery(*req)
		resp, err := queryable.Query(&sdkReq)
		if err != nil {
			return queryResult(err), nil
		}
		resp.Height = req.Height

		abciResp := abci.QueryResponse(*resp)

		return &abciResp, nil

	default:
		resp = queryResult(errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "unknown query path"))
	}

	return resp, nil
}

// NOTE: Copied from here: https://github.com/cosmos/cosmos-sdk/blob/960d44842b9e313cbe762068a67a894ac82060ab/baseapp/errors.go#L37-L46.
// This was made public in v0.53, under the sdkerrors module.
// NOTE: the debug parameter has been removed since Service does not expose this functionality.
//
// queryResult returns a ResponseQuery from an error. It will try to parse ABCI
// info from the error.
func queryResult(err error) *abci.QueryResponse {
	space, code, log := errorsmod.ABCIInfo(err, false)
	return &abci.QueryResponse{
		Codespace: space,
		Code:      code,
		Log:       log,
	}
}

// NOTE: Copied from here: https://github.com/cosmos/cosmos-sdk/blob/960d44842b9e313cbe762068a67a894ac82060ab/baseapp/abci.go#L1153-L1165
//
// SplitABCIQueryPath splits a string path using the delimiter '/'.
//
// e.g. "this/is/funny" becomes []string{"this", "is", "funny"}
func SplitABCIQueryPath(requestPath string) []string {
	path := strings.Split(requestPath, "/")

	// first element is empty string
	if len(path) > 0 && path[0] == "" {
		path = path[1:]
	}

	return path
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
