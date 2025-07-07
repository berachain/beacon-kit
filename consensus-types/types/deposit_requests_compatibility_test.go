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

// NOTE: DepositRequests is a slice type that uses EIP-7685 encoding format.
// While the individual DepositRequest elements were using karalabe/ssz,
// the slice encoding follows a different format specified by EIP-7685. These tests
// focus on the encoding/decoding functions rather than direct SSZ compatibility.
// The individual element type (DepositRequest) has its own compatibility test.

//go:build test

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

// TestDepositRequestsCompatibility tests that DepositRequests slice type
// behaves correctly with EIP-7685 encoding/decoding.
func TestDepositRequestsCompatibility(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() types.DepositRequests
	}{
		// Note: Empty slices are handled differently - the Decode function
		// requires at least one element due to the validation check.
		{
			name: "single deposit request",
			setup: func() types.DepositRequests {
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

				return types.DepositRequests{
					&types.DepositRequest{
						Pubkey:      pubkey,
						Credentials: creds,
						Amount:      amount,
						Signature:   sig,
						Index:       index,
					},
				}
			},
		},
		{
			name: "multiple deposit requests",
			setup: func() types.DepositRequests {
				deposits := make(types.DepositRequests, 5)
				
				for i := 0; i < 5; i++ {
					pubkey := crypto.BLSPubkey{}
					creds := types.WithdrawalCredentials{}
					sig := crypto.BLSSignature{}
					
					// Fill with unique data for each deposit
					for j := range pubkey {
						pubkey[j] = byte(i*48 + j)
					}
					creds[0] = 0x01 // ETH1_ADDRESS_WITHDRAWAL_PREFIX
					for j := 1; j < 32; j++ {
						creds[j] = byte(i*32 + j)
					}
					for j := range sig {
						sig[j] = byte(i*96 + j)
					}
					
					amount := math.Gwei(32000000000 + uint64(i)*1000000000)
					index := uint64(i * 1000)
					
					deposits[i] = &types.DepositRequest{
						Pubkey:      pubkey,
						Credentials: creds,
						Amount:      amount,
						Signature:   sig,
						Index:       index,
					}
				}
				
				return deposits
			},
		},
		{
			name: "deposits with various amounts",
			setup: func() types.DepositRequests {
				amounts := []uint64{
					1000000000,    // 1 ETH
					16000000000,   // 16 ETH
					32000000000,   // 32 ETH
					100000000000,  // 100 ETH
				}
				
				deposits := make(types.DepositRequests, len(amounts))
				for i, amount := range amounts {
					pubkey := crypto.BLSPubkey{}
					pubkey[0] = byte(i)
					
					deposits[i] = &types.DepositRequest{
						Pubkey:      pubkey,
						Credentials: types.WithdrawalCredentials{0x01},
						Amount:      math.Gwei(amount),
						Signature:   crypto.BLSSignature{byte(i), byte(i + 1)},
						Index:       uint64(i * 100),
					}
				}
				
				return deposits
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			deposits := tc.setup()

			// Test Marshal
			marshaled, err := deposits.MarshalSSZ()
			require.NoError(t, err, "MarshalSSZ should not error")

			// Test Decode (which is the unmarshal operation for this type)
			decoded, err := types.DecodeDepositRequests(marshaled)
			require.NoError(t, err, "DecodeDepositRequests should not error")
			
			// Verify lengths match
			require.Equal(t, len(deposits), len(decoded), "decoded length should match original")
			
			// Verify each deposit matches
			for i := range deposits {
				require.Equal(t, deposits[i].Pubkey, decoded[i].Pubkey, "pubkey at index %d should match", i)
				require.Equal(t, deposits[i].Credentials, decoded[i].Credentials, "credentials at index %d should match", i)
				require.Equal(t, deposits[i].Amount, decoded[i].Amount, "amount at index %d should match", i)
				require.Equal(t, deposits[i].Signature, decoded[i].Signature, "signature at index %d should match", i)
				require.Equal(t, deposits[i].Index, decoded[i].Index, "index at index %d should match", i)
			}

			// Test ValidateAfterDecodingSSZ
			err = decoded.ValidateAfterDecodingSSZ()
			require.NoError(t, err, "ValidateAfterDecodingSSZ should not error for valid data")
		})
	}
}

// TestDepositRequestsEmptySlice tests empty slice handling separately
func TestDepositRequestsEmptySlice(t *testing.T) {
	// Empty slice marshals to empty byte array
	empty := types.DepositRequests{}
	marshaled, err := empty.MarshalSSZ()
	require.NoError(t, err)
	require.Equal(t, 0, len(marshaled), "empty slice should marshal to empty bytes")
	
	// But decoding empty bytes fails due to validation
	_, err = types.DecodeDepositRequests(marshaled)
	require.Error(t, err, "decoding empty bytes should error")
	require.Contains(t, err.Error(), "invalid deposit requests SSZ size", "error should indicate size issue")
}

