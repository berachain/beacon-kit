// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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
	"testing/quick"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/berachain/beacon-kit/testing/utils"

	"github.com/stretchr/testify/require"
)

func TestBeaconBlockForDeneb(t *testing.T) {
	t.Parallel()
	deneb1 := version.Deneb1()
	block, err := types.NewBeaconBlockWithVersion(
		math.Slot(10),
		math.ValidatorIndex(5),
		common.Root{1, 2, 3, 4, 5}, // parent root
		deneb1,
	)
	require.NoError(t, err)
	require.NotNil(t, block)
	require.Equal(t, deneb1, block.GetForkVersion())
	require.Equal(t, deneb1, block.GetBody().GetForkVersion())
	require.Equal(t, deneb1, block.GetBody().GetExecutionPayload().GetForkVersion())
}

func TestBeaconBlock(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block := utils.GenerateValidBeaconBlock(t, v)

		require.NotNil(t, block.Body)
		require.Equal(t, math.U64(10), block.GetTimestamp())
		require.Equal(t, v, block.GetForkVersion())
		require.NotNil(t, block)

		// Set a new state root and test the SetStateRoot and GetBody methods
		newStateRoot := [32]byte{1, 1, 1, 1, 1}
		block.SetStateRoot(newStateRoot)
		require.Equal(t, newStateRoot, [32]byte(block.StateRoot))

		// Test the GetHeader method
		header := block.GetHeader()
		require.NotNil(t, header)
		require.Equal(t, block.Slot, header.Slot)
		require.Equal(t, block.ProposerIndex, header.ProposerIndex)
		require.Equal(t, block.ParentRoot, header.ParentBlockRoot)
		require.Equal(t, block.StateRoot, header.StateRoot)
		require.Equal(t, newStateRoot, [32]byte(block.GetStateRoot()))
	})
}

func TestBeaconBlock_MarshalUnmarshalSSZ(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block := utils.GenerateValidBeaconBlock(t, v)

		sszBlock, err := block.MarshalSSZ()
		require.NoError(t, err)
		require.NotNil(t, sszBlock)

		unmarshalledBlock := types.NewEmptyBeaconBlockWithVersion(v)
		err = ssz.Unmarshal(sszBlock, unmarshalledBlock)
		require.NoError(t, err)
		require.Equal(t, block, unmarshalledBlock)
	})
}

func TestBeaconBlock_HashTreeRoot(t *testing.T) {
	t.Parallel()
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block := utils.GenerateValidBeaconBlock(t, v)
		hashRoot := block.HashTreeRoot()
		require.NotNil(t, hashRoot)
	})
}

func TestBeaconBlock_IsNil(t *testing.T) {
	t.Parallel()
	var block *types.BeaconBlock
	require.Nil(t, block)
}

func TestNewWithVersion(t *testing.T) {
	t.Parallel()
	slot := math.Slot(10)
	proposerIndex := math.ValidatorIndex(5)
	parentBlockRoot := common.Root{1, 2, 3, 4, 5}

	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		block, err := types.NewBeaconBlockWithVersion(
			slot, proposerIndex, parentBlockRoot, v,
		)
		require.NoError(t, err)
		require.NotNil(t, block)

		// Check the block's fields
		require.NotNil(t, block)
		require.Equal(t, slot, block.GetSlot())
		require.Equal(t, proposerIndex, block.GetProposerIndex())
		require.Equal(t, parentBlockRoot, block.GetParentBlockRoot())
		require.Equal(t, v, block.GetForkVersion())
	})
}

func TestNewWithVersionInvalidForkVersion(t *testing.T) {
	t.Parallel()
	slot := math.Slot(10)
	proposerIndex := math.ValidatorIndex(5)
	parentBlockRoot := common.Root{1, 2, 3, 4, 5}

	_, err := types.NewBeaconBlockWithVersion(
		slot,
		proposerIndex,
		parentBlockRoot,
		common.Version{100, 0, 0, 0},
	) // 100 is an invalid fork version
	require.ErrorIs(t, err, types.ErrForkVersionNotSupported)
}

func TestPropertyBlockRootAndBlockHeaderRootEquivalence(t *testing.T) {
	t.Parallel()
	qc := &quick.Config{MaxCount: 100}
	runForAllSupportedVersions(t, func(t *testing.T, v common.Version) {
		f := func(
			slot math.Slot,
			proposerIdx math.ValidatorIndex,
			parentBlockRoot common.Root,
		) bool {
			blk, err := types.NewBeaconBlockWithVersion(
				slot,
				proposerIdx,
				parentBlockRoot,
				v,
			)
			require.NoError(t, err)
			return blk.GetHeader().HashTreeRoot().Equals(blk.HashTreeRoot())
		}
		require.NoError(t, quick.Check(f, qc))
	})
}
