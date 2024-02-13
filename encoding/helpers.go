package encoding

import (
	"unsafe"

	"github.com/protolambda/ztyp/tree"
)

// Hashable is an interface representing objects that implement HashTreeRoot()
type Hashable interface {
	HashTreeRoot() ([32]byte, error)
}

func convertTreeRootsToBytes(roots []tree.Root) [][32]byte {
	return *(*[][32]byte)(unsafe.Pointer(&roots))
}

func convertBytesToTreeRoots(bytes [][32]byte) []tree.Root {
	return *(*[]tree.Root)(unsafe.Pointer(&bytes))
}

// MerkleizeVectorSSZ hashes each element in the list and then returns the HTR
// of the corresponding list of roots
func MerkleizeVectorSSZ[T Hashable](elements []T, length uint64) ([32]byte, error) {
	roots := make([][32]byte, len(elements))
	var err error
	for i, el := range elements {
		roots[i], err = el.HashTreeRoot()
		if err != nil {
			return [32]byte{}, err
		}
	}

	return UnsafeMerkleizeVector(convertBytesToTreeRoots(roots), length), nil
}
