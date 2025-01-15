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
	"github.com/berachain/beacon-kit/node-core/components/storage"
)

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput[
	AvailabilityStoreT any,
	BeaconBlockStoreT any,
	BeaconStoreT any,
	DepositStoreT any,
] struct {
	depinject.In
	AvailabilityStore AvailabilityStoreT
	BlockStore        BeaconBlockStoreT
	ChainSpec         chain.ChainSpec
	DepositStore      DepositStoreT
	BeaconStore       BeaconStoreT
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend[
	AvailabilityStoreT any,
	BeaconBlockStoreT any,
	BeaconStoreT interface {
		WithContext(context.Context) BeaconStoreT
	},
	DepositStoreT any,
](
	in StorageBackendInput[
		AvailabilityStoreT, BeaconBlockStoreT, BeaconStoreT,
		DepositStoreT,
	],
) *storage.Backend[
	AvailabilityStoreT, BeaconBlockStoreT,
	DepositStoreT, BeaconStoreT,
] {
	return storage.NewBackend[
		AvailabilityStoreT,
		BeaconBlockStoreT,
		DepositStoreT,
		BeaconStoreT,
	](
		in.ChainSpec,
		in.AvailabilityStore,
		in.BeaconStore,
		in.DepositStore,
		in.BlockStore,
	)
}