// TestDepositRequestsEIP7685Encoding tests the specific EIP-7685 encoding format
func TestDepositRequestsEIP7685Encoding(t *testing.T) {
	// Create a simple deposit request
	deposit := &types.DepositRequest{
		Pubkey:      crypto.BLSPubkey{1, 2, 3},
		Credentials: types.WithdrawalCredentials{0x01},
		Amount:      32000000000,
		Signature:   crypto.BLSSignature{4, 5, 6},
		Index:       100,
	}
	
	deposits := types.DepositRequests{deposit}
	
	// Marshal using EIP-7685 format
	marshaled, err := deposits.MarshalSSZ()
	require.NoError(t, err)
	
	// Verify the marshaled data is exactly 192 bytes (one deposit)
	require.Equal(t, 192, len(marshaled), "single deposit should be 192 bytes")
	
	// Test with multiple deposits
	deposits = types.DepositRequests{deposit, deposit, deposit}
	marshaled, err = deposits.MarshalSSZ()
	require.NoError(t, err)
	
	// Verify the marshaled data is exactly 576 bytes (three deposits)
	require.Equal(t, 576, len(marshaled), "three deposits should be 576 bytes")
}

// TestDepositRequestsValidation tests the validation logic
func TestDepositRequestsValidation(t *testing.T) {
	// Create maximum allowed deposits
	maxDeposits := make(types.DepositRequests, 8192) // MaxDepositRequestsPerPayload
	for i := range maxDeposits {
		maxDeposits[i] = &types.DepositRequest{
			Pubkey:      crypto.BLSPubkey{byte(i % 256)},
			Credentials: types.WithdrawalCredentials{0x01},
			Amount:      32000000000,
			Signature:   crypto.BLSSignature{},
			Index:       uint64(i),
		}
	}
	
	// This should validate successfully
	err := maxDeposits.ValidateAfterDecodingSSZ()
	require.NoError(t, err, "max deposits should validate successfully")
	
	// Create one more than allowed
	tooManyDeposits := append(maxDeposits, &types.DepositRequest{
		Pubkey:      crypto.BLSPubkey{255},
		Credentials: types.WithdrawalCredentials{0x01},
		Amount:      32000000000,
		Signature:   crypto.BLSSignature{},
		Index:       9999,
	})
	
	// This should fail validation
	err = tooManyDeposits.ValidateAfterDecodingSSZ()
	require.Error(t, err, "too many deposits should fail validation")
	require.Contains(t, err.Error(), "invalid number of deposit requests", "error message should indicate too many deposits")
}

// TestDepositRequestsDecodeErrors tests error cases in decoding
func TestDepositRequestsDecodeErrors(t *testing.T) {
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
			errorContains: "invalid deposit requests SSZ size",
		},
		{
			name:          "partial deposit data",
			data:          make([]byte, 100), // Less than 192 bytes
			expectError:   true,
			errorContains: "invalid deposit requests SSZ size",
		},
		{
			name:          "data not multiple of deposit size",
			data:          make([]byte, 300), // Not a multiple of 192
			expectError:   true,
			errorContains: "is not a multiple of item size",
		},
		{
			name:        "valid single deposit",
			data:        make([]byte, 192), // Exactly one deposit
			expectError: false,
		},
		{
			name:        "valid multiple deposits",
			data:        make([]byte, 576), // Exactly three deposits
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			decoded, err := types.DecodeDepositRequests(tc.data)
			
			if tc.expectError {
				require.Error(t, err, "should error")
				if tc.errorContains != "" {
					require.Contains(t, err.Error(), tc.errorContains, "error message should contain expected text")
				}
			} else {
				require.NoError(t, err, "should not error")
				require.NotNil(t, decoded, "decoded should not be nil")
				// Verify we got the expected number of deposits
				expectedCount := len(tc.data) / 192
				require.Equal(t, expectedCount, len(decoded), "should decode correct number of deposits")
			}
		})
	}
}

// TestDepositRequestsOrdering tests that ordering is preserved
func TestDepositRequestsOrdering(t *testing.T) {
	// Create deposits with specific ordering based on index
	count := 10
	deposits := make(types.DepositRequests, count)
	
	for i := 0; i < count; i++ {
		pubkey := crypto.BLSPubkey{}
		// Make pubkey based on reverse order to ensure ordering is preserved
		pubkey[0] = byte(count - i - 1)
		
		// Use descending indices to test ordering
		index := uint64((count - i) * 1000)
		
		deposits[i] = &types.DepositRequest{
			Pubkey:      pubkey,
			Credentials: types.WithdrawalCredentials{0x01},
			Amount:      math.Gwei(32000000000),
			Signature:   crypto.BLSSignature{byte(i), byte(i + 1)},
			Index:       index,
		}
	}
	
	// Marshal
	marshaled, err := deposits.MarshalSSZ()
	require.NoError(t, err)
	
	// Decode
	decoded, err := types.DecodeDepositRequests(marshaled)
	require.NoError(t, err)
	
	// Verify ordering is preserved
	require.Equal(t, count, len(decoded))
	
	for i := 0; i < count; i++ {
		expectedPubkeyByte := byte(count - i - 1)
		expectedIndex := uint64((count - i) * 1000)
		
		require.Equal(t, expectedPubkeyByte, decoded[i].Pubkey[0], "ordering preserved at index %d", i)
		require.Equal(t, expectedIndex, decoded[i].Index, "index preserved at index %d", i)
	}
}