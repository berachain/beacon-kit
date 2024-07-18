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

package math_test

import (
	"math/big"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/hex"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/encoding/ssz/schema"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestU64_MarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected []byte
	}{
		{
			name:     "zero",
			value:    0,
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "max uint64",
			value:    math.U64(^uint64(0)),
			expected: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "arbitrary number",
			value:    math.U64(123456789),
			expected: []byte{21, 205, 91, 7, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.value.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_UnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected math.U64
		err      error
	}{
		{
			name:     "valid data",
			data:     []byte{21, 205, 91, 7, 0, 0, 0, 0},
			expected: math.U64(123456789),
		},
		{
			name: "invalid data - short buffer",
			data: []byte{0, 0, 0},
			err:  math.ErrUnexpectedInputLengthBase,
		},
		{
			name:     "valid data - max uint64",
			data:     []byte{255, 255, 255, 255, 255, 255, 255, 255},
			expected: math.U64(^uint64(0)),
		},
		{
			name:     "valid data - zero",
			data:     []byte{0, 0, 0, 0, 0, 0, 0, 0},
			expected: math.U64(0),
		},
		{
			name: "invalid data - long buffer",
			data: []byte{0, 0, 0, 0, 0, 0, 0, 0, 1},
			err:  math.ErrUnexpectedInputLengthBase,
		},
		{
			name:     "valid data - one",
			data:     []byte{1, 0, 0, 0, 0, 0, 0, 0},
			expected: math.U64(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u math.U64
			err := u.UnmarshalSSZ(tt.data)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, u)
			}
		})
	}
}

func TestU64_RoundTripSSZ(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected []byte
	}{
		{
			name:     "zero value",
			value:    math.U64(0),
			expected: []byte{0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			name:     "max uint64",
			value:    math.U64(^uint64(0)),
			expected: []byte{255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			name:     "arbitrary number",
			value:    math.U64(123456789),
			expected: []byte{21, 205, 91, 7, 0, 0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test MarshalSSZ
			marshaled, err := tt.value.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.expected, marshaled)

			// Test UnmarshalSSZ
			var unmarshaled math.U64
			err = unmarshaled.UnmarshalSSZ(tt.expected)
			require.NoError(t, err)
			require.Equal(t, tt.value, unmarshaled)
		})
	}
}

func TestU64_MarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{"Zero value", 0, "0x0"},
		{"Small value", 123, "0x7b"},
		{"Max uint64 value", ^uint64(0), "0xffffffffffffffff"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := math.U64(tt.input)
			result, err := u.MarshalText()
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

func TestU64_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected uint64
		err      error
	}{
		{"Valid hex string", "\"0x7b\"", 123, nil},
		{"Zero value", "\"0x0\"", 0, nil},
		{"Max uint64 value", "\"0xffffffffffffffff\"", ^uint64(0), nil},
		{"Invalid hex string", "\"0xxyz\"", 0,
			hex.ErrInvalidString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u math.U64
			err := u.UnmarshalJSON([]byte(tt.json))
			if tt.err != nil {
				require.Error(t, err)
				require.Equal(t, tt.err, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, math.U64(tt.expected), u)
			}
		})
	}
}

func TestU64_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
		err      error
	}{
		{"Valid hex string", "0x7b", 123, nil},
		{"Zero value", "0x0", 0, nil},
		{"Max uint64 value", "0xffffffffffffffff", ^uint64(0), nil},
		{"Invalid hex string", "0xxyz", 0, hex.ErrInvalidString},
		{"Overflow hex string", "0x10000000000000000", 0, hex.ErrUint64Range},
		{"Empty string", "", 0, hex.ErrEmptyString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u math.U64
			err := u.UnmarshalText([]byte(tt.input))
			if tt.err != nil {
				require.Error(t, err)
				require.EqualError(t, err, tt.err.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, math.U64(tt.expected), u)
			}
		})
	}
}

