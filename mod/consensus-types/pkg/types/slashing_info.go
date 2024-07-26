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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/constraints"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/karalabe/ssz"
)

// SlashingInfoSize is the size of the SlashingInfo object in SSZ encoding.
const SlashingInfoSize = 16 // 8 bytes for Slot + 8 bytes for Index

// Compile-time assertions to ensure SlashingInfo implements the correct
// interfaces.
var (
	_ ssz.StaticObject                    = (*SlashingInfo)(nil)
	_ constraints.SSZMarshallableRootable = (*SlashingInfo)(nil)
)

// Compile-time assertion to ensure SlashingInfoSize matches the SizeSSZ method.
var _ = [1]struct{}{}[16-SlashingInfoSize]

// SlashingInfo represents a slashing info.
type SlashingInfo struct {
	// Slot is the slot number of the slashing info.
	Slot math.Slot
	// ValidatorIndex is the validator index of the slashing info.
	Index math.U64
}

/* -------------------------------------------------------------------------- */
/*                                 Constructor                                */
/* -------------------------------------------------------------------------- */

// New creates a new slashing info instance.
func (s *SlashingInfo) New(slot, index math.U64) *SlashingInfo {
	s = &SlashingInfo{
		Slot:  slot,
		Index: index,
	}
	return s
}

/* -------------------------------------------------------------------------- */
/*                                     SSZ                                    */
/* -------------------------------------------------------------------------- */

// SizeSSZ returns the size of the SlashingInfo object in SSZ encoding.
func (*SlashingInfo) SizeSSZ() uint32 {
	return SlashingInfoSize
}

// DefineSSZ defines the SSZ encoding for the SlashingInfo object.
func (s *SlashingInfo) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &s.Slot)
	ssz.DefineUint64(codec, &s.Index)
}

// HashTreeRoot computes the SSZ hash tree root of the SlashingInfo object.
func (s *SlashingInfo) HashTreeRoot() ([32]byte, error) {
	return ssz.HashSequential(s), nil
}

// MarshalSSZ marshals the SlashingInfo object to SSZ format.
func (s *SlashingInfo) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, s.SizeSSZ())
	return buf, ssz.EncodeToBytes(buf, s)
}

// UnmarshalSSZ unmarshals the SlashingInfo object from SSZ format.
func (s *SlashingInfo) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

/* -------------------------------------------------------------------------- */
/*                                   FastSSZ                                  */
/* -------------------------------------------------------------------------- */

// MarshalSSZTo ssz marshals the SlashingInfo object into a pre-allocated byte
// slice.
func (s *SlashingInfo) MarshalSSZTo(dst []byte) ([]byte, error) {
	bz, err := s.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	dst = append(dst, bz...)
	return dst, nil
}

// HashTreeRootWith ssz hashes the SlashingInfo object with a hasher.
func (s *SlashingInfo) HashTreeRootWith(hh fastssz.HashWalker) error {
	indx := hh.Index()

	// Field (0) 'Slot'
	hh.PutUint64(uint64(s.Slot))

	// Field (1) 'Index'
	hh.PutUint64(uint64(s.Index))

	hh.Merkleize(indx)
	return nil
}

// GetTree ssz hashes the SlashingInfo object.
func (s *SlashingInfo) GetTree() (*fastssz.Node, error) {
	return fastssz.ProofTree(s)
}

/* -------------------------------------------------------------------------- */
/*                             Getters and Setters                            */
/* -------------------------------------------------------------------------- */

// GetSlot returns the slot of the slashing info.
func (s *SlashingInfo) GetSlot() math.Slot {
	return s.Slot
}

// GetIndex returns the index of the slashing info.
func (s *SlashingInfo) GetIndex() math.U64 {
	return s.Index
}

// SetSlot sets the slot of the slashing info.
func (s *SlashingInfo) SetSlot(slot math.Slot) {
	s.Slot = slot
}

// SetIndex sets the index of the slashing info.
func (s *SlashingInfo) SetIndex(index math.U64) {
	s.Index = index
}
