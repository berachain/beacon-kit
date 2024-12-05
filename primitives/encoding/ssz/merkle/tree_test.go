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

package merkle_test

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/encoding/ssz/merkle"
	"github.com/stretchr/testify/require"
)

func TestCalculateMerkleRoot(t *testing.T) {
	tests := []struct {
		name      string
		index     merkle.GeneralizedIndex
		leaf      [32]byte
		proof     [][32]byte
		expect    [32]byte
		expectErr bool
	}{
		{
			name:  "Valid Proof",
			index: merkle.GeneralizedIndex(3),
			leaf:  [32]byte{0x01},
			proof: [][32]byte{{0x02}},
			expect: [32]byte{
				0x95, 0xe7, 0x3e, 0x86, 0x16, 0xbb, 0x92, 0x7b, 0xb0, 0x74, 0xee,
				0x5, 0x5b, 0x12, 0x23, 0xf3, 0xa0, 0x85, 0xf7, 0x10, 0xc, 0x97,
				0x46, 0x8d, 0x92, 0xe6, 0x3a, 0x1c, 0x87, 0xaf, 0x1c, 0x1a},
			expectErr: false,
		},
		{
			name:      "Invalid Proof Length",
			index:     merkle.GeneralizedIndex(3),
			leaf:      [32]byte{0x01},
			proof:     [][32]byte{{0x02}, {0x03}},
			expect:    [32]byte{},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := merkle.CalculateRoot(
				tt.index,
				tt.leaf,
				tt.proof,
			)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expect, result)
			}
		})
	}
}

func TestVerifyMerkleProof(t *testing.T) {
	tests := []struct {
		name        string
		index       merkle.GeneralizedIndex
		leaf        [32]byte
		proof       [][32]byte
		root        [32]byte
		expectValid bool
		expectErr   bool
	}{
		{
			name:  "Valid Proof",
			index: merkle.GeneralizedIndex(3),
			leaf:  [32]byte{0x01},
			proof: [][32]byte{{0x02}},
			root: [32]byte{
				0x95, 0xe7, 0x3e, 0x86, 0x16, 0xbb, 0x92, 0x7b, 0xb0, 0x74, 0xee,
				0x5, 0x5b, 0x12, 0x23, 0xf3, 0xa0, 0x85, 0xf7, 0x10, 0xc, 0x97,
				0x46, 0x8d, 0x92, 0xe6, 0x3a, 0x1c, 0x87, 0xaf, 0x1c, 0x1a},
			expectErr:   false,
			expectValid: true,
		},
		{
			name:        "Invalid Proof",
			index:       merkle.GeneralizedIndex(3),
			leaf:        [32]byte{0x01},
			proof:       [][32]byte{{0x02}, {0x04}},
			root:        [32]byte{0x01},
			expectErr:   true,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := merkle.VerifyProof(
				tt.index,
				tt.leaf,
				tt.proof,
				tt.root,
			)
			require.Equal(t, tt.expectValid, result)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCalculateMultiMerkleRoot(t *testing.T) {
	tests := []struct {
		name      string
		indices   merkle.GeneralizedIndices
		leaves    [][32]byte
		proof     [][32]byte
		expect    [32]byte
		expectErr bool
	}{
		{
			name:    "Valid Multi Merkle Root",
			indices: []merkle.GeneralizedIndex{3, 5},
			leaves:  [][32]byte{{0x01}, {0x02}},
			proof:   [][32]byte{{0x03}},
			expect: [32]byte{
				0xcb, 0x54, 0x4b, 0x4c, 0x58, 0x15, 0x6f, 0x7d, 0x17, 0xd9,
				0x3b, 0xb4, 0x77, 0xbb, 0xb3, 0xeb, 0x50, 0x16, 0x47, 0x1a,
				0xd4, 0x58, 0xf5, 0x94, 0xb1, 0xf0, 0x29, 0x4a, 0x82, 0xcf,
				0x9f, 0x13,
			},
			expectErr: false,
		},
		{
			name:      "Mismatched Leaves and Indices Length",
			indices:   []merkle.GeneralizedIndex{3, 5},
			leaves:    [][32]byte{{0x01}},
			proof:     [][32]byte{{0x03}},
			expect:    [32]byte{},
			expectErr: true,
		},
		{
			name:      "Mismatched Proof and Helper Indices Length",
			indices:   []merkle.GeneralizedIndex{3, 5},
			leaves:    [][32]byte{{0x01}, {0x02}},
			proof:     [][32]byte{},
			expect:    [32]byte{},
			expectErr: true,
		},
		{
			name:      "Empty Indices and Leaves",
			indices:   []merkle.GeneralizedIndex{},
			leaves:    [][32]byte{},
			proof:     [][32]byte{},
			expect:    [32]byte{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := merkle.CalculateMultiRoot(
				tt.indices,
				tt.leaves,
				tt.proof,
			)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expect, result)
			}
		})
	}
}

func TestVerifyMerkleMultiproof(t *testing.T) {
	tests := []struct {
		name    string
		indices merkle.GeneralizedIndices
		leaves  [][32]byte
		proof   [][32]byte
		root    [32]byte
		expect  bool
	}{
		{
			name:    "Valid Merkle Multiproof",
			indices: []merkle.GeneralizedIndex{3, 5},
			leaves:  [][32]byte{{0x01}, {0x02}},
			proof:   [][32]byte{{0x03}},
			root: [32]byte{
				0xcb, 0x54, 0x4b, 0x4c, 0x58, 0x15, 0x6f, 0x7d, 0x17, 0xd9,
				0x3b, 0xb4, 0x77, 0xbb, 0xb3, 0xeb, 0x50, 0x16, 0x47, 0x1a,
				0xd4, 0x58, 0xf5, 0x94, 0xb1, 0xf0, 0x29, 0x4a, 0x82, 0xcf,
				0x9f, 0x13,
			},
			expect: true,
		},
		{
			name:    "Invalid Merkle Multiproof",
			indices: []merkle.GeneralizedIndex{3, 5},
			leaves:  [][32]byte{{0x01}, {0x02}},
			proof:   [][32]byte{{0x03}},
			root:    [32]byte{},
			expect:  false,
		},
		{
			name:    "Mismatched Leaves and Indices Length",
			indices: []merkle.GeneralizedIndex{3, 5},
			leaves:  [][32]byte{{0x01}},
			proof:   [][32]byte{{0x03}},
			root:    [32]byte{},
			expect:  false,
		},
		{
			name:    "Mismatched Proof and Helper Indices Length",
			indices: []merkle.GeneralizedIndex{3, 5},
			leaves:  [][32]byte{{0x01}, {0x02}},
			proof:   [][32]byte{},
			root:    [32]byte{},
			expect:  false,
		},
		{
			name:    "Empty Indices and Leaves",
			indices: []merkle.GeneralizedIndex{},
			leaves:  [][32]byte{},
			proof:   [][32]byte{},
			root:    [32]byte{},
			expect:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := merkle.VerifyMultiproof(
				tt.indices,
				tt.leaves,
				tt.proof,
				tt.root,
			)
			require.Equal(t, tt.expect, result)
		})
	}
}
