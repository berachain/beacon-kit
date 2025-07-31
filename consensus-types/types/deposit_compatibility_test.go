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
	karalabe "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

// depositSizeKaralabe is the size of the SSZ encoding of a Deposit.
const depositSizeKaralabe = 192 // 48 + 32 + 8 + 96 + 8

// Compile-time assertions to ensure DepositKaralabe implements necessary interfaces.
var _ karalabe.StaticObject = (*DepositKaralabe)(nil)

// DepositKaralabe is an exact copy of Deposit from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// This type uses karalabe/ssz for SSZ operations to ensure compatibility testing.
type DepositKaralabe struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials types.WithdrawalCredentials `json:"credentials"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
	// Signature of the deposit data.
	Signature crypto.BLSSignature `json:"signature"`
	// Index of the deposit in the deposit contract.
	Index uint64 `json:"index"`
}

// DefineSSZ defines the SSZ encoding for the Deposit object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (d *DepositKaralabe) DefineSSZ(c *karalabe.Codec) {
	karalabe.DefineStaticBytes(c, &d.Pubkey)
	karalabe.DefineStaticBytes(c, &d.Credentials)
	karalabe.DefineUint64(c, &d.Amount)
	karalabe.DefineStaticBytes(c, &d.Signature)
	karalabe.DefineUint64(c, &d.Index)
}

// MarshalSSZ marshals the Deposit object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (d *DepositKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabe.Size(d))
	return buf, karalabe.EncodeToBytes(buf, d)
}

func (*DepositKaralabe) ValidateAfterDecodingSSZ() error { return nil }

// SizeSSZ returns the SSZ encoded size of the Deposit object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (d *DepositKaralabe) SizeSSZ() uint32 {
	return depositSizeKaralabe
}

// HashTreeRoot computes the Merkleization of the Deposit object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (d *DepositKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashSequential(d)
}

// UnmarshalSSZ unmarshals the Deposit object from SSZ format.
// Note: karalabe/ssz doesn't have explicit UnmarshalSSZ, we use karalabe.DecodeFromBytes
func (d *DepositKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, d)
}

// TestDepositCompatibility tests that the current Deposit implementation
// produces identical SSZ encoding/decoding results as the original karalabe/ssz implementation.
func TestDepositCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.Deposit, *DepositKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.Deposit, *DepositKaralabe) {
				return &types.Deposit{}, &DepositKaralabe{}
			},
		},
		{
			name: "typical deposit",
			setup: func() (*types.Deposit, *DepositKaralabe) {
				pubkey := crypto.BLSPubkey{1, 2, 3, 4, 5, 6, 7, 8}
				creds := types.WithdrawalCredentials{0x01} // ETH1_ADDRESS_WITHDRAWAL_PREFIX
				for i := 1; i < 12; i++ {
					creds[i] = 0x00 // padding
				}
				// Set example address bytes
				for i := 12; i < 32; i++ {
					creds[i] = byte(i)
				}
				amount := math.Gwei(32000000000) // 32 ETH
				sig := crypto.BLSSignature{9, 10, 11, 12, 13, 14, 15, 16}
				index := uint64(12345)

				current := &types.Deposit{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
					Signature:   sig,
					Index:       index,
				}
				karalabe := &DepositKaralabe{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
					Signature:   sig,
					Index:       index,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.Deposit, *DepositKaralabe) {
				var pubkey crypto.BLSPubkey
				var creds types.WithdrawalCredentials
				var sig crypto.BLSSignature

				// Fill with max values
				for i := range pubkey {
					pubkey[i] = 0xFF
				}
				for i := range creds {
					creds[i] = 0xFF
				}
				for i := range sig {
					sig[i] = 0xFF
				}

				amount := math.Gwei(^uint64(0)) // max uint64
				index := ^uint64(0)             // max uint64

				current := &types.Deposit{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
					Signature:   sig,
					Index:       index,
				}
				karalabe := &DepositKaralabe{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
					Signature:   sig,
					Index:       index,
				}
				return current, karalabe
			},
		},
		{
			name: "BLS withdrawal credentials",
			setup: func() (*types.Deposit, *DepositKaralabe) {
				pubkey := crypto.BLSPubkey{100, 101, 102, 103, 104, 105}
				creds := types.WithdrawalCredentials{0x00} // BLS_WITHDRAWAL_PREFIX
				// Fill rest with example BLS pubkey hash
				for i := 1; i < 32; i++ {
					creds[i] = byte(i * 2)
				}
				amount := math.Gwei(1000000000) // 1 ETH
				sig := crypto.BLSSignature{200, 201, 202, 203, 204, 205}
				index := uint64(999999)

				current := &types.Deposit{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
					Signature:   sig,
					Index:       index,
				}
				karalabe := &DepositKaralabe{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
					Signature:   sig,
					Index:       index,
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
			require.Equal(t, int(depositSizeKaralabe), current.SizeSSZ(), "size should match")
			require.Equal(t, uint32(depositSizeKaralabe), karalabe.SizeSSZ(), "size should match")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.Deposit{}
			err := newCurrent.UnmarshalSSZ(karalableBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &DepositKaralabe{}
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

// TestDepositCompatibilityFuzz uses fuzzing to find edge cases in SSZ compatibility
func TestDepositCompatibilityFuzz(t *testing.T) {
	// Test with random valid SSZ data
	for i := 0; i < 100; i++ {
		// Create random but valid deposit data
		var pubkey crypto.BLSPubkey
		var creds types.WithdrawalCredentials
		var sig crypto.BLSSignature

		// Use deterministic "random" data based on iteration
		for j := range pubkey {
			pubkey[j] = byte((i + j) % 256)
		}
		for j := range creds {
			creds[j] = byte((i*2 + j) % 256)
		}
		for j := range sig {
			sig[j] = byte((i*3 + j) % 256)
		}

		amount := math.Gwei(uint64(i) * 1000000000)
		index := uint64(i * 12345)

		current := &types.Deposit{
			Pubkey:      pubkey,
			Credentials: creds,
			Amount:      amount,
			Signature:   sig,
			Index:       index,
		}
		karalabe := &DepositKaralabe{
			Pubkey:      pubkey,
			Credentials: creds,
			Amount:      amount,
			Signature:   sig,
			Index:       index,
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

// TestDepositCompatibilityInvalidData tests that both implementations handle invalid data the same way
func TestDepositCompatibilityInvalidData(t *testing.T) {
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
			data: make([]byte, 100), // less than required 192 bytes
		},
		{
			name: "excess data",
			data: make([]byte, 300), // more than required 192 bytes
		},
		{
			name: "exact size but invalid content",
			data: make([]byte, 192), // correct size, all zeros
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test unmarshal with current implementation
			current := &types.Deposit{}
			currentErr := current.UnmarshalSSZ(tc.data)

			// Test unmarshal with karalabe implementation
			karalabe := &DepositKaralabe{}
			karalabelErr := karalabe.UnmarshalSSZ(tc.data)

			// Both should handle errors consistently
			if currentErr != nil && karalabelErr != nil {
				// Both errored, which is expected for invalid data
				t.Logf("Both implementations correctly rejected invalid data: current=%v, karalabe=%v", currentErr, karalabelErr)
			} else if currentErr == nil && karalabelErr == nil {
				// Both succeeded, verify they decoded to the same values
				// Convert to same type for comparison
				require.Equal(t, current.Pubkey, karalabe.Pubkey, "pubkeys should match")
				require.Equal(t, current.Credentials, karalabe.Credentials, "credentials should match")
				require.Equal(t, current.Amount, karalabe.Amount, "amounts should match")
				require.Equal(t, current.Signature, karalabe.Signature, "signatures should match")
				require.Equal(t, current.Index, karalabe.Index, "indices should match")
			} else {
				// One errored and one didn't - this would be a compatibility issue
				t.Errorf("Inconsistent error handling: current error=%v, karalabe error=%v", currentErr, karalabelErr)
			}
		})
	}
}

// TestDepositCompatibilityRoundTrip verifies that data can round-trip between implementations
func TestDepositCompatibilityRoundTrip(t *testing.T) {
	// Create a deposit with specific values
	original := &types.Deposit{
		Pubkey:      crypto.BLSPubkey{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		Credentials: types.WithdrawalCredentials{0x01, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
		Amount:      math.Gwei(32000000000),
		Signature:   crypto.BLSSignature{40, 41, 42, 43, 44, 45, 46, 47, 48, 49},
		Index:       uint64(1337),
	}

	// Marshal with current implementation
	currentBytes, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with karalabe implementation
	karalabe := &DepositKaralabe{}
	err = karalabe.UnmarshalSSZ(currentBytes)
	require.NoError(t, err)

	// Marshal with karalabe implementation
	karalableBytes, err := karalabe.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal with current implementation
	roundTrip := &types.Deposit{}
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
