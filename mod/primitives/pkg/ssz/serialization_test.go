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

package ssz_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/stretchr/testify/require"
)

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
	unmarshaled, err := ssz.UnmarshalBool[bool](marshaled)
	require.NoError(t, err)
	require.Equal(t, original, unmarshaled, "Marshal/Unmarshal Bool failed")
}

func FuzzMarshalUnmarshalU64(f *testing.F) {
	f.Fuzz(func(t *testing.T, original uint64) {
		marshaled := ssz.MarshalU64(original)
		unmarshaled := ssz.UnmarshalU64[uint64](marshaled)
		require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U64 failed")
	})
}

func FuzzMarshalUnmarshalU32(f *testing.F) {
	f.Fuzz(func(t *testing.T, original uint32) {
		marshaled := ssz.MarshalU32[uint32](original)
		unmarshaled := ssz.UnmarshalU32[uint32](marshaled)
		require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U32 failed")
	})
}

func FuzzMarshalUnmarshalU16(f *testing.F) {
	f.Fuzz(func(t *testing.T, original uint16) {
		marshaled := ssz.MarshalU16[uint16](original)
		unmarshaled := ssz.UnmarshalU16[uint16](marshaled)
		require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U16 failed")
	})
}

func FuzzMarshalUnmarshalU8(f *testing.F) {
	f.Fuzz(func(t *testing.T, original uint8) {
		marshaled := ssz.MarshalU8(original)
		unmarshaled := ssz.UnmarshalU8[uint8](marshaled)
		require.Equal(t, original, unmarshaled, "Marshal/Unmarshal U8 failed")
	})
}

func FuzzMarshalUnmarshalBool(f *testing.F) {
	f.Fuzz(func(t *testing.T, original bool) {
		marshaled := ssz.MarshalBool(original)
		unmarshaled, err := ssz.UnmarshalBool[bool](marshaled)
		require.NoError(t, err)
		require.Equal(t, original, unmarshaled, "Marshal/Unmarshal Bool failed")
	})
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
			require.Equal(t, tt.expect, got, "MarshalBitVector failed")
		})
	}
}

func TestMarshalBitList(t *testing.T) {
	var testcases = []struct {
		name      string
		input     []bool
		expOutput []byte
	}{
		{
			name:      "empty input",
			input:     []bool{},
			expOutput: []byte{0b00000001},
		},
		{
			name:      "single true input",
			input:     []bool{true},
			expOutput: []byte{0b00000011},
		},
		{
			name:      "four elements input",
			input:     []bool{true, true, false, false},
			expOutput: []byte{0b00010011},
		},
		{
			name:      "seven elements input",
			input:     []bool{true, false, true, false, true, false, true},
			expOutput: []byte{0b11010101},
		},
		{
			name: "eight elements input",
			input: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
			},
			expOutput: []byte{0b01010101, 0b00000001},
		},
		{
			name: "nine elements input",
			input: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
				false,
			},
			expOutput: []byte{0b01010101, 0b00000010},
		},
		{
			name: "fifteen elements input",
			input: []bool{true, false, true, false, true, false, true, false,
				true, true, true, true, true, true, true,
			},
			expOutput: []byte{0b01010101, 0b11111111},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			output := ssz.MarshalBitList(tc.input)
			require.Equal(t, tc.expOutput, output, "Failed at "+tc.name)
		})
	}
}

func TestMostSignificantBitIndex(t *testing.T) {
	var tests = []struct {
		name     string
		original byte
		result   int
	}{
		{"0", byte('\x00'), -1},
		{"1", byte('\x01'), 0},
		{"2", byte('\x02'), 1},
		{"4", byte('\x04'), 2},
		{"8", byte('\x08'), 3},
		{"16", byte('\x10'), 4},
		{"32", byte('\x20'), 5},
		{"64", byte('\x40'), 6},
		{"128", byte('\x80'), 7},
		{"255", byte('\xFF'), 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ssz.MostSignificantBitIndex(tt.original)
			require.Equal(t, tt.result, result)
		})
	}
}

func FuzzMostSignificantBitIndex(f *testing.F) {
	f.Fuzz(func(t *testing.T, original byte) {
		result := ssz.MostSignificantBitIndex(original)

		// Basic bounds checking
		require.GreaterOrEqual(t, result, -1)
		require.LessOrEqual(t, result, 7)

		// Check each index edge for violations of spec
		switch {
		case int(original) == 0:
			require.Equal(t, -1, result)
		case int(original) < 2:
			require.Equal(t, 0, result)
		case int(original) < 4:
			require.Equal(t, 1, result)
		case int(original) < 8:
			require.Equal(t, 2, result)
		case int(original) < 16:
			require.Equal(t, 3, result)
		case int(original) < 32:
			require.Equal(t, 4, result)
		case int(original) < 64:
			require.Equal(t, 5, result)
		case int(original) < 128:
			require.Equal(t, 6, result)
		default:
			require.Equal(t, 7, result)
		}
	})
}

