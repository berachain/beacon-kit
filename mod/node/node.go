package node

import (
	"github.com/berachain/beacon-kit/mod/node-builder/service"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// BeaconNode is a beacon node.
type BeaconNode struct {
	// app represents the core abci application.
	app servertypes.Application

	// sr is a service registry filled with sidecar
	// services that are started and stopped with the node.
	sr *service.Registry
}

// New creates a new beacon node.
func New() *BeaconNode {
	// New creates a new beacon node.
	return &BeaconNode{}
}

// Start starts the beacon node.
func (bn *BeaconNode) Start(
	svrCtx *server.Context,
	clientCtx client.Context,
	appCreator servertypes.AppCreator[servertypes.Application],
	inProcessConsensus bool,
	opts server.StartCmdOptions[servertypes.Application],
) error {
	return nil
}

// Stop stops the beacon node.
func (bn *BeaconNode) Stop() {
	return
}
