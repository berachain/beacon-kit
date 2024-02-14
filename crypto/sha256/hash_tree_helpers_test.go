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
	"fmt"
	"sync"
	"testing"

	"github.com/itsdevbear/bolaris/crypto/sha256"
	"github.com/protolambda/ztyp/tree"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/stretchr/testify/require"
)

// requireGoHashTreeEquivalence is a helper function to ensure that the output of
// sha256.HashTreeRoot is equivalent to the output of gohashtree.Hash.
func requireGoHashTreeEquivalence(
	t *testing.T, inputList []tree.Root, numRoutines int, expectError bool,
) {
	expectedOutput := make([]tree.Root, len(inputList)/2)
	var output []tree.Root

	var wg sync.WaitGroup
	errChan := make(chan error, 2) // Buffer for 2 potential errors

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		output, err = sha256.HashTreeRootWithNRoutines(inputList, numRoutines)
		if err != nil {
			errChan <- fmt.Errorf("HashTreeRoot failed: %w", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := gohashtree.Hash(
			sha256.ConvertTreeRootsToBytes(expectedOutput),
			sha256.ConvertTreeRootsToBytes(inputList),
		)
		if err != nil {
			errChan <- fmt.Errorf("gohashtree.Hash failed: %w", err)
		}
	}()

	wg.Wait()      // Wait for both goroutines to finish
	close(errChan) // Close the channel

	// Check if there were any errors
	for err := range errChan {
		if !expectError {
			require.NoError(t, err, "Error occurred during hashing")
		} else {
			require.Error(t, err, "Expected error did not occur")
			return
		}
	}

	// Ensure the lengths are the same
	require.Equal(
		t, len(expectedOutput), len(output),
		fmt.Sprintf("Expected output length %d, got %d",
			len(expectedOutput), len(output)))

	// Compare the outputs element by element
	for i := range output {
		require.Equal(
			t, expectedOutput[i], output[i],
			fmt.Sprintf(
				"Output mismatch at index %d: expected %x, got %x",
				i, expectedOutput[i], output[i],
			),
		)
	}
}
