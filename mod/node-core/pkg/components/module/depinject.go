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
	"cosmossdk.io/core/appmodule"
	appmodulev2 "cosmossdk.io/core/appmodule/v2"
	"cosmossdk.io/core/transaction"
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	modulev2 "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module/api/module/v2"
)

// TODO: we don't allow generics here? Why? Is it fixable?
//
//nolint:gochecknoinits // required by sdk.
func init() {
	appconfig.RegisterModule(&modulev2.Module{},
		appconfig.Provide(
			components.ProvideKVStore,
			ProvideModule[
				transaction.Tx, appmodulev2.ValidatorUpdate, // TODO: idk man
			],
		),
	)
}

// ModuleInput is the input for the dep inject framework.
type ModuleInput[T transaction.Tx] struct {
	depinject.In
	ABCIMiddleware *components.ABCIMiddleware
	TxCodec        transaction.Codec[T]
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
		Module: NewAppModule[T, ValidatorUpdateT](
			in.ABCIMiddleware,
			in.TxCodec,
		),
	}, nil
}

func SupplyModuleDependencies() []any {
	return []any{
		&components.ABCIMiddleware{},
	}
}
