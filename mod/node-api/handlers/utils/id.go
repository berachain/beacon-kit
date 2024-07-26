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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
)

const (
	StateIDGenesis   = "genesis"
	StateIDFinalized = "finalized"
	StateIDJustified = "justified"
	StateIDHead      = "head"
)

const (
	Head uint64 = iota
	Genesis
)

// SlotFromStateID returns a slot from the state ID.
//
// NOTE: Right now, `stateID` only supports querying by "head" (all of "head",
// "finalized", "justified" are the same), "genesis", and <slot>. We do NOT
// support querying by <stateRoot>.
func SlotFromStateID(stateID string) (uint64, error) {
	switch stateID {
	case StateIDFinalized, StateIDJustified, StateIDHead:
		return Head, nil
	case StateIDGenesis:
		return Genesis, nil
	default:
		slot, err := hex.String(stateID).ToUint64()
		if err != nil {
			return 0, err
		}
		//#nosec:G701 // not an issue in practice.
		return slot, nil
	}
}

// BlockID shares the same semantics as StateID, with the addition of
// being able to query state by block hash.
func SlotFromBlockID[StorageBackendT interface {
	GetSlotByRoot(root [32]byte) (uint64, error)
}](blockID string, storage StorageBackendT) (uint64, error) {
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
