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

package internal

import (
	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/karalabe/ssz"
)

type BeaconBlockHeader struct {
	*types.BeaconBlockHeader
}

// SizeSSZ returns the size of the BeaconBlockHeader object in SSZ encoding.
func (b *BeaconBlockHeader) SizeSSZ() uint32 {
	return 112 // Total size: Slot (8) + ProposerIndex (8) + ParentBlockRoot (32) + StateRoot (32) + BodyRoot (32)
}

// DefineSSZ defines the SSZ encoding for the BeaconBlockHeader object.
func (b *BeaconBlockHeader) DefineSSZ(codec *ssz.Codec) {
	if b.BeaconBlockHeader == nil {
		b.BeaconBlockHeader = &types.BeaconBlockHeader{}
	}
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentBlockRoot)
	ssz.DefineStaticBytes(codec, &b.StateRoot)
	ssz.DefineStaticBytes(codec, &b.BodyRoot)
}

// MarshalSSZToBytes marshals the BeaconBlockHeader object to SSZ format.
func (b *BeaconBlockHeader) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, b)
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

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeader) GetSlot() math.Slot {
	return math.Slot(b.Slot)
}

// GetSlot retrieves the slot of the BeaconBlockBase.
func (b *BeaconBlockHeader) GetProposerIndex() math.ValidatorIndex {
	return math.ValidatorIndex(b.ProposerIndex)
}

// GetParentBlockRoot retrieves the parent block root of the BeaconBlockBase.
func (b *BeaconBlockHeader) GetParentBlockRoot() common.Root {
	return b.ParentBlockRoot
}

// GetStateRoot retrieves the state root of the BeaconBlockDeneb.
func (b *BeaconBlockHeader) GetStateRoot() common.Root {
	return b.StateRoot
}

// SetStateRoot sets the state root of the BeaconBlockHeader.
func (b *BeaconBlockHeader) SetStateRoot(stateRoot common.Root) {
	b.StateRoot = stateRoot
}
