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
	"strings"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
)

// SlotFromBlockID returns a slot from the block ID.
//
// NOTE: `blockID` shares the same semantics as `stateID`, with the modification
// of being able to query by beacon <blockRoot> instead of <stateRoot>.
func SlotFromBlockID[StorageBackendT interface {
	GetSlotByBlockRoot(root common.Root) (math.Slot, error)
}](blockID string, storage StorageBackendT) (math.Slot, error) {
	if slot, err := mapStateIDToSlot(blockID); err == nil {
		return slot, nil
	}

	// We assume that the block ID is a block hash.
	root, err := common.NewRootFromHex(blockID)
	if err != nil {
		return 0, err
	}
	return storage.GetSlotByBlockRoot(root)
}

// ParentSlotFromTimestampID returns the parent slot corresponding to the
// timestamp ID.
//
// NOTE: `timestampID` shares the same semantics as `stateID`, with the
// modification of being able to query by next block's <timestamp> instead of
// the current block's <stateRoot>.
//
// The <timestamp> must be prefixed by the 't', followed by the timestamp
// in decimal UNIX notation. For example 't1728681738' corresponds to the slot
// which has the next block with a timestamp of 1728681738. Providing just the
// string '1728681738' (without the prefix 't') will query for the beacon block
// for slot 1728681738.
func ParentSlotFromTimestampID[StorageBackendT interface {
	GetParentSlotByTimestamp(timestamp math.U64) (math.Slot, error)
}](timestampID string, storage StorageBackendT) (math.Slot, error) {
	if !IsTimestampIDPrefix(timestampID) {
		return mapStateIDToSlot(timestampID)
	}

	// Parse the timestamp from the timestampID.
	timestamp, err := math.U64FromString(timestampID[1:])
	if err != nil {
		return 0, errors.Wrapf(
			err, "failed to parse timestamp from timestampID: %s", timestampID,
		)
	}
	return storage.GetParentSlotByTimestamp(timestamp)
}

// IsTimestampIDPrefix checks if the given timestampID is prefixed with the
// correct prefix 't'.
func IsTimestampIDPrefix(timestampID string) bool {
	return strings.HasPrefix(timestampID, TimestampIDPrefix)
}

// slotFromStateID returns a slot number from the given state ID.
// Currently, when "genesis" is requested, we return the block 1 state.
// This is due to a CometBFT limitation that does not explicitly commit the
// genesis state (it accumulates block 1 state changes and flushes them together).
// Numeric requests are clamped so that slot 0 maps to Genesis (slot 1).
// TODO: Properly return the true genesis state when requested, instead of block 1.
func mapStateIDToSlot(id string) (math.Slot, error) {
	switch id {
	case StateIDFinalized, StateIDJustified, StateIDHead:
		return Head, nil
	case StateIDGenesis:
		return Genesis, nil
	default:
		slot, err := math.U64FromString(id)
		if err != nil {
			return math.Slot(0), errors.Wrapf(err, "failed mapping stateID %q to slot", id)
		}

		// Enforce here too the choice of mapping genesis to slot 1. Without this
		// the request of slot 0 would be mapped to tip of the chain.
		return max(slot, Genesis), nil
	}
}
