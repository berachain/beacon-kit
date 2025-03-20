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
	"github.com/berachain/beacon-kit/primitives/decoder"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

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
			err = decoder.SSZUnmarshal(prysmWithdrawalBytes, &recomputedWithdrawalRequest)
			require.NoError(t, err)

			// Compare that the original and recomputed values match.
			require.Equal(t, tc.withdrawalRequest, recomputedWithdrawalRequest)
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
				err = decoder.SSZUnmarshal(payload, &w)
				require.Error(t, err, "expected error for payload %v", payload)
			})
		})
	}
}
