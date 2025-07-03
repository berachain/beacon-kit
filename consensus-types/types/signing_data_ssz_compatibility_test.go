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

// KaralabeSigningData is a wrapper type with karalabe/ssz methods for testing compatibility
type KaralabeSigningData struct {
	SigningData // Embed to avoid duplicating fields
}

// SizeSSZ returns the size of the SigningData object in SSZ encoding.
func (*KaralabeSigningData) SizeSSZ(_ *ssz.Sizer) uint32 {
	//nolint:mnd // 32*2 = 64.
	return 64
}

// DefineSSZ defines the SSZ encoding for the SigningData object.
func (s *KaralabeSigningData) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &s.ObjectRoot)
	ssz.DefineStaticBytes(codec, &s.Domain)
}

// HashTreeRoot computes the SSZ hash tree root of the SigningData object.
func (s *KaralabeSigningData) HashTreeRoot() common.Root {
	return ssz.HashSequential(s)
}

// MarshalSSZTo marshals the SigningData object to SSZ format into the provided buffer.
func (s *KaralabeSigningData) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, s)
}

// MarshalSSZ marshals the SigningData object to SSZ format.
func (s *KaralabeSigningData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(s))
	return s.MarshalSSZTo(buf)
}

// UnmarshalSSZ unmarshals the SigningData object from SSZ format.
func (s *KaralabeSigningData) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, s)
}

// Test all operations produce identical results
func TestSigningDataSSZCompatibility(t *testing.T) {
	// Create test data
	objectRoot := common.Root{}
	copy(objectRoot[:], []byte("test-object-root-32-bytes-longrr"))

	domain := common.Domain{}
	copy(domain[:], []byte("test-domain-32-bytes-long-value!"))

	// Create instances
	signingData := &SigningData{
		ObjectRoot: objectRoot,
		Domain:     domain,
	}

	karalabeSigningData := &KaralabeSigningData{
		SigningData: *signingData,
	}

	t.Run("Size calculation", func(t *testing.T) {
		// Both should report the same size
		require.Equal(t, int(karalabeSigningData.SizeSSZ(nil)), signingData.SizeSSZ())
		require.Equal(t, 64, signingData.SizeSSZ())
	})

	t.Run("Marshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeSigningData.MarshalSSZ()
		require.NoError(t, err)

		// Marshal with fastssz
		fastsszBytes, err := signingData.MarshalSSZ()
		require.NoError(t, err)

		// Should produce identical bytes
		require.Equal(t, karalabeBytes, fastsszBytes, "Marshaled bytes should be identical")
	})

	t.Run("Unmarshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeSigningData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		newSigningData := &SigningData{}
		err = newSigningData.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)
		require.Equal(t, signingData, newSigningData)

		// Marshal with fastssz
		fastsszBytes, err := signingData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		newKaralabeSigningData := &KaralabeSigningData{}
		err = newKaralabeSigningData.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)
		require.Equal(t, karalabeSigningData.SigningData, newKaralabeSigningData.SigningData)
	})

	t.Run("Hash tree root compatibility", func(t *testing.T) {
		// Get HTR from karalabe
		karalabeHTR := karalabeSigningData.HashTreeRoot()

		// Get HTR from fastssz
		fastsszHTR, err := signingData.HashTreeRoot()
		require.NoError(t, err)

		// Compare HTRs - convert to same type for comparison
		require.Equal(t, [32]byte(karalabeHTR), fastsszHTR, "HashTreeRoot should be identical")
	})

	t.Run("Cross unmarshaling", func(t *testing.T) {
		// Create bytes with karalabe
		karalabeBytes, err := karalabeSigningData.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		fastsszSD := &SigningData{}
		err = fastsszSD.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)

		// Create bytes with fastssz
		fastsszBytes, err := fastsszSD.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		karalabeSD := &KaralabeSigningData{}
		err = karalabeSD.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)

		// Both should be equal to original
		require.Equal(t, signingData, fastsszSD)
		require.Equal(t, karalabeSigningData.SigningData, karalabeSD.SigningData)
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test unmarshaling with incorrect size
		shortBytes := make([]byte, 30) // Too short
		err := signingData.UnmarshalSSZ(shortBytes)
		require.Error(t, err)

		err = karalabeSigningData.UnmarshalSSZ(shortBytes)
		require.Error(t, err)
	})
}
