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
	karalabe "github.com/karalabe/ssz"
	"github.com/stretchr/testify/require"
)

const sszConsolidationRequestSizeKaralabe = 116

// Compile-time check to ensure ConsolidationRequestKaralabe implements the necessary interfaces.
var _ karalabe.StaticObject = (*ConsolidationRequestKaralabe)(nil)

// ConsolidationRequestKaralabe is introduced in Pectra but not used by us.
// Exact copy from commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
type ConsolidationRequestKaralabe struct {
	SourceAddress common.ExecutionAddress
	SourcePubKey  crypto.BLSPubkey
	TargetPubKey  crypto.BLSPubkey
}

func (c *ConsolidationRequestKaralabe) DefineSSZ(codec *karalabe.Codec) {
	karalabe.DefineStaticBytes(codec, &c.SourceAddress)
	karalabe.DefineStaticBytes(codec, &c.SourcePubKey)
	karalabe.DefineStaticBytes(codec, &c.TargetPubKey)
}

func (c *ConsolidationRequestKaralabe) SizeSSZ() uint32 {
	return sszConsolidationRequestSizeKaralabe
}

func (c *ConsolidationRequestKaralabe) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, karalabe.Size(c))
	return buf, karalabe.EncodeToBytes(buf, c)
}

// HashTreeRoot returns the hash tree root of the ConsolidationRequest.
func (c *ConsolidationRequestKaralabe) HashTreeRoot() common.Root {
	return karalabe.HashSequential(c)
}

func (c *ConsolidationRequestKaralabe) UnmarshalSSZ(buf []byte) error {
	return karalabe.DecodeFromBytes(buf, c)
}

