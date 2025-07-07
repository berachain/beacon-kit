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

package merkle_test

import (
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle/mock"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/berachain/beacon-kit/primitives/version"
	"github.com/stretchr/testify/require"
)

func TestValidatorBalanceProof(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		numValidators   int
		validatorIndex  math.U64
		fork            common.Version
		slot            math.Slot
		parentBlockRoot common.Root
		bodyRoot        common.Root
	}{
		{
			name:            "Single Validator Balance - Electra",
			numValidators:   1,
			validatorIndex:  0,
			fork:            version.Electra(),
			slot:            4,
			parentBlockRoot: common.Root{1, 2, 3},
			bodyRoot:        common.Root{3, 2, 1},
		},
		{
			name:            "Multiple Validators Balance - First Leaf - Electra",
			numValidators:   10,
			validatorIndex:  2, // Within first leaf (0-3)
			fork:            version.Electra(),
			slot:            5,
			parentBlockRoot: common.Root{1, 2, 3, 4},
			bodyRoot:        common.Root{4, 3, 2, 1},
		},
		{
			name:            "Multiple Validators Balance - Second Leaf - Electra",
			numValidators:   10,
			validatorIndex:  5, // In second leaf (4-7)
			fork:            version.Electra(),
			slot:            6,
			parentBlockRoot: common.Root{1, 2, 3, 4, 5},
			bodyRoot:        common.Root{5, 4, 3, 2, 1},
		},
		{
			name:            "Many Validators Balance - Electra",
			numValidators:   100,
			validatorIndex:  47, // Tests a validator in a later leaf
			fork:            version.Electra(),
			slot:            7,
			parentBlockRoot: common.Root{1, 2, 3, 4, 5, 6},
			bodyRoot:        common.Root{6, 5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create validators
			vals := make(ctypes.Validators, tt.numValidators)
			for i := 0; i < tt.numValidators; i++ {
				vals[i] = &ctypes.Validator{}
			}

			// Create beacon state with validators
			bs := mock.NewBeaconStateWith(
				tt.slot, vals, 0, common.ExecutionAddress{}, tt.fork,
			)

			// Set balances manually since the mock doesn't set them
			bs.Balances = make([]uint64, tt.numValidators)
			for i := 0; i < tt.numValidators; i++ {
				bs.Balances[i] = uint64(32000000000 + i*1000000000) // 32 ETH + i ETH
			}

			// Create beacon block header
			bbh := ctypes.NewBeaconBlockHeader(
				tt.slot,
				0, // proposer index
				tt.parentBlockRoot,
				bs.HashTreeRoot(),
				tt.bodyRoot,
			)

			// Generate the proof
			proof, beaconRoot, err := merkle.ProveBalanceInBlock(
				tt.validatorIndex,
				bbh,
				bs,
			)
			require.NoError(t, err)
			require.NotEmpty(t, proof)
			require.NotEqual(t, [32]byte{}, beaconRoot)

			// Verify the proof is valid (this is done internally in ProveBalanceInBlock)
			// but we can double-check the returned beacon root matches
			expectedRoot := bbh.HashTreeRoot()
			require.Equal(t, expectedRoot, beaconRoot)
		})
	}
}

func TestValidatorBalanceProofEdgeCases(t *testing.T) {
	t.Parallel()

	// Test with validators at leaf boundaries
	numValidators := 17 // This gives us 5 leaves (0-3, 4-7, 8-11, 12-15, 16)
	
	vals := make(ctypes.Validators, numValidators)
	for i := 0; i < numValidators; i++ {
		vals[i] = &ctypes.Validator{}
	}

	bs := mock.NewBeaconStateWith(
		10, vals, 0, common.ExecutionAddress{}, version.Electra(),
	)
	
	// Set balances
	bs.Balances = make([]uint64, numValidators)
	for i := 0; i < numValidators; i++ {
		bs.Balances[i] = uint64(32000000000)
	}

	bbh := ctypes.NewBeaconBlockHeader(
		10,
		0, // proposer index
		common.Root{1, 2, 3},
		bs.HashTreeRoot(),
		common.Root{3, 2, 1},
	)

	// Test validators at different positions within leaves
	testCases := []struct {
		name           string
		validatorIndex math.U64
		leafIndex      uint64
		positionInLeaf uint64
	}{
		{"First in first leaf", 0, 0, 0},
		{"Last in first leaf", 3, 0, 3},
		{"First in second leaf", 4, 1, 0},
		{"Last in second leaf", 7, 1, 3},
		{"First in last partial leaf", 16, 4, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proof, beaconRoot, err := merkle.ProveBalanceInBlock(
				tc.validatorIndex,
				bbh,
				bs,
			)
			require.NoError(t, err)
			require.NotEmpty(t, proof)
			require.NotEqual(t, [32]byte{}, beaconRoot)

			// Verify the leaf index calculation
			calculatedLeafIndex := tc.validatorIndex / 4
			require.Equal(t, tc.leafIndex, calculatedLeafIndex.Unwrap())

			// Verify position within leaf
			positionInLeaf := tc.validatorIndex % 4
			require.Equal(t, tc.positionInLeaf, positionInLeaf.Unwrap())
		})
	}
}