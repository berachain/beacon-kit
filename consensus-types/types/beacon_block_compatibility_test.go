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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

// NOTE: BeaconBlock was using karalabe/ssz in commit 787c4675581b3281fbaf45ca8d8c26ae6cd72934
// However, due to the complexity of extracting the entire type hierarchy (BeaconBlock depends on
// BeaconBlockBody which has many nested types), this test focuses on comprehensive round-trip
// and fork-specific testing rather than direct karalabe/ssz comparison.

package types

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestBeaconBlockSSZRoundTrip tests that BeaconBlock can be marshaled and unmarshaled
// with the new SSZ implementation and produces expected results
func TestBeaconBlockSSZRoundTrip(t *testing.T) {
	// Test with different fork versions
	testCases := []struct {
		name        string
		forkVersion common.Version
	}{
		{"Deneb", version.Deneb()},
		{"Deneb1", version.Deneb1()},
		{"Electra", version.Electra()},
		{"Electra1", version.Electra1()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test BeaconBlock with sample data
			slot := math.Slot(12345)
			proposerIndex := math.ValidatorIndex(42)
			parentRoot := common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32}
			stateRoot := common.Root{32, 31, 30, 29, 28, 27, 26, 25, 24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}

			// Create block
			block, err := NewBeaconBlockWithVersion(slot, proposerIndex, parentRoot, tc.forkVersion)
			require.NoError(t, err)
			block.StateRoot = stateRoot

			// Set some body fields
			block.Body.RandaoReveal = [96]byte{1, 2, 3}
			block.Body.Graffiti = [32]byte{4, 5, 6}
			block.Body.Eth1Data.DepositRoot = common.Root{7, 8, 9}
			block.Body.Eth1Data.DepositCount = 100
			block.Body.Eth1Data.BlockHash = common.ExecutionHash{10, 11, 12}

			// Test 1: Marshal
			bytes, err := block.MarshalSSZ()
			require.NoError(t, err)
			require.NotEmpty(t, bytes)

			// Test 2: Size calculation
			expectedSize := block.SizeSSZ()
			require.Equal(t, expectedSize, len(bytes))

			// Test 3: Unmarshal
			newBlock := NewEmptyBeaconBlockWithVersion(tc.forkVersion)
			err = newBlock.UnmarshalSSZ(bytes)
			require.NoError(t, err)

			// Test 4: Verify fields
			require.Equal(t, slot, newBlock.Slot)
			require.Equal(t, proposerIndex, newBlock.ProposerIndex)
			require.Equal(t, parentRoot, newBlock.ParentRoot)
			require.Equal(t, stateRoot, newBlock.StateRoot)
			require.Equal(t, block.Body.RandaoReveal, newBlock.Body.RandaoReveal)
			require.Equal(t, block.Body.Graffiti, newBlock.Body.Graffiti)
			require.Equal(t, block.Body.Eth1Data.DepositRoot, newBlock.Body.Eth1Data.DepositRoot)
			require.Equal(t, block.Body.Eth1Data.DepositCount, newBlock.Body.Eth1Data.DepositCount)
			require.Equal(t, block.Body.Eth1Data.BlockHash, newBlock.Body.Eth1Data.BlockHash)

			// Test 5: HashTreeRoot
			htr1, err := block.HashTreeRoot()
			require.NoError(t, err)
			
			htr2, err := newBlock.HashTreeRoot()
			require.NoError(t, err)
			
			require.Equal(t, htr1, htr2)
		})
	}
}

