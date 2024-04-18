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
	"fmt"
	"reflect"

	"github.com/berachain/beacon-kit/mod/merkle"
	"github.com/berachain/beacon-kit/mod/merkle/zero"
	"github.com/berachain/beacon-kit/mod/primitives"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
)

// TODO: implement the rest of these.
// merkleize(pack_bits(value), limit=chunk_count(type)) if value is a bitvector.
// mix_in_length(merkleize(pack_bits(value), limit=chunk_count(type)), len(value)) if value is a bitlist.
// mix_in_selector(hash_tree_root(value.value), value.selector) if value is of union type, and value.value is not None
// mix_in_selector(Bytes32(), 0) if value is of union type, and value.value is None

// Merkleize hashes a list of chunks and returns the HTR of the list of.
// As per the Ethereum 2.0 SSZ Specifcation:
// https://github.com/ethereum/consensus-specs/blob/dev/ssz/simple-serialize.md#merkleization
//
//nolint:lll
type SSZType interface {
	Marshallable
	Hashable[[32]byte]
}

// MerkelizeBasic hashes a basic object and returns the HTR.
//
// From Spec:
// merkleize(pack(value)) if value is a basic object or a vector of basic objects.
func MerkelizeBasic[T SSZType, RootT ~[32]byte](
	value T,
) (RootT, error) {
	return MerkleizeVectorBasic[T, RootT]([]T{value})
}

// MerkleizeVectorBasic hashes each element in the vector and then returns the HTR of the corresponding list of roots.
// From Spec:
// merkleize(pack(value)) if value is a basic object or a vector of basic objects.
func MerkleizeVectorBasic[T SSZType, RootT ~[32]byte](
	value []T,
) (RootT, error) {
	packed, err := Pack[T, RootT](value)
	if err != nil {
		return RootT{}, err
	}
	return Merkleize[RootT, RootT](packed)
}

// MerkleizeList hashes each element in the list and then returns the HTR of
// the list of corresponding roots, with the length mixed in.
//
// From Spec:
// mix_in_length(merkleize(pack(value), limit=chunk_count(type)), len(value)) if value is a list of basic objects.
func MerkleizeListBasic[T SSZType, RootT ~[32]byte](
	value []T, limit ...uint64,
) ([32]byte, error) {
	root, err := MerkleizeVectorBasic[T, [32]byte](value)
	if err != nil {
		return [32]byte{}, err
	}

	var effectiveLimit uint64
	if len(limit) == 0 {
		effectiveLimit = ChunkCount(value, "BasicVecList")
	} else {
		effectiveLimit = limit[0]
		if uint64(len(value)) > effectiveLimit {
			return [32]byte{}, fmt.Errorf("list length exceeds specified limit")
		}
	}

	return MixinLength(root, effectiveLimit), nil
}

// From Spec:
// mix_in_length(merkleize([hash_tree_root(element) for element in value], limit=chunk_count(type)), len(value)) if value is a list of composite objects.
func MerkleizeListComposite[T SSZType, RootT ~[32]byte](
	value []T, limit ...uint64,
) (RootT, error) {
	if len(value) == 0 {
		return zero.Hashes[0], nil
	}

	// Calculate the limit based on the chunk count for composite lists
	var effectiveLimit uint64
	if len(limit) == 0 {
		effectiveLimit = ChunkCount(value, "CompositeVecList")
	} else {
		effectiveLimit = limit[0]
		if uint64(len(value)) > effectiveLimit {
			return RootT{}, fmt.Errorf("list length exceeds specified limit")
		}
	}

	// Compute hash tree root for each element
	roots := make([][32]byte, len(value))
	for i, element := range value {
		var err error
		roots[i], err = element.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}

	// Merkleize the list of roots

	merkleRoot, err := Merkleize[[32]byte, RootT](roots, effectiveLimit)
	if err != nil {
		return RootT{}, err
	}

	// Mix in the length of the list
	return MixinLength(merkleRoot, uint64(len(value))), nil
}

// MerkleizeVector hashes each element in the vector and then returns the HTR of the corresponding list of roots.
// merkleize([hash_tree_root(element) for element in value]) if value is a vector of composite objects or a container.
func MerkleizeVector[T SSZType, RootT ~[32]byte](
	value []T,
) (RootT, error) {
	if len(value) == 0 {
		return zero.Hashes[0], nil
	}

	// Compute hash tree root for each element
	roots := make([][32]byte, len(value))
	for i, element := range value {
		var err error
		roots[i], err = element.HashTreeRoot()
		if err != nil {
			return RootT{}, err
		}
	}

	// Merkleize the list of roots
	merkleRoot, err := Merkleize[[32]byte, RootT](roots)
	if err != nil {
		return RootT{}, err
	}

	return merkleRoot, nil
}

