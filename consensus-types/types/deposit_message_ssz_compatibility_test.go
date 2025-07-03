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
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// KaralabeDepositMessage is a wrapper type with karalabe/ssz methods for testing compatibility
type KaralabeDepositMessage struct {
	DepositMessage // Embed to avoid duplicating fields
}

// SizeSSZ returns the size of the DepositMessage object in SSZ encoding.
func (*KaralabeDepositMessage) SizeSSZ(*ssz.Sizer) uint32 {
	//nolint:mnd // 48 + 32 + 8 = 88.
	return 88
}

// DefineSSZ defines the SSZ encoding for the DepositMessage object.
func (dm *KaralabeDepositMessage) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &dm.Pubkey)
	ssz.DefineStaticBytes(codec, &dm.Credentials)
	ssz.DefineUint64(codec, &dm.Amount)
}

// HashTreeRoot computes the SSZ hash tree root of the DepositMessage object.
func (dm *KaralabeDepositMessage) HashTreeRoot() common.Root {
	return ssz.HashSequential(dm)
}

// MarshalSSZTo marshals the DepositMessage object to SSZ format into the provided buffer.
func (dm *KaralabeDepositMessage) MarshalSSZTo(buf []byte) ([]byte, error) {
	return buf, ssz.EncodeToBytes(buf, dm)
}

// MarshalSSZ marshals the DepositMessage object to SSZ format.
func (dm *KaralabeDepositMessage) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(dm))
	return dm.MarshalSSZTo(buf)
}

// UnmarshalSSZ unmarshals the DepositMessage object from SSZ format.
func (dm *KaralabeDepositMessage) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, dm)
}

// Test all operations produce identical results
func TestDepositMessageSSZCompatibility(t *testing.T) {
	// Create test data
	pubkey := crypto.BLSPubkey{}
	copy(pubkey[:], []byte("test-pubkey-48-bytes-long-for-bls-signature-key"))

	credentials := WithdrawalCredentials{}
	copy(credentials[:], []byte("test-withdrawal-credentials-32by"))

	amount := math.Gwei(1000000000) // 1 ETH

	// Create instances
	depositMessage := &DepositMessage{
		Pubkey:      pubkey,
		Credentials: credentials,
		Amount:      amount,
	}

	karalabeDepositMessage := &KaralabeDepositMessage{
		DepositMessage: *depositMessage,
	}

	t.Run("Size calculation", func(t *testing.T) {
		// Both should report the same size
		require.Equal(t, int(karalabeDepositMessage.SizeSSZ(nil)), depositMessage.SizeSSZ())
		require.Equal(t, 88, depositMessage.SizeSSZ())
	})

	t.Run("Marshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeDepositMessage.MarshalSSZ()
		require.NoError(t, err)

		// Marshal with fastssz
		fastsszBytes, err := depositMessage.MarshalSSZ()
		require.NoError(t, err)

		// Should produce identical bytes
		require.Equal(t, karalabeBytes, fastsszBytes, "Marshaled bytes should be identical")
	})

	t.Run("Unmarshaling compatibility", func(t *testing.T) {
		// Marshal with karalabe
		karalabeBytes, err := karalabeDepositMessage.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		newDepositMessage := &DepositMessage{}
		err = newDepositMessage.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)
		require.Equal(t, depositMessage, newDepositMessage)

		// Marshal with fastssz
		fastsszBytes, err := depositMessage.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		newKaralabeDepositMessage := &KaralabeDepositMessage{}
		err = newKaralabeDepositMessage.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)
		require.Equal(t, karalabeDepositMessage.DepositMessage, newKaralabeDepositMessage.DepositMessage)
	})

	t.Run("Hash tree root compatibility", func(t *testing.T) {
		// Get HTR from karalabe
		karalabeHTR := karalabeDepositMessage.HashTreeRoot()

		// Get HTR from fastssz
		fastsszHTR, err := depositMessage.HashTreeRoot()
		require.NoError(t, err)

		// Compare HTRs - convert to same type for comparison
		require.Equal(t, [32]byte(karalabeHTR), fastsszHTR, "HashTreeRoot should be identical")
	})

	t.Run("Cross unmarshaling", func(t *testing.T) {
		// Create bytes with karalabe
		karalabeBytes, err := karalabeDepositMessage.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with fastssz
		fastsszDM := &DepositMessage{}
		err = fastsszDM.UnmarshalSSZ(karalabeBytes)
		require.NoError(t, err)

		// Create bytes with fastssz
		fastsszBytes, err := fastsszDM.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal with karalabe
		karalabeDM := &KaralabeDepositMessage{}
		err = karalabeDM.UnmarshalSSZ(fastsszBytes)
		require.NoError(t, err)

		// Both should be equal to original
		require.Equal(t, depositMessage, fastsszDM)
		require.Equal(t, karalabeDepositMessage.DepositMessage, karalabeDM.DepositMessage)
	})

	t.Run("Error cases", func(t *testing.T) {
		// Test unmarshaling with incorrect size
		shortBytes := make([]byte, 50) // Too short
		err := depositMessage.UnmarshalSSZ(shortBytes)
		require.Error(t, err)

		err = karalabeDepositMessage.UnmarshalSSZ(shortBytes)
		require.Error(t, err)
	})
}
