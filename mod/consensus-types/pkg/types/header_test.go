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
	ssz "github.com/ferranbt/fastssz"
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

func TestBeaconBlockHeaderBase_MarshalSSZUnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name   string
		header *types.BeaconBlockHeaderBase
		valid  bool
	}{
		{
			name: "normal case",
			header: &types.BeaconBlockHeaderBase{
				Slot:            100,
				ProposerIndex:   200,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
			valid: true,
		},
		{
			name: "zero values",
			header: &types.BeaconBlockHeaderBase{
				Slot:            0,
				ProposerIndex:   0,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
			valid: true,
		},
		{
			name: "invalid size",
			header: &types.BeaconBlockHeaderBase{
				Slot:            100,
				ProposerIndex:   200,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.header.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			var unmarshalled types.BeaconBlockHeaderBase
			if tt.valid {
				err = unmarshalled.UnmarshalSSZ(data)
				require.NoError(t, err)
				require.Equal(t, tt.header, &unmarshalled)
			} else {
				// Modify data to simulate invalid size
				data = data[:len(data)-1]
				err = unmarshalled.UnmarshalSSZ(data)
				require.ErrorIs(t, err, ssz.ErrSize)
			}
		})
	}
}

func TestBeaconBlockHeaderBase_MarshalSSZToUnmarshalSSZ(t *testing.T) {
	tests := []struct {
		name   string
		header *types.BeaconBlockHeaderBase
	}{
		{
			name: "normal case",
			header: &types.BeaconBlockHeaderBase{
				Slot:            100,
				ProposerIndex:   200,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
		{
			name: "zero values",
			header: &types.BeaconBlockHeaderBase{
				Slot:            0,
				ProposerIndex:   0,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 0, tt.header.SizeSSZ())
			data, err := tt.header.MarshalSSZTo(buf)
			require.NoError(t, err)
			require.NotNil(t, data)

			var unmarshalled types.BeaconBlockHeaderBase
			err = unmarshalled.UnmarshalSSZ(data)
			require.NoError(t, err)
			require.Equal(t, tt.header, &unmarshalled)
		})
	}
}

func TestBeaconBlockHeaderBase_HashTreeRoot(t *testing.T) {
	tests := []struct {
		name   string
		header *types.BeaconBlockHeaderBase
	}{
		{
			name: "HashTreeRoot normal case",
			header: &types.BeaconBlockHeaderBase{
				Slot:            100,
				ProposerIndex:   200,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
		{
			name: "HashTreeRoot zero values",
			header: &types.BeaconBlockHeaderBase{
				Slot:            0,
				ProposerIndex:   0,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := tt.header.HashTreeRoot()
			require.NoError(t, err)
			require.NotNil(t, root)
		})
	}
}

func TestBeaconBlockHeaderBase_HashTreeRootWith(t *testing.T) {
	tests := []struct {
		name   string
		header *types.BeaconBlockHeaderBase
	}{
		{
			name: "HashTreeRootWith normal case",
			header: &types.BeaconBlockHeaderBase{
				Slot:            100,
				ProposerIndex:   200,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
		{
			name: "HashTreeRootWith zero values",
			header: &types.BeaconBlockHeaderBase{
				Slot:            0,
				ProposerIndex:   0,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hh := ssz.NewHasher()
			err := tt.header.HashTreeRootWith(hh)
			require.NoError(t, err)
			require.NotNil(t, hh.Hash())
		})
	}
}

func TestBeaconBlockHeaderBase_GetTree(t *testing.T) {
	tests := []struct {
		name   string
		header *types.BeaconBlockHeaderBase
	}{
		{
			name: "GetTree normal case",
			header: &types.BeaconBlockHeaderBase{
				Slot:            100,
				ProposerIndex:   200,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
		{
			name: "GetTree zero values",
			header: &types.BeaconBlockHeaderBase{
				Slot:            0,
				ProposerIndex:   0,
				ParentBlockRoot: common.Root{},
				StateRoot:       common.Root{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := tt.header.GetTree()
			require.NoError(t, err)
			require.NotNil(t, tree)
		})
	}
}
