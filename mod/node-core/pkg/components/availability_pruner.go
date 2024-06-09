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
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	dastore "github.com/berachain/beacon-kit/mod/da/pkg/store"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/ethereum/go-ethereum/event"
)

// AvailabilityPrunerInput is the input for the ProviderAvailabilityPruner
// function for the depinject framework.
type AvailabilityPrunerInput struct {
	depinject.In
	Logger            log.Logger
	ChainSpec         primitives.ChainSpec
	BlockFeed         *event.FeedOf[feed.Event[*types.BeaconBlock]]
	AvailabilityStore *dastore.Store[types.BeaconBlockBody]
}

// ProvideAvailabilityPruner provides a availability pruner for the depinject
// framework.
func ProvideAvailabilityPruner(
	in AvailabilityPrunerInput,
) *pruner.DBPruner[
	*types.BeaconBlock,
	feed.Event[*types.BeaconBlock],
	pruner.Prunable,
	event.Subscription,
] {
	rangeDB, _ := in.AvailabilityStore.IndexDB.(*filedb.RangeDB)

	return pruner.NewPruner[
		*types.BeaconBlock,
		feed.Event[*types.BeaconBlock],
		pruner.Prunable,
		event.Subscription,
	](
		in.Logger.With("service", manager.AvailabilityPrunerName),
		rangeDB,
		manager.AvailabilityPrunerName,
		in.BlockFeed,
		dastore.BuildPruneRangeFn[
			*types.BeaconBlock,
			feed.Event[*types.BeaconBlock],
		](in.ChainSpec),
	)
}
