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

// NOTE: WithdrawalRequests is a slice type that uses EIP-7685 encoding format.
// While the individual WithdrawalRequest elements were using karalabe/ssz,
// the slice encoding follows a different format specified by EIP-7685. These tests
// focus on the encoding/decoding functions rather than direct SSZ compatibility.
// The individual element type (WithdrawalRequest) has its own compatibility test.

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

// TestWithdrawalRequestsCompatibility tests that WithdrawalRequests slice type
// behaves correctly with EIP-7685 encoding/decoding.
func TestWithdrawalRequestsCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() types.WithdrawalRequests
	}{
		// Note: Empty slices are handled differently - the Decode function
		// requires at least one element due to the validation check.
		{
			name: "single withdrawal request",
			setup: func() types.WithdrawalRequests {
				addr := common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
				pubkey := crypto.BLSPubkey{20, 21, 22, 23, 24, 25, 26, 27}
				amount := math.Gwei(32000000000) // 32 ETH

				return types.WithdrawalRequests{
					&types.WithdrawalRequest{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          amount,
					},
				}
			},
		},
		{
			name: "multiple withdrawal requests",
			setup: func() types.WithdrawalRequests {
				withdrawals := make(types.WithdrawalRequests, 5)

				for i := 0; i < 5; i++ {
					addr := common.ExecutionAddress{}
					pubkey := crypto.BLSPubkey{}

					// Fill with unique data for each withdrawal
					for j := range addr {
						addr[j] = byte(i*20 + j)
					}
					for j := range pubkey {
						pubkey[j] = byte(i*48 + j)
					}

					amount := math.Gwei(1000000000 * uint64(i+1)) // Varying amounts

					withdrawals[i] = &types.WithdrawalRequest{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          amount,
					}
				}

				return withdrawals
			},
		},
		{
			name: "withdrawals with various amounts",
			setup: func() types.WithdrawalRequests {
				amounts := []uint64{
					0,           // Full exit request
					1000000000,  // 1 ETH
					16000000000, // 16 ETH
					32000000000, // 32 ETH
				}

				withdrawals := make(types.WithdrawalRequests, len(amounts))
				for i, amount := range amounts {
					addr := common.ExecutionAddress{}
					addr[0] = byte(i)
					pubkey := crypto.BLSPubkey{}
					pubkey[0] = byte(i + 50)

					withdrawals[i] = &types.WithdrawalRequest{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          math.Gwei(amount),
					}
				}

				return withdrawals
			},
		},
		{
			name: "maximum withdrawal requests",
			setup: func() types.WithdrawalRequests {
				// Create maximum allowed (16)
				maxWithdrawals := make(types.WithdrawalRequests, 16)
				for i := range maxWithdrawals {
					addr := common.ExecutionAddress{}
					pubkey := crypto.BLSPubkey{}

					// Use different patterns for variety
					for j := range addr {
						addr[j] = byte((i + j) % 256)
					}
					for j := range pubkey {
						pubkey[j] = byte((i*2 + j) % 256)
					}

					maxWithdrawals[i] = &types.WithdrawalRequest{
						SourceAddress:   addr,
						ValidatorPubKey: pubkey,
						Amount:          math.Gwei(1000000000 * uint64(i+1)),
					}
				}

				return maxWithdrawals
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			withdrawals := tc.setup()

			// Test Marshal
			marshaled, err := withdrawals.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")

			// Test Decode (which is the unmarshal operation for this type)
			decoded, err := types.DecodeWithdrawalRequests(marshaled)
			require.NoError(t, err, "DecodeWithdrawalRequests should not error")

			// Verify lengths match
			require.Equal(t, len(withdrawals), len(decoded), "decoded length should match original")

			// Verify each withdrawal matches
			for i := range withdrawals {
				require.Equal(t, withdrawals[i].SourceAddress, decoded[i].SourceAddress, "source address at index %d should match", i)
				require.Equal(t, withdrawals[i].ValidatorPubKey, decoded[i].ValidatorPubKey, "validator pubkey at index %d should match", i)
				require.Equal(t, withdrawals[i].Amount, decoded[i].Amount, "amount at index %d should match", i)
			}

			// Test ValidateAfterDecodingSSZ
			err = decoded.ValidateAfterDecodingSSZ()
			require.NoError(t, err, "ValidateAfterDecodingSSZ should not error for valid data")
		})
	}
}

