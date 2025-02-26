// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/chain"
	"github.com/berachain/beacon-kit/consensus-types/types"
	dastore "github.com/berachain/beacon-kit/da/store"
	"github.com/berachain/beacon-kit/node-core/components/storage"
	"github.com/berachain/beacon-kit/storage/beacondb"
	"github.com/berachain/beacon-kit/storage/block"
	depositdb "github.com/berachain/beacon-kit/storage/deposit"
)

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput struct {
	depinject.In
	AvailabilityStore *dastore.Store
	BlockStore        *block.KVStore[*types.BeaconBlock]
	ChainSpec         chain.Spec
	DepositStore      *depositdb.KVStore
	BeaconStore       *beacondb.KVStore
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend(
	in StorageBackendInput,
) *storage.Backend {
	return storage.NewBackend(
		in.ChainSpec,
		in.AvailabilityStore,
		in.BeaconStore,
		in.DepositStore,
		in.BlockStore,
	)
}
