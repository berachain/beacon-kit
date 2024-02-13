// SPDX-License-Identifier: MIT
//
// Copyright (c) 2023 Berachain Foundation
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

package sha256_test

import (
	"testing"

	"github.com/itsdevbear/bolaris/crypto/sha256"
	"github.com/protolambda/ztyp/tree"
)

func FuzzHashTreeRoot(f *testing.F) {
	// Seed corpus with a variety of sizes, including edge cases
	f.Add([]byte{})                                      // Test with empty slice
	f.Add(make([]byte, 31))                              // Just below a single block size
	f.Add(make([]byte, 32))                              // Exactly one block size
	f.Add(make([]byte, 33))                              // Just above a single block size
	f.Add(make([]byte, 64))                              // Multiple blocks
	f.Add(make([]byte, 1024))                            // Larger input
	f.Add(make([]byte, sha256.MinParallelizationSize-2)) // Just below MinParallelizationSize
	f.Add(make([]byte, sha256.MinParallelizationSize))   // Exactly MinParallelizationSize
	f.Add(make([]byte, sha256.MinParallelizationSize+2)) // Just above MinParallelizationSize
	f.Add(make([]byte, 2*sha256.MinParallelizationSize)) // Double MinParallelizationSize

	f.Fuzz(func(t *testing.T, original []byte) {
		// Convert []byte to [][32]byte as required by HashTreeRoot
		var input []tree.Root
		for i := 0; i < len(original); i += 32 {
			var block tree.Root
			copy(block[:], original[i:min(i+32, len(original))])
			input = append(input, block)
		}

		// Ensure an even number of chunks for HashTreeRoot
		expectError := false
		if len(input)%2 != 0 {
			expectError = true
		}

		requireGoHashTreeEquivalence(t, input, expectError)
	})
}