func TestU64_NextPowerOfTwo(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected math.U64
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: math.U64(1),
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: math.U64(1),
		},
		{
			name:     "already a power of two",
			value:    math.U64(8),
			expected: math.U64(8),
		},
		{
			name:     "not a power of two",
			value:    math.U64(9),
			expected: math.U64(16),
		},
		{
			name:     "large number",
			value:    math.U64(1<<20 + 1<<46),
			expected: math.U64(1 << 47),
		},
		{
			name:     "large number with lots zeros",
			value:    math.U64(1<<62 + 1),
			expected: math.U64(1 << 63),
		},
		{
			name:     "large number at the limit",
			value:    math.U64(1 << 63),
			expected: math.U64(1 << 63),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.NextPowerOfTwo()
			require.Equal(t, tt.expected, result)
			require.Equal(
				t,
				tt.expected,
				math.U64(uint64(1)<<tt.value.ILog2Ceil()),
			)
		})
	}
}

func TestU64_NextPowerOfTwoPanic(t *testing.T) {
	u := ^math.U64(0)
	require.Panics(t, func() {
		_ = u.NextPowerOfTwo()
	})
}

func TestU64_ILog2Ceil(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected uint8
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: 0,
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: 0,
		},
		{
			name:     "power of two",
			value:    math.U64(8),
			expected: 3,
		},
		{
			name:     "not a power of two",
			value:    math.U64(9),
			expected: 4,
		},
		{
			name:     "max uint64",
			value:    math.U64(1<<64 - 1),
			expected: 64,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.ILog2Ceil()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_ILog2Floor(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected uint8
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: 0,
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: 0,
		},
		{
			name:     "power of two",
			value:    math.U64(8),
			expected: 3,
		},
		{
			name:     "not a power of two",
			value:    math.U64(9),
			expected: 3,
		},
		{
			name:     "max uint64",
			value:    math.U64(1<<64 - 1),
			expected: 63,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.ILog2Floor()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_PrevPowerOfTwo(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected math.U64
	}{
		{
			name:     "zero",
			value:    math.U64(0),
			expected: 1,
		},
		{
			name:     "one",
			value:    math.U64(1),
			expected: 1,
		},
		{
			name:     "two",
			value:    math.U64(2),
			expected: 2,
		},
		{
			name:     "three",
			value:    math.U64(3),
			expected: 2,
		},
		{
			name:     "four",
			value:    math.U64(4),
			expected: 4,
		},
		{
			name:     "five",
			value:    math.U64(5),
			expected: 4,
		},
		{
			name:     "eight",
			value:    math.U64(8),
			expected: 8,
		},
		{
			name:     "nine",
			value:    math.U64(9),
			expected: 8,
		},
		{
			name:     "thirty-two",
			value:    math.U64(32),
			expected: 32,
		},
		{
			name:     "thirty-three",
			value:    math.U64(33),
			expected: 32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.PrevPowerOfTwo()
			require.Equal(t, tt.expected, result)
			require.Equal(
				t,
				tt.expected,
				math.U64(uint64(1)<<tt.value.ILog2Floor()),
			)
		})
	}
}

