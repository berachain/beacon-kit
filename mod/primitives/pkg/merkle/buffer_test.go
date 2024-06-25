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
	"math/rand"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
)

func TestGet(t *testing.T) {
	buffer := merkle.NewBuffer[[32]byte]()

	testCases := []struct {
		size     int
		expected int
	}{
		{size: 0, expected: 0},
		{size: 1, expected: 1},
		{size: 5, expected: 5},
		{size: 16, expected: 16},
		{size: 33, expected: 33},
		{size: 17, expected: 17},
		{size: 100, expected: 100},
	}

	for i, tc := range testCases {
		result := buffer.Get(tc.size)

		if len(result) != tc.expected {
			t.Errorf(
				"Expected result size to be %d, got %d",
				tc.expected, len(result),
			)
		}

		// Ensure modifications to the underlying buffer persist.
		if i >= 1 {
			// Set the value in the previous iteration
			result[0][i-1] = byte(i - 1)

			// Check if the value persists in the current iteration
			newResult := buffer.Get(tc.size)
			if newResult[0][i-1] != byte(i-1) {
				t.Errorf(
					"Expected newResult[0][%d] to be %d, got %d",
					i-1, i-1, newResult[0][i-1],
				)
			}
		}
	}
}

// Benchmark for the Get method
//
// goos: darwin
// goarch: arm64
// pkg: github.com/berachain/beacon-kit/mod/primitives/pkg/merkle
// BenchmarkGet-12     173148679     6.917 ns/op     0 B/op     0 allocs/op
func BenchmarkGet(b *testing.B) {
	buffer := merkle.NewBuffer[[32]byte]()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for range b.N {
		size := r.Intn(100) + 1
		result := buffer.Get(size)

		// Peform some operation on the result to avoid compiler optimizations.
		result[0] = [32]byte{}
		index := r.Intn(32)
		result[0][index] = byte(index)
	}
}
