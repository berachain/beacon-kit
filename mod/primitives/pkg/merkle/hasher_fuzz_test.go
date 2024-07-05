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
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

func FuzzHashTreeRoot(f *testing.F) {
	// Seed corpus with a variety of sizes, including edge cases
	//
	// Test with empty slice
	f.Add(make([]byte, 0), true, merkle.MinParallelizationSize)
	// Just below a single block size
	f.Add(
		make([]byte, 31), true, merkle.MinParallelizationSize,
	)
	// Exactly one block size
	f.Add(
		make([]byte, 32), true, merkle.MinParallelizationSize,
	)
	// Just above a single block size
	f.Add(
		make([]byte, 33), true, merkle.MinParallelizationSize,
	)
	// Multiple blocks
	f.Add(
		make([]byte, 64), true, merkle.MinParallelizationSize,
	)
	// Larger input
	f.Add(
		make([]byte, 1024), true, merkle.MinParallelizationSize,
	)
	// Just below MinParallelizationSize leaves
	f.Add(
		make([]byte, merkle.MinParallelizationSize-2), false,
		merkle.MinParallelizationSize,
	)
	// Exactly MinParallelizationSize leaves
	f.Add(
		make([]byte, merkle.MinParallelizationSize), false,
		merkle.MinParallelizationSize,
	)
	// Just above MinParallelizationSize leaves
	f.Add(
		make([]byte, merkle.MinParallelizationSize+2), false,
		merkle.MinParallelizationSize,
	)
	// Double MinParallelizationSize leaves
	f.Add(
		make([]byte, 2*merkle.MinParallelizationSize), false,
		merkle.MinParallelizationSize,
	)
	// Max Txs leaves
	f.Add(
		make([]byte, int(constants.MaxTxsPerPayload)), false,
		merkle.MinParallelizationSize,
	)

	f.Fuzz(func(
		t *testing.T,
		original []byte, isLeaves bool, minParallelizationSize int,
	) {
		// Extend the fuzzed input to 32 byte leaves if not in leaves format.
		if !isLeaves {
			leavesBytes := make([]byte, len(original)*32)
			for i := range 32 {
				copy(
					leavesBytes[i*len(original):(i+1)*len(original)],
					original,
				)
			}
			original = leavesBytes
		}

		// Convert []byte to [][32]byte as required by HashTreeRoot.
		var input [][32]byte
		for i := 0; i < len(original); i += 32 {
			var block [32]byte
			copy(block[:], original[i:min(i+32, len(original))])
			input = append(input, block)
		}

		// Ensure an even number of chunks for HashTreeRoot.
		expectError := false
		if len(input)%2 != 0 {
			expectError = true
		}

		requireGoHashTreeEquivalence(
			t, input, minParallelizationSize, expectError,
		)
	})
}
