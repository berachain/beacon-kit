package encoding

import (
	"unsafe"

	"github.com/protolambda/ztyp/tree"
	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/crypto/hash/htr"
	"github.com/prysmaticlabs/prysm/v4/encoding/ssz"
)

// TransactionsRoot computes the HTR for the Transactions' property of the ExecutionPayload
// The code was largely copy/pasted from the code generated to compute the HTR of the entire
// ExecutionPayload.
func TransactionsRoot(txs [][]byte) ([32]byte, error) {
	txRoots := make([][32]byte, 0)
	for i := 0; i < len(txs); i++ {
		rt := tree.GetHashFn().ByteListHTR(txs[i], fieldparams.MaxBytesPerTxLength) // getting the transaction root here
		txRoots = append(txRoots, rt)
	}

	byteRoots, err := ssz.BitwiseMerkleize(txRoots, uint64(len(txRoots)), fieldparams.MaxTxsPerPayloadLength)
	if err != nil {
		return [32]byte{}, err
	}

	return tree.GetHashFn().Mixin(byteRoots, uint64(len(txRoots))), nil
}

func MerkleizeVector(roots []tree.Root, length uint64) tree.Root {
	depth := tree.CoverDepth(length)

	if len(roots) == 0 {
		return tree.ZeroHashes[depth]
	}

	// loop over i, depth
	for i := uint8(0); i < depth; i++ {
		oddLength := len(roots)%2 == 1
		if oddLength {
			roots = append(roots, tree.ZeroHashes[i])
		}

		// map htr.VectorizedSha256 to roots result
		res := htr.VectorizedSha256(*(*[][32]byte)(unsafe.Pointer(&roots)))
		roots = *(*[]tree.Root)(unsafe.Pointer(&res))
	}
	return roots[0]
}

// Hashable is an interface representing objects that implement HashTreeRoot()
type Hashable interface {
	HashTreeRoot() ([32]byte, error)
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

	return MerkleizeVector(convertBytesToTreeRoots(roots), length), nil
}

func convertTreeRootsToBytes(roots []tree.Root) [][32]byte {
	return *(*[][32]byte)(unsafe.Pointer(&roots))
}

func convertBytesToTreeRoots(bytes [][32]byte) []tree.Root {
	return *(*[]tree.Root)(unsafe.Pointer(&bytes))
}
