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

// TODO: BeaconBlockHeader needs manual fastssz migration to handle dual interface compatibility
// go:generate sszgen -path . -objs BeaconBlockHeader -output header_sszgen.go -include ../../primitives/common,../../primitives/math

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
)

// BeaconBlockHeaderSize is the size of the BeaconBlockHeader object in bytes.
//
// Total size: Slot (8) + ProposerIndex (8) +
// ParentBlockRoot (32) + StateRoot (32) + BodyRoot (32).
const BeaconBlockHeaderSize = 112

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZMarshallableRootable = (*BeaconBlockHeader)(nil)
// )

// BeaconBlockHeader represents the base of a beacon block header.
type BeaconBlockHeader struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot `json:"slot"`
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex `json:"proposer_index"`
	// ParentBlockRoot is the hash of the parent block
	ParentBlockRoot common.Root `json:"parent_block_root"`
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root `json:"state_root"`
	// BodyRoot is the root of the block body.
	BodyRoot common.Root `json:"body_root"`
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

// NewBeaconBlockHeader creates a new BeaconBlockHeader.
func NewBeaconBlockHeader(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	stateRoot common.Root,
	bodyRoot common.Root,
) *BeaconBlockHeader {
	return &BeaconBlockHeader{
		Slot:            slot,
		ProposerIndex:   proposerIndex,
		ParentBlockRoot: parentBlockRoot,
		StateRoot:       stateRoot,
		BodyRoot:        bodyRoot,
	}
}

func NewEmptyBeaconBlockHeader() *BeaconBlockHeader {
	return &BeaconBlockHeader{}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlockHeader object in SSZ encoding.
func (b *BeaconBlockHeader) SizeSSZ() int {
	return BeaconBlockHeaderSize
}


// MarshalSSZ marshals the BeaconBlockHeader object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, BeaconBlockHeaderSize)
	return b.MarshalSSZTo(buf)
}

func (*BeaconBlockHeader) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot computes the SSZ hash tree root of the BeaconBlockHeader object.
func (b *BeaconBlockHeader) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	b.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the BeaconBlockHeader object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Field (0) 'Slot'
	dst = fastssz.MarshalUint64(dst, uint64(b.Slot))

	// Field (1) 'ProposerIndex'
	dst = fastssz.MarshalUint64(dst, uint64(b.ProposerIndex))

	// Field (2) 'ParentBlockRoot'
	dst = append(dst, b.ParentBlockRoot[:]...)

	// Field (3) 'StateRoot'
	dst = append(dst, b.StateRoot[:]...)

	// Field (4) 'BodyRoot'
	dst = append(dst, b.BodyRoot[:]...)

	return dst, nil
}

// HashTreeRootWith ssz hashes the BeaconBlockHeader object with a hasher.
func (b *BeaconBlockHeader) HashTreeRootWith(
	hh fastssz.HashWalker,
) error {
	indx := hh.Index()

	// Field (0) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (1) 'ProposerIndex'
	hh.PutUint64(uint64(b.ProposerIndex))

	// Field (2) 'ParentBlockRoot'
	hh.PutBytes(b.ParentBlockRoot[:])

	// Field (3) 'StateRoot'
	hh.PutBytes(b.StateRoot[:])

	// Field (4) 'BodyRoot'
	hh.PutBytes(b.BodyRoot[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconBlockHeader object.
func (b *BeaconBlockHeader) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

// UnmarshalSSZ ssz unmarshals the BeaconBlockHeader object.
func (b *BeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	if len(buf) != BeaconBlockHeaderSize {
		return fastssz.ErrSize
	}

	// Field (0) 'Slot'
	b.Slot = math.Slot(fastssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'ProposerIndex'
	b.ProposerIndex = math.ValidatorIndex(fastssz.UnmarshallUint64(buf[8:16]))

	// Field (2) 'ParentBlockRoot'
	copy(b.ParentBlockRoot[:], buf[16:48])

	// Field (3) 'StateRoot'
	copy(b.StateRoot[:], buf[48:80])

	// Field (4) 'BodyRoot'
	copy(b.BodyRoot[:], buf[80:112])

	return nil
}


/* -------------------------------------------------------------------------- */
/*                            Getters and Setters                             */
/* -------------------------------------------------------------------------- */

// Equals returns true if the Withdrawal is equal to the other.
func (b *BeaconBlockHeader) Equals(rhs *BeaconBlockHeader) bool {
	switch {
	case b == nil && rhs == nil:
		return true
	case b != nil && rhs != nil:
		return b.Slot == rhs.Slot &&
			b.ProposerIndex == rhs.ProposerIndex &&
			b.ParentBlockRoot == rhs.ParentBlockRoot &&
			b.StateRoot == rhs.StateRoot &&
			b.BodyRoot == rhs.BodyRoot
	default:
		return false
	}
}

// GetSlot retrieves the slot of the BeaconBlockHeader.
func (b *BeaconBlockHeader) GetSlot() math.Slot {
	return b.Slot
}

// SetSlot sets the slot of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetSlot(slot math.Slot) {
	b.Slot = slot
}

// GetProposerIndex retrieves the proposer index of the BeaconBlockHeader.
func (b *BeaconBlockHeader) GetProposerIndex() math.ValidatorIndex {
	return b.ProposerIndex
}

// SetProposerIndex sets the proposer index of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetProposerIndex(
	proposerIndex math.ValidatorIndex,
) {
	b.ProposerIndex = proposerIndex
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) GetParentBlockRoot() common.Root {
	return b.ParentBlockRoot
}

// SetParentBlockRoot sets the parent block root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetParentBlockRoot(parentBlockRoot common.Root) {
	b.ParentBlockRoot = parentBlockRoot
}

// GetStateRoot retrieves the state root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) GetStateRoot() common.Root {
	return b.StateRoot
}

// SetStateRoot sets the state root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetStateRoot(stateRoot common.Root) {
	b.StateRoot = stateRoot
}

// GetBodyRoot retrieves the body root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) GetBodyRoot() common.Root {
	return b.BodyRoot
}

// SetBodyRoot sets the body root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetBodyRoot(bodyRoot common.Root) {
	b.BodyRoot = bodyRoot
}
