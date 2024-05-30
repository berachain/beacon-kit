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

func TestForkData_Serialization(t *testing.T) {
	original := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)

	var unmarshalled types.ForkData
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)

	require.Equal(t, original, &unmarshalled)
}

func TestForkData_SizeSSZ(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}

	size := forkData.SizeSSZ()

	require.Equal(t, 36, size)
}

func TestForkData_HashTreeRoot(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}
	_, err := forkData.HashTreeRoot()

	require.NoError(t, err)
}

func TestForkData_GetTree(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}

	tree, err := forkData.GetTree()

	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestForkData_ComputeDomain(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}
	domainType := common.DomainType{
		0x01, 0x00, 0x00, 0x00,
	}
	_, err := forkData.ComputeDomain(domainType)
	require.NoError(t, err)
}

func TestNewForkData(t *testing.T) {
	currentVersion := common.Version{}
	genesisValidatorsRoot := common.Root{}

	forkData := types.NewForkData(currentVersion, genesisValidatorsRoot)

	require.Equal(t, currentVersion, forkData.CurrentVersion)
	require.Equal(t, genesisValidatorsRoot, forkData.GenesisValidatorsRoot)
}

func TestNew(t *testing.T) {
	currentVersion := common.Version{}
	genesisValidatorsRoot := common.Root{}
	forkData := &types.ForkData{}

	newForkData := forkData.New(currentVersion, genesisValidatorsRoot)

	require.Equal(t, currentVersion, newForkData.CurrentVersion)
	require.Equal(t, genesisValidatorsRoot, newForkData.GenesisValidatorsRoot)
}
