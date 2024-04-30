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

package merkle

import (
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
)

// MerkleTree returns a Merkle tree of the given leaves.
// As defined in the Ethereum 2.0 Spec:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/merkle-proofs.md#generalized-merkle-tree-index
//
//nolint:lll
func Tree[LeafT ~[32]byte](
	leaves []LeafT,
	hashFn func([]byte) LeafT,
) []LeafT {
	/*
	   Return an array representing the tree nodes by generalized index:
	   [0, 1, 2, 3, 4, 5, 6, 7], where each layer is a power of 2. The 0 index is ignored. The 1 index is the root.
	   The result will be twice the size as the padded bottom layer for the input leaves.
	*/
	bottomLength := math.U64(len(leaves)).NextPowerOfTwo()
	//nolint:mnd // 2 is okay.
	o := make([]LeafT, bottomLength*2)
	copy(o[bottomLength:], leaves)
	for i := bottomLength - 1; i > 0; i-- {
		o[i] = hashFn(append(o[i*2][:], o[i*2+1][:]...))
	}
	return o
}
