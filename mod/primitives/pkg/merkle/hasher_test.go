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

	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/merkle/zero"
)

// Test NewRootWithMaxLeaves with empty leaves.
func TestNewRootWithMaxLeaves_EmptyLeaves(t *testing.T) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer)

	root, err := hasher.NewRootWithMaxLeaves(nil, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRoot := zero.Hashes[0]
	if root != expectedRoot {
		t.Errorf("Expected root to be %v, got %v", expectedRoot, root)
	}
}

// Test NewRootWithDepth with empty leaves.
func TestNewRootWithDepth_EmptyLeaves(t *testing.T) {
	buffer := getBuffer("reusable")
	hasher := merkle.NewHasher(buffer)

	root, err := hasher.NewRootWithDepth(nil, 0)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expectedRoot := zero.Hashes[0]
	if root != expectedRoot {
		t.Errorf("Expected root to be %v, got %v", expectedRoot, root)
	}
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
	hasher := merkle.NewHasher(buffer)

	leaf := createDummyLeaf(1)
	leaves := [][32]byte{leaf}

	root, err := hasher.NewRootWithMaxLeaves(leaves, 1)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if root != leaf {
		t.Errorf("Expected root to be %v, got %v", leaf, root)
	}
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
	hasher := merkle.NewHasher(buffer)

	leaves := make([][32]byte, 1000)
	for i := range 1000 {
		leaves[i] = createDummyLeaf(byte(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.NewRootWithMaxLeaves(leaves, uint64(len(leaves)))
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
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
	hasher := merkle.NewHasher(buffer)

	leaves := make([][32]byte, 1000)
	for i := range 1000 {
		leaves[i] = createDummyLeaf(byte(i))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := hasher.NewRootWithMaxLeaves(leaves, uint64(len(leaves)))
		if err != nil {
			b.Fatalf("Expected no error, got %v", err)
		}
	}
}
