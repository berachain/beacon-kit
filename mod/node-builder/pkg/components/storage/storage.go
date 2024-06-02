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
		*types.BeaconBlockHeader, *types.ExecutionPayloadHeader,
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
		*types.BeaconBlockHeader, *types.ExecutionPayloadHeader,
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
