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

package buffer_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes/buffer"
)

type bufferI interface {
	Get(size int) [][32]byte
}

// getBuffer returns a buffer of the given type.
func getBuffer(usageType string) bufferI {
	switch usageType {
	case "reusable":
		return buffer.NewReusableBuffer[[32]byte]()
	case "singleuse":
		return buffer.NewSingleuseBuffer[[32]byte]()
	default:
		panic("unknown usage type: " + usageType)
	}
}

// Test getting a slice of the internal re-usable buffer and modifying.
func TestReusableGet(t *testing.T) {
	buffer := getBuffer("reusable")

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
			if result[0][i-1] != byte(i-1) {
				t.Errorf(
					"Expected result[0][%d] to be %b, got %d",
					i-1, byte(i-1), result[0][i-1],
				)
			}

			result[0] = [32]byte{}
			result[0][i] = byte(i)
		}
	}
}

// Test getting a slice of the internal single-use buffer and modifying.
func TestSingleuseGet(t *testing.T) {
	buffer := getBuffer("singleuse")

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

		// Ensure modifications to the underlying buffer do not persist.
		if i >= 1 {
			if result[0][i-1] != byte(0) {
				t.Errorf(
					"Expected result[0][%d] to be %b, got %d",
					i-1, byte(0), result[0][i-1],
				)
			}

			result[0] = [32]byte{}
			result[0][i] = byte(i)
		}
	}
}

// Benchmark for the Get method on the re-usable buffer
//
// goos: darwin
// goarch: arm64
// pkg: github.com/berachain/beacon-kit/mod/primitives/pkg/merkle
// BenchmarkReusableGet-12  158002471  7.439 ns/op  0 B/op  0 allocs/op.
func BenchmarkReusableGet(b *testing.B) {
	buffer := getBuffer("reusable")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for range b.N {
		size := r.Intn(100) + 1
		result := buffer.Get(size)

		// Perform some operation on the result to avoid compiler optimizations.
		result[0] = [32]byte{}
		index := r.Intn(32)
		result[0][index] = byte(index)
	}
}

// Benchmark for the Get method on the single-use buffer
//
// goos: darwin
// goarch: arm64
// pkg: github.com/berachain/beacon-kit/mod/primitives/pkg/merkle
// BenchmarkSingleuseGet-12  5388972  215.0 ns/op  1700 B/op  1 allocs/op.
func BenchmarkSingleuseGet(b *testing.B) {
	buffer := getBuffer("singleuse")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b.ResetTimer()
	for range b.N {
		size := r.Intn(100) + 1
		result := buffer.Get(size)

		// Perform some operation on the result to avoid compiler optimizations.
		result[0] = [32]byte{}
		index := r.Intn(32)
		result[0][index] = byte(index)
	}
}
