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
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// Compile-time assertions to ensure SignedBeaconBlockHeaderKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*SignedBeaconBlockHeaderKaralabe)(nil)

// SignedBeaconBlockHeaderKaralabe is an exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// SignedBeaconBlockHeader is a struct that contains a BeaconBlockHeader and a BLSSignature.
//
// NOTE: This struct is only ever (un)marshalled with SSZ and NOT with JSON.
type SignedBeaconBlockHeaderKaralabe struct {
	Header    *BeaconBlockHeaderKaralabe
	Signature crypto.BLSSignature
}

// SizeSSZ returns the size of the SignedBeaconBlockHeader object
// in SSZ encoding. Total size: Header (112) + Signature (96).
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *SignedBeaconBlockHeaderKaralabe) SizeSSZ() uint32 {
	//nolint:mnd // no magic
	size := (*BeaconBlockHeaderKaralabe)(nil).SizeSSZ() + 96
	return size
}

// DefineSSZ defines the SSZ encoding for the SignedBeaconBlockHeader object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *SignedBeaconBlockHeaderKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticObject(codec, &b.Header)
	ssz.DefineStaticBytes(codec, &b.Signature)
}

// MarshalSSZ marshals the SignedBeaconBlockHeader object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *SignedBeaconBlockHeaderKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(b))
	return buf, ssz.EncodeToBytes(buf, b)
}

// HashTreeRoot computes the SSZ hash tree root of the
// SignedBeaconBlockHeader object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (b *SignedBeaconBlockHeaderKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(b)
}

func (*SignedBeaconBlockHeaderKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// UnmarshalSSZ unmarshals the SignedBeaconBlockHeader object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (b *SignedBeaconBlockHeaderKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, b)
}

