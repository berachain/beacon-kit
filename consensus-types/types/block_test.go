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
	engineprimitives "github.com/berachain/beacon-kit/engine-primitives/engine-primitives"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// generateValidBeaconBlock generates a valid beacon block for the Deneb.
func generateValidBeaconBlock(t *testing.T) *types.BeaconBlock {
	t.Helper()

	// Initialize your block here
	deneb1 := version.Deneb1()
	beaconBlock, err := types.NewBeaconBlockWithVersion(
		math.Slot(10),
		math.ValidatorIndex(5),
		common.Root{1, 2, 3, 4, 5}, // parent block root
		deneb1,
	)
	require.NoError(t, err)

	beaconBlock.StateRoot = common.Root{5, 4, 3, 2, 1}
	beaconBlock.Body = &types.BeaconBlockBody{
		ExecutionPayload: &types.ExecutionPayload{
			Timestamp: 10,
			ExtraData: []byte("dummy extra data for testing"),
			Transactions: [][]byte{
				[]byte("tx1"),
				[]byte("tx2"),
				[]byte("tx3"),
			},
			Withdrawals: engineprimitives.Withdrawals{
				{Index: 0, Amount: 100},
				{Index: 1, Amount: 200},
			},
			BaseFeePerGas: math.NewU256(0),
			EpVersion:     beaconBlock.GetBody().GetExecutionPayload().Version(),
		},
		Eth1Data: &types.Eth1Data{},
		Deposits: []*types.Deposit{
			{
				Index: 1,
			},
		},
		BlobKzgCommitments: []eip4844.KZGCommitment{
			{1, 2, 3},
		},
	}

	body := beaconBlock.GetBody()
	body.SetVersion(beaconBlock.BbVersion)
	body.SetProposerSlashings(types.ProposerSlashings{})
	body.SetAttesterSlashings(types.AttesterSlashings{})
	body.SetAttestations(types.Attestations{})
	body.SetSyncAggregate(&types.SyncAggregate{})
	body.SetVoluntaryExits(types.VoluntaryExits{})
	body.SetBlsToExecutionChanges(types.BlsToExecutionChanges{})
	return beaconBlock
}

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
	require.Equal(t, deneb1, block.Version())
}

func TestBeaconBlock(t *testing.T) {
	t.Parallel()
	block := generateValidBeaconBlock(t)

	require.NotNil(t, block.Body)
	require.Equal(t, math.U64(10), block.GetTimestamp())
	require.Equal(t, version.Deneb1(), block.Version())
	require.False(t, block.IsNil())

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
}

func TestBeaconBlock_MarshalUnmarshalSSZ(t *testing.T) {
	t.Parallel()
	block := *generateValidBeaconBlock(t)

	sszBlock, err := block.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, sszBlock)

	var unmarshalledBlock types.BeaconBlock
	err = unmarshalledBlock.UnmarshalSSZ(sszBlock)
	require.NoError(t, err)

	unmarshalledBlock.Body.ExecutionPayload.EpVersion = block.Version()
	unmarshalledBlock.Body.SetVersion(block.Version())
	unmarshalledBlock.BbVersion = block.Version()
	require.Equal(t, block, unmarshalledBlock)
}

func TestBeaconBlock_HashTreeRoot(t *testing.T) {
	t.Parallel()
	block := generateValidBeaconBlock(t)
	hashRoot := block.HashTreeRoot()
	require.NotNil(t, hashRoot)
}

func TestBeaconBlock_IsNil(t *testing.T) {
	t.Parallel()
	var block *types.BeaconBlock
	require.True(t, block.IsNil())
}

func TestNewWithVersion(t *testing.T) {
	t.Parallel()
	slot := math.Slot(10)
	proposerIndex := math.ValidatorIndex(5)
	parentBlockRoot := common.Root{1, 2, 3, 4, 5}

	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
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
			require.Equal(t, v, block.Version())
		})
	}
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
	for _, v := range version.GetSupportedVersions() {
		t.Run(v.String(), func(t *testing.T) {
			t.Parallel()
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
}
