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

// Compile-time assertions to ensure DepositMessageKaralabe implements necessary interfaces.
var _ ssz.StaticObject = (*DepositMessageKaralabe)(nil)

// DepositMessageKaralabe represents a deposit message as defined in the Ethereum 2.0
// specification - exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
type DepositMessageKaralabe struct {
	// Public key of the validator specified in the deposit.
	Pubkey crypto.BLSPubkey `json:"pubkey"`
	// A staking credentials with
	// 1 byte prefix + 11 bytes padding + 20 bytes address = 32 bytes.
	Credentials types.WithdrawalCredentials `json:"credentials"`
	// Deposit amount in gwei.
	Amount math.Gwei `json:"amount"`
}

// SizeSSZ returns the size of the DepositMessage object in SSZ encoding.
// Note: karalabe/ssz has different signatures for StaticObject vs DynamicObject
func (*DepositMessageKaralabe) SizeSSZ() uint32 {
	//nolint:mnd // 48 + 32 + 8 = 88.
	return 88
}

// DefineSSZ defines the SSZ encoding for the DepositMessage object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (dm *DepositMessageKaralabe) DefineSSZ(codec *ssz.Codec) {
	ssz.DefineStaticBytes(codec, &dm.Pubkey)
	ssz.DefineStaticBytes(codec, &dm.Credentials)
	ssz.DefineUint64(codec, &dm.Amount)
}

// HashTreeRoot computes the SSZ hash tree root of the DepositMessage object.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (dm *DepositMessageKaralabe) HashTreeRoot() common.Root {
	return ssz.HashSequential(dm)
}

// MarshalSSZ marshals the DepositMessage object to SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (dm *DepositMessageKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, ssz.Size(dm))
	return buf, ssz.EncodeToBytes(buf, dm)
}

// UnmarshalSSZ unmarshals the DepositMessage object from SSZ format.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
func (dm *DepositMessageKaralabe) UnmarshalSSZ(buf []byte) error {
	return ssz.DecodeFromBytes(buf, dm)
}

// TestDepositMessageSSZRegression ensures that the SSZ encoding for DepositMessage
// remains stable and backward compatible.
func TestDepositMessageSSZRegression(t *testing.T) {
	testCases := []struct {
		name        string
		depositMsg  *types.DepositMessage
		expectedSSZ []byte // Pre-computed expected SSZ encoding
	}{
		{
			name: "zero values",
			depositMsg: &types.DepositMessage{
				Pubkey:      crypto.BLSPubkey{},
				Credentials: types.WithdrawalCredentials{},
				Amount:      math.Gwei(0),
			},
			// Expected SSZ: 48 zero bytes (pubkey) + 32 zero bytes (credentials) + 8 zero bytes (amount) = 88 bytes
			expectedSSZ: make([]byte, 88),
		},
		{
			name: "typical deposit message",
			depositMsg: &types.DepositMessage{
				Pubkey: crypto.BLSPubkey{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				},
				Credentials: types.WithdrawalCredentials{
					0x00,                                                             // ETH1 withdrawal credentials prefix
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 11 bytes padding
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, // 20 bytes address
				},
				Amount: math.Gwei(32_000_000_000), // 32 ETH in Gwei
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 88)
				// Pubkey
				copy(ssz[0:48], []byte{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				})
				// Credentials
				copy(ssz[48:80], []byte{
					0x00,                                                             // ETH1 withdrawal credentials prefix
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // 11 bytes padding
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00, 0x11, 0x22, 0x33, 0x44, // 20 bytes address
				})
				// Amount (32_000_000_000 Gwei in little-endian)
				ssz[80] = 0x00
				ssz[81] = 0x40
				ssz[82] = 0x59
				ssz[83] = 0x73
				ssz[84] = 0x07
				ssz[85] = 0x00
				ssz[86] = 0x00
				ssz[87] = 0x00
				return ssz
			}(),
		},
		{
			name: "maximum amount",
			depositMsg: &types.DepositMessage{
				Pubkey: crypto.BLSPubkey{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				Credentials: types.WithdrawalCredentials{
					0x01, // BLS withdrawal credentials prefix
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				Amount: math.Gwei(^uint64(0)), // Max uint64
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 88)
				// Fill pubkey with 0xff
				for i := 0; i < 48; i++ {
					ssz[i] = 0xff
				}
				// Credentials
				ssz[48] = 0x01 // BLS prefix
				for i := 49; i < 80; i++ {
					ssz[i] = 0xff
				}
				// Max amount (all 0xff in little-endian)
				for i := 80; i < 88; i++ {
					ssz[i] = 0xff
				}
				return ssz
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Marshal
			actualSSZ, err := tc.depositMsg.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ, "SSZ encoding should match expected")

			// Test Size
			require.Equal(t, 88, tc.depositMsg.SizeSSZ(), "size should be 88 bytes")

			// Test Unmarshal
			unmarshaled := &types.DepositMessage{}
			err = unmarshaled.UnmarshalSSZ(tc.expectedSSZ)
			require.NoError(t, err, "UnmarshalSSZ should not error")
			require.Equal(t, tc.depositMsg, unmarshaled, "unmarshaled object should match original")

			// Test MarshalSSZTo
			buf := make([]byte, 0, tc.depositMsg.SizeSSZ())
			actualSSZ2, err := tc.depositMsg.MarshalSSZTo(buf)
			require.NoError(t, err, "MarshalSSZTo should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ2, "MarshalSSZTo should produce same output")

			// Test HashTreeRoot consistency
			root1, err := tc.depositMsg.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot should not error")
			root2, err := unmarshaled.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot of unmarshaled should not error")
			require.Equal(t, root1, root2, "hash tree roots should match")
		})
	}
}

