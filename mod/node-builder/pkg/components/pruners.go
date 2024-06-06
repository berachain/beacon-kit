// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is govered by the Business Source License included
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
	"cosmossdk.io/core/log"
	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/mod/beacon/blockchain"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/events"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	datypes "github.com/berachain/beacon-kit/mod/da/pkg/types"
	"github.com/berachain/beacon-kit/mod/execution/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/primitives"
	depositdb "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/ethereum/go-ethereum/event"
)

// PrunerInput is the input for the pruners through the dep inject framework.
type PrunerInput struct {
	depinject.In
	ChainSpec      primitives.ChainSpec
	StorageBackend blockchain.StorageBackend[
		*dastore.Store[types.BeaconBlockBody],
		types.BeaconBlockBody,
		BeaconState,
		*datypes.BlobSidecars,
		*types.Deposit,
		*depositdb.KVStore[*types.Deposit],
	]
	BlockFeed *event.FeedOf[events.Block[*types.BeaconBlock]]
	Logger    log.Logger
}

// TODO: this hood af xD
// ProvidePruners provides the pruners for db manager through the dep inject
// framework.
func ProvidePruners(
	in PrunerInput,
) []*pruner.Pruner[
	*types.BeaconBlock,
	events.Block[*types.BeaconBlock],
	event.Subscription] {
	// Build the deposit pruner.
	depositPruner := pruner.NewPruner[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription,
	](
		in.Logger.With("service", manager.DepositPrunerName),
		in.StorageBackend.DepositStore(nil),
		manager.DepositPrunerName,
		in.BlockFeed,
		deposit.BuildPruneRangeFn[
			types.BeaconBlockBody,
			*types.BeaconBlock,
			events.Block[*types.BeaconBlock],
			*types.Deposit,
			*types.ExecutionPayload,
			types.WithdrawalCredentials,
		](in.ChainSpec),
	)

	// build the availability pruner if IndexDB is available.
	avs := in.StorageBackend.AvailabilityStore(nil).IndexDB
	availabilityPruner := pruner.NewPruner[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription,
	](
		in.Logger.With("service", manager.AvailabilityPrunerName),
		avs.(*filedb.RangeDB),
		manager.AvailabilityPrunerName,
		in.BlockFeed,
		dastore.BuildPruneRangeFn[
			*types.BeaconBlock,
			events.Block[*types.BeaconBlock],
		](in.ChainSpec),
	)

	// slice of pruners to pass to the DBManager.
	return []*pruner.Pruner[
		*types.BeaconBlock,
		events.Block[*types.BeaconBlock],
		event.Subscription]{depositPruner, availabilityPruner}
}