func BenchmarkMostSignificantBitIndex(b *testing.B) {
	var table = []struct {
		input byte
	}{
		{input: 0},
		{input: 1},
		{input: 2},
		{input: 4},
		{input: 8},
		{input: 16},
		{input: 32},
		{input: 64},
		{input: 128},
		{input: 255},
	}

	for _, v := range table {
		b.Run(fmt.Sprintf("input_size_%d", v.input), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ssz.MostSignificantBitIndex(v.input)
			}
		})
	}
}

func TestUnmarshalBitList(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		expOutput []bool
	}{
		{
			name:      "Empty input",
			input:     []byte{},
			expOutput: []bool{},
		},
		{
			name:      "Input with sentinel bit set",
			input:     []byte{0b00000011},
			expOutput: []bool{true},
		},
		{
			name:      "Input with multiple bits set",
			input:     []byte{0b11001100},
			expOutput: []bool{false, false, true, true, false, false, true},
		},
		{
			name: "Input with multiple bits set - check both marshal and unmarshal",
			// noliint: lll
			input: ssz.MarshalBitList([]bool{true, false, true, false,
				true, false, true,
			}),
			expOutput: []bool{true, false, true, false, true, false, true},
		},
		{
			name:  "Input with 2 bytes set - check input and output",
			input: []byte{0b01010101, 0b11111111},
			expOutput: []bool{true, false, true, false, true, false,
				true, false, true, true, true, true, true, true, true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output := ssz.UnmarshalBitList(tc.input)
			require.Equal(t, tc.expOutput, output, "unmarshal failed")
		})
	}
}

func FuzzMarshalUnmarshalBitList(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) == 0 {
			return
		}

		totalBits := len(data) * 8
		randomBitLength := totalBits - rand.Intn(8)
		// Convert bytes to a bit list (bool slice) with the random length
		bitList := make([]bool, randomBitLength)
		for i, b := range data {
			for j := 0; j < 8 && i*8+j < randomBitLength; j++ {
				bitList[i*8+j] = (b & (1 << j)) != 0
			}
		}

		marshaled := ssz.MarshalBitList(bitList)
		unmarshaled := ssz.UnmarshalBitList(marshaled)

		// Check if the original and unmarshaled bit lists are the same
		require.Equal(t, bitList, unmarshaled,
			"Original and unmarshaled bit lists do not match")
	})
}

func TestMarshalUnmarshalBitList(t *testing.T) {
	var tests = []struct {
		name      string
		input     []bool
		expOutput []byte
	}{
		{
			name:      "empty input",
			input:     []bool{},
			expOutput: []byte{0b00000001},
		},
		{
			name:      "single true input",
			input:     []bool{true},
			expOutput: []byte{0b00000011},
		},
		{
			name:      "four elements input",
			input:     []bool{true, true, false, false},
			expOutput: []byte{0b00010011},
		},
		{
			name:      "seven elements input",
			input:     []bool{true, false, true, false, true, false, true},
			expOutput: []byte{0b11010101},
		},
		{
			name: "eight elements input",
			input: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
			},
			expOutput: []byte{0b01010101, 0b00000001},
		},
		{
			name: "nine elements input",
			input: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
				false,
			},
			expOutput: []byte{0b01010101, 0b00000010},
		},
		{
			name: "fifteen elements input",
			input: []bool{
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				false,
				true,
				true,
				true,
				true,
				true,
				true,
				true,
			},
			expOutput: []byte{0b01010101, 0b11111111},
		},
		{
			name: "alternating pattern",
			input: []bool{true, false, true, false, true, false, true, false,
				true, false},
			expOutput: []byte{0b01010101, 0b00000101},
		},
		{
			name: "all true",
			input: []bool{true, true, true, true, true, true, true, true,
				true, true},
			expOutput: []byte{0b11111111, 0b00000111},
		},
		{
			name: "all false",
			input: []bool{false, false, false, false, false, false, false,
				false, false, false},
			expOutput: []byte{0b00000000, 0b00000100},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			marshaled := ssz.MarshalBitList(tc.input)
			unmarshaled := ssz.UnmarshalBitList(marshaled)
			require.Equal(
				t,
				tc.input,
				unmarshaled,
				"marshal/unmarshal not equal",
			)
			require.Equal(t, tc.expOutput, marshaled)
		})
	}
}
