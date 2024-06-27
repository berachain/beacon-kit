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
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/ssz"
	"github.com/stretchr/testify/require"
)

func TestNewGeneralizedIndex(t *testing.T) {
	tests := []struct {
		depth  uint8
		index  uint64
		expect ssz.GeneralizedIndex[[32]byte]
	}{
		{depth: 0, index: 0, expect: 1},
		{depth: 1, index: 1, expect: 3},
		{depth: 2, index: 2, expect: 6},
		{depth: 3, index: 5, expect: 13},
	}

	for _, tt := range tests {
		result := ssz.NewGeneralizedIndex[[32]byte](tt.depth, tt.index)
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
		indices ssz.GeneralizedIndicies[[32]byte]
		expect  ssz.GeneralizedIndex[[32]byte]
	}{
		{indices: []ssz.GeneralizedIndex[[32]byte]{1, 2, 3}, expect: 0x05},
		{indices: []ssz.GeneralizedIndex[[32]byte]{4, 5, 6}, expect: 0x46},
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
	gi := ssz.GeneralizedIndex[[32]byte](12) // Example index

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
		ssz.GeneralizedIndex[[32]byte](13),
		gi.Sibling(),
		"Incorrect sibling index",
	)
	require.Equal(
		t,
		ssz.GeneralizedIndex[[32]byte](24),
		gi.LeftChild(),
		"Incorrect right child index",
	)
	require.Equal(
		t,
		ssz.GeneralizedIndex[[32]byte](25),
		gi.RightChild(),
		"Incorrect left child index",
	)
	require.Equal(
		t,
		ssz.GeneralizedIndex[[32]byte](6),
		gi.Parent(),
		"Incorrect parent index",
	)
}
