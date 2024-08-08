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
	"io"

	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/baseapp"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/runtime"
	dbm "github.com/cosmos/cosmos-db"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

var (
	_ runtime.AppI            = (*BeaconApp)(nil)
	_ servertypes.Application = (*BeaconApp)(nil)
)

// BeaconApp extends an ABCI application, but with most of its parameters
// exported. They are exported for convenience in creating helper
// functions, as object capabilities aren't needed for testing.
type BeaconApp struct {
	*runtime.App
}

// NewBeaconKitApp returns a reference to an initialized BeaconApp.
func NewBeaconKitApp(
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appBuilder *runtime.AppBuilder,
	baseAppOptions ...func(*baseapp.BaseApp),
) *BeaconApp {
	app := &BeaconApp{}

	// Build the runtime.App using the app builder.
	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// Load the app.
	if err := app.Load(loadLatest); err != nil {
		panic(err)
	}

	return app
}
