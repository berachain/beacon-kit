package tree

import "math"

type GeneralizedIndex = int

// Usage note: functions outside this section should manipulate generalized indices using only functions inside this section. This is to make it easier for developers to implement generalized indices with underlying representations other than bigints.

// concatGeneralizedIndices concatenates multiple generalized indices into a single generalized index.
func ConcatGeneralizedIndices(indices ...GeneralizedIndex) GeneralizedIndex {
	o := GeneralizedIndex(1)
	for _, i := range indices {
		o = GeneralizedIndex(o*GetPowerOfTwoFloor(i) + (i - GetPowerOfTwoFloor(i)))
	}
	return o
}

// GetGeneralizedIndexLength returns the length of a path represented by a generalized index.
func GetGeneralizedIndexLength(index GeneralizedIndex) int {
	return int(math.Log2(float64(index)))
}

// GetGeneralizedIndexBit returns the specified bit of a generalized index.
func GetGeneralizedIndexBit(index GeneralizedIndex, position int) bool {
	return (index & (1 << position)) > 0
}

// GeneralizedIndexSibling returns the sibling of a given generalized index.
func GeneralizedIndexSibling(index GeneralizedIndex) GeneralizedIndex {
	return GeneralizedIndex(index ^ 1)
}

// GeneralizedIndexChild returns the child index of a given generalized index, specifying if it's the right child.
func GeneralizedIndexChild(index GeneralizedIndex, rightSide bool) GeneralizedIndex {
	if rightSide {
		return GeneralizedIndex(index*2 + 1)
	}
	return GeneralizedIndex(index * 2)
}

// generalizedIndexParent returns the parent index of a given generalized index.
func generalizedIndexParent(index GeneralizedIndex) GeneralizedIndex {
	return GeneralizedIndex(index / 2)
}
