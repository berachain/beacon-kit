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

	byteslib "github.com/berachain/beacon-kit/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/primitives/pkg/merkle"
	"github.com/stretchr/testify/require"
)

const depth = uint8(16)

func FuzzTree_IsValidMerkleBranch(f *testing.F) {
	splitProofs := func(proofRaw []byte) [][32]byte {
		var proofs [][32]byte
		for i := 0; i < len(proofRaw); i += 32 {
			end := i + 32
			if end >= len(proofRaw) {
				end = len(proofRaw) - 1
			}
			var proofSegment [32]byte
			copy(proofSegment[:], proofRaw[i:end])
			proofs = append(proofs, proofSegment)
		}
		return proofs
	}

	items := make([][32]byte, 0)
	for _, v := range [][]byte{
		byteslib.ExtendToSize([]byte("A"), byteslib.B32Size),
		byteslib.ExtendToSize([]byte("B"), byteslib.B32Size),
		byteslib.ExtendToSize([]byte("C"), byteslib.B32Size),
		byteslib.ExtendToSize([]byte("D"), byteslib.B32Size),
		byteslib.ExtendToSize([]byte("E"), byteslib.B32Size),
		byteslib.ExtendToSize([]byte("F"), byteslib.B32Size),
		byteslib.ExtendToSize([]byte("G"), byteslib.B32Size),
	} {
		item, err := byteslib.ToBytes32(v)
		require.NoError(f, err)
		items = append(items, item)
	}
	require.NotEmpty(f, items) // appease nilaway

	m, err := merkle.NewTreeFromLeavesWithDepth(items, depth)
	require.NoError(f, err)
	proof, err := m.MerkleProofWithMixin(0)
	require.NoError(f, err)
	require.Len(f, proof, int(depth)+1)
	root := m.HashTreeRoot()
	var proofRaw []byte
	for _, p := range proof {
		proofRaw = append(proofRaw, p[:]...)
	}
	f.Add(root[:], items[0][:], uint64(0), proofRaw, depth)

	f.Fuzz(
		func(t *testing.T,
			root, item []byte, merkleIndex uint64,
			proofRaw []byte, depth uint8,
		) {
			var r, leaf byteslib.B32

			item = byteslib.ExtendToSize(item, byteslib.B32Size)[:byteslib.B32Size]
			leaf, err = byteslib.ToBytes32(item)
			require.NoError(t, err)

			root = byteslib.ExtendToSize(root, byteslib.B32Size)[:byteslib.B32Size]
			r, err = byteslib.ToBytes32(root)
			require.NoError(t, err)

			merkle.IsValidMerkleBranch(
				leaf,
				splitProofs(proofRaw),
				depth,
				merkleIndex,
				r,
			)
		},
	)
}
