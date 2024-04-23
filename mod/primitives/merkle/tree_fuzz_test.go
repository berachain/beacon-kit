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

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/merkle"
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

	items := [][32]byte{
		byteslib.ToBytes32([]byte("A")),
		byteslib.ToBytes32([]byte("B")),
		byteslib.ToBytes32([]byte("C")),
		byteslib.ToBytes32([]byte("D")),
		byteslib.ToBytes32([]byte("E")),
		byteslib.ToBytes32([]byte("F")),
		byteslib.ToBytes32([]byte("G")),
		byteslib.ToBytes32([]byte("H")),
	}
	m, err := merkle.NewTreeFromLeavesWithDepth[[32]byte, [32]byte](
		items,
		depth,
	)
	require.NoError(f, err)
	proof, err := m.MerkleProofWithMixin(0)
	require.NoError(f, err)
	require.Len(f, proof, int(depth)+1)
	root, err := m.HashTreeRoot()
	require.NoError(f, err)
	var proofRaw []byte
	for _, p := range proof {
		proofRaw = append(proofRaw, p[:]...)
	}
	f.Add(root[:], items[0][:], uint64(0), proofRaw, depth)

	f.Fuzz(
		func(_ *testing.T,
			root, item []byte, merkleIndex uint64,
			proofRaw []byte, depth uint8,
		) {
			merkle.IsValidMerkleBranch(
				byteslib.ToBytes32(item),
				splitProofs(proofRaw),
				depth,
				merkleIndex,
				byteslib.ToBytes32(root),
			)
		},
	)
}
