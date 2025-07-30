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

// Compile-time assertions to ensure SlashingInfoKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*SlashingInfoKaralabe)(nil)

// SlashingInfoKaralabe represents a slashing info - exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
type SlashingInfoKaralabe struct {
	// Slot is the slot number of the slashing info.
	Slot math.Slot
	// ValidatorIndex is the validator index of the slashing info.
	Index math.U64
}

// SizeSSZ returns the size of the SlashingInfo object in SSZ encoding.
func (*SlashingInfoKaralabe) SizeSSZ() uint32 {
	return 16 // 8 bytes for Slot + 8 bytes for Index
}

// DefineSSZ defines the SSZ encoding for the SlashingInfo object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SlashingInfoKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &s.Slot)
	ssz.DefineUint64(codec, &s.Index)
}

// HashTreeRoot computes the SSZ hash tree root of the SlashingInfo object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SlashingInfoKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(s)
}

// MarshalSSZ marshals the SlashingInfo object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (s *SlashingInfoKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(s))
	return buf, ssz.EncodeToBytes(buf, s)
}

// UnmarshalSSZ unmarshals the SlashingInfo object from SSZ format.
// Adding this method for completeness
func (s *SlashingInfoKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

// TestSlashingInfoSSZRegression ensures that the SSZ encoding for SlashingInfo
// remains stable and backward compatible.
func TestSlashingInfoSSZRegression(t *testing.T) {
	testCases := []struct {
		name         string
		slashingInfo *types.SlashingInfo
		expectedSSZ  []byte // Pre-computed expected SSZ encoding
	}{
		{
			name: "zero values",
			slashingInfo: &types.SlashingInfo{
				Slot:  math.Slot(0),
				Index: math.U64(0),
			},
			// Expected SSZ: 8 zero bytes (slot) + 8 zero bytes (index) = 16 bytes
			expectedSSZ: make([]byte, 16),
		},
		{
			name: "typical slashing info",
			slashingInfo: &types.SlashingInfo{
				Slot:  math.Slot(12345),
				Index: math.U64(67890),
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 16)
				// Slot (12345 in little-endian)
				ssz[0] = 0x39
				ssz[1] = 0x30
				ssz[2] = 0x00
				ssz[3] = 0x00
				ssz[4] = 0x00
				ssz[5] = 0x00
				ssz[6] = 0x00
				ssz[7] = 0x00
				// Index (67890 in little-endian)
				ssz[8] = 0x32
				ssz[9] = 0x09
				ssz[10] = 0x01
				ssz[11] = 0x00
				ssz[12] = 0x00
				ssz[13] = 0x00
				ssz[14] = 0x00
				ssz[15] = 0x00
				return ssz
			}(),
		},
		{
			name: "maximum values",
			slashingInfo: &types.SlashingInfo{
				Slot:  math.Slot(^uint64(0)),
				Index: math.U64(^uint64(0)),
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 16)
				// All 0xff
				for i := range ssz {
					ssz[i] = 0xff
				}
				return ssz
			}(),
		},
		{
			name: "specific slot and index",
			slashingInfo: &types.SlashingInfo{
				Slot:  math.Slot(1000000),
				Index: math.U64(999),
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 16)
				// Slot (1000000 = 0xF4240 in little-endian)
				ssz[0] = 0x40
				ssz[1] = 0x42
				ssz[2] = 0x0f
				ssz[3] = 0x00
				ssz[4] = 0x00
				ssz[5] = 0x00
				ssz[6] = 0x00
				ssz[7] = 0x00
				// Index (999 = 0x3E7 in little-endian)
				ssz[8] = 0xe7
				ssz[9] = 0x03
				ssz[10] = 0x00
				ssz[11] = 0x00
				ssz[12] = 0x00
				ssz[13] = 0x00
				ssz[14] = 0x00
				ssz[15] = 0x00
				return ssz
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Marshal
			actualSSZ, err := tc.slashingInfo.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ, "SSZ encoding should match expected")

			// Test Size
			require.Equal(t, types.SlashingInfoSize, tc.slashingInfo.SizeSSZ(), "size should be 16 bytes")

			// Test Unmarshal
			unmarshaled := &types.SlashingInfo{}
			err = unmarshaled.UnmarshalSSZ(tc.expectedSSZ)
			require.NoError(t, err, "UnmarshalSSZ should not error")
			require.Equal(t, tc.slashingInfo, unmarshaled, "unmarshaled object should match original")

			// Test MarshalSSZTo
			buf := make([]byte, 0, tc.slashingInfo.SizeSSZ())
			actualSSZ2, err := tc.slashingInfo.MarshalSSZTo(buf)
			require.NoError(t, err, "MarshalSSZTo should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ2, "MarshalSSZTo should produce same output")

			// Test HashTreeRoot consistency
			root1, err := tc.slashingInfo.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot should not error")
			root2, err := unmarshaled.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot of unmarshaled should not error")
			require.Equal(t, root1, root2, "hash tree roots should match")
		})
	}
}

