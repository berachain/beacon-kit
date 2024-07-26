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

func generateBeaconBlockBodyDeneb() types.BeaconBlockBodyDeneb {
	var byteArray [256]byte
	return types.BeaconBlockBodyDeneb{
		BeaconBlockBodyBase: types.BeaconBlockBodyBase{
			RandaoReveal: [96]byte{1, 2, 3},
			Eth1Data:     &types.Eth1Data{},
			Graffiti:     [32]byte{4, 5, 6},
			Deposits:     []*types.Deposit{},
		},
		ExecutionPayload: &types.ExecutionPayload{
			LogsBloom: byteArray,
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
}

func TestBeaconBlockBodyBase(t *testing.T) {
	body := types.BeaconBlockBodyBase{
		RandaoReveal: [96]byte{1, 2, 3},
		Eth1Data:     &types.Eth1Data{},
		Graffiti:     [32]byte{4, 5, 6},
		Deposits:     []*types.Deposit{},
	}

	require.Equal(t, bytes.B96{1, 2, 3}, body.GetRandaoReveal())
	require.NotNil(t, body.GetEth1Data())
	require.Equal(t, bytes.B32{4, 5, 6}, body.GetGraffiti())
	require.NotNil(t, body.GetDeposits())
}

func TestBeaconBlockBodyDeneb(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{
		BeaconBlockBodyBase: types.BeaconBlockBodyBase{
			RandaoReveal: [96]byte{1, 2, 3},
			Eth1Data:     &types.Eth1Data{},
			Graffiti:     [32]byte{4, 5, 6},
			Deposits:     []*types.Deposit{},
		},
		ExecutionPayload:   &types.ExecutionPayload{},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}

	require.False(t, body.IsNil())
	require.NotNil(t, body.GetExecutionPayload())
	require.NotNil(t, body.GetBlobKzgCommitments())
	require.Equal(t, types.BodyLengthDeneb, body.Length())
}

func TestBeaconBlockBodyDeneb_GetTree(t *testing.T) {
	body := generateBeaconBlockBodyDeneb()
	tree, err := body.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestBeaconBlockBodyDeneb_SetBlobKzgCommitments(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{}
	commitments := eip4844.KZGCommitments[gethprimitives.ExecutionHash]{}
	body.SetBlobKzgCommitments(commitments)

	require.Equal(t, commitments, body.GetBlobKzgCommitments())
}

func TestBeaconBlockBodyDeneb_SetRandaoReveal(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{}
	randaoReveal := crypto.BLSSignature{1, 2, 3}
	body.SetRandaoReveal(randaoReveal)

	require.Equal(t, randaoReveal, body.GetRandaoReveal())
}

func TestBeaconBlockBodyDeneb_SetEth1Data(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{}
	eth1Data := &types.Eth1Data{}
	body.SetEth1Data(eth1Data)

	require.Equal(t, eth1Data, body.GetEth1Data())
}

func TestBeaconBlockBodyDeneb_SetDeposits(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{}
	deposits := []*types.Deposit{}
	body.SetDeposits(deposits)

	require.Equal(t, deposits, body.GetDeposits())
}

func TestBeaconBlockBodyDeneb_MarshalSSZ(t *testing.T) {
	var byteArray [256]byte
	body := types.BeaconBlockBodyDeneb{
		BeaconBlockBodyBase: types.BeaconBlockBodyBase{
			RandaoReveal: [96]byte{1, 2, 3},
			Eth1Data:     &types.Eth1Data{},
			Graffiti:     [32]byte{4, 5, 6},
			Deposits:     []*types.Deposit{},
		},
		ExecutionPayload: &types.ExecutionPayload{
			LogsBloom: byteArray,
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{},
	}
	data, err := body.MarshalSSZ()

	require.NoError(t, err)
	require.NotNil(t, data)
}
func TestBeaconBlockBodyDeneb_GetTopLevelRoots(t *testing.T) {
	body := generateBeaconBlockBodyDeneb()
	roots, err := body.GetTopLevelRoots()
	require.NoError(t, err)
	require.NotNil(t, roots)
}

func TestBeaconBlockBody_Empty(t *testing.T) {
	blockBody := types.BeaconBlockBody{}
	body := blockBody.Empty(version.Deneb)
	require.NotNil(t, body)

	_, ok := body.RawBeaconBlockBody.(*types.BeaconBlockBodyDeneb)
	require.True(t, ok)
}
