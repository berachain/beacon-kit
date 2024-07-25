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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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
	Slot math.Slot
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex
	// ParentBlockRoot is the hash of the parent block
	ParentBlockRoot common.Root
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root
	// BodyRoot is the root of the block body.
	BodyRoot common.Root
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
func (b *BeaconBlockHeader) SizeSSZ() uint32 {
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
	buf := make([]byte, b.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the BeaconBlockBody object from SSZ format.
func (b *BeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the BeaconBlockHeader object.
func (b *BeaconBlockHeader) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(b), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZToBytes marshals the BeaconBlockHeader object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZTo(buf []byte) ([]byte, error) {
	bz, err := b.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(buf, bz...), nil
}

// HashTreeRootWith ssz hashes the BeaconBlockHeader object with a hasher
func (b *BeaconBlockHeader) HashTreeRootWith(hh fastssz.HashWalker) (err error) {
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
	return
}

// GetTree ssz hashes the BeaconBlockHeader object
func (b *BeaconBlockHeader) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

/* -------------------------------------------------------------------------- */
/*                            Getters and Setters                             */
/* -------------------------------------------------------------------------- */

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

// TODO: Deprecate
//
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
