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

	types "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	karalabessz "github.com/karalabe/ssz"
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

func TestForkData_Unmarshal(t *testing.T) {
	var unmarshalled types.ForkData
	err := unmarshalled.UnmarshalSSZ([]byte{})
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}

func TestForkData_SizeSSZ(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}

	size := karalabessz.Size(forkData)
	require.Equal(t, uint32(36), size)
}

func TestForkData_HashTreeRoot(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}
	require.NotPanics(t, func() {
		_ = forkData.HashTreeRoot()
	})
}

func TestForkData_ComputeDomain(t *testing.T) {
	forkData := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}
	domainType := common.DomainType{
		0x01, 0x00, 0x00, 0x00,
	}
	require.NotPanics(t, func() {
		_ = forkData.ComputeDomain(domainType)
	})
}

func TestForkData_ComputeRandaoSigningRoot(t *testing.T) {
	fd := &types.ForkData{
		CurrentVersion:        common.Version{},
		GenesisValidatorsRoot: common.Root{},
	}

	domainType := common.DomainType{0, 0, 0, 0}
	epoch := math.Epoch(1)

	require.NotPanics(t, func() {
		fd.ComputeRandaoSigningRoot(domainType, epoch)
	})
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
