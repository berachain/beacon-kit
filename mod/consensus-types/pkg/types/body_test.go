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
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

func generateBeaconBlockBody() types.BeaconBlockBody {
	return types.BeaconBlockBody{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data: &types.Eth1Data{
			DepositRoot:  [32]byte{7, 8, 9},
			DepositCount: 12345,
			BlockHash:    [32]byte{10, 11, 12},
		},
		Graffiti: [32]byte{4, 5, 6},
		Deposits: []*types.Deposit{
			{
				Pubkey:      [48]byte{16, 17, 18},
				Credentials: [32]byte{19, 20, 21},
				Amount:      1000,
				Signature:   [96]byte{22, 23, 24},
				Index:       1,
			},
		},
		ExecutionPayload: &types.ExecutionPayload{
			ParentHash:   [32]byte{25, 26, 27},
			FeeRecipient: [20]byte{28, 29, 30},
			StateRoot:    [32]byte{31, 32, 33},
			ReceiptsRoot: [32]byte{34, 35, 36},
			LogsBloom:    [256]byte{37, 38, 39},

			GasLimit:      8000000,
			GasUsed:       7500000,
			Timestamp:     1617181920,
			ExtraData:     []byte{43, 44, 45},
			BaseFeePerGas: [32]byte{46, 47, 48},
			BlockHash:     [32]byte{49, 50, 51},
			Transactions:  [][]byte{{52, 53, 54}},
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{
			{55, 56, 57},
			{58, 59, 60},
		},
	}
}

func TestBeaconBlockBodyBase(t *testing.T) {
	body := types.BeaconBlockBody{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &types.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*types.Deposit{},
	}

	require.Equal(t, bytes.B96{1, 2, 3}, body.GetRandaoReveal())
	require.NotNil(t, body.GetEth1Data())
	require.Equal(t, bytes.B32{4, 5, 6}, body.GetGraffiti())
	require.NotNil(t, body.GetDeposits())

	// Test SetExecutionPayload and GetExecutionPayload
	executionPayload := &types.ExecutionPayload{ExtraData: []byte{7, 8, 9}}
	body.SetExecutionPayload(executionPayload)
	require.Equal(t, executionPayload, body.GetExecutionPayload())

	// Test SetGraffiti and GetGraffiti
	newGraffiti := bytes.B32{10, 11, 12}
	body.SetGraffiti(newGraffiti)
	require.Equal(t, newGraffiti, body.GetGraffiti())
}

func TestBeaconBlockBody(t *testing.T) {
	body := types.BeaconBlockBody{
		RandaoReveal:       [96]byte{1, 2, 3},
		Eth1Data:           &types.Eth1Data{},
		Graffiti:           [32]byte{4, 5, 6},
		Deposits:           []*types.Deposit{},
		ExecutionPayload:   &types.ExecutionPayload{},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}

	require.False(t, body.IsNil())
	require.NotNil(t, body.GetExecutionPayload())
	require.NotNil(t, body.GetBlobKzgCommitments())
	require.Equal(t, types.BodyLengthDeneb, body.Length())
}

func TestBeaconBlockBody_GetTree(t *testing.T) {
	body := generateBeaconBlockBody()
	tree, err := body.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestBeaconBlockBody_SetBlobKzgCommitments(t *testing.T) {
	body := types.BeaconBlockBody{}
	commitments := eip4844.KZGCommitments[gethprimitives.ExecutionHash]{}
	body.SetBlobKzgCommitments(commitments)

	require.Equal(t, commitments, body.GetBlobKzgCommitments())
}

func TestBeaconBlockBody_SetRandaoReveal(t *testing.T) {
	body := types.BeaconBlockBody{}
	randaoReveal := crypto.BLSSignature{1, 2, 3}
	body.SetRandaoReveal(randaoReveal)

	require.Equal(t, randaoReveal, body.GetRandaoReveal())
}

func TestBeaconBlockBody_SetEth1Data(t *testing.T) {
	body := types.BeaconBlockBody{}
	eth1Data := &types.Eth1Data{}
	body.SetEth1Data(eth1Data)

	require.Equal(t, eth1Data, body.GetEth1Data())
}

func TestBeaconBlockBody_SetDeposits(t *testing.T) {
	body := types.BeaconBlockBody{}
	deposits := []*types.Deposit{}
	body.SetDeposits(deposits)

	require.Equal(t, deposits, body.GetDeposits())
}

func TestBeaconBlockBody_MarshalSSZ(t *testing.T) {
	body := types.BeaconBlockBody{
		RandaoReveal:       [96]byte{1, 2, 3},
		Eth1Data:           &types.Eth1Data{},
		Graffiti:           [32]byte{4, 5, 6},
		Deposits:           []*types.Deposit{},
		ExecutionPayload:   &types.ExecutionPayload{},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	data, err := body.MarshalSSZ()

	require.NoError(t, err)
	require.NotNil(t, data)
}
func TestBeaconBlockBody_GetTopLevelRoots(t *testing.T) {
	body := generateBeaconBlockBody()
	roots, err := body.GetTopLevelRoots()
	require.NoError(t, err)
	require.NotNil(t, roots)
}

func TestBeaconBlockBody_Empty(t *testing.T) {
	blockBody := types.BeaconBlockBody{}
	body := blockBody.Empty(version.Deneb)
	require.NotNil(t, body)
}

func TestBeaconBlockBody_MarshalSSZUnmarshalSSZ(t *testing.T) {
	originalBody := generateBeaconBlockBody()

	// Marshal the original body to SSZ
	data, err := originalBody.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	// Unmarshal the data back into a new body
	var newBody types.BeaconBlockBody
	err = newBody.UnmarshalSSZ(data)
	require.NoError(t, err)

	// Verify that the new body matches the original body
	require.Equal(t, originalBody, newBody)

	invalidData := []byte{0x00, 0x01, 0x02} // Invalid SSZ data
	var newBodyInvalid types.BeaconBlockBody
	err = newBodyInvalid.UnmarshalSSZ(invalidData)
	require.Error(t, err)
}

func TestBeaconBlockBody_MarshalSSZTo(t *testing.T) {
	originalBody := generateBeaconBlockBody()

	// Prepare a destination slice
	dst := make([]byte, 0)

	// Marshal the original body to SSZ and append to dst
	result, err := originalBody.MarshalSSZTo(dst)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify that the result is not empty and contains the marshaled data
	require.Greater(t, len(result), len(dst))
}
