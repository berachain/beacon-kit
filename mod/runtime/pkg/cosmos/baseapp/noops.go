package baseapp

import (
	"context"

	"cosmossdk.io/store/snapshots"
	abci "github.com/cometbft/cometbft/api/cometbft/abci/v1"
	gogogrpc "github.com/cosmos/gogoproto/grpc"
)

// RegisterGRPCServer registers gRPC services directly with the gRPC server.
func (app *BaseApp) RegisterGRPCServer(_ gogogrpc.Server) {}

// Query implements the ABCI interface. It delegates to CommitMultiStore if it
// implements Queryable.
func (app *BaseApp) Query(_ context.Context, _ *abci.QueryRequest) (*abci.QueryResponse, error) {
	return &abci.QueryResponse{}, nil
}

// ListSnapshots implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) ListSnapshots(_ *abci.ListSnapshotsRequest) (*abci.ListSnapshotsResponse, error) {
	return &abci.ListSnapshotsResponse{}, nil
}

// LoadSnapshotChunk implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) LoadSnapshotChunk(_ *abci.LoadSnapshotChunkRequest) (*abci.LoadSnapshotChunkResponse, error) {
	return &abci.LoadSnapshotChunkResponse{}, nil
}

// OfferSnapshot implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) OfferSnapshot(_ *abci.OfferSnapshotRequest) (*abci.OfferSnapshotResponse, error) {
	return &abci.OfferSnapshotResponse{}, nil
}

// ApplySnapshotChunk implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) ApplySnapshotChunk(_ *abci.ApplySnapshotChunkRequest) (*abci.ApplySnapshotChunkResponse, error) {
	return &abci.ApplySnapshotChunkResponse{}, nil
}

// SnapshotManager returns the snapshot manager.
func (app *BaseApp) SnapshotManager() *snapshots.Manager {
	return &snapshots.Manager{}
}

func (app *BaseApp) ExtendVote(_ context.Context, _ *abci.ExtendVoteRequest) (*abci.ExtendVoteResponse, error) {
	return &abci.ExtendVoteResponse{}, nil
}

// VerifyVoteExtension implements the VerifyVoteExtension ABCI method and returns
// a ResponseVerifyVoteExtension. It calls the applications' VerifyVoteExtension
// handler which is responsible for performing application-specific business
// logic in verifying a vote extension from another validator during the pre-commit
// phase. The response MUST be deterministic. An error is returned if vote
// extensions are not enabled or if verifyVoteExt fails or panics.
// We highly recommend a size validation due to performance degradation,
// see more here https://docs.cometbft.com/v1.0/references/qa/cometbft-qa-38#vote-extensions-testbed
func (app *BaseApp) VerifyVoteExtension(*abci.VerifyVoteExtensionRequest) (*abci.VerifyVoteExtensionResponse, error) {
	return &abci.VerifyVoteExtensionResponse{}, nil
}