// TestDepositMessageCompatibility tests that current and karalabe implementations produce identical results
func TestDepositMessageCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.DepositMessage, *DepositMessageKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.DepositMessage, *DepositMessageKaralabe) {
				return &types.DepositMessage{},
					&DepositMessageKaralabe{}
			},
		},
		{
			name: "typical deposit message",
			setup: func() (*types.DepositMessage, *DepositMessageKaralabe) {
				pubkey := crypto.BLSPubkey{1, 2, 3, 4, 5, 6}
				creds := types.WithdrawalCredentials{7, 8, 9, 10}
				amount := math.Gwei(32000000000) // 32 ETH

				current := &types.DepositMessage{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
				}
				karalabe := &DepositMessageKaralabe{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.DepositMessage, *DepositMessageKaralabe) {
				var pubkey crypto.BLSPubkey
				var creds types.WithdrawalCredentials
				for i := range pubkey {
					pubkey[i] = 0xFF
				}
				for i := range creds {
					creds[i] = 0xFF
				}
				amount := math.Gwei(^uint64(0))

				current := &types.DepositMessage{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
				}
				karalabe := &DepositMessageKaralabe{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      amount,
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
			require.Equal(t, 88, current.SizeSSZ(), "current size should be 88")
			require.Equal(t, uint32(88), karalabe.SizeSSZ(), "karalabe size should be 88")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.DepositMessage{}
			err := newCurrent.UnmarshalSSZ(karalabeBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &DepositMessageKaralabe{}
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

// TestDepositMessageSSZInvalidData tests error handling for invalid SSZ data
func TestDepositMessageSSZInvalidData(t *testing.T) {
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
			data:          make([]byte, 50), // less than required 88 bytes
			expectedError: "incorrect size",
		},
		{
			name:          "excess data",
			data:          make([]byte, 100), // more than required 88 bytes
			expectedError: "incorrect size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			depositMsg := &types.DepositMessage{}
			err := depositMsg.UnmarshalSSZ(tc.data)
			require.Error(t, err, "UnmarshalSSZ should error on invalid data")
			require.Contains(t, err.Error(), tc.expectedError, "error should contain expected message")
		})
	}
}

// TestDepositMessageSSZRoundTrip tests round-trip encoding/decoding with various data patterns
func TestDepositMessageSSZRoundTrip(t *testing.T) {
	// Test with various patterns
	patterns := []struct {
		name  string
		setup func() *types.DepositMessage
	}{
		{
			name: "all zeros",
			setup: func() *types.DepositMessage {
				return &types.DepositMessage{}
			},
		},
		{
			name: "incremental pattern",
			setup: func() *types.DepositMessage {
				var pubkey crypto.BLSPubkey
				for i := range pubkey {
					pubkey[i] = byte(i)
				}
				var creds types.WithdrawalCredentials
				for i := range creds {
					creds[i] = byte(i * 2)
				}
				return &types.DepositMessage{
					Pubkey:      pubkey,
					Credentials: creds,
					Amount:      math.Gwei(12345678),
				}
			},
		},
		{
			name: "specific values",
			setup: func() *types.DepositMessage {
				return &types.DepositMessage{
					Pubkey: crypto.BLSPubkey{
						0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
						0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
						0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
						0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
						0x87, 0x96, 0xa5, 0xb4, 0xc3, 0xd2, 0xe1, 0xf0,
						0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
					},
					Credentials: types.WithdrawalCredentials{
						0x00,                                                             // ETH1 prefix
						0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // padding
						0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe, 0xfa, 0xce, 0xdb, 0xad, 0xde, 0xed, 0xbe, 0xef, 0x12, 0x34, 0x56, 0x78,
					},
					Amount: math.Gwei(1_000_000_000), // 1 ETH
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
			decoded := &types.DepositMessage{}
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
