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

//nolint:testpackage // private functions.
package merkle

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constants"
	"github.com/berachain/beacon-kit/primitives/crypto/sha256"
	"github.com/berachain/beacon-kit/primitives/merkle"
	"github.com/stretchr/testify/require"
)

func Test_create(t *testing.T) {
	hasher := merkle.NewHasher[common.Root](sha256.Hash)

	tests := []struct {
		name   string
		leaves []common.Root
		depth  uint64
		want   TreeNode
	}{
		{
			name:   "empty tree",
			leaves: nil,
			depth:  0,
			want:   &ZeroNode{},
		},
		{
			name:   "zero depth",
			leaves: []common.Root{hexString(t, fmt.Sprintf("%064d", 0))},
			depth:  0,
			want:   &LeafNode{},
		},
		{
			name:   "depth of 1",
			leaves: []common.Root{hexString(t, fmt.Sprintf("%064d", 0))},
			depth:  1,
			want: &InnerNode{
				left:   &LeafNode{},
				right:  &ZeroNode{},
				hasher: hasher,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := create(hasher, tt.leaves, tt.depth); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				require.True(t, tt.want.Equals(got))
			}
		})
	}
}

func Test_fromSnapshotParts(t *testing.T) {
	tests := []struct {
		name      string
		finalized []common.Root
	}{
		{
			name: "multiple deposits and multiple Finalized",
			finalized: []common.Root{
				hexString(t, fmt.Sprintf("%064d", 1)),
				hexString(t, fmt.Sprintf("%064d", 2)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := NewDepositTree()
			for _, leaf := range tt.finalized {
				err := test.pushLeaf(leaf)
				require.NoError(t, err)
			}
			got := test.HashTreeRoot()

			generatedTree, err := merkle.NewTreeFromLeavesWithDepth(
				tt.finalized,
				uint8(constants.DepositContractDepth),
			)
			require.NoError(t, err)
			want := generatedTree.HashTreeRoot()
			require.True(t, want.Equals(common.NewRootFromBytes(got[:])))

			// Test finalization
			for i := range uint64(len(tt.finalized)) {
				err = test.Finalize(i, tt.finalized[i], 0)
				require.NoError(t, err)
			}

			sShot := test.GetSnapshot()
			got = sShot.CalculateRoot()

			require.Len(t, sShot.finalized, 1)
			require.True(t, want.Equals(common.NewRootFromBytes(got[:])))

			// Build from the snapshot once more
			recovered, err := fromSnapshot(sShot)
			require.NoError(t, err)
			got = recovered.HashTreeRoot()
			require.True(t, want.Equals(common.NewRootFromBytes(got[:])))
		})
	}
}

// TODO: add tests against spec for generateProof.

func hexString(t *testing.T, hexStr string) common.Root {
	t.Helper()
	b, err := hex.DecodeString(hexStr)
	require.NoError(t, err)
	if len(b) != 32 {
		require.Len(t, b, 32, "bad hash length, expected 32")
	}
	x := (*common.Root)(b)
	return *x
}
