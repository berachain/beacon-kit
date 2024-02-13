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
	"github.com/prysmaticlabs/gohashtree"
)

func FuzzVectorizedSha256(f *testing.F) {
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
		// Convert []byte to [][32]byte as required by VectorizedSha256
		var input [][32]byte
		for i := 0; i < len(original); i += 32 {
			var block [32]byte
			copy(block[:], original[i:min(i+32, len(original))])
			input = append(input, block)
		}

		// Ensure an even number of chunks for VectorizedSha256
		if len(input)%2 != 0 {
			// Add an extra block of zeros if the number of chunks is odd
			input = append(input, [32]byte{})
		}

		// Execute the function under test
		output, err := sha256.VectorizedSha256(input)
		if err != nil {
			t.Fatalf("VectorizedSha256 failed: %v", err)
		}

		// Compare the output of VectorizedSha256 with gohashtree.Hash
		expectedOutput := make([][32]byte, len(input)/2)
		err = gohashtree.Hash(expectedOutput, input)
		if err != nil {
			t.Fatalf("gohashtree.Hash failed: %v", err)
		}

		// Ensure the lengths are the same
		if len(output) != len(expectedOutput) {
			t.Fatalf("Expected output length %d, got %d", len(expectedOutput), len(output))
		}

		// Compare the outputs element by element
		for i := range output {
			if output[i] != expectedOutput[i] {
				t.Errorf("Output mismatch at index %d: expected %x, got %x", i, expectedOutput[i], output[i])
			}
		}
	})
}

// Helper function to get the minimum of two integers.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
