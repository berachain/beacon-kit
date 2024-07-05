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

//nolint:lll // long strings.
package math_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	huint256 "github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestLittleEndian_UInt256(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{1, 2, 3, 4, 5}, []byte{5, 4, 3, 2, 1}},
		{[]byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}},
		{[]byte{255, 255, 255, 255}, []byte{255, 255, 255, 255}},
	}

	for _, tc := range testCases {
		le, err := math.NewU256L(tc.input)
		require.NoError(t, err)
		expected := new(huint256.Int).SetBytes(tc.expected)
		require.Equal(t, expected, le.UnwrapU256().Unwrap())
	}
}

func TestLittleEndian_Big(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected []byte
	}{
		{[]byte{1, 2, 3, 4, 5}, []byte{5, 4, 3, 2, 1}},
		{[]byte{0, 0, 0, 0}, []byte{0, 0, 0, 0}},
		{[]byte{255, 255, 255, 255}, []byte{255, 255, 255, 255}},
	}

	for _, tc := range testCases {
		le, err := math.NewU256L(tc.input)
		require.NoError(t, err)
		expected := new(huint256.Int).SetBytes(tc.expected)
		require.Equal(t, expected.ToBig(), le.UnwrapBig())
	}
}

func TestLittleEndian_MarshalJSON(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{[]byte{1, 2, 3, 4, 5}, "\"0x504030201\""},
		{[]byte{0, 0, 0, 0}, "\"0x0\""},
		{[]byte{255, 255, 255, 255}, "\"0xffffffff\""},
	}

	for _, tc := range testCases {
		le, err := math.NewU256L(tc.input)
		require.NoError(t, err)
		result, err := le.MarshalJSON()
		require.NoError(t, err)
		require.JSONEq(t, tc.expected, string(result))
	}
}

func TestLittleEndian_UnmarshalJSON(t *testing.T) {
	testCases := []struct {
		json     string
		expected []byte
	}{
		{"\"0x504030201\"", []byte{1, 2, 3, 4, 5}},
		{"\"0x0\"", []byte{0, 0, 0, 0}},
		{"\"0xffffffff\"", []byte{255, 255, 255, 255}},
	}

	for _, tc := range testCases {
		le := new(math.U256L)
		err := le.UnmarshalJSON([]byte(tc.json))
		require.NoError(t, err)
		expected, err := math.NewU256L(tc.expected)
		require.NoError(t, err)
		require.Equal(t, expected, *le)
	}
}

func TestU256L_MarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
	}{
		{
			name: "zero",
			input: []byte{
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			expected: make([]byte, 32),
		},
		{
			name:     "max value",
			input:    bytes.Repeat([]byte{255}, 32),
			expected: bytes.Repeat([]byte{255}, 32),
		},
		{
			name: "arbitrary value",
			input: []byte{
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			expected: []byte{
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u256l, err := math.NewU256L(tt.input)
			require.NoError(t, err)
			result, err := u256l.MarshalSSZ()
			require.NoError(t, err)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestU256L_UnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected []byte
		err      error
	}{
		{
			name: "valid data",
			data: []byte{
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			expected: []byte{
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
		},
		{
			name: "invalid data - short buffer",
			data: []byte{0, 0},
			err:  math.ErrUnexpectedInputLengthBase,
		},
		{
			name:     "valid data - zero",
			data:     make([]byte, 32),
			expected: make([]byte, 32),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var u math.U256L
			err := u.UnmarshalSSZ(tt.data)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				var expected math.U256L
				expected, err = math.NewU256L(tt.expected)
				require.NoError(t, err)
				require.Equal(t, expected, u)
			}
		})
	}
}

func TestNewU256L_SilentTruncation(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected [32]byte
	}{
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
			13, 14, 15, 16, 16, 17, 18, 19, 20, 21, 22,
			23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34},
			[32]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13,
				14, 15, 16, 16, 17, 18, 19, 20, 21, 22, 23, 24,
				25, 26, 27, 28, 29, 30, 31}},
	}

	for _, tc := range testCases {
		_, err := math.NewU256L(tc.input)
		require.ErrorIs(t, err, math.ErrUnexpectedInputLengthBase)
	}
}

