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

	"github.com/berachain/beacon-kit/crypto/merkle"
	"github.com/prysmaticlabs/prysm/v5/container/trie"
	"github.com/stretchr/testify/require"
)

func TestMerkleTrie_VerifyMerkleProofWithDepth(t *testing.T) {
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
	var depth uint64 = 32
	m, err := merkle.GenerateTrieFromItems(items, depth)
	require.NoError(t, err)
	require.NotNil(t, m)
	root, err := m.HashTreeRoot()
	require.NoError(t, err)
	expectedRoot := [32]byte{
		0x90, 0x5e, 0xde, 0xcb, 0x94, 0x3e, 0x2e, 0x27,
		0xa6, 0x42, 0x0b, 0x16, 0x91, 0xff, 0xbf, 0xf2,
		0xc8, 0x38, 0xd3, 0x08, 0xd7, 0x48, 0xda, 0x31,
		0x74, 0xdb, 0x58, 0x9f, 0x5f, 0x6e, 0xa9, 0x23,
	}
	require.Equal(t, expectedRoot[:], root[:])
	proof, err := m.MerkleProof(0)
	require.NoError(t, err)
	require.Len(t, proof, int(depth)+1)
	require.True(t,
		merkle.VerifyMerkleProofWithDepth(
			root[:],
			items[0],
			0,
			proof,
			depth,
		),
	)
	proof, err = m.MerkleProof(3)
	require.NoError(t, err)
	require.True(t,
		trie.VerifyMerkleProofWithDepth(
			root[:],
			items[3],
			3,
			proof,
			depth,
		),
	)
	require.False(t,
		trie.VerifyMerkleProofWithDepth(
			root[:], []byte("buzz"), 3,
			proof, depth,
		),
	)
}
