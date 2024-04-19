package tree

import (
	"crypto/sha256"
	"sort"

	"github.com/berachain/beacon-kit/mod/primitives"
)

// GetBranchIndices returns the generalized indices of the sister chunks along the path from the chunk with the
// given tree index to the root.
func GetBranchIndices(treeIndex GeneralizedIndex) []GeneralizedIndex {
	var o []GeneralizedIndex
	o = append(o, GeneralizedIndexSibling(treeIndex))
	for o[len(o)-1] > 1 {
		o = append(o, GeneralizedIndexSibling(generalizedIndexParent(o[len(o)-1])))
	}
	if len(o) > 1 {
		return o[:len(o)-1]
	}
	return o
}

// GetPathIndices returns the generalized indices of the chunks along the path from the chunk with the
// given tree index to the root.
func GetPathIndices(treeIndex GeneralizedIndex) []GeneralizedIndex {
	var o []GeneralizedIndex
	o = append(o, treeIndex)
	for o[len(o)-1] > 1 {
		o = append(o, generalizedIndexParent(o[len(o)-1]))
	}
	if len(o) > 1 {
		return o[:len(o)-1]
	}
	return o
}

// GetHelperIndices returns the generalized indices of all "extra" chunks in the tree needed to prove the chunks with the given
// generalized indices. Note that the decreasing order is chosen deliberately to ensure equivalence to the
// order of hashes in a regular single-item Merkle proof in the single-item case.
func GetHelperIndices(indices []GeneralizedIndex) []GeneralizedIndex {
	allHelperIndices := make(map[GeneralizedIndex]struct{})
	allPathIndices := make(map[GeneralizedIndex]struct{})

	for _, index := range indices {
		for _, idx := range GetBranchIndices(index) {
			allHelperIndices[idx] = struct{}{}
		}
		for _, idx := range GetPathIndices(index) {
			allPathIndices[idx] = struct{}{}
		}
	}

	var diff []GeneralizedIndex
	for idx := range allHelperIndices {
		if _, exists := allPathIndices[idx]; !exists {
			diff = append(diff, idx)
		}
	}

	sort.Slice(diff, func(i, j int) bool {
		return diff[i] > diff[j]
	})

	return diff
}

// calculateMerkleRoot calculates the Merkle root from a leaf and a proof based on the generalized index.
func calculateMerkleRoot(leaf primitives.Bytes32, proof []primitives.Bytes32, index GeneralizedIndex) primitives.Bytes32 {
	if len(proof) != GetGeneralizedIndexLength(index) {
		panic("proof length does not match the expected length from index")
	}
	for i, h := range proof {
		if GetGeneralizedIndexBit(index, i) {
			leaf = sha256.Sum256(append(h[:], leaf[:]...))
		} else {
			leaf = sha256.Sum256(append(leaf[:], h[:]...))
		}
	}
	return leaf
}

// verifyMerkleProof verifies a Merkle proof for a single leaf and a given root.
func verifyMerkleProof(leaf primitives.Bytes32, proof []primitives.Bytes32, index GeneralizedIndex, root primitives.Bytes32) bool {
	return calculateMerkleRoot(leaf, proof, index) == root
}

// calculateMultiMerkleRoot calculates the Merkle root for multiple leaves with a given proof and indices.
func calculateMultiMerkleRoot(leaves []primitives.Bytes32, proof []primitives.Bytes32, indices []GeneralizedIndex) primitives.Bytes32 {
	if len(leaves) != len(indices) {
		panic("leaves and indices length mismatch")
	}
	helperIndices := GetHelperIndices(indices)
	if len(proof) != len(helperIndices) {
		panic("proof length does not match helper indices length")
	}
	objects := make(map[GeneralizedIndex]primitives.Bytes32)
	for i, leaf := range leaves {
		objects[indices[i]] = leaf
	}
	for i, h := range proof {
		objects[helperIndices[i]] = h
	}

	keys := make([]GeneralizedIndex, 0, len(objects))
	for k := range objects {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	pos := 0
	for pos < len(keys) {
		k := keys[pos]
		if _, ok := objects[k]; ok && (objects[k^1] != primitives.Bytes32{}) && objects[k/2] == (primitives.Bytes32{}) {
			// Assuming left child is k|1^1 and right child is k|1
			leftIndex := GeneralizedIndex(k | 1 ^ 1)
			rightIndex := GeneralizedIndex(k | 1)
			leftChild := objects[leftIndex]
			rightChild := objects[rightIndex]

			// Hashing the concatenation of left and right child data
			hashed := sha256.Sum256(append(leftChild[:], rightChild[:]...))
			parentIndex := GeneralizedIndex(k / 2)
			objects[parentIndex] = hashed

			// Adding the parent index to keys
			keys = append(keys, parentIndex)
		}
		pos++
	}
	return objects[1]
}

// verifyMerkleMultiproof verifies a Merkle multiproof for multiple leaves and a given root.
func verifyMerkleMultiproof(leaves []primitives.Bytes32, proof []primitives.Bytes32, indices []GeneralizedIndex, root primitives.Bytes32) bool {
	return calculateMultiMerkleRoot(leaves, proof, indices) == root
}
