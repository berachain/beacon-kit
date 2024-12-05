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
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

// TestProveBeaconStateInBlock tests the ProveBeaconStateInBlock function and
// that the generated proof correctly verifies.
func TestProveBeaconStateInBlock(t *testing.T) {
	bbh := (&types.BeaconBlockHeader{}).Empty()

	testCases := []struct {
		name              string
		slot              math.Slot
		proposerIndex     math.ValidatorIndex
		parentBlockRoot   common.Root
		bodyRoot          common.Root
		stateRoot         common.Root
		expectedProofFile string
	}{
		{
			name:              "Empty block with non-empty state root",
			stateRoot:         common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedProofFile: "empty_state_proof.json",
		},
		{
			name:              "Non-empty block with empty state root",
			slot:              4,
			proposerIndex:     5,
			parentBlockRoot:   common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9},
			bodyRoot:          common.Root{9, 8, 7, 6, 5, 4, 3, 2, 1},
			expectedProofFile: "non_empty_state_proof.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bbh.SetSlot(tc.slot)
			bbh.SetProposerIndex(tc.proposerIndex)
			bbh.SetParentBlockRoot(tc.parentBlockRoot)
			bbh.SetBodyRoot(tc.bodyRoot)
			bbh.SetStateRoot(tc.stateRoot)

			proof, err := merkle.ProveBeaconStateInBlock(bbh, true)
			require.NoError(t, err)
			expectedProof := ReadProofFromFile(t, tc.expectedProofFile)
			require.Equal(t, expectedProof, proof)
		})
	}
}
