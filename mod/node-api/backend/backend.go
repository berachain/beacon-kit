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
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/backend/storage"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

type Backend[
	AvailabilityStoreT storage.AvailabilityStore[
		BeaconBlockBodyT, BlobSidecarsT,
	],
	BeaconBlockT any,
	BeaconBlockBodyT types.RawBeaconBlockBody,
	BeaconBlockHeaderT core.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT core.BeaconState[
		BeaconStateT, BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, StateStoreT, ValidatorT, WithdrawalT,
	],
	BlobSidecarsT any,
	BlockStoreT storage.BlockStore[BeaconBlockT],
	DepositT storage.Deposit,
	DepositStoreT storage.DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT nodetypes.Node,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT storage.Validator[WithdrawalCredentialsT],
	WithdrawalT storage.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT storage.WithdrawalCredentials,
] struct {
	storage.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, StateStoreT,
		ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]
	cs common.ChainSpec
}

// New creates and returns a new Backend instance.
// TODO: need to add state_id resolver; possible values are: "head" (canonical
// head in node's view), "genesis", "finalized", "justified", <slot>, <hex
// encoded stateRoot with 0x prefix>.
func New[
	AvailabilityStoreT storage.AvailabilityStore[
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
	BlockStoreT storage.BlockStore[BeaconBlockT],
	DepositT storage.Deposit,
	DepositStoreT storage.DepositStore[DepositT],
	Eth1DataT,
	ExecutionPayloadHeaderT,
	ForkT any,
	NodeT nodetypes.Node,
	StateStoreT state.KVStore[
		StateStoreT, BeaconBlockHeaderT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	],
	ValidatorT storage.Validator[WithdrawalCredentialsT],
	WithdrawalT storage.Withdrawal[WithdrawalT],
	WithdrawalCredentialsT storage.WithdrawalCredentials,
](
	storageBackend storage.Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, StateStoreT,
		ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	],
	cs common.ChainSpec,
) *Backend[
	AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
	Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, StateStoreT,
	ValidatorT, WithdrawalT, WithdrawalCredentialsT,
] {
	return &Backend[
		AvailabilityStoreT, BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
		BeaconStateT, BlobSidecarsT, BlockStoreT, DepositT, DepositStoreT,
		Eth1DataT, ExecutionPayloadHeaderT, ForkT, NodeT, StateStoreT,
		ValidatorT, WithdrawalT, WithdrawalCredentialsT,
	]{
		Backend: storageBackend,
		cs:      cs,
	}
}
