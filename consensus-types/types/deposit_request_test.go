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
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/encoding/sszutil"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	"github.com/stretchr/testify/require"
)

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
			depositHTR, err := tc.depositRequest.HashTreeRoot()
			require.NoError(t, err)
			// Compare the HashTreeRoots to ensure all fields were correctly interpreted.
			require.Equal(t, depositHTR[:], prysmHTR[:])

			// Marshal the Prysm deposit request.
			prysmDepositBytes, err := prysmDeposit.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into a new DepositRequest.
			var recomputedDepositRequest types.DepositRequest
			err = sszutil.Unmarshal(prysmDepositBytes, &recomputedDepositRequest)
			require.NoError(t, err)

			// Compare that the original and recomputed deposit requests match.
			require.Equal(t, *tc.depositRequest, recomputedDepositRequest)
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
				var d types.DepositRequest
				err = sszutil.Unmarshal(payload, &d)
				// We expect an error for every invalid payload.
				require.Error(t, err, "expected error for payload %v", payload)
			})
		})
	}
}

func TestDepositRequests_ValidValuesSSZ(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		depositRequest types.DepositRequests
	}{
		{
			name: "basic",
			depositRequest: []*types.DepositRequest{
				{
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
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Marshal the original deposit request.
			depositRequestBytes, err := tc.depositRequest.MarshalSSZ()
			require.NoError(t, err)

			// Unmarshal back into a new DepositRequest.
			recomputedDepositRequest, err := types.DecodeDepositRequests(depositRequestBytes)
			require.NoError(t, err)

			// Compare that the original and recomputed deposit requests match.
			require.Equal(t, tc.depositRequest, recomputedDepositRequest)
		})
	}
}

// Tests below are adapted from Prysm
// https://github.com/prysmaticlabs/prysm/blob/develop/proto/engine/v1/electra_test.go#L198-L240

const depositRequestsSSZHex = "0x706b000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
	"0000000000077630000000000000000000000000000000000000000000000000000000000007b00000000000000736967000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" +
	"00000000000000000000000000000000000000000000000000000000000c801000000000000706b000000000000000000000000000000000" +
	"0000000000000000000000000000000000000000000000000000000000077630000000000000000000000000000000000000000000000000" +
	"000000000009001000000000000736967000000000000000000000000000000000000000000000000000000000000000000000000000000" +
	"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000020" +
	"00000000000000"

func TestUnmarshalItems_OK(t *testing.T) {
	t.Parallel()
	drb, err := hexutil.Decode(depositRequestsSSZHex)
	require.NoError(t, err)
	exampleRequest := &types.DepositRequest{}
	depositRequests, err := sszutil.UnmarshalItemsEIP7685(
		drb,
		exampleRequest.SizeSSZ(),
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
	drbs, err := sszutil.MarshalItemsEIP7685([]*types.DepositRequest{exampleRequest1, exampleRequest2})
	require.NoError(t, err)
	require.Equal(t, depositRequestsSSZHex, hexutil.Encode(drbs))
}
