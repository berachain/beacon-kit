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

package utils

import (
	"fmt"

	cometbft "github.com/berachain/beacon-kit/consensus/cometbft/service"
	"github.com/berachain/beacon-kit/errors"
	handlertypes "github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	statedb "github.com/berachain/beacon-kit/state-transition/core/state"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func MapStateIDToStateAndSlot[
	StorageBackendT interface {
		StateAtSlot(slot math.Slot) (*statedb.StateDB, math.Slot, error)
		GetSlotByStateRoot(root common.Root) (math.Slot, error)
	},
](backend StorageBackendT, stateID string) (*statedb.StateDB, math.Slot, error) {
	slot, err := slotFromStateID(stateID, backend)
	if err != nil {
		return nil, 0, err
	}
	st, resolvedSlot, err := backend.StateAtSlot(slot)
	if err != nil {
		if errors.Is(err, cometbft.ErrAppNotReady) {
			// chain not ready, like when genesis time is set in the future
			return nil, 0, handlertypes.ErrNotFound
		}
		if errors.Is(err, sdkerrors.ErrInvalidHeight) {
			// height requested too high
			return nil, 0, handlertypes.ErrNotFound
		}
		return nil, 0, fmt.Errorf("failed to get state from stateID %s, slot %d: %w", stateID, slot, err)
	}
	return st, resolvedSlot, nil
}

var errNoSlotForStateRoot = errors.New("slot not found at state root")

// TODO: define unique types for each of the query-able IDs (state & block from
// spec, execution unique to beacon-kit). For each type define validation
// functions and resolvers to slot number.

// SlotFromStateID returns a slot from the state ID.
//
// NOTE: Right now, `stateID` only supports querying by "head" (all of "head",
// "finalized", "justified" are the same), "genesis", and <slot>.
func slotFromStateID[
	StorageBackendT interface {
		GetSlotByStateRoot(root common.Root) (math.Slot, error)
	},
](stateID string, storage StorageBackendT) (math.Slot, error) {
	if slot, err := mapStateIDToSlot(stateID); err == nil {
		return slot, nil
	}

	// We assume that the state ID is a state hash.
	root, err := common.NewRootFromHex(stateID)
	if err != nil {
		return 0, err
	}
	slot, err := storage.GetSlotByStateRoot(root)
	if err != nil {
		return 0, errNoSlotForStateRoot
	}
	return slot, nil
}
