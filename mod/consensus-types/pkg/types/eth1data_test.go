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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	ssz "github.com/ferranbt/fastssz"
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
	require.ErrorIs(t, err, ssz.ErrSize)
}

func TestEth1Data_SizeSSZ(t *testing.T) {
	eth1Data := (&types.Eth1Data{}).New(
		common.Root{},
		10,
		gethprimitives.ExecutionHash{},
	)

	size := eth1Data.SizeSSZ()
	require.Equal(t, 72, size)
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
