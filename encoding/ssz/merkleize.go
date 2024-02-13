// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package ssz

import (
	"errors"

	"github.com/protolambda/ztyp/tree"
	"github.com/prysmaticlabs/prysm/v4/crypto/hash/htr"

	// We need to import this package to use the VectorizedSha256 function.
	_ "github.com/minio/sha256-simd"
)

func SafeMerkleizeVector(roots []tree.Root, length, maxLength uint64) (tree.Root, error) {
	if length > maxLength {
		return tree.Root{}, errors.New("merkleizing list that is too large, over limit")
	}
	return UnsafeMerkleizeVector(roots, maxLength), nil
}

func UnsafeMerkleizeVector(roots []tree.Root, length uint64) tree.Root {
	depth := tree.CoverDepth(length)

	if len(roots) == 0 {
		return tree.ZeroHashes[depth]
	}

	// loop over i, depth
	for i := uint8(0); i < depth; i++ {
		oddLength := len(roots)%2 == 1 //nolint:gomnd // 2 is the divisor.
		if oddLength {
			x := tree.ZeroHashes[i]
			roots = append(roots, x)
		}

		// TODO: move this because gpl
		res := htr.VectorizedSha256(convertTreeRootsToBytes(roots))
		roots = convertBytesToTreeRoots(res)
	}
	return roots[0]
}

// SafeMerkelizeVectorAndMixinLength hashes each element in the list and then returns the HTR
// with the length mixed in.
func SafeMerkelizeVectorAndMixinLength(txRoots []tree.Root, limit uint64) ([32]byte, error) {
	txRootLen := uint64(len(txRoots))
	byteRoots, err := SafeMerkleizeVector(txRoots, txRootLen, limit)
	if err != nil {
		return [32]byte{}, err
	}
	return tree.GetHashFn().Mixin(byteRoots, txRootLen), nil
}