func TestSignedness_Big(t *testing.T) {
	for z := -100; z < 0; z++ {
		a := big.NewInt(int64(z))
		_, err := math.NewU256LFromBigInt(a)
		require.Error(t, err)
	}
}

// UnwrapBig() should never return a negative number.
func TestUnwrapBigSign(t *testing.T) {
	tests := []struct {
		name  string
		input math.U256L
	}{
		{
			name:  "Zero value",
			input: math.U256L{},
		},
		{
			name: "Maximum value",
			input: func() math.U256L {
				var maxVal math.U256L
				for i := range maxVal {
					maxVal[i] = 0xff
				}
				return maxVal
			}(),
		},
		{
			name: "Random value",
			input: math.U256L{
				0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad,
				0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef,
				0xca, 0xfe, 0xba, 0xbe, 0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe,
				0xba, 0xbe},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.UnwrapBig()
			if result.Sign() < 0 {
				t.Errorf(
					"UnwrapBig() resulted in a negative number: %v",
					result,
				)
			}
		})
	}
}

func TestMustNewU256L(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		panic    bool
	}{
		{
			name:     "valid input",
			input:    []byte{1, 2, 3, 4, 5},
			expected: []byte{5, 4, 3, 2, 1},
			panic:    false,
		},
		{
			name:  "invalid input - too long",
			input: make([]byte, 33),
			panic: true,
		},
		{
			name:     "valid input - single byte",
			input:    []byte{1},
			expected: []byte{1},
			panic:    false,
		},
		{
			name:     "valid input - all zeros",
			input:    []byte{0, 0, 0, 0},
			expected: []byte{0, 0, 0, 0},
			panic:    false,
		},
		{
			name:     "valid input - max uint256",
			input:    bytes.Repeat([]byte{255}, 32),
			expected: bytes.Repeat([]byte{255}, 32),
			panic:    false,
		},
		{
			name:     "valid input - arbitrary value",
			input:    []byte{0xde, 0xad, 0xbe, 0xef},
			expected: []byte{0xef, 0xbe, 0xad, 0xde},
			panic:    false,
		},
		{
			name:     "valid input - 32 bytes",
			input:    bytes.Repeat([]byte{1}, 32),
			expected: bytes.Repeat([]byte{1}, 32),
			panic:    false,
		},
		{
			name:     "valid input - 31 bytes",
			input:    bytes.Repeat([]byte{1}, 31),
			expected: bytes.Repeat([]byte{1}, 31),
			panic:    false,
		},
		{
			name:     "valid input - 30 bytes",
			input:    bytes.Repeat([]byte{1}, 30),
			expected: bytes.Repeat([]byte{1}, 30),
			panic:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				require.Panics(t, func() {
					math.MustNewU256L(tt.input)
				})
			} else {
				result := math.MustNewU256L(tt.input)
				expected := new(huint256.Int).SetBytes(tt.expected)
				require.Equal(t, expected, result.UnwrapU256().Unwrap(), "Test case %s", tt.name)
			}
		})
	}
}

func TestNewU256LFromBigEndian(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		err      error
	}{
		{
			name:     "valid input",
			input:    []byte{1, 2, 3, 4, 5},
			expected: []byte{1, 2, 3, 4, 5},
			err:      nil,
		},
		{
			name:  "invalid input - too long",
			input: make([]byte, 33),
			err:   math.ErrUnexpectedInputLengthBase,
		},
		{
			name:     "valid input - all zeros",
			input:    []byte{0, 0, 0, 0},
			expected: []byte{0, 0, 0, 0},
			err:      nil,
		},
		{
			name:     "valid input - single byte",
			input:    []byte{1},
			expected: []byte{1},
			err:      nil,
		},
		{
			name:     "valid input - max uint256",
			input:    bytes.Repeat([]byte{255}, 32),
			expected: bytes.Repeat([]byte{255}, 32),
			err:      nil,
		},
		{
			name:     "valid input - arbitrary value",
			input:    []byte{0xde, 0xad, 0xbe, 0xef},
			expected: []byte{0xde, 0xad, 0xbe, 0xef},
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := math.NewU256LFromBigEndian(tt.input)
			expected := new(huint256.Int).SetBytes(tt.expected)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err, "Test case: %s", tt.name)
			} else {
				require.NoError(t, err, "Test case: %s", tt.name)
				require.Equal(t, expected, result.UnwrapU256().Unwrap(), "Test case: %s", tt.name)
			}
		})
	}
}

