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

package merkle

import "github.com/berachain/beacon-kit/mod/primitives/math"

type GeneralizedIndex uint64

// NewGeneralizedIndex calculates the generalized index from the depth and
// index.
func NewGeneralizedIndex(depth uint8, index uint64) GeneralizedIndex {
	return GeneralizedIndex((1 << depth) + index)
}

func ConcatGeneralizedIndex(a, b GeneralizedIndex) GeneralizedIndex {
	return GeneralizedIndex(uint64(a)<<b.Length() + uint64(b))
}

// Length returns the length of the generalized index.
func (g GeneralizedIndex) Length() uint64 {
	return uint64(math.U64(g).ILog2Ceil())
}

// IndexBit returns the bit at the specified position in a generalized index.
func (g GeneralizedIndex) IndexBit(position uint) bool {
	return (g & (1 << position)) > 0
}

// Sibling returns the sibling index of the current generalized index.
func (g GeneralizedIndex) Sibling() GeneralizedIndex {
	return g ^ 1
}

// Child returns the child index of the current generalized index, specifying if
// it should return the right child.
func (g GeneralizedIndex) Child(rightSide bool) GeneralizedIndex {
	if rightSide {
		return 2*g + 1
	}
	return 2 * g
}

// Parent returns the parent index of the current generalized index.
func (g GeneralizedIndex) Parent() GeneralizedIndex {
	return g / 2
}
