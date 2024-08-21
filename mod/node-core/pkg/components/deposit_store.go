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
	storev2 "cosmossdk.io/store/v2/db"
	servertypes "github.com/berachain/beacon-kit/mod/consensus/pkg/cometbft/service/server/types"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/log"
	"github.com/berachain/beacon-kit/mod/node-core/pkg/components/storage"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	depositstore "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
)

// DepositStoreInput is the input for the dep inject framework.
type DepositStoreInput struct {
	depinject.In
	AppOpts servertypes.AppOptions
}

// ProvideDepositStore is a function that provides the module to the
// application.
func ProvideDepositStore[
	DepositT Deposit[
		DepositT, *ForkData, WithdrawalCredentials,
	],
](
	in DepositStoreInput,
) (*depositstore.KVStore[DepositT], error) {
	name := "deposits"
	dir := cast.ToString(in.AppOpts.Get(flags.FlagHome)) + "/data"
	kvp, err := storev2.NewDB(storev2.DBTypePebbleDB, name, dir, nil)
	if err != nil {
		return nil, err
	}

	return depositstore.NewStore[DepositT](storage.NewKVStoreProvider(kvp)), nil
}

// DepositPrunerInput is the input for the deposit pruner.
type DepositPrunerInput[
	BeaconBlockT any,
	DepositStoreT any,
	LoggerT any,
] struct {
	depinject.In
	ChainSpec    common.ChainSpec
	DepositStore DepositStoreT
	Dispatcher   *Dispatcher
	Logger       LoggerT
}

// ProvideDepositPruner provides a deposit pruner for the depinject framework.
func ProvideDepositPruner[
	BeaconBlockT BeaconBlock[
		BeaconBlockT, BeaconBlockBodyT, BeaconBlockHeaderT,
	],
	BeaconBlockBodyT interface {
		GetDeposits() []DepositT
	},
	BeaconBlockHeaderT any,
	DepositT Deposit[
		DepositT, *ForkData, WithdrawalCredentials,
	],
	DepositStoreT DepositStore[DepositT],
	LoggerT log.AdvancedLogger[any, LoggerT],
](
	in DepositPrunerInput[BeaconBlockT, DepositStoreT, LoggerT],
) (pruner.Pruner[DepositStoreT], error) {
	// initialize a subscription for finalized blocks.
	subFinalizedBlocks := make(chan async.Event[BeaconBlockT])
	if err := in.Dispatcher.Subscribe(
		async.BeaconBlockFinalizedEvent, subFinalizedBlocks,
	); err != nil {
		in.Logger.Error("failed to subscribe to event", "event",
			async.BeaconBlockFinalizedEvent, "err", err)
		return nil, err
	}

	return pruner.NewPruner[BeaconBlockT, DepositStoreT](
		in.Logger.With("service", manager.DepositPrunerName),
		in.DepositStore,
		manager.DepositPrunerName,
		subFinalizedBlocks,
		deposit.BuildPruneRangeFn[
			BeaconBlockT,
			BeaconBlockBodyT,
			DepositT,
			WithdrawalCredentials,
		](in.ChainSpec),
	), nil
}
