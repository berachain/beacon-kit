// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package nodebuilder

import (
	"io"
	"reflect"

	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/node-builder/pkg/app"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/comet"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// NodeBuilder is a struct that holds the.
func (nb *NodeBuilder[T]) AppCreator(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) T {
	// Check for goleveldb cause bad project.
	if appOpts.Get("app-db-backend") == "goleveldb" {
		panic("goleveldb is not supported")
	}

	app := *app.NewBeaconKitApp(
		logger, db, traceStore, true,
		appOpts,
		nb.appInfo.DepInjectConfig,
		nb.chainSpec,
		append(
			server.DefaultBaseappOptions(appOpts),
			func(bApp *baseapp.BaseApp) {
				bApp.SetParamStore(comet.NewConsensusParamsStore(nb.chainSpec))
			})...,
	)
	return reflect.ValueOf(app).Convert(
		reflect.TypeOf((*T)(nil)).Elem()).Interface().(T)
}
