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

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// ForkSize is the size of the Fork object in bytes.
// 4 bytes for PreviousVersion + 4 bytes for CurrentVersion + 8 bytes for Epoch.
const ForkSize = 16

// KaralabeFork embeds Fork and adds karalabe/ssz methods
type KaralabeFork struct {
	Fork
}

// SizeSSZ returns the SSZ encoded size of the Fork object in bytes.
func (f *KaralabeFork) SizeSSZ(*ssz.Sizer) uint32 {
	return ForkSize
}

// DefineSSZ defines the SSZ encoding for the Fork object.
func (f *KaralabeFork) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &f.PreviousVersion)
	ssz.DefineStaticBytes(codec, &f.CurrentVersion)
	ssz.DefineUint64(codec, &f.Epoch)
}

// MarshalSSZ marshals the Fork object to SSZ format.
func (f *KaralabeFork) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(f))
	return buf, ssz.EncodeToBytes(buf, f)
}

// HashTreeRoot computes the SSZ hash tree root of the Fork object.
func (f *KaralabeFork) HashTreeRoot() common.Root {
	return ssz.HashSequential(f)
}

func TestForkSSZCompatibility(t *testing.T) {
	tests := []struct {
		name            string
		previousVersion common.Version
		currentVersion  common.Version
		epoch           math.Epoch
	}{
		{
			name:            "zero values",
			previousVersion: common.Version{0, 0, 0, 0},
			currentVersion:  common.Version{0, 0, 0, 0},
			epoch:           0,
		},
		{
			name:            "typical fork values",
			previousVersion: common.Version{0x01, 0x00, 0x00, 0x00},
			currentVersion:  common.Version{0x02, 0x00, 0x00, 0x00},
			epoch:           100,
		},
		{
			name:            "max values",
			previousVersion: common.Version{0xFF, 0xFF, 0xFF, 0xFF},
			currentVersion:  common.Version{0xFF, 0xFF, 0xFF, 0xFF},
			epoch:           math.Epoch(^uint64(0)),
		},
		{
			name:            "different version bytes",
			previousVersion: common.Version{0x12, 0x34, 0x56, 0x78},
			currentVersion:  common.Version{0x9A, 0xBC, 0xDE, 0xF0},
			epoch:           12345,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create Fork instances
			fork := NewFork(tt.previousVersion, tt.currentVersion, tt.epoch)
			karalabeFork := &KaralabeFork{Fork: *fork}

			// Test marshaling
			t.Run("marshal", func(t *testing.T) {
				// Marshal using fastssz
				bytes, err := fork.MarshalSSZ()
				require.NoError(t, err)

				// Marshal using karalabe/ssz
				karalabeBytes, err := karalabeFork.MarshalSSZ()
				require.NoError(t, err)

				// Compare marshaled bytes
				require.Equal(t, karalabeBytes, bytes, "Marshaled bytes should be identical between fastssz and karalabe/ssz")

				// Verify expected size
				require.Equal(t, 16, len(bytes), "Fork should marshal to exactly 16 bytes")
				require.Equal(t, 16, len(karalabeBytes), "Fork should marshal to exactly 16 bytes")
			})

			// Test unmarshaling
			t.Run("unmarshal", func(t *testing.T) {
				// Marshal with karalabe to get bytes
				karalabeBytes, err := karalabeFork.MarshalSSZ()
				require.NoError(t, err)

				// Unmarshal using fastssz
				fastsszResult := &Fork{}
				err = fastsszResult.UnmarshalSSZ(karalabeBytes)
				require.NoError(t, err)

				// Compare unmarshaled values
				require.Equal(t, fork.PreviousVersion[:], fastsszResult.PreviousVersion[:], "PreviousVersion should be identical")
				require.Equal(t, fork.CurrentVersion[:], fastsszResult.CurrentVersion[:], "CurrentVersion should be identical")
				require.Equal(t, fork.Epoch, fastsszResult.Epoch, "Epoch should be identical")

				// Now test the reverse - marshal with fastssz and unmarshal with karalabe
				fastsszBytes, err := fork.MarshalSSZ()
				require.NoError(t, err)

				// Unmarshal with karalabe
				karalabeResult := &KaralabeFork{}
				err = ssz.DecodeFromBytes(fastsszBytes, karalabeResult)
				require.NoError(t, err)

				// Compare values
				require.Equal(t, fork.PreviousVersion[:], karalabeResult.PreviousVersion[:])
				require.Equal(t, fork.CurrentVersion[:], karalabeResult.CurrentVersion[:])
				require.Equal(t, fork.Epoch, karalabeResult.Epoch)
			})

			// Test size methods
			t.Run("size", func(t *testing.T) {
				// Get size from fastssz
				fastsszSize := fork.SizeSSZ()

				// Get size from karalabe
				karalabeSize := karalabeFork.SizeSSZ(nil)

				// Compare sizes
				require.Equal(t, uint32(fastsszSize), karalabeSize, "Size should be identical between implementations")
				require.Equal(t, uint32(16), karalabeSize, "Fork should always be 16 bytes")
			})

			// Test HashTreeRoot
			t.Run("hash_tree_root", func(t *testing.T) {
				// Get HTR from fastssz
				fastsszHTR, err := fork.HashTreeRoot()
				require.NoError(t, err)

				// Get HTR from karalabe
				karalabeHTR := karalabeFork.HashTreeRoot()

				// Compare HTRs - convert to same type for comparison
				require.Equal(t, [32]byte(karalabeHTR), fastsszHTR, "HashTreeRoot should be identical between implementations")
			})
		})
	}
}

// TestForkSSZErrorCases tests error handling in both implementations
func TestForkSSZErrorCases(t *testing.T) {
	t.Run("unmarshal wrong size", func(t *testing.T) {
		// Test with too few bytes
		shortBytes := make([]byte, 10)

		// Test fastssz
		fork := &Fork{}
		err := fork.UnmarshalSSZ(shortBytes)
		require.Error(t, err, "fastssz should error on wrong size")

		// Test karalabe/ssz
		karalabeFork := &KaralabeFork{}
		err = ssz.DecodeFromBytes(shortBytes, karalabeFork)
		require.Error(t, err, "karalabe/ssz should error on wrong size")

		// Test with too many bytes
		longBytes := make([]byte, 20)

		// Test fastssz
		err = fork.UnmarshalSSZ(longBytes)
		require.Error(t, err, "fastssz should error on wrong size")

		// Test karalabe/ssz
		err = ssz.DecodeFromBytes(longBytes, karalabeFork)
		require.Error(t, err, "karalabe/ssz should error on wrong size")
	})
}
