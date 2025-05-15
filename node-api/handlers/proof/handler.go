// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
)

// Handler is the handler for the proof API.
type Handler struct {
	*handlers.BaseHandler
	backend Backend
}

// NewHandler creates a new handler for the proof API.
func NewHandler(backend Backend) *Handler {
	h := &Handler{
		BaseHandler: handlers.NewBaseHandler(
			handlers.NewRouteSet(""),
		),
		backend: backend,
	}
	return h
}

// Get the slot from the given input of timestamp id, beacon state, and beacon
// block header for the resolved slot.
func (h *Handler) resolveTimestampID(timestampID string) (
	math.Slot, *statedb.StateDB, *ctypes.BeaconBlockHeader, error,
) {
	var (
		beaconState *statedb.StateDB
		blockHeader *ctypes.BeaconBlockHeader
	)

	slot, err := utils.ParentSlotFromTimestampID(timestampID, h.backend)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	beaconState, slot, err = h.backend.StateAtSlot(slot)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	blockHeader, err = h.backend.BlockHeaderAtSlot(slot)
	if err != nil {
		return 0, beaconState, blockHeader, err
	}

	return slot, beaconState, blockHeader, nil
}