// TestSlashingInfoSSZInvalidData tests error handling for invalid SSZ data
func TestSlashingInfoSSZInvalidData(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		expectedError string
	}{
		{
			name:          "empty data",
			data:          []byte{},
			expectedError: "incorrect size",
		},
		{
			name:          "insufficient data",
			data:          make([]byte, 10), // less than required 16 bytes
			expectedError: "incorrect size",
		},
		{
			name:          "excess data",
			data:          make([]byte, 20), // more than required 16 bytes
			expectedError: "incorrect size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			slashingInfo := &types.SlashingInfo{}
			err := slashingInfo.UnmarshalSSZ(tc.data)
			require.Error(t, err, "UnmarshalSSZ should error on invalid data")
			require.Contains(t, err.Error(), tc.expectedError, "error should contain expected message")
		})
	}
}

// TestSlashingInfoSSZRoundTrip tests round-trip encoding/decoding with various data patterns
func TestSlashingInfoSSZRoundTrip(t *testing.T) {
	patterns := []struct {
		name  string
		setup func() *types.SlashingInfo
	}{
		{
			name: "all zeros",
			setup: func() *types.SlashingInfo {
				return &types.SlashingInfo{}
			},
		},
		{
			name: "incremental pattern",
			setup: func() *types.SlashingInfo {
				return &types.SlashingInfo{
					Slot:  math.Slot(12345678),
					Index: math.U64(87654321),
				}
			},
		},
		{
			name: "specific values",
			setup: func() *types.SlashingInfo {
				return &types.SlashingInfo{
					Slot:  math.Slot(999999999),
					Index: math.U64(1234567890),
				}
			},
		},
		{
			name: "boundary values",
			setup: func() *types.SlashingInfo {
				return &types.SlashingInfo{
					Slot:  math.Slot(1),
					Index: math.U64(^uint64(0) - 1),
				}
			},
		},
	}

	for _, pattern := range patterns {
		t.Run(pattern.name, func(t *testing.T) {
			original := pattern.setup()

			// Marshal
			data, err := original.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")

			// Unmarshal
			decoded := &types.SlashingInfo{}
			err = decoded.UnmarshalSSZ(data)
			require.NoError(t, err, "UnmarshalSSZ should not error")

			// Compare
			require.Equal(t, original, decoded, "round trip should preserve data")

			// Verify hash tree roots match
			root1, err := original.HashTreeRoot()
			require.NoError(t, err)
			root2, err := decoded.HashTreeRoot()
			require.NoError(t, err)
			require.Equal(t, root1, root2, "hash tree roots should match after round trip")
		})
	}
}

