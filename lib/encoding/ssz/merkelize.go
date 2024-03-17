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

package ssz

import (
	"encoding/binary"
	"errors"

	"github.com/berachain/beacon-kit/crypto/sha256"
	"github.com/protolambda/ztyp/tree"
	"github.com/prysmaticlabs/gohashtree"
)

var errInvalidNilSlice = errors.New("invalid empty slice")

const (
	mask0 = ^uint64((1 << (1 << iota)) - 1)
	mask1
	mask2
	mask3
	mask4
	mask5
)

const (
	bit0 = uint8(1 << iota)
	bit1
	bit2
	bit3
	bit4
	bit5
)

// Depth retrieves the appropriate depth for the provided trie size.
func Depth(v uint64) uint8 {
	// bitmagic: binary search through a uint32, offset down by 1 to not round
	// powers of 2 up. Then adding 1 to it to not get the index of the first
	// bit, but the length of the bits (depth of tree)
	// Zero is a special case, it has a 0 depth.
	// Example:
	// (in out): (0 0), (1 0), (2 1), (3 2), (4 2), (5 3), (6 3), (7 3), (8 3),
	// (9 4)
	out := uint8(0)
	if v <= 1 {
		return 0
	}
	v--
	if v&mask5 != 0 {
		v >>= bit5
		out |= bit5
	}
	if v&mask4 != 0 {
		v >>= bit4
		out |= bit4
	}
	if v&mask3 != 0 {
		v >>= bit3
		out |= bit3
	}
	if v&mask2 != 0 {
		v >>= bit2
		out |= bit2
	}
	if v&mask1 != 0 {
		v >>= bit1
		out |= bit1
	}
	if v&mask0 != 0 {
		out |= bit0
	}
	out++
	return out
}

// MerkleizeVector uses our optimized routine to hash a list of 32-byte
// elements.
func MerkleizeVector(elements [][32]byte, length uint64) [32]byte {
	depth := Depth(length)
	// Return zerohash at depth
	if len(elements) == 0 {
		return tree.ZeroHashes[depth]
	}
	for i := uint8(0); i < depth; i++ {
		layerLen := len(elements)
		oddNodeLength := layerLen%2 == 1
		if oddNodeLength {
			zerohash := tree.ZeroHashes[i]
			elements = append(elements, zerohash)
		}
		var err error
		elements, err = sha256.BuildParentTreeRoots(elements)
		if err != nil {
			return tree.ZeroHashes[depth]
		}
	}
	return elements[0]
}

// MerkleizeByteSliceSSZ hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func MerkleizeByteSliceSSZ(input []byte) ([32]byte, error) {
	numChunks := (len(input) + 31) / 32
	if numChunks == 0 {
		return [32]byte{}, errInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return MerkleizeVector(chunks, uint64(numChunks)), nil
}

// Hashable is an interface representing objects that implement HashTreeRoot().
type Hashable interface {
	HashTreeRoot() ([32]byte, error)
}

// MerkleizeListSSZ hashes each element in the list and then returns the HTR of
// the list of corresponding roots, with the length mixed in.
func MerkleizeListSSZ[T Hashable](
	elements []T,
	limit uint64,
) ([32]byte, error) {
	body, err := MerkleizeVectorSSZ(elements, limit)
	if err != nil {
		return [32]byte{}, err
	}
	chunks := make([][32]byte, 2)
	chunks[0] = body
	binary.LittleEndian.PutUint64(chunks[1][:], uint64(len(elements)))
	if err = gohashtree.Hash(chunks, chunks); err != nil {
		return [32]byte{}, err
	}
	return chunks[0], err
}

// MerkleizeVectorSSZ hashes each element in the list and then returns the HTR
// of the corresponding list of roots.
func MerkleizeVectorSSZ[T Hashable](
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
	return MerkleizeVector(roots, length), nil
}
