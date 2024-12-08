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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// BeaconBlockHeaderSize is the size of the BeaconBlockHeader object in bytes.
//
// Total size: Slot (8) + ProposerIndex (8) +
// ParentBlockRoot (32) + StateRoot (32) + BodyRoot (32).
const BeaconBlockHeaderSize = 112

var (
	_ ssz.StaticObject                    = (*BeaconBlockHeader)(nil)
	_ constraints.SSZMarshallableRootable = (*BeaconBlockHeader)(nil)
)

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

// Empty creates an empty BeaconBlockHeader instance.
func (*BeaconBlockHeader) Empty() *BeaconBlockHeader {
	return &BeaconBlockHeader{}
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

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlockHeader object in SSZ encoding.
func (b *BeaconBlockHeader) SizeSSZ(*ssz.Sizer) uint32 {
	return BeaconBlockHeaderSize
}

// DefineSSZ defines the SSZ encoding for the BeaconBlockHeader object.
func (b *BeaconBlockHeader) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentBlockRoot)
	ssz.DefineStaticBytes(codec, &b.StateRoot)
	ssz.DefineStaticBytes(codec, &b.BodyRoot)
}

// MarshalSSZ marshals the BeaconBlockBody object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the BeaconBlockBody object from SSZ format.
func (b *BeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the BeaconBlockHeader object.
func (b *BeaconBlockHeader) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the BeaconBlockHeader object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := b.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
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
