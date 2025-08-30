// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

	"github.com/berachain/beacon-kit/primitives/encoding/hex"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestU64_MarshalText(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			u := math.U64(tt.input)
			result, err := u.MarshalText()
			require.NoError(t, err)
			require.Equal(t, tt.expected, string(result))
		})
	}
}

func TestU64_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		json     string
		expected uint64
		err      error
	}{
		{"Valid hex string", "\"0x7b\"", 123, nil},
		{"Zero value", "\"0x0\"", 0, nil},
		{"Max uint64 value", "\"0xffffffffffffffff\"", ^uint64(0), nil},
		{"Invalid hex string", "\"0xxyz\"", 0, hex.ErrInvalidString},
		{"Invalid JSON text", "", 0, hex.ErrNonQuotedString},
		{"Invalid quoted JSON text", `"0x`, 0, hex.ErrNonQuotedString},
		{"Empty JSON text", `""`, 0, hex.ErrEmptyString},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()
	u := ^math.U64(0)
	require.Panics(t, func() {
		_ = u.NextPowerOfTwo()
	})
}

func TestU64_ILog2Ceil(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := tt.value.ILog2Ceil()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_ILog2Floor(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := tt.value.ILog2Floor()
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_PrevPowerOfTwo(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
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
	t.Parallel()

	tests := []struct {
		name        string
		value       math.U64
		expectedHex string
	}{
		{
			name:        "zero",
			value:       math.U64(0),
			expectedHex: "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			// https://eth2book.info/capella/part2/building_blocks/merkleization/#the-data-root
			name:        "nine",
			value:       math.U64(9),
			expectedHex: "0x0900000000000000000000000000000000000000000000000000000000000000",
		},
		{
			name:        "max uint64",
			value:       math.U64(1<<64 - 1),
			expectedHex: "0xffffffffffffffff000000000000000000000000000000000000000000000000",
		},
		{
			// https://eth2book.info/capella/part2/building_blocks/merkleization/#the-data-root
			name:        "large number",
			value:       math.U64(3080829),
			expectedHex: "0x7d022f0000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expectedHex, tt.value.HashTreeRoot().String())
		})
	}
}

func TestGweiFromWei(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       func(t *testing.T) *big.Int
		expectedErr error
		expectedRes math.Gwei
	}{
		{
			name: "invalid negative gwei",
			input: func(t *testing.T) *big.Int {
				t.Helper()
				b, _ := new(big.Int).SetString("-1", 10)
				return b
			},
			expectedErr: math.ErrGweiOverflow,
			expectedRes: math.Gwei(0),
		},
		{
			name: "invalid huge gwei",
			input: func(t *testing.T) *big.Int {
				t.Helper()
				b, _ := new(
					big.Int,
				).SetString("18446744073709551616000000000", 10)
				return b
			},
			expectedErr: math.ErrGweiOverflow,
			expectedRes: math.Gwei(0),
		},
		{
			name: "zero wei",
			input: func(t *testing.T) *big.Int {
				t.Helper()
				return big.NewInt(0)
			},
			expectedErr: nil,
			expectedRes: math.Gwei(0),
		},
		{
			name: "one gwei",
			input: func(t *testing.T) *big.Int {
				t.Helper()
				return big.NewInt(params.GWei)
			},
			expectedErr: nil,
			expectedRes: math.Gwei(1),
		},
		{
			name: "arbitrary wei",
			input: func(t *testing.T) *big.Int {
				t.Helper()
				return big.NewInt(params.GWei * 123456789)
			},
			expectedErr: nil,
			expectedRes: math.Gwei(123456789),
		},
		{
			name: "max uint64 wei",
			input: func(t *testing.T) *big.Int {
				t.Helper()
				return new(big.Int).Mul(
					big.NewInt(params.GWei),
					new(big.Int).SetUint64(^uint64(0)),
				)
			},
			expectedErr: nil,
			expectedRes: math.Gwei(1<<64 - 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := math.GweiFromWei(tt.input(t))
			if tt.expectedErr != nil {
				require.ErrorIs(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedRes, result, "Test case: %s", tt.name)
			}
		})
	}
}

