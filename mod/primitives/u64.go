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
	"encoding/binary"
	"math/bits"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// U64NumBytes is the number of bytes in a U64.
const U64NumBytes = 8

// U64 represents a 64-bit unsigned integer that is both SSZ and JSON
// marshallable. We marshal U64 as hex strings in JSON in order to keep the
// execution client apis happy, and we marshal U64 as little-endian in SSZ to be
// compatible with the spec.
type U64 uint64

// -------------------------- SSZMarshallable --------------------------

// MarshalSSZTo serializes the U64 into a byte slice.
func (u U64) MarshalSSZTo(buf []byte) ([]byte, error) {
	binary.LittleEndian.PutUint64(buf, uint64(u))
	return buf, nil
}

// MarshalSSZ serializes the U64 into a byte slice.
func (u U64) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, U64NumBytes)
	if _, err := u.MarshalSSZTo(buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// UnmarshalSSZ deserializes the U64 from a byte slice.
func (u *U64) UnmarshalSSZ(buf []byte) error {
	if len(buf) != U64NumBytes {
		return ErrInvalidSSZLength
	}
	if u == nil {
		u = new(U64)
	}
	*u = U64(binary.LittleEndian.Uint64(buf))
	return nil
}

// SizeSSZ returns the size of the U64 in bytes.
func (u U64) SizeSSZ() int {
	return U64NumBytes
}

// -------------------------- JSONMarshallable -------------------------

// UnmarshalJSON parses a blob in hex syntax.
func (u *U64) UnmarshalJSON(input []byte) error {
	return (*hexutil.Uint64)(u).UnmarshalJSON(input)
}

// MarshalText returns the hex representation of b.
func (u U64) MarshalText() ([]byte, error) {
	return hexutil.Uint64(u).MarshalText()
}

// ---------------------------- U64 Methods ----------------------------

// Unwrap returns the underlying uint64 value of U64.
func (u U64) Unwrap() uint64 {
	return uint64(u)
}

// NextPowerOfTwo returns the next power of two greater than or equal to the.
//
//nolint:gomnd // powers of 2.
func (u U64) NextPowerOfTwo() U64 {
	u--
	u |= u >> 1
	u |= u >> 2
	u |= u >> 4
	u |= u >> 8
	u |= u >> 16
	u++
	return u
}

// ILog2Ceil returns the ceiling of the base 2 logarithm of the U64.
func (u U64) ILog2Ceil() uint8 {
	// Log2(0) is undefined, should we panic?
	if u == 0 {
		return 0
	}
	//#nosec:G701 // we handle the case of u == 0 above, so this is safe.
	return 64 - uint8(bits.LeadingZeros64(uint64(u-1)))
}
