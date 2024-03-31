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

package trie

import (
	"bytes"
	"encoding/binary"
	"fmt"

	byteslib "github.com/berachain/beacon-kit/mod/primitives/bytes"
	sha256 "github.com/minio/sha256-simd"
	"github.com/protolambda/ztyp/tree"
)

// VerifyMerkleProof given a trie root, a leaf, the generalized merkle index
// of the leaf in the trie, and the proof itself.
func VerifyMerkleProof(
	root, leaf []byte,
	merkleIndex uint64,
	proof [][]byte,
) bool {
	if len(proof) == 0 {
		return false
	}
	return VerifyMerkleProofWithDepth(
		root,
		leaf,
		merkleIndex,
		proof,
		uint64(len(proof))-1,
	)
}

// VerifyMerkleProofWithDepth verifies a Merkle branch against a root of a trie.
func VerifyMerkleProofWithDepth(
	root, item []byte,
	merkleIndex uint64,
	proof [][]byte,
	depth uint64,
) bool {
	if uint64(len(proof)) != depth+1 {
		return false
	}
	node := byteslib.ToBytes32(item)
	for i := uint64(0); i <= depth; i++ {
		if (merkleIndex & 1) == 1 {
			node = sha256.Sum256(append(proof[i], node[:]...))
		} else {
			node = sha256.Sum256(append(node[:], proof[i]...))
		}
		merkleIndex /= 2
	}
	return bytes.Equal(root, node[:])
}

// MerkleProof computes a proof from a trie's branches using a Merkle index.
func (m *SparseMerkleTrie) MerkleProof(index uint64) ([][]byte, error) {
	numLeaves := uint64(len(m.branches[0]))
	if index >= numLeaves {
		return nil, fmt.Errorf(
			"merkle index out of range in trie, max range: %d, received: %d",
			numLeaves,
			index,
		)
	}
	merkleIndex := index
	proof := make([][]byte, m.depth+1)
	for i := uint64(0); i < m.depth; i++ {
		subIndex := (merkleIndex / (1 << i)) ^ 1
		if subIndex < uint64(len(m.branches[i])) {
			item := byteslib.ToBytes32(m.branches[i][subIndex])
			proof[i] = item[:]
		} else {
			proof[i] = tree.ZeroHashes[i][:]
		}
	}
	var enc [32]byte
	binary.LittleEndian.PutUint64(enc[:], uint64(len(m.originalItems)))
	proof[len(proof)-1] = enc[:]
	return proof, nil
}
