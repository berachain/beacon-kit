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
	"github.com/berachain/beacon-kit/mod/async/pkg/publisher"
	"github.com/berachain/beacon-kit/mod/async/pkg/server"
	asynctypes "github.com/berachain/beacon-kit/mod/async/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/events"
)

// ProvideEventServer provides an event server.
func ProvideEventServer() *EventServer {
	return server.NewEventServer()
}

// ProvidePublishers provides a publisher for beacon block
// finalized events.
func ProvidePublishers() []asynctypes.Publisher {
	return []asynctypes.Publisher{
		publisher.New[GenesisDataReceivedEvent](
			events.GenesisDataReceived,
		),
		publisher.New[GenesisDataProcessedEvent](
			events.GenesisDataProcessed,
		),
		publisher.New[NewSlotEvent](
			events.NewSlot,
		),
		publisher.New[BuiltBeaconBlockEvent](
			events.BuiltBeaconBlock,
		),
		publisher.New[BuiltSidecarsEvent](
			events.BuiltSidecars,
		),
		publisher.New[BeaconBlockReceivedEvent](
			events.BeaconBlockReceived,
		),
		publisher.New[SidecarsReceivedEvent](
			events.SidecarsReceived,
		),
		publisher.New[BeaconBlockVerifiedEvent](
			events.BeaconBlockVerified,
		),
		publisher.New[SidecarsVerifiedEvent](
			events.SidecarsVerified,
		),
		publisher.New[FinalBeaconBlockReceivedEvent](
			events.FinalBeaconBlockReceived,
		),
		publisher.New[FinalSidecarsReceivedEvent](
			events.FinalSidecarsReceived,
		),
		publisher.New[FinalValidatorUpdatesProcessedEvent](
			events.FinalValidatorUpdatesProcessed,
		),
		publisher.New[FinalizedBlockEvent](
			events.BeaconBlockFinalizedEvent,
		),
	}
}
