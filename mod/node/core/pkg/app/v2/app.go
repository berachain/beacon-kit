// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	corestore "cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"cosmossdk.io/runtime/v2"
	"cosmossdk.io/server/v2/cometbft"
	"cosmossdk.io/server/v2/cometbft/flags"
	"cosmossdk.io/store/v2"
	"cosmossdk.io/store/v2/commitment/iavl"
	"cosmossdk.io/store/v2/root"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	bkcomponents "github.com/berachain/beacon-kit/mod/node/core/pkg/components"
	beacon "github.com/berachain/beacon-kit/mod/node/core/pkg/components/module/v2"
	"github.com/berachain/beacon-kit/mod/primitives"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

var (
// _ runtime.AppI[types.Tx] = (*BeaconApp[types.Tx])(nil)
)

type BeaconApp[TransactionT Tx] struct {
	*runtime.App
	CmtServer *cometbft.CometBFTServer[types.Tx]
}

func NewBeaconApp[TransactionT Tx](
	logger log.Logger,
	dCfg depinject.Config,
	db corestore.KVStoreWithBatch,
	appOpts servertypes.AppOptions,
	chainSpec primitives.ChainSpec,
) *BeaconApp[TransactionT] {
	var err error
	homeDir := cast.ToString(appOpts.Get(flags.FlagHome))
	app := &BeaconApp[TransactionT]{}
	appBuilder := &runtime.AppBuilder{}
	if err = depinject.Inject(
		depinject.Configs(
			dCfg,
			depinject.Provide(
				bkcomponents.ProvideAvailibilityStore,
				bkcomponents.ProvideBlsSigner,
				bkcomponents.ProvideTrustedSetup,
				bkcomponents.ProvideDepositStore,
				bkcomponents.ProvideConfig,
				bkcomponents.ProvideEngineClient,
				bkcomponents.ProvideJWTSecret,
				bkcomponents.ProvideTelemetrySink,
			),
			depinject.Supply(
				appOpts,
				logger,
				chainSpec,
				&root.FactoryOptions{
					Logger:  logger,
					RootDir: homeDir,
					SSType:  0,
					SCType:  0,
					SCPruneOptions: &store.PruneOptions{
						KeepRecent: 0,
						Interval:   0,
					},
					IavlConfig: &iavl.Config{
						CacheSize:              100_000,
						SkipFastStorageUpgrade: true,
					},
					SCRawDB: db,
				},
			),
		),
		&appBuilder,
	); err != nil {
		panic(err)
	}

	if app.App, err = appBuilder.Build(); err != nil {
		panic(err)
	}

	if err = app.LoadLatest(); err != nil {
		panic(err)
	}

	// TODO: tx Config

	return app
}

func (app *BeaconApp[TransactionT]) setupModule() {
	module, ok := app.ModuleManager().
		Modules()[beacon.ModuleName].(beacon.AppModule)
	if !ok {
		panic("module not found")
	}

	app.CmtServer.App.SetPrepareProposalHandler(
		module.ABCIValidatorMiddleware().PrepareProposalHandler(),
	)

	app.CmtServer.App.SetProcessProposalHandler(
		module.ABCIValidatorMiddleware().ProcessProposalHandler(),
	)

	app.App.ModuleManager().PreBlocker()
}
