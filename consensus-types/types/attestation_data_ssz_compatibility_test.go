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

// KaralabeAttestationData is a wrapper type with karalabe/ssz methods for testing compatibility
type KaralabeAttestationData struct {
	AttestationData // Embed to avoid duplicating fields
}

// SizeSSZ returns the size of the AttestationData object in SSZ encoding.
func (*KaralabeAttestationData) SizeSSZ(*ssz.Sizer) uint32 {
	return AttestationDataSize
}

// DefineSSZ defines the SSZ encoding for the AttestationData object.
func (a *KaralabeAttestationData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &a.Slot)
	ssz.DefineUint64(codec, &a.Index)
	ssz.DefineStaticBytes(codec, &a.BeaconBlockRoot)
}

// HashTreeRoot computes the SSZ hash tree root of the AttestationData object.
func (a *KaralabeAttestationData) HashTreeRoot() common.Root {
	return ssz.HashSequential(a)
}

// MarshalSSZ marshals the AttestationData object to SSZ format.
func (a *KaralabeAttestationData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(a))
	return buf, ssz.EncodeToBytes(buf, a)
}

// UnmarshalSSZ unmarshals the AttestationData object from SSZ format.
func (a *KaralabeAttestationData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, a)
}

// Test all operations produce identical results
func TestAttestationDataSSZCompatibility(t *testing.T) {
	// Create test data
	slot := math.U64(12345)
	index := math.U64(67890)
	beaconBlockRoot := common.Root{}
	copy(beaconBlockRoot[:], []byte("test-beacon-block-root-32-bytes!"))

	// Create instances
	attestationData := &AttestationData{
		Slot:            slot,
		Index:           index,
		BeaconBlockRoot: beaconBlockRoot,
	}

	karalabeAttestationData := &KaralabeAttestationData{
		AttestationData: *attestationData,
	}

	t.Run("Size calculation", func(t *testing.T) {
		// Both should report the same size
		require.Equal(t, int(karalabeAttestationData.SizeSSZ(nil)), attestationData.SizeSSZ())
		require.Equal(t, 48, attestationData.SizeSSZ())
	})

	t.Run("Marshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeAttestationData.MarshalSSZ()
		require.NoError(t, err)

		// Marshal with fastssz
		fastsszBytes, err := attestationData.MarshalSSZ()
		require.NoError(t, err)

		// Should produce identical bytes
		require.Equal(t, karalabeBytes, fastsszBytes, "Marshaled bytes should be identical")
	})

	t.Run("Unmarshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeAttestationData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		newAttestationData := &AttestationData{}
		err = newAttestationData.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)
		require.Equal(t, attestationData, newAttestationData)

		// Marshal with fastssz
		fastsszBytes, err := attestationData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		newKaralabeAttestationData := &KaralabeAttestationData{}
		err = newKaralabeAttestationData.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)
		require.Equal(t, karalabeAttestationData.AttestationData, newKaralabeAttestationData.AttestationData)
	})

	t.Run("Hash tree root compatibility", func(t *testing.T) {
		// Get HTR from karalabe
		karalabeHTR := karalabeAttestationData.HashTreeRoot()

		// Get HTR from fastssz
		fastsszHTR, err := attestationData.HashTreeRoot()
		require.NoError(t, err)

		// Compare HTRs - convert to same type for comparison
		require.Equal(t, [32]byte(karalabeHTR), fastsszHTR, "HashTreeRoot should be identical")
	})

	t.Run("Cross unmarshaling", func(t *testing.T) {
		// Create bytes with karalabe
		karalabeBytes, err := karalabeAttestationData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		fastsszAD := &AttestationData{}
		err = fastsszAD.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)

		// Create bytes with fastssz
		fastsszBytes, err := fastsszAD.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		karalabeAD := &KaralabeAttestationData{}
		err = karalabeAD.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)

		// Both should be equal to original
		require.Equal(t, attestationData, fastsszAD)
		require.Equal(t, karalabeAttestationData.AttestationData, karalabeAD.AttestationData)
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test unmarshaling with incorrect size
		shortBytes := make([]byte, 20) // Too short
		err := attestationData.UnmarshalSSZ(shortBytes)
		require.Error(t, err)

		err = karalabeAttestationData.UnmarshalSSZ(shortBytes)
		require.Error(t, err)
	})

	t.Run("Getters", func(t *testing.T) {
		// Test getters
		require.Equal(t, slot, attestationData.GetSlot())
		require.Equal(t, index, attestationData.GetIndex())
		require.Equal(t, beaconBlockRoot, attestationData.GetBeaconBlockRoot())
	})
}