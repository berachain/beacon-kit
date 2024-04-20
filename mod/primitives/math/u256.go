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

package math

import (
	"math/big"

	"github.com/berachain/beacon-kit/mod/primitives/bytes"
	"github.com/holiman/uint256"
)

// U256NumBytes is the number of bytes in a uint256.
const U256NumBytes = 32

// U256 represents a uint256 number stored as big-endian
// format.
type U256 [32]byte

// Int is an alias for U256, it is used to prevent fastssz
// from freaking out -_- smfh.
type Int = U256

// Wei is the smallest unit of Ether, we store the value as LittleEndian for
// the best compatibility with the SSZ spec.
type Wei = U256

// --------------------------- Constructors ----------------------------

// NewU256 creates a new U256 from a byte slice.
func NewU256(bz []byte) (*U256, error) {
	// Ensure that we are not silently truncating the input.
	if len(bz) > U256NumBytes {
		return nil, ErrUnexpectedInputLength(U256NumBytes, len(bz))
	}
	n := U256([32]byte(bz))
	return &n, nil
}

// MustNewU256 creates a new U256 from a byte slice.
// Panics if the input is invalid.
func MustNewU256(bz []byte) *U256 {
	n, err := NewU256(bz)
	if err != nil {
		panic(err)
	}
	return n
}

// NewU256FromBigInt creates a new U256 from a big.Int.
func NewU256FromBigInt(b *big.Int) (*U256, error) {
	if b == nil {
		return nil, ErrNilBigInt
	}
	u := new(uint256.Int)
	if u.SetFromBig(b) {
		return nil, ErrOverflowBigInt
	}
	return (*U256)(u.Bytes()), nil
}

// MustNewU256FromBigInt creates a new U256 from a big.Int.
// Panics if the input is invalid.
func MustNewU256FromBigInt(b *big.Int) *U256 {
	n, err := NewU256FromBigInt(b)
	if err != nil {
		panic(err)
	}
	return n
}

// ---------------------------- Unwrappers -----------------------------

// UnwrapBig converts the U256 type to a *big.Int.
func (u U256) UnwrapBig() *big.Int {
	return new(big.Int).SetBytes(u[:])
}

// -------------------------- JSONMarshallable -------------------------

// MarshalJSON marshals a U256 to JSON, it flips the endianness
// before encoding it to hex such that it is marshalled as big-endian.
func (u *U256) MarshalJSON() ([]byte, error) {
	return new(uint256.Int).SetBytes(u[:]).MarshalJSON()
}

// UnmarshalJSON unmarshals a U256 from JSON by decoding the hex
// string and flipping the endianness, such that it is unmarshalled as
// big-endian.
func (u *U256) UnmarshalJSON(input []byte) error {
	n := new(uint256.Int)
	if err := n.UnmarshalJSON(input); err != nil {
		return err
	}
	*u = U256(n.Bytes())
	return nil
}

// -------------------------- SSZMarshallable --------------------------

// MarshalSSZTo serializes the U64 into a byte slice.
func (u U256) MarshalSSZTo(buf []byte) ([]byte, error) {
	if len(buf) != U256NumBytes {
		return nil, ErrUnexpectedInputLength(U256NumBytes, len(buf))
	}
	// Reverse incoming BigEndian into LittleEndian for serialization.
	return bytes.CopyAndReverseEndianess(buf), nil
}

// MarshalSSZ serializes a U256 into a byte slice.
func (u *U256) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, U256NumBytes)
	return u.MarshalSSZTo(buf)
}

// UnmarshalSSZ deserializes a U256 from a byte slice.
func (u *U256) UnmarshalSSZ(buf []byte) error {
	if len(buf) != U256NumBytes {
		return ErrUnexpectedInputLength(U256NumBytes, len(buf))
	}
	// Reverse incoming LittleEndian into BigEndian for storage.
	x := [32]byte(bytes.CopyAndReverseEndianess(buf))
	*u = U256(x)
	return nil
}

// SizeSSZ returns the size of the U256 in bytes.
func (s U256) SizeSSZ() int {
	return U256NumBytes
}

// String returns the string representation of a U256.
func (s *U256) String() string {
	return new(uint256.Int).SetBytes(s[:]).String()
}
