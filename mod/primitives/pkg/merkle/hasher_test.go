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

	"github.com/berachain/beacon-kit/mod/errors"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle/zero"
	"github.com/prysmaticlabs/gohashtree"
	"github.com/stretchr/testify/require"
)

// Test NewRootWithMaxLeaves with empty leaves.
func TestNewRootWithMaxLeaves_EmptyLeaves(t *testing.T) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer, gohashtree.Hash)

	root, err := hasher.NewRootWithMaxLeaves(nil, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRoot := zero.Hashes[0]
	require.Equal(t, expectedRoot, root)
}

// Test NewRootWithDepth with empty leaves.
func TestNewRootWithDepth_EmptyLeaves(t *testing.T) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer, gohashtree.Hash)

	root, err := hasher.NewRootWithDepth([][32]byte{}, 0, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRoot := zero.Hashes[0]
	require.Equal(t, expectedRoot, root)
}

// Helper function to create a dummy leaf.
func createDummyLeaf(value byte) [32]byte {
	var leaf [32]byte
	leaf[0] = value
	return leaf
}

// Test NewRootWithMaxLeaves with one leaf.
func TestNewRootWithMaxLeaves_OneLeaf(t *testing.T) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer, gohashtree.Hash)

	leaf := createDummyLeaf(1)
	leaves := [][32]byte{leaf}

	root, err := hasher.NewRootWithMaxLeaves(leaves, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	require.Equal(t, leaf, root)
}

// Benchmark using a reusable buffer
//
// goos: darwin
// goarch: arm64
// pkg: github.com/berachain/beacon-kit/mod/primitives/pkg/merkle
// BenchmarkHasherWithReusableBuffer-12
// 29875  37987 ns/op  0 B/op  0 allocs/op.
func BenchmarkHasherWithReusableBuffer(b *testing.B) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer, gohashtree.Hash)

	leaves := make([][32]byte, 1000)
	for i := range 1000 {
		leaves[i] = createDummyLeaf(byte(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.NewRootWithMaxLeaves(leaves, math.U64(len(leaves)))
		require.NoError(b, err)
	}
}

// Benchmark using a single-use buffer
//
// goos: darwin
// goarch: arm64
// pkg: github.com/berachain/beacon-kit/mod/primitives/pkg/merkle
// BenchmarkHasherWithSingleUseBuffer-12
// 29114  38953 ns/op  16384 B/op  1 allocs/op.
func BenchmarkHasherWithSingleUseBuffer(b *testing.B) {
	buffer := getBuffer("singleuse")
	hasher := merkle.NewHasher(buffer, gohashtree.Hash)

	leaves := make([][32]byte, 1000)
	for i := range 1000 {
		leaves[i] = createDummyLeaf(byte(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.NewRootWithMaxLeaves(leaves, math.U64(len(leaves)))
		require.NoError(b, err)
	}
}

// getBuffer returns a buffer of the given type.
func getBuffer(usageType string) bytes.Buffer[[32]byte] {
	switch usageType {
	case "reusable":
		return bytes.NewReusableBuffer[[32]byte]()
	case "singleuse":
		return bytes.NewSingleuseBuffer[[32]byte]()
	default:
		panic("unknown usage type: " + usageType)
	}
}

func Test_HashTreeRootEqualInputs(t *testing.T) {
	// Test with slices of varying sizes to ensure robustness across different
	// conditions
	sliceSizes := []int{16, 32, 64}
	for _, size := range sliceSizes {
		t.Run(
			fmt.Sprintf("Size%d", size*merkle.MinParallelizationSize),
			func(t *testing.T) {
				largeSlice := make(
					[][32]byte, size*merkle.MinParallelizationSize,
				)
				secondLargeSlice := make(
					[][32]byte, size*merkle.MinParallelizationSize,
				)
				hash1 := make([][32]byte, size*merkle.MinParallelizationSize)
				hash2 := make([][32]byte, size*merkle.MinParallelizationSize)
				var err error

				err = merkle.BuildParentTreeRoots(hash1, largeSlice)
				require.NoError(t, err)

				err = merkle.BuildParentTreeRoots(hash2, secondLargeSlice)
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

func TestNewRootWithDepth(t *testing.T) {
	tests := []struct {
		name     string
		leaves   [][32]byte
		depth    int
		expected [32]byte
		wantErr  bool
	}{
		{
			name: "even number of leaves",
			leaves: [][32]byte{
				createDummyLeaf(1),
				createDummyLeaf(2),
			},
			depth: 1,
			expected: [32]uint8{
				0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				0x0, 0x0, 0x0, 0x0},
			wantErr: false,
		},
		{
			name: "not enough depth",
			leaves: [][32]byte{
				createDummyLeaf(1),
				createDummyLeaf(2),
				createDummyLeaf(3),
			},
			depth:    1,
			expected: zero.Hashes[1],
			wantErr:  true,
		},
		{
			name: "odd leaves",
			leaves: [][32]byte{
				createDummyLeaf(1),
				createDummyLeaf(2),
				createDummyLeaf(3),
			},
			depth:    2,
			expected: zero.Hashes[1],
			wantErr:  true,
		},
		{
			name: "hasher returns error",
			leaves: [][32]byte{
				createDummyLeaf(1),
			},
			depth:    1,
			expected: zero.Hashes[1],
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := getBuffer("reusable")
			hasher := merkle.NewHasher(buffer, func(dst, src [][32]byte) error {
				if tt.wantErr {
					return errors.New("hasher error")
				}
				copy(dst, src)
				return nil
			})

			root, err := hasher.NewRootWithDepth(
				tt.leaves,
				uint8(tt.depth),
				uint8(tt.depth),
			)
			if tt.wantErr {
				require.Error(t, err,
					"Test case %s", tt.name)
			} else {
				require.NoError(t, err,
					"Test case %s", tt.name)
				require.Equal(t, tt.expected, root,
					"Test case %s", tt.name)
			}
		})
	}
}

func TestNewRootWithMaxLeaves(t *testing.T) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer, gohashtree.Hash)

	tests := []struct {
		name     string
		leaves   [][32]byte
		limit    uint64
		expected error
	}{
		{
			name:     "Exceeds limit",
			leaves:   createLeaves(11), // Assuming limit is 10
			limit:    10,
			expected: errors.New("number of leaves exceeds limit"),
		},
		{
			name:     "Valid leaves within limit",
			leaves:   createLeaves(8), // Within limit
			limit:    10,
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := hasher.NewRootWithMaxLeaves(tc.leaves, math.U64(tc.limit))
			if tc.expected != nil {
				require.Error(t, err)
				require.Equal(t, tc.expected.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func createLeaves(count int) [][32]byte {
	leaves := make([][32]byte, count)
	for i := range leaves {
		leaves[i] = createDummyLeaf(byte(i))
	}
	return leaves
}
