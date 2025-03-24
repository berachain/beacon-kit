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

package types

import (
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// AttestationDataSize is the size of the AttestationData object in bytes.
// 8 bytes for Slot + 8 bytes for Index + 32 bytes for BeaconBlockRoot.
const AttestationDataSize = 48

var (
	_ ssz.StaticObject                    = (*AttestationData)(nil)
	_ constraints.SSZMarshallableRootable = (*AttestationData)(nil)
)

// AttestationData represents an attestation data.
type AttestationData struct {
	// Slot is the slot number of the attestation data.
	Slot math.U64 `json:"slot"`
	// Index is the index of the validator.
	Index math.U64 `json:"index"`
	// BeaconBlockRoot is the root of the beacon block.
	BeaconBlockRoot common.Root `json:"beaconBlockRoot"`
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the AttestationData object in SSZ encoding.
func (*AttestationData) SizeSSZ(*ssz.Sizer) uint32 {
	return AttestationDataSize
}

// DefineSSZ defines the SSZ encoding for the AttestationData object.
func (a *AttestationData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &a.Slot)
	ssz.DefineUint64(codec, &a.Index)
	ssz.DefineStaticBytes(codec, &a.BeaconBlockRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the AttestationData object.
func (a *AttestationData) HashTreeRoot() common.Root {
	return ssz.HashSequential(a)
}

// MarshalSSZ marshals the AttestationData object to SSZ format.
func (a *AttestationData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(a))
	return buf, ssz.EncodeToBytes(buf, a)
}

// UnmarshalSSZ unmarshals the AttestationData object from SSZ format.
func (a *AttestationData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, a)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo marshals the AttestationData object into a pre-allocated byte
// slice.
func (a *AttestationData) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := a.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, err
}

// HashTreeRootWith ssz hashes the AttestationData object with a hasher.
func (a *AttestationData) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Slot'
	hh.PutUint64(uint64(a.Slot))

	// Field (1) 'Index'
	hh.PutUint64(uint64(a.Index))

	// Field (2) 'BeaconBlockRoot'
	hh.PutBytes(a.BeaconBlockRoot[:])

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the AttestationData object.
func (a *AttestationData) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(a)
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// GetSlot returns the slot of the attestation data.
func (a *AttestationData) GetSlot() math.U64 {
	return a.Slot
}

// GetIndex returns the index of the attestation data.
func (a *AttestationData) GetIndex() math.U64 {
	return a.Index
}

// GetBeaconBlockRoot returns the beacon block root of the attestation data.
func (a *AttestationData) GetBeaconBlockRoot() common.Root {
	return a.BeaconBlockRoot
}
