// SPDX-License-Identifier: MIT
//
// Copyright (c) 2024 Berachain Foundation
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
// WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
// OTHER DEALINGS IN THE SOFTWARE.

package types_test

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/crypto"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/stretchr/testify/require"
)

func generateBeaconBlockBodyDeneb() types.BeaconBlockBodyDeneb {
	var byteArray [256]byte
	byteSlice := byteArray[:]
	return types.BeaconBlockBodyDeneb{
		BeaconBlockBodyBase: types.BeaconBlockBodyBase{
			RandaoReveal: [96]byte{1, 2, 3},
			Eth1Data:     &types.Eth1Data{},
			Graffiti:     [32]byte{4, 5, 6},
			Deposits:     []*types.Deposit{},
		},
		ExecutionPayload: &types.ExecutableDataDeneb{
			LogsBloom: byteSlice,
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
		ExecutionPayload:   &types.ExecutableDataDeneb{},
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

func TestBeaconBlockBodyDeneb_SetExecutionData_Error(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{}
	executionData := &types.ExecutionPayload{}
	err := body.SetExecutionData(executionData)

	require.ErrorContains(t, err, "invalid execution data type")
}

func TestBeaconBlockBodyDeneb_SetBlobKzgCommitments(t *testing.T) {
	body := types.BeaconBlockBodyDeneb{}
	commitments := eip4844.KZGCommitments[common.ExecutionHash]{}
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

func TestBeaconBlockBodyDeneb_GetTopLevelRoots(t *testing.T) {
	body := generateBeaconBlockBodyDeneb()
	roots, err := body.GetTopLevelRoots()
	require.NoError(t, err)
	require.NotNil(t, roots)
}
