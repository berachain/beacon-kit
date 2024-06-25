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
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/stretchr/testify/require"
)

func Test_HashTreeRootEqualInputs(t *testing.T) {
	// Test with slices of varying sizes to ensure robustness across different
	// conditions
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
				hash1 := make([][32]byte, size*merkle.MinParallelizationSize)
				hash2 := make([][32]byte, size*merkle.MinParallelizationSize)
				var err error

				err = merkle.BuildParentTreeRoots[[32]byte](
					hash1,
					largeSlice,
				)
				require.NoError(t, err)

				err = merkle.BuildParentTreeRoots[[32]byte](
					hash2,
					secondLargeSlice,
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
	// Attempt to call BuildParentTreeRootsWithNRoutines with n set to 0
	// to test handling of division by zero.
	inputList := make([][32]byte, 10) // Arbitrary size larger than 0
	output := make([][32]byte, 8)     // Arbitrary size smaller than inputList
	err := merkle.BuildParentTreeRootsWithNRoutines(
		output,
		inputList,
		0,
	)
	require.NoError(
		t,
		err,
		"BuildParentTreeRootsWithNRoutines should handle n=0 without error",
	)
}

// requireGoHashTreeEquivalence is a helper function to ensure that the output
// of merkle.BuildParentTreeRootsWithNRoutines is equivalent to the output of
// gohashtree.Hash.
func requireGoHashTreeEquivalence(
	t *testing.T, inputList [][32]byte, numRoutines int, expectError bool,
) {
	t.Helper()

	// Deep copy inputList
	inputListCopy := make([][32]byte, len(inputList))
	copy(inputListCopy, inputList)

	expectedOutput := make([][32]byte, len(inputListCopy)/2)
	output := make([][32]byte, len(inputListCopy)/2)
	var err1, err2 error

	// Run merkle.BuildParentTreeRootsWithNRoutines
	err1 = merkle.BuildParentTreeRootsWithNRoutines(
		output,
		inputListCopy,
		numRoutines,
	)

	// Run gohashtree.Hash
	err2 = gohashtree.Hash(
		expectedOutput,
		inputListCopy,
	)

	// Check for errors
	if !expectError {
		require.NoError(t, err1, "BuildParentTreeRootsWithNRoutines failed")
		require.NoError(t, err2, "gohashtree.Hash failed")
	} else {
		if err1 == nil && err2 == nil {
			t.Error("Expected error did not occur")
		}
		return
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
