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
	"github.com/stretchr/testify/require"
)

// TestWithdrawalRequestSSZRegression ensures that the SSZ encoding for WithdrawalRequest
// remains stable and backward compatible.
func TestWithdrawalRequestSSZRegression(t *testing.T) {
	testCases := []struct {
		name        string
		request     *types.WithdrawalRequest
		expectedSSZ []byte // Pre-computed expected SSZ encoding
	}{
		{
			name: "zero values",
			request: &types.WithdrawalRequest{
				SourceAddress:   common.ExecutionAddress{},
				ValidatorPubKey: crypto.BLSPubkey{},
				Amount:          math.Gwei(0),
			},
			// Expected SSZ: 20 zero bytes (address) + 48 zero bytes (pubkey) + 8 zero bytes (amount) = 76 bytes
			expectedSSZ: make([]byte, 76),
		},
		{
			name: "typical withdrawal request",
			request: &types.WithdrawalRequest{
				SourceAddress: common.ExecutionAddress{
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
					0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
					0x11, 0x22, 0x33, 0x44,
				},
				ValidatorPubKey: crypto.BLSPubkey{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				},
				Amount: math.Gwei(32_000_000_000), // 32 ETH in Gwei
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 76)
				// SourceAddress
				copy(ssz[0:20], []byte{
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
					0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
					0x11, 0x22, 0x33, 0x44,
				})
				// ValidatorPubKey
				copy(ssz[20:68], []byte{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				})
				// Amount (32_000_000_000 Gwei in little-endian)
				ssz[68] = 0x00
				ssz[69] = 0x40
				ssz[70] = 0x59
				ssz[71] = 0x73
				ssz[72] = 0x07
				ssz[73] = 0x00
				ssz[74] = 0x00
				ssz[75] = 0x00
				return ssz
			}(),
		},
		{
			name: "maximum amount",
			request: &types.WithdrawalRequest{
				SourceAddress: common.ExecutionAddress{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff,
				},
				ValidatorPubKey: crypto.BLSPubkey{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				Amount: math.Gwei(^uint64(0)), // Max uint64
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 76)
				for i := range ssz {
					ssz[i] = 0xff
				}
				return ssz
			}(),
		},
		{
			name: "partial withdrawal",
			request: &types.WithdrawalRequest{
				SourceAddress: common.ExecutionAddress{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0xfa, 0xce, 0xdb, 0xad, 0xde, 0xed, 0xbe, 0xef,
					0x12, 0x34, 0x56, 0x78,
				},
				ValidatorPubKey: crypto.BLSPubkey{
					0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
					0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
					0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
					0x87, 0x96, 0xa5, 0xb4, 0xc3, 0xd2, 0xe1, 0xf0,
					0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				},
				Amount: math.Gwei(1_000_000_000), // 1 ETH
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 76)
				// SourceAddress
				copy(ssz[0:20], []byte{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0xfa, 0xce, 0xdb, 0xad, 0xde, 0xed, 0xbe, 0xef,
					0x12, 0x34, 0x56, 0x78,
				})
				// ValidatorPubKey
				copy(ssz[20:68], []byte{
					0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
					0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
					0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
					0x87, 0x96, 0xa5, 0xb4, 0xc3, 0xd2, 0xe1, 0xf0,
					0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				})
				// Amount (1_000_000_000 = 0x3B9ACA00 in little-endian)
				ssz[68] = 0x00
				ssz[69] = 0xca
				ssz[70] = 0x9a
				ssz[71] = 0x3b
				ssz[72] = 0x00
				ssz[73] = 0x00
				ssz[74] = 0x00
				ssz[75] = 0x00
				return ssz
			}(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test Marshal
			actualSSZ, err := tc.request.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ, "SSZ encoding should match expected")

			// Test Size
			require.Equal(t, 76, tc.request.SizeSSZ(), "size should be 76 bytes")

			// Test Unmarshal
			unmarshaled := &types.WithdrawalRequest{}
			err = unmarshaled.UnmarshalSSZ(tc.expectedSSZ)
			require.NoError(t, err, "UnmarshalSSZ should not error")
			require.Equal(t, tc.request, unmarshaled, "unmarshaled object should match original")

			// Test MarshalSSZTo
			buf := make([]byte, 0, tc.request.SizeSSZ())
			actualSSZ2, err := tc.request.MarshalSSZTo(buf)
			require.NoError(t, err, "MarshalSSZTo should not error")
			require.Equal(t, tc.expectedSSZ, actualSSZ2, "MarshalSSZTo should produce same output")

			// Test HashTreeRoot consistency
			root1, err := tc.request.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot should not error")
			root2, err := unmarshaled.HashTreeRoot()
			require.NoError(t, err, "HashTreeRoot of unmarshaled should not error")
			require.Equal(t, root1, root2, "hash tree roots should match")
		})
	}
}

