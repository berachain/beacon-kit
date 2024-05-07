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

package merkle_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz/merkle"
	"github.com/stretchr/testify/require"
)

func TestNewGeneralizedIndex(t *testing.T) {
	tests := []struct {
		depth  uint8
		index  uint64
		expect merkle.GeneralizedIndex[[32]byte]
	}{
		{depth: 0, index: 0, expect: 1},
		{depth: 1, index: 1, expect: 3},
		{depth: 2, index: 2, expect: 6},
		{depth: 3, index: 5, expect: 13},
	}

	for _, tt := range tests {
		result := merkle.NewGeneralizedIndex[[32]byte](tt.depth, tt.index)
		require.Equal(
			t,
			tt.expect,
			result,
			"Failed at depth %d and index %d",
			tt.depth,
			tt.index,
		)
	}
}

func TestConcatGeneralizedIndices(t *testing.T) {
	tests := []struct {
		indices merkle.GeneralizedIndicies[[32]byte]
		expect  merkle.GeneralizedIndex[[32]byte]
	}{
		{indices: []merkle.GeneralizedIndex[[32]byte]{1, 2, 3}, expect: 0x05},
		{indices: []merkle.GeneralizedIndex[[32]byte]{4, 5, 6}, expect: 0x46},
	}

	for _, tt := range tests {
		result := tt.indices.Concat()
		require.Equal(
			t,
			tt.expect,
			result,
			"Failed with indices %v",
			tt.indices,
		)
	}
}

func TestGeneralizedIndexMethods(t *testing.T) {
	gi := merkle.GeneralizedIndex[[32]byte](12) // Example index

	require.Equal(
		t,
		uint64(3),
		gi.Length(),
		"Incorrect length for GeneralizedIndex",
	)
	require.True(
		t,
		gi.IndexBit(2),
		"IndexBit should return true for bit position 2",
	)
	require.False(
		t,
		gi.IndexBit(1),
		"IndexBit should return false for bit position 1",
	)
	require.Equal(
		t,
		merkle.GeneralizedIndex[[32]byte](13),
		gi.Sibling(),
		"Incorrect sibling index",
	)
	require.Equal(
		t,
		merkle.GeneralizedIndex[[32]byte](24),
		gi.LeftChild(),
		"Incorrect right child index",
	)
	require.Equal(
		t,
		merkle.GeneralizedIndex[[32]byte](25),
		gi.RightChild(),
		"Incorrect left child index",
	)
	require.Equal(
		t,
		merkle.GeneralizedIndex[[32]byte](6),
		gi.Parent(),
		"Incorrect parent index",
	)
}
