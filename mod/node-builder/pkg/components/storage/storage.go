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

	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives"
	engineprimitives "github.com/berachain/beacon-kit/mod/primitives-engine"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/consensus"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// KVStore is a type alias for the beacon store with
// the generics defined using primitives.
type KVStore = beacondb.KVStore[
	*consensus.Fork,
	*consensus.BeaconBlockHeader,
	engineprimitives.ExecutionPayloadHeader,
	*consensus.Eth1Data,
	*consensus.Validator,
]

// Backend is a struct that holds the storage backend. It
// provides a simply interface to access all types of storage
// required by the runtime.
type Backend struct {
	cs                primitives.ChainSpec
	availabilityStore *dastore.Store[consensus.ReadOnlyBeaconBlockBody]
	beaconStore       *KVStore
	depositStore      *deposit.KVStore
}

func NewBackend(
	cs primitives.ChainSpec,
	availabilityStore *dastore.Store[consensus.ReadOnlyBeaconBlockBody],
	beaconStore *KVStore,
	depositStore *deposit.KVStore,
) *Backend {
	return &Backend{
		cs:                cs,
		availabilityStore: availabilityStore,
		beaconStore:       beaconStore,
		depositStore:      depositStore,
	}
}

// AvailabilityStore returns the availability store struct initialized with a.
func (k *Backend) AvailabilityStore(
	_ context.Context,
) core.AvailabilityStore[
	consensus.ReadOnlyBeaconBlockBody, *datypes.BlobSidecars,
] {
	return k.availabilityStore
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k *Backend) BeaconState(
	ctx context.Context,
) state.BeaconState {
	return state.NewBeaconStateFromDB(k.beaconStore.WithContext(ctx), k.cs)
}

// BeaconStore returns the beacon store struct.
func (k *Backend) BeaconStore() *KVStore {
	return k.beaconStore
}

// DepositStore returns the deposit store struct initialized with a.
func (k *Backend) DepositStore(
	_ context.Context,
) *deposit.KVStore {
	return k.depositStore
}
