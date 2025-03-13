// SPDX-License-Identifier: MIT
//
// Copyright (c) 2025 Berachain Foundation
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use,
// copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the
// Software is furnished to do so, subject to the following
// conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
// OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
// WdeHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

// TODO(pectra): Add tests
// 1. Marshalling/Unmarshalling invalid values.
// 3. NewSignedBeaconBlockFromSSZ tests.

func TestDepositRequest_ValidValuesSSZ(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		depositRequest *types.DepositRequest
	}{
		{
			name: "basic",
			depositRequest: &types.DepositRequest{
				// 48-byte public key
				Pubkey: crypto.BLSPubkey{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
					21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
					31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
					41, 42, 43, 44, 45, 46, 47, 48,
				},
				// 32-byte withdrawal credentials
				WithdrawalCredentials: [32]byte{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
					21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
					31, 32,
				},
				Amount: 1000,
				// 96-byte BLS signature
				Signature: crypto.BLSSignature{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
					21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
					31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
					41, 42, 43, 44, 45, 46, 47, 48, 49, 50,
					51, 52, 53, 54, 55, 56, 57, 58, 59, 60,
					61, 62, 63, 64, 65, 66, 67, 68, 69, 70,
					71, 72, 73, 74, 75, 76, 77, 78, 79, 80,
					81, 82, 83, 84, 85, 86, 87, 88, 89, 90,
					91, 92, 93, 94, 95, 96,
				},
				Index: 1,
			},
		},
		{
			name: "zero amount",
			depositRequest: &types.DepositRequest{
				Pubkey: crypto.BLSPubkey{
					10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
					20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
					30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
					40, 41, 42, 43, 44, 45, 46, 47, 48,
				},
				WithdrawalCredentials: [32]byte{
					10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
					20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
					30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
					40, 41,
				},
				Amount: 0,
				Signature: crypto.BLSSignature{
					10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
					20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
					30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
					40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
					50, 51, 52, 53, 54, 55, 56, 57, 58, 59,
					60, 61, 62, 63, 64, 65, 66, 67, 68, 69,
					70, 71, 72, 73, 74, 75, 76, 77, 78, 79,
					80, 81, 82, 83, 84, 85, 86, 87, 88, 89,
					90, 91, 92, 93, 94, 95, 96,
				},
				Index: 2,
			},
		},
		{
			name: "max values",
			depositRequest: &types.DepositRequest{
				Pubkey: crypto.BLSPubkey{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255,
				},
				WithdrawalCredentials: [32]byte{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255,
				},
				Amount: 1<<64 - 1,
				Signature: crypto.BLSSignature{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255,
				},
				Index: 3,
			},
		},
		{
			name: "random-ish values",
			depositRequest: &types.DepositRequest{
				Pubkey: crypto.BLSPubkey{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
					37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
					47, 48, 49, 50, 51, 52, 53, 54,
				},
				WithdrawalCredentials: [32]byte{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
					37, 38,
				},
				Amount: 54321,
				Signature: crypto.BLSSignature{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
					37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
					47, 48, 49, 50, 51, 52, 53, 54, 55, 56,
					57, 58, 59, 60, 61, 62, 63, 64, 65, 66,
					67, 68, 69, 70, 71, 72, 73, 74, 75, 76,
					77, 78, 79, 80, 81, 82, 83, 84, 85, 86,
					87, 88, 89, 90, 91, 92, 93, 94, 95, 96,
				},
				Index: 4,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the original deposit request.
			depositRequestBytes, err := tc.depositRequest.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal into a Prysm deposit request.
			var prysmDeposit enginev1.DepositRequest
			err = prysmDeposit.UnmarshalSSZ(depositRequestBytes)
			require.NoError(t, err)

			// Compare the HashTreeRoots: first compute the HTRs.
			prysmHTR, err := prysmDeposit.HashTreeRoot()
			require.NoError(t, err)
			depositHTR := tc.depositRequest.HashTreeRoot()
			// Compare the HashTreeRoots to ensure all fields were correctly interpreted.
			require.Equal(t, depositHTR[:], prysmHTR[:])

			// Marshal the Prysm deposit request.
			prysmDepositBytes, err := prysmDeposit.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into a new DepositRequest.
			var recomputedDepositRequest types.DepositRequest
			err = recomputedDepositRequest.UnmarshalSSZ(prysmDepositBytes)
			require.NoError(t, err)

			// Compare that the original and recomputed deposit requests match.
			require.Equal(t, *tc.depositRequest, recomputedDepositRequest)
		})
	}
}

