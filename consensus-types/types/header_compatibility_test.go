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
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// BeaconBlockHeaderSizeKaralabe is the size of the BeaconBlockHeader object in bytes.
//
// Total size: Slot (8) + ProposerIndex (8) +
// ParentBlockRoot (32) + StateRoot (32) + BodyRoot (32).
const BeaconBlockHeaderSizeKaralabe = 112

// Compile-time assertions to ensure BeaconBlockHeaderKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*BeaconBlockHeaderKaralabe)(nil)

// BeaconBlockHeaderKaralabe is an exact copy of BeaconBlockHeader from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// BeaconBlockHeader represents the base of a beacon block header.
type BeaconBlockHeaderKaralabe struct {
	// Slot represents the position of the block in the chain.
	Slot math.Slot `json:"slot"`
	// ProposerIndex is the index of the validator who proposed the block.
	ProposerIndex math.ValidatorIndex `json:"proposer_index"`
	// ParentBlockRoot is the hash of the parent block
	ParentBlockRoot common.Root `json:"parent_block_root"`
	// StateRoot is the hash of the state at the block.
	StateRoot common.Root `json:"state_root"`
	// BodyRoot is the root of the block body.
	BodyRoot common.Root `json:"body_root"`
}

// SizeSSZ returns the size of the BeaconBlockHeader object in SSZ encoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *BeaconBlockHeaderKaralabe) SizeSSZ() uint32 {
	return BeaconBlockHeaderSizeKaralabe
}

// DefineSSZ defines the SSZ encoding for the BeaconBlockHeader object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *BeaconBlockHeaderKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &b.Slot)
	ssz.DefineUint64(codec, &b.ProposerIndex)
	ssz.DefineStaticBytes(codec, &b.ParentBlockRoot)
	ssz.DefineStaticBytes(codec, &b.StateRoot)
	ssz.DefineStaticBytes(codec, &b.BodyRoot)
}

// MarshalSSZ marshals the BeaconBlockBody object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *BeaconBlockHeaderKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

