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

// ForkSize is the size of the Fork object in bytes.
// 4 bytes for PreviousVersion + 4 bytes for CurrentVersion + 8 bytes for Epoch.
const ForkSizeKaralabe = 16

// Compile-time assertions to ensure ForkKaralabe implements necessary interfaces.
var (
	_ ssz.StaticObject = (*ForkKaralabe)(nil)
)

// ForkKaralabe is an exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// Fork as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#fork
type ForkKaralabe struct {
	// PreviousVersion is the last version before the fork.
	PreviousVersion common.Version `json:"previous_version"`
	// CurrentVersion is the first version after the fork.
	CurrentVersion common.Version `json:"current_version"`
	// Epoch is the epoch at which the fork occurred.
	Epoch math.Epoch `json:"epoch"`
}

// NewForkKaralabe creates a new fork using karalabe/ssz implementation.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func NewForkKaralabe(
	previousVersion common.Version,
	currentVersion common.Version,
	epoch math.Epoch,
) *ForkKaralabe {
	return &ForkKaralabe{
		PreviousVersion: previousVersion,
		CurrentVersion:  currentVersion,
		Epoch:           epoch,
	}
}

// SizeSSZ returns the SSZ encoded size of the Fork object in bytes.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (f *ForkKaralabe) SizeSSZ() uint32 {
	return ForkSizeKaralabe
}

// DefineSSZ defines the SSZ encoding for the Fork object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (f *ForkKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &f.PreviousVersion)
	ssz.DefineStaticBytes(codec, &f.CurrentVersion)
	ssz.DefineUint64(codec, &f.Epoch)
}

// MarshalSSZ marshals the Fork object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (f *ForkKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(f))
	return buf, ssz.EncodeToBytes(buf, f)
}

// ValidateAfterDecodingSSZ validates the fork after decoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (*ForkKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// HashTreeRoot computes the SSZ hash tree root of the Fork object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (f *ForkKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(f)
}

// UnmarshalSSZ unmarshals the Fork object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (f *ForkKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, f)
}

// Compile-time assertions to ensure ForkDataKaralabe implements necessary interfaces.
var (
	_ ssz.StaticObject = (*ForkDataKaralabe)(nil)
)

// ForkDataKaralabe is an exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// ForkData as defined in the Ethereum 2.0 specification:
// https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#forkdata
type ForkDataKaralabe struct {
	// CurrentVersion is the current version of the fork.
	CurrentVersion common.Version
	// GenesisValidatorsRoot is the root of the genesis validators.
	GenesisValidatorsRoot common.Root
}

// NewForkDataKaralabe creates a new ForkData struct using karalabe/ssz implementation.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func NewForkDataKaralabe(
	currentVersion common.Version, genesisValidatorsRoot common.Root,
) *ForkDataKaralabe {
	return &ForkDataKaralabe{
		CurrentVersion:        currentVersion,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}
}

// SizeSSZ returns the size of the SigningData object in SSZ encoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (*ForkDataKaralabe) SizeSSZ() uint32 {
	//nolint:mnd // 32+4 = 36.
	return 36
}

// DefineSSZ defines the SSZ encoding for the ForkData object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (fd *ForkDataKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &fd.CurrentVersion)
	ssz.DefineStaticBytes(codec, &fd.GenesisValidatorsRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the ForkData object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (fd *ForkDataKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(fd)
}

// MarshalSSZ marshals the ForkData object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (fd *ForkDataKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(fd))
	return buf, ssz.EncodeToBytes(buf, fd)
}

// ValidateAfterDecodingSSZ validates the fork data after decoding.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (*ForkDataKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// UnmarshalSSZ unmarshals the ForkData object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use ssz.DecodeFromBytes
func (fd *ForkDataKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, fd)
}

