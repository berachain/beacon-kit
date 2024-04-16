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

package bitlen_test

import (
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/mod/merkle/bitlen"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	v uint64
	d uint8
	l uint8
	i uint8
}

//nolint:gochecknoglobals // test cases.
var testCases = []testCase{
	{v: 0, d: 0, l: 0, i: 0}, // 0
	{v: 1, d: 0, l: 1, i: 0}, // 1
	{v: 2, d: 1, l: 2, i: 1}, // 10
	{v: 3, d: 2, l: 2, i: 1}, // 11
	{v: 4, d: 2, l: 3, i: 2}, // 100
	{v: 5, d: 3, l: 3, i: 2}, // 101
	{v: 6, d: 3, l: 3, i: 2}, // 110
	{v: 7, d: 3, l: 3, i: 2}, // 111
	{v: 8, d: 3, l: 4, i: 3}, // 1000
	{v: 9, d: 4, l: 4, i: 3}, // 1001
	{v: ^uint64(0), d: 64, l: 64, i: 63},
}

//nolint:gochecknoinits // test cases.
func init() {
	for i := uint8(4); i < 64; i++ {
		testCases = append(testCases,
			testCase{v: (1 << i) - 1, d: i, l: i, i: i - 1},
			testCase{v: 1 << i, d: i, l: i + 1, i: i},
			testCase{v: (1 << i) + 1, d: i + 1, l: i + 1, i: i},
		)
	}
}

func TestCoverDepth(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d_%d", tc.v, tc.d), func(t *testing.T) {
			require.Equal(
				t,
				tc.d,
				bitlen.CoverDepth(tc.v),
				"Expected depth for v %d (bin %b)",
				tc.v,
				tc.v,
			)
		})
	}
}

func TestBitLength(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d_%d", tc.v, tc.l), func(t *testing.T) {
			require.Equal(
				t,
				tc.l,
				bitlen.BitLength(tc.v),
				"Expected length for v %d (bin %b)",
				tc.v,
				tc.v,
			)
		})
	}
}

func TestBitIndex(t *testing.T) {
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("case_%d_%d", tc.v, tc.i), func(t *testing.T) {
			require.Equal(
				t,
				tc.i,
				bitlen.BitIndex(tc.v),
				"Expected index for v %d (bin %b)",
				tc.v,
				tc.v,
			)
		})
	}
}
