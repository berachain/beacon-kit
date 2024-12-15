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

package merkle_test

import (
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/node-api/handlers/proof/merkle/mock"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

// TestExecutionNumberProof tests the ProveExecutionNumberInBlock
// function and that the generated proof correctly verifies.
func TestExecutionNumberProof(t *testing.T) {
	var proof []common.Root

	testCases := []struct {
		name              string
		slot              math.Slot
		proposerIndex     math.ValidatorIndex
		parentBlockRoot   common.Root
		bodyRoot          common.Root
		executionNumber   math.U64
		expectedProof     []common.Root
		expectedProofFile string
	}{
		{
			name:              "Empty Execution Number",
			slot:              4,
			proposerIndex:     0,
			parentBlockRoot:   common.Root{1, 2, 3},
			bodyRoot:          common.Root{3, 2, 1},
			executionNumber:   0,
			expectedProofFile: "empty_execution_number_proof.json",
		},
		{
			name:              "Non-empty Execution Number",
			slot:              5,
			proposerIndex:     95,
			parentBlockRoot:   common.Root{1, 2, 3, 4, 5, 6},
			bodyRoot:          common.Root{3, 2, 1, 9, 8, 7},
			executionNumber:   69420,
			expectedProofFile: "non_empty_execution_number_proof.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := mock.NewBeaconState(
				tc.slot, nil, tc.executionNumber, common.ExecutionAddress{},
			)
			require.NoError(t, err)

			bbh := (&types.BeaconBlockHeader{}).New(
				tc.slot,
				tc.proposerIndex,
				tc.parentBlockRoot,
				bs.HashTreeRoot(),
				tc.bodyRoot,
			)

			proof, _, err = merkle.ProveExecutionNumberInBlock(bbh, bs)
			require.NoError(t, err)
			expectedProof := ReadProofFromFile(t, tc.expectedProofFile)
			require.Equal(t, expectedProof, proof)
		})
	}
}
