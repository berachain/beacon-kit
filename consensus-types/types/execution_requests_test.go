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
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

func TestExecutionRequests_ValidValuesSSZ(t *testing.T) {
	t.Parallel()
	// Create a few helper instances to reuse in test cases.
	// You can reuse your existing tests' values for deposit, withdrawal, and consolidation.
	depositBasic := &types.DepositRequest{
		// 48-byte public key
		Pubkey: crypto.BLSPubkey{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 46, 47, 48,
		},
		Credentials: [32]byte{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32,
		},
		Amount: 1000,
		Signature: crypto.BLSSignature{1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
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
	}

	withdrawalBasic := &types.WithdrawalRequest{
		SourceAddress: common.ExecutionAddress{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		},
		ValidatorPubKey: crypto.BLSPubkey{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 46, 47, 48,
		},
		Amount: 1000,
	}

	consolidationBasic := &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		},
		SourcePubKey: crypto.BLSPubkey{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 46, 47, 48,
		},
		TargetPubKey: crypto.BLSPubkey{
			48, 47, 46, 45, 44, 43, 42, 41, 40, 39,
			38, 37, 36, 35, 34, 33, 32, 31, 30, 29,
			28, 27, 26, 25, 24, 23, 22, 21, 20, 19,
			18, 17, 16, 15, 14, 13, 12, 11, 10, 9,
			8, 7, 6, 5, 4, 3, 2, 1,
		},
	}

	// Define test cases. We vary the content of each slice.
	testCases := []struct {
		name              string
		executionRequests *types.ExecutionRequests
	}{
		{
			name: "all basic",
			executionRequests: &types.ExecutionRequests{
				Deposits:       []*types.DepositRequest{depositBasic},
				Withdrawals:    []*types.WithdrawalRequest{withdrawalBasic},
				Consolidations: []*types.ConsolidationRequest{consolidationBasic},
			},
		},
		{
			name: "empty slices",
			executionRequests: &types.ExecutionRequests{
				Deposits:       []*types.DepositRequest{},
				Withdrawals:    []*types.WithdrawalRequest{},
				Consolidations: []*types.ConsolidationRequest{},
			},
		},
		{
			name: "multiple entries",
			executionRequests: &types.ExecutionRequests{
				Deposits:       []*types.DepositRequest{depositBasic, depositBasic},
				Withdrawals:    []*types.WithdrawalRequest{withdrawalBasic, withdrawalBasic, withdrawalBasic},
				Consolidations: []*types.ConsolidationRequest{consolidationBasic, consolidationBasic},
			},
		},
		{
			name: "random-ish values",
			executionRequests: &types.ExecutionRequests{
				Deposits: []*types.DepositRequest{
					{
						Pubkey: crypto.BLSPubkey{
							7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
							17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
							27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
							37, 38, 39, 40, 41, 42, 43, 44, 45, 46,
							47, 48, 49, 50, 51, 52, 53, 54,
						},
						Credentials: [32]byte{
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
				Withdrawals: []*types.WithdrawalRequest{
					{
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
				Consolidations: []*types.ConsolidationRequest{
					{
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
			},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the original ExecutionRequests.
			execReqBytes, err := tc.executionRequests.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal into a Prysm ExecutionRequests.
			var prysmER enginev1.ExecutionRequests
			err = prysmER.UnmarshalSSZ(execReqBytes)
			require.NoError(t, err)

			prysmHTR, err := prysmER.HashTreeRoot()
			require.NoError(t, err)
			execReqHTR := tc.executionRequests.HashTreeRoot()

			// Compare the HashTreeRoots to ensure encoding was correct.
			require.Equal(t, execReqHTR[:], prysmHTR[:])

			// Marshal the Prysm ExecutionRequests.
			prysmERBytes, err := prysmER.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into a new ExecutionRequests.
			var recomputedER *types.ExecutionRequests
			recomputedER, err = recomputedER.NewFromSSZ(prysmERBytes)
			require.NoError(t, err)

			// Compare that the original and recomputed ExecutionRequests match.
			require.Equal(t, *tc.executionRequests, recomputedER)
		})
	}
}

// TestExecutionRequests_InvalidValuesUnmarshalSSZ ensures that Unmarshal must never panic.
//
//nolint:paralleltest // Invalid SSZ values cannot be run in parallel due to zeroalloc, which is global shared memory.
func TestExecutionRequests_InvalidValuesUnmarshalSSZ(t *testing.T) {
	// Define several invalid payloads.
	invalidPayloads := [][]byte{
		nil,                    // nil slice
		{},                     // empty slice
		[]byte("invalid data"), // arbitrary string data
		{0x00, 0x01},           // too short to be valid
		// A random 50-byte slice (likely invalid)
		func() []byte {
			b := make([]byte, 50)
			for i := range b {
				b[i] = byte((i * 3) % 256)
			}
			return b
		}(),
		// A truncated valid payload: marshal a valid empty ExecutionRequests and drop last 4 bytes.
		func() []byte {
			er := types.ExecutionRequests{
				Deposits:       []*types.DepositRequest{},
				Withdrawals:    []*types.WithdrawalRequest{},
				Consolidations: []*types.ConsolidationRequest{},
			}
			validBytes, err := er.MarshalSSZ()
			require.NoError(t, err)
			if len(validBytes) > 4 {
				return validBytes[:len(validBytes)-4]
			}
			return validBytes
		}(),
		// A valid payload with extra trailing bytes.
		func() []byte {
			er := types.ExecutionRequests{
				Deposits:       []*types.DepositRequest{},
				Withdrawals:    []*types.WithdrawalRequest{},
				Consolidations: []*types.ConsolidationRequest{},
			}
			validBytes, err := er.MarshalSSZ()
			require.NoError(t, err)
			// Append extra bytes that should make the payload invalid.
			extra := []byte{0xFF, 0xEE, 0xDD, 0xCC}
			return append(validBytes, extra...)
		}(),
	}

	for i, payload := range invalidPayloads {
		i, payload := i, payload // capture loop variables
		t.Run(fmt.Sprintf("invalidPayload_%d", i), func(t *testing.T) {
			var er *types.ExecutionRequests
			// Ensure that calling UnmarshalSSZ with an invalid payload does not panic
			// and returns a non-nil error.
			require.NotPanics(t, func() {
				var err error
				er, err = er.NewFromSSZ(payload)
				require.Error(t, err, "Expected error for payload %v", payload)
			})
		})
	}
}

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
				Credentials: [32]byte{
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
				Credentials: [32]byte{
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
				Credentials: [32]byte{
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
				Credentials: [32]byte{
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
			recomputedDepositRequest, err := (&types.DepositRequest{}).NewFromSSZ(prysmDepositBytes)
			require.NoError(t, err)

			// Compare that the original and recomputed deposit requests match.
			require.Equal(t, *tc.depositRequest, *recomputedDepositRequest)
		})
	}
}

//nolint:paralleltest // Invalid SSZ values cannot be run in parallel due to zeroalloc, which is global shared memory.
func TestDepositRequest_InvalidValuesUnmarshalSSZ(t *testing.T) {
	// Build a valid deposit request and marshal it to obtain a baseline valid payload.
	validDeposit := &types.DepositRequest{
		// 48-byte public key
		Pubkey: crypto.BLSPubkey{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 46, 47, 48,
		},
		// 32-byte withdrawal credentials
		Credentials: [32]byte{
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
	}
	validBytes, err := validDeposit.MarshalSSZ()
	require.NoError(t, err)

	// Build a slice of invalid payloads.
	invalidPayloads := [][]byte{
		nil,                       // nil slice
		{},                        // empty slice
		[]byte("this is not ssz"), // arbitrary non-SSZ data
		{0x00, 0x01, 0x02},        // too short to be valid
		// A truncated valid payload.
		func() []byte {
			if len(validBytes) > 5 {
				return validBytes[:len(validBytes)-5]
			}
			return validBytes
		}(),
		// A valid payload with extra trailing bytes.
		func() []byte {
			extra := []byte{0xAA, 0xBB, 0xCC, 0xDD}
			return append(validBytes, extra...)
		}(),
	}

	// Iterate over each invalid payload.
	for i, payload := range invalidPayloads {
		i, payload := i, payload // capture range variables
		t.Run(fmt.Sprintf("invalidPayload_%d", i), func(t *testing.T) {
			require.NotPanics(t, func() {
				_, err = (&types.DepositRequest{}).NewFromSSZ(payload)
				// We expect an error for every invalid payload.
				require.Error(t, err, "expected error for payload %v", payload)
			})
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

//nolint:paralleltest // Invalid SSZ values cannot be run in parallel due to zeroalloc, which is global shared memory.
func TestWithdrawalRequest_InvalidValuesUnmarshalSSZ(t *testing.T) {
	// Build a valid withdrawal request to obtain a baseline valid payload.
	validWithdrawal := &types.WithdrawalRequest{
		SourceAddress: common.ExecutionAddress{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		},
		ValidatorPubKey: crypto.BLSPubkey{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 46, 47, 48,
		},
		Amount: 1000,
	}
	validBytes, err := validWithdrawal.MarshalSSZ()
	require.NoError(t, err)

	// Build a slice of invalid payloads.
	invalidPayloads := [][]byte{
		nil,                       // nil slice
		{},                        // empty slice
		[]byte("this is not ssz"), // arbitrary non-SSZ data
		{0x00, 0x01},              // too short to be valid
		// A truncated valid payload: remove last 5 bytes.
		func() []byte {
			if len(validBytes) > 5 {
				return validBytes[:len(validBytes)-5]
			}
			return validBytes
		}(),
		// A valid payload with extra trailing bytes.
		func() []byte {
			extra := []byte{0xAA, 0xBB, 0xCC, 0xDD}
			return append(validBytes, extra...)
		}(),
	}

	for i, payload := range invalidPayloads {
		i, payload := i, payload // capture loop variables
		t.Run(fmt.Sprintf("invalidWithdrawal_%d", i), func(t *testing.T) {
			// Ensure that calling UnmarshalSSZ does not panic and returns an error.
			require.NotPanics(t, func() {
				var w types.WithdrawalRequest
				err = w.UnmarshalSSZ(payload)
				require.Error(t, err, "expected error for payload %v", payload)
			})
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

//nolint:paralleltest // Invalid SSZ values cannot be run in parallel due to zeroalloc, which is global shared memory.
func TestConsolidationRequest_InvalidValuesUnmarshalSSZ(t *testing.T) {
	// Build a valid consolidation request to get a baseline valid payload.
	validConsolidation := &types.ConsolidationRequest{
		SourceAddress: common.ExecutionAddress{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
		},
		SourcePubKey: crypto.BLSPubkey{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
			11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
			21, 22, 23, 24, 25, 26, 27, 28, 29, 30,
			31, 32, 33, 34, 35, 36, 37, 38, 39, 40,
			41, 42, 43, 44, 45, 46, 47, 48,
		},
		TargetPubKey: crypto.BLSPubkey{
			48, 47, 46, 45, 44, 43, 42, 41, 40, 39,
			38, 37, 36, 35, 34, 33, 32, 31, 30, 29,
			28, 27, 26, 25, 24, 23, 22, 21, 20, 19,
			18, 17, 16, 15, 14, 13, 12, 11, 10, 9,
			8, 7, 6, 5, 4, 3, 2, 1,
		},
	}
	validBytes, err := validConsolidation.MarshalSSZ()
	require.NoError(t, err)

	// Build a slice of invalid payloads.
	invalidPayloads := [][]byte{
		nil,                       // nil slice
		{},                        // empty slice
		[]byte("this is not ssz"), // arbitrary non-SSZ data
		{0x00, 0x01, 0x02},        // too short to be valid
		// A truncated valid payload.
		func() []byte {
			if len(validBytes) > 5 {
				return validBytes[:len(validBytes)-5]
			}
			return validBytes
		}(),
		// A valid payload with extra trailing bytes.
		func() []byte {
			extra := []byte{0xAA, 0xBB, 0xCC, 0xDD}
			return append(validBytes, extra...)
		}(),
	}

	for i, payload := range invalidPayloads {
		i, payload := i, payload // capture loop variables
		t.Run(fmt.Sprintf("invalidConsolidation_%d", i), func(t *testing.T) {
			// Ensure that calling UnmarshalSSZ does not panic and returns an error.
			require.NotPanics(t, func() {
				var c types.ConsolidationRequest
				err = c.UnmarshalSSZ(payload)
				require.Error(t, err, "expected error for payload %v", payload)
			})
		})
	}
}

func TestDecodeExecutionRequests_PrysmTests(t *testing.T) {
	t.Run("All requests decode successfully", func(t *testing.T) {
		depositRequestBytes, err := hexutil.Decode("0x610000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"620000000000000000000000000000000000000000000000000000000000000000" +
			"4059730700000063000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		withdrawalRequestBytes, err := hexutil.Decode("0x6400000000000000000000000000000000000000" +
			"6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040597307000000")
		require.NoError(t, err)
		consolidationRequestBytes, err := hexutil.Decode("0x6600000000000000000000000000000000000000" +
			"670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{append([]byte{uint8(enginev1.DepositRequestType)}, depositRequestBytes...),
				append([]byte{uint8(enginev1.WithdrawalRequestType)}, withdrawalRequestBytes...),
				append([]byte{uint8(enginev1.ConsolidationRequestType)}, consolidationRequestBytes...)},
		}
		requests, err := types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.NoError(t, err)
		require.Len(t, requests.Deposits, 1)
		require.Len(t, requests.Withdrawals, 1)
		require.Len(t, requests.Consolidations, 1)
	})
}
