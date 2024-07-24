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

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// BeaconBlockHeaderBase represents the base of a beacon block header.
type BeaconBlockHeaderBase struct {
	// Slot represents the position of the block in the chain.
	// TODO: Put back to math.Slot after fastssz fixes.
	Slot uint64
	// ProposerIndex is the index of the validator who proposed the block.
	// TODO: Put back to math.ProposerIndex after fastssz fixes.
	ProposerIndex uint64
	// ParentBlockRoot is the hash of the parent block
	ParentBlockRoot common.Root
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeaderBase) GetSlot() math.Slot {
	return math.Slot(b.Slot)
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeaderBase) GetProposerIndex() math.ValidatorIndex {
	return math.ValidatorIndex(b.ProposerIndex)
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockBase.
func (b *BeaconBlockHeaderBase) GetParentBlockRoot() common.Root {
	return b.ParentBlockRoot
}

// GetStateRoot retrieves the state root of the BeaconBlockDeneb.
func (b *BeaconBlockHeaderBase) GetStateRoot() common.Root {
	return b.StateRoot
}

// BeaconBlockHeader is the header of a beacon block.
type BeaconBlockHeader struct {
	// BeaconBlockHeaderBase is the base of the block.
	BeaconBlockHeaderBase
	// 	// BodyRoot is the root of the block body.
	BodyRoot common.Root `json:"bodyRoot"`
}

// NewBeaconBlockHeader creates a new BeaconBlockHeader.
func NewBeaconBlockHeader(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	stateRoot common.Root,
	bodyRoot common.Root,
) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		BeaconBlockHeaderBase: BeaconBlockHeaderBase{
			Slot:            uint64(slot),
			ProposerIndex:   uint64(proposerIndex),
			ParentBlockRoot: parentBlockRoot,
			StateRoot:       stateRoot,
		},
		BodyRoot: bodyRoot,
	}
}

// New creates a new BeaconBlockHeader.
func (b *BeaconBlockHeader) New(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	stateRoot common.Root,
	bodyRoot common.Root,
) *BeaconBlockHeader {
	return NewBeaconBlockHeader(
		slot, proposerIndex, parentBlockRoot, stateRoot, bodyRoot,
	)
}

// SetStateRoot sets the state root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetStateRoot(stateRoot common.Root) {
	b.StateRoot = stateRoot
}
