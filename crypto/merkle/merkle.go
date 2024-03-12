package merkle

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"

	bytesutil "github.com/berachain/beacon-kit/lib/bytes"
)

// SparseMerkleTrie implements a sparse, general purpose Merkle trie to be used
// across Ethereum consensus functionality.
type SparseMerkleTrie struct {
	depth         uint
	branches      [][][]byte
	originalItems [][]byte // list of provided items before hashing them into leaves.
}

// GenerateTrieFromItems constructs a Merkle trie from a sequence of byte slices.
func GenerateTrieFromItems(items [][]byte,
	depth uint64) (*SparseMerkleTrie, error) {
	if len(items) == 0 {
		return nil, errors.New("no items provided to generate Merkle trie")
	}
	if depth >= 63 {
		return nil, errors.New("supported merkle trie depth exceeded (max uint64 depth is 63, " +
			"theoretical max sparse merkle trie depth is 64)") // PowerOf2 would overflow
	}

	leaves := items
	layers := make([][][]byte, depth+1)
	transformedLeaves := make([][]byte, len(leaves))
	for i := range leaves {
		arr := bytesutil.ToBytes32(leaves[i])
		transformedLeaves[i] = arr[:]
	}
	layers[0] = transformedLeaves
	for i := uint64(0); i < depth; i++ {
		if len(layers[i])%2 == 1 {
			layers[i] = append(layers[i], ZeroHashes[i][:])
		}
		updatedValues := make([][]byte, 0)
		for j := 0; j < len(layers[i]); j += 2 {
			concat := sha256.Sum256(append(layers[i][j], layers[i][j+1]...))
			updatedValues = append(updatedValues, concat[:])
		}
		layers[i+1] = updatedValues
	}
	return &SparseMerkleTrie{
		branches:      layers,
		originalItems: items,
		depth:         uint(depth),
	}, nil
}

// MerkleProof computes a proof from a trie's branches using a Merkle index.
func (m *SparseMerkleTrie) MerkleProof(index int) ([][]byte, error) {
	if index < 0 {
		return nil, fmt.Errorf("merkle index is negative: %d", index)
	}
	leaves := m.branches[0]
	if index >= len(leaves) {
		return nil, fmt.Errorf("merkle index out of range in trie, max range: %d, received: %d", len(leaves), index)
	}
	merkleIndex := uint(index)
	proof := make([][]byte, m.depth+1)
	for i := uint(0); i < m.depth; i++ {
		subIndex := (merkleIndex / (1 << i)) ^ 1
		if subIndex < uint(len(m.branches[i])) {
			item := bytesutil.ToBytes32(m.branches[i][subIndex])
			proof[i] = item[:]
		} else {
			proof[i] = ZeroHashes[i][:]
		}
	}
	var enc [32]byte
	binary.LittleEndian.PutUint64(enc[:], uint64(len(m.originalItems)))
	proof[len(proof)-1] = enc[:]
	return proof, nil
}

// VerifyMerkleProofWithDepth verifies a Merkle branch against a root of a trie.
func VerifyMerkleProofWithDepth(root, item []byte, merkleIndex uint64,
	proof [][]byte, depth uint64) bool {
	if uint64(len(proof)) != depth+1 {
		return false
	}
	node := bytesutil.ToBytes32(item)
	for i := uint64(0); i <= depth; i++ {
		if (merkleIndex & 1) == 1 {
			node = sha256.Sum256(append(proof[i], node[:]...))
		} else {
			node = sha256.Sum256(append(node[:], proof[i]...))
		}
		merkleIndex /= 2
	}
	return bytes.Equal(root, node[:])
}

// VerifyMerkleProof given a trie root, a leaf, the generalized merkle index
// of the leaf in the trie, and the proof itself.
func VerifyMerkleProof(root, item []byte, merkleIndex uint64,
	proof [][]byte) bool {
	if len(proof) == 0 {
		return false
	}
	return VerifyMerkleProofWithDepth(root, item, merkleIndex,
		proof, uint64(len(proof)-1))
}
