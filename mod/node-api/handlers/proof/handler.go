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

package proof

import (
	"errors"

	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/types/context"
	types "github.com/berachain/beacon-kit/mod/node-api/types/proof"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
)

// Handler is the handler for the proof API.
type Handler[
	ContextT context.Context,
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT types.BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT,
	],
	Eth1DataT constraints.SSZMarshallable,
	ExecutionPayloadHeaderT constraints.SSZMarshallable,
	ForkT constraints.SSZMarshallable,
	ValidatorT any,
] struct {
	backend Backend[
		BeaconBlockHeaderT, BeaconStateT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	]
	routes *handlers.RouteSet[ContextT]
}

// NewHandler creates a new handler for the proof API.
func NewHandler[
	ContextT context.Context,
	BeaconBlockHeaderT types.BeaconBlockHeader[BeaconBlockHeaderT],
	BeaconStateT types.BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT,
	],
	Eth1DataT constraints.SSZMarshallable,
	ExecutionPayloadHeaderT constraints.SSZMarshallable,
	ForkT constraints.SSZMarshallable,
	ValidatorT any,
](
	backend Backend[
		BeaconBlockHeaderT, BeaconStateT, Eth1DataT, ExecutionPayloadHeaderT,
		ForkT, ValidatorT,
	],
) *Handler[
	ContextT, BeaconBlockHeaderT, BeaconStateT, Eth1DataT,
	ExecutionPayloadHeaderT, ForkT, ValidatorT,
] {
	h := &Handler[
		ContextT, BeaconBlockHeaderT, BeaconStateT, Eth1DataT,
		ExecutionPayloadHeaderT, ForkT, ValidatorT,
	]{
		backend: backend,
		routes:  handlers.NewRouteSet[ContextT](""),
	}
	return h
}

func (
	h *Handler[ContextT, _, _, _, _, _, _],
) RouteSet() *handlers.RouteSet[ContextT] {
	return h.routes
}

// NotImplemented is a placeholder for the proof API.
func (
	h *Handler[ContextT, _, _, _, _, _, _],
) NotImplemented(_ ContextT) (any, error) {
	return nil, errors.New("not implemented")
}
