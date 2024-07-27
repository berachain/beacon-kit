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
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

func TestBeaconBlockHeader_Serialization(t *testing.T) {
	original := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	data, err := original.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, data)
	var unmarshalled types.BeaconBlockHeader
	err = unmarshalled.UnmarshalSSZ(data)
	require.NoError(t, err)
	require.Equal(t, original, &unmarshalled)
}

func TestBeaconBlockHeader_SizeSSZ(t *testing.T) {
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	size := header.SizeSSZ()
	require.Equal(t, uint32(112), size)
}

func TestBeaconBlockHeader_HashTreeRoot(t *testing.T) {
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	_, err := header.HashTreeRoot()
	require.NoError(t, err)
}

func TestBeaconBlockHeader_GetTree(t *testing.T) {
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	tree, err := header.GetTree()

	require.NoError(t, err)
	require.NotNil(t, tree)
}

func TestBeaconBlockHeader_SetStateRoot(t *testing.T) {
	header := types.NewBeaconBlockHeader(
		math.Slot(100),
		math.ValidatorIndex(200),
		common.Root{},
		common.Root{},
		common.Root{},
	)

	newStateRoot := common.Root{}
	header.SetStateRoot(newStateRoot)

	require.Equal(t, newStateRoot, header.GetStateRoot())
}

func TestBeaconBlockHeader_New(t *testing.T) {
	slot := math.Slot(100)
	proposerIndex := math.ValidatorIndex(200)
	parentBlockRoot := common.Root{}
	stateRoot := common.Root{}
	bodyRoot := common.Root{}

	header := types.NewBeaconBlockHeader(
		slot,
		proposerIndex,
		parentBlockRoot,
		stateRoot,
		bodyRoot,
	)

	newHeader := header.New(
		slot,
		proposerIndex,
		parentBlockRoot,
		stateRoot,
		bodyRoot,
	)
	require.Equal(t, slot, newHeader.GetSlot())
	require.Equal(t, proposerIndex, newHeader.GetProposerIndex())
	require.Equal(t, parentBlockRoot, newHeader.GetParentBlockRoot())
	require.Equal(t, stateRoot, newHeader.GetStateRoot())
	require.Equal(t, bodyRoot, newHeader.BodyRoot)
}

func TestBeaconBlockHeader_UnmarshalSSZ_ErrSize(t *testing.T) {
	header := &types.BeaconBlockHeader{}
	buf := make([]byte, 100) // Incorrect size

	err := header.UnmarshalSSZ(buf)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
