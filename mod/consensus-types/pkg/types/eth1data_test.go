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
	"io"
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	gethprimitives "github.com/berachain/beacon-kit/mod/geth-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestEth1Data_Serialization(t *testing.T) {
	original := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    gethprimitives.ExecutionHash{},
	}

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.Eth1Data
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original, &unmarshalled)
}

func TestEth1Data_UnmarshalError(t *testing.T) {
	var unmarshalled types.Eth1Data
	err := unmarshalled.UnmarshalSSZ([]byte{})
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestEth1Data_SizeSSZ(t *testing.T) {
	eth1Data := (&types.Eth1Data{}).New(
		common.Root{},
		10,
		gethprimitives.ExecutionHash{},
	)

	size := eth1Data.SizeSSZ()
	require.Equal(t, uint32(72), size)
}

func TestEth1Data_HashTreeRoot(t *testing.T) {
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    gethprimitives.ExecutionHash{},
	}

	_, err := eth1Data.HashTreeRoot()
	require.NoError(t, err)
}

func TestEth1Data_GetTree(t *testing.T) {
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    gethprimitives.ExecutionHash{},
	}

	tree, err := eth1Data.GetTree()

	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestEth1Data_GetDepositCount(t *testing.T) {
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    gethprimitives.ExecutionHash{},
	}

	count := eth1Data.GetDepositCount()

	require.Equal(t, uint64(10), count.Unwrap())
}

func TestEth1DataMarshalSSZTo(t *testing.T) {
	eth1Data := &types.Eth1Data{
		DepositRoot:  [32]byte{1, 2, 3},
		DepositCount: 12345,
		BlockHash:    [32]byte{4, 5, 6},
	}

	dst := make([]byte, 0, types.Eth1DataSize)
	result, err := eth1Data.MarshalSSZTo(dst)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, types.Eth1DataSize, len(result))

	// Unmarshal the result to verify correctness
	var unmarshaledEth1Data types.Eth1Data
	err = unmarshaledEth1Data.UnmarshalSSZ(result)
	require.NoError(t, err)
	require.Equal(t, eth1Data, &unmarshaledEth1Data)
}
