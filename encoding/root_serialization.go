package encoding

import (
	"errors"

	_ "github.com/minio/sha256-simd"
	"github.com/protolambda/ztyp/tree"

	fieldparams "github.com/prysmaticlabs/prysm/v4/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v4/crypto/hash/htr"
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

	byteRoots, err := SafeMerkleizeVector(
		convertBytesToTreeRoots(txRoots), uint64(len(txRoots)), fieldparams.MaxTxsPerPayloadLength,
	)
	if err != nil {
		return [32]byte{}, err
	}
	return tree.GetHashFn().Mixin(byteRoots, uint64(len(txRoots))), nil
}

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
		oddLength := len(roots)%2 == 1
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
