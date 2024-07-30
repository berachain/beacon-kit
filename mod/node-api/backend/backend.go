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

package backend

import (
	"context"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
)

// Backend is the db access layer for the beacon node-api.
// It serves as a wrapper around the storage backend and provides an abstraction
// over building the query context for a given state.
type Backend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	ContextT context.Context,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT Node[ContextT],
	StateStoreT any,
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	ValidatorsT ~[]ValidatorT,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	sb   StorageBackendT
	cs   common.ChainSpec
	node NodeT
}

// New creates and returns a new Backend instance.
func New[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BeaconBlockHeaderT BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT, ValidatorsT, WithdrawalT,
	],
	BeaconStateMarshallableT any,
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	ContextT context.Context,
	DepositT any,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT Node[ContextT],
	StateStoreT any,
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	ValidatorsT ~[]ValidatorT,
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
](
	storageBackend StorageBackendT, cs common.ChainSpec,
) *Backend[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
	ContextT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
	NodeT, StateStoreT, StorageBackendT, ValidatorT, ValidatorsT, WithdrawalT,
	WithdrawalCredentialsT,
] {
	return &Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		ContextT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		NodeT, StateStoreT, StorageBackendT, ValidatorT, ValidatorsT, WithdrawalT,
		WithdrawalCredentialsT,
	]{
		sb: storageBackend,
		cs: cs,
	}
}

func (b *Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, NodeT, _, _, _, _, _, _,
]) AttachNode(node NodeT) {
	b.node = node
}

// ChainSpec returns the chain spec from the backend.
func (b *Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, NodeT, _, _, _, _, _, _,
]) ChainSpec() common.ChainSpec {
	return b.cs
}

func (b *Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) GetSlotByRoot(root [32]byte) (uint64, error) {
	return b.sb.BlockStore().GetSlotByRoot(root)
}
