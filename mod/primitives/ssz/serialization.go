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

package ssz

import (
	"encoding/binary"
)

// ----------------------------- Unmarshal -----------------------------

// UnmarshalU256 unmarshals a big endian U256 from the src input.
func UnmarshalU256L[U256LT ~[32]byte](src []byte) U256LT {
	var u256 U256LT
	copy(u256[:], src)
	return u256
}

// MarshalU256 marshals a little endian U256 into a byte slice.
func UnmarshalU128L[U128LT ~[16]byte](src []byte) U128LT {
	var u128 U128LT
	copy(u128[:], src)
	return u128
}

// UnmarshalU64 unmarshals a little endian U64 from the src input.
func UnmarshalU64[U64T ~uint64](src []byte) U64T {
	return U64T(binary.LittleEndian.Uint64(src))
}

// UnmarshalU32 unmarshals a little endian U32 from the src input.
func UnmarshalU32[U32T ~uint32](src []byte) U32T {
	return U32T(binary.LittleEndian.Uint32(src[:4]))
}

// UnmarshalU16 unmarshals a little endian U16 from the src input.
func UnmarshalU16[U16T ~uint16](src []byte) U16T {
	return U16T(binary.LittleEndian.Uint16(src[:2]))
}

// UnmarshalU8 unmarshals a little endian U8 from the src input.
func UnmarshalU8[U8T ~uint8](src []byte) U8T {
	return U8T(src[0])
}

// UnmarshalBool unmarshals a boolean from the src input.
func UnmarshalBool[BoolT ~bool](src []byte) BoolT {
	return src[0] == 1
}

// ----------------------------- Marshal ------------------------------

// MarshalU256 marshals a big endian U256 into a byte slice.
func MarshalU256[U256LT ~[32]byte](u256 U256LT) []byte {
	var dst U256LT
	copy(dst[:], u256[:])
	return dst[:]
}

// MarshalU128 marshals a little endian U128 into a byte slice.
func MarshalU128[U128LT ~[16]byte](u128 U128LT) []byte {
	var dst [16]byte
	copy(dst[:], u128[:])
	return dst[:]
}

// MarshalU64 marshals a little endian U64 into a byte slice.
func MarshalU64[U64T ~uint64](u64 U64T) []byte {
	//nolint:mnd // 8 is the size of a U64.
	dst := make([]byte, 8)
	//#nosec:G701 // we are using the same size as the U64.
	binary.LittleEndian.PutUint64(dst, uint64(u64))
	return dst
}

// MarshalU32 marshals a little endian U32 into a byte slice.
func MarshalU32[U32T ~uint32](u32 U32T) []byte {
	//nolint:mnd // 4 is the size of a U32.
	dst := make([]byte, 4)
	//#nosec:G701 // we are using the same size as the U32.
	binary.LittleEndian.PutUint32(dst, uint32(u32))
	return dst
}

// MarshalU16 marshals a little endian U16 into a byte slice.
func MarshalU16[U16T ~uint16](u16 U16T) []byte {
	//nolint:mnd // 2 is the size of a U16.
	dst := make([]byte, 2)
	//#nosec:G701 // we are using the same size as the U16.
	binary.LittleEndian.PutUint16(dst, uint16(u16))
	return dst
}

// MarshalU8 marshals a little endian U8 into a byte slice.
func MarshalU8[U8T ~uint8](u8 U8T) []byte {
	return []byte{byte(u8)}
}

// MarshalBool marshals a boolean into a byte slice.
func MarshalBool[BoolT ~bool](b BoolT) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}
