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
	engineprimitives "github.com/berachain/beacon-kit/mod/engine-primitives/pkg/engine-primitives"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/bytes"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/eip4844"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/version"
	"github.com/stretchr/testify/require"
)

// generateValidBeaconBlock generates a valid beacon block for the Deneb.
func generateValidBeaconBlock() *types.BeaconBlock {
	// Initialize your block here
	return &types.BeaconBlock{
		Slot:          10,
		ProposerIndex: 5,
		ParentRoot:    bytes.B32{1, 2, 3, 4, 5},
		StateRoot:     bytes.B32{5, 4, 3, 2, 1},
		Body: &types.BeaconBlockBody{
			ExecutionPayload: &types.ExecutionPayload{
				Number:    10,
				ExtraData: []byte("dummy extra data for testing"),
				Transactions: [][]byte{
					[]byte("tx1"),
					[]byte("tx2"),
					[]byte("tx3"),
				},
				Withdrawals: []*engineprimitives.Withdrawal{
					{Index: 0, Amount: 100},
					{Index: 1, Amount: 200},
				},
				BaseFeePerGas: math.NewU256(0),
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
		},
	}
}

func TestBeaconBlockForDeneb(t *testing.T) {
	block := &types.BeaconBlock{
		Slot:          10,
		ProposerIndex: 5,
		ParentRoot:    bytes.B32{1, 2, 3, 4, 5},
	}
	require.NotNil(t, block)
}

func TestBeaconBlockFromSSZ(t *testing.T) {
	originalBlock := generateValidBeaconBlock()

	sszBlock, err := originalBlock.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, sszBlock)

	wrappedBlock := &types.BeaconBlock{}
	wrappedBlock, err = wrappedBlock.NewFromSSZ(sszBlock, version.Deneb)
	require.NoError(t, err)
	require.NotNil(t, wrappedBlock)
	require.Equal(t, originalBlock, wrappedBlock)
}

func TestBeaconBlockFromSSZForkVersionNotSupported(t *testing.T) {
	wrappedBlock := &types.BeaconBlock{}
	_, err := wrappedBlock.NewFromSSZ([]byte{}, 1)
	require.ErrorIs(t, err, types.ErrForkVersionNotSupported)
}

func TestBeaconBlock(t *testing.T) {
	block := generateValidBeaconBlock()

	require.NotNil(t, block.Body)
	require.Equal(t, math.U64(10), block.GetExecutionNumber())
	require.Equal(t, version.Deneb, block.Version())
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
	block := *generateValidBeaconBlock()

	sszBlock, err := block.MarshalSSZ()
	require.NoError(t, err)
	require.NotNil(t, sszBlock)

	var unmarshalledBlock types.BeaconBlock
	err = unmarshalledBlock.UnmarshalSSZ(sszBlock)
	require.NoError(t, err)

	require.Equal(t, block, unmarshalledBlock)

	var buf []byte
	buf, err = block.MarshalSSZTo(buf)
	require.NoError(t, err)

	// The two byte slices should be equal
	require.Equal(t, sszBlock, buf)
}

func TestBeaconBlock_HashTreeRoot(t *testing.T) {
	block := generateValidBeaconBlock()
	hashRoot := block.HashTreeRoot()
	require.NotNil(t, hashRoot)
}

func TestBeaconBlockEmpty(t *testing.T) {
	block := &types.BeaconBlock{}
	emptyBlock := block.Empty()
	require.NotNil(t, emptyBlock)
	require.IsType(t, &types.BeaconBlock{}, emptyBlock)
}

func TestBeaconBlock_IsNil(t *testing.T) {
	var block *types.BeaconBlock
	require.True(t, block.IsNil())
}

func TestNewWithVersion(t *testing.T) {
	slot := math.Slot(10)
	proposerIndex := math.ValidatorIndex(5)
	parentBlockRoot := bytes.B32{1, 2, 3, 4, 5}

	block, err := (&types.BeaconBlock{}).NewWithVersion(
		slot, proposerIndex, parentBlockRoot, version.Deneb,
	)
	require.NoError(t, err)
	require.NotNil(t, block)

	// Check the block's fields
	require.NotNil(t, block)
	require.Equal(t, slot, block.GetSlot())
	require.Equal(t, proposerIndex, block.GetProposerIndex())
	require.Equal(t, parentBlockRoot, block.GetParentBlockRoot())
	require.Equal(t, version.Deneb, block.Version())
}

func TestNewWithVersionInvalidForkVersion(t *testing.T) {
	slot := math.Slot(10)
	proposerIndex := math.ValidatorIndex(5)
	parentBlockRoot := bytes.B32{1, 2, 3, 4, 5}

	_, err := (&types.BeaconBlock{}).NewWithVersion(
		slot,
		proposerIndex,
		parentBlockRoot,
		100,
	) // 100 is an invalid fork version
	require.ErrorIs(t, err, types.ErrForkVersionNotSupported)
}

func TestBeaconBlock_GetTree(t *testing.T) {
	block := generateValidBeaconBlock()
	tree, err := block.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)
}
