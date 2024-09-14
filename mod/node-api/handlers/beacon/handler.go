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

package beacon

import (
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/server/context"
)

// Handler is the handler for the beacon API.
type Handler[
	BeaconBlockHeaderT types.BeaconBlockHeader,
	ContextT context.Context,
	ForkT types.Fork,
	ValidatorT types.Validator,
] struct {
	*handlers.BaseHandler[ContextT]
	backend Backend[BeaconBlockHeaderT, ForkT, any]
}

// NewHandler creates a new handler for the beacon API.
func NewHandler[
	BeaconBlockHeaderT types.BeaconBlockHeader,
	ContextT context.Context,
	ForkT types.Fork,
	ValidatorT types.Validator,
](
	backend Backend[BeaconBlockHeaderT, ForkT, ValidatorT],
) *Handler[BeaconBlockHeaderT, ContextT, ForkT, ValidatorT] {
	h := &Handler[BeaconBlockHeaderT, ContextT, ForkT, ValidatorT]{
		BaseHandler: handlers.NewBaseHandler(
			handlers.NewRouteSet[ContextT](""),
		),
		backend: backend,
	}
	return h
}
