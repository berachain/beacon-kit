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

package storage

import (
	"context"

	"github.com/berachain/beacon-kit/chain-spec/chain"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	"github.com/berachain/beacon-kit/storage/beacondb"
)

// Backend is a struct that holds the storage backend. It provides a simple
// interface to access all types of storage required by the runtime.
type Backend[
	AvailabilityStoreT any,
	BlockStoreT any,
	DepositStoreT any,
	KVStoreT KVStore[KVStoreT],
] struct {
	chainSpec         chain.ChainSpec
	availabilityStore AvailabilityStoreT
	kvStore           KVStoreT
	depositStore      DepositStoreT
	blockStore        BlockStoreT
}

func NewBackend[
	AvailabilityStoreT any,
	BlockStoreT any,
	DepositStoreT any,
	KVStoreT KVStore[KVStoreT],
](
	chainSpec chain.ChainSpec,
	availabilityStore AvailabilityStoreT,
	kvStore KVStoreT,
	depositStore DepositStoreT,
	blockStore BlockStoreT,
) *Backend[
	AvailabilityStoreT, BlockStoreT, DepositStoreT, KVStoreT,
] {
	return &Backend[
		AvailabilityStoreT, BlockStoreT, DepositStoreT, KVStoreT,
	]{
		chainSpec:         chainSpec,
		availabilityStore: availabilityStore,
		kvStore:           kvStore,
		depositStore:      depositStore,
		blockStore:        blockStore,
	}
}

// AvailabilityStore returns the availability store struct initialized with a
// given context.
func (k Backend[
	AvailabilityStoreT, _, _, _,
]) AvailabilityStore() AvailabilityStoreT {
	return k.availabilityStore
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k Backend[
	_, _, _, KVStoreT,
]) StateFromContext(
	ctx context.Context,
) *statedb.StateDB {
	kvstore, ok := any(k.kvStore.WithContext(ctx)).(*beacondb.KVStore)
	if !ok {
		panic("failed to cast KVStoreT to beacondb.KVStore")
	}

	var st *statedb.StateDB
	return st.NewFromDB(kvstore, k.chainSpec)
}

// BeaconStore returns the beacon store struct.
func (k Backend[
	_, _, _, KVStoreT,
]) BeaconStore() KVStoreT {
	return k.kvStore
}

func (k Backend[
	_, BlockStoreT, _, _,
]) BlockStore() BlockStoreT {
	return k.blockStore
}

// DepositStore returns the deposit store struct initialized with a.
func (k Backend[
	_, _, DepositStoreT, _,
]) DepositStore() DepositStoreT {
	return k.depositStore
}
