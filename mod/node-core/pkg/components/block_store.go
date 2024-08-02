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
	"cosmossdk.io/log"
	storev2 "cosmossdk.io/store/v2/db"
	blockservice "github.com/berachain/beacon-kit/mod/beacon/block_store"
	"github.com/berachain/beacon-kit/mod/config"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/storage/pkg/block"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/spf13/cast"
)

// BlockStoreInput is the input for the dep inject framework.
type BlockStoreInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// ProvideBlockStore is a function that provides the module to the
// application.
func ProvideBlockStore(
	in BlockStoreInput,
) (*BlockStore, error) {
	name := "blocks"
	dir := cast.ToString(in.AppOpts.Get(flags.FlagHome)) + "/data"
	kvp, err := storev2.NewDB(storev2.DBTypePebbleDB, name, dir, nil)
	if err != nil {
		return nil, err
	}

	return block.NewStore[*BeaconBlock](storage.NewKVStoreProvider(kvp)), nil
}

// BlockPrunerInput is the input for the block pruner.
type BlockPrunerInput struct {
	depinject.In

	BlockBroker *BlockBroker
	BlockStore  *BlockStore
	Config      *config.Config
	Logger      log.Logger
}

// ProvideBlockPruner provides a block pruner for the depinject framework.
func ProvideBlockPruner(
	in BlockPrunerInput,
) (BlockPruner, error) {
	subCh, err := in.BlockBroker.Subscribe()
	if err != nil {
		in.Logger.Error("failed to subscribe to block feed", "err", err)
		return nil, err
	}

	return pruner.NewPruner[
		*BeaconBlock,
		*BlockEvent,
		*BlockStore,
	](
		in.Logger.With("service", manager.BlockPrunerName),
		in.BlockStore,
		manager.BlockPrunerName,
		subCh,
		blockservice.BuildPruneRangeFn[
			*BeaconBlock,
			*BlockEvent,
		](in.Config.BlockStoreService),
	), nil
}
