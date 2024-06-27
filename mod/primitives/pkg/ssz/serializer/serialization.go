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

package serializer

import (
	"fmt"
	"math/bits"
)

// bitsPerByte is the number of bits in a byte.
const bitsPerByte = 8

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

// MarshalVectorFixed converts a slice of basic values into a byte slice.
func MarshalVectorFixed[T interface{ MarshalSSZ() ([]byte, error) }](
	out []byte, v []T,
) ([]byte, error) {
	// From the Spec:
	// fixed_parts = [
	// 		serialize(element)
	// 			if not is_variable_size(element)
	//			else None for element in value,
	// 		]
	// VectorBasic has all fixed types, so we simply
	// serialize each element and pack them together.
	for _, val := range v {
		bytes, err := val.MarshalSSZ()
		if err != nil {
			return out, err
		}
		out = append(out, bytes...)
	}
	return out, nil
}

// UnmarshalVectorFixed converts a byte slice into a slice of basic values.
func UnmarshalVectorFixed[
	T interface {
		NewFromSSZ([]byte) (T, error)
		SizeSSZ() int
	},
](
	buf []byte,
) ([]T, error) {
	var (
		err error
		t   T
	)
	elementSize := t.SizeSSZ()
	if len(buf)%elementSize != 0 {
		return nil, fmt.Errorf(
			"invalid buffer length %d for element size %d",
			len(buf),
			elementSize,
		)
	}

	result := make([]T, 0, len(buf)/elementSize)
	for i := 0; i < len(buf); i += elementSize {
		if t, err = t.NewFromSSZ(buf[i : i+elementSize]); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