// TestForkCompatibility tests that the current Fork implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestForkCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.Fork, *ForkKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.Fork, *ForkKaralabe) {
				current := types.NewFork(
					common.Version{0, 0, 0, 0},
					common.Version{0, 0, 0, 0},
					0,
				)
				karalabe := NewForkKaralabe(
					common.Version{0, 0, 0, 0},
					common.Version{0, 0, 0, 0},
					0,
				)
				return current, karalabe
			},
		},
		{
			name: "genesis fork",
			setup: func() (*types.Fork, *ForkKaralabe) {
				current := types.NewFork(
					common.Version{0, 0, 0, 0},
					common.Version{1, 0, 0, 0},
					0,
				)
				karalabe := NewForkKaralabe(
					common.Version{0, 0, 0, 0},
					common.Version{1, 0, 0, 0},
					0,
				)
				return current, karalabe
			},
		},
		{
			name: "fork transition",
			setup: func() (*types.Fork, *ForkKaralabe) {
				current := types.NewFork(
					common.Version{1, 0, 0, 0},
					common.Version{2, 0, 0, 0},
					math.Epoch(74240),
				)
				karalabe := NewForkKaralabe(
					common.Version{1, 0, 0, 0},
					common.Version{2, 0, 0, 0},
					math.Epoch(74240),
				)
				return current, karalabe
			},
		},
		{
			name: "deneb fork",
			setup: func() (*types.Fork, *ForkKaralabe) {
				current := types.NewFork(
					common.Version{3, 0, 0, 0},
					common.Version{4, 0, 0, 0},
					132608,
				)
				karalabe := NewForkKaralabe(
					common.Version{3, 0, 0, 0},
					common.Version{4, 0, 0, 0},
					132608,
				)
				return current, karalabe
			},
		},
		{
			name: "max values",
			setup: func() (*types.Fork, *ForkKaralabe) {
				current := types.NewFork(
					common.Version{0xFF, 0xFF, 0xFF, 0xFF},
					common.Version{0xFF, 0xFF, 0xFF, 0xFF},
					math.Epoch(0xFFFFFFFFFFFFFFFF),
				)
				karalabe := NewForkKaralabe(
					common.Version{0xFF, 0xFF, 0xFF, 0xFF},
					common.Version{0xFF, 0xFF, 0xFF, 0xFF},
					math.Epoch(0xFFFFFFFFFFFFFFFF),
				)
				return current, karalabe
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current, karalabe := tc.setup()

			// Test Marshal
			currentBytes, err := current.MarshalSSZ()
			require.NoError(t, err, "current MarshalSSZ should not error")

			karalabelBytes, err := karalabe.MarshalSSZ()
			require.NoError(t, err, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabelBytes, currentBytes, "marshaled bytes should be identical")

			// Test Unmarshal
			currentUnmarshaled := &types.Fork{}
			err = currentUnmarshaled.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "current UnmarshalSSZ should not error")

			karalabelUnmarshaled := &ForkKaralabe{}
			err = karalabelUnmarshaled.UnmarshalSSZ(karalabelBytes)
			require.NoError(t, err, "karalabe UnmarshalSSZ should not error")

			// Verify unmarshaled data matches
			require.Equal(t, current.PreviousVersion, currentUnmarshaled.PreviousVersion)
			require.Equal(t, karalabe.PreviousVersion, karalabelUnmarshaled.PreviousVersion)
			require.Equal(t, current.CurrentVersion, currentUnmarshaled.CurrentVersion)
			require.Equal(t, karalabe.CurrentVersion, karalabelUnmarshaled.CurrentVersion)
			require.Equal(t, current.Epoch, currentUnmarshaled.Epoch)
			require.Equal(t, karalabe.Epoch, karalabelUnmarshaled.Epoch)

			// Test Size
			require.Equal(t, int(karalabe.SizeSSZ()), current.SizeSSZ(), "sizes should be identical")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestForkDataCompatibility tests that the current ForkData implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestForkDataCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.ForkData, *ForkDataKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.ForkData, *ForkDataKaralabe) {
				current := &types.ForkData{
					CurrentVersion:        common.Version{0, 0, 0, 0},
					GenesisValidatorsRoot: common.Root{},
				}
				karalabe := &ForkDataKaralabe{
					CurrentVersion:        common.Version{0, 0, 0, 0},
					GenesisValidatorsRoot: common.Root{},
				}
				return current, karalabe
			},
		},
		{
			name: "with version",
			setup: func() (*types.ForkData, *ForkDataKaralabe) {
				current := &types.ForkData{
					CurrentVersion:        common.Version{1, 0, 0, 0},
					GenesisValidatorsRoot: common.Root{},
				}
				karalabe := &ForkDataKaralabe{
					CurrentVersion:        common.Version{1, 0, 0, 0},
					GenesisValidatorsRoot: common.Root{},
				}
				return current, karalabe
			},
		},
		{
			name: "with root",
			setup: func() (*types.ForkData, *ForkDataKaralabe) {
				root := common.Root{
					1, 2, 3, 4, 5, 6, 7, 8,
					9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24,
					25, 26, 27, 28, 29, 30, 31, 32,
				}
				current := &types.ForkData{
					CurrentVersion:        common.Version{0, 0, 0, 0},
					GenesisValidatorsRoot: root,
				}
				karalabe := &ForkDataKaralabe{
					CurrentVersion:        common.Version{0, 0, 0, 0},
					GenesisValidatorsRoot: root,
				}
				return current, karalabe
			},
		},
		{
			name: "mainnet genesis",
			setup: func() (*types.ForkData, *ForkDataKaralabe) {
				root := common.Root{
					0x4b, 0x36, 0x3d, 0xb9, 0x4e, 0x28, 0x61, 0x20,
					0xd9, 0x70, 0x7e, 0xba, 0x52, 0x04, 0x80, 0xb8,
					0x6c, 0x90, 0x27, 0xb4, 0x5b, 0x25, 0x4d, 0x74,
					0x25, 0x60, 0xba, 0xa0, 0xa5, 0x70, 0x3f, 0x8f,
				}
				current := &types.ForkData{
					CurrentVersion:        common.Version{0, 0, 0, 0},
					GenesisValidatorsRoot: root,
				}
				karalabe := &ForkDataKaralabe{
					CurrentVersion:        common.Version{0, 0, 0, 0},
					GenesisValidatorsRoot: root,
				}
				return current, karalabe
			},
		},
		{
			name: "complex values",
			setup: func() (*types.ForkData, *ForkDataKaralabe) {
				root := common.Root{
					0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
					0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x00,
					0x0f, 0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09, 0x08,
					0x07, 0x06, 0x05, 0x04, 0x03, 0x02, 0x01, 0x00,
				}
				current := &types.ForkData{
					CurrentVersion:        common.Version{0xaa, 0xbb, 0xcc, 0xdd},
					GenesisValidatorsRoot: root,
				}
				karalabe := &ForkDataKaralabe{
					CurrentVersion:        common.Version{0xaa, 0xbb, 0xcc, 0xdd},
					GenesisValidatorsRoot: root,
				}
				return current, karalabe
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			current, karalabe := tc.setup()

			// Test Marshal
			currentBytes, err := current.MarshalSSZ()
			require.NoError(t, err, "current MarshalSSZ should not error")

			karalabelBytes, err := karalabe.MarshalSSZ()
			require.NoError(t, err, "karalabe MarshalSSZ should not error")

			require.Equal(t, karalabelBytes, currentBytes, "marshaled bytes should be identical")

			// Test Unmarshal
			currentUnmarshaled := &types.ForkData{}
			err = currentUnmarshaled.UnmarshalSSZ(currentBytes)
			require.NoError(t, err, "current UnmarshalSSZ should not error")

			karalabelUnmarshaled := &ForkDataKaralabe{}
			err = karalabelUnmarshaled.UnmarshalSSZ(karalabelBytes)
			require.NoError(t, err, "karalabe UnmarshalSSZ should not error")

			// Verify unmarshaled data matches
			require.Equal(t, current.CurrentVersion, currentUnmarshaled.CurrentVersion)
			require.Equal(t, karalabe.CurrentVersion, karalabelUnmarshaled.CurrentVersion)
			require.Equal(t, current.GenesisValidatorsRoot, currentUnmarshaled.GenesisValidatorsRoot)
			require.Equal(t, karalabe.GenesisValidatorsRoot, karalabelUnmarshaled.GenesisValidatorsRoot)

			// Test Size
			require.Equal(t, int(karalabe.SizeSSZ()), current.SizeSSZ(), "sizes should be identical")

			// Test HashTreeRoot
			currentRoot, err := current.HashTreeRoot()
			require.NoError(t, err, "current HashTreeRoot should not error")
			karalabelRoot := karalabe.HashTreeRoot()
			require.Equal(t, [32]byte(karalabelRoot), currentRoot, "hash tree roots should be identical")
		})
	}
}