// TestSignedBeaconBlockHeaderCompatibility tests that the current SignedBeaconBlockHeader implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestSignedBeaconBlockHeaderCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.SignedBeaconBlockHeader, *SignedBeaconBlockHeaderKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.SignedBeaconBlockHeader, *SignedBeaconBlockHeaderKaralabe) {
				return &types.SignedBeaconBlockHeader{
						Header: &types.BeaconBlockHeader{},
					}, &SignedBeaconBlockHeaderKaralabe{
						Header: &BeaconBlockHeaderKaralabe{},
					}
			},
		},
		{
			name: "typical signed header",
			setup: func() (*types.SignedBeaconBlockHeader, *SignedBeaconBlockHeaderKaralabe) {
				header := &types.BeaconBlockHeader{
					Slot:            math.Slot(12345),
					ProposerIndex:   math.ValidatorIndex(567),
					ParentBlockRoot: common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					StateRoot:       common.Root{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
					BodyRoot:        common.Root{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
				}
				headerKaralabe := &BeaconBlockHeaderKaralabe{
					Slot:            math.Slot(12345),
					ProposerIndex:   math.ValidatorIndex(567),
					ParentBlockRoot: common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32},
					StateRoot:       common.Root{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
					BodyRoot:        common.Root{16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
				}

				var sig crypto.BLSSignature
				for i := range sig {
					sig[i] = byte(i % 256)
				}

				current := &types.SignedBeaconBlockHeader{
					Header:    header,
					Signature: sig,
				}
				karalabe := &SignedBeaconBlockHeaderKaralabe{
					Header:    headerKaralabe,
					Signature: sig,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.SignedBeaconBlockHeader, *SignedBeaconBlockHeaderKaralabe) {
				var parentRoot, stateRoot, bodyRoot common.Root
				var sig crypto.BLSSignature

				// Fill with max values
				for i := range parentRoot {
					parentRoot[i] = 0xFF
					stateRoot[i] = 0xFF
					bodyRoot[i] = 0xFF
				}
				for i := range sig {
					sig[i] = 0xFF
				}

				header := &types.BeaconBlockHeader{
					Slot:            math.Slot(^uint64(0)),
					ProposerIndex:   math.ValidatorIndex(^uint64(0)),
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}
				headerKaralabe := &BeaconBlockHeaderKaralabe{
					Slot:            math.Slot(^uint64(0)),
					ProposerIndex:   math.ValidatorIndex(^uint64(0)),
					ParentBlockRoot: parentRoot,
					StateRoot:       stateRoot,
					BodyRoot:        bodyRoot,
				}

				current := &types.SignedBeaconBlockHeader{
					Header:    header,
					Signature: sig,
				}
				karalabe := &SignedBeaconBlockHeaderKaralabe{
					Header:    headerKaralabe,
					Signature: sig,
				}
				return current, karalabe
			},
		},
		{
			name: "genesis signed header",
			setup: func() (*types.SignedBeaconBlockHeader, *SignedBeaconBlockHeaderKaralabe) {
				header := &types.BeaconBlockHeader{
					Slot:            math.Slot(0),
					ProposerIndex:   math.ValidatorIndex(0),
					ParentBlockRoot: common.Root{},
					StateRoot:       common.Root{0x01},
					BodyRoot:        common.Root{0x02},
				}
				headerKaralabe := &BeaconBlockHeaderKaralabe{
					Slot:            math.Slot(0),
					ProposerIndex:   math.ValidatorIndex(0),
					ParentBlockRoot: common.Root{},
					StateRoot:       common.Root{0x01},
					BodyRoot:        common.Root{0x02},
				}

				// Empty signature for genesis
				var sig crypto.BLSSignature

				current := &types.SignedBeaconBlockHeader{
					Header:    header,
					Signature: sig,
				}
				karalabe := &SignedBeaconBlockHeaderKaralabe{
					Header:    headerKaralabe,
					Signature: sig,
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

			// Test Size - SignedBeaconBlockHeader has size 208 (112 header + 96 signature)
			require.Equal(t, 208, current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(208), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.SignedBeaconBlockHeader{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current.Header.Slot, newCurrent.Header.Slot)
			require.Equal(t, current.Header.ProposerIndex, newCurrent.Header.ProposerIndex)
			require.Equal(t, current.Header.ParentBlockRoot, newCurrent.Header.ParentBlockRoot)
			require.Equal(t, current.Header.StateRoot, newCurrent.Header.StateRoot)
			require.Equal(t, current.Header.BodyRoot, newCurrent.Header.BodyRoot)
			require.Equal(t, current.Signature, newCurrent.Signature)

			// Test Unmarshal with current marshaled data
			newKaralabe := &SignedBeaconBlockHeaderKaralabe{}
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, karalabe.Header.Slot, newKaralabe.Header.Slot)
			require.Equal(t, karalabe.Header.ProposerIndex, newKaralabe.Header.ProposerIndex)
			require.Equal(t, karalabe.Header.ParentBlockRoot, newKaralabe.Header.ParentBlockRoot)
			require.Equal(t, karalabe.Header.StateRoot, newKaralabe.Header.StateRoot)
			require.Equal(t, karalabe.Header.BodyRoot, newKaralabe.Header.BodyRoot)
			require.Equal(t, karalabe.Signature, newKaralabe.Signature)

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestSignedBeaconBlockHeaderCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestSignedBeaconBlockHeaderCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid signed header data
		var parentRoot, stateRoot, bodyRoot common.Root
		var sig crypto.BLSSignature

		// Use deterministic "random" data based on iteration
		for j := range parentRoot {
			parentRoot[j] = byte((i + j) % 256)
			stateRoot[j] = byte((i*2 + j) % 256)
			bodyRoot[j] = byte((i*3 + j) % 256)
		}
		for j := range sig {
			sig[j] = byte((i*4 + j) % 256)
		}

		header := &types.BeaconBlockHeader{
			Slot:            math.Slot(uint64(i) * 12345),
			ProposerIndex:   math.ValidatorIndex(uint64(i) % 1000),
			ParentBlockRoot: parentRoot,
			StateRoot:       stateRoot,
			BodyRoot:        bodyRoot,
		}
		headerKaralabe := &BeaconBlockHeaderKaralabe{
			Slot:            math.Slot(uint64(i) * 12345),
			ProposerIndex:   math.ValidatorIndex(uint64(i) % 1000),
			ParentBlockRoot: parentRoot,
			StateRoot:       stateRoot,
			BodyRoot:        bodyRoot,
		}

		current := &types.SignedBeaconBlockHeader{
			Header:    header,
			Signature: sig,
		}
		karalabe := &SignedBeaconBlockHeaderKaralabe{
			Header:    headerKaralabe,
			Signature: sig,
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

// TestSignedBeaconBlockHeaderCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestSignedBeaconBlockHeaderCompatibilityInvalidData(t *testing.T) {
	testCases := []struct {
		name string
		data []byte
	}{
		{
			name: "empty data",
			data: []byte{},
		},
		{
			name: "insufficient data - header only",
			data: make([]byte, 112), // Only header, no signature
		},
		{
			name: "insufficient data - partial signature",
			data: make([]byte, 150), // Header + partial signature
		},
		{
			name: "excess data",
			data: make([]byte, 300), // more than required 208 bytes
		},
		{
			name: "exact size but all zeros",
			data: make([]byte, 208), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.SignedBeaconBlockHeader{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test unmarshal with karalabe implementation
			karalabe := &SignedBeaconBlockHeaderKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				require.NotNil(t, current.Header, "current header should not be nil")
				require.NotNil(t, karalabe.Header, "karalabe header should not be nil")
				require.Equal(t, current.Header.Slot, karalabe.Header.Slot, "slots should match")
				require.Equal(t, current.Header.ProposerIndex, karalabe.Header.ProposerIndex, "proposer indices should match")
				require.Equal(t, current.Header.ParentBlockRoot, karalabe.Header.ParentBlockRoot, "parent block roots should match")
				require.Equal(t, current.Header.StateRoot, karalabe.Header.StateRoot, "state roots should match")
				require.Equal(t, current.Header.BodyRoot, karalabe.Header.BodyRoot, "body roots should match")
				require.Equal(t, current.Signature, karalabe.Signature, "signatures should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestSignedBeaconBlockHeaderCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestSignedBeaconBlockHeaderCompatibilityRoundTrip(t *testing.T) {
	// Create a signed header with specific values
	header := &types.BeaconBlockHeader{
		Slot:            math.Slot(98765),
		ProposerIndex:   math.ValidatorIndex(42),
		ParentBlockRoot: common.Root{10, 20, 30, 40, 50, 60, 70, 80, 90, 100, 110, 120, 130, 140, 150, 160, 170, 180, 190, 200, 210, 220, 230, 240, 250, 1, 2, 3, 4, 5, 6, 7},
		StateRoot:       common.Root{7, 6, 5, 4, 3, 2, 1, 250, 240, 230, 220, 210, 200, 190, 180, 170, 160, 150, 140, 130, 120, 110, 100, 90, 80, 70, 60, 50, 40, 30, 20, 10},
		BodyRoot:        common.Root{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119, 120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131},
	}

	var sig crypto.BLSSignature
	for i := range sig {
		sig[i] = byte(255 - i)
	}

	original := &types.SignedBeaconBlockHeader{
		Header:    header,
		Signature: sig,
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &SignedBeaconBlockHeaderKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.SignedBeaconBlockHeader{}
	err = roundTrip.UnmarshalSSZ(karalableBytes)
	require.NoError(t, err)

	// Verify round trip preserved all data
	require.Equal(t, original.Header.Slot, roundTrip.Header.Slot)
	require.Equal(t, original.Header.ProposerIndex, roundTrip.Header.ProposerIndex)
	require.Equal(t, original.Header.ParentBlockRoot, roundTrip.Header.ParentBlockRoot)
	require.Equal(t, original.Header.StateRoot, roundTrip.Header.StateRoot)
	require.Equal(t, original.Header.BodyRoot, roundTrip.Header.BodyRoot)
	require.Equal(t, original.Signature, roundTrip.Signature)

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

// TestSignedBeaconBlockHeaderCompatibilityFieldOrdering verifies correct field ordering in SSZ
func TestSignedBeaconBlockHeaderCompatibilityFieldOrdering(t *testing.T) {
	// Create a signed header with easily identifiable patterns
	header := &types.BeaconBlockHeader{
		Slot:          math.Slot(0x0102030405060708),
		ProposerIndex: math.ValidatorIndex(0x090A0B0C0D0E0F10),
		ParentBlockRoot: common.Root{
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
			0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		},
		StateRoot: common.Root{
			0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22,
			0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22,
			0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22,
			0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22, 0x22,
		},
		BodyRoot: common.Root{
			0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
		},
	}

	var sig crypto.BLSSignature
	for i := range sig {
		sig[i] = 0x44 // All signature bytes are 0x44
	}

	current := &types.SignedBeaconBlockHeader{
		Header:    header,
		Signature: sig,
	}

	// Marshal and verify field ordering
	bytes, err := current.MarshalSSZ()
	require.NoError(t, err)
	require.Len(t, bytes, 208, "marshaled size should be 208 bytes")

	// Verify header fields are in correct order and little-endian
	// Slot (8 bytes, little-endian) at offset 0
	require.Equal(t, []byte{0x08, 0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01}, bytes[0:8], "slot should be at offset 0 in little-endian")

	// ProposerIndex (8 bytes, little-endian) at offset 8
	require.Equal(t, []byte{0x10, 0x0F, 0x0E, 0x0D, 0x0C, 0x0B, 0x0A, 0x09}, bytes[8:16], "proposer index should be at offset 8 in little-endian")

	// ParentBlockRoot (32 bytes) at offset 16
	for i := 0; i < 32; i++ {
		require.Equal(t, byte(0x11), bytes[16+i], "parent block root byte %d should be 0x11", i)
	}

	// StateRoot (32 bytes) at offset 48
	for i := 0; i < 32; i++ {
		require.Equal(t, byte(0x22), bytes[48+i], "state root byte %d should be 0x22", i)
	}

	// BodyRoot (32 bytes) at offset 80
	for i := 0; i < 32; i++ {
		require.Equal(t, byte(0x33), bytes[80+i], "body root byte %d should be 0x33", i)
	}

	// Signature (96 bytes) at offset 112
	for i := 0; i < 96; i++ {
		require.Equal(t, byte(0x44), bytes[112+i], "signature byte %d should be 0x44", i)
	}
}
