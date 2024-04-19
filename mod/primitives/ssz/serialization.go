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
func UnmarshalU256L[U256 ~[32]byte, B ~[]byte](src B) U256 {
	var u256 U256
	copy(u256[:], src)
	return u256
}

// MarshalU256 marshals a little endian U256 into a byte slice.
func UnmarshalU128L[U128 ~[16]byte, B ~[]byte](src B) U128 {
	var u128 U128
	copy(u128[:], src)
	return u128
}

// UnmarshalU64 unmarshals a little endian U64 from the src input.
func UnmarshalU64[U64 ~uint64, B ~[]byte](src B) U64 {
	return U64(binary.LittleEndian.Uint64(src))
}

// UnmarshalU32 unmarshals a little endian U32 from the src input.
func UnmarshalU32[U32 ~uint32, B ~[]byte](src B) U32 {
	return U32(binary.LittleEndian.Uint32(src[:4]))
}

// UnmarshalU16 unmarshals a little endian U16 from the src input.
func UnmarshalU16[U16 ~uint16, B ~[]byte](src B) U16 {
	return U16(binary.LittleEndian.Uint16(src[:2]))
}

// UnmarshalU8 unmarshals a little endian U8 from the src input.
func UnmarshalU8[U8 ~uint8, B ~[]byte](src B) U8 {
	return U8(src[0])
}

// ----------------------------- Marshal ------------------------------

// MarshalU256 marshals a big endian U256 into a byte slice.
func MarshalU256[U256 ~[32]byte, B ~[]byte](u256 U256) B {
	var dst [32]byte
	copy(dst[:], u256[:])
	return dst[:]
}

// MarshalU128 marshals a little endian U128 into a byte slice.
func MarshalU128[U128 ~[16]byte, B ~[]byte](u128 U128) B {
	var dst [16]byte
	copy(dst[:], u128[:])
	return dst[:]
}

// MarshalU64 marshals a little endian U64 into a byte slice.
func MarshalU64[U64 ~uint64, B ~[]byte](u64 U64) B {
	dst := make([]byte, 8)
	binary.LittleEndian.PutUint64(dst, uint64(u64))
	return dst
}

// MarshalU32 marshals a little endian U32 into a byte slice.
func MarshalU32[U32 ~uint32, B ~[]byte](u32 U32) B {
	dst := make([]byte, 4)
	binary.LittleEndian.PutUint32(dst, uint32(u32))
	return dst
}

// MarshalU16 marshals a little endian U16 into a byte slice.
func MarshalU16[U16 ~uint16, B ~[]byte](u16 U16) B {
	dst := make([]byte, 2)
	binary.LittleEndian.PutUint16(dst, uint16(u16))
	return dst
}

// MarshalU8 marshals a little endian U8 into a byte slice.
func MarshalU8[U8 ~uint8, B ~[]byte](u8 U8) B {
	return []byte{byte(u8)}
}
