// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
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

// New creates a new AttestationData.
func (a *AttestationData) New(
	slot math.U64,
	index math.U64,
	beaconBlockRoot common.Root,
) *AttestationData {
	a = &AttestationData{
		Slot:            slot,
		Index:           index,
		BeaconBlockRoot: beaconBlockRoot,
	}
	return a
}

// SizeSSZ returns the size of the AttestationData object in SSZ encoding.
func (*AttestationData) SizeSSZ() uint32 {
	return AttestationDataSize
}

// DefineSSZ defines the SSZ encoding for the AttestationData object.
func (a *AttestationData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &a.Slot)
	ssz.DefineUint64(codec, &a.Index)
	ssz.DefineStaticBytes(codec, &a.BeaconBlockRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the AttestationData object.
func (a *AttestationData) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(a), nil
}

// MarshalSSZ marshals the AttestationData object to SSZ format.
func (a *AttestationData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, a.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, a)
}

// UnmarshalSSZ unmarshals the AttestationData object from SSZ format.
func (a *AttestationData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, a)
}

// MarshalSSZTo marshals the AttestationData object into a pre-allocated byte slice.
func (a *AttestationData) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := a.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(dst, bz...), err
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
