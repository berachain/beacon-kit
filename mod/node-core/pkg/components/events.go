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
	"github.com/berachain/beacon-kit/mod/async/pkg/broker"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/async"
)

// ProvidePublishers provides a publisher for beacon block
// finalized events.
func ProvidePublishers[
	BeaconBlockT any,
	BlobSidecarsT any,
	GenesisDataT any,
]() []asynctypes.Broker {
	return []asynctypes.Broker{
		broker.New[async.Event[GenesisDataT]](
			async.GenesisDataReceived,
		),
		broker.New[GenesisDataProcessedEvent](
			async.GenesisDataProcessed,
		),
		broker.New[NewSlotEvent](
			async.NewSlot,
		),
		broker.New[async.Event[BeaconBlockT]](
			async.BuiltBeaconBlock,
		),
		broker.New[async.Event[BlobSidecarsT]](
			async.BuiltSidecars,
		),
		broker.New[async.Event[BeaconBlockT]](
			async.BeaconBlockReceived,
		),
		broker.New[async.Event[BlobSidecarsT]](
			async.SidecarsReceived,
		),
		broker.New[async.Event[BeaconBlockT]](
			async.BeaconBlockVerified,
		),
		broker.New[async.Event[BlobSidecarsT]](
			async.SidecarsVerified,
		),
		broker.New[async.Event[BeaconBlockT]](
			async.FinalBeaconBlockReceived,
		),
		broker.New[async.Event[BlobSidecarsT]](
			async.FinalSidecarsReceived,
		),
		broker.New[FinalValidatorUpdatesProcessedEvent](
			async.FinalValidatorUpdatesProcessed,
		),
		broker.New[async.Event[BeaconBlockT]](
			async.BeaconBlockFinalizedEvent,
		),
	}
}
