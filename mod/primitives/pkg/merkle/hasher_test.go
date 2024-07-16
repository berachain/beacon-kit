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
	"crypto/md5"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/stretchr/testify/require"
)

func TestCombi(t *testing.T) {
	// Initialize the hasher function
	hashFunc := func(data []byte) [32]byte {
		var result [32]byte
		hash := md5.Sum(data)
		copy(result[:], hash[:])
		return result
	}
	hasher := merkle.NewHasher[[32]byte](hashFunc)

	tests := []struct {
		name     string
		a, b     [32]byte
		expected [32]byte
	}{
		{
			name: "Simple combination",
			a:    [32]byte{1, 2, 3, 4},
			b:    [32]byte{5, 6, 7, 8},
			expected: [32]uint8{0xa3, 0x56, 0x78, 0x7f, 0x44, 0x4f, 0x5c, 0xa3,
				0x51, 0x4a, 0xfd, 0x73, 0x23, 0x18, 0xa7, 0x40},
		},
		{
			name: "Another combination",
			a:    [32]byte{9, 10, 11, 12},
			b:    [32]byte{13, 14, 15, 16},
			expected: [32]uint8{0xa, 0x63, 0x9a, 0x14, 0x35, 0xed, 0x3d, 0x9a,
				0x39, 0xfa, 0x9a, 0xa4, 0x60, 0xdd, 0xda},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := hasher.Combi(tc.a, tc.b)
			require.Equal(t, tc.expected, result,
				"TestCase %s", tc.name)
		})
	}
}

func TestMixIn(t *testing.T) {
	// Initialize the hasher function
	hashFunc := func(data []byte) [32]byte {
		var result [32]byte
		hash := md5.Sum(data)
		copy(result[:], hash[:])
		return result
	}
	hasher := merkle.NewHasher[[32]byte](hashFunc)

	tests := []struct {
		name     string
		a        [32]byte
		i        uint64
		expected [32]byte
	}{
		{
			name: "MixIn with integer 1",
			a:    [32]byte{1, 2, 3, 4},
			i:    1,
			expected: [32]uint8{0x90, 0xde, 0x5, 0xc9, 0x96, 0x7a, 0xc0, 0xdb,
				0x74, 0xf2, 0x8e, 0xf1, 0xd4, 0x62, 0xaf, 0x4d},
		},
		{
			name: "MixIn with integer 2",
			a:    [32]byte{5, 6, 7, 8},
			i:    2,
			expected: [32]uint8{0xaa, 0x93, 0x16, 0x9f, 0x1c, 0x2d, 0xc5, 0xb8,
				0xaa, 0x40, 0x39, 0x5, 0xd7, 0x80, 0x6f, 0xa9},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := hasher.MixIn(tc.a, tc.i)
			require.Equal(t, tc.expected, result,
				"TestCase %s", tc.name)
		})
	}
}
