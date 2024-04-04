package merkleize

import (
	"github.com/berachain/beacon-kit/mod/primitives/constants"
	"github.com/protolambda/ztyp/tree"
)

// Vector uses our optimized routine to hash a list of 32-byte
// elements.
func Vector(elements [][32]byte, length uint64) [32]byte {
	depth := tree.CoverDepth(length)
	// Return zerohash at depth
	if len(elements) == 0 {
		return tree.ZeroHashes[depth]
	}
	for i := uint8(0); i < depth; i++ {
		layerLen := len(elements)
		oddNodeLength := layerLen%two == 1
		if oddNodeLength {
			zerohash := tree.ZeroHashes[i]
			elements = append(elements, zerohash)
		}
		var err error
		elements, err = HashChunks(elements)
		if err != nil {
			return tree.ZeroHashes[depth]
		}
	}
	if len(elements) != 1 {
		return tree.ZeroHashes[depth]
	}
	return elements[0]
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

// ByteSliceVectorSSZ hashes a byteslice by chunkifying it and returning the
// corresponding HTR as if it were a fixed vector of bytes of the given length.
func ByteSliceVectorSSZ(input []byte) ([32]byte, error) {
	//nolint:gomnd // we add 31 in order to round up the division.
	numChunks := (uint64(len(input)) + 31) / constants.RootLength
	if numChunks == 0 {
		return [32]byte{}, errInvalidNilSlice
	}
	chunks := make([][32]byte, numChunks)
	for i := range chunks {
		copy(chunks[i][:], input[32*i:])
	}
	return Vector(chunks, numChunks), nil
}
