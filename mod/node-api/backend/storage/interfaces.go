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

	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core"
	"github.com/berachain/beacon-kit/mod/state-transition/pkg/core/state"
)

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
] interface {
	AttachNode(node NodeT)

	AvailabilityStore(ctx context.Context) AvailabilityStoreT
	BlockStore() BlockStoreT
	StateFromContext(ctx context.Context, stateID string) (BeaconStateT, error)
}
