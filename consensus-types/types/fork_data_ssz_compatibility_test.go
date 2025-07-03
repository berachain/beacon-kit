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
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// KaralabeForkData is a wrapper type with karalabe/ssz methods for testing compatibility
type KaralabeForkData struct {
	ForkData // Embed to avoid duplicating fields
}

// SizeSSZ returns the size of the ForkData object in SSZ encoding.
func (*KaralabeForkData) SizeSSZ(*ssz.Sizer) uint32 {
	//nolint:mnd // 32+4 = 36.
	return 36
}

// DefineSSZ defines the SSZ encoding for the ForkData object.
func (fd *KaralabeForkData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &fd.CurrentVersion)
	ssz.DefineStaticBytes(codec, &fd.GenesisValidatorsRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the ForkData object.
func (fd *KaralabeForkData) HashTreeRoot() common.Root {
	return ssz.HashSequential(fd)
}

// MarshalSSZTo marshals the ForkData object to SSZ format into the provided buffer.
func (fd *KaralabeForkData) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, fd)
}

// MarshalSSZ marshals the ForkData object to SSZ format.
func (fd *KaralabeForkData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(fd))
	return fd.MarshalSSZTo(buf)
}

// UnmarshalSSZ unmarshals the ForkData object from SSZ format.
func (fd *KaralabeForkData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, fd)
}

// Test all operations produce identical results
func TestForkDataSSZCompatibility(t *testing.T) {
	// Create test data
	version := common.Version{0x01, 0x02, 0x03, 0x04}
	genesisValidatorsRoot := common.Root{}
	copy(genesisValidatorsRoot[:], []byte("test-genesis-validators-root-32b"))

	// Create instances
	forkData := &ForkData{
		CurrentVersion:        version,
		GenesisValidatorsRoot: genesisValidatorsRoot,
	}

	karalabeForkData := &KaralabeForkData{
		ForkData: *forkData,
	}

	t.Run("Size calculation", func(t *testing.T) {
		// Both should report the same size
		require.Equal(t, int(karalabeForkData.SizeSSZ(nil)), forkData.SizeSSZ())
		require.Equal(t, 36, forkData.SizeSSZ())
	})

	t.Run("Marshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeForkData.MarshalSSZ()
		require.NoError(t, err)

		// Marshal with fastssz
		fastsszBytes, err := forkData.MarshalSSZ()
		require.NoError(t, err)

		// Should produce identical bytes
		require.Equal(t, karalabeBytes, fastsszBytes, "Marshaled bytes should be identical")
	})

	t.Run("Unmarshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeForkData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		newForkData := &ForkData{}
		err = newForkData.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)
		require.Equal(t, forkData, newForkData)

		// Marshal with fastssz
		fastsszBytes, err := forkData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		newKaralabeForkData := &KaralabeForkData{}
		err = newKaralabeForkData.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)
		require.Equal(t, karalabeForkData.ForkData, newKaralabeForkData.ForkData)
	})

	t.Run("Hash tree root compatibility", func(t *testing.T) {
		// Get HTR from karalabe
		karalabeHTR := karalabeForkData.HashTreeRoot()

		// Get HTR from fastssz
		fastsszHTR, err := forkData.HashTreeRoot()
		require.NoError(t, err)

		// Compare HTRs - convert to same type for comparison
		require.Equal(t, [32]byte(karalabeHTR), fastsszHTR, "HashTreeRoot should be identical")
	})

	t.Run("Cross unmarshaling", func(t *testing.T) {
		// Create bytes with karalabe
		karalabeBytes, err := karalabeForkData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		fastsszFD := &ForkData{}
		err = fastsszFD.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)

		// Create bytes with fastssz
		fastsszBytes, err := fastsszFD.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		karalabeFD := &KaralabeForkData{}
		err = karalabeFD.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)

		// Both should be equal to original
		require.Equal(t, forkData, fastsszFD)
		require.Equal(t, karalabeForkData.ForkData, karalabeFD.ForkData)
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test unmarshaling with incorrect size
		shortBytes := make([]byte, 20) // Too short
		err := forkData.UnmarshalSSZ(shortBytes)
		require.Error(t, err)

		err = karalabeForkData.UnmarshalSSZ(shortBytes)
		require.Error(t, err)
	})
}
