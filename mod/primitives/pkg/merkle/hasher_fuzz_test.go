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
	"runtime"
	"testing"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/constants"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

func FuzzHashTreeRoot(f *testing.F) {
	// Seed corpus with a variety of sizes, including edge cases
	//
	// Test with empty slice
	f.Add(0, 1, merkle.MinParallelizationSize)
	// Just below a single block size
	f.Add(31, runtime.GOMAXPROCS(0)-1, merkle.MinParallelizationSize)
	// Exactly one block size
	f.Add(32, runtime.GOMAXPROCS(0)+1, merkle.MinParallelizationSize)
	// Just above a single block size
	f.Add(33, runtime.GOMAXPROCS(0)*2, merkle.MinParallelizationSize)
	// Multiple blocks
	f.Add(64, runtime.GOMAXPROCS(0)*4, merkle.MinParallelizationSize)
	// Larger input
	f.Add(1024, 3, merkle.MinParallelizationSize)
	// Just below MinParallelizationSize
	f.Add(merkle.MinParallelizationSize-2, 300, merkle.MinParallelizationSize)
	// Exactly MinParallelizationSize
	f.Add(merkle.MinParallelizationSize, 1, merkle.MinParallelizationSize)
	// Just above MinParallelizationSize
	f.Add(merkle.MinParallelizationSize+2, 64, merkle.MinParallelizationSize)
	// Double MinParallelizationSize
	f.Add(
		2*merkle.MinParallelizationSize,
		runtime.GOMAXPROCS(0)-1,
		merkle.MinParallelizationSize,
	)
	// Really large inputs
	f.Add(
		// Max Txs leaves
		int(constants.MaxTxsPerPayload),
		runtime.GOMAXPROCS(0)-1,
		merkle.MinParallelizationSize,
	)
	// f.Add(
	// 	// NOTE: Testing Max Bytes Per Tx leaves times out
	// 	int(constants.MaxBytesPerTx),
	// 	runtime.GOMAXPROCS(0)-1,
	//  merkle.MinParallelizationSize,
	// )

	f.Fuzz(func(
		t *testing.T,
		leaveSize,
		numRoutines,
		minParallelizationSize int,
	) {
		original := make([]byte, 32*leaveSize)

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

		requireGoHashTreeEquivalence(
			t, input, numRoutines, minParallelizationSize, expectError,
		)
	})
}
