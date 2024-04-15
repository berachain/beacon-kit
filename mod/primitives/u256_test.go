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

package primitives_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives"
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
		le := primitives.NewU256L(tc.input)
		expected := new(huint256.Int).SetBytes(tc.expected)
		require.Equal(t, expected, le.ToU256())
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
		le := primitives.NewU256L(tc.input)
		expected := new(huint256.Int).SetBytes(tc.expected)
		require.Equal(t, expected.ToBig(), le.ToBig())
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
		le := primitives.NewU256L(tc.input)
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
		le := new(primitives.U256L)
		err := le.UnmarshalJSON([]byte(tc.json))
		require.NoError(t, err)
		expected := primitives.NewU256L(tc.expected)
		require.Equal(t, expected, *le)
	}
}
