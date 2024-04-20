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

package math_test

import (
	"bytes"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/math"
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
				expected, err := math.NewU256L(tt.expected)
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
				14, 15, 16, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}},
	}

	for _, tc := range testCases {
		_, err := math.NewU256L(tc.input)
		require.ErrorIs(t, err, math.ErrUnexpectedInputLengthBase)
	}
}
