package baseapp

import (
	gogogrpc "github.com/cosmos/gogoproto/grpc"
)

// RegisterGRPCServer registers gRPC services directly with the gRPC server.
func (app *BaseApp) RegisterGRPCServer(_ gogogrpc.Server) {}
