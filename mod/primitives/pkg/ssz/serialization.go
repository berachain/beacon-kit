// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package ssz

import (
	"encoding/binary"
	"math/bits"

	"github.com/berachain/beacon-kit/mod/errors"
)

// bitsPerByte is the number of bits in a byte.
const bitsPerByte = 8

// ----------------------------- Unmarshal -----------------------------

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
func UnmarshalBool[BoolT ~bool](src []byte) (BoolT, error) {
	if len(src) != 1 {
		return false, errors.Wrapf(ErrInvalidLength,
			"expected 1 byte, got %d", len(src))
	}

	switch src[0] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, errors.Wrapf(ErrInvalidByteValue,
			"expected 0 or 1, got %d", src[0])
	}
}

// MostSignificantBitIndex uses a lookup table for fast determination of the
// most significant bit's index in a byte.
func MostSignificantBitIndex(x byte) int {
	return bits.Len8(x) - 1
}

// UnmarshalBitList converts a byte slice into a boolean slice where each bit
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

	lastByteStartIdx := bitsPerByte * (len(bv) - 1)
	arrLen := lastByteStartIdx + msbi
	// we use the pre-calculated array size to edit in place
	var newArray = make([]bool, arrLen)

	// use a bitmask to get the bit value from the byte for all bytes in the
	// slice
	// note: this reverses the order of the bits in a byte, as higher bits come
	// later in the array
	for j := range len(bv) - 1 {
		for i := range bitsPerByte {
			val := bv[j] & (1 << i)
			newArray[(bitsPerByte*j)+i] = val > 0
		}
	}

	lastByte := bv[len(bv)-1]
	for i := range msbi {
		val := lastByte & (1 << i)
		newArray[lastByteStartIdx+i] = val > 0
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
