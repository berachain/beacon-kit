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

package types

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// TestSyncAggregateSSZCompatibility tests that SyncAggregate's
// fastssz and karalabe/ssz implementations produce identical results.
func TestSyncAggregateSSZCompatibility(t *testing.T) {
	// Test case 1: Empty SyncAggregate
	t.Run("EmptySyncAggregate", func(t *testing.T) {
		sa := &SyncAggregate{}

		// Test marshaling
		karalabeBytes, err := marshalSyncAggregateKaralabe(sa)
		require.NoError(t, err)

		fastSSZBytes, err := sa.MarshalSSZ()
		require.NoError(t, err)

		require.Equal(t, karalabeBytes, fastSSZBytes,
			"karalabe/ssz and fastssz should produce identical bytes for empty SyncAggregate")

		// Test hash tree root
		karalabeHTR := hashTreeRootSyncAggregateKaralabe(sa)
		fastSSZHTR := sa.HashTreeRoot()

		require.Equal(t, [32]byte(karalabeHTR), [32]byte(fastSSZHTR),
			"karalabe/ssz and fastssz should produce identical hash tree roots for empty SyncAggregate")

		// Test size  
		require.Equal(t, 160, sa.SizeSSZFastSSZ(),
			"SyncAggregate should have size 160 bytes")
	})

	// Test case 2: SyncAggregate with data
	t.Run("SyncAggregateWithData", func(t *testing.T) {
		sa := &SyncAggregate{
			SyncCommitteeBits: [64]byte{
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
				0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
				0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x11, 0x22,
				0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
			},
			SyncCommitteeSignature: crypto.BLSSignature{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
				0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10,
				0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
				0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20,
				0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
				0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30,
				0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38,
				0x39, 0x3A, 0x3B, 0x3C, 0x3D, 0x3E, 0x3F, 0x40,
				0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48,
				0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50,
				0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57, 0x58,
				0x59, 0x5A, 0x5B, 0x5C, 0x5D, 0x5E, 0x5F, 0x60,
			},
		}

		// Test marshaling
		karalabeBytes, err := marshalSyncAggregateKaralabe(sa)
		require.NoError(t, err)

		fastSSZBytes, err := sa.MarshalSSZ()
		require.NoError(t, err)

		require.Equal(t, karalabeBytes, fastSSZBytes,
			"karalabe/ssz and fastssz should produce identical bytes for populated SyncAggregate")

		// Test hash tree root
		karalabeHTR := hashTreeRootSyncAggregateKaralabe(sa)
		fastSSZHTR := sa.HashTreeRoot()

		require.Equal(t, [32]byte(karalabeHTR), [32]byte(fastSSZHTR),
			"karalabe/ssz and fastssz should produce identical hash tree roots for populated SyncAggregate")
	})

	// Test case 3: Unmarshal with empty data (to pass EnforceUnused)
	t.Run("UnmarshalEmpty", func(t *testing.T) {
		// Create empty SyncAggregate
		originalSA := &SyncAggregate{}

		// Marshal with fastssz
		data, err := originalSA.MarshalSSZ()
		require.NoError(t, err)
		require.Len(t, data, 160, "Marshaled data should be 160 bytes")

		// Unmarshal with fastssz
		unmarshaledSA := &SyncAggregate{}
		err = unmarshaledSA.UnmarshalSSZ(data)
		require.NoError(t, err)

		// Compare
		require.Equal(t, originalSA, unmarshaledSA,
			"Unmarshaled empty SyncAggregate should equal original")
	})

	// Test case 4: EnforceUnused
	t.Run("EnforceUnused", func(t *testing.T) {
		// Empty should pass
		emptySA := &SyncAggregate{}
		require.NoError(t, emptySA.EnforceUnused(),
			"Empty SyncAggregate should pass EnforceUnused")

		// Non-empty should fail
		nonEmptySA := &SyncAggregate{
			SyncCommitteeBits: [64]byte{0x01}, // Just one bit set
		}
		require.Error(t, nonEmptySA.EnforceUnused(),
			"Non-empty SyncAggregate should fail EnforceUnused")
	})

	// Test case 5: Error cases
	t.Run("ErrorCases", func(t *testing.T) {
		sa := &SyncAggregate{}

		// Test unmarshal with wrong size
		err := sa.UnmarshalSSZ(make([]byte, 159)) // One byte short
		require.Error(t, err, "Should error on incorrect buffer size")

		err = sa.UnmarshalSSZ(make([]byte, 161)) // One byte too many
		require.Error(t, err, "Should error on incorrect buffer size")
	})
}

// Helper function to marshal SyncAggregate using karalabe/ssz
func marshalSyncAggregateKaralabe(sa *SyncAggregate) ([]byte, error) {
	buf := make([]byte, ssz.Size(sa))
	return buf, ssz.EncodeToBytes(buf, sa)
}

// Helper function to get hash tree root using karalabe/ssz
func hashTreeRootSyncAggregateKaralabe(sa *SyncAggregate) [32]byte {
	return ssz.HashSequential(sa)
}