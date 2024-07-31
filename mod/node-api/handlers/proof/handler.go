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
	"github.com/berachain/beacon-kit/mod/node-api/handlers"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/mod/node-api/server/context"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// Handler is the handler for the proof API.
type Handler[
	ContextT context.Context,
	BeaconBlockHeaderT types.BeaconBlockHeader,
	BeaconStateT types.BeaconState[
		BeaconStateMarshallableT, ExecutionPayloadHeaderT, ValidatorT,
	],
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader,
	ValidatorT types.Validator,
] struct {
	*handlers.BaseHandler[ContextT]
	backend Backend[BeaconBlockHeaderT, BeaconStateT, ValidatorT]
}

// NewHandler creates a new handler for the proof API.
func NewHandler[
	ContextT context.Context,
	BeaconBlockHeaderT types.BeaconBlockHeader,
	BeaconStateT types.BeaconState[
		BeaconStateMarshallableT, ExecutionPayloadHeaderT, ValidatorT,
	],
	BeaconStateMarshallableT types.BeaconStateMarshallable,
	ExecutionPayloadHeaderT types.ExecutionPayloadHeader,
	ValidatorT types.Validator,
](
	backend Backend[BeaconBlockHeaderT, BeaconStateT, ValidatorT],
) *Handler[
	ContextT, BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
	ExecutionPayloadHeaderT, ValidatorT,
] {
	h := &Handler[
		ContextT, BeaconBlockHeaderT, BeaconStateT, BeaconStateMarshallableT,
		ExecutionPayloadHeaderT, ValidatorT,
	]{
		BaseHandler: handlers.NewBaseHandler(
			handlers.NewRouteSet[ContextT](""),
		),
		backend: backend,
	}
	return h
}

// Get the slot from the given input of timestamp id, beacon state, and beacon
// block header for the resolved slot.
func (h *Handler[
	ContextT, BeaconBlockHeaderT, BeaconStateT, _, _, _,
]) resolveTimestampID(timestampID string) (
	math.Slot, BeaconStateT, BeaconBlockHeaderT, error,
) {
	var (
		beaconState BeaconStateT
		blockHeader BeaconBlockHeaderT
	)

	slot, err := utils.SlotFromTimestampID(timestampID, h.backend)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	beaconState, slot, err = h.backend.StateFromSlotForProof(slot)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	blockHeader, err = h.backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	return slot, beaconState, blockHeader, nil
}