// TestWithdrawalRequestSSZInvalidData tests error handling for invalid SSZ data
func TestWithdrawalRequestSSZInvalidData(t *testing.T) {
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
			data:          make([]byte, 50), // less than required 76 bytes
			expectedError: "incorrect size",
		},
		{
			name:          "excess data",
			data:          make([]byte, 100), // more than required 76 bytes
			expectedError: "incorrect size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &types.WithdrawalRequest{}
			err := request.UnmarshalSSZ(tc.data)
			require.Error(t, err, "UnmarshalSSZ should error on invalid data")
			require.Contains(t, err.Error(), tc.expectedError, "error should contain expected message")
		})
	}
}

// TestWithdrawalRequestSSZRoundTrip tests round-trip encoding/decoding with various data patterns
func TestWithdrawalRequestSSZRoundTrip(t *testing.T) {
	patterns := []struct {
		name  string
		setup func() *types.WithdrawalRequest
	}{
		{
			name: "all zeros",
			setup: func() *types.WithdrawalRequest {
				return &types.WithdrawalRequest{}
			},
		},
		{
			name: "incremental pattern",
			setup: func() *types.WithdrawalRequest {
				var addr common.ExecutionAddress
				for i := range addr {
					addr[i] = byte(i)
				}
				var pubkey crypto.BLSPubkey
				for i := range pubkey {
					pubkey[i] = byte(i + 20)
				}
				return &types.WithdrawalRequest{
					SourceAddress:   addr,
					ValidatorPubKey: pubkey,
					Amount:          math.Gwei(12345678),
				}
			},
		},
		{
			name: "specific values",
			setup: func() *types.WithdrawalRequest {
				return &types.WithdrawalRequest{
					SourceAddress: common.ExecutionAddress{
						0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
						0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
						0x87, 0x96, 0xa5, 0xb4,
					},
					ValidatorPubKey: crypto.BLSPubkey{
						0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11,
						0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99,
						0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11,
						0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99,
						0x88, 0x77, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11,
						0x00, 0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99,
					},
					Amount: math.Gwei(1337_000_000_000), // 1337 ETH
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
			decoded := &types.WithdrawalRequest{}
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

// TestWithdrawalRequestsSSZ tests the WithdrawalRequests slice type SSZ encoding
func TestWithdrawalRequestsSSZ(t *testing.T) {
	// Create test withdrawal requests
	req1 := &types.WithdrawalRequest{
		SourceAddress:   common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		ValidatorPubKey: crypto.BLSPubkey{1},
		Amount:          math.Gwei(1_000_000_000),
	}
	req2 := &types.WithdrawalRequest{
		SourceAddress:   common.ExecutionAddress{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		ValidatorPubKey: crypto.BLSPubkey{2},
		Amount:          math.Gwei(2_000_000_000),
	}

	requests := types.WithdrawalRequests{req1, req2}

	// Test Marshal
	data, err := requests.MarshalSSZ()
	require.NoError(t, err, "MarshalSSZ should not error")
	require.Equal(t, 152, len(data), "marshaled data should be 152 bytes (2 * 76)")

	// Test Unmarshal
	decoded, err := types.DecodeWithdrawalRequests(data)
	require.NoError(t, err, "DecodeWithdrawalRequests should not error")
	require.Equal(t, len(requests), len(decoded), "decoded should have same length")
	require.Equal(t, requests[0], decoded[0], "first request should match")
	require.Equal(t, requests[1], decoded[1], "second request should match")
}

// TestWithdrawalRequestsSSZInvalidData tests error handling for invalid WithdrawalRequests data
func TestWithdrawalRequestsSSZInvalidData(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		expectedError string
	}{
		{
			name:          "data not multiple of withdrawal request size",
			data:          make([]byte, 100), // not a multiple of 76
			expectedError: "invalid data length",
		},
		{
			name:          "too small",
			data:          make([]byte, 50), // less than one withdrawal request
			expectedError: "invalid withdrawal requests SSZ size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := types.DecodeWithdrawalRequests(tc.data)
			require.Error(t, err, "DecodeWithdrawalRequests should error on invalid data")
			require.Contains(t, err.Error(), tc.expectedError, "error should contain expected message")
		})
	}
}

// TestWithdrawalRequestCompatibility tests that current and karalabe implementations produce identical results
func TestWithdrawalRequestCompatibility(t *testing.T) {
	testCases := []struct {
		name string
		setup func() (*types.WithdrawalRequest, *WithdrawalRequestKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.WithdrawalRequest, *WithdrawalRequestKaralabe) {
				return &types.WithdrawalRequest{},
					&WithdrawalRequestKaralabe{}
			},
		},
		{
			name: "typical withdrawal request",
			setup: func() (*types.WithdrawalRequest, *WithdrawalRequestKaralabe) {
				sourceAddr := common.ExecutionAddress{
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
					0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
					0x11, 0x22, 0x33, 0x44,
				}
				pubkey := crypto.BLSPubkey{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				}
				amount := math.Gwei(32000000000) // 32 ETH
				
				current := &types.WithdrawalRequest{
					SourceAddress:   sourceAddr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}
				karalabe := &WithdrawalRequestKaralabe{
					SourceAddress:   sourceAddr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.WithdrawalRequest, *WithdrawalRequestKaralabe) {
				var sourceAddr common.ExecutionAddress
				var pubkey crypto.BLSPubkey
				for i := range sourceAddr {
					sourceAddr[i] = 0xFF
				}
				for i := range pubkey {
					pubkey[i] = 0xFF
				}
				amount := math.Gwei(^uint64(0))
				
				current := &types.WithdrawalRequest{
					SourceAddress:   sourceAddr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}
				karalabe := &WithdrawalRequestKaralabe{
					SourceAddress:   sourceAddr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}
				return current, karalabe
			},
		},
		{
			name: "partial withdrawal",
			setup: func() (*types.WithdrawalRequest, *WithdrawalRequestKaralabe) {
				sourceAddr := common.ExecutionAddress{
					0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
					0xfa, 0xce, 0xdb, 0xad, 0xde, 0xed, 0xbe, 0xef,
					0x12, 0x34, 0x56, 0x78,
				}
				pubkey := crypto.BLSPubkey{
					0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
					0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
					0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
					0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
					0x87, 0x96, 0xa5, 0xb4, 0xc3, 0xd2, 0xe1, 0xf0,
					0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
				}
				amount := math.Gwei(1_000_000_000) // 1 ETH
				
				current := &types.WithdrawalRequest{
					SourceAddress:   sourceAddr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
				}
				karalabe := &WithdrawalRequestKaralabe{
					SourceAddress:   sourceAddr,
					ValidatorPubKey: pubkey,
					Amount:          amount,
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
			require.Equal(t, 76, current.SizeSSZ(), "current size should be 76")
			require.Equal(t, uint32(76), karalabe.SizeSSZ(), "karalabe size should be 76")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.WithdrawalRequest{}
			err := newCurrent.UnmarshalSSZ(karalabeBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &WithdrawalRequestKaralabe{}
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