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
	"cosmossdk.io/depinject"
	"cosmossdk.io/depinject/appconfig"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	modulev1alpha1 "github.com/berachain/beacon-kit/mod/node-core/pkg/components/module/api/module/v1alpha1"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/config"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb/encoding"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// TODO: we don't allow generics here? Why? Is it fixable?
//
//nolint:gochecknoinits // required by sdk.
func init() {
	appconfig.RegisterModule(&modulev1alpha1.Module{},
		// TODO: make storage backend its own module and remove the
		// coupling between construction of runtime and module
		appconfig.Provide(
			ProvideStorageBackend,
			ProvideModule,
		),
	)
}

// DepInjectInput is the input for the dep inject framework.
type DepInjectInput struct {
	depinject.In
	BeaconConfig *config.Config
	Runtime      *components.BeaconKitRuntime
}

// DepInjectOutput is the output for the dep inject framework.
type DepInjectOutput struct {
	depinject.Out
	Module appmodule.AppModule
}

// ProvideModule is a function that provides the module to the application.
func ProvideModule(in DepInjectInput) (DepInjectOutput, error) {
	return DepInjectOutput{
		Module: NewAppModule(in.Runtime),
	}, nil
}

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput struct {
	depinject.In
	ChainSpec         primitives.ChainSpec
	AvailabilityStore *dastore.Store[*types.BeaconBlockBody]
	Environment       appmodule.Environment
	DepositStore      *depositdb.KVStore[*types.Deposit]
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend(
	in StorageBackendInput,
) *storage.Backend[
	*dastore.Store[*types.BeaconBlockBody],
	*types.BeaconBlock,
	*types.BeaconBlockBody,
	core.BeaconState[
		*types.BeaconBlockHeader, *types.Eth1Data,
		*types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal,
	],
	*depositdb.KVStore[*types.Deposit],
] {
	payloadCodec := &encoding.
		SSZInterfaceCodec[*types.ExecutionPayloadHeader]{}
	return storage.NewBackend[
		*dastore.Store[*types.BeaconBlockBody],
		*types.BeaconBlock,
		*types.BeaconBlockBody,
		core.BeaconState[
			*types.BeaconBlockHeader, *types.Eth1Data,
			*types.ExecutionPayloadHeader, *types.Fork,
			*types.Validator, *engineprimitives.Withdrawal,
		],
	](
		in.ChainSpec,
		in.AvailabilityStore,
		beacondb.New[
			*types.Fork,
			*types.BeaconBlockHeader,
			*types.ExecutionPayloadHeader,
			*types.Eth1Data,
			*types.Validator,
		](in.Environment.KVStoreService, payloadCodec),
		in.DepositStore,
	)
}
