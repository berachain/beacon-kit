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
func (app *BaseApp) Query(_ context.Context, req *abci.QueryRequest) (resp *abci.QueryResponse, err error) {
	return resp, nil
}

// ListSnapshots implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) ListSnapshots(req *abci.ListSnapshotsRequest) (*abci.ListSnapshotsResponse, error) {
	return nil, nil
}

// LoadSnapshotChunk implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) LoadSnapshotChunk(req *abci.LoadSnapshotChunkRequest) (*abci.LoadSnapshotChunkResponse, error) {
	return nil, nil
}

// OfferSnapshot implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) OfferSnapshot(req *abci.OfferSnapshotRequest) (*abci.OfferSnapshotResponse, error) {
	return nil, nil
}

// ApplySnapshotChunk implements the ABCI interface. It delegates to app.snapshotManager if set.
func (app *BaseApp) ApplySnapshotChunk(req *abci.ApplySnapshotChunkRequest) (*abci.ApplySnapshotChunkResponse, error) {
	return nil, nil
}

// SnapshotManager returns the snapshot manager.
func (app *BaseApp) SnapshotManager() *snapshots.Manager {
	return nil
}
