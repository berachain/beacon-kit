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
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// BeaconBlockForValidator represents a block in the beacon chain with the
// minimally required values to prove an element in the state exists.
type BeaconBlockForStateProof[
	BeaconBlockHeaderT any,
	BeaconStateT BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT,
	],
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ValidatorT any,
] struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex
	// ParentRoot is the hash of the parent block.
	ParentRoot common.Root
	// State is full BeaconState type to prove elements inside the state.
	State BeaconStateT
	// BodyRoot is the root of the block body.
	BodyRoot common.Root
}

// NewBeaconBlockForStateProof creates a new BeaconBlock SSZ summary with only
// the required raw values to prove an element in the beacon state exists in
// this block.
func NewBeaconBlockForStateProof[
	BeaconBlockHeaderT any,
	BeaconStateT BeaconState[
		BeaconBlockHeaderT, Eth1DataT, ExecutionPayloadHeaderT, ForkT,
		ValidatorT,
	],
	Eth1DataT any,
	ExecutionPayloadHeaderT any,
	ForkT any,
	ValidatorT any,
](
	bbh BeaconBlockHeader[BeaconBlockHeaderT],
	bsv BeaconStateT,
) (
	*BeaconBlockForStateProof[
		BeaconBlockHeaderT,
		BeaconStateT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
	], error) {
	return &BeaconBlockForStateProof[
		BeaconBlockHeaderT,
		BeaconStateT,
		Eth1DataT,
		ExecutionPayloadHeaderT,
		ForkT,
		ValidatorT,
	]{
		Slot:          bbh.GetSlot(),
		ProposerIndex: bbh.GetProposerIndex(),
		ParentRoot:    bbh.GetParentBlockRoot(),
		State:         bsv,
		BodyRoot:      bbh.GetBodyRoot(),
	}, nil
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the BeaconBlockForStateProof object in SSZ 
// encoding.
func (b *BeaconBlockForStateProof[_, _, _, _, _, _]) SizeSSZ(fixed bool) uint32 {
	//nolint:mnd // todo fix.
	var size = uint32(8 + 8 + 32 + 32 + 4)
	if fixed {
		return size
	}
	size += ssz.SizeDynamicObject(b.State)
	return size
}

// DefineSSZ defines the SSZ encoding for the BeaconBlock object.
func (b *BeaconBlockForStateProof[_, BeaconStateT, _, _, _, _]) DefineSSZ(
	codec *ssz.Codec,
) {
	// Define the static data (fields and dynamic offsets)
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentRoot)
	ssz.DefineDynamicObjectOffset(codec, &b.State)
	ssz.DefineStaticBytes(codec, &b.BodyRoot)
	

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

// HashTreeRoot computes the Merkleization of the BeaconBlockForStateProof object.
func (b *BeaconBlockForStateProof[_, _, _, _, _, _]) HashTreeRoot() (
	[32]byte, error,
) {
	return ssz.HashConcurrent(b), nil
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// GetTree ssz hashes the BeaconBlockForStateProof object.
func (b *BeaconBlockForStateProof[_, _, _, _, _, _]) GetTree() (
	*fastssz.Node, error,
) {
	return fastssz.ProofTree(b)
}

// HashTreeRootWith ssz hashes the BeaconBlockForStateProof object with a hasher.
func (b *BeaconBlockForStateProof[_, _, _, _, _, _]) HashTreeRootWith(
	hh fastssz.HashWalker,
) error {
	indx := hh.Index()

	// Field (0) 'Slot'
	hh.PutUint64(uint64(b.Slot))

	// Field (1) 'ProposerIndex'
	hh.PutUint64(uint64(b.ProposerIndex))

	// Field (2) 'ParentRoot'
	hh.PutBytes(b.ParentRoot[:])

	// Field (3) 'State'
	if err := b.State.HashTreeRootWith(hh); err != nil {
		return err
	}

	// Field (4) 'BodyRoot'
	hh.PutBytes(b.BodyRoot[:])

	hh.Merkleize(indx)
	return nil
}
