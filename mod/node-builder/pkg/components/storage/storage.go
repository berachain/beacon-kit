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

package storage

import (
	"context"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/runtime/pkg/runtime"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// KVStore is a type alias for the beacon store with the generics defined using
// primitives.
type KVStore = beacondb.KVStore[
	*types.Fork, *types.BeaconBlockHeader, *types.ExecutionPayloadHeader,
	*types.Eth1Data, *types.Validator,
]

// Backend is a struct that holds the storage backend. It provides a simple
// interface to access all types of storage required by the runtime.
type Backend[
	AvailabilityStoreT runtime.AvailabilityStore[
		BeaconBlockBodyT, *datypes.BlobSidecars,
	],
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader, *types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal],
	DepositStoreT *deposit.KVStore[*types.Deposit],
] struct {
	cs primitives.ChainSpec
	as AvailabilityStoreT
	bs *KVStore
	ds DepositStoreT
}

func NewBackend[
	AvailabilityStoreT runtime.AvailabilityStore[
		BeaconBlockBodyT, *datypes.BlobSidecars,
	],
	BeaconBlockBodyT types.BeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader, *types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal],
	DepositStoreT *deposit.KVStore[*types.Deposit],
](
	cs primitives.ChainSpec,
	as AvailabilityStoreT,
	bs *KVStore,
	ds DepositStoreT,
) *Backend[AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT, DepositStoreT] {
	return &Backend[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT, DepositStoreT,
	]{
		cs: cs,
		as: as,
		bs: bs,
		ds: ds,
	}
}

// AvailabilityStore returns the availability store struct initialized with a.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT, DepositT,
]) AvailabilityStore(
	_ context.Context,
) AvailabilityStoreT {
	return k.as
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT, DepositT,
]) StateFromContext(
	ctx context.Context,
) BeaconStateT {
	return state.NewBeaconStateFromDB[BeaconStateT](
		k.bs.WithContext(ctx), k.cs,
	)
}

// BeaconStore returns the beacon store struct.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT, DepositStoreT,
]) BeaconStore() *KVStore {
	return k.bs
}

// DepositStore returns the deposit store struct initialized with a.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT, DepositStoreT,
]) DepositStore(
	_ context.Context,
) DepositStoreT {
	return k.ds
}
