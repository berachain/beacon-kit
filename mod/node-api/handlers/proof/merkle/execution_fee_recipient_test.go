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

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/merkle"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/proof/merkle/mock"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/require"
)

// TestExecutionFeeRecipientProof tests the ProveExecutionFeeRecipientInBlock
// function and that the generated proof correctly verifies.
func TestExecutionFeeRecipientProof(t *testing.T) {
	var proof []common.Root

	testCases := []struct {
		name                  string
		slot                  math.Slot
		proposerIndex         math.ValidatorIndex
		parentBlockRoot       common.Root
		bodyRoot              common.Root
		executionFeeRecipient common.ExecutionAddress
		expectedProofFile     string
	}{
		{
			name:                  "Empty Fee Recipient",
			slot:                  4,
			proposerIndex:         0,
			parentBlockRoot:       common.Root{1, 2, 3},
			bodyRoot:              common.Root{3, 2, 1},
			executionFeeRecipient: common.ExecutionAddress{},
			expectedProofFile:     "empty_fee_recipient_proof.json",
		},
		{
			name:            "Non-empty Fee Recipient",
			slot:            5,
			proposerIndex:   95,
			parentBlockRoot: common.Root{1, 2, 3, 4, 5, 6},
			bodyRoot:        common.Root{3, 2, 1, 9, 8, 7},
			executionFeeRecipient: common.NewExecutionAddressFromHex(
				"0x20f33ce90a13a4b5e7697e3544c3083b8f8a51d4",
			),
			expectedProofFile: "non_empty_fee_recipient_proof.json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bs, err := mock.NewBeaconState(
				tc.slot, nil, 0, tc.executionFeeRecipient,
			)
			require.NoError(t, err)

			bbh := (&types.BeaconBlockHeader{}).New(
				tc.slot,
				tc.proposerIndex,
				tc.parentBlockRoot,
				bs.HashTreeRoot(),
				tc.bodyRoot,
			)

			proof, _, err = merkle.ProveExecutionFeeRecipientInBlock(bbh, bs)
			require.NoError(t, err)
			expectedProof := ReadProofFromFile(t, tc.expectedProofFile)
			require.Equal(t, expectedProof, proof)
		})
	}
}
