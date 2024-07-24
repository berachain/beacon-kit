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
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
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
	BeaconBlockHeaderT core.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT core.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, StateStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	ContextT context.Context,
	DepositT Deposit,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT Node[ContextT],
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	sb   StorageBackendT
	cs   common.ChainSpec
	node NodeT
}

// New creates and returns a new Backend instance.
// TODO: need to add state_id resolver; possible values are: "head" (canonical
// head in node's view), "genesis", "finalized", "justified", <slot>, <hex
// encoded stateRoot with 0x prefix>.
func New[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT any,
	BeaconBlockHeaderT core.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT core.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, StateStoreT, ValidatorT, WithdrawalT,
	],
	BeaconStateMarshallableT state.BeaconStateMarshallable[
		BeaconStateMarshallableT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	BlobSidecarsT any,
	BlockStoreT BlockStore[BeaconBlockT],
	ContextT context.Context,
	DepositT Deposit,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT Node[ContextT],
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	StorageBackendT StorageBackend[
		AvailabilityStoreT, BeaconStateT, BlockStoreT, DepositStoreT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
](
	storageBackend StorageBackendT,
	cs common.ChainSpec,
) *Backend[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
	ContextT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
	NodeT, StateStoreT, StorageBackendT, ValidatorT, WithdrawalT,
	WithdrawalCredentialsT,
] {
	return &Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		ContextT, DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		NodeT, StateStoreT, StorageBackendT, ValidatorT, WithdrawalT,
		WithdrawalCredentialsT,
	]{
		sb: storageBackend,
		cs: cs,
	}
}

func (b *Backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, _, NodeT, _, _, _, _, _,
]) AttachNode(node NodeT) {
	b.node = node
}

// stateFromSlot returns the state at the given slot using query context.
func (b *Backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) stateFromSlot(
	slot uint64,
) (BeaconStateT, error) {
	var state BeaconStateT
	//#nosec:G701 // not an issue in practice.
	queryCtx, err := b.node.CreateQueryContext(int64(slot), false)
	if err != nil {
		return state, err
	}

	return b.sb.StateFromContext(queryCtx), nil
}
