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

package types_test

import (
	"fmt"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/constraints"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip7685"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
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
			var recomputedER types.ExecutionRequests
			err = constraints.SSZUnmarshal(prysmERBytes, &recomputedER)
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
			var er types.ExecutionRequests
			// Ensure that calling UnmarshalSSZ with an invalid payload does not panic
			// and returns a non-nil error.
			require.NotPanics(t, func() {
				err := constraints.SSZUnmarshal(payload, &er)
				require.Error(t, err, "Expected error for payload %v", payload)
			})
		})
	}
}

// All tests below are adapted from Prysm
// https://github.com/prysmaticlabs/prysm/blob/e0e735470809df29c5404f64102ffbae5a574e0a/proto/engine/v1/electra_test.go#L13-L240

var depositRequestsSSZHex = "0x706b000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
	"0000000000077630000000000000000000000000000000000000000000000000000000000007b00000000000000736967000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
	"00000000000000000000000000000000000000000000000000000000000c801000000000000706b000000000000000000000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000077630000000000000000000000000000000000000000000000000" +
	"000000000009001000000000000736967000000000000000000000000000000000000000000000000000000000000000000000000000000" +
	"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020" +
	"00000000000000"

func TestDecodeExecutionRequests(t *testing.T) {
	t.Parallel()
	t.Run("All requests decode successfully", func(t *testing.T) {
		depositRequestBytes, err := hexutil.Decode("0x610000000000000000000000000000000000000000000000000000000" +
			"000000000000000000000000000000000000000" +
			"620000000000000000000000000000000000000000000000000000000000000000" +
			"40597307000000630000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000" +
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
	t.Run("Excluded requests still decode successfully when one request is missing", func(t *testing.T) {
		depositRequestBytes, err := hexutil.Decode("0x610000000000000000000000000000000000000000000000000000000" +
			"000000000000000000000000000000000000000" +
			"620000000000000000000000000000000000000000000000000000000000000000" +
			"405973070000006300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		consolidationRequestBytes, err := hexutil.Decode("0x6600000000000000000000000000000000000000" +
			"670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.DepositRequestType)}, depositRequestBytes...),
				append([]byte{uint8(enginev1.ConsolidationRequestType)}, consolidationRequestBytes...),
			},
		}
		requests, err := types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.NoError(t, err)
		require.Len(t, requests.Deposits, 1)
		require.Empty(t, requests.Withdrawals)
		require.Len(t, requests.Consolidations, 1)
	})
	t.Run("Decode execution requests should fail if ordering is not sorted", func(t *testing.T) {
		depositRequestBytes, err := hexutil.Decode("0x61000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000" +
			"620000000000000000000000000000000000000000000000000000000000000000" +
			"405973070000006300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		consolidationRequestBytes, err := hexutil.Decode("0x6600000000000000000000000000000000000000" +
			"670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.ConsolidationRequestType)}, consolidationRequestBytes...),
				append([]byte{uint8(enginev1.DepositRequestType)}, depositRequestBytes...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "requests should be in sorted order and unique")
	})
	t.Run("Requests should error if the request type is shorter than 1 byte", func(t *testing.T) {
		consolidationRequestBytes, err := hexutil.Decode("0x6600000000000000000000000000000000000000" +
			"670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{}, []byte{}...),
				append([]byte{uint8(enginev1.ConsolidationRequestType)}, consolidationRequestBytes...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "invalid execution request, length less than 1")
	})
	t.Run("a duplicate request should fail", func(t *testing.T) {
		withdrawalRequestBytes, err := hexutil.Decode("0x6400000000000000000000000000000000000000" +
			"6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040597307000000")
		require.NoError(t, err)
		withdrawalRequestBytes2, err := hexutil.Decode("0x6400000000000000000000000000000000000000" +
			"6500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040597307000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.WithdrawalRequestType)}, withdrawalRequestBytes...),
				append([]byte{uint8(enginev1.WithdrawalRequestType)}, withdrawalRequestBytes2...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "requests should be in sorted order and unique")
	})
	t.Run("a duplicate withdrawals ( non 0 request type )request should fail", func(t *testing.T) {
		depositRequestBytes, err := hexutil.Decode("0x61000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000" +
			"620000000000000000000000000000000000000000000000000000000000000000" +
			"4059730700000063000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		depositRequestBytes2, err := hexutil.Decode("0x61000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000" +
			"620000000000000000000000000000000000000000000000000000000000000000" +
			"405973070000006300000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.DepositRequestType)}, depositRequestBytes...),
				append([]byte{uint8(enginev1.DepositRequestType)}, depositRequestBytes2...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "requests should be in sorted order and unique")
	})
	t.Run("If a request type is provided, but the request list is shorter than the ssz of 1 request we error", func(t *testing.T) {
		consolidationRequestBytes, err := hexutil.Decode("0x6600000000000000000000000000000000000000" +
			"670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
			"680000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.DepositRequestType)}, []byte{}...),
				append([]byte{uint8(enginev1.ConsolidationRequestType)}, consolidationRequestBytes...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "invalid deposit requests SSZ size")
	})
	t.Run("If deposit requests are over the max allowed per payload then we should error", func(t *testing.T) {
		requests := make([]*enginev1.DepositRequest, types.MaxDepositRequestsPerPayload+1)
		for i := range requests {
			requests[i] = &enginev1.DepositRequest{
				Pubkey:                bytesutil.PadTo([]byte("pk"), 48),
				WithdrawalCredentials: bytesutil.PadTo([]byte("wc"), 32),
				Amount:                123,
				Signature:             bytesutil.PadTo([]byte("sig"), 96),
				Index:                 456,
			}
		}
		by, err := eip7685.MarshalItems(requests)
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.DepositRequestType)}, by...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "invalid deposit requests SSZ size, requests should not be more than the max per payload")
	})
	t.Run("If withdrawal requests are over the max allowed per payload then we should error", func(t *testing.T) {
		requests := make([]*enginev1.WithdrawalRequest, params.BeaconConfig().MaxWithdrawalRequestsPerPayload+1)
		for i := range requests {
			requests[i] = &enginev1.WithdrawalRequest{
				SourceAddress:   bytesutil.PadTo([]byte("sa"), 20),
				ValidatorPubkey: bytesutil.PadTo([]byte("pk"), 48),
				Amount:          55555,
			}
		}
		by, err := eip7685.MarshalItems(requests)
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.WithdrawalRequestType)}, by...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "invalid withdrawal requests SSZ size, requests should not be more than the max per payload")
	})
	t.Run("If consolidation requests are over the max allowed per payload then we should error", func(t *testing.T) {
		requests := make([]*enginev1.ConsolidationRequest, params.BeaconConfig().MaxConsolidationsRequestsPerPayload+1)
		for i := range requests {
			requests[i] = &enginev1.ConsolidationRequest{
				SourceAddress: bytesutil.PadTo([]byte("sa"), 20),
				SourcePubkey:  bytesutil.PadTo([]byte("pk"), 48),
				TargetPubkey:  bytesutil.PadTo([]byte("pk"), 48),
			}
		}
		by, err := eip7685.MarshalItems(requests)
		require.NoError(t, err)
		ebe := &enginev1.ExecutionBundleElectra{
			ExecutionRequests: [][]byte{
				append([]byte{uint8(enginev1.ConsolidationRequestType)}, by...),
			},
		}
		_, err = types.DecodeExecutionRequests(ebe.GetExecutionRequests())
		require.ErrorContains(t, err, "invalid consolidation requests SSZ size, requests should not be more than the max per payload")
	})
}

func TestGetExecutionRequestsList(t *testing.T) {
	t.Parallel()
	t.Run("Empty execution requests should return an empty response and not nil", func(t *testing.T) {
		ebe := &types.ExecutionRequests{}
		b, err := types.GetExecutionRequestsList(ebe)
		require.NoError(t, err)
		require.NotNil(t, b)
		require.Empty(t, b)
	})
}

func TestUnmarshalItems_OK(t *testing.T) {
	t.Parallel()
	drb, err := hexutil.Decode(depositRequestsSSZHex)
	require.NoError(t, err)
	exampleRequest := &types.DepositRequest{}
	depositRequests, err := eip7685.UnmarshalItems(
		drb,
		int(exampleRequest.SizeSSZ(nil)),
		func() *types.DepositRequest { return &types.DepositRequest{} })
	require.NoError(t, err)

	exampleRequest1 := &types.DepositRequest{
		Pubkey:      crypto.BLSPubkey(bytesutil.PadTo([]byte("pk"), 48)),
		Credentials: types.WithdrawalCredentials(bytesutil.PadTo([]byte("wc"), 32)),
		Amount:      123,
		Signature:   crypto.BLSSignature(bytesutil.PadTo([]byte("sig"), 96)),
		Index:       456,
	}
	exampleRequest2 := &types.DepositRequest{
		Pubkey:      crypto.BLSPubkey(bytesutil.PadTo([]byte("pk"), 48)),
		Credentials: types.WithdrawalCredentials(bytesutil.PadTo([]byte("wc"), 32)),
		Amount:      400,
		Signature:   crypto.BLSSignature(bytesutil.PadTo([]byte("sig"), 96)),
		Index:       32,
	}
	require.Equal(t, []*types.DepositRequest{exampleRequest1, exampleRequest2}, depositRequests)
}

func TestMarshalItems_OK(t *testing.T) {
	t.Parallel()
	exampleRequest1 := &types.DepositRequest{
		Pubkey:      crypto.BLSPubkey(bytesutil.PadTo([]byte("pk"), 48)),
		Credentials: types.WithdrawalCredentials(bytesutil.PadTo([]byte("wc"), 32)),
		Amount:      123,
		Signature:   crypto.BLSSignature(bytesutil.PadTo([]byte("sig"), 96)),
		Index:       456,
	}
	exampleRequest2 := &types.DepositRequest{
		Pubkey:      crypto.BLSPubkey(bytesutil.PadTo([]byte("pk"), 48)),
		Credentials: types.WithdrawalCredentials(bytesutil.PadTo([]byte("wc"), 32)),
		Amount:      400,
		Signature:   crypto.BLSSignature(bytesutil.PadTo([]byte("sig"), 96)),
		Index:       32,
	}
	drbs, err := eip7685.MarshalItems([]*types.DepositRequest{exampleRequest1, exampleRequest2})
	require.NoError(t, err)
	require.Equal(t, depositRequestsSSZHex, hexutil.Encode(drbs))
}
