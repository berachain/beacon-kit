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

// NOTE: ConsolidationRequests is a slice type that uses EIP-7685 encoding format.
// While the individual ConsolidationRequest elements were using karalabe/ssz,
// the slice encoding follows a different format specified by EIP-7685. These tests
// focus on the encoding/decoding functions rather than direct SSZ compatibility.
// The individual element type (ConsolidationRequest) has its own compatibility test.

//go:build test

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/stretchr/testify/require"
)

// TestConsolidationRequestsCompatibility tests that ConsolidationRequests slice type
// behaves correctly with EIP-7685 encoding/decoding.
func TestConsolidationRequestsCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() types.ConsolidationRequests
	}{
		// Note: Empty slices are handled differently - the Decode function
		// requires at least one element due to the validation check.
		{
			name: "single consolidation request",
			setup: func() types.ConsolidationRequests {
				addr := common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				srcPubkey := crypto.BLSPubkey{20, 21, 22, 23, 24, 25, 26, 27}
				tgtPubkey := crypto.BLSPubkey{30, 31, 32, 33, 34, 35, 36, 37}

				return types.ConsolidationRequests{
					&types.ConsolidationRequest{
						SourceAddress: addr,
						SourcePubKey:  srcPubkey,
						TargetPubKey:  tgtPubkey,
					},
				}
			},
		},
		{
			name: "multiple consolidation requests",
			setup: func() types.ConsolidationRequests {
				consolidations := make(types.ConsolidationRequests, 2) // Max is 2

				for i := 0; i < 2; i++ {
					addr := common.ExecutionAddress{}
					srcPubkey := crypto.BLSPubkey{}
					tgtPubkey := crypto.BLSPubkey{}

					// Fill with unique data for each consolidation
					for j := range addr {
						addr[j] = byte(i*20 + j)
					}
					for j := range srcPubkey {
						srcPubkey[j] = byte(i*48 + j)
					}
					for j := range tgtPubkey {
						tgtPubkey[j] = byte(i*48 + j + 100)
					}

					consolidations[i] = &types.ConsolidationRequest{
						SourceAddress: addr,
						SourcePubKey:  srcPubkey,
						TargetPubKey:  tgtPubkey,
					}
				}

				return consolidations
			},
		},
		{
			name: "consolidations with different patterns",
			setup: func() types.ConsolidationRequests {
				return types.ConsolidationRequests{
					{
						SourceAddress: common.ExecutionAddress{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
						SourcePubKey:  crypto.BLSPubkey{0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF},
						TargetPubKey:  crypto.BLSPubkey{0x11, 0x22, 0x33, 0x44, 0x55, 0x66},
					},
					{
						SourceAddress: common.ExecutionAddress{0x00, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33},
						SourcePubKey:  crypto.BLSPubkey{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08},
						TargetPubKey:  crypto.BLSPubkey{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80},
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			consolidations := tc.setup()

			// Test Marshal
			marshaled, err := consolidations.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")

			// Test Decode (which is the unmarshal operation for this type)
			decoded, err := types.DecodeConsolidationRequests(marshaled)
			require.NoError(t, err, "DecodeConsolidationRequests should not error")

			// Verify lengths match
			require.Equal(t, len(consolidations), len(decoded), "decoded length should match original")

			// Verify each consolidation matches
			for i := range consolidations {
				require.Equal(t, consolidations[i].SourceAddress, decoded[i].SourceAddress, "source address at index %d should match", i)
				require.Equal(t, consolidations[i].SourcePubKey, decoded[i].SourcePubKey, "source pubkey at index %d should match", i)
				require.Equal(t, consolidations[i].TargetPubKey, decoded[i].TargetPubKey, "target pubkey at index %d should match", i)
			}

			// Test ValidateAfterDecodingSSZ
			err = decoded.ValidateAfterDecodingSSZ()
			require.NoError(t, err, "ValidateAfterDecodingSSZ should not error for valid data")
		})
	}
}

// TestConsolidationRequestsEmptySlice tests empty slice handling separately
func TestConsolidationRequestsEmptySlice(t *testing.T) {
	// Empty slice marshals to empty byte array
	empty := types.ConsolidationRequests{}
	marshaled, err := empty.MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, 0, len(marshaled), "empty slice should marshal to empty bytes")

	// But decoding empty bytes fails due to validation
	_, err = types.DecodeConsolidationRequests(marshaled)
	require.Error(t, err, "decoding empty bytes should error")
	require.Contains(t, err.Error(), "invalid consolidation requests SSZ size", "error should indicate size issue")
}

// TestConsolidationRequestsEIP7685Encoding tests the specific EIP-7685 encoding format
func TestConsolidationRequestsEIP7685Encoding(t *testing.T) {
	// Create a simple consolidation request
	consolidation := &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		SourcePubKey:  crypto.BLSPubkey{1, 2, 3},
		TargetPubKey:  crypto.BLSPubkey{4, 5, 6},
	}

	consolidations := types.ConsolidationRequests{consolidation}

	// Marshal using EIP-7685 format
	marshaled, err := consolidations.MarshalSSZ()
	require.NoError(t, err)

	// Verify the marshaled data is exactly 116 bytes (one consolidation)
	require.Equal(t, 116, len(marshaled), "single consolidation should be 116 bytes")

	// Test with maximum consolidations (2)
	consolidations = types.ConsolidationRequests{consolidation, consolidation}
	marshaled, err = consolidations.MarshalSSZ()
	require.NoError(t, err)

	// Verify the marshaled data is exactly 232 bytes (two consolidations)
	require.Equal(t, 232, len(marshaled), "two consolidations should be 232 bytes")
}