func TestU64_HashTreeRoot(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected [32]byte
	}{
		{
			name:  "zero value",
			value: math.U64(0),
			expected: [32]byte{
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
			},
		},
		{
			name:  "max uint64 value",
			value: math.U64(^uint64(0)),
			expected: [32]byte{
				255,
				255,
				255,
				255,
				255,
				255,
				255,
				255,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
			},
		},
		{
			name:  "arbitrary number",
			value: math.U64(123456789),
			expected: [32]byte{
				21,
				205,
				91,
				7,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
			},
		},
		{
			name:  "one",
			value: math.U64(1),
			expected: [32]byte{
				1,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
				0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.value.HashTreeRoot()
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestGweiFromWei(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected math.Gwei
	}{
		{
			name:     "zero wei",
			input:    big.NewInt(0),
			expected: math.Gwei(0),
		},
		{
			name:     "one gwei",
			input:    big.NewInt(math.GweiPerWei),
			expected: math.Gwei(1),
		},
		{
			name:     "arbitrary wei",
			input:    big.NewInt(math.GweiPerWei * 123456789),
			expected: math.Gwei(123456789),
		},
		{
			name: "max uint64 wei",
			input: new(
				big.Int,
			).Mul(big.NewInt(math.GweiPerWei), new(big.Int).SetUint64(^uint64(0))),
			expected: math.Gwei(1<<64 - 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := math.GweiFromWei(tt.input)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestGwei_ToWei(t *testing.T) {
	tests := []struct {
		name     string
		input    math.Gwei
		expected *big.Int
	}{
		{
			name:     "zero gwei",
			input:    math.Gwei(0),
			expected: big.NewInt(0),
		},
		{
			name:     "one gwei",
			input:    math.Gwei(1),
			expected: big.NewInt(math.GweiPerWei),
		},
		{
			name:  "arbitrary gwei",
			input: math.Gwei(123456789),
			expected: new(
				big.Int,
			).Mul(big.NewInt(math.GweiPerWei), big.NewInt(123456789)),
		},
		{
			name:  "max uint64 gwei",
			input: math.Gwei(1<<64 - 1),
			expected: new(
				big.Int,
			).Mul(big.NewInt(math.GweiPerWei), new(big.Int).SetUint64(1<<64-1)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.ToWei()
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestU64_Base10(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected string
	}{
		{
			name:     "zero value",
			value:    math.U64(0),
			expected: "0",
		},
		{
			name:     "small value",
			value:    math.U64(123),
			expected: "123",
		},
		{
			name:     "large value",
			value:    math.U64(123456789),
			expected: "123456789",
		},
		{
			name:     "max uint64 value",
			value:    math.U64(^uint64(0)),
			expected: "18446744073709551615",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.Base10()
			require.Equal(t, tt.expected, result,
				"Test case: %s", tt.name)
		})
	}
}

func TestU64_UnwrapPtr(t *testing.T) {
	tests := []struct {
		name     string
		value    math.U64
		expected uint64
	}{
		{
			name:     "zero value",
			value:    math.U64(0),
			expected: 0,
		},
		{
			name:     "small value",
			value:    math.U64(123),
			expected: 123,
		},
		{
			name:     "large value",
			value:    math.U64(123456789),
			expected: 123456789,
		},
		{
			name:     "max uint64 value",
			value:    math.U64(^uint64(0)),
			expected: ^uint64(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.value.UnwrapPtr()
			require.NotNil(t, result)
			require.Equal(t, tt.expected, *result,
				"Test case: %s", tt.name)
		})
	}
}

func TestU64(t *testing.T) {
	var u math.U64
	require.Equal(t, constants.U64Size, u.SizeSSZ())
	require.True(t, u.IsFixed())
	require.Equal(t, schema.U64(), u.Type())
	require.Equal(t, uint64(1), u.ChunkCount())
}

func TestU64_NewFromSSZ(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected math.U64
		err      bool
	}{
		{name: "Valid1", input: []byte{1, 0, 0, 0, 0, 0, 0, 0},
			expected: 1, err: false},
		{name: "Valid0", input: []byte{0, 0, 0, 0, 0, 0, 0, 0},
			expected: 0, err: false},
		{name: "InvalidLength", input: []byte{1, 0, 0, 0, 0, 0, 0},
			expected: 0,
			err:      true},
	}

	var u64 math.U64 = 1

	for _, tt := range tests {
		result, err := u64.NewFromSSZ(tt.input)
		if tt.err {
			require.Error(t, err, "Test name %s", tt.name)
		} else {
			require.NoError(t, err, "Test name %s", tt.name)
			require.Equal(t, tt.expected, result,
				"Test name %s", tt.name)
		}
	}
}
