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
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	blockstore "github.com/berachain/beacon-kit/mod/beacon/block_store"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/log"
)

// BlockServiceInput is the input for the block service.
type BlockServiceInput[
	BeaconBlockT any,
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	LoggerT log.AdvancedLogger[any, LoggerT],
] struct {
	depinject.In

	BlockBroker *broker.Broker[*asynctypes.Event[BeaconBlockT]]
	BlockStore  BeaconBlockStoreT
	Config      *config.Config
	Logger      LoggerT
}

// ProvideBlockStoreService provides the block service.
func ProvideBlockStoreService[
	BeaconBlockT BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	],
	BeaconBlockBodyT any,
	BeaconBlockHeaderT any,
	BeaconBlockStoreT BlockStore[BeaconBlockT],
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	in BlockServiceInput[BeaconBlockT, BeaconBlockStoreT, LoggerT],
) *blockstore.Service[
	BeaconBlockT, BeaconBlockStoreT,
] {
	return blockstore.NewService(
		in.Config.BlockStoreService,
		in.Logger,
		in.BlockBroker,
		in.BlockStore,
	)
}
