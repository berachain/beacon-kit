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

package runtime

import (
	"io"

	"cosmossdk.io/log"

	"github.com/berachain/beacon-kit/mod/runtime/pkg/cosmos/baseapp"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/version"
)

// AppBuilder is a type that is injected into a container by the runtime module
// (as *AppBuilder) which can be used to create an app which is compatible with
// the existing app.go initialization conventions.
type AppBuilder struct {
	App        *App
	Middleware Middleware
}

// Build builds an *App instance.
func (a *AppBuilder) Build(
	db dbm.DB,
	_ io.Writer,
	logger log.Logger,
	baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	bApp := baseapp.NewBaseApp(
		"BeaconKit",
		logger,
		db,
		baseAppOptions...)
	bApp.SetVersion(version.Version)
	bApp.MountStores(a.App.StoreKeys...)
	a.App.BaseApp = bApp
	return a.App
}
