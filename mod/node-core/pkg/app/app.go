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

package app

import (
	coreapp "cosmossdk.io/core/app"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/runtime/v2"
	serverv2 "cosmossdk.io/server/v2"
	"cosmossdk.io/server/v2/appmanager"
	bkcomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	"github.com/cosmos/cosmos-sdk/client"
)

var (
	_ runtime.AppI[transaction.Tx]  = (*BeaconApp[transaction.Tx])(nil)
	_ serverv2.AppI[transaction.Tx] = (*BeaconApp[transaction.Tx])(nil)
)

// BeaconApp extends an ABCI application, but with most of its parameters
// exported. They are exported for convenience in creating helper
// functions, as object capabilities aren't needed for testing.
type BeaconApp[T transaction.Tx] struct {
	*runtime.App[T]
	middleware *bkcomponents.ABCIMiddleware
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func NewBeaconKitApp[T transaction.Tx](
	sdkApp *runtime.App[T],
	middleware *bkcomponents.ABCIMiddleware,
) *BeaconApp[T] {
	app := &BeaconApp[T]{
		App:        sdkApp,
		middleware: middleware,
	}

	// app.SetTxDecoder(bkcomponents.NoOpTxConfig{}.TxDecoder())

	// Load the app.
	if err := app.LoadLatest(); err != nil {
		panic(err)
	}

	return app
}

// InterfaceRegistry returns BeaconApp's InterfaceRegistry.
func (app *BeaconApp[T]) InterfaceRegistry() coreapp.InterfaceRegistry {
	return nil
}

// TxConfig returns BeaconApp's TxConfig.
func (app *BeaconApp[T]) TxConfig() client.TxConfig {
	return nil
}

// GetConsensusAuthority gets the consensus authority.
func (app *BeaconApp[T]) GetConsensusAuthority() string {
	return "gov"
}

// GetStore gets the app store.
func (app *BeaconApp[T]) GetStore() any {
	return app.App.GetStore()
}

func (app *BeaconApp[T]) GetAppManager() appmanager.AppManager[T] {
	return app.App.GetAppManager()
}
