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

package primitives

import (
	byteslib "github.com/itsdevbear/bolaris/lib/bytes"
	fssz "github.com/prysmaticlabs/fastssz"
)

var (
	// Ensure SSZUInt256 implements the fssz.HashRoot interface.
	_ fssz.HashRoot = SSZUInt256{}
	// Ensure SSZUInt256 implements the fssz.Marshaler interface.
	_ fssz.Marshaler = (*SSZUInt256)(nil)
	// Ensure SSZUInt256 implements the fssz.Unmarshaler interface.
	_ fssz.Unmarshaler = (*SSZUInt256)(nil)
)

const thirtyTwo = 32

// SSZUInt256 represents a ssz-able uint64.
type SSZUInt256 []byte

// SizeSSZ returns the fixed SSZ size of a SSZUInt256, which is 32 bytes.
func (s *SSZUInt256) SizeSSZ() int {
	return thirtyTwo
}

// MarshalSSZTo appends the SSZ-encoded byte slice of SSZUInt256 to the provided
// destination slice. It returns the resulting slice and an error if any occurs
// during marshaling.
func (s *SSZUInt256) MarshalSSZTo(dst []byte) ([]byte, error) {
	marshalled, err := s.MarshalSSZ()
	if err != nil {
		return nil, err
	}
	padded := byteslib.ToBytes32(marshalled)
	return append(dst, padded[:]...), nil
}

// MarshalSSZ marshals a SSZUInt256 into a byte slice.
// It returns the byte slice and an error if any.
func (s *SSZUInt256) MarshalSSZ() ([]byte, error) {
	x := byteslib.ToBytes32(*s)
	return x[:], nil
}

// UnmarshalSSZ unmarshals a SSZUInt256 from SSZ-encoded data. It returns an
// error
// if the buffer size is incorrect.
func (s *SSZUInt256) UnmarshalSSZ(buf []byte) error {
	x := byteslib.ToBytes32(buf)
	if s == nil {
		s = new(SSZUInt256)
	}
	*s = x[:]
	return nil
}

// HashTreeRoot computes the hash tree root of the SSZUInt256.
// It returns a fixed-size byte array and an error if any.
func (s SSZUInt256) HashTreeRoot() ([32]byte, error) {
	return fssz.HashWithDefaultHasher(s)
}

// HashTreeRootWith computes the hash tree root of the SSZUInt256 using the
// provided hasher.
// It modifies the hasher's state and returns an error if any.
func (s SSZUInt256) HashTreeRootWith(hh *fssz.Hasher) error {
	hh.AppendBytes32(s[:])
	return nil
}

// String returns the string representation of the SSZUInt256.
func (s SSZUInt256) String() string {
	return string(s[:])
}
