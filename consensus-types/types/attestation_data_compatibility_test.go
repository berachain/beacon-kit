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
	"github.com/berachain/beacon-kit/primitives/math"
	karalabe "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// AttestationDataSize is the size of the AttestationData object in bytes.
// 8 bytes for Slot + 8 bytes for Index + 32 bytes for BeaconBlockRoot.
const AttestationDataSize = 48

// Compile-time assertions to ensure AttestationDataKaralabe implements necessary interfaces.
var _ karalabe.StaticObject = (*AttestationDataKaralabe)(nil)

// AttestationDataKaralabe represents an attestation data.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
type AttestationDataKaralabe struct {
	// Slot is the slot number of the attestation data.
	Slot math.U64 `json:"slot"`
	// Index is the index of the validator.
	Index math.U64 `json:"index"`
	// BeaconBlockRoot is the root of the beacon block.
	BeaconBlockRoot common.Root `json:"beaconBlockRoot"`
}

// SizeSSZ returns the size of the AttestationData object in SSZ encoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (*AttestationDataKaralabe) SizeSSZ() uint32 {
	return AttestationDataSize
}

// DefineSSZ defines the SSZ encoding for the AttestationData object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (a *AttestationDataKaralabe) DefineSSZ(codec *karalabe.Codec) {
	karalabe.DefineUint64(codec, &a.Slot)
	karalabe.DefineUint64(codec, &a.Index)
	karalabe.DefineStaticBytes(codec, &a.BeaconBlockRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the AttestationData object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (a *AttestationDataKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashSequential(a)
}

// MarshalSSZ marshals the AttestationData object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (a *AttestationDataKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabe.Size(a))
	return buf, karalabe.EncodeToBytes(buf, a)
}

func (*AttestationDataKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// UnmarshalSSZ unmarshals the AttestationData object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use karalabe.DecodeFromBytes
func (a *AttestationDataKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, a)
}

// TestAttestationDataCompatibility tests that the current AttestationData implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestAttestationDataCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.AttestationData, *AttestationDataKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.AttestationData, *AttestationDataKaralabe) {
				return &types.AttestationData{}, &AttestationDataKaralabe{}
			},
		},
		{
			name: "typical attestation data",
			setup: func() (*types.AttestationData, *AttestationDataKaralabe) {
				root := common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}

				current := &types.AttestationData{
					Slot:            math.U64(12345),
					Index:           math.U64(67890),
					BeaconBlockRoot: root,
				}
				karalabe := &AttestationDataKaralabe{
					Slot:            math.U64(12345),
					Index:           math.U64(67890),
					BeaconBlockRoot: root,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.AttestationData, *AttestationDataKaralabe) {
				var root common.Root
				for i := range root {
					root[i] = 0xFF
				}

				current := &types.AttestationData{
					Slot:            math.U64(^uint64(0)),
					Index:           math.U64(^uint64(0)),
					BeaconBlockRoot: root,
				}
				karalabe := &AttestationDataKaralabe{
					Slot:            math.U64(^uint64(0)),
					Index:           math.U64(^uint64(0)),
					BeaconBlockRoot: root,
				}
				return current, karalabe
			},
		},
		{
			name: "specific slot and index",
			setup: func() (*types.AttestationData, *AttestationDataKaralabe) {
				root := common.Root{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77,
					0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
				}

				current := &types.AttestationData{
					Slot:            math.U64(1000),
					Index:           math.U64(5),
					BeaconBlockRoot: root,
				}
				karalabe := &AttestationDataKaralabe{
					Slot:            math.U64(1000),
					Index:           math.U64(5),
					BeaconBlockRoot: root,
				}
				return current, karalabe
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
			require.Equal(t, int(AttestationDataSize), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(AttestationDataSize), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.AttestationData{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &AttestationDataKaralabe{}
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

// TestAttestationDataCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestAttestationDataCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid attestation data
		var root common.Root

		// Use deterministic "random" data based on iteration
		for j := range root {
			root[j] = byte((i + j) % 256)
		}

		slot := math.U64(uint64(i) * 12345)
		index := math.U64(uint64(i) * 67)

		current := &types.AttestationData{
			Slot:            slot,
			Index:           index,
			BeaconBlockRoot: root,
		}
		karalabe := &AttestationDataKaralabe{
			Slot:            slot,
			Index:           index,
			BeaconBlockRoot: root,
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

// TestAttestationDataCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestAttestationDataCompatibilityInvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "insufficient data",
			data: make([]byte, 30), // less than required 48 bytes
		},
		{
			name: "excess data",
			data: make([]byte, 100), // more than required 48 bytes
		},
		{
			name: "exact size but invalid content",
			data: make([]byte, 48), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.AttestationData{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test unmarshal with karalabe implementation
			karalabe := &AttestationDataKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.Equal(t, current.Slot, karalabe.Slot, "slots should match")
				require.Equal(t, current.Index, karalabe.Index, "indices should match")
				require.Equal(t, current.BeaconBlockRoot, karalabe.BeaconBlockRoot, "roots should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestAttestationDataCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestAttestationDataCompatibilityRoundTrip(t *testing.T) {
	// Create attestation data with specific values
	original := &types.AttestationData{
		Slot:  math.U64(999999),
		Index: math.U64(12345),
		BeaconBlockRoot: common.Root{
			0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
			0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10,
			0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
			0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
		},
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &AttestationDataKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.AttestationData{}
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
