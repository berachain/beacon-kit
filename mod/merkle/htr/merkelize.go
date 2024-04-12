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

package htr

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/mod/merkle/zero"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	ztyp "github.com/protolambda/ztyp/tree"
	"github.com/prysmaticlabs/gohashtree"
)

// Vector uses our optimized routine to hash a list of 32-byte
// elements.
func Vector(elements [][32]byte, length uint64) [32]byte {
	depth := ztyp.CoverDepth(length)
	// Return zerohash at depth
	if len(elements) == 0 {
		return zero.Hashes[depth]
	}
	for i := range depth {
		layerLen := len(elements)
		oddNodeLength := layerLen%two == 1
		if oddNodeLength {
			zerohash := zero.Hashes[i]
			elements = append(elements, zerohash)
		}
		var err error
		elements, err = BuildParentTreeRoots(elements)
		if err != nil {
			return zero.Hashes[depth]
		}
	}
	if len(elements) != 1 {
		return zero.Hashes[depth]
	}
	return elements[0]
}

// ByteSliceSSZ hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func ByteSliceSSZ(input []byte) ([32]byte, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	numChunks := (uint64(len(input)) + 31) / constants.RootLength
	if numChunks == 0 {
		return [32]byte{}, ErrInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return Vector(chunks, numChunks), nil
}

// ListSSZ hashes each element in the list and then returns the HTR of
// the list of corresponding roots, with the length mixed in.
func ListSSZ[T Hashable](
	elements []T,
	limit uint64,
) ([32]byte, error) {
	body, err := VectorSSZ(elements, limit)
	if err != nil {
		return [32]byte{}, err
	}
	chunks := make([][32]byte, two)
	chunks[0] = body
	binary.LittleEndian.PutUint64(chunks[1][:], uint64(len(elements)))
	if err = gohashtree.Hash(chunks, chunks); err != nil {
		return [32]byte{}, err
	}
	return chunks[0], err
}

// VectorSSZ hashes each element in the list and then returns the HTR
// of the corresponding list of roots.
func VectorSSZ[T Hashable](
	elements []T,
	length uint64,
) ([32]byte, error) {
	roots := make([][32]byte, len(elements))
	var err error
	for i, el := range elements {
		roots[i], err = el.HashTreeRoot()
		if err != nil {
			return [32]byte{}, err
		}
	}
	return Vector(roots, length), nil
}