// TestConsolidationRequestsValidation tests the validation logic
func TestConsolidationRequestsValidation(t *testing.T) {
	// Create maximum allowed consolidations
	maxConsolidations := make(types.ConsolidationRequests, 2) // MaxConsolidationRequestsPerPayload
	for i := range maxConsolidations {
		addr := common.ExecutionAddress{}
		addr[0] = byte(i)

		maxConsolidations[i] = &types.ConsolidationRequest{
			SourceAddress: addr,
			SourcePubKey:  crypto.BLSPubkey{byte(i)},
			TargetPubKey:  crypto.BLSPubkey{byte(i + 100)},
		}
	}

	// This should validate successfully
	err := maxConsolidations.ValidateAfterDecodingSSZ()
	require.NoError(t, err, "max consolidations should validate successfully")

	// Create one more than allowed
	tooManyConsolidations := append(maxConsolidations, &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{255},
		SourcePubKey:  crypto.BLSPubkey{255},
		TargetPubKey:  crypto.BLSPubkey{254},
	})

	// This should fail validation
	err = tooManyConsolidations.ValidateAfterDecodingSSZ()
	require.Error(t, err, "too many consolidations should fail validation")
	require.Contains(t, err.Error(), "invalid number of consolidation requests", "error message should indicate too many consolidations")
}

// TestConsolidationRequestsDecodeErrors tests error cases in decoding
func TestConsolidationRequestsDecodeErrors(t *testing.T) {
	testCases := []struct {
		name          string
		data          []byte
		expectError   bool
		errorContains string
	}{
		{
			name:          "empty data",
			data:          []byte{},
			expectError:   true,
			errorContains: "invalid consolidation requests SSZ size",
		},
		{
			name:          "partial consolidation data",
			data:          make([]byte, 50), // Less than 116 bytes
			expectError:   true,
			errorContains: "invalid consolidation requests SSZ size",
		},
		{
			name:          "data not multiple of consolidation size",
			data:          make([]byte, 200), // Not a multiple of 116
			expectError:   true,
			errorContains: "is not a multiple of consolidation request size",
		},
		{
			name:        "valid single consolidation",
			data:        make([]byte, 116), // Exactly one consolidation
			expectError: false,
		},
		{
			name:        "valid two consolidations",
			data:        make([]byte, 232), // Exactly two consolidations
			expectError: false,
		},
		{
			name:          "exceeds maximum size",
			data:          make([]byte, 116*3), // 3 consolidations (max is 2)
			expectError:   true,
			errorContains: "invalid consolidation requests SSZ size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := types.DecodeConsolidationRequests(tc.data)

			if tc.expectError {
				require.Error(t, err, "should error")
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains, "error message should contain expected text")
				}
			} else {
				require.NoError(t, err, "should not error")
				require.NotNil(t, decoded, "decoded should not be nil")
				// Verify we got the expected number of consolidations
				expectedCount := len(tc.data) / 116
				require.Equal(t, expectedCount, len(decoded), "should decode correct number of consolidations")
			}
		})
	}
}