func (*BeaconBlockHeaderKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot computes the SSZ hash tree root of the BeaconBlockHeader object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *BeaconBlockHeaderKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

// UnmarshalSSZ unmarshals the BeaconBlockHeader object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (b *BeaconBlockHeaderKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// TestBeaconBlockHeaderCompatibility tests that the current BeaconBlockHeader implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestBeaconBlockHeaderCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.BeaconBlockHeader, *BeaconBlockHeaderKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.BeaconBlockHeader, *BeaconBlockHeaderKaralabe) {
				return &types.BeaconBlockHeader{}, &BeaconBlockHeaderKaralabe{}
			},
		},
		{
			name: "typical header",
			setup: func() (*types.BeaconBlockHeader, *BeaconBlockHeaderKaralabe) {
				slot := math.Slot(12345)
				proposerIndex := math.ValidatorIndex(567)
				parentRoot := common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
				stateRoot := common.Root{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
				bodyRoot := common.Root{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47}

				current := &types.BeaconBlockHeader{
					Slot:            slot,
					ProposerIndex:   proposerIndex,
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}
				karalabe := &BeaconBlockHeaderKaralabe{
					Slot:            slot,
					ProposerIndex:   proposerIndex,
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.BeaconBlockHeader, *BeaconBlockHeaderKaralabe) {
				slot := math.Slot(^uint64(0))
				proposerIndex := math.ValidatorIndex(^uint64(0))
				var parentRoot, stateRoot, bodyRoot common.Root
				for i := range parentRoot {
					parentRoot[i] = 0xFF
					stateRoot[i] = 0xFF
					bodyRoot[i] = 0xFF
				}

				current := &types.BeaconBlockHeader{
					Slot:            slot,
					ProposerIndex:   proposerIndex,
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}
				karalabe := &BeaconBlockHeaderKaralabe{
					Slot:            slot,
					ProposerIndex:   proposerIndex,
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}
				return current, karalabe
			},
		},
		{
			name: "genesis header",
			setup: func() (*types.BeaconBlockHeader, *BeaconBlockHeaderKaralabe) {
				// Genesis typically has slot 0 and special root values
				slot := math.Slot(0)
				proposerIndex := math.ValidatorIndex(0)
				parentRoot := common.Root{} // Zero root for genesis
				stateRoot := common.Root{0x01} // Some genesis state root
				bodyRoot := common.Root{0x02}  // Some genesis body root

				current := &types.BeaconBlockHeader{
					Slot:            slot,
					ProposerIndex:   proposerIndex,
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}
				karalabe := &BeaconBlockHeaderKaralabe{
					Slot:            slot,
					ProposerIndex:   proposerIndex,
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
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
			require.Equal(t, int(BeaconBlockHeaderSizeKaralabe), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(BeaconBlockHeaderSizeKaralabe), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.BeaconBlockHeader{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &BeaconBlockHeaderKaralabe{}
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

// TestBeaconBlockHeaderCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestBeaconBlockHeaderCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid header data
		slot := math.Slot(uint64(i) * 12345)
		proposerIndex := math.ValidatorIndex(uint64(i) % 1000)
		
		var parentRoot, stateRoot, bodyRoot common.Root
		// Use deterministic "random" data based on iteration
		for j := range parentRoot {
			parentRoot[j] = byte((i + j) % 256)
			stateRoot[j] = byte((i * 2 + j) % 256)
			bodyRoot[j] = byte((i * 3 + j) % 256)
		}

		current := &types.BeaconBlockHeader{
			Slot:            slot,
			ProposerIndex:   proposerIndex,
			ParentBlockRoot: parentRoot,
			StateRoot:       stateRoot,
			BodyRoot:        bodyRoot,
		}
		karalabe := &BeaconBlockHeaderKaralabe{
			Slot:            slot,
			ProposerIndex:   proposerIndex,
			ParentBlockRoot: parentRoot,
			StateRoot:       stateRoot,
			BodyRoot:        bodyRoot,
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

// TestBeaconBlockHeaderCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestBeaconBlockHeaderCompatibilityInvalidData(t *testing.T) {
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
			data: make([]byte, 50), // less than required 112 bytes
		},
		{
			name: "excess data",
			data: make([]byte, 200), // more than required 112 bytes
		},
		{
			name: "exact size but all zeros",
			data: make([]byte, 112), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.BeaconBlockHeader{}
			currentErr := current.UnmarshalSSZ(tc.data)
			
			// Test unmarshal with karalabe implementation
			karalabe := &BeaconBlockHeaderKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)
			
			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.Equal(t, current.Slot, karalabe.Slot, "slots should match")
				require.Equal(t, current.ProposerIndex, karalabe.ProposerIndex, "proposer indices should match")
				require.Equal(t, current.ParentBlockRoot, karalabe.ParentBlockRoot, "parent block roots should match")
				require.Equal(t, current.StateRoot, karalabe.StateRoot, "state roots should match")
				require.Equal(t, current.BodyRoot, karalabe.BodyRoot, "body roots should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestBeaconBlockHeaderCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestBeaconBlockHeaderCompatibilityRoundTrip(t *testing.T) {
	// Create a header with specific values
	original := &types.BeaconBlockHeader{
		Slot:            math.Slot(98765),
		ProposerIndex:   math.ValidatorIndex(42),
		ParentBlockRoot: common.Root{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210, 220, 230, 240, 250, 1, 2, 3, 4, 5, 6, 7},
		StateRoot:       common.Root{7, 6, 5, 4, 3, 2, 1, 250, 240, 230, 220, 210, 200, 190, 180, 170, 160, 150, 140, 130, 120, 110, 100, 90, 80, 70, 60, 50, 40, 30, 20, 10},
		BodyRoot:        common.Root{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131},
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &BeaconBlockHeaderKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.BeaconBlockHeader{}
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