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
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

func NewBackend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT types.RawBeaconBlockBody,
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
	DepositT Deposit,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT nodetypes.Node,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
](
	node NodeT,
	sb *storage.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		StateStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	],
) Backend[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, StateStoreT, ValidatorT,
	WithdrawalT, WithdrawalCredentialsT,
] {
	return &backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		NodeT, StateStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]{
		node:    node,
		Backend: sb,
	}
}

type backend[
	AvailabilityStoreT AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT types.RawBeaconBlockBody,
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
	DepositT Deposit,
	DepositStoreT DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT nodetypes.Node,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT Validator[WithdrawalCredentialsT],
	WithdrawalT Withdrawal[WithdrawalT],
	WithdrawalCredentialsT WithdrawalCredentials,
] struct {
	*storage.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BeaconStateMarshallableT, BlobSidecarsT, BlockStoreT,
		DepositT, DepositStoreT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		StateStoreT, ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]
	node NodeT
}

func (b *backend[
	_, _, _, _, _, _, _, _, _, _, _, _, _, NodeT, _, _, _, _,
]) AttachNode(node NodeT) {
	b.node = node
}

// StateFromContext returns a state from the context.
// It wraps the StorageBackend.StateFromContext method to allow for state ID
// querying.
func (b *backend[
	_, _, _, _, BeaconStateT, _, _, _, _, _, _, _, _, _, _, _, _, _,
]) StateFromContext(
	_ context.Context,
	stateID string,
) (BeaconStateT, error) {
	var state BeaconStateT
	height, err := heightFromStateID(stateID)
	if err != nil {
		return state, err
	}
	queryCtx, err := b.node.CreateQueryContext(height, false)
	if err != nil {
		panic(err)
	}
	//nolint: contextcheck // we want this context
	return b.Backend.StateFromContext(queryCtx), nil
}
