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

package components

import (
	"context"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/chain-spec/chain"
	dastore "github.com/berachain/beacon-kit/da/store"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
)

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput[

	BeaconBlockStoreT any,
	BeaconStoreT any,
] struct {
	depinject.In
	AvailabilityStore *dastore.Store
	BlockStore        BeaconBlockStoreT
	ChainSpec         chain.ChainSpec
	DepositStore      *depositdb.KVStore
	BeaconStore       BeaconStoreT
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend[
	BeaconBlockStoreT any,
	BeaconStoreT interface {
		WithContext(context.Context) BeaconStoreT
	},
](
	in StorageBackendInput[
		BeaconBlockStoreT, BeaconStoreT,
	],
) *storage.Backend[
	BeaconBlockStoreT,
	BeaconStoreT,
] {
	return storage.NewBackend[
		BeaconBlockStoreT,
		BeaconStoreT,
	](
		in.ChainSpec,
		in.AvailabilityStore,
		in.BeaconStore,
		in.DepositStore,
		in.BlockStore,
	)
}
