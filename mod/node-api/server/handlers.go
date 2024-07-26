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

package server

import (
	"github.com/berachain/beacon-kit/mod/node-api/backend"
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/beacon"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/builder"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/config"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/debug"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/events"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/node"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/node-api/types/context"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// DefaultHandlers returns the default handlers for the node API.
func DefaultHandlers[
	ContextT context.Context,
	BeaconBlockHeaderT backend.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT types.BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT,
	],
	Eth1DataT constraints.SSZRootable,
	ExecutionPayloadHeaderT constraints.SSZRootable,
	ForkT constraints.SSZRootable,
	ValidatorT any,
](
	backend Backend[
		BeaconBlockHeaderT, BeaconStateT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	],
) []handlers.Handlers[ContextT] {
	return []handlers.Handlers[ContextT]{
		beacon.NewHandler[ContextT](backend),
		builder.NewHandler[ContextT](),
		config.NewHandler[ContextT](),
		debug.NewHandler[ContextT](),
		events.NewHandler[ContextT](),
		node.NewHandler[ContextT](),
		proof.NewHandler[ContextT](backend),
	}
}
