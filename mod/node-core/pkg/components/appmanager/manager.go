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

package appmanager

import (
	"context"

	appmanager "cosmossdk.io/core/app"
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/core/transaction"
	appmanagerv2 "cosmossdk.io/server/v2/appmanager"
)

// AppManager is a wrapper around the AppManager from the Cosmos SDK.
// It is a wrapper around the ABCIMiddleware.
type AppManager[T transaction.Tx] struct {
	appmanagerv2.AppManager[T]
	abciMiddleware *ABCIMiddleware
}

func NewAppManager[T transaction.Tx](
	am appmanagerv2.AppManager[T],
	middleware *ABCIMiddleware,
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
	resp, writerMap, err := am.AppManager.DeliverBlock(ctx, block)
	if err != nil {
		return nil, nil, err
	}

	// apply the block state changes to the writer map
	if err := writerMap.ApplyStateChanges(am.abciMiddleware.GetBlockStateChanges()); err != nil {
		return nil, nil, err
	}
	am.abciMiddleware.FlushBlockStore()
	return resp, writerMap, nil
}
