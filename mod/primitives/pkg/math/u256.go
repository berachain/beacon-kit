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

	byteslib "github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/hex"
	"github.com/holiman/uint256"
)

// U256NumBytes is the number of bytes in a uint256.
const U256NumBytes = 32

// U256 represents a uint256 number stored as big-endian
// format.
type U256 = uint256.Int

// Wei is the smallest unit of Ether, we store the value as LittleEndian for
// the best compatibility with the SSZ spec.
type Wei = U256L

// U256L represents a uint256 number stored as little-endian. It
// is designed to marshal and unmarshal JSON in big-endian
// format, while under the hood storing the value as little-endian
// for compatibility with the SSZ spec.
type U256L [32]byte

// --------------------------- Constructors ----------------------------

// NewU256L creates a new U256L from a byte slice.
func NewU256L(bz []byte) (U256L, error) {
	// Ensure that we are not silently truncating the input.
	if len(bz) > U256NumBytes {
		return U256L{}, ErrUnexpectedInputLength(U256NumBytes, len(bz))
	}
	return U256L(byteslib.ExtendToSize(bz, U256NumBytes)), nil
}

// MustNewU256L creates a new U256L from a byte slice.
// Panics if the input is invalid.
func MustNewU256L(bz []byte) U256L {
	n, err := NewU256L(bz)
	if err != nil {
		panic(err)
	}
	return n
}

// NewU256LFromBigEndian creates a new U256L from a big-endian
// byte slice.
func NewU256LFromBigEndian(b []byte) (U256L, error) {
	return NewU256L(byteslib.CopyAndReverseEndianess(b))
}

// MustNewU256LFromBigEndian creates a new U256L from a big-endian
// byte slice. Panics if the input is invalid.
func MustNewU256LFromBigEndian(b []byte) U256L {
	n, err := NewU256L(byteslib.CopyAndReverseEndianess(b))
	if err != nil {
		panic(err)
	}
	return n
}

// NewU256LFromBigInt creates a new U256L from a big.Int.
func NewU256LFromBigInt(b *big.Int) (U256L, error) {
	if b == nil {
		return U256L{}, ErrNilBigInt
	} else if b.Sign() < 0 {
		return U256L{}, ErrNegativeBigInt(b)
	}
	return NewU256LFromBigEndian(b.Bytes())
}

// MustNewU256LFromBigInt creates a new U256L from a big.Int.
// Panics if the input is invalid.
func MustNewU256LFromBigInt(b *big.Int) U256L {
	n, err := NewU256LFromBigInt(b)
	if err != nil {
		panic(err)
	}
	return n
}

// ------------------------------ Unwraps ------------------------------

// UnwrapU256 converts an U256L to a raw [32]byte chunk.
func (s U256L) Unwrap() [32]byte {
	return s
}

// UnwrapU256 converts an U256L to a *U256.
func (s U256L) UnwrapU256() *U256 {
	return new(uint256.Int).SetBytes(byteslib.CopyAndReverseEndianess(s[:]))
}

// UnwrapBig converts a U256 to a non-negative big.Int.
func (s U256L) UnwrapBig() *big.Int {
	// SetBytes treats byte slice as unsigned int in big-endian, so result
	// will always be non-negative.
	return new(big.Int).SetBytes(byteslib.CopyAndReverseEndianess(s[:]))
}

// -------------------------- JSONMarshallable -------------------------

// MarshalJSON marshals a U256L to JSON, it flips the endianness
// before encoding it to hex such that it is marshalled as big-endian.
func (s U256L) MarshalJSON() ([]byte, error) {
	return []byte(hex.FromBigInt(s.UnwrapBig()).AddQuotes().Unwrap()), nil
}

// UnmarshalJSON unmarshals a U256L from JSON by decoding the hex
// string and flipping the endianness, such that it is unmarshalled as
// big-endian.
func (s *U256L) UnmarshalJSON(input []byte) error {
	baseFee, err := hex.FromJSONString(input).ToBigInt()
	if err != nil {
		return err
	}
	*s = U256L(
		byteslib.ExtendToSize(
			byteslib.CopyAndReverseEndianess(
				baseFee.Bytes()), U256NumBytes),
	)
	return nil
}

// -------------------------- SSZMarshallable --------------------------

// MarshalSSZTo serializes the U64 into a byte slice.
func (s U256L) MarshalSSZTo(buf []byte) ([]byte, error) {
	copy(buf, s[:])
	return buf, nil
}

// MarshalSSZ serializes a U256L into a byte slice.
func (s U256L) MarshalSSZ() ([]byte, error) {
	return s[:], nil
}

// UnmarshalSSZ deserializes a U256L from a byte slice.
func (s *U256L) UnmarshalSSZ(buf []byte) error {
	if len(buf) != U256NumBytes {
		return ErrUnexpectedInputLength(U256NumBytes, len(buf))
	}
	copy(s[:], buf)
	return nil
}

// SizeSSZ returns the size of the U256L in bytes.
func (s U256L) SizeSSZ() int {
	return U256NumBytes
}

// String returns the string representation of a U256L.
func (s *U256L) String() string {
	return s.UnwrapU256().String()
}