// TestConsolidationRequestsOrdering tests that ordering is preserved
func TestConsolidationRequestsOrdering(t *testing.T) {
	// Create consolidations with specific ordering
	consolidations := types.ConsolidationRequests{
		{
			SourceAddress: common.ExecutionAddress{100}, // Higher value first
			SourcePubKey:  crypto.BLSPubkey{200},
			TargetPubKey:  crypto.BLSPubkey{201},
		},
		{
			SourceAddress: common.ExecutionAddress{50}, // Lower value second
			SourcePubKey:  crypto.BLSPubkey{100},
			TargetPubKey:  crypto.BLSPubkey{101},
		},
	}

	// Marshal
	marshaled, err := consolidations.MarshalSSZ()
	require.NoError(t, err)

	// Decode
	decoded, err := types.DecodeConsolidationRequests(marshaled)
	require.NoError(t, err)

	// Verify ordering is preserved
	require.Equal(t, 2, len(decoded))
	require.Equal(t, byte(100), decoded[0].SourceAddress[0], "first consolidation address should be preserved")
	require.Equal(t, byte(200), decoded[0].SourcePubKey[0], "first consolidation source pubkey should be preserved")
	require.Equal(t, byte(50), decoded[1].SourceAddress[0], "second consolidation address should be preserved")
	require.Equal(t, byte(100), decoded[1].SourcePubKey[0], "second consolidation source pubkey should be preserved")
}

// TestConsolidationRequestsRoundTrip tests round-trip encoding/decoding
func TestConsolidationRequestsRoundTrip(t *testing.T) {
	// Create various consolidation requests
	original := types.ConsolidationRequests{
		{
			SourceAddress: common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			SourcePubKey:  crypto.BLSPubkey{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
			TargetPubKey:  crypto.BLSPubkey{50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 76, 77},
		},
		{
			SourceAddress: common.ExecutionAddress{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119},
			SourcePubKey:  crypto.BLSPubkey{120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147},
			TargetPubKey:  crypto.BLSPubkey{150, 151, 152, 153, 154, 155, 156, 157, 158, 159, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169, 170, 171, 172, 173, 174, 175, 176, 177},
		},
	}

	// Marshal
	marshaled, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Decode
	decoded, err := types.DecodeConsolidationRequests(marshaled)
	require.NoError(t, err)

	// Verify round trip preserved all data
	require.Equal(t, len(original), len(decoded), "lengths should match")
	for i := range original {
		require.Equal(t, original[i].SourceAddress, decoded[i].SourceAddress, "source address %d should match", i)
		require.Equal(t, original[i].SourcePubKey, decoded[i].SourcePubKey, "source pubkey %d should match", i)
		require.Equal(t, original[i].TargetPubKey, decoded[i].TargetPubKey, "target pubkey %d should match", i)
	}
}

// TestConsolidationRequestsBoundaryValues tests boundary values and edge cases
func TestConsolidationRequestsBoundaryValues(t *testing.T) {
	// Test with all zeros
	zeroConsolidation := &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{},
		SourcePubKey:  crypto.BLSPubkey{},
		TargetPubKey:  crypto.BLSPubkey{},
	}

	consolidations := types.ConsolidationRequests{zeroConsolidation}
	marshaled, err := consolidations.MarshalSSZ()
	require.NoError(t, err)

	decoded, err := types.DecodeConsolidationRequests(marshaled)
	require.NoError(t, err)
	require.Equal(t, 1, len(decoded))
	require.Equal(t, zeroConsolidation.SourceAddress, decoded[0].SourceAddress)
	require.Equal(t, zeroConsolidation.SourcePubKey, decoded[0].SourcePubKey)
	require.Equal(t, zeroConsolidation.TargetPubKey, decoded[0].TargetPubKey)

	// Test with all max values (0xFF)
	maxConsolidation := &types.ConsolidationRequest{}
	for i := range maxConsolidation.SourceAddress {
		maxConsolidation.SourceAddress[i] = 0xFF
	}
	for i := range maxConsolidation.SourcePubKey {
		maxConsolidation.SourcePubKey[i] = 0xFF
	}
	for i := range maxConsolidation.TargetPubKey {
		maxConsolidation.TargetPubKey[i] = 0xFF
	}

	consolidations = types.ConsolidationRequests{maxConsolidation}
	marshaled, err = consolidations.MarshalSSZ()
	require.NoError(t, err)

	decoded, err = types.DecodeConsolidationRequests(marshaled)
	require.NoError(t, err)
	require.Equal(t, 1, len(decoded))
	require.Equal(t, maxConsolidation.SourceAddress, decoded[0].SourceAddress)
	require.Equal(t, maxConsolidation.SourcePubKey, decoded[0].SourcePubKey)
	require.Equal(t, maxConsolidation.TargetPubKey, decoded[0].TargetPubKey)
}
