package appmanager

import (
	"context"
	"fmt"

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
	// var genTxs []T
	// var rawMessages map[string]json.RawMessage
	// if err := json.Unmarshal(initGenesisJSON, &rawMessages); err != nil {
	// 	return nil, nil, fmt.Errorf("failed to unmarshal initGenesisJSON: %w", err)
	// }
	// // decode transactions
	// for _, rawMessage := range rawMessages {
	// 	tx, err := txDecoder.DecodeJSON(rawMessage)
	// 	if err != nil {
	// 		return nil, nil, fmt.Errorf("failed to decode tx: %w", err)
	// 	}
	// 	genTxs = append(genTxs, tx)
	// }

	// blockRequest.Txs = genTxs

	am.abciMiddleware.SetRequest(blockToABCIRequest(blockRequest))
	resp, writerMap, err := am.AppManager.InitGenesis(ctx, blockRequest, initGenesisJSON, txDecoder)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("DONE INIT GENESIS")

	// run block
	// TODO: in an ideal world, genesis state is simply an initial state being applied
	// unaware of what that state means in relation to every other, so here we can
	// chain genesis
	// blockRequest.Txs = genTxs

	return resp, writerMap, nil
}

func (am *AppManager[T]) DeliverBlock(
	ctx context.Context,
	block *appmanager.BlockRequest[T],
) (*appmanager.BlockResponse, corestore.WriterMap, error) {
	am.abciMiddleware.SetRequest(blockToABCIRequest(block))
	return am.AppManager.DeliverBlock(ctx, block)
}