// TestForkCompatibilityFuzz tests Fork with random data
func TestForkCompatibilityFuzz(t *testing.T) {
	// Test with pseudo-random data
	for i := 0; i < 100; i++ {
		// Create pseudo-random values
		prevVersion := common.Version{byte(i), byte(i + 1), byte(i + 2), byte(i + 3)}
		currVersion := common.Version{byte(i + 4), byte(i + 5), byte(i + 6), byte(i + 7)}
		epoch := math.Epoch(uint64(i) * 12345)

		current := types.NewFork(prevVersion, currVersion, epoch)
		karalabe := NewForkKaralabe(prevVersion, currVersion, epoch)

		// Marshal both
		currentBytes, err := current.MarshalSSZ()
		require.NoError(t, err)
		karalabelBytes, err := karalabe.MarshalSSZ()
		require.NoError(t, err)

		// They should be identical
		require.Equal(t, karalabelBytes, currentBytes, "fuzzing: marshaled bytes should be identical for iteration %d", i)
	}
}

// TestForkDataCompatibilityFuzz tests ForkData with random data
func TestForkDataCompatibilityFuzz(t *testing.T) {
	// Test with pseudo-random data
	for i := 0; i < 100; i++ {
		// Create pseudo-random values
		version := common.Version{byte(i), byte(i + 1), byte(i + 2), byte(i + 3)}
		var root common.Root
		for j := range root {
			root[j] = byte((i + j) % 256)
		}

		current := &types.ForkData{
			CurrentVersion:        version,
			GenesisValidatorsRoot: root,
		}
		karalabe := &ForkDataKaralabe{
			CurrentVersion:        version,
			GenesisValidatorsRoot: root,
		}

		// Marshal both
		currentBytes, err := current.MarshalSSZ()
		require.NoError(t, err)
		karalabelBytes, err := karalabe.MarshalSSZ()
		require.NoError(t, err)

		// They should be identical
		require.Equal(t, karalabelBytes, currentBytes, "fuzzing: marshaled bytes should be identical for iteration %d", i)
	}
}

