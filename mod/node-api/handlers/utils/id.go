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

package utils

import (
	"strings"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SlotFromStateID returns a slot from the state ID.
//
// NOTE: Right now, `stateID` only supports querying by "head" (all of "head",
// "finalized", "justified" are the same), "genesis", and <slot>.
func SlotFromStateID[StorageBackendT interface {
	GetSlotByStateRoot(root common.Root) (math.Slot, error)
}](stateID string, storage StorageBackendT) (math.Slot, error) {
	if slot, err := slotFromID(stateID); err == nil {
		return slot, nil
	}

	// We assume that the state ID is a state hash.
	root, err := common.NewRootFromHex(stateID)
	if err != nil {
		return 0, err
	}
	return storage.GetSlotByStateRoot(root)
}

// SlotFromBlockID returns a slot from the block ID.
//
// NOTE: `blockID` shares the same semantics as `stateID`, with the modification
// of being able to query by beacon <blockRoot> instead of <stateRoot>.
func SlotFromBlockID[StorageBackendT interface {
	GetSlotByBlockRoot(root common.Root) (math.Slot, error)
}](blockID string, storage StorageBackendT) (math.Slot, error) {
	if slot, err := slotFromID(blockID); err == nil {
		return slot, nil
	}

	// We assume that the block ID is a block hash.
	root, err := common.NewRootFromHex(blockID)
	if err != nil {
		return 0, err
	}
	return storage.GetSlotByBlockRoot(root)
}

// SlotFromExecutionID returns a slot from the execution number ID.
//
// NOTE: `executionID` shares the same semantics as `stateID`, with the
// modification of being able to query by beacon block <executionNumber>
// instead of <stateRoot>.
//
// The <executionNumber> must be prefixed by the 'n', followed by the execution
// number in hexadecimal notation. For example 'n0x66aab3ef' corresponds to
// the slot with execution number 1722463215. Providing just the string
// '0x66aab3ef' (without the prefix 'n') will query for the beacon block with
// slot 1722463215.
func SlotFromExecutionID[StorageBackendT interface {
	GetSlotByExecutionNumber(executionNumber math.U64) (math.Slot, error)
}](executionID string, storage StorageBackendT) (math.Slot, error) {
	if !IsExecutionNumberPrefix(executionID) {
		return slotFromID(executionID)
	}

	// Parse the execution number from the executionID.
	executionNumber, err := U64FromString(executionID[1:])
	if err != nil {
		return 0, err
	}
	return storage.GetSlotByExecutionNumber(executionNumber)
}

// IsExecutionNumberPrefix checks if the given executionID is prefixed
// with the execution number prefix.
func IsExecutionNumberPrefix(executionID string) bool {
	return strings.HasPrefix(executionID, ExecutionIDPrefix)
}

// U64FromString returns a math.U64 from the given string. Errors if the given
// string is not in proper hexadecimal notation.
func U64FromString(id string) (math.U64, error) {
	var u64 math.U64
	return u64, u64.UnmarshalText([]byte(id))
}

func slotFromID(id string) (math.Slot, error) {
	switch id {
	case StateIDFinalized, StateIDJustified, StateIDHead:
		return Head, nil
	case StateIDGenesis:
		return Genesis, nil
	default:
		return U64FromString(id)
	}
}