func TestGwei_ToWei(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    math.Gwei
		expected func(t *testing.T) *math.U256
	}{
		{
			name:  "zero gwei",
			input: math.Gwei(0),
			expected: func(t *testing.T) *math.U256 {
				t.Helper()
				res, err := math.NewU256FromBigInt(big.NewInt(0))
				require.NoError(t, err)
				return res
			},
		},
		{
			name:  "one gwei",
			input: math.Gwei(1),
			expected: func(t *testing.T) *math.U256 {
				t.Helper()
				res, err := math.NewU256FromBigInt(big.NewInt(params.GWei))
				require.NoError(t, err)
				return res
			},
		},
		{
			name:  "arbitrary gwei",
			input: math.Gwei(123456789),
			expected: func(t *testing.T) *math.U256 {
				t.Helper()
				n := new(big.Int).Mul(
					big.NewInt(params.GWei),
					big.NewInt(123456789),
				)
				res, err := math.NewU256FromBigInt(n)
				require.NoError(t, err)
				return res
			},
		},
		{
			name:  "max uint64 gwei",
			input: math.Gwei(1<<64 - 1),
			expected: func(t *testing.T) *math.U256 {
				t.Helper()
				n := new(big.Int).Mul(
					big.NewInt(params.GWei),
					new(big.Int).SetUint64(1<<64-1),
				)
				res, err := math.NewU256FromBigInt(n)
				require.NoError(t, err)
				return res
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.input.ToWei()
			require.Equal(t, tt.expected(t), result)
		})
	}
}

func TestU64_Base10(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := tt.value.Base10()
			require.Equal(t, tt.expected, result,
				"Test case: %s", tt.name)
		})
	}
}

func TestU64_UnwrapPtr(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			result := tt.value.UnwrapPtr()
			require.NotNil(t, result)
			require.Equal(t, tt.expected, *result,
				"Test case: %s", tt.name)
		})
	}
}

func TestU64_SizeSSZ(t *testing.T) {
	t.Parallel()
	u := math.U64(123)
	require.Equal(t, 8, u.SizeSSZ())
}

func TestU64_MarshalSSZ(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    math.U64
		expected []byte
	}{
		{
			name:     "zero value",
			value:    math.U64(0),
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "small value",
			value:    math.U64(123),
			expected: []byte{0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "large value",
			value:    math.U64(0x123456789ABCDEF0),
			expected: []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
		},
		{
			name:     "max uint64 value",
			value:    math.U64(^uint64(0)),
			expected: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result, err := tt.value.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU64_MarshalSSZTo(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		value    math.U64
		bufSize  int
		expected []byte
	}{
		{
			name:     "zero value with exact buffer",
			value:    math.U64(0),
			bufSize:  8,
			expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "small value with larger buffer",
			value:    math.U64(123),
			bufSize:  16,
			expected: []byte{0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		},
		{
			name:     "value with small buffer (should allocate new)",
			value:    math.U64(0x123456789ABCDEF0),
			bufSize:  4,
			expected: []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			buf := make([]byte, tt.bufSize)
			result, err := tt.value.MarshalSSZTo(buf)
			require.NoError(t, err)
			// The function should append 8 bytes to the buffer
			require.Equal(t, tt.bufSize+8, len(result))
			require.Equal(t, tt.expected, result[tt.bufSize:])
		})
	}
}

func TestU64_UnmarshalSSZ(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []byte
		expected math.U64
		wantErr  bool
	}{
		{
			name:     "zero value",
			input:    []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: math.U64(0),
			wantErr:  false,
		},
		{
			name:     "small value",
			input:    []byte{0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			expected: math.U64(123),
			wantErr:  false,
		},
		{
			name:     "large value",
			input:    []byte{0xF0, 0xDE, 0xBC, 0x9A, 0x78, 0x56, 0x34, 0x12},
			expected: math.U64(0x123456789ABCDEF0),
			wantErr:  false,
		},
		{
			name:     "max uint64 value",
			input:    []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
			expected: math.U64(^uint64(0)),
			wantErr:  false,
		},
		{
			name:     "buffer too short",
			input:    []byte{0x7b, 0x00, 0x00, 0x00},
			expected: math.U64(0),
			wantErr:  true,
		},
		{
			name:     "buffer with extra bytes",
			input:    []byte{0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF},
			expected: math.U64(123),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var u math.U64
			err := u.UnmarshalSSZ(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expected, u)
			}
		})
	}
}

func TestGwei_SSZ(t *testing.T) {
	t.Parallel()
	// Test that Gwei inherits SSZ methods from U64
	g := math.Gwei(1000000000) // 1 Gwei

	// Test SizeSSZ
	require.Equal(t, 8, g.SizeSSZ())

	// Test MarshalSSZ
	marshaled, err := g.MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, []byte{0x00, 0xca, 0x9a, 0x3b, 0x00, 0x00, 0x00, 0x00}, marshaled)

	// Test UnmarshalSSZ
	var g2 math.Gwei
	err = g2.UnmarshalSSZ(marshaled)
	require.NoError(t, err)
	require.Equal(t, g, g2)
}