// TestForkCompatibilityRoundTrip tests that data can be marshaled/unmarshaled between implementations
func TestForkCompatibilityRoundTrip(t *testing.T) {
	// Create a fork with specific values
	prevVersion := common.Version{1, 2, 3, 4}
	currVersion := common.Version{5, 6, 7, 8}
	epoch := math.Epoch(99999)

	// Create using current implementation
	current := types.NewFork(prevVersion, currVersion, epoch)

	// Marshal with current
	currentBytes, err := current.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe
	karalabe := &ForkKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Verify values
	require.Equal(t, prevVersion, karalabe.PreviousVersion)
	require.Equal(t, currVersion, karalabe.CurrentVersion)
	require.Equal(t, epoch, karalabe.Epoch)

	// Marshal with karalabe
	karalabelBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current
	newCurrent := &types.Fork{}
	err = newCurrent.UnmarshalSSZ(karalabelBytes)
	require.NoError(t, err)

	// Verify values
	require.Equal(t, prevVersion, newCurrent.PreviousVersion)
	require.Equal(t, currVersion, newCurrent.CurrentVersion)
	require.Equal(t, epoch, newCurrent.Epoch)

	// Bytes should be identical
	require.Equal(t, currentBytes, karalabelBytes)
}

// TestForkDataCompatibilityRoundTrip tests that data can be marshaled/unmarshaled between implementations
func TestForkDataCompatibilityRoundTrip(t *testing.T) {
	// Create fork data with specific values
	version := common.Version{9, 10, 11, 12}
	var root common.Root
	for i := range root {
		root[i] = byte(255 - i)
	}

	// Create using current implementation
	current := &types.ForkData{
		CurrentVersion:        version,
		GenesisValidatorsRoot: root,
	}

	// Marshal with current
	currentBytes, err := current.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe
	karalabe := &ForkDataKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Verify values
	require.Equal(t, version, karalabe.CurrentVersion)
	require.Equal(t, root, karalabe.GenesisValidatorsRoot)

	// Marshal with karalabe
	karalabelBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current
	newCurrent := &types.ForkData{}
	err = newCurrent.UnmarshalSSZ(karalabelBytes)
	require.NoError(t, err)

	// Verify values
	require.Equal(t, version, newCurrent.CurrentVersion)
	require.Equal(t, root, newCurrent.GenesisValidatorsRoot)

	// Bytes should be identical
	require.Equal(t, currentBytes, karalabelBytes)
}
