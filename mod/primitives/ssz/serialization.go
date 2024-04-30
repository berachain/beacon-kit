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

// bitsPerByte is the number of bits in a byte.
const bitsPerByte = 8

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

// MostSignificantBitIndex uses bitwise operations for a fast determination
// of the most significant bit's index in a byte.
//
//nolint:mnd // lots of random numbers in bit math.
func MostSignificantBitIndex(x byte) int {
	if x == 0 {
		return -1
	}

	// Initialize the result index
	r := 0
	// Check if the upper half of the byte (higher 4 bits) is non-zero
	if x >= 0x10 {
		// Right shift by 4 bits to focus on the higher half
		x >>= 4
		// Increment result index by 4 because we've shifted half the byte
		r += 4
	}
	// Check if the upper quarter of the byte (bits 4-5) is non-zero
	if x >= 0x4 {
		// Right shift by 2 bits to focus on the next significant bits
		x >>= 2
		// Increment result index by 2 because we've shifted two bits
		r += 2
	}
	// Check if the second bit is set
	if x >= 0x2 {
		// Increment result index by 1 because the second bit is significant
		r++
	}
	return r
}

// TODO: May be buggy, see test case 3 TestUnMarshalBitList
// UnMarshalBitList converts a byte slice into a boolean slice where each bit
// represents a boolean value. The function assumes the input byte slice
// represents a bit list in a compact form, where the presence of a sentinel bit
// (most significant bit of the last byte in the array) can be used to deduce
// the length of the bitlist (not the limit). It returns a slice of booleans
// representing the bit list, excluding the sentinel bit.
func UnmarshalBitList(bv []byte) []bool {
	if len(bv) == 0 {
		return make([]bool, 0)
	}

	msbi := MostSignificantBitIndex(bv[len(bv)-1])
	if msbi == -1 {
		// if no msbi found, its most likely all padding/malformed, return an
		// empty []bool of len 0
		return make([]bool, 0)
	}
	arrL := bitsPerByte*(len(bv)-1) + msbi
	var newArray = make([]bool, arrL)

	// use a bitmask to get the bit value from the byte for all bytes in the
	// slice
	// note: this reverses the order of the bits as highest bit is last
	// we use the pre-calculated array size using msbi to only read whats
	// relevant
	for j := range len(bv) {
		limit := bitsPerByte
		if j == len(bv)-1 {
			limit = msbi
		}
		for i := range limit {
			val := ((bv[j] & (1 << i)) >> i)
			newArray[(bitsPerByte*j)+i] = (val == 1)
		}
	}

	return newArray
}

func UnmarshalBitVector([]byte) []bool {
	// Bit vectors cannot be unmarshalled as there is no sentinel bit to denote
	// its initial length
	panic("not implemented")
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

// MarshalNull takes any type T and returns an empty byte slice.
// This function is useful when you need to represent a null value in byte slice
// form.
func MarshalNull[T any](T) []byte {
	return []byte{}
}

// MarshalBitVector converts a slice of boolean values into a byte slice where
// each bit represents a boolean value.
func MarshalBitVector(bv []bool) []byte {
	// Calculate the necessary byte length to represent the bit vector.
	//nolint:mnd // per spec.
	array := make([]byte, (len(bv)+7)/bitsPerByte)
	for i, val := range bv {
		if val {
			// set the corresponding bit in the byte slice.
			array[i/bitsPerByte] |= 1 << (i % bitsPerByte)
		}
	}
	// Return the byte slice representation of the bit vector.
	return array
}

// MarshalBitList converts a slice of boolean values into a byte slice where
// each bit represents a boolean value, with an additional bit set at the end.
// Note that from the offset coding, the length (in bytes) of the bitlist is
// known. An additional 1 bit is added to the end, at index e where e is the
// length of the bitlist (not the limit), so that the length in bits will also
// be known.
func MarshalBitList(bv []bool) []byte {
	// Allocate enough bytes to represent the bit list, plus one for the end
	// bit.
	array := make([]byte, (len(bv)/bitsPerByte)+1)
	for i, val := range bv {
		if val {
			// Set the bit at the appropriate position if the boolean is true.
			array[i/8] |= 1 << (i % bitsPerByte)
		}
	}
	// Set the additional bit at the end.
	array[len(bv)/8] |= 1 << (len(bv) % bitsPerByte)
	return array
}
