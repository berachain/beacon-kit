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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestEth1Data_Serialization(t *testing.T) {
	original := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    common.ExecutionHash{},
	}

	data, err := original.MarshalSSZ()
	require.NoError(t, err)

	var unmarshalled types.Eth1Data
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)

	// The original and unmarshalled Eth1Data should be the same
	require.Equal(t, original, &unmarshalled)
}

func TestEth1Data_SizeSSZ(t *testing.T) {
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    common.ExecutionHash{},
	}

	// Get the SSZ size of the Eth1Data
	size := eth1Data.SizeSSZ()

	// The size should be 72
	require.Equal(t, 72, size)
}

func TestEth1Data_HashTreeRoot(t *testing.T) {
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    common.ExecutionHash{},
	}

	// Get the hash tree root of the Eth1Data
	_, err := eth1Data.HashTreeRoot()

	// There should be no error
	require.NoError(t, err)
}

func TestEth1Data_GetTree(t *testing.T) {
	// Create an Eth1Data
	eth1Data := &types.Eth1Data{
		DepositRoot:  common.Root{},
		DepositCount: 10,
		BlockHash:    common.ExecutionHash{},
	}

	// Get the SSZ tree of the Eth1Data
	tree, err := eth1Data.GetTree()

	// There should be no error
	require.NoError(t, err)

	// The tree should not be nil
	require.NotNil(t, tree)
}