func TestBeaconBlockSSZEdgeCases(t *testing.T) {
	forkVersion := version.Deneb()

	t.Run("EmptyBlock", func(t *testing.T) {
		// Create empty block
		block := NewEmptyBeaconBlockWithVersion(forkVersion)

		// Marshal
		bytes, err := block.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal
		newBlock := NewEmptyBeaconBlockWithVersion(forkVersion)
		err = newBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)

		// Verify empty values
		require.Equal(t, math.Slot(0), newBlock.Slot)
		require.Equal(t, math.ValidatorIndex(0), newBlock.ProposerIndex)
		require.Equal(t, common.Root{}, newBlock.ParentRoot)
		require.Equal(t, common.Root{}, newBlock.StateRoot)
	})

	t.Run("MaxValues", func(t *testing.T) {
		// Test with maximum values
		slot := math.Slot(^uint64(0))
		proposerIndex := math.ValidatorIndex(^uint64(0))
		
		// Fill roots with max values
		var maxRoot common.Root
		for i := range maxRoot {
			maxRoot[i] = 0xFF
		}

		block, err := NewBeaconBlockWithVersion(slot, proposerIndex, maxRoot, forkVersion)
		require.NoError(t, err)
		block.StateRoot = maxRoot

		// Marshal and unmarshal
		bytes, err := block.MarshalSSZ()
		require.NoError(t, err)

		newBlock := NewEmptyBeaconBlockWithVersion(forkVersion)
		err = newBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)

		require.Equal(t, slot, newBlock.Slot)
		require.Equal(t, proposerIndex, newBlock.ProposerIndex)
		require.Equal(t, maxRoot, newBlock.ParentRoot)
		require.Equal(t, maxRoot, newBlock.StateRoot)
	})

	t.Run("InvalidUnmarshal", func(t *testing.T) {
		// Test unmarshaling invalid data
		block := NewEmptyBeaconBlockWithVersion(forkVersion)
		
		// Too short
		err := block.UnmarshalSSZ(make([]byte, 50))
		require.Error(t, err)
		
		// Empty
		err = block.UnmarshalSSZ([]byte{})
		require.Error(t, err)
		
		// Invalid offset
		invalidData := make([]byte, 84)
		// Set an invalid offset at position 80-84
		invalidData[80] = 0xFF
		invalidData[81] = 0xFF
		invalidData[82] = 0xFF
		invalidData[83] = 0xFF
		err = block.UnmarshalSSZ(invalidData)
		require.Error(t, err)
	})
}

// TestBeaconBlockSSZConsistency verifies that the SSZ implementation
// is consistent across marshal/unmarshal operations
func TestBeaconBlockSSZConsistency(t *testing.T) {
	forkVersion := version.Deneb()

	// Create a block with various data
	block, err := NewBeaconBlockWithVersion(
		math.Slot(54321),
		math.ValidatorIndex(999),
		common.Root{99, 98, 97, 96, 95},
		forkVersion,
	)
	require.NoError(t, err)
	
	// Set additional fields
	block.StateRoot = common.Root{10, 20, 30, 40, 50}
	block.Body.RandaoReveal = [96]byte{100, 101, 102}
	block.Body.Graffiti = [32]byte{200, 201, 202}
	
	// Set execution payload fields
	block.Body.ExecutionPayload.ParentHash = common.ExecutionHash{1, 1, 1}
	block.Body.ExecutionPayload.FeeRecipient = common.ExecutionAddress{2, 2, 2}
	block.Body.ExecutionPayload.StateRoot = common.Bytes32{3, 3, 3}
	block.Body.ExecutionPayload.ReceiptsRoot = common.Bytes32{4, 4, 4}
	block.Body.ExecutionPayload.LogsBloom = [256]byte{5, 5, 5}
	block.Body.ExecutionPayload.Random = common.Bytes32{6, 6, 6}
	block.Body.ExecutionPayload.Number = 123456
	block.Body.ExecutionPayload.GasLimit = 30000000
	block.Body.ExecutionPayload.GasUsed = 15000000
	block.Body.ExecutionPayload.Timestamp = 1234567890
	block.Body.ExecutionPayload.BlockHash = common.ExecutionHash{7, 7, 7}
	block.Body.ExecutionPayload.BaseFeePerGas = math.NewU256(1000000000)
	
	// Add some transactions
	block.Body.ExecutionPayload.Transactions = [][]byte{
		{0x01, 0x02, 0x03},
		{0x04, 0x05, 0x06, 0x07},
		{0x08, 0x09},
	}
	
	// Add deposits
	block.Body.Deposits = []*Deposit{
		{
			Pubkey:      [48]byte{11, 12, 13},
			Credentials: [32]byte{14, 15, 16},
			Amount:      32000000000, // 32 ETH in Gwei
			Signature:   [96]byte{17, 18, 19},
			Index:       1,
		},
	}

	// Multiple round trips to ensure consistency
	for i := 0; i < 5; i++ {
		// Marshal
		bytes, err := block.MarshalSSZ()
		require.NoError(t, err)
		
		// Create new block and unmarshal
		newBlock := NewEmptyBeaconBlockWithVersion(forkVersion)
		err = newBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)
		
		// Marshal again
		newBytes, err := newBlock.MarshalSSZ()
		require.NoError(t, err)
		
		// Bytes should be identical
		require.Equal(t, bytes, newBytes, "Round trip %d produced different bytes", i)
		
		// HashTreeRoot should be identical
		htr1, err := block.HashTreeRoot()
		require.NoError(t, err)
		
		htr2, err := newBlock.HashTreeRoot()
		require.NoError(t, err)
		
		require.Equal(t, htr1, htr2, "Round trip %d produced different HTR", i)
		
		// Use the new block for the next iteration
		block = newBlock
	}
}