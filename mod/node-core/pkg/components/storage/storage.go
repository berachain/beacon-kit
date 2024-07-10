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

	types "github.com/berachain/beacon-kit/mod/interfaces/pkg/consensus-types"
	dastore "github.com/berachain/beacon-kit/mod/interfaces/pkg/da/store"
	engineprimitives "github.com/berachain/beacon-kit/mod/interfaces/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/state-transition/state"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/storage/beacondb"
	"github.com/berachain/beacon-kit/mod/interfaces/pkg/storage/deposit"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// Backend is a struct that holds the storage backend. It provides a simple
// interface to access all types of storage required by the runtime.
type Backend[
	AvailabilityStoreT dastore.AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT state.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.Marshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	DepositT types.Deposit[
		DepositT, ForkDataT, WithdrawalCredentialsT,
	],
	DepositStoreT deposit.Store[DepositT],
	Eth1DataT any,
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ForkDataT any,
	KVStoreT beacondb.KVStore[
		KVStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT types.Validator[ValidatorT, WithdrawalCredentialsT],
	WithdrawalT engineprimitives.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT types.WithdrawalCredentials[WithdrawalCredentialsT],
] struct {
	cs common.ChainSpec
	as AvailabilityStoreT
	bs KVStoreT
	ds DepositStoreT
}

func NewBackend[
	AvailabilityStoreT dastore.AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockBodyT types.BeaconBlockBody[
		BeaconBlockBodyT, DepositT, Eth1DataT, ExecutionPayloadT,
	],
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT state.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.Marshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	DepositT types.Deposit[
		DepositT, ForkDataT, WithdrawalCredentialsT,
	],
	DepositStoreT deposit.Store[DepositT],
	Eth1DataT any,
	ExecutionPayloadT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ForkDataT any,
	KVStoreT beacondb.KVStore[
		KVStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT types.Validator[ValidatorT, WithdrawalCredentialsT],
	WithdrawalT engineprimitives.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT types.WithdrawalCredentials[WithdrawalCredentialsT],
](
	cs common.ChainSpec,
	as AvailabilityStoreT,
	bs KVStoreT,
	ds DepositStoreT,
) *Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT,
	BeaconStateMarshallableT, BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT,
	ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT,
	ValidatorT, WithdrawalT, WithdrawalCredentialsT,
] {
	return &Backend[
		AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT,
		BeaconStateMarshallableT, BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT,
		ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT,
		ValidatorT, WithdrawalT, WithdrawalCredentialsT,
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
	AvailabilityStoreT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) AvailabilityStore(
	_ context.Context,
) AvailabilityStoreT {
	return k.as
}

// BeaconState returns the beacon state struct initialized with a given
// context and the store key.
func (k Backend[
	AvailabilityStoreT, BeaconBlockBodyT, BeaconBlockHeaderT, BeaconStateT,
	BeaconStateMarshallableT, BlobSidecarsT, DepositT, DepositStoreT, Eth1DataT,
	ExecutionPayloadT, ExecutionPayloadHeaderT, ForkT, ForkDataT, KVStoreT,
	ValidatorT, WithdrawalT, WithdrawalCredentialsT,
]) StateFromContext(
	ctx context.Context,
) BeaconStateT {
	var st BeaconStateT
	return st.NewFromDB(
		k.bs.WithContext(ctx), k.cs,
	)
}

// DepositStore returns the deposit store struct initialized with a.
func (k Backend[
	_, _, _, _, _, _, _, DepositStoreT, _, _, _, _, _, _, _, _, _,
]) DepositStore(
	_ context.Context,
) DepositStoreT {
	return k.ds
}
