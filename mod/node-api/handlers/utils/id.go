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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// SlotFromStateID returns a slot from the state ID.
//
// NOTE: Right now, `stateID` only supports querying by "head" (all of "head",
// "finalized", "justified" are the same), "genesis", and <slot>. We do NOT
// support querying by <stateRoot>.
func SlotFromStateID(stateID string) (math.Slot, error) {
	switch stateID {
	case StateIDFinalized, StateIDJustified, StateIDHead:
		return Head, nil
	case StateIDGenesis:
		return Genesis, nil
	default:
		return U64FromString(stateID)
	}
}

// SlotFromBlockID returns a slot from the block ID.
//
// NOTE: `blockID` shares the same semantics as `stateID`, with the modification
// of being able to query by beacon <blockRoot> instead of <stateRoot>.
func SlotFromBlockID[StorageBackendT interface {
	GetSlotByRoot(root common.Root) (math.Slot, error)
}](blockID string, storage StorageBackendT) (math.Slot, error) {
	if slot, err := SlotFromStateID(blockID); err == nil {
		return slot, nil
	}

	// We assume that the block ID is a block hash.
	root, err := hex.String(blockID).ToBytes()
	if err != nil {
		return 0, err
	}
	return storage.GetSlotByRoot(bytes.ToBytes32(root))
}

// SlotFromTimestampID returns a slot from the timestamp ID.
//
// NOTE: `timestampID` shares the same semantics as `stateID`, with the
// modification of being able to query by beacon block <timestamp> instead of
// <stateRoot>.
//
// The <timestamp> must be prefixed by the 't', followed by the Unix timestamp
// number in hexadecimal notation. For example 't0x66aab3ef' corresponds to
// the block at Unix timestamp 1722463215. Providing just the string
// '0x66aab3ef' (without the prefix 't') will query for the beacon block with
// slot 1722463215.
func SlotFromTimestampID[StorageBackendT interface {
	GetSlotByTimestamp(timestamp math.U64) (math.Slot, error)
}](timestampID string, storage StorageBackendT) (math.Slot, error) {
	if !IsTimestampPrefix(timestampID) {
		return SlotFromStateID(timestampID)
	}

	// Parse the timestamp from the timestampID.
	timestamp, err := U64FromString(timestampID[1:])
	if err != nil {
		return 0, err
	}
	return storage.GetSlotByTimestamp(timestamp)
}

// IsTimestampPrefix checks if the given timestampID is prefixed with the
// timestamp prefix.
func IsTimestampPrefix(timestampID string) bool {
	return strings.HasPrefix(timestampID, TimestampIDPrefix)
}

// U64FromString returns a math.U64 from the given string. Errors if the given
// string is not in proper hexadecimal notation.
func U64FromString(id string) (math.U64, error) {
	var u64 math.U64
	return u64, u64.UnmarshalText([]byte(id))
}
