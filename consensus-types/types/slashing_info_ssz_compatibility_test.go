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

// KaralabeSlashingInfo is a wrapper type with karalabe/ssz methods for testing compatibility
type KaralabeSlashingInfo struct {
	SlashingInfo // Embed to avoid duplicating fields
}

// SizeSSZ returns the size of the SlashingInfo object in SSZ encoding.
func (*KaralabeSlashingInfo) SizeSSZ(*ssz.Sizer) uint32 {
	return SlashingInfoSize
}

// DefineSSZ defines the SSZ encoding for the SlashingInfo object.
func (s *KaralabeSlashingInfo) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineUint64(codec, &s.Slot)
	ssz.DefineUint64(codec, &s.Index)
}

// HashTreeRoot computes the SSZ hash tree root of the SlashingInfo object.
func (s *KaralabeSlashingInfo) HashTreeRoot() common.Root {
	return ssz.HashSequential(s)
}

// MarshalSSZ marshals the SlashingInfo object to SSZ format.
func (s *KaralabeSlashingInfo) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(s))
	return buf, ssz.EncodeToBytes(buf, s)
}

// UnmarshalSSZ unmarshals the SlashingInfo object from SSZ format.
func (s *KaralabeSlashingInfo) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

// Test all operations produce identical results
func TestSlashingInfoSSZCompatibility(t *testing.T) {
	// Create test data
	testSlot := math.Slot(12345)
	testIndex := math.U64(67890)

	// Create instances
	slashingInfo := &SlashingInfo{
		Slot:  testSlot,
		Index: testIndex,
	}

	karalabeSlashingInfo := &KaralabeSlashingInfo{
		SlashingInfo: *slashingInfo,
	}

	t.Run("Size calculation", func(t *testing.T) {
		// Both should report the same size
		require.Equal(t, int(karalabeSlashingInfo.SizeSSZ(nil)), slashingInfo.SizeSSZ())
		require.Equal(t, 16, slashingInfo.SizeSSZ())
	})

	t.Run("Marshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeSlashingInfo.MarshalSSZ()
		require.NoError(t, err)

		// Marshal with fastssz
		fastsszBytes, err := slashingInfo.MarshalSSZ()
		require.NoError(t, err)

		// Should produce identical bytes
		require.Equal(t, karalabeBytes, fastsszBytes, "Marshaled bytes should be identical")
	})

	t.Run("Unmarshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeSlashingInfo.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		newSlashingInfo := &SlashingInfo{}
		err = newSlashingInfo.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)
		require.Equal(t, slashingInfo, newSlashingInfo)

		// Marshal with fastssz
		fastsszBytes, err := slashingInfo.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		newKaralabeSlashingInfo := &KaralabeSlashingInfo{}
		err = newKaralabeSlashingInfo.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)
		require.Equal(t, karalabeSlashingInfo.SlashingInfo, newKaralabeSlashingInfo.SlashingInfo)
	})

	t.Run("Hash tree root compatibility", func(t *testing.T) {
		// Get HTR from karalabe
		karalabeHTR := karalabeSlashingInfo.HashTreeRoot()

		// Get HTR from fastssz
		fastsszHTR, err := slashingInfo.HashTreeRoot()
		require.NoError(t, err)

		// Compare HTRs - convert to same type for comparison
		require.Equal(t, [32]byte(karalabeHTR), fastsszHTR, "HashTreeRoot should be identical")
	})

	t.Run("Cross unmarshaling", func(t *testing.T) {
		// Create bytes with karalabe
		karalabeBytes, err := karalabeSlashingInfo.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		fastsszSI := &SlashingInfo{}
		err = fastsszSI.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)

		// Create bytes with fastssz
		fastsszBytes, err := fastsszSI.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		karalabeSI := &KaralabeSlashingInfo{}
		err = karalabeSI.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)

		// Both should be equal to original
		require.Equal(t, slashingInfo, fastsszSI)
		require.Equal(t, karalabeSlashingInfo.SlashingInfo, karalabeSI.SlashingInfo)
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test unmarshaling with incorrect size
		shortBytes := make([]byte, 8) // Too short
		err := slashingInfo.UnmarshalSSZ(shortBytes)
		require.Error(t, err)

		err = karalabeSlashingInfo.UnmarshalSSZ(shortBytes)
		require.Error(t, err)
	})

	t.Run("Getters and Setters", func(t *testing.T) {
		// Test getters
		require.Equal(t, testSlot, slashingInfo.GetSlot())
		require.Equal(t, testIndex, slashingInfo.GetIndex())

		// Test setters
		newSlot := math.Slot(99999)
		newIndex := math.U64(11111)
		
		slashingInfo.SetSlot(newSlot)
		slashingInfo.SetIndex(newIndex)
		
		require.Equal(t, newSlot, slashingInfo.GetSlot())
		require.Equal(t, newIndex, slashingInfo.GetIndex())
	})
}