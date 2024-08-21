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

package types

import (
	"context"

	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
)

// ABCI is an interface that enables any finite, deterministic state machine
// to be driven by a blockchain-based replication engine via the ABCI.
type ABCI interface {
	// Info/Query Connection

	// Info returns application info
	Info(*abci.InfoRequest) (*abci.InfoResponse, error)
	// Query  returns application state
	Query(context.Context, *abci.QueryRequest) (*abci.QueryResponse, error)

	// Mempool Connection

	// CheckTx validate a tx for the mempool
	CheckTx(*abci.CheckTxRequest) (*abci.CheckTxResponse, error)

	// Consensus Connection

	// InitChain Initialize blockchain w validators/other info from CometBFT
	InitChain(*abci.InitChainRequest) (*abci.InitChainResponse, error)
	PrepareProposal(
		*abci.PrepareProposalRequest,
	) (*abci.PrepareProposalResponse, error)
	ProcessProposal(
		*abci.ProcessProposalRequest,
	) (*abci.ProcessProposalResponse, error)
	// FinalizeBlock deliver the decided block with its txs to the Application
	FinalizeBlock(
		*abci.FinalizeBlockRequest,
	) (*abci.FinalizeBlockResponse, error)
	// ExtendVote create application specific vote extension
	ExtendVote(
		context.Context,
		*abci.ExtendVoteRequest,
	) (*abci.ExtendVoteResponse, error)
	// VerifyVoteExtension verify application's vote extension data
	VerifyVoteExtension(
		*abci.VerifyVoteExtensionRequest,
	) (*abci.VerifyVoteExtensionResponse, error)
	// Commit the state and return the application Merkle root hash
	Commit() (*abci.CommitResponse, error)

	// State Sync Connection

	// ListSnapshots list available snapshots
	ListSnapshots(
		*abci.ListSnapshotsRequest,
	) (*abci.ListSnapshotsResponse, error)
	// OfferSnapshot offer a snapshot to the application
	OfferSnapshot(
		*abci.OfferSnapshotRequest,
	) (*abci.OfferSnapshotResponse, error)
	// LoadSnapshotChunk load a snapshot chunk
	LoadSnapshotChunk(
		*abci.LoadSnapshotChunkRequest,
	) (*abci.LoadSnapshotChunkResponse, error)
	// ApplySnapshotChunk apply a snapshot chunk
	ApplySnapshotChunk(
		*abci.ApplySnapshotChunkRequest,
	) (*abci.ApplySnapshotChunkResponse, error)
}