// MerkleizeContainer hashes each field in the container and then returns the HTR of the corresponding list of roots.
func MerkleizeContainer[T SSZType, RootT ~[32]byte](
	value T,
) (RootT, error) {
	fields := make([]SSZType, 0)
	v := reflect.ValueOf(value)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i).Interface().(SSZType)
		fields = append(fields, field)
	}

	return MerkleizeVector[SSZType, RootT](fields)
}

// Merkleize hashes a list of chunks and returns the HTR of the list of
// corresponding roots.
func Merkleize[ChunkT, RootT ~[32]byte](
	chunks []ChunkT,
	limit ...uint64,
) (RootT, error) {
	var effectiveLimit primitives.U64
	if len(limit) == 0 {
		effectiveLimit = primitives.U64(len(chunks))
	} else {
		limit := limit[0]
		if limit < uint64(len(chunks)) {
			return RootT{}, fmt.Errorf("input exceeds limit")
		}
		effectiveLimit = primitives.U64(limit)
	}

	if effectiveLimit == 0 {
		return zero.Hashes[0], nil
	} else if effectiveLimit == 1 {
		if len(chunks) == 0 {
			return zero.Hashes[0], nil
		}
		return RootT(chunks[0]), nil
	}

	paddedChunks := make([]ChunkT, effectiveLimit)
	copy(paddedChunks, chunks)
	for i := len(chunks); i < len(paddedChunks); i++ {
		paddedChunks[i] = zero.Hashes[0]
	}

	// TODO: This is very inefficient. We can build the merkle root without building the tree.
	tree, err := merkle.NewTreeWithMaxLeaves[ChunkT, RootT](paddedChunks, uint64(effectiveLimit))
	if err != nil {
		return RootT{}, err
	}
	return tree.Root(), nil
}

// ------------------------------ Helpers ------------------------------

// MixinLength mixes the length into the root.
func MixinLength(root [32]byte, length uint64) [32]byte {
	return merkle.MixinLength(root, length)
}

// Pack packs a list of SSZ-marshallable elements into a single byte slice.
func Pack[S Marshallable, RootT ~[32]byte](s []S) ([]RootT, error) {
	// Pack each element into separate buffers.
	var buffers []RootT
	for _, el := range s {
		packed, err := el.MarshalSSZ()
		if err != nil {
			return nil, err
		}

		// Right pad each buffer to ensure it is a multiple of 32 bytes
		paddingSize := 32 - (len(packed) % 32)
		if paddingSize > 0 && paddingSize < 32 {
			padded := make([]byte, len(packed)+paddingSize)
			copy(padded, packed)
			buffers = append(buffers, RootT(padded))
		} else {
			buffers = append(buffers, RootT(packed))
		}
	}
	return buffers, nil
}

// chunk_count(type): calculate the amount of leafs for merkleization of the type.
// all basic types: 1
// Bitlist[N] and Bitvector[N]: (N + 255) // 256 (dividing by chunk size, rounding up)
// List[B, N] and Vector[B, N], where B is a basic type: (N * size_of(B) + 31) // 32 (dividing by chunk size, rounding up)
// List[C, N] and Vector[C, N], where C is a composite type: N
// containers: len(fields)
func ChunkCount[S Marshallable](obj []S, sszType string) uint64 {
	switch sszType {
	case "Basic":
		return 1
	case "Bitlist":
		return (uint64(len(obj)) + 255) / 256
	case "BasicVecList":
		size := int(0)
		for _, el := range obj {
			size += el.SizeSSZ()
		}
		return (uint64(size) + 31) / constants.RootLength
	case "CompositeVecList":
		return uint64(len(obj))
	case "Container":
		return uint64(reflect.TypeOf(obj).NumField())
	default:
		return 0
	}
}

// // MerkleizeVector hashes each element in the list and then returns the HTR
// // of the corresponding list of roots.
// func MerkleizeVector[T Hashable[[32]byte]](
// 	elements []T, length uint64,
// ) ([32]byte, error) {
// 	roots := make([][32]byte, len(elements))
// 	var err error
// 	for i, el := range elements {
// 		roots[i], err = el.HashTreeRoot()
// 		if err != nil {
// 			return [32]byte{}, err
// 		}
// 	}
// 	return Merkleize[[32]byte, [32]byte](roots, length)
// }

// MerkleizeByteSlice hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
// MerkleizeByteSlice hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func MerkleizeByteSlice(input []byte) ([32]byte, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	numChunks := (uint64(len(input)) + 31) / constants.RootLength
	if numChunks == 0 {
		return [32]byte{}, ErrInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return merkle.NewRootWithMaxLeaves[[32]byte, [32]byte](chunks, numChunks)
}
