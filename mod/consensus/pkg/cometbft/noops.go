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

	"cosmossdk.io/store/snapshots"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
)

// RegisterGRPCServer registers gRPC services directly with the gRPC server.
func (BaseApp) RegisterGRPCServer(_ gogogrpc.Server) {}

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (BaseApp) Query(
	_ context.Context,
	_ *abci.QueryRequest,
) (*abci.QueryResponse, error) {
	return &abci.QueryResponse{}, nil
}

// ListSnapshots implements the ABCI interface. It delegates to
// app.snapshotManager if set.
func (BaseApp) ListSnapshots(
	_ *abci.ListSnapshotsRequest,
) (*abci.ListSnapshotsResponse, error) {
	return &abci.ListSnapshotsResponse{}, nil
}

// LoadSnapshotChunk implements the ABCI interface. It delegates to
// app.snapshotManager if set.
func (BaseApp) LoadSnapshotChunk(
	_ *abci.LoadSnapshotChunkRequest,
) (*abci.LoadSnapshotChunkResponse, error) {
	return &abci.LoadSnapshotChunkResponse{}, nil
}

// OfferSnapshot implements the ABCI interface. It delegates to
// app.snapshotManager if set.
func (BaseApp) OfferSnapshot(
	_ *abci.OfferSnapshotRequest,
) (*abci.OfferSnapshotResponse, error) {
	return &abci.OfferSnapshotResponse{}, nil
}

// ApplySnapshotChunk implements the ABCI interface. It delegates to
// app.snapshotManager if set.
func (BaseApp) ApplySnapshotChunk(
	_ *abci.ApplySnapshotChunkRequest,
) (*abci.ApplySnapshotChunkResponse, error) {
	return &abci.ApplySnapshotChunkResponse{}, nil
}

// SnapshotManager returns the snapshot manager.
func (BaseApp) SnapshotManager() *snapshots.Manager {
	return &snapshots.Manager{}
}

func (BaseApp) ExtendVote(
	_ context.Context,
	_ *abci.ExtendVoteRequest,
) (*abci.ExtendVoteResponse, error) {
	return &abci.ExtendVoteResponse{}, nil
}

// VerifyVoteExtension implements the VerifyVoteExtension ABCI method and
// returns
// a ResponseVerifyVoteExtension. It calls the applications' VerifyVoteExtension
// handler which is responsible for performing application-specific business
// logic in verifying a vote extension from another validator during the
// pre-commit
// phase. The response MUST be deterministic. An error is returned if vote
// extensions are not enabled or if verifyVoteExt fails or panics.
// We highly recommend a size validation due to performance degradation,
// see more here
// https://docs.cometbft.com/v1.0/references/qa/cometbft-qa-38#vote-extensions-testbed
func (BaseApp) VerifyVoteExtension(
	*abci.VerifyVoteExtensionRequest,
) (*abci.VerifyVoteExtensionResponse, error) {
	return &abci.VerifyVoteExtensionResponse{}, nil
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (BaseApp) RegisterAPIRoutes(apiSvr *api.Server, _ config.APIConfig) {}

// RegisterTxService implements the Application.RegisterTxService method.
func (BaseApp) RegisterTxService(client.Context) {}

// RegisterTendermintService implements the
// Application.RegisterTendermintService method.
func (BaseApp) RegisterTendermintService(client.Context) {}

// RegisterNodeService registers the node gRPC service on the app gRPC router.
func (BaseApp) RegisterNodeService(_ client.Context, _ config.Config) {}

// CheckTx implements the ABCI interface and executes a tx in CheckTx mode. In
// CheckTx mode, messages are not executed. This means messages are only
// validated
// and only the AnteHandler is executed. State is persisted to the BaseApp's
// internal CheckTx state if the AnteHandler passes. Otherwise, the
// ResponseCheckTx
// will contain relevant error information. Regardless of tx execution outcome,
// the ResponseCheckTx will contain relevant gas execution context.
func (*BaseApp) CheckTx(
	*abci.CheckTxRequest,
) (*abci.CheckTxResponse, error) {
	return &abci.CheckTxResponse{}, nil
}
