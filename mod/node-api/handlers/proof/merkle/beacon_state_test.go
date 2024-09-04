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

package merkle

import (
	"testing"

	"github.com/berachain/beacon-kit/mod/consensus-types/pkg/types"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"github.com/berachain/beacon-kit/mod/primitives/pkg/math"
	"github.com/stretchr/testify/assert"
)

// TestProveBeaconStateInBlock tests the ProveBeaconStateInBlock function and
// that the generated proof correctly verifies.
func TestProveBeaconStateInBlock(t *testing.T) {
	// Create an empty BeaconBlockHeader
	bbh := (&types.BeaconBlockHeader{}).Empty()

	// Set up test cases
	testCases := []struct {
		name            string
		slot            math.Slot
		proposerIndex   math.ValidatorIndex
		parentBlockRoot common.Root
		bodyRoot        common.Root
		stateRoot       common.Root
		expectedProof   []common.Root
		expectedError   error
	}{
		{
			name:      "Empty block with non-empty state root",
			stateRoot: common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9},
			expectedProof: []common.Root{
				{
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
				{
					0xf5, 0xa5, 0xfd, 0x42, 0xd1, 0x6a, 0x20, 0x30, 0x27, 0x98,
					0xef, 0x6e, 0xd3, 0x09, 0x97, 0x9b, 0x43, 0x00, 0x3d, 0x23,
					0x20, 0xd9, 0xf0, 0xe8, 0xea, 0x98, 0x31, 0xa9, 0x27, 0x59,
					0xfb, 0x4b,
				},
				{
					0xdb, 0x56, 0x11, 0x4e, 0x00, 0xfd, 0xd4, 0xc1, 0xf8, 0x5c,
					0x89, 0x2b, 0xf3, 0x5a, 0xc9, 0xa8, 0x92, 0x89, 0xaa, 0xec,
					0xb1, 0xeb, 0xd0, 0xa9, 0x6c, 0xde, 0x60, 0x6a, 0x74, 0x8b,
					0x5d, 0x71,
				},
			},
			expectedError: nil,
		},
		{
			name:            "Non-empty block with empty state root",
			slot:            4,
			proposerIndex:   5,
			parentBlockRoot: common.Root{1, 2, 3, 4, 5, 6, 7, 8, 9},
			bodyRoot:        common.Root{9, 8, 7, 6, 5, 4, 3, 2, 1},
			expectedProof: []common.Root{
				{
					0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
					0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
				},
				{
					0xbd, 0x50, 0x45, 0x6d, 0x5a, 0xd1, 0x75, 0xae, 0x99, 0xa1,
					0x61, 0x2a, 0x53, 0xca, 0x22, 0x91, 0x24, 0xb6, 0x5d, 0x3e,
					0xaa, 0xbd, 0x9f, 0xf9, 0xc7, 0xab, 0x97, 0x9a, 0x38, 0x5c,
					0xf6, 0xb3,
				},
				{
					0x21, 0x37, 0xa0, 0xff, 0x62, 0x2b, 0xd6, 0xe3, 0x7, 0x28,
					0xbc, 0x64, 0xe7, 0xde, 0xed, 0x6e, 0xa4, 0x18, 0x59, 0x6c,
					0x77, 0x21, 0x43, 0xe2, 0x32, 0xcf, 0xea, 0x9d, 0x88, 0xba,
					0xbd, 0x58,
				},
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bbh.SetSlot(tc.slot)
			bbh.SetProposerIndex(tc.proposerIndex)
			bbh.SetParentBlockRoot(tc.parentBlockRoot)
			bbh.SetBodyRoot(tc.bodyRoot)
			bbh.SetStateRoot(tc.stateRoot)
			proof, err := ProveBeaconStateInBlock(bbh, true)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedProof, proof)
		})
	}
}