// TestWithdrawalRequestsEmptySlice tests empty slice handling separately
func TestWithdrawalRequestsEmptySlice(t *testing.T) {
	// Empty slice marshals to empty byte array
	empty := types.WithdrawalRequests{}
	marshaled, err := empty.MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, 0, len(marshaled), "empty slice should marshal to empty bytes")

	// But decoding empty bytes fails due to validation
	_, err = types.DecodeWithdrawalRequests(marshaled)
	require.Error(t, err, "decoding empty bytes should error")
	require.Contains(t, err.Error(), "invalid withdrawal requests SSZ size", "error should indicate size issue")
}

// TestWithdrawalRequestsEIP7685Encoding tests the specific EIP-7685 encoding format
func TestWithdrawalRequestsEIP7685Encoding(t *testing.T) {
	// Create a simple withdrawal request
	withdrawal := &types.WithdrawalRequest{
		SourceAddress:   common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
		ValidatorPubKey: crypto.BLSPubkey{1, 2, 3},
		Amount:          32000000000,
	}

	withdrawals := types.WithdrawalRequests{withdrawal}

	// Marshal using EIP-7685 format
	marshaled, err := withdrawals.MarshalSSZ()
	require.NoError(t, err)

	// Verify the marshaled data is exactly 76 bytes (one withdrawal)
	require.Equal(t, 76, len(marshaled), "single withdrawal should be 76 bytes")

	// Test with multiple withdrawals
	withdrawals = types.WithdrawalRequests{withdrawal, withdrawal, withdrawal}
	marshaled, err = withdrawals.MarshalSSZ()
	require.NoError(t, err)

	// Verify the marshaled data is exactly 228 bytes (three withdrawals)
	require.Equal(t, 228, len(marshaled), "three withdrawals should be 228 bytes")
}

// TestWithdrawalRequestsValidation tests the validation logic
func TestWithdrawalRequestsValidation(t *testing.T) {
	// Create maximum allowed withdrawals
	maxWithdrawals := make(types.WithdrawalRequests, 16) // MaxWithdrawalRequestsPerPayload
	for i := range maxWithdrawals {
		addr := common.ExecutionAddress{}
		addr[0] = byte(i)

		maxWithdrawals[i] = &types.WithdrawalRequest{
			SourceAddress:   addr,
			ValidatorPubKey: crypto.BLSPubkey{byte(i)},
			Amount:          32000000000,
		}
	}

	// This should validate successfully
	err := maxWithdrawals.ValidateAfterDecodingSSZ()
	require.NoError(t, err, "max withdrawals should validate successfully")

	// Create one more than allowed
	tooManyWithdrawals := append(maxWithdrawals, &types.WithdrawalRequest{
		SourceAddress:   common.ExecutionAddress{255},
		ValidatorPubKey: crypto.BLSPubkey{255},
		Amount:          32000000000,
	})

	// This should fail validation
	err = tooManyWithdrawals.ValidateAfterDecodingSSZ()
	require.Error(t, err, "too many withdrawals should fail validation")
	require.Contains(t, err.Error(), "invalid number of withdrawal requests", "error message should indicate too many withdrawals")
}

// TestWithdrawalRequestsDecodeErrors tests error cases in decoding
func TestWithdrawalRequestsDecodeErrors(t *testing.T) {
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
			errorContains: "invalid withdrawal requests SSZ size",
		},
		{
			name:          "partial withdrawal data",
			data:          make([]byte, 50), // Less than 76 bytes
			expectError:   true,
			errorContains: "invalid withdrawal requests SSZ size",
		},
		{
			name:          "data not multiple of withdrawal size",
			data:          make([]byte, 100), // Not a multiple of 76
			expectError:   true,
			errorContains: "is not a multiple of item size",
		},
		{
			name:        "valid single withdrawal",
			data:        make([]byte, 76), // Exactly one withdrawal
			expectError: false,
		},
		{
			name:        "valid multiple withdrawals",
			data:        make([]byte, 228), // Exactly three withdrawals
			expectError: false,
		},
		{
			name:          "exceeds maximum size",
			data:          make([]byte, 76*17), // 17 withdrawals (max is 16)
			expectError:   true,
			errorContains: "invalid withdrawal requests SSZ size",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := types.DecodeWithdrawalRequests(tc.data)

			if tc.expectError {
				require.Error(t, err, "should error")
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains, "error message should contain expected text")
				}
			} else {
				require.NoError(t, err, "should not error")
				require.NotNil(t, decoded, "decoded should not be nil")
				// Verify we got the expected number of withdrawals
				expectedCount := len(tc.data) / 76
				require.Equal(t, expectedCount, len(decoded), "should decode correct number of withdrawals")
			}
		})
	}
}