func TestWithdrawalRequest_ValidValuesSSZ(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name              string
		withdrawalRequest *types.WithdrawalRequest
	}{
		{
			name: "basic",
			withdrawalRequest: &types.WithdrawalRequest{
				// 20-byte execution address
				SourceAddress: common.ExecutionAddress{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				// 48-byte public key
				ValidatorPubKey: crypto.BLSPubkey{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
					21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
					31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
					41, 42, 43, 44, 45, 46, 47, 48,
				},
				Amount: 1000,
			},
		},
		{
			name: "zero amount",
			withdrawalRequest: &types.WithdrawalRequest{
				SourceAddress: common.ExecutionAddress{
					10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
					20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
				},
				ValidatorPubKey: crypto.BLSPubkey{
					10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
					20, 21, 22, 23, 24, 25, 26, 27, 28, 29,
					30, 31, 32, 33, 34, 35, 36, 37, 38, 39,
					40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
					50, 51, 52, 53, 54, 55, 56, 57,
				},
				Amount: 0,
			},
		},
		{
			name: "max values",
			withdrawalRequest: &types.WithdrawalRequest{
				SourceAddress: common.ExecutionAddress{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
				},
				ValidatorPubKey: crypto.BLSPubkey{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255,
				},
				Amount: 1<<64 - 1,
			},
		},
		{
			name: "random-ish values",
			withdrawalRequest: &types.WithdrawalRequest{
				SourceAddress: common.ExecutionAddress{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
				},
				ValidatorPubKey: crypto.BLSPubkey{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
					37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
					47, 48, 49, 50, 51, 52, 53, 54,
				},
				Amount: 54321,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the original withdrawal request.
			withdrawalRequestBytes, err := tc.withdrawalRequest.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal into a Prysm withdrawal request.
			var prysmWithdrawal enginev1.WithdrawalRequest
			err = prysmWithdrawal.UnmarshalSSZ(withdrawalRequestBytes)
			require.NoError(t, err)

			prysmHTR, err := prysmWithdrawal.HashTreeRoot()
			require.NoError(t, err)
			withdrawalHTR := tc.withdrawalRequest.HashTreeRoot()

			// Compare the HashTreeRoots. Effectively a test for comparing all field values.
			require.Equal(t, withdrawalHTR[:], prysmHTR[:])

			// Marshal the Prysm withdrawal request.
			prysmWithdrawalBytes, err := prysmWithdrawal.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into a new WithdrawalRequest.
			var recomputedWithdrawalRequest types.WithdrawalRequest
			err = recomputedWithdrawalRequest.UnmarshalSSZ(prysmWithdrawalBytes)
			require.NoError(t, err)

			// Compare that the original and recomputed values match.
			require.Equal(t, *tc.withdrawalRequest, recomputedWithdrawalRequest)
		})
	}
}

func TestConsolidationRequest_ValidValuesSSZ(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		consolidationRequest *types.ConsolidationRequest
	}{
		{
			name: "basic",
			consolidationRequest: &types.ConsolidationRequest{
				// 20-byte execution address for SourceAddress
				SourceAddress: common.ExecutionAddress{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
				},
				// 48-byte public key for SourcePubKey
				SourcePubKey: crypto.BLSPubkey{
					1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
					11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
					21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
					31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
					41, 42, 43, 44, 45, 46, 47, 48,
				},
				// 48-byte public key for TargetPubKey
				TargetPubKey: crypto.BLSPubkey{
					48, 47, 46, 45, 44, 43, 42, 41, 40, 39,
					38, 37, 36, 35, 34, 33, 32, 31, 30, 29,
					28, 27, 26, 25, 24, 23, 22, 21, 20, 19,
					18, 17, 16, 15, 14, 13, 12, 11, 10, 9,
					8, 7, 6, 5, 4, 3, 2, 1,
				},
			},
		},
		{
			name: "max values",
			consolidationRequest: &types.ConsolidationRequest{
				SourceAddress: common.ExecutionAddress{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
				},
				SourcePubKey: crypto.BLSPubkey{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255,
				},
				TargetPubKey: crypto.BLSPubkey{
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255, 255, 255,
					255, 255, 255, 255, 255, 255, 255, 255,
				},
			},
		},
		{
			name: "random-ish values",
			consolidationRequest: &types.ConsolidationRequest{
				SourceAddress: common.ExecutionAddress{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
				},
				SourcePubKey: crypto.BLSPubkey{
					7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
					17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
					27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
					37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
					47, 48, 49, 50, 51, 52, 53, 54,
				},
				TargetPubKey: crypto.BLSPubkey{
					14, 15, 16, 17, 18, 19, 20, 21, 22, 23,
					24, 25, 26, 27, 28, 29, 30, 31, 32, 33,
					34, 35, 36, 37, 38, 39, 40, 41, 42, 43,
					44, 45, 46, 47, 48, 49, 50, 51, 52, 53,
					54, 55, 56, 57, 58, 59, 60, 61,
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the original consolidation request.
			crBytes, err := tc.consolidationRequest.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal into a Prysm consolidation request.
			var prysmCR enginev1.ConsolidationRequest
			err = prysmCR.UnmarshalSSZ(crBytes)
			require.NoError(t, err)

			prysmHTR, err := prysmCR.HashTreeRoot()
			require.NoError(t, err)
			crHTR := tc.consolidationRequest.HashTreeRoot()

			// Compare the HashTreeRoots. This effectively tests that all fields were encoded correctly.
			require.Equal(t, crHTR[:], prysmHTR[:])

			// Marshal the Prysm consolidation request.
			prysmCRBytes, err := prysmCR.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into a new ConsolidationRequest.
			var recomputedCR types.ConsolidationRequest
			err = recomputedCR.UnmarshalSSZ(prysmCRBytes)
			require.NoError(t, err)

			// Compare that the original and recomputed consolidation requests match.
			require.Equal(t, *tc.consolidationRequest, recomputedCR)
		})
	}
}
