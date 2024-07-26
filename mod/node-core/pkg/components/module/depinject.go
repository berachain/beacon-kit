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

package beacon

import (
	// runtimev2 "cosmossdk.io/api/cosmos/app/runtime/v2"
	// appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appmodule"
	appmodulev2 "cosmossdk.io/core/appmodule/v2"

	// "cosmossdk.io/core/legacy"
	// "cosmossdk.io/core/registry"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"

	// "cosmossdk.io/log"
	// "cosmossdk.io/runtime/v2"
	// rootstorev2 "cosmossdk.io/store/v2/root"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	modulev2 "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module/api/module/v2"

	cmtruntime "github.com/berachain/beacon-kit/mod/runtime/pkg/cometbft"
)

// TODO: we don't allow generics here? Why? Is it fixable?
//
//nolint:gochecknoinits // required by sdk.
func init() {
	appconfig.RegisterModule(&modulev2.Module{},
		appconfig.Provide(
			components.ProvideKVStore,
			components.ProvideMessageServer,
			ProvideModule[
				transaction.Tx,
				appmodulev2.ValidatorUpdate, // TODO: idk man
			],
		),
	)
}

// ModuleInput is the input for the dep inject framework.
type ModuleInput[T transaction.Tx] struct {
	depinject.In
	ABCIMiddleware *components.ABCIMiddleware
	TxCodec        transaction.Codec[T]
	MsgServer      *cmtruntime.MsgServer
	QueryServer    *cmtruntime.QueryServer
	StorageBackend *components.StorageBackend
}

// ModuleOutput is the output for the dep inject framework.
type ModuleOutput struct {
	depinject.Out
	Module appmodule.AppModule
}

// ProvideModule is a function that provides the module to the application.
func ProvideModule[T transaction.Tx, ValidatorUpdateT any](
	in ModuleInput[T],
) (ModuleOutput, error) {
	return ModuleOutput{
		Module: NewAppModule(
			in.ABCIMiddleware,
			in.TxCodec,
			in.MsgServer,
			in.QueryServer,
			in.StorageBackend,
		),
	}, nil
}

func SupplyModuleDependencies() []any {
	return []any{
		&components.ABCIMiddleware{},
	}
}
