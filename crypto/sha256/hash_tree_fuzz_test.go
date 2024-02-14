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

package sha256_test

import (
	"runtime"
	"testing"

	"github.com/itsdevbear/bolaris/crypto/sha256"
)

func FuzzHashTreeRoot(f *testing.F) {
	// Seed corpus with a variety of sizes, including edge cases
	//
	// Test with empty slice
	f.Add([]byte{}, 1)
	// Just below a single block size
	f.Add(make([]byte, 31), runtime.GOMAXPROCS(0)-1)
	// Exactly one block size
	f.Add(make([]byte, 32), runtime.GOMAXPROCS(0)+1)
	// Just above a single block size
	f.Add(make([]byte, 33), runtime.GOMAXPROCS(0)*2)
	// Multiple blocks
	f.Add(make([]byte, 64), runtime.GOMAXPROCS(0)*4)
	// Larger input
	f.Add(make([]byte, 1024), 3)
	// Just below MinParallelizationSize
	f.Add(make([]byte, sha256.MinParallelizationSize-2), 300)
	// Exactly MinParallelizationSize
	f.Add(make([]byte, sha256.MinParallelizationSize), 1)
	// Just above MinParallelizationSize
	f.Add(make([]byte, sha256.MinParallelizationSize+2), 64)
	// Double MinParallelizationSize
	f.Add(make([]byte, 2*sha256.MinParallelizationSize), runtime.GOMAXPROCS(0)-1)

	f.Fuzz(func(t *testing.T, original []byte, numRoutines int) {
		// Convert []byte to [][32]byte as required by HashTreeRoot
		var input [][32]byte
		for i := 0; i < len(original); i += 32 {
			var block [32]byte
			copy(block[:], original[i:min(i+32, len(original))])
			input = append(input, block)
		}

		// Ensure an even number of chunks for HashTreeRoot
		expectError := false
		if len(input)%2 != 0 {
			expectError = true
		}

		requireGoHashTreeEquivalence(t, input, numRoutines, expectError)
	})
}
