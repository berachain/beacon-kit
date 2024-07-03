package appmanager

import (
	"context"

	appmanager "cosmossdk.io/core/app"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/core/transaction"
	appmanagerv2 "cosmossdk.io/server/v2/appmanager"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
)

// AppManager is a wrapper around the AppManager from the Cosmos SDK.
// It is a wrapper around the ABCIMiddleware.
type AppManager[T transaction.Tx] struct {
	appmanagerv2.AppManager[T]
	abciMiddleware *components.ABCIMiddleware
}

func NewAppManager[T transaction.Tx](
	am appmanagerv2.AppManager[T],
	middleware *components.ABCIMiddleware,
) *AppManager[T] {
	return &AppManager[T]{
		am,
		middleware,
	}
}

func (am *AppManager[T]) InitGenesis(
	ctx context.Context,
	blockRequest *appmanager.BlockRequest[T],
	initGenesisJSON []byte,
	txDecoder transaction.Codec[T],
) (*appmanager.BlockResponse, corestore.WriterMap, error) {
	am.abciMiddleware.SetRequest(blockToABCIRequest(blockRequest))
	resp, writerMap, err := am.AppManager.InitGenesis(ctx, blockRequest, initGenesisJSON, txDecoder)
	if err != nil {
		return nil, nil, err
	}

	// run block
	// TODO: in an ideal world, genesis state is simply an initial state being applied
	// unaware of what that state means in relation to every other, so here we can
	// chain genesis
	return resp, writerMap, nil
}

func (am *AppManager[T]) DeliverBlock(
	ctx context.Context,
	block *appmanager.BlockRequest[T],
) (*appmanager.BlockResponse, corestore.WriterMap, error) {
	am.abciMiddleware.SetRequest(blockToABCIRequest(block))
	return am.AppManager.DeliverBlock(ctx, block)
}
