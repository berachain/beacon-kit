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

//go:build test

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

const (
	syncCommitteeSizeKaralabe       = 512
	syncCommitteeBitsLengthKaralabe = syncCommitteeSizeKaralabe / 8
)

// Compile-time assertions to ensure SyncAggregateKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*SyncAggregateKaralabe)(nil)

// SyncAggregateKaralabe is an exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
type SyncAggregateKaralabe struct {
	SyncCommitteeBits      [syncCommitteeBitsLengthKaralabe]byte
	SyncCommitteeSignature crypto.BLSSignature
}

// SizeSSZ returns the SSZ encoded size in bytes for the SyncAggregate.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SyncAggregateKaralabe) SizeSSZ() uint32 {
	return syncCommitteeBitsLengthKaralabe + 96 // 64 bytes for bits + 96 bytes for BLS signature
}

// DefineSSZ defines the SSZ encoding for the SyncAggregate object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SyncAggregateKaralabe) DefineSSZ(c *ssz.Codec) {
	ssz.DefineStaticBytes(c, &s.SyncCommitteeBits)
	ssz.DefineStaticBytes(c, &s.SyncCommitteeSignature)
}

// MarshalSSZ marshals the SyncAggregate object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SyncAggregateKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(s))
	return buf, ssz.EncodeToBytes(buf, s)
}

func (*SyncAggregateKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot returns the hash tree root of the Deposits.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SyncAggregateKaralabe) HashTreeRoot() common.Root {
	htr := ssz.HashSequential(s)
	return htr
}

// UnmarshalSSZ unmarshals the SyncAggregate object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (s *SyncAggregateKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

// TestSyncAggregateCompatibility tests that the current SyncAggregate implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
// Note: Since SyncAggregate must be unused (all zeros) in the current implementation,
// we can only test with zero values.
func TestSyncAggregateCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.SyncAggregate, *SyncAggregateKaralabe)
	}{
		{
			name: "zero values (required for unused SyncAggregate)",
			setup: func() (*types.SyncAggregate, *SyncAggregateKaralabe) {
				return &types.SyncAggregate{}, &SyncAggregateKaralabe{}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current, karalabe := tc.setup()

			// Test Marshal
			currentBytes, err1 := current.MarshalSSZ()
			require.NoError(t, err1, "current MarshalSSZ should not error")

			karalableBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalableBytes, currentBytes, "marshaled bytes should be identical")

			// Test Size
			require.Equal(t, int(syncCommitteeBitsLengthKaralabe+96), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(syncCommitteeBitsLengthKaralabe+96), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.SyncAggregate{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &SyncAggregateKaralabe{}
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, karalabe, newKaralabe, "unmarshaled karalabe should match original")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestSyncAggregateCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
// Note: This test is commented out because SyncAggregate must be unused (all zeros)
/*
func TestSyncAggregateCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid sync aggregate data
		var bits [64]byte
		var sig crypto.BLSSignature

		// Use deterministic "random" data based on iteration
		for j := range bits {
			bits[j] = byte((i * j) % 256)
		}
		for j := range sig {
			sig[j] = byte((i + j * 2) % 256)
		}

		current := &types.SyncAggregate{
			SyncCommitteeBits:      bits,
			SyncCommitteeSignature: sig,
		}
		karalabe := &SyncAggregateKaralabe{
			SyncCommitteeBits:      bits,
			SyncCommitteeSignature: sig,
		}

		// Compare marshaling
		currentBytes, err1 := current.MarshalSSZ()
		require.NoError(t, err1)
		karalableBytes, err2 := karalabe.MarshalSSZ()
		require.NoError(t, err2)
		require.Equal(t, karalableBytes, currentBytes, "fuzzing iteration %d: marshaled bytes should be identical", i)

		// Compare roots
		currentRoot, err := current.HashTreeRoot()
		require.NoError(t, err)
		require.Equal(t, [32]byte(karalabe.HashTreeRoot()), currentRoot, "fuzzing iteration %d: roots should be identical", i)
	}
}
*/

// TestSyncAggregateCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestSyncAggregateCompatibilityInvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "insufficient data - missing signature",
			data: make([]byte, 64), // Only bits, no signature
		},
		{
			name: "insufficient data - partial signature",
			data: make([]byte, 100), // 64 bits + 36 bytes (partial signature)
		},
		{
			name: "excess data",
			data: make([]byte, 200), // more than required 160 bytes
		},
		{
			name: "exact size but all zeros",
			data: make([]byte, 160), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.SyncAggregate{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test unmarshal with karalabe implementation
			karalabe := &SyncAggregateKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.Equal(t, current.SyncCommitteeBits, karalabe.SyncCommitteeBits, "sync committee bits should match")
				require.Equal(t, current.SyncCommitteeSignature, karalabe.SyncCommitteeSignature, "sync committee signatures should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestSyncAggregateCompatibilityRoundTrip verifies that data can round-trip between implementations
// Note: Since SyncAggregate must be unused (all zeros), we can only test with zero values
func TestSyncAggregateCompatibilityRoundTrip(t *testing.T) {
	// Create a sync aggregate with zero values (required for unused SyncAggregate)
	original := &types.SyncAggregate{}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &SyncAggregateKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.SyncAggregate{}
	err = roundTrip.UnmarshalSSZ(karalableBytes)
	require.NoError(t, err)

	// Verify round trip preserved all data
	require.Equal(t, original, roundTrip, "round trip should preserve all data")

	// Verify both serializations are identical
	require.Equal(t, currentBytes, karalableBytes, "both serializations should be identical")

	// Verify hash roots match throughout
	originalRoot, err := original.HashTreeRoot()
	require.NoError(t, err)
	roundTripRoot, err := roundTrip.HashTreeRoot()
	require.NoError(t, err)
	require.Equal(t, originalRoot, [32]byte(karalabe.HashTreeRoot()), "hash roots should match")
	require.Equal(t, originalRoot, roundTripRoot, "hash roots should match after round trip")
}
