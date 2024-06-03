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
	require.Equal(t, 112, size)
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
	require.ErrorIs(t, err, ssz.ErrSize)
}
