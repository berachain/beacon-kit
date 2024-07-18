// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
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
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	ssz "github.com/ferranbt/fastssz"
	"github.com/stretchr/testify/require"
)

func generateBeaconBlockBodyDenebPlus() *types.BeaconBlockBodyDenebPlus {
	var byteArray [256]byte
	byteSlice := byteArray[:]
	return &types.BeaconBlockBodyDenebPlus{
		BeaconBlockBodyBase: types.BeaconBlockBodyBase{
			RandaoReveal: [96]byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
				20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36,
				37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53,
				54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70,
				71, 72, 73, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87,
				88, 89, 90, 91, 92, 93, 94, 95, 96},
			Eth1Data: &types.Eth1Data{
				DepositRoot: common.Root{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
					0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D,
					0x1E, 0x1F,
				},
				DepositCount: 12345,
				BlockHash: gethprimitives.ExecutionHash{
					0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09,
					0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13,
					0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D,
					0x1E, 0x1F,
				},
			},
			Graffiti: [32]byte{
				1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19,
				20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32,
			},
			Deposits: []*types.Deposit{},
		},
		ExecutionPayload: &types.ExecutableDataDeneb{
			LogsBloom:    byteSlice,
			ExtraData:    make([]byte, 0),
			Transactions: make([][]byte, 0),
			Withdrawals:  make([]*engineprimitives.Withdrawal, 0),
		},
		Attestations:       []*types.AttestationData{},
		SlashingInfo:       []*types.SlashingInfo{},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
}

func TestBeaconBlockBodyDenebPlus_MarshalSSZ_UnmarshalSSZ(t *testing.T) {
	var byteArray [256]byte
	byteSlice := byteArray[:]
	testCases := []struct {
		name     string
		data     *types.BeaconBlockBodyDenebPlus
		expected *types.BeaconBlockBodyDenebPlus
		err      error
	}{
		{
			name:     "Valid BeaconBlockBodyDenebPlus",
			data:     generateBeaconBlockBodyDenebPlus(),
			expected: generateBeaconBlockBodyDenebPlus(),
			err:      nil,
		},
		{
			name: "Empty BeaconBlockBodyDenebPlus",
			data: &types.BeaconBlockBodyDenebPlus{
				BeaconBlockBodyBase: types.BeaconBlockBodyBase{
					RandaoReveal: [96]byte{},
					Eth1Data:     &types.Eth1Data{},
					Graffiti:     [32]byte{},
					Deposits:     []*types.Deposit{},
				},
				ExecutionPayload: &types.ExecutableDataDeneb{
					LogsBloom:    byteSlice,
					ExtraData:    make([]byte, 0),
					Transactions: make([][]byte, 0),
					Withdrawals:  make([]*engineprimitives.Withdrawal, 0),
				},
				Attestations:       []*types.AttestationData{},
				SlashingInfo:       []*types.SlashingInfo{},
				BlobKzgCommitments: []eip4844.KZGCommitment{},
			},
			expected: &types.BeaconBlockBodyDenebPlus{
				BeaconBlockBodyBase: types.BeaconBlockBodyBase{
					RandaoReveal: [96]byte{},
					Eth1Data:     &types.Eth1Data{},
					Graffiti:     [32]byte{},
					Deposits:     []*types.Deposit{},
				},
				ExecutionPayload: &types.ExecutableDataDeneb{
					LogsBloom:    byteSlice,
					ExtraData:    make([]byte, 0),
					Transactions: make([][]byte, 0),
					Withdrawals:  make([]*engineprimitives.Withdrawal, 0),
				},
				Attestations:       []*types.AttestationData{},
				SlashingInfo:       []*types.SlashingInfo{},
				BlobKzgCommitments: []eip4844.KZGCommitment{},
			},
			err: nil,
		},
		{
			name:     "Invalid Buffer Size",
			data:     generateBeaconBlockBodyDenebPlus(),
			expected: nil,
			err:      ssz.ErrSize,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := tc.data.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			var unmarshalled types.BeaconBlockBodyDenebPlus
			if tc.name == "Invalid Buffer Size" {
				err = unmarshalled.UnmarshalSSZ(data[:32])
				require.Error(t, err)
				require.Equal(t, tc.err, err)
			} else {
				err = unmarshalled.UnmarshalSSZ(data)
				require.NoError(t, err)
				require.Equal(t, tc.expected, &unmarshalled)
			}
		})
	}
}

func TestBeaconBlockBodyDenebPlus_GetTree(t *testing.T) {
	data := generateBeaconBlockBodyDenebPlus()

	tree, err := data.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)

	expectedRoot, err := data.HashTreeRoot()
	require.NoError(t, err)

	// Compare the tree root with the expected root
	actualRoot := tree.Hash()
	require.Equal(t, string(expectedRoot[:]), string(actualRoot))
}

// ... existing imports ...

func TestBeaconBlockBodyDenebPlus_IsNil(t *testing.T) {
	var body *types.BeaconBlockBodyDenebPlus
	require.True(t, body.IsNil())

	body = generateBeaconBlockBodyDenebPlus()
	require.False(t, body.IsNil())
}

func TestBeaconBlockBodyDenebPlus_GetExecutionPayload(t *testing.T) {
	body := generateBeaconBlockBodyDenebPlus()
	payload := body.GetExecutionPayload()
	require.NotNil(t, payload)
	require.Equal(t, body.ExecutionPayload, payload.InnerExecutionPayload)
}

func TestBeaconBlockBodyDenebPlus_SetExecutionData(t *testing.T) {
	body := generateBeaconBlockBodyDenebPlus()
	payload := &types.ExecutionPayload{
		InnerExecutionPayload: &types.ExecutableDataDeneb{},
	}

	err := body.SetExecutionData(payload)
	require.NoError(t, err)
	require.Equal(t, payload.InnerExecutionPayload, body.ExecutionPayload)

	invalidPayload := &types.ExecutionPayload{
		InnerExecutionPayload: nil,
	}
	err = body.SetExecutionData(invalidPayload)
	require.Error(t, err)
	require.Equal(t, "invalid execution data type", err.Error())
}

func TestBeaconBlockBodyDenebPlus_GetBlobKzgCommitments(t *testing.T) {
	body := generateBeaconBlockBodyDenebPlus()
	commitments := body.GetBlobKzgCommitments()
	require.NotNil(t, commitments)
	expected := eip4844.KZGCommitments[gethprimitives.ExecutionHash]{}
	require.Equal(t, expected, commitments)
}

func TestBeaconBlockBodyDenebPlus_SetBlobKzgCommitments(t *testing.T) {
	body := generateBeaconBlockBodyDenebPlus()
	commitments := eip4844.KZGCommitments[gethprimitives.ExecutionHash]{}
	body.SetBlobKzgCommitments(commitments)
	require.Equal(t, commitments, body.GetBlobKzgCommitments())
}

func TestBeaconBlockBodyDenebPlus_GetTopLevelRoots(t *testing.T) {
	body := generateBeaconBlockBodyDenebPlus()
	roots, err := body.GetTopLevelRoots()
	require.NoError(t, err)
	require.NotNil(t, roots)
	require.Equal(t, types.BodyLengthDeneb, uint64(len(roots)))
}

func TestBeaconBlockBodyDenebPlus_Length(t *testing.T) {
	body := generateBeaconBlockBodyDenebPlus()
	length := body.Length()
	require.Equal(t, types.BodyLengthDeneb, length)
}
