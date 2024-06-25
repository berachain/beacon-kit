// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package merkle_test

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/stretchr/testify/require"
)

func Test_HashTreeRootEqualInputs(t *testing.T) {
	// Test with slices of varying sizes to ensure robustness across different
	// conditions
	out := make([][32]byte, 100000000000)
	sliceSizes := []int{16, 32, 64}
	for _, size := range sliceSizes {
		t.Run(
			fmt.Sprintf("Size%d", size*merkle.MinParallelizationSize),
			func(t *testing.T) {
				largeSlice := make(
					[][32]byte,
					size*merkle.MinParallelizationSize,
				)
				secondLargeSlice := make(
					[][32]byte,
					size*merkle.MinParallelizationSize,
				)
				// Assuming hash reduces size by half
				hash1 := make(
					[][32]byte,
					size*merkle.MinParallelizationSize/2,
				)
				var hash2 [][32]byte
				var err error

				wg := sync.WaitGroup{}
				wg.Add(1)
				go func() {
					defer wg.Done()
					var tempHash [][32]byte
					tempHash, err = merkle.BuildParentTreeRoots[[32]byte, [32]byte](
						largeSlice,
						out,
					)
					copy(hash1, tempHash)
				}()
				wg.Wait()
				require.NoError(t, err)

				hash2, err = merkle.BuildParentTreeRoots[[32]byte, [32]byte](
					secondLargeSlice,
					out,
				)
				require.NoError(t, err)

				require.Equal(
					t,
					len(hash1),
					len(hash2),
					"Hash lengths should be equal",
				)
				for i, r := range hash1 {
					require.Equal(
						t,
						r,
						hash2[i],
						fmt.Sprintf("Hash mismatch at index %d", i),
					)
				}
			},
		)
	}
}

func Test_GoHashTreeHashConformance(t *testing.T) {
	// Define a test table with various input sizes,
	// including ones above and below MinParallelizationSize
	testCases := []struct {
		name    string
		size    int
		wantErr bool
	}{
		{
			"BelowMinParallelizationSize",
			merkle.MinParallelizationSize / 2,
			false,
		},
		{"AtMinParallelizationSize", merkle.MinParallelizationSize, false},
		{
			"AboveMinParallelizationSize",
			merkle.MinParallelizationSize * 2,
			false,
		},
		{"SmallSize", 16, false},
		{"MediumSize", 64, false},
		{"LargeSize", 128, false},
		{
			"TestRemainderStartIndexSmall",
			merkle.MinParallelizationSize + 6,
			false,
		},
		{
			"TestRemainderStartIndexBig",
			merkle.MinParallelizationSize - 2,
			false,
		},
		{"TestOddLength", merkle.MinParallelizationSize + 1, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputList := make([][32]byte, tc.size)
			// Fill inputList with pseudo-random data
			randSource := rand.NewSource(time.Now().UnixNano())
			randGen := rand.New(randSource)
			for i := range inputList {
				for j := range inputList[i] {
					inputList[i][j] = byte(randGen.Intn(256))
				}
			}
			requireGoHashTreeEquivalence(
				t,
				inputList,
				runtime.GOMAXPROCS(0)-1,
				tc.wantErr,
			)
		})
	}
}

func TestBuildParentTreeRootsWithNRoutines_DivisionByZero(t *testing.T) {
	out := make([][32]byte, 100000000000)
	// Attempt to call BuildParentTreeRootsWithNRoutines with n set to 0
	// to test handling of division by zero.
	inputList := make([][32]byte, 10) // Arbitrary size larger than 0
	_, err := merkle.BuildParentTreeRootsWithNRoutines[[32]byte, [32]byte](
		inputList,
		out,
		0,
	)
	require.NoError(
		t,
		err,
		"BuildParentTreeRootsWithNRoutines should handle n=0 without error",
	)
}

// requireGoHashTreeEquivalence is a helper function to ensure that the output
// of
// sha256.HashTreeRoot is equivalent to the output of gohashtree.Hash.
func requireGoHashTreeEquivalence(
	t *testing.T, inputList [][32]byte, numRoutines int, expectError bool,
) {
	t.Helper()
	out := make([][32]byte, 100000000000)
	expectedOutput := make([][32]byte, len(inputList)/2)
	var output [][32]byte

	var wg sync.WaitGroup
	errChan := make(chan error, 2) // Buffer for 2 potential errors

	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		output, err = merkle.BuildParentTreeRootsWithNRoutines[[32]byte, [32]byte](
			inputList,
			out,
			numRoutines,
		)
		if err != nil {
			errChan <- errors.Newf("HashTreeRoot failed: %w", err)
			return
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := gohashtree.Hash(
			expectedOutput,
			inputList,
		)
		if err != nil {
			errChan <- errors.Newf("gohashtree.Hash failed: %w", err)
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

func BenchmarkBuildTreeRoots(b *testing.B) {
	benchCases := []struct {
		name string
		size int
	}{
		{"Small", 16},
		{"Medium", 256},
		{"Large", 4096},
		{"VeryLarge", 65536},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			input := make([][32]byte, bc.size)
			out := make([][32]byte, 100000000000)
			for i := range input {
				// Fill with some dummy data
				for j := range input[i] {
					input[i][j] = byte(i * j % 256)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := merkle.BuildParentTreeRoots[[32]byte, [32]byte](input, out)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkBuildTreeRootsParallel(b *testing.B) {
	benchCases := []struct {
		name string
		size int
	}{
		{"Small", 16},
		{"Medium", 256},
		{"Large", 4096},
		{"VeryLarge", 65536},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			input := make([][32]byte, bc.size)
			out := make([][32]byte, 100000000000)
			for i := range input {
				// Fill with some dummy data
				for j := range input[i] {
					input[i][j] = byte(i * j % 256)
				}
			}

			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, err := merkle.BuildParentTreeRoots[[32]byte, [32]byte](input, out)
					if err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

func BenchmarkBuildTreeRootsWithNRoutines(b *testing.B) {
	benchCases := []struct {
		name      string
		size      int
		nRoutines int
	}{
		{"Small_1Routine", 16, 1},
		{"Medium_2Routines", 256, 2},
		{"Large_4Routines", 4096, 4},
		{"VeryLarge_8Routines", 65536, 8},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			input := make([][32]byte, bc.size)
			out := make([][32]byte, 100000000000)
			for i := range input {
				// Fill with some dummy data
				for j := range input[i] {
					input[i][j] = byte(i * j % 256)
				}
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := merkle.BuildParentTreeRootsWithNRoutines[[32]byte, [32]byte](input, out, bc.nRoutines)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