// TestConsolidationRequestSSZRegression ensures that the SSZ encoding for ConsolidationRequest
// remains stable and backward compatible.
func TestConsolidationRequestSSZRegression(t *testing.T) {
	testCases := []struct {
		name        string
		request     *types.ConsolidationRequest
		expectedSSZ []byte // Pre-computed expected SSZ encoding
	}{
		{
			name: "zero values",
			request: &types.ConsolidationRequest{
				SourceAddress: common.ExecutionAddress{},
				SourcePubKey:  crypto.BLSPubkey{},
				TargetPubKey:  crypto.BLSPubkey{},
			},
			// Expected SSZ: 20 zero bytes (address) + 48 zero bytes (source) + 48 zero bytes (target) = 116 bytes
			expectedSSZ: make([]byte, 116),
		},
		{
			name: "typical consolidation request",
			request: &types.ConsolidationRequest{
				SourceAddress: common.ExecutionAddress{
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
					0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
					0x11, 0x22, 0x33, 0x44,
				},
				SourcePubKey: crypto.BLSPubkey{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				},
				TargetPubKey: crypto.BLSPubkey{
					0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8,
					0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb0,
					0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
					0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf, 0xc0,
					0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8,
					0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf, 0xd0,
				},
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 116)
				// SourceAddress
				copy(ssz[0:20], []byte{
					0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
					0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x00,
					0x11, 0x22, 0x33, 0x44,
				})
				// SourcePubKey
				copy(ssz[20:68], []byte{
					0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08,
					0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10,
					0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
					0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
					0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
					0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
				})
				// TargetPubKey
				copy(ssz[68:116], []byte{
					0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6, 0xa7, 0xa8,
					0xa9, 0xaa, 0xab, 0xac, 0xad, 0xae, 0xaf, 0xb0,
					0xb1, 0xb2, 0xb3, 0xb4, 0xb5, 0xb6, 0xb7, 0xb8,
					0xb9, 0xba, 0xbb, 0xbc, 0xbd, 0xbe, 0xbf, 0xc0,
					0xc1, 0xc2, 0xc3, 0xc4, 0xc5, 0xc6, 0xc7, 0xc8,
					0xc9, 0xca, 0xcb, 0xcc, 0xcd, 0xce, 0xcf, 0xd0,
				})
				return ssz
			}(),
		},
		{
			name: "all max values",
			request: &types.ConsolidationRequest{
				SourceAddress: common.ExecutionAddress{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff,
				},
				SourcePubKey: crypto.BLSPubkey{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
				TargetPubKey: crypto.BLSPubkey{
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
					0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				},
			},
			expectedSSZ: func() []byte {
				ssz := make([]byte, 116)
				for i := range ssz {
					ssz[i] = 0xff
				}
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
			require.Equal(t, 116, tc.request.SizeSSZ(), "size should be 116 bytes")

			// Test Unmarshal
			unmarshaled := &types.ConsolidationRequest{}
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

// TestConsolidationRequestCompatibility tests that current and karalabe implementations produce identical results
func TestConsolidationRequestCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (*types.ConsolidationRequest, *ConsolidationRequestKaralabe)
	}{
		{
			name: "zero values",
			setup: func() (*types.ConsolidationRequest, *ConsolidationRequestKaralabe) {
				return &types.ConsolidationRequest{},
					&ConsolidationRequestKaralabe{}
			},
		},
		{
			name: "typical consolidation request",
			setup: func() (*types.ConsolidationRequest, *ConsolidationRequestKaralabe) {
				sourceAddr := common.ExecutionAddress{1, 2, 3, 4, 5}
				sourcePubkey := crypto.BLSPubkey{6, 7, 8, 9, 10}
				targetPubkey := crypto.BLSPubkey{11, 12, 13, 14, 15}

				current := &types.ConsolidationRequest{
					SourceAddress: sourceAddr,
					SourcePubKey:  sourcePubkey,
					TargetPubKey:  targetPubkey,
				}
				karalabe := &ConsolidationRequestKaralabe{
					SourceAddress: sourceAddr,
					SourcePubKey:  sourcePubkey,
					TargetPubKey:  targetPubkey,
				}
				return current, karalabe
			},
		},
		{
			name: "maximum values",
			setup: func() (*types.ConsolidationRequest, *ConsolidationRequestKaralabe) {
				var sourceAddr common.ExecutionAddress
				var sourcePubkey, targetPubkey crypto.BLSPubkey
				for i := range sourceAddr {
					sourceAddr[i] = 0xFF
				}
				for i := range sourcePubkey {
					sourcePubkey[i] = 0xFF
				}
				for i := range targetPubkey {
					targetPubkey[i] = 0xFF
				}

				current := &types.ConsolidationRequest{
					SourceAddress: sourceAddr,
					SourcePubKey:  sourcePubkey,
					TargetPubKey:  targetPubkey,
				}
				karalabe := &ConsolidationRequestKaralabe{
					SourceAddress: sourceAddr,
					SourcePubKey:  sourcePubkey,
					TargetPubKey:  targetPubkey,
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
			require.Equal(t, 116, current.SizeSSZ(), "current size should be 116")
			require.Equal(t, uint32(116), karalabe.SizeSSZ(), "karalabe size should be 116")

			// Test Unmarshal with karalabe marshaled data
			newCurrent := &types.ConsolidationRequest{}
			err := newCurrent.UnmarshalSSZ(karalabeBytes)
			require.NoError(t, err, "unmarshal karalabe data into current should not error")
			require.Equal(t, current, newCurrent, "unmarshaled current should match original")

			// Test Unmarshal with current marshaled data
			newKaralabe := &ConsolidationRequestKaralabe{}
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

// TestConsolidationRequestSSZInvalidData tests error handling for invalid SSZ data
func TestConsolidationRequestSSZInvalidData(t *testing.T) {
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
			data:          make([]byte, 50), // less than required 116 bytes
			expectedError: "incorrect size",
		},
		{
			name:          "excess data",
			data:          make([]byte, 150), // more than required 116 bytes
			expectedError: "incorrect size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			request := &types.ConsolidationRequest{}
			err := request.UnmarshalSSZ(tc.data)
			require.Error(t, err, "UnmarshalSSZ should error on invalid data")
			require.Contains(t, err.Error(), tc.expectedError, "error should contain expected message")
		})
	}
}

// TestConsolidationRequestSSZRoundTrip tests round-trip encoding/decoding with various data patterns
func TestConsolidationRequestSSZRoundTrip(t *testing.T) {
	patterns := []struct {
		name  string
		setup func() *types.ConsolidationRequest
	}{
		{
			name: "all zeros",
			setup: func() *types.ConsolidationRequest {
				return &types.ConsolidationRequest{}
			},
		},
		{
			name: "incremental pattern",
			setup: func() *types.ConsolidationRequest {
				var addr common.ExecutionAddress
				for i := range addr {
					addr[i] = byte(i)
				}
				var source crypto.BLSPubkey
				for i := range source {
					source[i] = byte(i + 20)
				}
				var target crypto.BLSPubkey
				for i := range target {
					target[i] = byte(i + 68)
				}
				return &types.ConsolidationRequest{
					SourceAddress: addr,
					SourcePubKey:  source,
					TargetPubKey:  target,
				}
			},
		},
		{
			name: "specific values",
			setup: func() *types.ConsolidationRequest {
				return &types.ConsolidationRequest{
					SourceAddress: common.ExecutionAddress{
						0xde, 0xad, 0xbe, 0xef, 0xca, 0xfe, 0xba, 0xbe,
						0xfa, 0xce, 0xdb, 0xad, 0xde, 0xed, 0xbe, 0xef,
						0x12, 0x34, 0x56, 0x78,
					},
					SourcePubKey: crypto.BLSPubkey{
						0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22,
						0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0x00,
						0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0,
						0x0f, 0x1e, 0x2d, 0x3c, 0x4b, 0x5a, 0x69, 0x78,
						0x87, 0x96, 0xa5, 0xb4, 0xc3, 0xd2, 0xe1, 0xf0,
						0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa, 0x99, 0x88,
					},
					TargetPubKey: crypto.BLSPubkey{
						0x11, 0x11, 0x11, 0x11, 0x22, 0x22, 0x22, 0x22,
						0x33, 0x33, 0x33, 0x33, 0x44, 0x44, 0x44, 0x44,
						0x55, 0x55, 0x55, 0x55, 0x66, 0x66, 0x66, 0x66,
						0x77, 0x77, 0x77, 0x77, 0x88, 0x88, 0x88, 0x88,
						0x99, 0x99, 0x99, 0x99, 0xaa, 0xaa, 0xaa, 0xaa,
						0xbb, 0xbb, 0xbb, 0xbb, 0xcc, 0xcc, 0xcc, 0xcc,
					},
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
			decoded := &types.ConsolidationRequest{}
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

// TestConsolidationRequestsSSZ tests the ConsolidationRequests slice type SSZ encoding
func TestConsolidationRequestsSSZ(t *testing.T) {
	// Create test consolidation requests
	req1 := &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		SourcePubKey:  crypto.BLSPubkey{1},
		TargetPubKey:  crypto.BLSPubkey{2},
	}
	req2 := &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40},
		SourcePubKey:  crypto.BLSPubkey{3},
		TargetPubKey:  crypto.BLSPubkey{4},
	}

	requests := types.ConsolidationRequests{req1, req2}

	// Test Marshal
	data, err := requests.MarshalSSZ()
	require.NoError(t, err, "MarshalSSZ should not error")
	require.Equal(t, 232, len(data), "marshaled data should be 232 bytes (2 * 116)")

	// Test Unmarshal
	decoded, err := types.DecodeConsolidationRequests(data)
	require.NoError(t, err, "DecodeConsolidationRequests should not error")
	require.Equal(t, len(requests), len(decoded), "decoded should have same length")
	require.Equal(t, requests[0], decoded[0], "first request should match")
	require.Equal(t, requests[1], decoded[1], "second request should match")
}
