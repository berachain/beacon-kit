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

package htr

import (
	"github.com/berachain/beacon-kit/mod/merkle/bitlen"
	"github.com/berachain/beacon-kit/mod/merkle/zero"
)

// Vector uses our optimized routine to hash a list of 32-byte
// elements.
func Vector(elements [][32]byte, length uint64) [32]byte {
	depth := bitlen.CoverDepth(length)
	// Return zerohash at depth
	if len(elements) == 0 {
		return zero.Hashes[depth]
	}
	for i := range depth {
		layerLen := len(elements)
		oddNodeLength := layerLen%two == 1
		if oddNodeLength {
			zerohash := zero.Hashes[i]
			elements = append(elements, zerohash)
		}
		var err error
		elements, err = BuildParentTreeRoots(elements)
		if err != nil {
			return zero.Hashes[depth]
		}
	}
	if len(elements) != 1 {
		return zero.Hashes[depth]
	}
	return elements[0]
}