func TestMustNewU256LFromBigEndian(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []byte
		panic    bool
	}{
		{
			name:     "valid input",
			input:    []byte{1, 2, 3, 4, 5},
			expected: []byte{1, 2, 3, 4, 5},
			panic:    false,
		},
		{
			name:  "invalid input - too long",
			input: make([]byte, 33),
			panic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				require.Panics(t, func() {
					math.MustNewU256LFromBigEndian(tt.input)
				})
			} else {
				result := math.MustNewU256LFromBigEndian(tt.input)
				expected := new(huint256.Int).SetBytes(tt.expected)
				require.Equal(t, expected, result.UnwrapU256().Unwrap(), "Test case: %s", tt.name)
			}
		})
	}
}

func TestMustNewU256LFromBigInt(t *testing.T) {
	tests := []struct {
		name     string
		input    *big.Int
		expected []byte
		panic    bool
	}{
		{
			name:     "valid input",
			input:    big.NewInt(12345),
			expected: []byte{48, 57},
			panic:    false,
		},
		{
			name:     "valid input - zero",
			input:    big.NewInt(0),
			expected: []byte{0},
			panic:    false,
		},
		{
			name: "valid input - max uint256",
			input: new(
				big.Int,
			).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1)),
			expected: bytes.Repeat([]byte{255}, 32),
			panic:    false,
		},
		{
			name:  "invalid input - negative value",
			input: big.NewInt(-1),
			panic: true,
		},
		{
			name:  "invalid input - too large",
			input: new(big.Int).Lsh(big.NewInt(1), 256),
			panic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panic {
				require.Panics(t, func() {
					math.MustNewU256LFromBigInt(tt.input)
				})
			} else {
				result := math.MustNewU256LFromBigInt(tt.input)
				expected := new(huint256.Int).SetBytes(tt.expected)
				require.Equal(t, expected, result.UnwrapU256().Unwrap(), "Test case: %s", tt.name)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	tests := []struct {
		name     string
		input    math.U256L
		expected [32]byte
	}{
		{
			name:     "unwrap zero value",
			input:    math.U256L{},
			expected: [32]byte{},
		},
		{
			name: "unwrap non-zero value",
			input: math.U256L{
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
				12,
				13,
				14,
				15,
				16,
				17,
				18,
				19,
				20,
				21,
				22,
				23,
				24,
				25,
				26,
				27,
				28,
				29,
				30,
				31,
				32,
			},
			expected: [32]byte{
				1,
				2,
				3,
				4,
				5,
				6,
				7,
				8,
				9,
				10,
				11,
				12,
				13,
				14,
				15,
				16,
				17,
				18,
				19,
				20,
				21,
				22,
				23,
				24,
				25,
				26,
				27,
				28,
				29,
				30,
				31,
				32,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.Unwrap()
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestMarshalSSZTo(t *testing.T) {
	tests := []struct {
		name     string
		input    math.U256L
		expected []byte
	}{
		{
			name:     "marshal zero value",
			input:    math.U256L{},
			expected: make([]byte, 32),
		},
		{
			name: "marshal non-zero value",
			input: math.U256L{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
				17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
			},
			expected: []byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
				17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 32)
			result, err := tt.input.MarshalSSZTo(buf)
			require.NoError(t, err, "Test case: %s", tt.name)
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestSizeSSZ(t *testing.T) {
	tests := []struct {
		name     string
		input    math.U256L
		expected int
	}{
		{
			name:     "size of U256L",
			input:    math.U256L{},
			expected: 32, // U256NumBytes is 32
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.SizeSSZ()
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		input    math.U256L
		expected string
	}{
		{
			name:     "string representation of zero value",
			input:    math.U256L{},
			expected: "0",
		},
		{
			name: "string representation of non-zero value",
			input: math.U256L{
				1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			},
			expected: "1",
		},
		{
			name: "string representation of large value",
			input: math.U256L{
				255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
				255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
			},
			expected: "115792089237316195423570985008687907853269984665640564039457584007913129639935",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			require.Equal(t, tt.expected, result, "Test case: %s", tt.name)
		})
	}
}
