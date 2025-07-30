// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	fastssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func TestSyncAggregate_FastSSZ(t *testing.T) {
	t.Parallel()

	// Create a test SyncAggregate with zero values (as it must be unused)
	sa := &types.SyncAggregate{
		SyncCommitteeBits:      [64]byte{},
		SyncCommitteeSignature: crypto.BLSSignature{},
	}

	t.Run("MarshalSSZTo", func(t *testing.T) {
		dst := make([]byte, 0)
		result, err := sa.MarshalSSZTo(dst)
		require.NoError(t, err)
		require.NotNil(t, result)
		require.Equal(t, 160, len(result)) // 64 + 96 bytes
	})

	t.Run("UnmarshalSSZ", func(t *testing.T) {
		// Create a valid buffer with zero values
		buf := make([]byte, 160)
		newSA := &types.SyncAggregate{}
		err := newSA.UnmarshalSSZ(buf)
		require.NoError(t, err)
		require.Equal(t, sa, newSA)
	})

	t.Run("UnmarshalSSZ_InvalidSize", func(t *testing.T) {
		// Test with invalid buffer size
		buf := make([]byte, 100) // Wrong size
		newSA := &types.SyncAggregate{}
		err := newSA.UnmarshalSSZ(buf)
		require.Error(t, err)
		require.Contains(t, err.Error(), "incorrect size")
	})

	t.Run("UnmarshalSSZ_NonZeroData", func(t *testing.T) {
		// Test with non-zero data (should fail EnforceUnused)
		buf := make([]byte, 160)
		buf[0] = 1 // Set a non-zero bit
		newSA := &types.SyncAggregate{}
		err := newSA.UnmarshalSSZ(buf)
		require.Error(t, err)
		require.Contains(t, err.Error(), "SyncAggregate must be unused")
	})

	t.Run("SizeSSZ", func(t *testing.T) {
		size := sa.SizeSSZ()
		require.Equal(t, 160, size)
	})

	t.Run("HashTreeRootWith", func(t *testing.T) {
		hh := fastssz.NewHasher()
		err := sa.HashTreeRootWith(hh)
		require.NoError(t, err)
	})

	t.Run("GetTree", func(t *testing.T) {
		tree, err := sa.GetTree()
		require.NoError(t, err)
		require.NotNil(t, tree)
	})

	t.Run("CompareHashTreeRoot", func(t *testing.T) {
		// Compare karalabe/ssz HashTreeRoot with fastssz HashTreeRootWith
		karalabRoot, err := sa.HashTreeRoot()
		require.NoError(t, err)

		hh := fastssz.NewHasher()
		err = sa.HashTreeRootWith(hh)
		require.NoError(t, err)
		fastsszRoot, err := hh.HashRoot()
		require.NoError(t, err)

		require.Equal(t, karalabRoot[:], fastsszRoot[:],
			"HashTreeRoot results should match between karalabe/ssz and fastssz")
	})
}

func TestSyncAggregate_EnforceUnused(t *testing.T) {
	t.Parallel()

	t.Run("ZeroValues", func(t *testing.T) {
		sa := &types.SyncAggregate{}
		err := sa.EnforceUnused()
		require.NoError(t, err)
	})

	t.Run("NonZeroBits", func(t *testing.T) {
		sa := &types.SyncAggregate{
			SyncCommitteeBits: [64]byte{1}, // Non-zero bit
		}
		err := sa.EnforceUnused()
		require.Error(t, err)
		require.Contains(t, err.Error(), "SyncAggregate must be unused")
	})

	t.Run("NonZeroSignature", func(t *testing.T) {
		sa := &types.SyncAggregate{
			SyncCommitteeSignature: crypto.BLSSignature{1}, // Non-zero signature
		}
		err := sa.EnforceUnused()
		require.Error(t, err)
		require.Contains(t, err.Error(), "SyncAggregate must be unused")
	})
}
