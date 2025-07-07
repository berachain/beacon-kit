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

package types

import (
	"testing"

	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

// TestSignedBeaconBlockSSZRoundTrip tests that SignedBeaconBlock can be marshaled and unmarshaled
// with the new SSZ implementation and produces expected results
func TestSignedBeaconBlockSSZRoundTrip(t *testing.T) {
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
			signature := crypto.BLSSignature{99, 98, 97} // Sample signature

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

			// Create SignedBeaconBlock
			signedBlock := &SignedBeaconBlock{
				BeaconBlock: block,
				Signature:   signature,
			}

			// Test 1: Marshal
			bytes, err := signedBlock.MarshalSSZ()
			require.NoError(t, err)
			require.NotEmpty(t, bytes)

			// Test 2: Size calculation
			expectedSize := signedBlock.SizeSSZ()
			require.Equal(t, expectedSize, len(bytes))

			// Test 3: Unmarshal
			newSignedBlock, err := NewEmptySignedBeaconBlockWithVersion(tc.forkVersion)
			require.NoError(t, err)
			err = newSignedBlock.UnmarshalSSZ(bytes)
			require.NoError(t, err)

			// Test 4: Verify fields
			require.Equal(t, slot, newSignedBlock.Slot)
			require.Equal(t, proposerIndex, newSignedBlock.ProposerIndex)
			require.Equal(t, parentRoot, newSignedBlock.ParentRoot)
			require.Equal(t, stateRoot, newSignedBlock.StateRoot)
			require.Equal(t, signature, newSignedBlock.Signature)
			require.Equal(t, block.Body.RandaoReveal, newSignedBlock.Body.RandaoReveal)
			require.Equal(t, block.Body.Graffiti, newSignedBlock.Body.Graffiti)

			// Test 5: HashTreeRoot
			htr1, err := signedBlock.HashTreeRoot()
			require.NoError(t, err)
			
			htr2, err := newSignedBlock.HashTreeRoot()
			require.NoError(t, err)
			
			require.Equal(t, htr1, htr2)
		})
	}
}

func TestSignedBeaconBlockSSZEdgeCases(t *testing.T) {
	forkVersion := version.Deneb()

	t.Run("EmptySignedBlock", func(t *testing.T) {
		// Create empty signed block
		signedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
		require.NoError(t, err)

		// Marshal
		bytes, err := signedBlock.MarshalSSZ()
		require.NoError(t, err)

		// Unmarshal
		newSignedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
		require.NoError(t, err)
		err = newSignedBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)

		// Verify empty values
		require.Equal(t, math.Slot(0), newSignedBlock.Slot)
		require.Equal(t, crypto.BLSSignature{}, newSignedBlock.Signature)
	})

	t.Run("MaxSignature", func(t *testing.T) {
		// Test with maximum signature values
		var maxSig crypto.BLSSignature
		for i := range maxSig {
			maxSig[i] = 0xFF
		}

		block := NewEmptyBeaconBlockWithVersion(forkVersion)
		signedBlock := &SignedBeaconBlock{
			BeaconBlock: block,
			Signature:   maxSig,
		}

		// Marshal and unmarshal
		bytes, err := signedBlock.MarshalSSZ()
		require.NoError(t, err)

		newSignedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
		require.NoError(t, err)
		err = newSignedBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)

		require.Equal(t, maxSig, newSignedBlock.Signature)
	})

	t.Run("InvalidUnmarshal", func(t *testing.T) {
		// Test unmarshaling invalid data
		signedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
		require.NoError(t, err)
		
		// Too short
		err = signedBlock.UnmarshalSSZ(make([]byte, 50))
		require.Error(t, err)
		
		// Empty
		err = signedBlock.UnmarshalSSZ([]byte{})
		require.Error(t, err)
		
		// Invalid offset (less than minimum required for offset + signature)
		invalidData := make([]byte, 100)
		// Set an invalid offset at position 0-4
		invalidData[0] = 0xFF
		invalidData[1] = 0xFF
		invalidData[2] = 0xFF
		invalidData[3] = 0xFF
		err = signedBlock.UnmarshalSSZ(invalidData)
		require.Error(t, err)
	})

	t.Run("NestedBlockData", func(t *testing.T) {
		// Test with more complex nested data
		slot := math.Slot(999999)
		proposerIndex := math.ValidatorIndex(12345)
		
		// Create block with nested data
		block, err := NewBeaconBlockWithVersion(slot, proposerIndex, common.Root{}, forkVersion)
		require.NoError(t, err)
		
		// Add blob commitments
		block.Body.BlobKzgCommitments = []eip4844.KZGCommitment{
			{1, 2, 3},
			{4, 5, 6},
		}
		
		// Add deposits
		block.Body.Deposits = []*Deposit{
			{
				Pubkey:      [48]byte{11, 12, 13},
				Credentials: [32]byte{14, 15, 16},
				Amount:      32000000000,
				Signature:   [96]byte{17, 18, 19},
				Index:       1,
			},
		}

		signedBlock := &SignedBeaconBlock{
			BeaconBlock: block,
			Signature:   crypto.BLSSignature{7, 8, 9},
		}

		// Marshal and unmarshal
		bytes, err := signedBlock.MarshalSSZ()
		require.NoError(t, err)

		newSignedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
		require.NoError(t, err)
		err = newSignedBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)

		// Verify nested data
		require.Equal(t, len(block.Body.BlobKzgCommitments), len(newSignedBlock.Body.BlobKzgCommitments))
		require.Equal(t, len(block.Body.Deposits), len(newSignedBlock.Body.Deposits))
		require.Equal(t, signedBlock.Signature, newSignedBlock.Signature)
	})
}

