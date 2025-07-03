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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types

import (
	"fmt"

	"github.com/berachain/beacon-kit/errors"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	fastssz "github.com/ferranbt/fastssz"
)

// TODO: Re-enable interface assertion once constraints are updated
// var (
// 	_ constraints.SSZVersionedMarshallableRootable = (*BeaconBlock)(nil)
// )

// BeaconBlock represents a block in the beacon chain.
type BeaconBlock struct {
	Versionable `json:"-"`

	// Slot represents the position of the block in the chain.
	Slot math.Slot `json:"slot"`
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex `json:"proposer_index"`
	// ParentRoot is the hash of the parent block
	ParentRoot common.Root `json:"parent_root"`
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root `json:"state_root"`
	// Body is the body of the BeaconBlock, containing the block's operations.
	Body *BeaconBlockBody `json:"body"`
}

// NewBeaconBlockWithVersion assembles a new beacon block from the given parameters.
func NewBeaconBlockWithVersion(
	slot math.Slot,
	proposerIndex math.ValidatorIndex,
	parentBlockRoot common.Root,
	forkVersion common.Version,
) (*BeaconBlock, error) {
	switch forkVersion {
	case version.Deneb(), version.Deneb1(), version.Electra(), version.Electra1():
		block := NewEmptyBeaconBlockWithVersion(forkVersion)
		block.Slot = slot
		block.ProposerIndex = proposerIndex
		block.ParentRoot = parentBlockRoot

		// StateRoot is left empty as it is not ready at this time.
		block.StateRoot = common.Root{}
		return block, nil
	default:
		// We return block here to appease nilaway.
		block := &BeaconBlock{}
		err := errors.Wrap(ErrForkVersionNotSupported, fmt.Sprintf("fork %d", forkVersion))
		return block, err
	}
}

func NewEmptyBeaconBlockWithVersion(version common.Version) *BeaconBlock {
	return &BeaconBlock{
		Versionable: NewVersionable(version),
		Body:        NewEmptyBeaconBlockBodyWithVersion(version),
	}
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlock object in SSZ encoding.
func (b *BeaconBlock) SizeSSZ() int {
	// Fixed part: slot(8) + proposerIndex(8) + parentRoot(32) + stateRoot(32) + offset(4) = 84
	return 84 + b.Body.SizeSSZ()
}


// MarshalSSZ marshals the BeaconBlock object to SSZ format.
func (b *BeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, 0, b.SizeSSZ())
	return b.MarshalSSZTo(buf)
}

func (b *BeaconBlock) ValidateAfterDecodingSSZ() error {
	return b.Body.ValidateAfterDecodingSSZ()
}

// HashTreeRoot computes the Merkleization of the BeaconBlock object.
func (b *BeaconBlock) HashTreeRoot() common.Root {
	hh := fastssz.DefaultHasherPool.Get()
	defer fastssz.DefaultHasherPool.Put(hh)
	b.HashTreeRootWith(hh)
	root, _ := hh.HashRoot()
	return common.Root(root)
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlock) GetSlot() math.Slot {
	return b.Slot
}

// GetProposerIndex retrieves the proposer index.
func (b *BeaconBlock) GetProposerIndex() math.ValidatorIndex {
	return b.ProposerIndex
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockBase.
func (b *BeaconBlock) GetParentBlockRoot() common.Root {
	return b.ParentRoot
}

// SetParentBlockRoot sets the parent block root of the BeaconBlockBase.
func (b *BeaconBlock) SetParentBlockRoot(parentBlockRoot common.Root) {
	b.ParentRoot = parentBlockRoot
}

// GetStateRoot retrieves the state root of the BeaconBlock.
func (b *BeaconBlock) GetStateRoot() common.Root {
	return b.StateRoot
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
	return &BeaconBlockHeader{
		Slot:            b.Slot,
		ProposerIndex:   b.ProposerIndex,
		ParentBlockRoot: b.ParentRoot,
		StateRoot:       b.StateRoot,
		BodyRoot:        b.GetBody().HashTreeRoot(),
	}
}

// GetTimestamp retrieves the timestamp of the BeaconBlock from
// the ExecutionPayload.
func (b *BeaconBlock) GetTimestamp() math.U64 {
	return b.Body.ExecutionPayload.Timestamp
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the BeaconBlock object to a target array.
func (b *BeaconBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	// Fixed part
	dst = fastssz.MarshalUint64(dst, uint64(b.Slot))
	dst = fastssz.MarshalUint64(dst, uint64(b.ProposerIndex))
	dst = append(dst, b.ParentRoot[:]...)
	dst = append(dst, b.StateRoot[:]...)

	// Offset for body
	offset := 84
	dst = fastssz.MarshalUint32(dst, uint32(offset))

	// Dynamic part: Body
	dst, err := b.Body.MarshalSSZTo(dst)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// UnmarshalSSZ ssz unmarshals the BeaconBlock object.
func (b *BeaconBlock) UnmarshalSSZ(buf []byte) error {
	if len(buf) < 84 {
		return fastssz.ErrSize
	}

	// Fixed part
	b.Slot = math.Slot(fastssz.UnmarshallUint64(buf[0:8]))
	b.ProposerIndex = math.ValidatorIndex(fastssz.UnmarshallUint64(buf[8:16]))
	copy(b.ParentRoot[:], buf[16:48])
	copy(b.StateRoot[:], buf[48:80])

	// Check offset
	offset := fastssz.UnmarshallUint32(buf[80:84])
	if offset != 84 {
		return fastssz.ErrInvalidVariableOffset
	}

	// Dynamic part: Body
	if b.Body == nil {
		b.Body = &BeaconBlockBody{}
	}
	if err := b.Body.UnmarshalSSZ(buf[84:]); err != nil {
		return err
	}

	return b.ValidateAfterDecodingSSZ()
}


// HashTreeRootWith ssz hashes the BeaconBlock object with a hasher.
func (b *BeaconBlock) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (1) 'ProposerIndex'
	hh.PutUint64(uint64(b.ProposerIndex))

	// Field (2) 'ParentRoot'
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
