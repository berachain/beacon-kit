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

package trie_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/trie"
	"github.com/stretchr/testify/require"
)

const depth = uint64(16)

func FuzzSparseMerkleTrie_VerifyMerkleProofWithDepth(f *testing.F) {
	splitProofs := func(proofRaw []byte) [][]byte {
		var proofs [][]byte
		for i := 0; i < len(proofRaw); i += 32 {
			end := i + 32
			if end >= len(proofRaw) {
				end = len(proofRaw) - 1
			}
			proofs = append(proofs, proofRaw[i:end])
		}
		return proofs
	}

	items := [][]byte{
		[]byte("A"),
		[]byte("B"),
		[]byte("C"),
		[]byte("D"),
		[]byte("E"),
		[]byte("F"),
		[]byte("G"),
		[]byte("H"),
	}
	m, err := trie.NewFromItems(items, depth)
	require.NoError(f, err)
	proof, err := m.MerkleProof(0)
	require.NoError(f, err)
	require.Len(f, proof, int(depth)+1)
	root, err := m.HashTreeRoot()
	require.NoError(f, err)
	var proofRaw []byte
	for _, p := range proof {
		proofRaw = append(proofRaw, p...)
	}
	f.Add(root[:], items[0], uint64(0), proofRaw, depth)

	f.Fuzz(
		func(_ *testing.T,
			root, item []byte, merkleIndex uint64,
			proofRaw []byte, depth uint64,
		) {
			trie.VerifyMerkleProofWithDepth(
				root,
				item,
				merkleIndex,
				splitProofs(proofRaw),
				depth,
			)
		},
	)
}