// TestSignedBeaconBlockSSZConsistency verifies that the SSZ implementation
// is consistent across marshal/unmarshal operations
func TestSignedBeaconBlockSSZConsistency(t *testing.T) {
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
	
	// Create signed block
	signedBlock := &SignedBeaconBlock{
		BeaconBlock: block,
		Signature:   crypto.BLSSignature{111, 112, 113},
	}

	// Multiple round trips to ensure consistency
	for i := 0; i < 5; i++ {
		// Marshal
		bytes, err := signedBlock.MarshalSSZ()
		require.NoError(t, err)
		
		// Create new signed block and unmarshal
		newSignedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
		require.NoError(t, err)
		err = newSignedBlock.UnmarshalSSZ(bytes)
		require.NoError(t, err)
		
		// Marshal again
		newBytes, err := newSignedBlock.MarshalSSZ()
		require.NoError(t, err)
		
		// Bytes should be identical
		require.Equal(t, bytes, newBytes, "Round trip %d produced different bytes", i)
		
		// HashTreeRoot should be identical
		htr1, err := signedBlock.HashTreeRoot()
		require.NoError(t, err)
		
		htr2, err := newSignedBlock.HashTreeRoot()
		require.NoError(t, err)
		
		require.Equal(t, htr1, htr2, "Round trip %d produced different HTR", i)
		
		// Use the new signed block for the next iteration
		signedBlock = newSignedBlock
	}
}

// TestSignedBeaconBlockElectraFeatures tests Electra-specific features
func TestSignedBeaconBlockElectraFeatures(t *testing.T) {
	forkVersion := version.Electra()

	// Create a block
	block, err := NewBeaconBlockWithVersion(
		math.Slot(12345),
		math.ValidatorIndex(42),
		common.Root{1, 2, 3},
		forkVersion,
	)
	require.NoError(t, err)

	// Set Electra-specific fields
	if version.EqualsOrIsAfter(forkVersion, version.Electra()) {
		execRequests := &ExecutionRequests{
			Deposits: []*DepositRequest{
				{
					Pubkey:      [48]byte{1, 2, 3},
					Credentials: [32]byte{4, 5, 6},
					Amount:      32000000000,
					Signature:   [96]byte{7, 8, 9},
					Index:       0,
				},
			},
			Withdrawals: []*WithdrawalRequest{
				{
					SourceAddress:   [20]byte{10, 11, 12},
					ValidatorPubKey: [48]byte{13, 14, 15},
					Amount:          1000000000,
				},
			},
		}
		err = block.Body.SetExecutionRequests(execRequests)
		require.NoError(t, err)
	}

	// Create signed block
	signedBlock := &SignedBeaconBlock{
		BeaconBlock: block,
		Signature:   crypto.BLSSignature{50, 51, 52},
	}

	// Marshal
	bytes, err := signedBlock.MarshalSSZ()
	require.NoError(t, err)

	// Unmarshal
	newSignedBlock, err := NewEmptySignedBeaconBlockWithVersion(forkVersion)
	require.NoError(t, err)
	err = newSignedBlock.UnmarshalSSZ(bytes)
	require.NoError(t, err)

	// Verify Electra fields
	if version.EqualsOrIsAfter(forkVersion, version.Electra()) {
		execRequests, err := newSignedBlock.Body.GetExecutionRequests()
		require.NoError(t, err)
		require.NotNil(t, execRequests)
		require.Len(t, execRequests.Deposits, 1)
		require.Len(t, execRequests.Withdrawals, 1)
	}
}