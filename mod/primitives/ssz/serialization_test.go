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

package ssz_test

import (
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/math"
	"github.com/berachain/beacon-kit/mod/primitives/ssz"
	"github.com/stretchr/testify/require"
)

func TestMarshalUnmarshalU256(t *testing.T) {
	original := math.U256L{
		0x01,
		0x02,
		0x03,
		0x04,
		0x05,
		0x06,
		0x07,
		0x08,
		0x09,
		0x0A,
		0x0B,
		0x0C,
		0x0D,
		0x0E,
		0x0F,
		0x10,
		0x11,
		0x12,
		0x13,
		0x14,
		0x15,
		0x16,
		0x17,
		0x18,
		0x19,
		0x1A,
		0x1B,
		0x1C,
		0x1D,
		0x1E,
		0x1F,
		0x20,
	}
	marshaled := ssz.MarshalU256(original)
	unmarshaled := ssz.UnmarshalU256L[[32]byte](marshaled)
	require.Equal(t, marshaled, unmarshaled[:], "Marshal/Unmarshal U256 failed")
}

func TestMarshalUnmarshalU128(t *testing.T) {
	original := [16]byte{
		0x01,
		0x02,
		0x03,
		0x04,
		0x05,
		0x06,
		0x07,
		0x08,
		0x09,
		0x0A,
		0x0B,
		0x0C,
		0x0D,
		0x0E,
		0x0F,
		0x10,
	}
	marshaled := ssz.MarshalU128(original)
	unmarshaled := ssz.UnmarshalU128L[[16]byte](marshaled)
	require.Equal(t, marshaled, unmarshaled[:], "Marshal/Unmarshal U128 failed")
}

func TestMarshalUnmarshalU64(t *testing.T) {
	original := uint64(0x0102030405060708)
	marshaled := ssz.MarshalU64(original)
	unmarshaled := ssz.UnmarshalU64[uint64](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U64 failed")
}

func TestMarshalUnmarshalU32(t *testing.T) {
	original := uint32(0x01020304)
	marshaled := ssz.MarshalU32[uint32](original)
	unmarshaled := ssz.UnmarshalU32[uint32](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U32 failed")
}

func TestMarshalUnmarshalU16(t *testing.T) {
	original := uint16(0x0102)
	marshaled := ssz.MarshalU16[uint16](original)
	unmarshaled := ssz.UnmarshalU16[uint16](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U16 failed")
}

func TestMarshalUnmarshalU8(t *testing.T) {
	original := uint8(0x01)
	marshaled := ssz.MarshalU8(original)
	unmarshaled := ssz.UnmarshalU8[uint8](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U8 failed")
}

func TestMarshalUnmarshalBool(t *testing.T) {
	original := true
	marshaled := ssz.MarshalBool(original)
	unmarshaled := ssz.UnmarshalBool[bool](marshaled)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal Bool failed")
}

func TestMarshalBitVector(t *testing.T) {
	var tests = []struct {
		name   string
		bv     []bool
		expect []byte
	}{
		{
			"empty bitvector",
			[]bool{},
			[]byte{},
		},
		{
			"single true value",
			[]bool{true},
			[]byte{1},
		},
		{
			"single false value",
			[]bool{false},
			[]byte{0},
		},
		{
			"multiple values with true at end",
			[]bool{false, false, true, false, false, false, true, true},
			[]byte{0b11000100},
		},
		{
			"multiple values with false at end",
			[]bool{true, true, false, true, true, false, false, false},
			[]byte{0b00011011},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ssz.MarshalBitVector(tt.bv)
			if !reflect.DeepEqual(got, tt.expect) {
				t.Errorf(
					"MarshalBitVector(%v) = %08b; expect %08b",
					tt.bv,
					got,
					tt.expect,
				)
			}
		})
	}
}

func TestMarshalBitList(t *testing.T) {
	// Create a slice of booleans to pass as input
	input := []bool{true, false, true, false, true, false, true}

	output := ssz.MarshalBitList(input)
	// Create a byte slice from a list of binary literals. 0b11010101 is the
	// binary representation of the input slice
	expectedOutput := []byte{0b11010101}
	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("Expected output %08b, got %08b", expectedOutput, output)
	}
}
