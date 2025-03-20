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
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/decoder"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

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
			err = decoder.SSZUnmarshal(prysmCRBytes, &recomputedCR)
			require.NoError(t, err)

			// Compare that the original and recomputed consolidation requests match.
			require.Equal(t, tc.consolidationRequest, recomputedCR)
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
				err = decoder.SSZUnmarshal(validBytes, &c)
				require.Error(t, err, "expected error for payload %v", payload)
			})
		})
	}
}
