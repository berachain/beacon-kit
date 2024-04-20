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

package tree

import "github.com/berachain/beacon-kit/mod/primitives"

type GeneralizedIndex = int

// Usage note: functions outside this section should manipulate generalized
// indices using only functions inside this section. This is to make it easier
// for developers to implement generalized indices with underlying
// representations other than bigints.

// concatGeneralizedIndices concatenates multiple generalized indices into a
// single generalized index.
func ConcatGeneralizedIndices(indices ...GeneralizedIndex) GeneralizedIndex {
	o := GeneralizedIndex(1)
	for _, i := range indices {
		o = GeneralizedIndex(
			o*GetPowerOfTwoFloor(i) + (i - GetPowerOfTwoFloor(i)),
		)
	}
	return o
}

// GetGeneralizedIndexLength returns the length of a path represented by a
// generalized index.
func GetGeneralizedIndexLength(index GeneralizedIndex) int {
	return int(primitives.U64(uint64(index)).ILog2Ceil())
}

// GetGeneralizedIndexBit returns the specified bit of a generalized index.
func GetGeneralizedIndexBit(index GeneralizedIndex, position int) bool {
	return (index & (1 << position)) > 0
}

// GeneralizedIndexSibling returns the sibling of a given generalized index.
func GeneralizedIndexSibling(index GeneralizedIndex) GeneralizedIndex {
	return GeneralizedIndex(index ^ 1)
}

// GeneralizedIndexChild returns the child index of a given generalized index,
// specifying if it's the right child.
func GeneralizedIndexChild(
	index GeneralizedIndex,
	rightSide bool,
) GeneralizedIndex {
	if rightSide {
		return GeneralizedIndex(index*2 + 1)
	}
	return GeneralizedIndex(index * 2)
}

// GeneralizedIndexParent returns the parent index of a given generalized index.
func GeneralizedIndexParent(index GeneralizedIndex) GeneralizedIndex {
	return GeneralizedIndex(index / 2)
}
