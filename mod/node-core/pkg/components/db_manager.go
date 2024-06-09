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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/feed"
	dastore "github.com/berachain/beacon-kit/mod/storage/pkg/deposit"
	"github.com/berachain/beacon-kit/mod/storage/pkg/filedb"
	"github.com/berachain/beacon-kit/mod/storage/pkg/manager"
	"github.com/berachain/beacon-kit/mod/storage/pkg/pruner"
	"github.com/ethereum/go-ethereum/event"
)

// DBManagerInput is the input for the dep inject framework.
type DBManagerInput struct {
	depinject.In
	Logger             log.Logger
	DepositPruner      pruner.Pruner[*dastore.KVStore[*types.Deposit]]
	AvailabilityPruner pruner.Pruner[*filedb.RangeDB]
}

// ProvideDBManager provides a DBManager for the depinject framework.
func ProvideDBManager(
	in DBManagerInput,
) (*manager.DBManager[*types.BeaconBlock,
	*feed.Event[*types.BeaconBlock],
	event.Subscription,
], error) {
	return manager.NewDBManager[
		*types.BeaconBlock,
		*feed.Event[*types.BeaconBlock],
		event.Subscription,
	](
		in.Logger.With("service", "db-manager"),
		in.DepositPruner,
		in.AvailabilityPruner,
	)
}