// TestWithdrawalRequestsOrdering tests that ordering is preserved
func TestWithdrawalRequestsOrdering(t *testing.T) {
	// Create withdrawals with specific ordering based on amount
	count := 10
	withdrawals := make(types.WithdrawalRequests, count)

	for i := 0; i < count; i++ {
		addr := common.ExecutionAddress{}
		pubkey := crypto.BLSPubkey{}

		// Make address based on reverse order to ensure ordering is preserved
		addr[0] = byte(count - i - 1)
		pubkey[0] = byte(count - i - 1)

		// Use descending amounts to test ordering
		amount := math.Gwei(uint64(count-i) * 1000000000)

		withdrawals[i] = &types.WithdrawalRequest{
			SourceAddress:   addr,
			ValidatorPubKey: pubkey,
			Amount:          amount,
		}
	}

	// Marshal
	marshaled, err := withdrawals.MarshalSSZ()
	require.NoError(t, err)

	// Decode
	decoded, err := types.DecodeWithdrawalRequests(marshaled)
	require.NoError(t, err)

	// Verify ordering is preserved
	require.Equal(t, count, len(decoded))

	for i := 0; i < count; i++ {
		expectedByte := byte(count - i - 1)
		expectedAmount := math.Gwei(uint64(count-i) * 1000000000)

		require.Equal(t, expectedByte, decoded[i].SourceAddress[0], "address ordering preserved at index %d", i)
		require.Equal(t, expectedByte, decoded[i].ValidatorPubKey[0], "pubkey ordering preserved at index %d", i)
		require.Equal(t, expectedAmount, decoded[i].Amount, "amount ordering preserved at index %d", i)
	}
}

// TestWithdrawalRequestsRoundTrip tests round-trip encoding/decoding
func TestWithdrawalRequestsRoundTrip(t *testing.T) {
	// Create various withdrawal requests
	original := types.WithdrawalRequests{
		{
			SourceAddress:   common.ExecutionAddress{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			ValidatorPubKey: crypto.BLSPubkey{20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47},
			Amount:          0, // Full exit
		},
		{
			SourceAddress:   common.ExecutionAddress{50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69},
			ValidatorPubKey: crypto.BLSPubkey{70, 71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97},
			Amount:          16000000000, // 16 ETH
		},
		{
			SourceAddress:   common.ExecutionAddress{100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115, 116, 117, 118, 119},
			ValidatorPubKey: crypto.BLSPubkey{120, 121, 122, 123, 124, 125, 126, 127, 128, 129, 130, 131, 132, 133, 134, 135, 136, 137, 138, 139, 140, 141, 142, 143, 144, 145, 146, 147},
			Amount:          32000000000, // 32 ETH
		},
	}

	// Marshal
	marshaled, err := original.MarshalSSZ()
	require.NoError(t, err)

	// Decode
	decoded, err := types.DecodeWithdrawalRequests(marshaled)
	require.NoError(t, err)

	// Verify round trip preserved all data
	require.Equal(t, len(original), len(decoded), "lengths should match")
	for i := range original {
		require.Equal(t, original[i].SourceAddress, decoded[i].SourceAddress, "source address %d should match", i)
		require.Equal(t, original[i].ValidatorPubKey, decoded[i].ValidatorPubKey, "validator pubkey %d should match", i)
		require.Equal(t, original[i].Amount, decoded[i].Amount, "amount %d should match", i)
	}
}
