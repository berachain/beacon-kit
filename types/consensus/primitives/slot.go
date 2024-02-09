// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package primitives

import (
	"errors"
	"fmt"

	fssz "github.com/prysmaticlabs/fastssz"
)

// eight is the fixed size of a Slot in bytes.
const eight = 8

var (
	// Ensure Slot implements the fssz.HashRoot interface.
	_ fssz.HashRoot = (Slot)(0)
	// Ensure Slot implements the fssz.Marshaler interface.
	_ fssz.Marshaler = (*Slot)(nil)
	// Ensure Slot implements the fssz.Unmarshaler interface.
	_ fssz.Unmarshaler = (*Slot)(nil)
	// ErrInvalidBufferSize is an error indicating that the provided buffer size is invalid.
	ErrInvalidBufferSize = errors.New("invalid buffer size")
)

// Slot represents a single slot in the blockchain context, implemented as a uint64.
type Slot uint64

// HashTreeRoot computes the hash tree root of the Slot.
// It returns a fixed-size byte array and an error if any.
func (s Slot) HashTreeRoot() ([32]byte, error) {
	return fssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith computes the hash tree root of the Slot using the provided hasher.
// It modifies the hasher's state and returns an error if any.
func (s Slot) HashTreeRootWith(hh *fssz.Hasher) error {
	hh.PutUint64(uint64(s))
	return nil
}

// UnmarshalSSZ unmarshals a Slot from SSZ-encoded data. It returns an error
// if the buffer size is incorrect.
func (s *Slot) UnmarshalSSZ(buf []byte) error {
	if len(buf) != s.SizeSSZ() {
		return fmt.Errorf(
			"%w: expected length %d, received %d",
			ErrInvalidBufferSize, s.SizeSSZ(), len(buf),
		)
	}
	*s = Slot(fssz.UnmarshallUint64(buf))
	return nil
}

// MarshalSSZTo appends the SSZ-encoded byte slice of Slot to the provided destination slice.
// It returns the resulting slice and an error if any occurs during marshaling.
func (s *Slot) MarshalSSZTo(dst []byte) ([]byte, error) {
	marshalled, err := s.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	return append(dst, marshalled...), nil
}

// MarshalSSZ marshals a Slot into a byte slice. It returns the byte slice and an error if any.
func (s *Slot) MarshalSSZ() ([]byte, error) {
	marshalled := fssz.MarshalUint64([]byte{}, uint64(*s))
	return marshalled, nil
}

// SizeSSZ returns the fixed SSZ size of a Slot, which is 8 bytes.
func (s *Slot) SizeSSZ() int {
	return eight
}
