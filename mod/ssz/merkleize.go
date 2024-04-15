package ssz

import (
	"encoding/binary"

	"github.com/berachain/beacon-kit/mod/merkle/htr"
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/prysmaticlabs/gohashtree"
)

// two is a commonly used constant.
const two = 2

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
	return htr.Vector(chunks, numChunks), nil
}

// MerkleizeList hashes each element in the list and then returns the HTR of
// the list of corresponding roots, with the length mixed in.
func MerkleizeList[T Hashable[[32]byte]](
	elements []T, limit uint64,
) ([32]byte, error) {
	body, err := MerkleizeVector(elements, limit)
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

// MerkleizeVector hashes each element in the list and then returns the HTR
// of the corresponding list of roots.
func MerkleizeVector[T Hashable[[32]byte]](
	elements []T, length uint64,
) ([32]byte, error) {
	roots := make([][32]byte, len(elements))
	var err error
	for i, el := range elements {
		roots[i], err = el.HashTreeRoot()
		if err != nil {
			return [32]byte{}, err
		}
	}
	return htr.Vector(roots, length), nil
}
