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

package builder

import (
	"context"
	"io"

	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/app"
	"github.com/berachain/beacon-kit/mod/primitives"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
)

// AppCreator is a function that creates an application and starts the bkRuntime
// services.
// It is necessary to adhere to the types.AppCreator[T] interface.
func (nb *NodeBuilder[NodeT]) AppCreator(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	appOpts servertypes.AppOptions,
) NodeT {
	// Check for goleveldb cause bad project.
	if appOpts.Get("app-db-backend") == "goleveldb" {
		panic("goleveldb is not supported")
	}

	// variables to hold the components needed to set up BeaconApp
	var chainSpec primitives.ChainSpec
	appBuilder := emptyAppBuilder()
	validatorMiddleware := emptyValidatorMiddleware()
	finalizeBlockMiddleware := emptyFinalizeBlockMiddlware()
	serviceRegistry := emptyServiceRegistry()

	// build all node components using depinject
	if err := depinject.Inject(
		depinject.Configs(
			nb.depInjectCfg,
			depinject.Provide(
				nb.components...,
			),
			depinject.Supply(
				appOpts,
				logger,
			),
		),
		&appBuilder,
		&chainSpec,
		&validatorMiddleware,
		&finalizeBlockMiddleware,
		&serviceRegistry,
	); err != nil {
		panic(err)
	}

	// set the application to a new BeaconApp with necessary ABCI handlers
	nb.node.SetApplication(
		app.NewBeaconKitApp(
			db, traceStore, true, appBuilder,
			append(
				server.DefaultBaseappOptions(appOpts),
				WithCometParamStore(chainSpec),
				WithPrepareProposal(validatorMiddleware.PrepareProposalHandler),
				WithProcessProposal(validatorMiddleware.ProcessProposalHandler),
				WithPreBlocker(finalizeBlockMiddleware.PreBlock),
			)...,
		),
	)

	// TODO: put this in some post node creation hook/listener
	// start all services
	if err := serviceRegistry.StartAll(context.Background()); err != nil {
		logger.Error("failed to start runtime services", "err", err)
		panic(err)
	}
	return nb.node
}