// TestSlashingInfoGettersSetters tests the getter and setter methods
func TestSlashingInfoGettersSetters(t *testing.T) {
	slashingInfo := &types.SlashingInfo{}

	// Test initial values
	require.Equal(t, math.Slot(0), slashingInfo.GetSlot(), "initial slot should be 0")
	require.Equal(t, math.U64(0), slashingInfo.GetIndex(), "initial index should be 0")

	// Test setters
	testSlot := math.Slot(12345)
	testIndex := math.U64(67890)

	slashingInfo.SetSlot(testSlot)
	slashingInfo.SetIndex(testIndex)

	// Test getters
	require.Equal(t, testSlot, slashingInfo.GetSlot(), "slot should match set value")
	require.Equal(t, testIndex, slashingInfo.GetIndex(), "index should match set value")

	// Test SSZ encoding after setting values
	data, err := slashingInfo.MarshalSSZ()
	require.NoError(t, err, "MarshalSSZ should not error after setting values")

	// Unmarshal and verify
	decoded := &types.SlashingInfo{}
	err = decoded.UnmarshalSSZ(data)
	require.NoError(t, err, "UnmarshalSSZ should not error")
	require.Equal(t, testSlot, decoded.GetSlot(), "slot should match after unmarshal")
	require.Equal(t, testIndex, decoded.GetIndex(), "index should match after unmarshal")
}

// TestSlashingInfoSSZFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestSlashingInfoSSZFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create deterministic "random" data based on iteration
		slot := math.Slot(uint64(i) * 1234567)
		index := math.U64(uint64(i) * 7654321)

		slashingInfo := &types.SlashingInfo{
			Slot:  slot,
			Index: index,
		}

		// Marshal
		data, err := slashingInfo.MarshalSSZ()
		require.NoError(t, err, "iteration %d: MarshalSSZ should not error", i)

		// Unmarshal
		decoded := &types.SlashingInfo{}
		err = decoded.UnmarshalSSZ(data)
		require.NoError(t, err, "iteration %d: UnmarshalSSZ should not error", i)

		// Verify
		require.Equal(t, slashingInfo, decoded, "iteration %d: round trip should preserve data", i)

		// Verify hash tree roots
		root1, err := slashingInfo.HashTreeRoot()
		require.NoError(t, err, "iteration %d: HashTreeRoot should not error", i)
		root2, err := decoded.HashTreeRoot()
		require.NoError(t, err, "iteration %d: decoded HashTreeRoot should not error", i)
		require.Equal(t, root1, root2, "iteration %d: hash tree roots should match", i)
	}
}

// TestSlashingInfoCompatibility tests that current and karalabe implementations produce identical results
func TestSlashingInfoCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.SlashingInfo, *SlashingInfoKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.SlashingInfo, *SlashingInfoKaralabe) {
				return &types.SlashingInfo{},
					&SlashingInfoKaralabe{}
			},
		},
		{
			name: "typical slashing info",
			setup: func() (*types.SlashingInfo, *SlashingInfoKaralabe) {
				slot := math.Slot(12345)
				index := math.U64(67890)

				current := &types.SlashingInfo{
					Slot:  slot,
					Index: index,
				}
				karalabe := &SlashingInfoKaralabe{
					Slot:  slot,
					Index: index,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.SlashingInfo, *SlashingInfoKaralabe) {
				slot := math.Slot(^uint64(0))
				index := math.U64(^uint64(0))

				current := &types.SlashingInfo{
					Slot:  slot,
					Index: index,
				}
				karalabe := &SlashingInfoKaralabe{
					Slot:  slot,
					Index: index,
				}
				return current, karalabe
			},
		},
		{
			name: "specific values",
			setup: func() (*types.SlashingInfo, *SlashingInfoKaralabe) {
				slot := math.Slot(1000000)
				index := math.U64(999)

				current := &types.SlashingInfo{
					Slot:  slot,
					Index: index,
				}
				karalabe := &SlashingInfoKaralabe{
					Slot:  slot,
					Index: index,
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

			karalabeBytes, err2 := karalabe.MarshalSSZ()
			require.NoError(t, err2, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabeBytes, currentBytes, "marshaled bytes should be identical")

			// Test Size
			require.Equal(t, types.SlashingInfoSize, current.SizeSSZ(), "current size should be 16")
			require.Equal(t, uint32(16), karalabe.SizeSSZ(), "karalabe size should be 16")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.SlashingInfo{}
			err := newCurrent.UnmarshalSSZ(karalabeBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &SlashingInfoKaralabe{}
			err = newKaralabe.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "unmarshal current data into karalabe should not error")
			require.Equal(t, karalabe, newKaralabe, "unmarshaled karalabe should match original")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabeRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabeRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}
