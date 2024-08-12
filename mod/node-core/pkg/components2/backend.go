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
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// StorageBackendInput is the input for the ProvideStorageBackend function.
type StorageBackendInput[
	AvailabilityStoreT any,
	BlockStoreT any,
	DepositStoreT any,
	KVStoreT any,
] struct {
	depinject.In
	AvailabilityStore AvailabilityStoreT
	BlockStore        BlockStoreT
	ChainSpec         common.ChainSpec
	DepositStore      DepositStoreT
	KVStore           KVStoreT
}

// ProvideStorageBackend is the depinject provider that returns a beacon storage
// backend.
func ProvideStorageBackend[
	AvailabilityStoreT AvailabilityStore[BeaconBlockBodyT, BlobSidecarsT],
	BeaconBlockT any,
	BeaconBlockBodyT constraints.SSZMarshallable,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, KVStoreT, ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	DepositT Deposit[ForkDataT, WithdrawalCredentialsT],
	DepositStoreT DepositStore[DepositT],
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ForkDataT any,
	KVStoreT BeaconStore[
		KVStoreT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT, ValidatorsT, WithdrawalT,
	],
	ValidatorT any,
	ValidatorsT any,
	WithdrawalT any,
	WithdrawalCredentialsT ~[32]byte,
](
	in StorageBackendInput[
		AvailabilityStoreT, BlockStoreT, DepositStoreT, KVStoreT,
	],
) *StorageBackend {
	return storage.NewBackend[
		AvailabilityStoreT,
		BeaconBlockT,
		BeaconBlockBodyT,
		BeaconBlockHeaderT,
		BeaconStateT,
		BeaconStateMarshallableT,
		BlobSidecarsT,
		BlockStoreT,
		DepositT,
		DepositStoreT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		KVStoreT,
		ValidatorT,
		ValidatorsT,
		WithdrawalT,
		WithdrawalCredentialsT,
	](
		in.ChainSpec,
		in.AvailabilityStore,
		in.KVStore,
		in.DepositStore,
		in.BlockStore,
	)
}
