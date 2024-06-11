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
		require.Equal(t, expected, le.UnwrapU256())
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
