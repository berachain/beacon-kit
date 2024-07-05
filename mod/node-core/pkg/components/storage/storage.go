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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
	"github.com/berachain/beacon-kit/mod/storage/pkg/beacondb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
)

// KVStore is a type alias for the beacon store with the generics defined using
// primitives.
type KVStore = beacondb.KVStore[
	*types.BeaconBlockHeader, *types.Eth1Data, *types.ExecutionPayloadHeader,
	*types.Fork, *types.Validator,
]

// The AvailabilityStore interface is responsible for validating and storing
// sidecars for specific blocks, as well as verifying sidecars that have already
// been stored.
type AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT any] interface {
	// IsDataAvailable ensures that all blobs referenced in the block are
	// securely stored before it returns without an error.
	IsDataAvailable(
		context.Context, math.Slot, BeaconBlockBodyT,
	) bool
	// Persist makes sure that the sidecar remains accessible for data
	// availability checks throughout the beacon node's operation.
	Persist(math.Slot, BlobSidecarsT) error
}

// Backend is a struct that holds the storage backend. It provides a simple
// interface to access all types of storage required by the runtime.
type Backend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, *datypes.BlobSidecars,
	],
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader, *types.Eth1Data, *types.ExecutionPayloadHeader,
		*types.Fork, *types.Validator, *engineprimitives.Withdrawal],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, *types.BeaconBlockHeader, *types.Eth1Data,
		*types.ExecutionPayloadHeader, *types.Fork, *types.Validator,
	],
	DepositStoreT *deposit.KVStore[*types.Deposit],
] struct {
	cs common.ChainSpec
	as AvailabilityStoreT
	bs *KVStore
	ds DepositStoreT
}

func NewBackend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, *datypes.BlobSidecars,
	],
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconStateT core.BeaconState[
		*types.BeaconBlockHeader, *types.Eth1Data,
		*types.ExecutionPayloadHeader, *types.Fork,
		*types.Validator, *engineprimitives.Withdrawal,
	],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, *types.BeaconBlockHeader, *types.Eth1Data,
		*types.ExecutionPayloadHeader, *types.Fork, *types.Validator,
	],
	DepositStoreT *deposit.KVStore[*types.Deposit],
](
	cs common.ChainSpec,
	as AvailabilityStoreT,
	bs *KVStore,
	ds DepositStoreT,
) *Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BeaconStateMarshallableT, DepositStoreT,
] {
	return &Backend[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
		BeaconStateMarshallableT, DepositStoreT,
	]{
		cs: cs,
		as: as,
		bs: bs,
		ds: ds,
	}
}

// AvailabilityStore returns the availability store struct initialized with a
// given context.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BeaconStateMarshallableT, DepositStoreT,
]) AvailabilityStore(
	_ context.Context,
) AvailabilityStoreT {
	return k.as
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BeaconStateMarshallableT, DepositStoreT,
]) StateFromContext(
	ctx context.Context,
) BeaconStateT {
	return state.NewBeaconStateFromDB[
		BeaconStateT, BeaconStateMarshallableT,
	](
		k.bs.WithContext(ctx), k.cs, // k.bs = kvstore.go / KVStore
	)
}

// BeaconStore returns the beacon store struct.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BeaconStateMarshallableT, DepositStoreT,
]) BeaconStore() *KVStore {
	return k.bs
}

// DepositStore returns the deposit store struct initialized with a.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconStateT,
	BeaconStateMarshallableT, DepositStoreT,
]) DepositStore(
	_ context.Context,
) DepositStoreT {
	return k.ds
}
