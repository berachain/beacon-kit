package tree

import (
	"crypto/sha256"

	"github.com/berachain/beacon-kit/mod/primitives"
)

func MerkleTree(leaves []primitives.Bytes32) []primitives.Bytes32 {
	/*
	   Return an array representing the tree nodes by generalized index:
	   [0, 1, 2, 3, 4, 5, 6, 7], where each layer is a power of 2. The 0 index is ignored. The 1 index is the root.
	   The result will be twice the size as the padded bottom layer for the input leaves.
	*/
	bottomLength := GetPowerOfTwoCeil(len(leaves))
	o := make([]primitives.Bytes32, bottomLength*2)
	copy(o[bottomLength:], leaves)

	for i := bottomLength - 1; i > 0; i-- {
		o[i] = sha256.Sum256(append(o[i*2][:], o[i*2+1][:]...))
	}
	return o
}

func GetPowerOfTwoCeil(x int) int {
	/*
	   Get the power of 2 for given input, or the closest higher power of 2 if the input is not a power of 2.
	   Commonly used for "how many nodes do I need for a bottom tree layer fitting x elements?"
	   Example: 0->1, 1->1, 2->2, 3->4, 4->4, 5->8, 6->8, 7->8, 8->8, 9->16.
	*/
	if x <= 1 {
		return 1
	} else if x == 2 {
		return 2
	} else {
		return 2 * GetPowerOfTwoCeil((x+1)/2)
	}
}

func GetPowerOfTwoFloor(x int) int {
	/*
	   Get the power of 2 for given input, or the closest lower power of 2 if the input is not a power of 2.
	   The zero case is a placeholder and not used for math with generalized indices.
	   Commonly used for "what power of two makes up the root bit of the generalized index?"
	   Example: 0->1, 1->1, 2->2, 3->2, 4->4, 5->4, 6->4, 7->4, 8->8, 9->8
	*/
	if x <= 1 {
		return 1
	}
	if x == 2 {
		return x
	} else {
		return 2 * GetPowerOfTwoFloor(x/2)
	}
}
