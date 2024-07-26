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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// BeaconBlock represents a block in the beacon chain during
// the Deneb fork.
type BeaconBlock struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot `json:"slot"`
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.Slot `json:"proposer_index"`
	// ParentRoot is the hash of the parent block
	ParentRoot common.Root `json:"parent_root"`
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root `json:"state_root"`
	// Body is the body of the BeaconBlock, containing the block's
	// operations.
	Body *BeaconBlockBody `json:"body"`
}

// Empty creates an empty beacon block.
func (b *BeaconBlock) Empty(forkVersion uint32) *BeaconBlock {
	switch forkVersion {
	case version.Deneb:
		return &BeaconBlock{}
	case version.DenebPlus:
		panic("unsupported fork version")
	default:
		panic("fork version not supported")
	}
}

// NewWithVersion assembles a new beacon block from the given.
func (b *BeaconBlock) NewWithVersion(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	forkVersion uint32,
) (*BeaconBlock, error) {
	var (
		block *BeaconBlock
	)

	switch forkVersion {
	case version.Deneb:
		block = &BeaconBlock{
			Slot:          slot,
			ProposerIndex: proposerIndex,
			ParentRoot:    parentBlockRoot,
			StateRoot:     bytes.B32{},
			Body:          &BeaconBlockBody{},
		}
	default:
		return &BeaconBlock{}, ErrForkVersionNotSupported
	}

	return block, nil
}

// NewFromSSZ creates a new beacon block from the given SSZ bytes.
func (b *BeaconBlock) NewFromSSZ(
	bz []byte,
	forkVersion uint32,
) (*BeaconBlock, error) {
	var block = new(BeaconBlock)
	switch forkVersion {
	case version.Deneb:
		block = &BeaconBlock{}
	case version.DenebPlus:
		panic("unsupported fork version")
	default:
		return block, ErrForkVersionNotSupported
	}

	return block, block.UnmarshalSSZ(bz)
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlock object in SSZ encoding.
func (b *BeaconBlock) SizeSSZ(fixed bool) uint32 {
	//nolint:mnd // todo fix.
	var size = uint32(8 + 8 + 32 + 32 + 4)
	if fixed {
		return size
	}
	size += ssz.SizeDynamicObject(b.Body)
	return size
}

// DefineSSZ defines the SSZ encoding for the BeaconBlock object.
func (b *BeaconBlock) DefineSSZ(codec *ssz.Codec) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentRoot)
	ssz.DefineStaticBytes(codec, &b.StateRoot)
	ssz.DefineDynamicObjectOffset(codec, &b.Body)

	// Define the dynamic data (fields)
	ssz.DefineDynamicObjectContent(codec, &b.Body)
}

// MarshalSSZ marshals the BeaconBlock object to SSZ format.
func (b *BeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, b.SizeSSZ(false))
	return buf, ssz.EncodeToBytes(buf, b)
}

// UnmarshalSSZ unmarshals the BeaconBlock object from SSZ format.
func (b *BeaconBlock) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// HashTreeRoot computes the Merkleization of the BeaconBlock object.
func (b *BeaconBlock) HashTreeRoot() ([32]byte, error) {
	return ssz.HashConcurrent(b), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the BeaconBlock object to the provided buffer in SSZ
// format.
func (b *BeaconBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := b.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the BeaconBlock object with a hasher.
func (b *BeaconBlock) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (1) 'ProposerIndex'
	hh.PutUint64(uint64(b.ProposerIndex))

	// Field (2) 'ParentBlockRoot'
	hh.PutBytes(b.ParentRoot[:])

	// Field (3) 'StateRoot'
	hh.PutBytes(b.StateRoot[:])

	// Field (4) 'Body'
	if err := b.Body.HashTreeRootWith(hh); err != nil {
		return err
	}

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the BeaconBlock object.
func (b *BeaconBlock) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(b)
}

// IsNil checks if the beacon block is nil.
func (b *BeaconBlock) IsNil() bool {
	return b == nil
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlock) GetSlot() math.Slot {
	return b.Slot
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlock) GetProposerIndex() math.ValidatorIndex {
	return b.ProposerIndex
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockBase.
func (b *BeaconBlock) GetParentBlockRoot() common.Root {
	return b.ParentRoot
}

// GetStateRoot retrieves the state root of the BeaconBlock.
func (b *BeaconBlock) GetStateRoot() common.Root {
	return b.StateRoot
}

// Version identifies the version of the BeaconBlock.
func (b *BeaconBlock) Version() uint32 {
	return version.Deneb
}

// SetStateRoot sets the state root of the BeaconBlock.
func (b *BeaconBlock) SetStateRoot(root common.Root) {
	b.StateRoot = root
}

// GetBody retrieves the body of the BeaconBlock.
func (b *BeaconBlock) GetBody() *BeaconBlockBody {
	return b.Body
}

// GetHeader builds a BeaconBlockHeader from the BeaconBlock.
func (b *BeaconBlock) GetHeader() *BeaconBlockHeader {
	bodyRoot, err := b.GetBody().HashTreeRoot()
	if err != nil {
		return nil
	}

	return &BeaconBlockHeader{
		Slot:            b.Slot,
		ProposerIndex:   b.ProposerIndex,
		ParentBlockRoot: b.ParentRoot,
		StateRoot:       b.StateRoot,
		BodyRoot:        bodyRoot,
	}
}
